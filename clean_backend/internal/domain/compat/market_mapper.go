package compat

import (
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// This package provides compatibility between deprecated models and canonical models
// It defines the deprecated models and provides conversion functions

// Deprecated models

// Interval represents a time interval for candlestick data
type Interval string

// Interval constants
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

// Candle represents candlestick data for a symbol
type Candle struct {
	Symbol      string    `json:"symbol"`
	Exchange    string    `json:"exchange"`
	Interval    Interval  `json:"interval"`
	OpenTime    time.Time `json:"openTime"`
	CloseTime   time.Time `json:"closeTime"`
	Open        float64   `json:"open"`
	High        float64   `json:"high"`
	Low         float64   `json:"low"`
	Close       float64   `json:"close"`
	Volume      float64   `json:"volume"`
	QuoteVolume float64   `json:"quoteVolume"`
	TradeCount  int64     `json:"tradeCount"`
	Complete    bool      `json:"complete"`
}

// OrderBookEntry represents a single entry in an order book
type OrderBookEntry struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

// OrderBook represents an order book for a symbol
type OrderBook struct {
	Symbol    string           `json:"symbol"`
	Exchange  string           `json:"exchange"`
	Bids      []OrderBookEntry `json:"bids"`
	Asks      []OrderBookEntry `json:"asks"`
	Timestamp time.Time        `json:"timestamp"`
}

// Ticker represents ticker data for a symbol
type Ticker struct {
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

// MarketTrade represents a trade that occurred on the market
type MarketTrade struct {
	ID            string    `json:"id"`
	Symbol        string    `json:"symbol"`
	Exchange      string    `json:"exchange"`
	Price         float64   `json:"price"`
	Quantity      float64   `json:"quantity"`
	QuoteQuantity float64   `json:"quoteQuantity"`
	Time          time.Time `json:"time"`
	IsBuyerMaker  bool      `json:"isBuyerMaker"`
	IsBestMatch   bool      `json:"isBestMatch"`
}

// Conversion functions

// IntervalToKlineInterval converts a deprecated Interval to a canonical KlineInterval
func IntervalToKlineInterval(interval Interval) model.KlineInterval {
	return model.KlineInterval(interval)
}

// KlineIntervalToInterval converts a canonical KlineInterval to a deprecated Interval
func KlineIntervalToInterval(interval model.KlineInterval) Interval {
	return Interval(interval)
}

// CandleToKline converts a deprecated Candle to a canonical Kline
func CandleToKline(candle *Candle) *model.Kline {
	if candle == nil {
		return nil
	}
	return &model.Kline{
		Symbol:      candle.Symbol,
		Exchange:    candle.Exchange,
		Interval:    model.KlineInterval(candle.Interval),
		OpenTime:    candle.OpenTime,
		CloseTime:   candle.CloseTime,
		Open:        candle.Open,
		High:        candle.High,
		Low:         candle.Low,
		Close:       candle.Close,
		Volume:      candle.Volume,
		QuoteVolume: candle.QuoteVolume,
		TradeCount:  candle.TradeCount,
		Complete:    candle.Complete,
	}
}

// KlineToCandle converts a canonical Kline to a deprecated Candle
func KlineToCandle(kline *model.Kline) *Candle {
	if kline == nil {
		return nil
	}
	return &Candle{
		Symbol:      kline.Symbol,
		Exchange:    kline.Exchange,
		Interval:    Interval(kline.Interval),
		OpenTime:    kline.OpenTime,
		CloseTime:   kline.CloseTime,
		Open:        kline.Open,
		High:        kline.High,
		Low:         kline.Low,
		Close:       kline.Close,
		Volume:      kline.Volume,
		QuoteVolume: kline.QuoteVolume,
		TradeCount:  kline.TradeCount,
		Complete:    kline.Complete,
	}
}

// OrderBookToCanonical converts a deprecated OrderBook to a canonical OrderBook
func OrderBookToCanonical(orderBook *OrderBook) *model.OrderBook {
	if orderBook == nil {
		return nil
	}

	// Convert bids
	bids := make([]model.OrderBookEntry, len(orderBook.Bids))
	for i, bid := range orderBook.Bids {
		bids[i] = model.OrderBookEntry{
			Price:    bid.Price,
			Quantity: bid.Quantity,
		}
	}

	// Convert asks
	asks := make([]model.OrderBookEntry, len(orderBook.Asks))
	for i, ask := range orderBook.Asks {
		asks[i] = model.OrderBookEntry{
			Price:    ask.Price,
			Quantity: ask.Quantity,
		}
	}

	return &model.OrderBook{
		Symbol:    orderBook.Symbol,
		Exchange:  orderBook.Exchange,
		Bids:      bids,
		Asks:      asks,
		Timestamp: orderBook.Timestamp,
	}
}

// CanonicalToOrderBook converts a canonical OrderBook to a deprecated OrderBook
func CanonicalToOrderBook(orderBook *model.OrderBook) *OrderBook {
	if orderBook == nil {
		return nil
	}

	// Convert bids
	bids := make([]OrderBookEntry, len(orderBook.Bids))
	for i, bid := range orderBook.Bids {
		bids[i] = OrderBookEntry{
			Price:    bid.Price,
			Quantity: bid.Quantity,
		}
	}

	// Convert asks
	asks := make([]OrderBookEntry, len(orderBook.Asks))
	for i, ask := range orderBook.Asks {
		asks[i] = OrderBookEntry{
			Price:    ask.Price,
			Quantity: ask.Quantity,
		}
	}

	return &OrderBook{
		Symbol:    orderBook.Symbol,
		Exchange:  orderBook.Exchange,
		Bids:      bids,
		Asks:      asks,
		Timestamp: orderBook.Timestamp,
	}
}

// TickerToCanonical converts a deprecated Ticker to a canonical Ticker
func TickerToCanonical(ticker *Ticker) *model.Ticker {
	if ticker == nil {
		return nil
	}
	return &model.Ticker{
		Symbol:             ticker.Symbol,
		Exchange:           ticker.Exchange,
		LastPrice:          ticker.LastPrice,
		PriceChange:        ticker.PriceChange,
		PriceChangePercent: ticker.PriceChangePercent,
		HighPrice:          ticker.HighPrice,
		LowPrice:           ticker.LowPrice,
		Volume:             ticker.Volume,
		QuoteVolume:        ticker.QuoteVolume,
		OpenPrice:          ticker.OpenPrice,
		PrevClosePrice:     ticker.PrevClosePrice,
		BidPrice:           ticker.BidPrice,
		BidQty:             ticker.BidQty,
		AskPrice:           ticker.AskPrice,
		AskQty:             ticker.AskQty,
		Count:              ticker.Count,
		Timestamp:          ticker.Timestamp,
	}
}

// CanonicalToTicker converts a canonical Ticker to a deprecated Ticker
func CanonicalToTicker(ticker *model.Ticker) *Ticker {
	if ticker == nil {
		return nil
	}
	return &Ticker{
		Symbol:             ticker.Symbol,
		Exchange:           ticker.Exchange,
		LastPrice:          ticker.LastPrice,
		PriceChange:        ticker.PriceChange,
		PriceChangePercent: ticker.PriceChangePercent,
		HighPrice:          ticker.HighPrice,
		LowPrice:           ticker.LowPrice,
		Volume:             ticker.Volume,
		QuoteVolume:        ticker.QuoteVolume,
		OpenPrice:          ticker.OpenPrice,
		PrevClosePrice:     ticker.PrevClosePrice,
		BidPrice:           ticker.BidPrice,
		BidQty:             ticker.BidQty,
		AskPrice:           ticker.AskPrice,
		AskQty:             ticker.AskQty,
		Count:              ticker.Count,
		Timestamp:          ticker.Timestamp,
	}
}

// TickerToMarketTicker converts a deprecated Ticker to a canonical MarketTicker
func TickerToMarketTicker(ticker *Ticker) *model.MarketTicker {
	if ticker == nil {
		return nil
	}
	return &model.MarketTicker{
		Symbol:        ticker.Symbol,
		Exchange:      ticker.Exchange,
		Price:         ticker.LastPrice,
		Volume:        ticker.Volume,
		High24h:       ticker.HighPrice,
		Low24h:        ticker.LowPrice,
		PriceChange:   ticker.PriceChange,
		PercentChange: ticker.PriceChangePercent,
		LastUpdated:   ticker.Timestamp,
	}
}

// MarketTickerToTicker converts a canonical MarketTicker to a deprecated Ticker
func MarketTickerToTicker(ticker *model.MarketTicker) *Ticker {
	if ticker == nil {
		return nil
	}
	return &Ticker{
		Symbol:             ticker.Symbol,
		Exchange:           ticker.Exchange,
		LastPrice:          ticker.Price,
		Volume:             ticker.Volume,
		HighPrice:          ticker.High24h,
		LowPrice:           ticker.Low24h,
		PriceChange:        ticker.PriceChange,
		PriceChangePercent: ticker.PercentChange,
		Timestamp:          ticker.LastUpdated,
	}
}
