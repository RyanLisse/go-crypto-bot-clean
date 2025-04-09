package models

import "time"

// NewCoinEvent represents an event for a newly listed coin.
type NewCoinEvent struct {
	Symbol    string    `json:"symbol"`
	Timestamp time.Time `json:"timestamp"`
}

// PurchaseOptions defines the options for purchasing a coin
type PurchaseOptions struct {
	StopLossPercent float64
	OrderType       string
	Price           float64
}
