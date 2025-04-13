package port

import (
	"context"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model/market"
)

// MarketDataService defines the interface for market data operations
type MarketDataService interface {
	// GetTicker retrieves the current ticker for a symbol
	GetTicker(ctx context.Context, symbol string) (*market.Ticker, error)

	// GetCandles retrieves historical candlestick data
	GetCandles(ctx context.Context, symbol string, interval string, limit int) ([]*market.Candle, error)

	// GetOrderBook retrieves the current order book for a symbol
	GetOrderBook(ctx context.Context, symbol string, depth int) (*market.OrderBook, error)

	// GetAllSymbols retrieves all available trading symbols
	GetAllSymbols(ctx context.Context) ([]*market.Symbol, error)

	// GetSymbolInfo retrieves detailed information about a specific symbol
	GetSymbolInfo(ctx context.Context, symbol string) (*market.Symbol, error)

	// GetHistoricalPrices retrieves historical prices for a symbol
	GetHistoricalPrices(ctx context.Context, symbol string, from, to time.Time, interval string) ([]*market.Candle, error)
}
