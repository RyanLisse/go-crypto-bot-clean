package models

import "time"

// TradeSide represents the side of a trade (buy/sell)
type TradeSide string

const (
	TradeSideBuy  TradeSide = "BUY"
	TradeSideSell TradeSide = "SELL"
)

// Trade represents a trade execution on the exchange
type Trade struct {
	ID        string    // Unique trade ID
	Symbol    string    // Trading pair symbol
	Price     float64   // Trade execution price
	Quantity  float64   // Trade quantity
	Side      TradeSide // Trade side (buy/sell)
	Timestamp time.Time // Trade execution time
}
