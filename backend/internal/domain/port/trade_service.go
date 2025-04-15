package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// TradeService defines the interface for trading operations
type TradeService interface {
	// PlaceOrder creates and submits a new order to the exchange
	PlaceOrder(ctx context.Context, request *model.OrderRequest) (*model.OrderResponse, error)

	// CancelOrder cancels an existing order
	CancelOrder(ctx context.Context, symbol, orderID string) error

	// GetOrderStatus retrieves the current status of an order
	GetOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error)

	// GetOpenOrders retrieves all open orders
	GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error)

	// GetOrderHistory retrieves historical orders for a symbol
	GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error)

	// CalculateRequiredQuantity calculates the required quantity for an order based on amount
	CalculateRequiredQuantity(ctx context.Context, symbol string, side model.OrderSide, amount float64) (float64, error)
}
