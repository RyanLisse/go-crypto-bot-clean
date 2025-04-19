package gateway

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/shopspring/decimal" // Assuming usage for precision
	// Implied import for the new GetExchangeInfo method
)

// MEXCGateway defines the interface for interacting with the MEXC exchange
// This acts as an adapter between the domain and the specific MEXC client implementation
type MEXCGateway interface {
	// GetSymbols retrieves all available trading symbols
	GetSymbols(ctx context.Context) ([]model.SymbolInfo, error)

	// Market Data
	GetTicker(ctx context.Context, symbol string) (model.Ticker, error)
	GetOrderBook(ctx context.Context, symbol string, depth int) (model.OrderBook, error)
	// TODO: Add other market data methods (Kline, Trades, etc.) as needed

	// Account Data
	GetAccountInfo(ctx context.Context) (model.AccountInfo, error)              // Requires authentication
	GetAssetBalance(ctx context.Context, asset string) (decimal.Decimal, error) // Requires authentication
	// TODO: Add other account methods (Open Orders, Order History, etc.) as needed

	// Trading
	PlaceOrder(ctx context.Context, order model.OrderRequest) (model.OrderResponse, error)
	CancelOrder(ctx context.Context, symbol, orderID string) error // Requires authentication
	// TODO: Add other trading methods (Batch Orders, etc.) as needed

	// WebSocket (Optional - could be separate interface if complex)
	// SubscribeTicker(ctx context.Context, symbol string, handler func(model.Ticker)) error
	// SubscribeOrderBook(ctx context.Context, symbol string, depth int, handler func(model.OrderBook)) error
	// Unsubscribe(ctx context.Context, streamName string) error

	// GetExchangeInfo retrieves general exchange information
	GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error)
}
