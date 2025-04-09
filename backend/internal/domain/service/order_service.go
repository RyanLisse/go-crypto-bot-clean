package service

import (
	"context"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// OrderService defines the interface for order execution
type OrderService interface {
	// ExecuteOrder executes a trade order
	ExecuteOrder(ctx context.Context, order *models.Order) (*models.Order, error)

	// CancelOrder cancels an existing order
	CancelOrder(ctx context.Context, orderID string) error

	// GetOrderStatus retrieves the status of an order
	GetOrderStatus(ctx context.Context, orderID string) (*models.Order, error)

	// GetOpenOrders retrieves all open orders
	GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error)
}
