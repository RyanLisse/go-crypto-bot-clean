package entity

import (
	"time"
)

// MexcTickerEntity represents a ticker from MEXC exchange in the database
type MexcTickerEntity struct {
	ID                 string `gorm:"primaryKey;type:varchar(50)"`
	Symbol             string `gorm:"index;type:varchar(20)"`
	Exchange           string `gorm:"index;type:varchar(20)"`
	LastPrice          float64
	Volume             float64
	HighPrice          float64
	LowPrice           float64
	PriceChange        float64
	PriceChangePercent float64
	Timestamp          time.Time `gorm:"index"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
}

// TableName specifies the table name for MexcTickerEntity
func (MexcTickerEntity) TableName() string {
	return "mexc_tickers"
}

// MexcCandleEntity represents candle (kline) data stored in the database
type MexcCandleEntity struct {
	ID          string    `gorm:"primaryKey;type:varchar(50)"`
	Symbol      string    `gorm:"index;type:varchar(20)"`
	Exchange    string    `gorm:"index;type:varchar(20)"`
	Interval    string    `gorm:"index;type:varchar(10)"`
	OpenTime    time.Time `gorm:"index"`
	CloseTime   time.Time
	Open        float64
	High        float64
	Low         float64
	Close       float64
	Volume      float64
	QuoteVolume float64
	TradeCount  int64
	Complete    bool
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

// TableName specifies the table name for MexcCandleEntity
func (MexcCandleEntity) TableName() string {
	return "mexc_candles"
}

// MexcSymbolEntity represents symbol information from MEXC
type MexcSymbolEntity struct {
	Symbol            string    `gorm:"primaryKey;type:varchar(50)"`
	BaseAsset         string    `gorm:"not null;index;type:varchar(20)"`
	QuoteAsset        string    `gorm:"not null;index;type:varchar(20)"`
	Exchange          string    `gorm:"index;type:varchar(20)"`
	Status            string    `gorm:"not null;type:varchar(20)"` // e.g., "TRADING", "BREAK", etc.
	MinPrice          float64   `gorm:"not null"`
	MaxPrice          float64   `gorm:"not null"`
	PricePrecision    int       `gorm:"not null"`
	MinQty            float64   `gorm:"not null"`
	MaxQty            float64   `gorm:"not null"`
	QtyPrecision      int       `gorm:"not null"`
	AllowedOrderTypes string    `gorm:"type:text"` // Comma-separated list of allowed order types
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for MexcSymbolEntity
func (MexcSymbolEntity) TableName() string {
	return "mexc_symbols"
}

// MexcSyncStateEntity tracks the last successful sync with MEXC API
type MexcSyncStateEntity struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Exchange     string    `gorm:"not null;index;type:varchar(20)"`
	EntityType   string    `gorm:"not null;index;type:varchar(20)"` // "symbols", "tickers", "candles", etc.
	LastSyncTime time.Time `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for MexcSyncStateEntity
func (MexcSyncStateEntity) TableName() string {
	return "mexc_sync_states"
}
