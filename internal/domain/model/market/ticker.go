package market

import (
	"time"
)

// Ticker represents current market data for a trading pair
type Ticker struct {
	// Exchange is the name of the exchange (e.g., "binance", "kucoin")
	Exchange string `json:"exchange"`

	// Symbol is the trading pair (e.g., "BTCUSDT")
	Symbol string `json:"symbol"`

	// Price is the current price
	Price float64 `json:"price"`

	// Volume is the 24h trading volume
	Volume float64 `json:"volume"`

	// High is the 24h highest price
	High float64 `json:"high"`

	// Low is the 24h lowest price
	Low float64 `json:"low"`

	// ChangePercent is the 24h price change percent
	ChangePercent float64 `json:"changePercent"`

	// QuoteVolume is the 24h quote asset volume
	QuoteVolume float64 `json:"quoteVolume"`

	// Bid is the best bid price
	Bid float64 `json:"bid"`

	// Ask is the best ask price
	Ask float64 `json:"ask"`

	// UpdateTime is when this ticker was last updated
	UpdateTime time.Time `json:"updateTime"`
}
