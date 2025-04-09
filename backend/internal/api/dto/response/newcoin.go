package response

import "time"

// NewCoinResponse represents a newly detected coin
type NewCoinResponse struct {
	ID            uint       `json:"id"`
	Symbol        string     `json:"symbol"`
	Name          string     `json:"name,omitempty"`
	FoundAt       time.Time  `json:"found_at"`
	FirstOpenTime *time.Time `json:"first_open_time,omitempty"`
	QuoteVolume   float64    `json:"quote_volume"`
	IsProcessed   bool       `json:"is_processed"`
	IsUpcoming    bool       `json:"is_upcoming"`
}

// NewCoinsListResponse represents a list of newly detected coins
type NewCoinsListResponse struct {
	Coins     []NewCoinResponse `json:"coins"`
	Count     int               `json:"count"`
	Timestamp time.Time         `json:"timestamp"`
}

// ProcessNewCoinsResponse represents the result of processing new coins
type ProcessNewCoinsResponse struct {
	ProcessedCoins []NewCoinResponse `json:"processed_coins"`
	Count          int               `json:"count"`
	Timestamp      time.Time         `json:"timestamp"`
}
