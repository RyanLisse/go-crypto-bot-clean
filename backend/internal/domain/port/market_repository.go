package port

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
)

// SymbolRepository handles storage and retrieval of trading pair information
type SymbolRepository interface {
	// Create stores a new Symbol
	Create(ctx context.Context, symbol *model.Symbol) error

	// GetBySymbol returns a Symbol by its symbol string (e.g., "BTCUSDT")
	GetBySymbol(ctx context.Context, symbol string) (*model.Symbol, error)

	// GetByExchange returns all Symbols from a specific exchange
	GetByExchange(ctx context.Context, exchange string) ([]*model.Symbol, error)

	// GetAll returns all available Symbols
	GetAll(ctx context.Context) ([]*model.Symbol, error)

	// Update updates an existing Symbol
	Update(ctx context.Context, symbol *model.Symbol) error

	// Delete removes a Symbol
	Delete(ctx context.Context, symbol string) error

	// GetSymbolsByStatus returns symbols by status with pagination
	GetSymbolsByStatus(ctx context.Context, status string, limit int, offset int) ([]*model.Symbol, error)
}

// MarketRepository defines methods for storing and retrieving market data
type MarketRepository interface {
	// SaveTicker stores a ticker in the database
	SaveTicker(ctx context.Context, ticker *model.Ticker) error

	// GetTicker retrieves the latest ticker for a symbol from a specific exchange
	GetTicker(ctx context.Context, symbol, exchange string) (*model.Ticker, error)

	// GetAllTickers retrieves all latest tickers from a specific exchange
	GetAllTickers(ctx context.Context, exchange string) ([]*model.Ticker, error)

	// GetTickerHistory retrieves ticker history for a symbol within a time range
	GetTickerHistory(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*model.Ticker, error)

	// SaveKline stores a kline/candle in the database
	SaveKline(ctx context.Context, kline *model.Kline) error

	// SaveKlines stores multiple klines/candles in the database
	SaveKlines(ctx context.Context, klines []*model.Kline) error

	// GetKline retrieves a specific kline/candle for a symbol, interval, and time
	GetKline(ctx context.Context, symbol, exchange string, interval model.KlineInterval, openTime time.Time) (*model.Kline, error)

	// GetKlines retrieves klines/candles for a symbol within a time range
	GetKlines(ctx context.Context, symbol, exchange string, interval model.KlineInterval, start, end time.Time, limit int) ([]*model.Kline, error)

	// GetLatestKline retrieves the most recent kline/candle for a symbol and interval
	GetLatestKline(ctx context.Context, symbol, exchange string, interval model.KlineInterval) (*model.Kline, error)

	// PurgeOldData removes market data older than the specified retention period
	PurgeOldData(ctx context.Context, olderThan time.Time) error

	// GetLatestTickers retrieves the latest tickers for all symbols
	GetLatestTickers(ctx context.Context, limit int) ([]*model.Ticker, error)

	// GetTickersBySymbol retrieves tickers for a specific symbol with optional time range
	GetTickersBySymbol(ctx context.Context, symbol string, limit int) ([]*model.Ticker, error)

	// GetOrderBook retrieves the order book for a symbol
	GetOrderBook(ctx context.Context, symbol, exchange string, depth int) (*model.OrderBook, error)

	// Legacy methods for backward compatibility
	// These methods will be removed in a future version

	// SaveTickerLegacy stores a ticker in the database using the legacy model
	SaveTickerLegacy(ctx context.Context, ticker *market.Ticker) error

	// GetTickerLegacy retrieves the latest ticker for a symbol from a specific exchange using the legacy model
	GetTickerLegacy(ctx context.Context, symbol, exchange string) (*market.Ticker, error)

	// GetAllTickersLegacy retrieves all latest tickers from a specific exchange using the legacy model
	GetAllTickersLegacy(ctx context.Context, exchange string) ([]*market.Ticker, error)

	// GetTickerHistoryLegacy retrieves ticker history for a symbol within a time range using the legacy model
	GetTickerHistoryLegacy(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*market.Ticker, error)

	// SaveCandleLegacy stores a candle in the database using the legacy model
	SaveCandleLegacy(ctx context.Context, candle *market.Candle) error

	// SaveCandlesLegacy stores multiple candles in the database using the legacy model
	SaveCandlesLegacy(ctx context.Context, candles []*market.Candle) error

	// GetCandleLegacy retrieves a specific candle for a symbol, interval, and time using the legacy model
	GetCandleLegacy(ctx context.Context, symbol, exchange string, interval market.Interval, openTime time.Time) (*market.Candle, error)

	// GetCandlesLegacy retrieves candles for a symbol within a time range using the legacy model
	GetCandlesLegacy(ctx context.Context, symbol, exchange string, interval market.Interval, start, end time.Time, limit int) ([]*market.Candle, error)

	// GetLatestCandleLegacy retrieves the most recent candle for a symbol and interval using the legacy model
	GetLatestCandleLegacy(ctx context.Context, symbol, exchange string, interval market.Interval) (*market.Candle, error)

	// GetLatestTickersLegacy retrieves the latest tickers for all symbols using the legacy model
	GetLatestTickersLegacy(ctx context.Context, limit int) ([]*market.Ticker, error)

	// GetTickersBySymbolLegacy retrieves tickers for a specific symbol with optional time range using the legacy model
	GetTickersBySymbolLegacy(ctx context.Context, symbol string, limit int) ([]*market.Ticker, error)

	// GetOrderBookLegacy retrieves the order book for a symbol using the legacy model
	GetOrderBookLegacy(ctx context.Context, symbol, exchange string, depth int) (*market.OrderBook, error)
}
