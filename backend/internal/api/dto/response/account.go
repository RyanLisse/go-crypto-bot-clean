package response

import "time"

// AccountResponse represents the response for the account endpoint
type AccountResponse struct {
	TotalBalance    float64         `json:"total_balance"`
	AvailableFunds  float64         `json:"available_funds"`
	CurrentExposure float64         `json:"current_exposure"`
	Assets          []AssetResponse `json:"assets"`
	Timestamp       time.Time       `json:"timestamp"`
}

// AssetResponse represents a single asset in the account
type AssetResponse struct {
	Symbol    string  `json:"symbol"`
	Free      float64 `json:"free"`
	Locked    float64 `json:"locked"`
	Total     float64 `json:"total"`
	Price     float64 `json:"price"`
	ValueUSDT float64 `json:"value_usdt"`
}

// APIKeyValidationResponse represents the response for the API key validation endpoint
type APIKeyValidationResponse struct {
	Valid     bool      `json:"valid"`
	Timestamp time.Time `json:"timestamp"`
}

// ListenKeyResponse represents the response for the listen key endpoint
type ListenKeyResponse struct {
	ListenKey string    `json:"listen_key"`
	Expires   time.Time `json:"expires"`
	Timestamp time.Time `json:"timestamp"`
}
