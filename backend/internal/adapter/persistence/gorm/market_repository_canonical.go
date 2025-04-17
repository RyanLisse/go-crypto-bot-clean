package gorm

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Ensure MarketRepositoryCanonical implements the proper interfaces
var _ port.MarketRepository = (*MarketRepositoryCanonical)(nil)
var _ port.SymbolRepository = (*MarketRepositoryCanonical)(nil)

// MarketRepositoryCanonical implements the port.MarketRepository interface
// by forwarding calls to the legacy MarketRepository implementation
type MarketRepositoryCanonical struct {
	legacy *MarketRepository // For method implementations
}

// NewMarketRepositoryCanonical creates a new MarketRepositoryCanonical
func NewMarketRepositoryCanonical(db *gorm.DB, logger *zerolog.Logger) *MarketRepositoryCanonical {
	return &MarketRepositoryCanonical{
		legacy: NewMarketRepository(db, logger),
	}
}

// SaveTicker stores a ticker in the database
func (r *MarketRepositoryCanonical) SaveTicker(ctx context.Context, ticker *model.Ticker) error {
	return r.legacy.SaveTicker(ctx, ticker)
}

// GetTicker retrieves the latest ticker for a symbol from a specific exchange
func (r *MarketRepositoryCanonical) GetTicker(ctx context.Context, symbol, exchange string) (*model.Ticker, error) {
	return r.legacy.GetTicker(ctx, symbol, exchange)
}

// GetAllTickers retrieves all latest tickers from a specific exchange
func (r *MarketRepositoryCanonical) GetAllTickers(ctx context.Context, exchange string) ([]*model.Ticker, error) {
	return r.legacy.GetAllTickers(ctx, exchange)
}

// GetTickerHistory retrieves ticker history for a symbol within a time range
func (r *MarketRepositoryCanonical) GetTickerHistory(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*model.Ticker, error) {
	return r.legacy.GetTickerHistory(ctx, symbol, exchange, start, end)
}

// SaveKline stores a kline/candle in the database
func (r *MarketRepositoryCanonical) SaveKline(ctx context.Context, kline *model.Kline) error {
	return r.legacy.SaveKline(ctx, kline)
}

// SaveKlines stores multiple klines/candles in the database
func (r *MarketRepositoryCanonical) SaveKlines(ctx context.Context, klines []*model.Kline) error {
	return r.legacy.SaveKlines(ctx, klines)
}

// GetKline retrieves a specific kline/candle for a symbol, interval, and time
func (r *MarketRepositoryCanonical) GetKline(ctx context.Context, symbol, exchange string, interval model.KlineInterval, openTime time.Time) (*model.Kline, error) {
	return r.legacy.GetKline(ctx, symbol, exchange, interval, openTime)
}

// GetKlines retrieves klines/candles for a symbol within a time range
func (r *MarketRepositoryCanonical) GetKlines(ctx context.Context, symbol, exchange string, interval model.KlineInterval, start, end time.Time, limit int) ([]*model.Kline, error) {
	return r.legacy.GetKlines(ctx, symbol, exchange, interval, start, end, limit)
}

// GetLatestKline retrieves the most recent kline/candle for a symbol and interval
func (r *MarketRepositoryCanonical) GetLatestKline(ctx context.Context, symbol, exchange string, interval model.KlineInterval) (*model.Kline, error) {
	return r.legacy.GetLatestKline(ctx, symbol, exchange, interval)
}

// PurgeOldData removes market data older than the specified retention period
func (r *MarketRepositoryCanonical) PurgeOldData(ctx context.Context, olderThan time.Time) error {
	return r.legacy.PurgeOldData(ctx, olderThan)
}

// GetLatestTickers retrieves the latest tickers for all symbols
func (r *MarketRepositoryCanonical) GetLatestTickers(ctx context.Context, limit int) ([]*model.Ticker, error) {
	return r.legacy.GetLatestTickers(ctx, limit)
}

// GetTickersBySymbol retrieves tickers for a specific symbol with optional time range
func (r *MarketRepositoryCanonical) GetTickersBySymbol(ctx context.Context, symbol string, limit int) ([]*model.Ticker, error) {
	return r.legacy.GetTickersBySymbol(ctx, symbol, limit)
}

// GetOrderBook retrieves the order book for a symbol
func (r *MarketRepositoryCanonical) GetOrderBook(ctx context.Context, symbol, exchange string, depth int) (*model.OrderBook, error) {
	return r.legacy.GetOrderBook(ctx, symbol, exchange, depth)
}

// Symbol Repository implementation

// Create stores a new Symbol
func (r *MarketRepositoryCanonical) Create(ctx context.Context, symbol *model.Symbol) error {
	return r.legacy.Create(ctx, symbol)
}

// GetBySymbol returns a Symbol by its symbol string (e.g., "BTCUSDT")
func (r *MarketRepositoryCanonical) GetBySymbol(ctx context.Context, symbol string) (*model.Symbol, error) {
	return r.legacy.GetBySymbol(ctx, symbol)
}

