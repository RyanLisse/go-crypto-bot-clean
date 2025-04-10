package services

import (
	"context"
	"time"
)

// Candlestick represents OHLCV data for a time period
type Candlestick struct {
	Symbol    string    `json:"symbol"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

// MarketTrade represents a trade from market data
type MarketTrade struct {
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	Side      string    `json:"side"`
	Timestamp time.Time `json:"timestamp"`
}

// OrderBookEntry represents a single entry in the order book
type OrderBookEntry struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

// OrderBook represents the current state of the order book
type OrderBook struct {
	Symbol string           `json:"symbol"`
	Bids   []OrderBookEntry `json:"bids"`
	Asks   []OrderBookEntry `json:"asks"`
	Time   time.Time        `json:"time"`
}

// MarketDataService defines the interface for market data operations
type MarketDataService interface {
	// Real-time price data
	GetCurrentPrice(ctx context.Context, symbol string) (float64, error)
	SubscribePriceUpdates(ctx context.Context, symbol string) (<-chan float64, error)

	// Historical data
	GetCandles(ctx context.Context, symbol string, interval string, start, end time.Time) ([]Candlestick, error)
	GetHistoricalTrades(ctx context.Context, symbol string, start, end time.Time) ([]MarketTrade, error)

	// Order book data
	GetOrderBook(ctx context.Context, symbol string, depth int) (*OrderBook, error)
	SubscribeOrderBookUpdates(ctx context.Context, symbol string) (<-chan *OrderBook, error)

	// Market statistics
	GetVolume24h(ctx context.Context, symbol string) (float64, error)
	GetPriceChange24h(ctx context.Context, symbol string) (float64, error)
}
