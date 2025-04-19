package model

import (
	"time"
)

// TradeRecord represents a record of an executed trade
type TradeRecord struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Symbol        string    `json:"symbol"`
	Side          OrderSide `json:"side"`
	Type          OrderType `json:"type"`
	Quantity      float64   `json:"quantity"`
	Price         float64   `json:"price"`
	Amount        float64   `json:"amount"`
	Fee           float64   `json:"fee"`
	FeeCurrency   string    `json:"fee_currency"`
	OrderID       string    `json:"order_id"`
	TradeID       string    `json:"trade_id"`
	ExecutionTime time.Time `json:"execution_time"`
	Strategy      string    `json:"strategy"`
	Notes         string    `json:"notes"`
	Tags          []string  `json:"tags"`
	Metadata      map[string]interface{} `json:"metadata"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// DetectionLog represents a log of a market event detection
type DetectionLog struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Symbol      string    `json:"symbol"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	Description string    `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
	DetectedAt  time.Time `json:"detected_at"`
	ProcessedAt *time.Time `json:"processed_at"`
	Processed   bool      `json:"processed"`
	Result      string    `json:"result"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
