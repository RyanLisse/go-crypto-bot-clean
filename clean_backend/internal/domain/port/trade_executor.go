package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// TradeExecutor defines the interface for executing trades on an exchange.
type TradeExecutor interface {
	// ExecuteOrder places a new order on the exchange.
	ExecuteOrder(ctx context.Context, order *model.OrderRequest) (*model.OrderResponse, error)

	// CancelOrder attempts to cancel an existing order.
	CancelOrder(ctx context.Context, symbol, orderID string) error

	// CancelOrderWithRetry attempts to cancel an existing order with retries.
	CancelOrderWithRetry(ctx context.Context, symbol, orderID string) error

	// GetOrderStatus retrieves the status of an existing order.
	GetOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error)

	// GetOrderStatusWithRetry retrieves the status of an existing order with retries.
	GetOrderStatusWithRetry(ctx context.Context, symbol, orderID string) (*model.Order, error)
}
