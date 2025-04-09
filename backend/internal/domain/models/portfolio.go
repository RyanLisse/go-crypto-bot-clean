package models

// Portfolio represents a user's portfolio
type Portfolio struct {
	TotalValue float64  `json:"total_value"`
	Assets     []*Asset `json:"assets"`
}

// Asset represents a single asset in a portfolio
type Asset struct {
	Symbol    string  `json:"symbol"`
	Quantity  float64 `json:"quantity"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}
