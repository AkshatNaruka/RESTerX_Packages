package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	
	"RestCLI/pkg"
)

var (
	subscriptionService = pkg.NewSubscriptionService()
	paymentService      = pkg.NewPaymentService()
)

// GetPlansHandler returns all available subscription plans
func GetPlansHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}
	
	plans := subscriptionService.GetAllPlans()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"plans": plans,
	})
}

// GetSubscriptionHandler returns the user's current subscription
func GetSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}
	
	userID := getUserIDFromRequest(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	subscription, err := subscriptionService.GetOrCreateSubscription(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// CreateRazorpayOrderHandler creates a Razorpay order for subscription
func CreateRazorpayOrderHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}
	
	userID := getUserIDFromRequest(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	var req struct {
		PlanType string `json:"planType"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.PlanType != "pro" && req.PlanType != "enterprise" {
		http.Error(w, "Invalid plan type", http.StatusBadRequest)
		return
	}
	
	order, err := paymentService.CreateRazorpayOrder(userID, req.PlanType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create order: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order": order,
		"keyId": paymentService.GetRazorpayKeyID(),
	})
}

// VerifyRazorpayPaymentHandler verifies and processes Razorpay payment
func VerifyRazorpayPaymentHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}
	
	userID := getUserIDFromRequest(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	var req struct {
		PlanType          string `json:"planType"`
		RazorpayOrderID   string `json:"razorpay_order_id"`
		RazorpayPaymentID string `json:"razorpay_payment_id"`
		RazorpaySignature string `json:"razorpay_signature"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	verification := pkg.RazorpayPaymentVerification{
		OrderID:   req.RazorpayOrderID,
		PaymentID: req.RazorpayPaymentID,
		Signature: req.RazorpaySignature,
	}
	
	// Process payment
	payment, err := paymentService.HandleRazorpayPayment(userID, req.PlanType, verification)
	if err != nil {
		http.Error(w, fmt.Sprintf("Payment verification failed: %v", err), http.StatusBadRequest)
		return
	}
	
	// Upgrade subscription
	subscription, err := subscriptionService.UpgradeSubscription(userID, req.PlanType, "razorpay", req.RazorpayOrderID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upgrade subscription: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"payment":      payment,
		"subscription": subscription,
	})
}

// CreatePayPalOrderHandler creates a PayPal order for subscription
func CreatePayPalOrderHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}
	
	userID := getUserIDFromRequest(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	var req struct {
		PlanType string `json:"planType"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.PlanType != "pro" && req.PlanType != "enterprise" {
		http.Error(w, "Invalid plan type", http.StatusBadRequest)
		return
	}
	
	order, err := paymentService.CreatePayPalOrder(userID, req.PlanType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create PayPal order: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order": order,
	})
}

// CapturePayPalOrderHandler captures and processes PayPal payment
func CapturePayPalOrderHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}
	
	userID := getUserIDFromRequest(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	var req struct {
		PlanType string `json:"planType"`
		OrderID  string `json:"orderId"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Process payment
	payment, err := paymentService.HandlePayPalPayment(userID, req.PlanType, req.OrderID)
	if err != nil {
		http.Error(w, fmt.Sprintf("PayPal payment processing failed: %v", err), http.StatusBadRequest)
		return
	}
	
	// Upgrade subscription
	subscription, err := subscriptionService.UpgradeSubscription(userID, req.PlanType, "paypal", req.OrderID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upgrade subscription: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"payment":      payment,
		"subscription": subscription,
	})
}

// CancelSubscriptionHandler cancels a user's subscription
func CancelSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}
	
	userID := getUserIDFromRequest(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	if err := subscriptionService.CancelSubscription(userID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to cancel subscription: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Subscription cancelled successfully",
	})
}

// GetPaymentHistoryHandler returns user's payment history
func GetPaymentHistoryHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}
	
	userID := getUserIDFromRequest(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	var payments []pkg.Payment
	if err := pkg.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&payments).Error; err != nil {
		http.Error(w, "Failed to fetch payment history", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"payments": payments,
	})
}

// Helper function to get user ID from request
func getUserIDFromRequest(r *http.Request) uint {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0
	}
	
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0
	}
	
	return uint(userID)
}
