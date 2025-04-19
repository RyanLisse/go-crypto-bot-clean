package model

import (
	"time"
)

// OrderBook represents an order book for a symbol
type OrderBook struct {
	Symbol       string           `json:"symbol"`
	Exchange     string           `json:"exchange"`
	LastUpdateID int64            `json:"lastUpdateId"`
	SequenceNum  int64            `json:"sequenceNum"`
	Bids         []OrderBookEntry `json:"bids"`
	Asks         []OrderBookEntry `json:"asks"`
	Timestamp    time.Time        `json:"timestamp"`
}

// OrderBookEntry represents a single entry in an order book
type OrderBookEntry struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}