// GetByExchange returns all Symbols from a specific exchange
func (r *MarketRepositoryCanonical) GetByExchange(ctx context.Context, exchange string) ([]*model.Symbol, error) {
	return r.legacy.GetByExchange(ctx, exchange)
}

// GetAll returns all available Symbols
func (r *MarketRepositoryCanonical) GetAll(ctx context.Context) ([]*model.Symbol, error) {
	return r.legacy.GetAll(ctx)
}

// Update updates an existing Symbol
func (r *MarketRepositoryCanonical) Update(ctx context.Context, symbol *model.Symbol) error {
	return r.legacy.Update(ctx, symbol)
}

// Delete removes a Symbol
func (r *MarketRepositoryCanonical) Delete(ctx context.Context, symbol string) error {
	return r.legacy.Delete(ctx, symbol)
}

// GetSymbolsByStatus returns symbols by status with pagination
func (r *MarketRepositoryCanonical) GetSymbolsByStatus(ctx context.Context, status string, limit int, offset int) ([]*model.Symbol, error) {
	return r.legacy.GetSymbolsByStatus(ctx, status, limit, offset)
}

// Legacy methods for backward compatibility

// SaveTickerLegacy stores a ticker in the database using the legacy model
func (r *MarketRepositoryCanonical) SaveTickerLegacy(ctx context.Context, ticker *market.Ticker) error {
	return r.legacy.SaveTickerLegacy(ctx, ticker)
}

// GetTickerLegacy retrieves the latest ticker for a symbol from a specific exchange using the legacy model
func (r *MarketRepositoryCanonical) GetTickerLegacy(ctx context.Context, symbol, exchange string) (*market.Ticker, error) {
	return r.legacy.GetTickerLegacy(ctx, symbol, exchange)
}

// GetAllTickersLegacy retrieves all latest tickers from a specific exchange using the legacy model
func (r *MarketRepositoryCanonical) GetAllTickersLegacy(ctx context.Context, exchange string) ([]*market.Ticker, error) {
	return r.legacy.GetAllTickersLegacy(ctx, exchange)
}

// GetTickerHistoryLegacy retrieves ticker history for a symbol within a time range using the legacy model
func (r *MarketRepositoryCanonical) GetTickerHistoryLegacy(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*market.Ticker, error) {
	return r.legacy.GetTickerHistoryLegacy(ctx, symbol, exchange, start, end)
}

// SaveCandleLegacy stores a candle in the database using the legacy model
func (r *MarketRepositoryCanonical) SaveCandleLegacy(ctx context.Context, candle *market.Candle) error {
	return r.legacy.SaveCandleLegacy(ctx, candle)
}

// SaveCandlesLegacy stores multiple candles in the database using the legacy model
func (r *MarketRepositoryCanonical) SaveCandlesLegacy(ctx context.Context, candles []*market.Candle) error {
	return r.legacy.SaveCandlesLegacy(ctx, candles)
}

// GetCandleLegacy retrieves a specific candle for a symbol, interval, and time using the legacy model
func (r *MarketRepositoryCanonical) GetCandleLegacy(ctx context.Context, symbol, exchange string, interval market.Interval, openTime time.Time) (*market.Candle, error) {
	return r.legacy.GetCandleLegacy(ctx, symbol, exchange, interval, openTime)
}

// GetCandlesLegacy retrieves candles for a symbol within a time range using the legacy model
func (r *MarketRepositoryCanonical) GetCandlesLegacy(ctx context.Context, symbol, exchange string, interval market.Interval, start, end time.Time, limit int) ([]*market.Candle, error) {
	return r.legacy.GetCandlesLegacy(ctx, symbol, exchange, interval, start, end, limit)
}

// GetLatestCandleLegacy retrieves the most recent candle for a symbol and interval using the legacy model
func (r *MarketRepositoryCanonical) GetLatestCandleLegacy(ctx context.Context, symbol, exchange string, interval market.Interval) (*market.Candle, error) {
	return r.legacy.GetLatestCandleLegacy(ctx, symbol, exchange, interval)
}

// GetLatestTickersLegacy retrieves the latest tickers for all symbols using the legacy model
func (r *MarketRepositoryCanonical) GetLatestTickersLegacy(ctx context.Context, limit int) ([]*market.Ticker, error) {
	return r.legacy.GetLatestTickersLegacy(ctx, limit)
}

// GetTickersBySymbolLegacy retrieves tickers for a specific symbol with optional time range using the legacy model
func (r *MarketRepositoryCanonical) GetTickersBySymbolLegacy(ctx context.Context, symbol string, limit int) ([]*market.Ticker, error) {
	return r.legacy.GetTickersBySymbolLegacy(ctx, symbol, limit)
}

// GetOrderBookLegacy retrieves the order book for a symbol using the legacy model
func (r *MarketRepositoryCanonical) GetOrderBookLegacy(ctx context.Context, symbol, exchange string, depth int) (*market.OrderBook, error) {
	return r.legacy.GetOrderBookLegacy(ctx, symbol, exchange, depth)
}
