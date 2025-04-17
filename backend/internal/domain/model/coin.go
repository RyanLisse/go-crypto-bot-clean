package model

import "time"

// Status represents the status of a coin
type Status string

const (
	// StatusPending indicates a coin is pending listing
	StatusPending Status = "pending"
	// StatusTrading indicates a coin is available for trading
	StatusTrading Status = "trading"
	// StatusDelisted indicates a coin has been delisted
	StatusDelisted Status = "delisted"
)

// Coin represents a cryptocurrency coin
type Coin struct {
	ID          string    `json:"id"`
	Symbol      string    `json:"symbol"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	ListedAt    time.Time `json:"listed_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CoinEvent represents an event related to a new coin
type CoinEvent struct {
	CoinID     string                 `json:"coin_id"`
	EventType  string                 `json:"event_type"` // "new_listing", "status_change", etc.
	OldStatus  CoinStatus             `json:"old_status,omitempty"`
	NewStatus  CoinStatus             `json:"new_status,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Exchange   string                 `json:"exchange"`
	Additional map[string]interface{} `json:"additional,omitempty"`
}
