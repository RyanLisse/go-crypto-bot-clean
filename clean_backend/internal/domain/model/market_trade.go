package model

import (
	"time"
)

// MarketTrade represents a trade that occurred on the market
type MarketTrade struct {
	ID            string    `json:"id"`
	Symbol        string    `json:"symbol"`
	Exchange      string    `json:"exchange"`
	Price         float64   `json:"price"`
	Quantity      float64   `json:"quantity"`
	QuoteQuantity float64   `json:"quoteQuantity"`
	Time          time.Time `json:"time"`
	IsBuyerMaker  bool      `json:"isBuyerMaker"`
	IsBestMatch   bool      `json:"isBestMatch"`
}
