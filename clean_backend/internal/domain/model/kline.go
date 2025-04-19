package model

import (
	"time"
)

// KlineInterval represents a time interval for candle/kline data
type KlineInterval string

// KlineInterval constants
const (
	KlineInterval1m  KlineInterval = "1m"
	KlineInterval3m  KlineInterval = "3m"
	KlineInterval5m  KlineInterval = "5m"
	KlineInterval15m KlineInterval = "15m"
	KlineInterval30m KlineInterval = "30m"
	KlineInterval1h  KlineInterval = "1h"
	KlineInterval2h  KlineInterval = "2h"
	KlineInterval4h  KlineInterval = "4h"
	KlineInterval6h  KlineInterval = "6h"
	KlineInterval8h  KlineInterval = "8h"
	KlineInterval12h KlineInterval = "12h"
	KlineInterval1d  KlineInterval = "1d"
	KlineInterval3d  KlineInterval = "3d"
	KlineInterval1w  KlineInterval = "1w"
	KlineInterval1M  KlineInterval = "1M"
)

// Kline represents a candlestick/kline for a symbol
type Kline struct {
	// Symbol is the trading pair identifier (e.g., "BTCUSDT")
	Symbol string `json:"symbol"`

	// Exchange indicates which exchange this kline is from
	Exchange string `json:"exchange"`

	// Interval is the time interval for this kline
	Interval KlineInterval `json:"interval"`

	// OpenTime is the opening time of this kline
	OpenTime time.Time `json:"openTime"`

	// CloseTime is the closing time of this kline
	CloseTime time.Time `json:"closeTime"`

	// Open is the opening price
	Open float64 `json:"open"`

	// High is the highest price during this interval
	High float64 `json:"high"`

	// Low is the lowest price during this interval
	Low float64 `json:"low"`

	// Close is the closing price
	Close float64 `json:"close"`

	// Volume is the trading volume in the base asset
	Volume float64 `json:"volume"`

	// QuoteVolume is the trading volume in the quote asset
	QuoteVolume float64 `json:"quoteVolume"`

	// TradeCount is the number of trades during this interval
	TradeCount int64 `json:"tradeCount"`

	// Complete indicates if this kline is complete (interval has ended)
	Complete bool `json:"complete"`

	// IsClosed indicates if this kline is closed (no more updates)
	IsClosed bool `json:"isClosed"`
}
