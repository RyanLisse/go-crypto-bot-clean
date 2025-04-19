package model

import (
	"time"
)

// MarketTicker represents a simplified ticker for legacy compatibility
type MarketTicker struct {
	// Symbol is the trading pair identifier (e.g., "BTCUSDT")
	Symbol string `json:"symbol"`

	// Exchange indicates which exchange this ticker is from
	Exchange string `json:"exchange"`

	// Price is the current price
	Price float64 `json:"price"`

	// Volume is the trading volume in the last 24 hours
	Volume float64 `json:"volume"`

	// High24h is the highest price in the last 24 hours
	High24h float64 `json:"high24h"`

	// Low24h is the lowest price in the last 24 hours
	Low24h float64 `json:"low24h"`

	// PriceChange is the absolute price change in the last 24 hours
	PriceChange float64 `json:"priceChange"`

	// PercentChange is the percentage price change in the last 24 hours
	PercentChange float64 `json:"percentChange"`

	// LastUpdated is when this ticker was last updated
	LastUpdated time.Time `json:"lastUpdated"`
}
