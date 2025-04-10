package storage

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// MarketDataType represents the type of market data
type MarketDataType string

const (
	MarketDataTypeCandle     MarketDataType = "candle"
	MarketDataTypeTrade      MarketDataType = "trade"
	MarketDataTypeOrderBook  MarketDataType = "orderbook"
	MarketDataTypeTicker     MarketDataType = "ticker"
	MarketDataTypePriceLevel MarketDataType = "price_level"
)

// StorageLevel represents the storage tier level
type StorageLevel string

const (
	StorageLevelHot  StorageLevel = "hot"  // In-memory, fastest access
	StorageLevelWarm StorageLevel = "warm" // Recent data, fast access (e.g., SQLite)
	StorageLevelCold StorageLevel = "cold" // Historical data, slower access (e.g., compressed files)
)

// MarketDataFilter defines criteria for querying market data
type MarketDataFilter struct {
	Symbol    string
	DataType  MarketDataType
	StartTime time.Time
	EndTime   time.Time
	Interval  string // For candle data
	Limit     int    // Maximum number of records to return
	OrderBy   string // Sorting field
	OrderDir  string // Sort direction (asc/desc)
}

// MarketDataStorage defines the interface for storing and retrieving market data
type MarketDataStorage interface {
	// Store methods
	StoreCandles(ctx context.Context, symbol string, candles []*models.Candle) error
	StoreTrades(ctx context.Context, symbol string, trades []*models.MarketTrade) error
	StoreOrderBook(ctx context.Context, symbol string, orderBook *models.OrderBook) error
	StoreTicker(ctx context.Context, symbol string, ticker *models.Ticker) error

	// Retrieve methods
	GetCandles(ctx context.Context, filter MarketDataFilter) ([]*models.Candle, error)
	GetTrades(ctx context.Context, filter MarketDataFilter) ([]*models.MarketTrade, error)
	GetOrderBook(ctx context.Context, symbol string, timestamp time.Time) (*models.OrderBook, error)
	GetTicker(ctx context.Context, symbol string) (*models.Ticker, error)

	// Aggregation methods
	GetVWAP(ctx context.Context, symbol string, startTime, endTime time.Time) (float64, error)
	GetOHLCV(ctx context.Context, symbol string, interval string, startTime, endTime time.Time) ([]*models.Candle, error)
	GetVolume(ctx context.Context, symbol string, startTime, endTime time.Time) (float64, error)

	// Maintenance methods
	Cleanup(ctx context.Context, olderThan time.Time, dataType MarketDataType) error
	Migrate(ctx context.Context, fromLevel, toLevel StorageLevel, olderThan time.Time) error
	Vacuum(ctx context.Context) error

	// Utility methods
	GetDataTypes(ctx context.Context, symbol string) ([]MarketDataType, error)
	GetTimeRange(ctx context.Context, symbol string, dataType MarketDataType) (start, end time.Time, err error)
	GetSymbols(ctx context.Context) ([]string, error)
}
