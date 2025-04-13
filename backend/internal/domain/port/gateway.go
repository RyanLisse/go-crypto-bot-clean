package port

import (
	"context"

	"github.com/neo/crypto-bot/internal/domain/model"
)

// MexcAPI defines the interface for interacting with the MEXC cryptocurrency exchange API
type MexcAPI interface {
	// GetAccount retrieves account information from MEXC
	GetAccount(ctx context.Context) (*model.Wallet, error)

	// GetMarketData retrieves market data for a symbol
	GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error)

	// GetKlines retrieves kline (candlestick) data for a symbol
	GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error)

	// GetOrderBook retrieves the order book for a symbol
	GetOrderBook(ctx context.Context, symbol string, limit int) (*model.OrderBook, error)

	// PlaceOrder places a new order on the exchange
	PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error)

	// CancelOrder cancels an existing order
	CancelOrder(ctx context.Context, symbol string, orderID string) error

	// GetOrderStatus checks the status of an order
	GetOrderStatus(ctx context.Context, symbol string, orderID string) (*model.Order, error)
}
