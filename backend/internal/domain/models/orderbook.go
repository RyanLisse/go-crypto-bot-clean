package models

import "time"

// OrderBookEntry represents a single entry in the order book (bid or ask)
type OrderBookEntry struct {
	Price    float64 // Price level
	Quantity float64 // Quantity at this price level
}

// OrderBookUpdate represents an update to the order book
type OrderBookUpdate struct {
	Symbol        string           // Trading pair symbol
	LastUpdateID  int64            // Last update ID
	FirstUpdateID int64            // First update ID
	Bids          []OrderBookEntry // Bid price/quantity pairs
	Asks          []OrderBookEntry // Ask price/quantity pairs
	Timestamp     time.Time        // Update time
}
