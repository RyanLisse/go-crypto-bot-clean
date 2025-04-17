package port

import (
   "context"
   "time"

   "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// MarketDataUseCaseInterface defines the interface for market data use cases
// This is used to allow for easier testing with mocks
type MarketDataUseCaseInterface interface {
   // GetTicker retrieves the current ticker for a symbol on a given exchange
   GetTicker(ctx context.Context, exchange, symbol string) (*model.Ticker, error)

   // GetLatestTickers retrieves the latest tickers across exchanges
   GetLatestTickers(ctx context.Context) ([]*model.Ticker, error)

   // GetAllTickers retrieves all tickers for an exchange
   GetAllTickers(ctx context.Context, exchange string) ([]*model.Ticker, error)

   // GetCandles retrieves historical candlestick data
   GetCandles(ctx context.Context, exchange, symbol string, interval model.KlineInterval, start, end time.Time, limit int) ([]*model.Kline, error)

   // GetSymbols retrieves all available trading symbols
   GetSymbols(ctx context.Context) ([]*model.Symbol, error)

   // GetSymbol retrieves detailed information about a specific symbol
   GetSymbol(ctx context.Context, symbol string) (*model.Symbol, error)

   // GetOrderBook retrieves the current order book for a symbol on a given exchange
   GetOrderBook(ctx context.Context, exchange, symbol string) (*model.OrderBook, error)
}
