package api

import (
	"RestCLI/pkg"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// StorageCollectionsHandler handles collection storage operations
func StorageCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "GET":
		// Get all collections
		collections, err := pkg.StorageMgr.GetAllCollections()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(collections)

	case "POST":
		// Create new collection
		var collection pkg.StorageCollection
		if err := json.NewDecoder(r.Body).Decode(&collection); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := pkg.StorageMgr.SaveCollection(&collection); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(collection)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// StorageCollectionHandler handles single collection operations
func StorageCollectionHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	switch r.Method {
	case "GET":
		collection, err := pkg.StorageMgr.GetCollection(id)
		if err != nil {
			http.Error(w, "Collection not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(collection)

	case "PUT":
		var collection pkg.StorageCollection
		if err := json.NewDecoder(r.Body).Decode(&collection); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		collection.ID = id

		if err := pkg.StorageMgr.SaveCollection(&collection); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(collection)

	case "DELETE":
		if err := pkg.StorageMgr.DeleteCollection(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// StorageEnvironmentsHandler handles environment storage operations
func StorageEnvironmentsHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "GET":
		// Get all environments
		environments, err := pkg.StorageMgr.GetAllEnvironments()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(environments)

	case "POST":
		// Create new environment
		var environment pkg.StorageEnvironment
		if err := json.NewDecoder(r.Body).Decode(&environment); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := pkg.StorageMgr.SaveEnvironment(&environment); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(environment)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// StorageEnvironmentHandler handles single environment operations
func StorageEnvironmentHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	switch r.Method {
	case "GET":
		environment, err := pkg.StorageMgr.GetEnvironment(id)
		if err != nil {
			http.Error(w, "Environment not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(environment)

	case "PUT":
		var environment pkg.StorageEnvironment
		if err := json.NewDecoder(r.Body).Decode(&environment); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		environment.ID = id

		if err := pkg.StorageMgr.SaveEnvironment(&environment); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(environment)

	case "DELETE":
		if err := pkg.StorageMgr.DeleteEnvironment(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// StorageRequestsHandler handles request storage operations
func StorageRequestsHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "GET":
		// Get all requests or filter by collection
		collectionID := r.URL.Query().Get("collectionId")
		var requests []*pkg.StorageRequest
		var err error

		if collectionID != "" {
			requests, err = pkg.StorageMgr.GetRequestsByCollection(collectionID)
		} else {
			requests, err = pkg.StorageMgr.GetAllRequests()
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(requests)

	case "POST":
		// Create new request
		var request pkg.StorageRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := pkg.StorageMgr.SaveRequest(&request); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(request)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// StorageRequestHandler handles single request operations
func StorageRequestHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	switch r.Method {
	case "GET":
		request, err := pkg.StorageMgr.GetRequest(id)
		if err != nil {
			http.Error(w, "Request not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(request)

	case "PUT":
		var request pkg.StorageRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		request.ID = id

		if err := pkg.StorageMgr.SaveRequest(&request); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(request)

	case "DELETE":
		if err := pkg.StorageMgr.DeleteRequest(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// StorageHistoryHandler handles history storage operations
func StorageHistoryHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "GET":
		// Get history with optional room filter and limit
		roomID := r.URL.Query().Get("roomId")
		limit := 100 // default limit

		var histories []*pkg.StorageHistory
		var err error

		if roomID != "" {
			histories, err = pkg.StorageMgr.GetHistoryByRoom(roomID, limit)
		} else {
			histories, err = pkg.StorageMgr.GetAllHistory(limit)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(histories)

	case "POST":
		// Save history entry
		var history pkg.StorageHistory
		if err := json.NewDecoder(r.Body).Decode(&history); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := pkg.StorageMgr.SaveHistory(&history); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(history)

	case "DELETE":
		// Delete history entry by ID
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing history ID", http.StatusBadRequest)
			return
		}

		if err := pkg.StorageMgr.DeleteHistory(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HealthHandler checks backend and storage health
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	health := map[string]interface{}{
		"status":  "healthy",
		"storage": "unknown",
	}

	if pkg.StorageMgr != nil {
		if err := pkg.StorageMgr.Ping(); err != nil {
			health["storage"] = "unhealthy"
			health["error"] = err.Error()
		} else {
			health["storage"] = "healthy"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// ProxyRequestHandler proxies HTTP requests through the backend
// This provides CORS support and logging for frontend requests
func ProxyRequestHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request from frontend
	var proxyReq struct {
		Method  string            `json:"method"`
		URL     string            `json:"url"`
		Headers map[string]string `json:"headers"`
		Body    string            `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&proxyReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Execute the HTTP request using existing pkg functions
	apiReq := pkg.APIRequest{
		Method:  proxyReq.Method,
		URL:     proxyReq.URL,
		Headers: proxyReq.Headers,
		Body:    proxyReq.Body,
	}

	// Execute the HTTP request using unified MakeHTTPRequest
	response := pkg.MakeHTTPRequest(proxyReq.Method, apiReq.URL, apiReq.Body, apiReq.Headers)

	// Check if there was an error
	if response.Error != "" {
		// Return error as a response object
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Still 200, but indicate error in response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":        true,
			"message":      response.Error,
			"statusCode":   response.StatusCode,
			"responseTime": response.ResponseTime.Milliseconds(),
		})
		return
	}

	// Return successful response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
