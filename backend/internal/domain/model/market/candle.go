package market

import "time"

// Interval represents the timeframe of a candle
type Interval string

// Available candle intervals
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

// Candle represents OHLCV (Open, High, Low, Close, Volume) data for a trading pair
type Candle struct {
	// Symbol is the trading pair identifier (e.g., "BTCUSDT")
	Symbol string `json:"symbol"`

	// Exchange indicates which exchange this candle is from
	Exchange string `json:"exchange"`

	// Interval represents the timeframe of this candle
	Interval Interval `json:"interval"`

	// OpenTime is the opening time of this candle
	OpenTime time.Time `json:"openTime"`

	// CloseTime is the closing time of this candle
	CloseTime time.Time `json:"closeTime"`

	// Open is the opening price of this candle
	Open float64 `json:"open"`

	// High is the highest price during this candle period
	High float64 `json:"high"`

	// Low is the lowest price during this candle period
	Low float64 `json:"low"`

	// Close is the closing price of this candle
	Close float64 `json:"close"`

	// Volume is the trading volume in base asset during this candle period
	Volume float64 `json:"volume"`

	// QuoteVolume is the trading volume in quote asset during this candle period
	QuoteVolume float64 `json:"quoteVolume"`

	// TradeCount is the number of trades during this candle period
	TradeCount int64 `json:"tradeCount"`

	// Complete indicates if this candle is completed (true) or still in progress (false)
	Complete bool `json:"complete"`
}
