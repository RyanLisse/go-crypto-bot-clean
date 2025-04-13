package port

import (
	"context"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model/market"
)

// MarketDataUseCaseInterface defines the interface for market data use cases
// This is used to allow for easier testing with mocks
type MarketDataUseCaseInterface interface {
	GetTicker(ctx context.Context, exchange, symbol string) (*market.Ticker, error)
	GetLatestTickers(ctx context.Context) ([]market.Ticker, error)
	GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, error)
	GetCandles(ctx context.Context, exchange, symbol string, interval market.Interval, start, end time.Time, limit int) ([]market.Candle, error)
	GetSymbols(ctx context.Context) ([]*market.Symbol, error)
	GetSymbol(ctx context.Context, symbol string) (*market.Symbol, error)
	GetOrderBook(ctx context.Context, exchange, symbol string) (*market.OrderBook, error)
}
