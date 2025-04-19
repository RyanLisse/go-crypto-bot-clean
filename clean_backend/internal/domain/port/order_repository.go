package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// OrderRepository defines methods for order persistence operations.
type OrderRepository interface {
	// Create persists a new order.
	Create(ctx context.Context, order *model.Order) error
	// GetByID retrieves an order by its unique ID.
	GetByID(ctx context.Context, id string) (*model.Order, error)
	// GetByClientOrderID retrieves an order by client-provided ID.
	GetByClientOrderID(ctx context.Context, clientOrderID string) (*model.Order, error)
	// Update modifies an existing order.
	Update(ctx context.Context, order *model.Order) error
	// GetBySymbol lists orders for a symbol with pagination.
	GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error)
	// GetByUserID lists orders for a user with pagination.
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Order, error)
	// GetByStatus lists orders by status.
	GetByStatus(ctx context.Context, status model.OrderStatus, limit, offset int) ([]*model.Order, error)
	// Count returns the number of orders matching filters.
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	// Delete removes an order by its ID.
	Delete(ctx context.Context, id string) error
	// GetByOrderID retrieves an order by its exchange-specific order ID.
	GetByOrderID(ctx context.Context, orderID string) (*model.Order, error)
}
