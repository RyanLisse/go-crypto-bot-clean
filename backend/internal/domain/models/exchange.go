package models

import (
	"time"
)

// Ticker is defined in ticker.go

// Kline represents a candlestick for a specific time interval
type Kline struct {
	Symbol    string    `json:"symbol"`
	Interval  string    `json:"interval"`
	OpenTime  time.Time `json:"openTime"`
	CloseTime time.Time `json:"closeTime"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
	IsClosed  bool      `json:"isClosed"`
}

// OrderSide, OrderType, and OrderStatus are defined in enums.go

// Order represents a trading order

// AssetBalance represents the balance of a single asset
type AssetBalance struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free"`
	Locked float64 `json:"locked"`
	Total  float64 `json:"total"`
	Price  float64 `json:"price"` // Current price in USDT
}

// Wallet represents the account's wallet with balances for multiple assets
// Commented out to avoid redeclaration with models.go
// type Wallet struct {
// 	Balances  map[string]*AssetBalance `json:"balances"`
// 	UpdatedAt time.Time                `json:"updatedAt"`
// }

// Wallet represents the balance information for a specific currency.
// type Wallet struct {
// 	Currency  string    `json:"currency"`
// 	Available float64   `json:"available,string"`
// 	Frozen    float64   `json:"frozen,string"`
// 	Total     float64   `json:"total,string"`
// 	UpdatedAt time.Time `json:"updatedAt"` // Assuming we want to track when it was last updated
// }
