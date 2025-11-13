package pkg

import (
	"errors"
	"time"
	"gorm.io/gorm"
)

type SubscriptionService struct{}

func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{}
}

// CreateFreeSubscription creates a free subscription for a user
func (ss *SubscriptionService) CreateFreeSubscription(userID uint) (*Subscription, error) {
	freePlan := GetPlanFeatures("free")
	
	subscription := Subscription{
		UserID:            userID,
		PlanType:          "free",
		Status:            "active",
		StartDate:         time.Now(),
		EndDate:           time.Now().AddDate(100, 0, 0), // 100 years for free plan
		NextBillingDate:   time.Now().AddDate(100, 0, 0),
		Amount:            0,
		Currency:          "USD",
		AutoRenew:         false,
		APIRequestLimit:   freePlan.APIRequestLimit,
		APIRequestsUsed:   0,
		CollectionLimit:   freePlan.CollectionLimit,
		EnvironmentLimit:  freePlan.EnvironmentLimit,
		TeamMemberLimit:   freePlan.TeamMemberLimit,
	}
	
	if err := DB.Create(&subscription).Error; err != nil {
		return nil, err
	}
	
	return &subscription, nil
}

// GetActiveSubscription returns the active subscription for a user
func (ss *SubscriptionService) GetActiveSubscription(userID uint) (*Subscription, error) {
	var subscription Subscription
	err := DB.Where("user_id = ? AND status = ?", userID, "active").
		Order("created_at DESC").
		First(&subscription).Error
	
	if err != nil {
		return nil, err
	}
	
	return &subscription, nil
}

// GetOrCreateSubscription gets active subscription or creates free one
func (ss *SubscriptionService) GetOrCreateSubscription(userID uint) (*Subscription, error) {
	sub, err := ss.GetActiveSubscription(userID)
	if err != nil {
		// Create free subscription if none exists
		return ss.CreateFreeSubscription(userID)
	}
	return sub, nil
}

// UpgradeSubscription upgrades a user's subscription to a paid plan
func (ss *SubscriptionService) UpgradeSubscription(userID uint, planType string, paymentProvider string, providerSubID string) (*Subscription, error) {
	// Cancel existing active subscriptions
	DB.Model(&Subscription{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Update("status", "cancelled")
	
	// Get plan features
	plan := GetPlanFeatures(planType)
	
	// Create new subscription
	subscription := Subscription{
		UserID:            userID,
		PlanType:          planType,
		Status:            "active",
		StartDate:         time.Now(),
		EndDate:           time.Now().AddDate(1, 0, 0), // 1 year
		NextBillingDate:   time.Now().AddDate(1, 0, 0),
		Amount:            plan.Price,
		Currency:          plan.Currency,
		PaymentProvider:   paymentProvider,
		ProviderSubID:     providerSubID,
		AutoRenew:         true,
		APIRequestLimit:   plan.APIRequestLimit,
		APIRequestsUsed:   0,
		CollectionLimit:   plan.CollectionLimit,
		EnvironmentLimit:  plan.EnvironmentLimit,
		TeamMemberLimit:   plan.TeamMemberLimit,
	}
	
	if err := DB.Create(&subscription).Error; err != nil {
		return nil, err
	}
	
	return &subscription, nil
}

// CancelSubscription cancels a user's subscription
func (ss *SubscriptionService) CancelSubscription(userID uint) error {
	result := DB.Model(&Subscription{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Update("status", "cancelled")
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("no active subscription found")
	}
	
	// Create a free subscription
	_, err := ss.CreateFreeSubscription(userID)
	return err
}

// CheckUsageLimit checks if the user has exceeded their usage limit
func (ss *SubscriptionService) CheckUsageLimit(userID uint, usageType string) (bool, error) {
	subscription, err := ss.GetOrCreateSubscription(userID)
	if err != nil {
		return false, err
	}
	
	switch usageType {
	case "api_request":
		if subscription.APIRequestLimit == -1 {
			return true, nil // unlimited
		}
		return subscription.APIRequestsUsed < subscription.APIRequestLimit, nil
	case "collection":
		if subscription.CollectionLimit == -1 {
			return true, nil // unlimited
		}
		var count int64
		DB.Model(&DBCollection{}).Where("created_by = ?", userID).Count(&count)
		return int(count) < subscription.CollectionLimit, nil
	case "environment":
		if subscription.EnvironmentLimit == -1 {
			return true, nil // unlimited
		}
		var count int64
		DB.Model(&DBEnvironment{}).Where("created_by = ?", userID).Count(&count)
		return int(count) < subscription.EnvironmentLimit, nil
	}
	
	return true, nil
}

// IncrementUsage increments the usage counter for API requests
func (ss *SubscriptionService) IncrementUsage(userID uint) error {
	return DB.Model(&Subscription{}).
		Where("user_id = ? AND status = ?", userID, "active").
		UpdateColumn("api_requests_used", gorm.Expr("api_requests_used + ?", 1)).
		Error
}

// ResetMonthlyUsage resets the monthly usage counters (called by cron job)
func (ss *SubscriptionService) ResetMonthlyUsage() error {
	return DB.Model(&Subscription{}).
		Where("status = ?", "active").
		Update("api_requests_used", 0).
		Error
}

// GetAllPlans returns all available plans
func (ss *SubscriptionService) GetAllPlans() []PlanFeatures {
	return []PlanFeatures{
		GetPlanFeatures("free"),
		GetPlanFeatures("pro"),
		GetPlanFeatures("enterprise"),
	}
}
