package pkg

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

// Database models for persistence
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	FullName  string    `json:"fullName"`
	Avatar    string    `json:"avatar"`
	Role      string    `json:"role" gorm:"default:user"`
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	
	// Relationships - Use proper GORM associations
	Workspaces    []UserWorkspace   `json:"workspaces" gorm:"foreignKey:UserID"`
	Requests      []RequestHistory  `json:"requests" gorm:"foreignKey:UserID"`
	Subscriptions []Subscription    `json:"subscriptions" gorm:"foreignKey:UserID"`
	Payments      []Payment         `json:"payments" gorm:"foreignKey:UserID"`
}

type WorkspaceDB struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Type        string    `json:"type" gorm:"default:team"` // personal, team, enterprise
	IsActive    bool      `json:"isActive" gorm:"default:true"`
	CreatedBy   uint      `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	
	// Relationships - Use proper GORM associations
	Members     []UserWorkspace   `json:"members" gorm:"foreignKey:WorkspaceID"`
	Collections []DBCollection    `json:"collections" gorm:"foreignKey:WorkspaceID"`
	Environments []DBEnvironment  `json:"environments" gorm:"foreignKey:WorkspaceID"`
}

type UserWorkspace struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"userId" gorm:"not null"`
	WorkspaceID uint      `json:"workspaceId" gorm:"not null"`
	Role        string    `json:"role" gorm:"default:member"` // admin, editor, viewer, member
	JoinedAt    time.Time `json:"joinedAt"`
	
	// Relationships - Use proper GORM associations
	User      User        `json:"user" gorm:"foreignKey:UserID"`
	Workspace WorkspaceDB `json:"workspace" gorm:"foreignKey:WorkspaceID"`
}

type DBCollection struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	WorkspaceID uint      `json:"workspaceId"`
	RoomID      uint      `json:"roomId"` // Support for room-based collections
	CreatedBy   uint      `json:"createdBy"`
	Version     int       `json:"version" gorm:"default:1"`
	IsPublic    bool      `json:"isPublic" gorm:"default:false"`
	Tags        string    `json:"tags"` // JSON string for tags
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	
	// Relationships - Use proper GORM associations
	Workspace WorkspaceDB `json:"workspace" gorm:"foreignKey:WorkspaceID"`
	Room      Room        `json:"room" gorm:"foreignKey:RoomID"`
	Requests  []DBRequest `json:"requests" gorm:"foreignKey:CollectionID"`
}

type DBRequest struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	CollectionID uint      `json:"collectionId"`
	Name         string    `json:"name" gorm:"not null"`
	Method       string    `json:"method" gorm:"not null"`
	URL          string    `json:"url" gorm:"not null"`
	Headers      string    `json:"headers"` // JSON string
	Body         string    `json:"body"`
	AuthType     string    `json:"authType"`
	AuthData     string    `json:"authData"` // JSON string
	Tests        string    `json:"tests"` // JSON string
	PreScript    string    `json:"preScript"`
	PostScript   string    `json:"postScript"`
	Variables    string    `json:"variables"` // JSON string
	Order        int       `json:"order"`
	CreatedBy    uint      `json:"createdBy"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	
	// Relationships - Use proper GORM associations
	Collection DBCollection `json:"collection" gorm:"foreignKey:CollectionID"`
}

type DBEnvironment struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	WorkspaceID uint      `json:"workspaceId"`
	RoomID      uint      `json:"roomId"` // Support for room-based environments
	Name        string    `json:"name" gorm:"not null"`
	Variables   string    `json:"variables"` // JSON string
	IsActive    bool      `json:"isActive" gorm:"default:false"`
	CreatedBy   uint      `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	
	// Relationships - Use proper GORM associations
	Workspace WorkspaceDB `json:"workspace" gorm:"foreignKey:WorkspaceID"`
	Room      Room        `json:"room" gorm:"foreignKey:RoomID"`
}

