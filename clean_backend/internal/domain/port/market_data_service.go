package port

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// MarketDataService defines the interface for fetching market data.
type MarketDataService interface {
	// GetTicker fetches the latest ticker information for a symbol.
	GetTicker(ctx context.Context, symbol string) (*model.Ticker, error)
	// GetTickers fetches the latest ticker information for multiple symbols.
	GetTickers(ctx context.Context, symbols []string) ([]*model.Ticker, error)
	// GetKlines fetches candlestick data for a symbol.
	GetKlines(ctx context.Context, symbol, interval string, startTime, endTime time.Time, limit int) ([]*model.Kline, error)
	// GetSymbolInfo fetches information about a specific trading symbol.
	GetSymbolInfo(ctx context.Context, symbol string) (*model.SymbolInfo, error)
	// GetAllSymbols fetches information about all available trading symbols.
	GetAllSymbols(ctx context.Context) ([]*model.SymbolInfo, error)
}

// CandlestickProvider defines the interface for fetching candlestick data.
type CandlestickProvider interface {
	// ... existing code ...
}

// TickerProvider defines the interface for fetching ticker data.
type TickerProvider interface {
	// ... existing code ...
}

// SymbolProvider defines the interface for fetching symbol information.
type SymbolProvider interface {
	// ... existing code ...
}
