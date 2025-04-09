package models

import "time"

// PurchaseDecision represents a decision to purchase a cryptocurrency
type PurchaseDecision struct {
	Symbol     string    `json:"symbol"`
	Decision   bool      `json:"decision"`
	Reason     string    `json:"reason"`
	Strategy   string    `json:"strategy"`
	Confidence float64   `json:"confidence"`
	Timestamp  time.Time `json:"timestamp"`
}
