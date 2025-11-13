package pkg

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type PaymentService struct {
	razorpayKeyID     string
	razorpayKeySecret string
	paypalClientID    string
	paypalSecret      string
	paypalMode        string // "sandbox" or "live"
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		razorpayKeyID:     os.Getenv("RAZORPAY_KEY_ID"),
		razorpayKeySecret: os.Getenv("RAZORPAY_KEY_SECRET"),
		paypalClientID:    os.Getenv("PAYPAL_CLIENT_ID"),
		paypalSecret:      os.Getenv("PAYPAL_SECRET"),
		paypalMode:        getEnvOrDefault("PAYPAL_MODE", "sandbox"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Razorpay Types
type RazorpayOrder struct {
	ID              string  `json:"id"`
	Entity          string  `json:"entity"`
	Amount          int     `json:"amount"` // Amount in paise
	AmountPaid      int     `json:"amount_paid"`
	AmountDue       int     `json:"amount_due"`
	Currency        string  `json:"currency"`
	Receipt         string  `json:"receipt"`
	Status          string  `json:"status"`
	Attempts        int     `json:"attempts"`
	Notes           map[string]string `json:"notes"`
	CreatedAt       int64   `json:"created_at"`
}

type RazorpayPaymentVerification struct {
	OrderID   string `json:"razorpay_order_id"`
	PaymentID string `json:"razorpay_payment_id"`
	Signature string `json:"razorpay_signature"`
}

// PayPal Types
type PayPalAccessToken struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppID       string `json:"app_id"`
	ExpiresIn   int    `json:"expires_in"`
	Nonce       string `json:"nonce"`
}

type PayPalOrder struct {
	ID            string                 `json:"id"`
	Status        string                 `json:"status"`
	Intent        string                 `json:"intent"`
	PurchaseUnits []PayPalPurchaseUnit   `json:"purchase_units"`
	Links         []PayPalLink           `json:"links"`
	CreateTime    string                 `json:"create_time"`
	UpdateTime    string                 `json:"update_time"`
}

type PayPalPurchaseUnit struct {
	ReferenceID string        `json:"reference_id"`
	Amount      PayPalAmount  `json:"amount"`
	Description string        `json:"description"`
}

type PayPalAmount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type PayPalLink struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

// ==================== RAZORPAY METHODS ====================

// CreateRazorpayOrder creates a Razorpay order
func (ps *PaymentService) CreateRazorpayOrder(userID uint, planType string) (*RazorpayOrder, error) {
	if ps.razorpayKeyID == "" || ps.razorpayKeySecret == "" {
		return nil, errors.New("razorpay credentials not configured")
	}
	
	plan := GetPlanFeatures(planType)
	amountInPaise := int(plan.Price * 100) // Convert to paise
	
	orderData := map[string]interface{}{
		"amount":   amountInPaise,
		"currency": "INR", // Razorpay primarily works with INR
		"receipt":  fmt.Sprintf("order_%d_%d", userID, time.Now().Unix()),
		"notes": map[string]string{
			"user_id":   fmt.Sprintf("%d", userID),
			"plan_type": planType,
		},
	}
	
	jsonData, err := json.Marshal(orderData)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("POST", "https://api.razorpay.com/v1/orders", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.SetBasicAuth(ps.razorpayKeyID, ps.razorpayKeySecret)
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("razorpay api error: %s", string(body))
	}
	
	var order RazorpayOrder
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, err
	}
	
	return &order, nil
}

