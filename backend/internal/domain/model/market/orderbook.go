package market

import "time"

// OrderBookEntry represents a single price level in the order book
type OrderBookEntry struct {
	Price    float64 
	Quantity float64 
}

// OrderBook represents the market depth for a trading pair
type OrderBook struct {
	Symbol       string           `json:"symbol"`
	LastUpdated  time.Time        `json:"last_updated"`
	Bids         []OrderBookEntry `json:"bids"`
	Asks         []OrderBookEntry `json:"asks"`
	Exchange     string           `json:"exchange"`
	SequenceNum  int64            `json:"sequence_num,omitempty"` // For consistency checking
	LastUpdateID int64            `json:"last_update_id,omitempty"`
}
