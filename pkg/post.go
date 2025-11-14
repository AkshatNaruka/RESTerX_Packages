package pkg

import (
	"time"
)

// HandlePostRequest sends a POST request to the specified URL and prints the response body
func HandlePostRequest(url string) {
	// Example payload, modify as needed
	payload := `{"key": "value"}`

	req := APIRequest{
		Method:  "POST",
		URL:     url,
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    payload,
	}

	response := MakePostRequest(url, payload, req.Headers)
	FormatAndPrintResponse(req, response)
}

// HandlePostRequestAdvanced sends a POST request with custom headers and body
func HandlePostRequestAdvanced(url string, headers map[string]string, body string) APIResponse {
	start := time.Now()
	
	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return APIResponse{
			Error: err.Error(),
		}
	}
	
	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return APIResponse{
			Error:        err.Error(),
			ResponseTime: time.Since(start),
		}
	}
	defer resp.Body.Close()
	
	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{
			StatusCode:   resp.StatusCode,
			Status:       resp.Status,
			Headers:      convertHeaders(resp.Header),
			Error:        err.Error(),
			ResponseTime: time.Since(start),
		}
	}
	
	return APIResponse{
		StatusCode:   resp.StatusCode,
		Status:       resp.Status,
		Headers:      convertHeaders(resp.Header),
		Body:         string(respBody),
		ResponseTime: time.Since(start),
	}
}

// MakePostRequest sends a POST request and returns structured response data
func MakePostRequest(url, body string, headers map[string]string) APIResponse {
	start := time.Now()
	
	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return APIResponse{
			Error:        err.Error(),
			ResponseTime: time.Since(start),
		}
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return APIResponse{
			Error:        err.Error(),
			ResponseTime: time.Since(start),
		}
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{
			StatusCode:   resp.StatusCode,
			Status:       resp.Status,
			Headers:      convertHeaders(resp.Header),
			Error:        "Error reading response body: " + err.Error(),
			ResponseTime: time.Since(start),
		}
	}

	return APIResponse{
		StatusCode:   resp.StatusCode,
		Status:       resp.Status,
		Headers:      convertHeaders(resp.Header),
		Body:         string(responseBody),
		ResponseTime: time.Since(start),
	}
}
