package model

import (
	"time"
)

// Ticker represents real-time market data for a symbol
type Ticker struct {
	ID                 string    `json:"id,omitempty"`
	Symbol             string    `json:"symbol"`
	Exchange           string    `json:"exchange"`
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

// NewSimpleTicker creates a new ticker with minimal information
func NewSimpleTicker(symbol string, price float64) *Ticker {
	return &Ticker{
		Symbol:    symbol,
		LastPrice: price,
		Timestamp: time.Now(),
	}
}

// ToMarketTicker converts a full Ticker to the simplified market.Ticker format
// This is provided for backward compatibility during the transition period
func (t *Ticker) ToMarketTicker() *MarketTicker {
	return &MarketTicker{
		ID:            t.ID,
		Symbol:        t.Symbol,
		Price:         t.LastPrice,
		Volume:        t.Volume,
		High24h:       t.HighPrice,
		Low24h:        t.LowPrice,
		PriceChange:   t.PriceChange,
		PercentChange: t.PriceChangePercent,
		LastUpdated:   t.Timestamp,
		Exchange:      t.Exchange,
	}
}

// MarketTicker represents the legacy market/ticker.go model
// This is provided for backward compatibility during transition
type MarketTicker struct {
	ID            string    `json:"id,omitempty"`
	Symbol        string    `json:"symbol"`
	Price         float64   `json:"price"`
	Volume        float64   `json:"volume"`
	High24h       float64   `json:"high24h"`
	Low24h        float64   `json:"low24h"`
	PriceChange   float64   `json:"priceChange"`
	PercentChange float64   `json:"percentChange"`
	LastUpdated   time.Time `json:"lastUpdated"`
	Exchange      string    `json:"exchange"`
}

// ToTicker converts a MarketTicker to the full Ticker model
func (mt *MarketTicker) ToTicker() *Ticker {
	return &Ticker{
		ID:                 mt.ID,
		Symbol:             mt.Symbol,
		LastPrice:          mt.Price,
		Volume:             mt.Volume,
		HighPrice:          mt.High24h,
		LowPrice:           mt.Low24h,
		PriceChange:        mt.PriceChange,
		PriceChangePercent: mt.PercentChange,
		Timestamp:          mt.LastUpdated,
		Exchange:           mt.Exchange,
	}
}

// KlineInterval represents a time interval for candle/kline data
type KlineInterval string

// Kline intervals
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

// Kline represents candle/kline data for a symbol
type Kline struct {
	Symbol      string        `json:"symbol"`
	Exchange    string        `json:"exchange"`
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
	Complete    bool          `json:"complete"`
	IsClosed    bool          `json:"isClosed"`
}

// OrderBook represents an order book for a symbol
type OrderBook struct {
	Symbol       string           `json:"symbol"`
	Exchange     string           `json:"exchange"`
	LastUpdateID int64            `json:"lastUpdateId"`
	SequenceNum  int64            `json:"sequenceNum"`
	Bids         []OrderBookEntry `json:"bids"`
	Asks         []OrderBookEntry `json:"asks"`
	Timestamp    time.Time        `json:"timestamp"`
}

// OrderBookEntry represents a single entry in an order book
type OrderBookEntry struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

// MarketTrade represents a market trade
type MarketTrade struct {
	ID            int64     `json:"id"`
	Symbol        string    `json:"symbol"`
	Price         float64   `json:"price"`
	Quantity      float64   `json:"quantity"`
	QuoteQuantity float64   `json:"quoteQuantity"`
	Time          time.Time `json:"time"`
	IsBuyerMaker  bool      `json:"isBuyerMaker"`
}