type RequestHistory struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"userId"`
	WorkspaceID  uint      `json:"workspaceId"`
	Method       string    `json:"method" gorm:"not null"`
	URL          string    `json:"url" gorm:"not null"`
	Headers      string    `json:"headers"` // JSON string
	Body         string    `json:"body"`
	StatusCode   int       `json:"statusCode"`
	ResponseTime int64     `json:"responseTime"` // milliseconds
	ResponseSize int64     `json:"responseSize"` // bytes
	Success      bool      `json:"success"`
	CreatedAt    time.Time `json:"createdAt"`
	
	// Relationships - Use proper GORM associations
	User      User        `json:"user" gorm:"foreignKey:UserID"`
	Workspace WorkspaceDB `json:"workspace" gorm:"foreignKey:WorkspaceID"`
}

type APIMonitor struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	WorkspaceID  uint      `json:"workspaceId"`
	Name         string    `json:"name" gorm:"not null"`
	URL          string    `json:"url" gorm:"not null"`
	Method       string    `json:"method" gorm:"default:GET"`
	Headers      string    `json:"headers"` // JSON string
	Interval     int       `json:"interval" gorm:"default:300"` // seconds
	Timeout      int       `json:"timeout" gorm:"default:30"` // seconds
	IsActive     bool      `json:"isActive" gorm:"default:true"`
	AlertEmail   string    `json:"alertEmail"`
	CreatedBy    uint      `json:"createdBy"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	
	// Relationships - Use proper GORM associations
	Workspace WorkspaceDB    `json:"workspace" gorm:"foreignKey:WorkspaceID"`
	Checks    []MonitorCheck `json:"checks" gorm:"foreignKey:MonitorID"`
}

type MonitorCheck struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	MonitorID   uint      `json:"monitorId"`
	StatusCode  int       `json:"statusCode"`
	ResponseTime int64     `json:"responseTime"` // milliseconds
	Success     bool      `json:"success"`
	Error       string    `json:"error"`
	CheckedAt   time.Time `json:"checkedAt"`
	
	// Relationships - Use proper GORM associations
	Monitor APIMonitor `json:"monitor" gorm:"foreignKey:MonitorID"`
}

type Room struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	RoomID      string    `json:"roomId" gorm:"uniqueIndex;not null"` // User-friendly room identifier
	Name        string    `json:"name" gorm:"not null"`
	Password    string    `json:"-" gorm:"not null"` // Hashed password
	Description string    `json:"description"`
	RoomType    string    `json:"roomType" gorm:"default:free"` // free, premium for SaaS
	DeviceLimit int       `json:"deviceLimit" gorm:"default:3"` // Number of devices allowed
	IsActive    bool      `json:"isActive" gorm:"default:true"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	
	// Relationships
	Devices      []RoomDevice     `json:"devices" gorm:"foreignKey:RoomID"`
	Collections  []DBCollection   `json:"collections" gorm:"foreignKey:RoomID"`
	Environments []DBEnvironment  `json:"environments" gorm:"foreignKey:RoomID"`
}

type RoomDevice struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	RoomID     uint      `json:"roomId" gorm:"not null"`
	DeviceName string    `json:"deviceName" gorm:"not null"`
	DeviceID   string    `json:"deviceId" gorm:"not null"` // Unique device identifier
	LastActive time.Time `json:"lastActive"`
	CreatedAt  time.Time `json:"createdAt"`
	
	// Relationships
	Room Room `json:"room" gorm:"foreignKey:RoomID"`
}

// Subscription models for SaaS
type Subscription struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	UserID            uint      `json:"userId" gorm:"not null;index"`
	PlanType          string    `json:"planType" gorm:"not null"` // free, pro, enterprise
	Status            string    `json:"status" gorm:"default:active"` // active, cancelled, expired, past_due
	StartDate         time.Time `json:"startDate"`
	EndDate           time.Time `json:"endDate"`
	NextBillingDate   time.Time `json:"nextBillingDate"`
	Amount            float64   `json:"amount"`
	Currency          string    `json:"currency" gorm:"default:USD"`
	PaymentProvider   string    `json:"paymentProvider"` // razorpay, paypal
	ProviderSubID     string    `json:"providerSubId"` // Subscription ID from payment provider
	AutoRenew         bool      `json:"autoRenew" gorm:"default:true"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	
	// Usage limits based on plan
	APIRequestLimit   int       `json:"apiRequestLimit" gorm:"default:100"` // per month
	APIRequestsUsed   int       `json:"apiRequestsUsed" gorm:"default:0"`
	CollectionLimit   int       `json:"collectionLimit" gorm:"default:5"`
	EnvironmentLimit  int       `json:"environmentLimit" gorm:"default:2"`
	TeamMemberLimit   int       `json:"teamMemberLimit" gorm:"default:1"`
	
	// Relationships
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Payments []Payment `json:"payments" gorm:"foreignKey:SubscriptionID"`
}

