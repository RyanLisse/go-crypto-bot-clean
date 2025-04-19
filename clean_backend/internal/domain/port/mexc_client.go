package port

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// MEXCClient defines the interface for interacting with the MEXC exchange API
type MEXCClient interface {
	// GetMarketData retrieves current market data for a symbol
	GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error)

	// GetKlines retrieves kline/candlestick data for a symbol
	GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error)

	// GetOrderBook retrieves the order book for a symbol
	GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error)

	// GetSymbols retrieves all available symbols
	GetSymbols(ctx context.Context) ([]*model.Symbol, error)

	// GetSymbol retrieves information about a specific symbol
	GetSymbol(ctx context.Context, symbol string) (*model.Symbol, error)

	// GetServerTime retrieves the current server time
	GetServerTime(ctx context.Context) (time.Time, error)

	// GetExchangeInfo retrieves exchange information
	GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error)

	// PlaceOrder places an order on the exchange
	PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error)
}