// VerifyRazorpaySignature verifies the Razorpay payment signature
func (ps *PaymentService) VerifyRazorpaySignature(orderID, paymentID, signature string) bool {
	message := orderID + "|" + paymentID
	mac := hmac.New(sha256.New, []byte(ps.razorpayKeySecret))
	mac.Write([]byte(message))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// HandleRazorpayPayment processes a Razorpay payment
func (ps *PaymentService) HandleRazorpayPayment(userID uint, planType string, verification RazorpayPaymentVerification) (*Payment, error) {
	// Verify signature
	if !ps.VerifyRazorpaySignature(verification.OrderID, verification.PaymentID, verification.Signature) {
		return nil, errors.New("invalid payment signature")
	}
	
	plan := GetPlanFeatures(planType)
	
	// Create payment record
	payment := Payment{
		UserID:          userID,
		Amount:          plan.Price,
		Currency:        "INR",
		Provider:        "razorpay",
		ProviderTxnID:   verification.PaymentID,
		ProviderOrderID: verification.OrderID,
		Status:          "completed",
		PaymentMethod:   "razorpay",
		Description:     fmt.Sprintf("%s plan subscription", planType),
		PaidAt:          timePtr(time.Now()),
	}
	
	if err := DB.Create(&payment).Error; err != nil {
		return nil, err
	}
	
	return &payment, nil
}

// ==================== PAYPAL METHODS ====================

// getPayPalAccessToken gets an access token from PayPal
func (ps *PaymentService) getPayPalAccessToken() (string, error) {
	if ps.paypalClientID == "" || ps.paypalSecret == "" {
		return "", errors.New("paypal credentials not configured")
	}
	
	baseURL := "https://api-m.paypal.com"
	if ps.paypalMode == "sandbox" {
		baseURL = "https://api-m.sandbox.paypal.com"
	}
	
	req, err := http.NewRequest("POST", baseURL+"/v1/oauth2/token", bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}
	
	req.SetBasicAuth(ps.paypalClientID, ps.paypalSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("paypal auth error: %s", string(body))
	}
	
	var tokenResp PayPalAccessToken
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}
	
	return tokenResp.AccessToken, nil
}

// CreatePayPalOrder creates a PayPal order
func (ps *PaymentService) CreatePayPalOrder(userID uint, planType string) (*PayPalOrder, error) {
	accessToken, err := ps.getPayPalAccessToken()
	if err != nil {
		return nil, err
	}
	
	plan := GetPlanFeatures(planType)
	
	baseURL := "https://api-m.paypal.com"
	if ps.paypalMode == "sandbox" {
		baseURL = "https://api-m.sandbox.paypal.com"
	}
	
	orderData := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"reference_id": fmt.Sprintf("user_%d_plan_%s", userID, planType),
				"amount": map[string]string{
					"currency_code": "USD",
					"value":         fmt.Sprintf("%.2f", plan.Price),
				},
				"description": fmt.Sprintf("RESTerX %s Plan - Yearly Subscription", planType),
			},
		},
		"application_context": map[string]string{
			"return_url": os.Getenv("PAYPAL_RETURN_URL"),
			"cancel_url": os.Getenv("PAYPAL_CANCEL_URL"),
		},
	}
	
	jsonData, err := json.Marshal(orderData)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("POST", baseURL+"/v2/checkout/orders", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("paypal order creation error: %s", string(body))
	}
	
	var order PayPalOrder
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, err
	}
	
	return &order, nil
}

// CapturePayPalOrder captures a PayPal order after approval
func (ps *PaymentService) CapturePayPalOrder(orderID string) (*PayPalOrder, error) {
	accessToken, err := ps.getPayPalAccessToken()
	if err != nil {
		return nil, err
	}
	
	baseURL := "https://api-m.paypal.com"
	if ps.paypalMode == "sandbox" {
		baseURL = "https://api-m.sandbox.paypal.com"
	}
	
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/checkout/orders/%s/capture", baseURL, orderID), nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("paypal capture error: %s", string(body))
	}
	
	var order PayPalOrder
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, err
	}
	
	return &order, nil
}

// HandlePayPalPayment processes a PayPal payment
func (ps *PaymentService) HandlePayPalPayment(userID uint, planType string, orderID string) (*Payment, error) {
	// Capture the order
	order, err := ps.CapturePayPalOrder(orderID)
	if err != nil {
		return nil, err
	}
	
	if order.Status != "COMPLETED" {
		return nil, errors.New("payment not completed")
	}
	
	plan := GetPlanFeatures(planType)
	
	// Create payment record
	payment := Payment{
		UserID:          userID,
		Amount:          plan.Price,
		Currency:        "USD",
		Provider:        "paypal",
		ProviderTxnID:   order.ID,
		ProviderOrderID: orderID,
		Status:          "completed",
		PaymentMethod:   "paypal",
		Description:     fmt.Sprintf("%s plan subscription", planType),
		PaidAt:          timePtr(time.Now()),
	}
	
	if err := DB.Create(&payment).Error; err != nil {
		return nil, err
	}
	
	return &payment, nil
}

// Helper function to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}

// GetRazorpayKeyID returns the Razorpay key ID (public key)
func (ps *PaymentService) GetRazorpayKeyID() string {
	return ps.razorpayKeyID
}