type Payment struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"userId" gorm:"not null;index"`
	SubscriptionID  uint      `json:"subscriptionId" gorm:"index"`
	Amount          float64   `json:"amount" gorm:"not null"`
	Currency        string    `json:"currency" gorm:"default:USD"`
	Provider        string    `json:"provider" gorm:"not null"` // razorpay, paypal
	ProviderTxnID   string    `json:"providerTxnId" gorm:"uniqueIndex"` // Transaction ID from provider
	ProviderOrderID string    `json:"providerOrderId"` // Order ID from provider
	Status          string    `json:"status" gorm:"default:pending"` // pending, completed, failed, refunded
	PaymentMethod   string    `json:"paymentMethod"` // card, upi, netbanking, paypal, etc.
	Description     string    `json:"description"`
	InvoiceURL      string    `json:"invoiceUrl"`
	ReceiptURL      string    `json:"receiptUrl"`
	FailureReason   string    `json:"failureReason"`
	PaidAt          *time.Time `json:"paidAt"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	
	// Relationships
	User         User         `json:"user" gorm:"foreignKey:UserID"`
	Subscription Subscription `json:"subscription" gorm:"foreignKey:SubscriptionID"`
}

// PlanFeatures defines features for each plan type
type PlanFeatures struct {
	PlanType         string  `json:"planType"`
	Price            float64 `json:"price"`
	Currency         string  `json:"currency"`
	Interval         string  `json:"interval"` // yearly
	APIRequestLimit  int     `json:"apiRequestLimit"`
	CollectionLimit  int     `json:"collectionLimit"`
	EnvironmentLimit int     `json:"environmentLimit"`
	TeamMemberLimit  int     `json:"teamMemberLimit"`
	Features         []string `json:"features"`
}

// GetPlanFeatures returns the features for each plan
func GetPlanFeatures(planType string) PlanFeatures {
	plans := map[string]PlanFeatures{
		"free": {
			PlanType:         "free",
			Price:            0,
			Currency:         "USD",
			Interval:         "monthly",
			APIRequestLimit:  500,
			CollectionLimit:  -1, // unlimited
			EnvironmentLimit: -1, // unlimited
			TeamMemberLimit:  1,
			Features: []string{
				"500 API requests/month",
				"2 Rooms",
				"Basic collaboration (2 devices)",
				"Unlimited collections",
				"Unlimited environments",
				"Local storage + limited sync",
				"Community support",
			},
		},
		"pro": {
			PlanType:         "pro",
			Price:            49.99,
			Currency:         "USD",
			Interval:         "yearly",
			APIRequestLimit:  10000,
			CollectionLimit:  -1, // unlimited
			EnvironmentLimit: -1, // unlimited
			TeamMemberLimit:  5,
			Features: []string{
				"10,000 API requests/month",
				"20 Rooms",
				"Full room sync + device management (10 devices)",
				"Unlimited collections",
				"Unlimited environments",
				"Priority API proxy bandwidth",
				"Basic usage analytics",
				"Team collaboration (up to 5 members)",
				"Priority email support",
			},
		},
		"enterprise": {
			PlanType:         "enterprise",
			Price:            199.99,
			Currency:         "USD",
			Interval:         "yearly",
			APIRequestLimit:  -1, // unlimited
			CollectionLimit:  -1, // unlimited
			EnvironmentLimit: -1, // unlimited
			TeamMemberLimit:  -1, // unlimited
			Features: []string{
				"Unlimited API requests",
				"Unlimited Rooms & Devices",
				"Multi-user Teams + Role management",
				"Unlimited collections",
				"Unlimited environments",
				"Admin dashboard",
				"Private backend sync server",
				"Full analytics suite",
				"24/7 priority support",
				"Custom domains and SLA support",
				"Dedicated account manager",
			},
		},
	}
	
	if plan, exists := plans[planType]; exists {
		return plan
	}
	return plans["free"] // default to free
}

var DB *gorm.DB

// InitDatabase initializes the database connection and runs migrations
func InitDatabase(dsn string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Run auto-migration
	err = DB.AutoMigrate(
		&User{},
		&WorkspaceDB{},
		&UserWorkspace{},
		&DBCollection{},
		&DBRequest{},
		&DBEnvironment{},
		&RequestHistory{},
		&APIMonitor{},
		&MonitorCheck{},
		&Room{},
		&RoomDevice{},
		&Subscription{},
		&Payment{},
	)
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}