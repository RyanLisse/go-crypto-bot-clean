package port

import (
	"context"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model/market"
)

// SymbolRepository handles storage and retrieval of trading pair information
type SymbolRepository interface {
	// Create stores a new Symbol
	Create(ctx context.Context, symbol *market.Symbol) error

	// GetBySymbol returns a Symbol by its symbol string (e.g., "BTCUSDT")
	GetBySymbol(ctx context.Context, symbol string) (*market.Symbol, error)

	// GetByExchange returns all Symbols from a specific exchange
	GetByExchange(ctx context.Context, exchange string) ([]*market.Symbol, error)

	// GetAll returns all available Symbols
	GetAll(ctx context.Context) ([]*market.Symbol, error)

	// Update updates an existing Symbol
	Update(ctx context.Context, symbol *market.Symbol) error

	// Delete removes a Symbol
	Delete(ctx context.Context, symbol string) error
}

// MarketRepository defines methods for storing and retrieving market data
type MarketRepository interface {
	// SaveTicker stores a ticker in the database
	SaveTicker(ctx context.Context, ticker *market.Ticker) error

	// GetTicker retrieves the latest ticker for a symbol from a specific exchange
	GetTicker(ctx context.Context, symbol, exchange string) (*market.Ticker, error)

	// GetAllTickers retrieves all latest tickers from a specific exchange
	GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, error)

	// GetTickerHistory retrieves ticker history for a symbol within a time range
	GetTickerHistory(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*market.Ticker, error)

	// SaveCandle stores a candle in the database
	SaveCandle(ctx context.Context, candle *market.Candle) error

	// SaveCandles stores multiple candles in the database
	SaveCandles(ctx context.Context, candles []*market.Candle) error

	// GetCandle retrieves a specific candle for a symbol, interval, and time
	GetCandle(ctx context.Context, symbol, exchange string, interval market.Interval, openTime time.Time) (*market.Candle, error)

	// GetCandles retrieves candles for a symbol within a time range
	GetCandles(ctx context.Context, symbol, exchange string, interval market.Interval, start, end time.Time, limit int) ([]*market.Candle, error)

	// GetLatestCandle retrieves the most recent candle for a symbol and interval
	GetLatestCandle(ctx context.Context, symbol, exchange string, interval market.Interval) (*market.Candle, error)

	// PurgeOldData removes market data older than the specified retention period
	PurgeOldData(ctx context.Context, olderThan time.Time) error

	// GetLatestTickers retrieves the latest tickers for all symbols
	GetLatestTickers(ctx context.Context, limit int) ([]*market.Ticker, error)

	// GetTickersBySymbol retrieves tickers for a specific symbol with optional time range
	GetTickersBySymbol(ctx context.Context, symbol string, limit int) ([]*market.Ticker, error)
}
