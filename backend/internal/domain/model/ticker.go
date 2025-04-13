package model

import (
	"time"
)

// Ticker represents a market data ticker
type Ticker struct {
	Symbol             string    `json:"symbol"`
	LastPrice          float64   `json:"lastPrice"`
	PriceChange        float64   `json:"priceChange"`
	PriceChangePercent float64   `json:"priceChangePercent"`
	HighPrice          float64   `json:"highPrice"`
	LowPrice           float64   `json:"lowPrice"`
	Volume             float64   `json:"volume"`
	QuoteVolume        float64   `json:"quoteVolume"`
	OpenPrice          float64   `json:"openPrice"`
	PrevClosePrice     float64   `json:"prevClosePrice"`
	BidPrice           float64   `json:"bidPrice"`
	BidQty             float64   `json:"bidQty"`
	AskPrice           float64   `json:"askPrice"`
	AskQty             float64   `json:"askQty"`
	Count              int64     `json:"count"`
	Timestamp          time.Time `json:"timestamp"`
}

// KlineInterval represents a kline/candlestick interval
type KlineInterval string

// Common kline intervals
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

// Kline represents a candlestick/kline data point
type Kline struct {
	Symbol      string        `json:"symbol"`
	Interval    KlineInterval `json:"interval"`
	OpenTime    time.Time     `json:"openTime"`
	CloseTime   time.Time     `json:"closeTime"`
	Open        float64       `json:"open"`
	High        float64       `json:"high"`
	Low         float64       `json:"low"`
	Close       float64       `json:"close"`
	Volume      float64       `json:"volume"`
	QuoteVolume float64       `json:"quoteVolume"`
	TradeCount  int64         `json:"tradeCount"`
	IsClosed    bool          `json:"isClosed"`
}

// OrderBook represents an order book snapshot
type OrderBook struct {
	Symbol       string           `json:"symbol"`
	LastUpdateID int64            `json:"lastUpdateId"`
	Bids         []OrderBookEntry `json:"bids"`
	Asks         []OrderBookEntry `json:"asks"`
	Timestamp    time.Time        `json:"timestamp"`
}

// OrderBookEntry represents a single order book entry (price level)
type OrderBookEntry struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

// MarketTrade represents a public trade reported by the exchange
type MarketTrade struct {
	ID            int64     `json:"id"`
	Symbol        string    `json:"symbol"`
	Price         float64   `json:"price"`
	Quantity      float64   `json:"quantity"`
	QuoteQuantity float64   `json:"quoteQuantity"`
	Time          time.Time `json:"time"`
	IsBuyerMaker  bool      `json:"isBuyerMaker"`
}

// TickerCache is used to store multiple tickers for quick access
type TickerCache struct {
	Tickers    map[string]*Ticker // Symbol -> Ticker
	LastUpdate time.Time
}

// NewTickerCache creates a new ticker cache
func NewTickerCache() *TickerCache {
	return &TickerCache{
		Tickers:    make(map[string]*Ticker),
		LastUpdate: time.Now(),
	}
}

// UpdateTicker updates or adds a ticker to the cache
func (c *TickerCache) UpdateTicker(ticker *Ticker) {
	c.Tickers[ticker.Symbol] = ticker
	c.LastUpdate = time.Now()
}

// GetTicker gets a ticker for a symbol from the cache
func (c *TickerCache) GetTicker(symbol string) *Ticker {
	ticker, exists := c.Tickers[symbol]
	if !exists {
		return nil
	}
	return ticker
}
