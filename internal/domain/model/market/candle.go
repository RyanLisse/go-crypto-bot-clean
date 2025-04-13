package market

import (
	"time"
)

// Interval represents a candle time interval
type Interval string

// Candle interval constants
const (
	Interval1m  Interval = "1m"
	Interval3m  Interval = "3m"
	Interval5m  Interval = "5m"
	Interval15m Interval = "15m"
	Interval30m Interval = "30m"
	Interval1h  Interval = "1h"
	Interval2h  Interval = "2h"
	Interval4h  Interval = "4h"
	Interval6h  Interval = "6h"
	Interval8h  Interval = "8h"
	Interval12h Interval = "12h"
	Interval1d  Interval = "1d"
	Interval3d  Interval = "3d"
	Interval1w  Interval = "1w"
	Interval1M  Interval = "1M"
)

// Candle represents OHLCV data for a trading pair
type Candle struct {
	// Exchange is the name of the exchange (e.g., "binance")
	Exchange string `json:"exchange"`

	// Symbol is the trading pair (e.g., "BTCUSDT")
	Symbol string `json:"symbol"`

	// Interval is the time interval of this candle
	Interval Interval `json:"interval"`

	// OpenTime is when this candle opened
	OpenTime time.Time `json:"openTime"`

	// CloseTime is when this candle closed
	CloseTime time.Time `json:"closeTime"`

	// Open is the opening price
	Open float64 `json:"open"`

	// High is the highest price during the interval
	High float64 `json:"high"`

	// Low is the lowest price during the interval
	Low float64 `json:"low"`

	// Close is the closing price
	Close float64 `json:"close"`

	// Volume is the trading volume
	Volume float64 `json:"volume"`

	// QuoteVolume is the quote asset volume
	QuoteVolume float64 `json:"quoteVolume"`

	// TradeCount is the number of trades during the interval
	TradeCount int64 `json:"tradeCount"`
}
