package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// TradeExecutor defines the interface for executing trades with error handling and rate limiting
type TradeExecutor interface {
	// ExecuteOrder places an order with error handling and rate limiting
	ExecuteOrder(ctx context.Context, request *model.OrderRequest) (*model.OrderResponse, error)
	
	// CancelOrderWithRetry attempts to cancel an order with retries
	CancelOrderWithRetry(ctx context.Context, symbol, orderID string) error
	
	// GetOrderStatusWithRetry attempts to get order status with retries
	GetOrderStatusWithRetry(ctx context.Context, symbol, orderID string) (*model.Order, error)
}
