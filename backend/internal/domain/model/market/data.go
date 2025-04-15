package market

import (
	"time"
)

// KlineInterval represents the time interval of a candle/kline
type KlineInterval string

const (
	// KlineInterval1m represents 1 minute interval
	KlineInterval1m KlineInterval = "1m"
	// KlineInterval5m represents 5 minute interval
	KlineInterval5m KlineInterval = "5m"
	// KlineInterval15m represents 15 minute interval
	KlineInterval15m KlineInterval = "15m"
	// KlineInterval1h represents 1 hour interval
	KlineInterval1h KlineInterval = "1h"
	// KlineInterval4h represents 4 hour interval
	KlineInterval4h KlineInterval = "4h"
	// KlineInterval1d represents 1 day interval
	KlineInterval1d KlineInterval = "1d"
)

// Kline represents a candlestick chart data point
type Kline struct {
	Symbol    string        `json:"symbol"`
	Interval  KlineInterval `json:"interval"`
	OpenTime  time.Time     `json:"openTime"`
	CloseTime time.Time     `json:"closeTime"`
	Open      float64       `json:"open"`
	High      float64       `json:"high"`
	Low       float64       `json:"low"`
	Close     float64       `json:"close"`
	Volume    float64       `json:"volume"`
}

// HistoricalData contains historical market data like klines (candles)
type HistoricalData struct {
	// Klines contains historical candle data
	Klines []Kline `json:"klines"`
}

// Data represents combined market data for a symbol including both current and historical data
type Data struct {
	// Symbol identifier
	Symbol string `json:"symbol"`

	// CurrentData contains the most recent ticker data
	CurrentData Ticker `json:"currentData"`

	// HistoricalData contains various historical data points
	HistoricalData HistoricalData `json:"historicalData"`
}
