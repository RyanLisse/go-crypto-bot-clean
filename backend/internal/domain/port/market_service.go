package port

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
)

// MarketDataService defines the interface for market data operations
type MarketDataService interface {
	// GetTicker retrieves the current ticker for a symbol
	GetTicker(ctx context.Context, symbol string) (*model.Ticker, error)

	// GetCandles retrieves historical candlestick data
	GetCandles(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error)

	// GetOrderBook retrieves the current order book for a symbol
	GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error)

	// GetAllSymbols retrieves all available trading symbols
	GetAllSymbols(ctx context.Context) ([]*model.Symbol, error)

	// GetSymbolInfo retrieves detailed information about a specific symbol
	GetSymbolInfo(ctx context.Context, symbol string) (*model.Symbol, error)

	// GetHistoricalPrices retrieves historical prices for a symbol
	GetHistoricalPrices(ctx context.Context, symbol string, from, to time.Time, interval model.KlineInterval) ([]*model.Kline, error)

	// Deprecated methods for backward compatibility
	// These methods will be removed in a future version

	// GetTickerLegacy retrieves the current ticker for a symbol using the legacy model
	GetTickerLegacy(ctx context.Context, symbol string) (*market.Ticker, error)

	// GetCandlesLegacy retrieves historical candlestick data using the legacy model
	GetCandlesLegacy(ctx context.Context, symbol string, interval string, limit int) ([]*market.Candle, error)

	// GetOrderBookLegacy retrieves the current order book for a symbol using the legacy model
	GetOrderBookLegacy(ctx context.Context, symbol string, depth int) (*market.OrderBook, error)

	// GetAllSymbolsLegacy retrieves all available trading symbols using the legacy model
	GetAllSymbolsLegacy(ctx context.Context) ([]*market.Symbol, error)

	// GetSymbolInfoLegacy retrieves detailed information about a specific symbol using the legacy model
	GetSymbolInfoLegacy(ctx context.Context, symbol string) (*market.Symbol, error)

	// GetHistoricalPricesLegacy retrieves historical prices for a symbol using the legacy model
	GetHistoricalPricesLegacy(ctx context.Context, symbol string, from, to time.Time, interval string) ([]*market.Candle, error)
}
