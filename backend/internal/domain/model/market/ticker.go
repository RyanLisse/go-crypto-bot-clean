package market

import (
	"time"
)

// Ticker represents market data for a symbol
type Ticker struct {
	ID            string    `json:"id"`
	Symbol        string    `json:"symbol"`
	Price         float64   `json:"price"`
	Volume        float64   `json:"volume"`
	High24h       float64   `json:"high24h"`
	Low24h        float64   `json:"low24h"`
	PriceChange   float64   `json:"priceChange"`
	PercentChange float64   `json:"percentChange"`
	LastUpdated   time.Time `json:"lastUpdated"`
	Exchange      string    `json:"exchange"`
}

// NewTicker creates a new ticker with the current time as LastUpdated
func NewTicker(symbol string, price float64) *Ticker {
	return &Ticker{
		Symbol:      symbol,
		Price:       price,
		LastUpdated: time.Now(),
	}
}
