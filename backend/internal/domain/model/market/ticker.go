package market

import (
	"time"
)

// Ticker represents market data for a symbol
type Ticker struct {
	ID            string    
	Symbol        string    
	Price         float64   
	Volume        float64   
	High24h       float64   
	Low24h        float64   
	PriceChange   float64   
	PercentChange float64   
	LastUpdated   time.Time 
	Exchange      string    
}

// NewTicker creates a new ticker with the current time as LastUpdated
func NewTicker(symbol string, price float64) *Ticker {
	return &Ticker{
		Symbol:      symbol,
		Price:       price,
		LastUpdated: time.Now(),
	}
}
