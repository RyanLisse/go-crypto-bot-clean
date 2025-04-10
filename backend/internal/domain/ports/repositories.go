package ports

import (
	"context"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id string) (*models.Order, error)
	List(ctx context.Context, symbol string, status models.OrderStatus) ([]*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
	Delete(ctx context.Context, id string) error
}

// PositionRepository defines the interface for position data operations
type PositionRepository interface {
	Create(ctx context.Context, position *models.Position) error
	GetByID(ctx context.Context, id string) (*models.Position, error)
	List(ctx context.Context, status models.PositionStatus) ([]*models.Position, error)
	Update(ctx context.Context, position *models.Position) error
	Delete(ctx context.Context, id string) error
	GetOpenPositionBySymbol(ctx context.Context, symbol string) (*models.Position, error)
}
