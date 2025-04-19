package port

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// MarketService defines the interface for market data services
type MarketService interface {
	// Ticker operations
	GetTicker(ctx context.Context, symbol, exchange string) (*model.Ticker, error)
	GetAllTickers(ctx context.Context, exchange string) ([]*model.Ticker, error)

	// Kline operations
	GetKlines(ctx context.Context, symbol, exchange string, interval model.KlineInterval, limit int) ([]*model.Kline, error)
	GetHistoricalKlines(ctx context.Context, symbol, exchange string, interval model.KlineInterval, startTime, endTime time.Time, limit int) ([]*model.Kline, error)

	// Symbol operations
	GetSymbol(ctx context.Context, symbol, exchange string) (*model.Symbol, error)
	GetAllSymbols(ctx context.Context, exchange string) ([]*model.Symbol, error)

	// Market data synchronization
	SyncMarketData(ctx context.Context, exchange string) error
	SyncSymbols(ctx context.Context, exchange string) error
}

// MarketDataRepository defines the interface for market data persistence
type MarketDataRepository interface {
	// Ticker operations
	SaveTicker(ctx context.Context, ticker *model.Ticker) error
	GetTicker(ctx context.Context, symbol, exchange string) (*model.Ticker, error)
	GetAllTickers(ctx context.Context, exchange string) ([]*model.Ticker, error)
	GetLatestTickers(ctx context.Context, limit int) ([]*model.Ticker, error)

	// Kline operations
	SaveKline(ctx context.Context, kline *model.Kline) error
	SaveKlines(ctx context.Context, klines []*model.Kline) error
	GetKline(ctx context.Context, symbol, exchange string, interval model.KlineInterval, openTime time.Time) (*model.Kline, error)
	GetKlines(ctx context.Context, symbol, exchange string, interval model.KlineInterval, start, end time.Time, limit int) ([]*model.Kline, error)
	GetLatestKline(ctx context.Context, symbol, exchange string, interval model.KlineInterval) (*model.Kline, error)

	// Data maintenance
	PurgeOldData(ctx context.Context, olderThan time.Time) error
}

// Note: SymbolRepository has been moved to symbol_repository.go

// TickerRepository defines the interface for ticker data persistence
type TickerRepository interface {
	Save(ctx context.Context, ticker *model.Ticker) error
	GetBySymbol(ctx context.Context, symbol string) (*model.Ticker, error)
	GetAll(ctx context.Context) ([]*model.Ticker, error)
	GetRecent(ctx context.Context, limit int) ([]*model.Ticker, error)
	SaveKline(ctx context.Context, kline *model.Kline) error
	GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, from, to time.Time, limit int) ([]*model.Kline, error)
}
