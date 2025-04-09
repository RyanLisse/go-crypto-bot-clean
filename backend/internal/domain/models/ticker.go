package models

import (
	"time"
)

// Ticker represents real-time price information for a trading pair
type Ticker struct {
	Symbol         string    `json:"symbol"`
	Price          float64   `json:"price"`
	PriceChange    float64   `json:"priceChange"`
	PriceChangePct float64   `json:"priceChangePercent"`
	Volume         float64   `json:"volume"`
	QuoteVolume    float64   `json:"quoteVolume"`
	High24h        float64   `json:"high24h"`
	Low24h         float64   `json:"low24h"`
	Timestamp      time.Time `json:"timestamp"`
}
