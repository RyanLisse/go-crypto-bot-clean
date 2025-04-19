package entity

import (
	"time"

	domainmodel "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// MexcTickerEntityCanonical represents market ticker data stored in the database
type MexcTickerEntityCanonical struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement"`
	Symbol             string    `gorm:"not null;index:idx_mexc_tickers_canonical_symbol_time"`
	Exchange           string    `gorm:"not null;index:idx_mexc_tickers_canonical_exchange"`
	LastPrice          float64   `gorm:"not null"`
	Volume             float64   `gorm:"not null"`
	HighPrice          float64   `gorm:"not null"`
	LowPrice           float64   `gorm:"not null"`
	PriceChange        float64   `gorm:"not null"`
	PriceChangePercent float64   `gorm:"not null"`
	Timestamp          time.Time `gorm:"not null;index:idx_mexc_tickers_canonical_symbol_time,priority:2"`
	CreatedAt          time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName specifies the table name for MexcTickerEntityCanonical
func (MexcTickerEntityCanonical) TableName() string {
	return "mexc_tickers_canonical"
}

// FromModel converts a domain model to an entity
func (e *MexcTickerEntityCanonical) FromModel(ticker *domainmodel.Ticker) {
	e.Symbol = ticker.Symbol
	e.Exchange = ticker.Exchange
	e.LastPrice = ticker.LastPrice
	e.Volume = ticker.Volume
	e.HighPrice = ticker.HighPrice
	e.LowPrice = ticker.LowPrice
	e.PriceChange = ticker.PriceChange
	e.PriceChangePercent = ticker.PriceChangePercent
	e.Timestamp = ticker.Timestamp
}

// ToModel converts an entity to a domain model
func (e *MexcTickerEntityCanonical) ToModel() *domainmodel.Ticker {
	return &domainmodel.Ticker{
		Symbol:             e.Symbol,
		Exchange:           e.Exchange,
		LastPrice:          e.LastPrice,
		Volume:             e.Volume,
		HighPrice:          e.HighPrice,
		LowPrice:           e.LowPrice,
		PriceChange:        e.PriceChange,
		PriceChangePercent: e.PriceChangePercent,
		Timestamp:          e.Timestamp,
	}
}

// KlineEntityCanonical represents candle (kline) data stored in the database
type KlineEntityCanonical struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Symbol      string    `gorm:"not null;index:idx_klines_canonical_symbol_interval_time"`
	Exchange    string    `gorm:"not null;index:idx_klines_canonical_exchange"`
	Interval    string    `gorm:"not null;index:idx_klines_canonical_symbol_interval_time,priority:2"`
	OpenTime    time.Time `gorm:"not null;index:idx_klines_canonical_symbol_interval_time,priority:3"`
	CloseTime   time.Time `gorm:"not null"`
	Open        float64   `gorm:"not null"`
	High        float64   `gorm:"not null"`
	Low         float64   `gorm:"not null"`
	Close       float64   `gorm:"not null"`
	Volume      float64   `gorm:"not null"`
	QuoteVolume float64   `gorm:"not null"`
	TradeCount  int64     `gorm:"not null"`
	Complete    bool      `gorm:"not null;default:false"`
	IsClosed    bool      `gorm:"not null;default:false"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName specifies the table name for KlineEntityCanonical
func (KlineEntityCanonical) TableName() string {
	return "klines_canonical"
}

// FromModel converts a domain model to an entity
func (e *KlineEntityCanonical) FromModel(kline *domainmodel.Kline) {
	e.Symbol = kline.Symbol
	e.Exchange = kline.Exchange
	e.Interval = string(kline.Interval)
	e.OpenTime = kline.OpenTime
	e.CloseTime = kline.CloseTime
	e.Open = kline.Open
	e.High = kline.High
	e.Low = kline.Low
	e.Close = kline.Close
	e.Volume = kline.Volume
	e.QuoteVolume = kline.QuoteVolume
	e.TradeCount = kline.TradeCount
	e.Complete = kline.Complete
	e.IsClosed = kline.IsClosed
}

// ToModel converts an entity to a domain model
func (e *KlineEntityCanonical) ToModel() *domainmodel.Kline {
	return &domainmodel.Kline{
		Symbol:      e.Symbol,
		Exchange:    e.Exchange,
		Interval:    domainmodel.KlineInterval(e.Interval),
		OpenTime:    e.OpenTime,
		CloseTime:   e.CloseTime,
		Open:        e.Open,
		High:        e.High,
		Low:         e.Low,
		Close:       e.Close,
		Volume:      e.Volume,
		QuoteVolume: e.QuoteVolume,
		TradeCount:  e.TradeCount,
		Complete:    e.Complete,
		IsClosed:    e.IsClosed,
	}
}

// OrderBookEntityCanonical represents order book data stored in the database
type OrderBookEntityCanonical struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Symbol       string    `gorm:"not null;index:idx_orderbooks_canonical_symbol_time"`
	Exchange     string    `gorm:"not null;index:idx_orderbooks_canonical_exchange"`
	LastUpdateID int64     `gorm:"not null"`
	SequenceNum  int64     `gorm:"not null"`
	BidsJSON     string    `gorm:"type:text;not null"` // JSON string of bids
	AsksJSON     string    `gorm:"type:text;not null"` // JSON string of asks
	Timestamp    time.Time `gorm:"not null;index:idx_orderbooks_canonical_symbol_time,priority:2"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName specifies the table name for OrderBookEntityCanonical
func (OrderBookEntityCanonical) TableName() string {
	return "orderbooks_canonical"
}

// ToModel converts an entity to a domain model
func (e *OrderBookEntityCanonical) ToModel() *domainmodel.OrderBook {
	return &domainmodel.OrderBook{
		Symbol:       e.Symbol,
		Exchange:     e.Exchange,
		LastUpdateID: e.LastUpdateID,
		SequenceNum:  e.SequenceNum,
		Timestamp:    e.Timestamp,
		// Bids and Asks are handled separately due to JSON serialization
	}
}
