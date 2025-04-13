package market

import (
	"time"
)

// OrderBookEntry represents a price level in the order book
type OrderBookEntry struct {
	// Price is the price level
	Price float64 `json:"price"`

	// Quantity is the quantity at this price level
	Quantity float64 `json:"quantity"`
}

// OrderBook represents the current state of the order book for a trading pair
type OrderBook struct {
	// Exchange is the name of the exchange (e.g., "binance")
	Exchange string `json:"exchange"`

	// Symbol is the trading pair (e.g., "BTCUSDT")
	Symbol string `json:"symbol"`

	// Bids are buy orders, sorted by price in descending order
	Bids []OrderBookEntry `json:"bids"`

	// Asks are sell orders, sorted by price in ascending order
	Asks []OrderBookEntry `json:"asks"`

	// UpdateTime is when this order book was last updated
	UpdateTime time.Time `json:"updateTime"`

	// UpdateID is a sequence ID used for synchronization
	UpdateID int64 `json:"updateID"`
}
