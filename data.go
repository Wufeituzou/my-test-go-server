package main

import "gorm.io/gorm"

type User struct {
	gorm.Model
	GoogleOpenID string `json:"googleOpenId"`
	DisplayName  string `json:"displayName"`
	AvatarUrl    string `json:"avatarUrl"`
}

type AIQuota struct {
	ID            string `json"id"`
	UserID        string `json:"userId"`
	Quota         string `json:"quota"`
	LastGrantTime string `json:"lastGrantTime"`
}

type PurchaseProduct struct {
	PurchaseToken        string `json:"purchaseToken"`
	ProductID            string `json:"productId"`
	AcknowledgementState string `json:"acknowledgementState"`
}

type PurchaseSubscription struct {
	PurchaseToken     string `json:"purchaseToken"`
	ProductID         string `json:"productId"`
	SubscriptionState string `json:"subscriptionState"`
}

type PurchaseRecord struct {
	UserID         string   `json:"userId"`
	PurchaseTokens []string `json:"purchaseTokens"`
}
