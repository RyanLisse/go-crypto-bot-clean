package model

import (
	"time"
)

// Exchange represents a cryptocurrency exchange
type Exchange struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Active    bool      `json:"active"`
	APIKey    string    `json:"api_key,omitempty"`
	APISecret string    `json:"api_secret,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Additional order types extending the basic ones in order.go
const (
	// OrderTypeStopLoss represents a stop loss order
	OrderTypeStopLoss OrderType = "STOP_LOSS"
	// OrderTypeStopLossLimit represents a stop loss limit order
	OrderTypeStopLossLimit OrderType = "STOP_LOSS_LIMIT"
	// OrderTypeTakeProfit represents a take profit order
	OrderTypeTakeProfit OrderType = "TAKE_PROFIT"
	// OrderTypeTakeProfitLimit represents a take profit limit order
	OrderTypeTakeProfitLimit OrderType = "TAKE_PROFIT_LIMIT"
)

// Additional order status extending the basic ones in order.go
const (
	// OrderStatusExpired represents an expired order
	OrderStatusExpired OrderStatus = "EXPIRED"
)

// MarketTicker represents a market ticker for a symbol
type MarketTicker struct {
	Symbol        string    `json:"symbol"`
	ExchangeID    string    `json:"exchange_id"`
	Price         float64   `json:"price"`
	Volume        float64   `json:"volume"`
	High          float64   `json:"high"`
	Low           float64   `json:"low"`
	Change        float64   `json:"change"`
	ChangePercent float64   `json:"change_percent"`
	Bid           float64   `json:"bid"`
	Ask           float64   `json:"ask"`
	LastUpdated   time.Time `json:"last_updated"`
}

// CandleInterval represents the interval of a candle
type CandleInterval string

const (
	// CandleInterval1m represents a 1-minute candle
	CandleInterval1m CandleInterval = "1m"
	// CandleInterval5m represents a 5-minute candle
	CandleInterval5m CandleInterval = "5m"
	// CandleInterval15m represents a 15-minute candle
	CandleInterval15m CandleInterval = "15m"
	// CandleInterval30m represents a 30-minute candle
	CandleInterval30m CandleInterval = "30m"
	// CandleInterval1h represents a 1-hour candle
	CandleInterval1h CandleInterval = "1h"
	// CandleInterval4h represents a 4-hour candle
	CandleInterval4h CandleInterval = "4h"
	// CandleInterval1d represents a 1-day candle
	CandleInterval1d CandleInterval = "1d"
	// CandleInterval1w represents a 1-week candle
	CandleInterval1w CandleInterval = "1w"
	// CandleInterval1M represents a 1-month candle
	CandleInterval1M CandleInterval = "1M"
)

// Candle represents a candlestick for a symbol
type Candle struct {
	Symbol     string         `json:"symbol"`
	ExchangeID string         `json:"exchange_id"`
	Interval   CandleInterval `json:"interval"`
	OpenTime   time.Time      `json:"open_time"`
	CloseTime  time.Time      `json:"close_time"`
	Open       float64        `json:"open"`
	High       float64        `json:"high"`
	Low        float64        `json:"low"`
	Close      float64        `json:"close"`
	Volume     float64        `json:"volume"`
}
