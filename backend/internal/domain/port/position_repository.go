package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// PositionRepository defines the interface for position data persistence
type PositionRepository interface {
	// Create creates a new position
	Create(ctx context.Context, position *model.Position) error

	// Update updates an existing position
	Update(ctx context.Context, position *model.Position) error

	// GetByID retrieves a position by its ID
	GetByID(ctx context.Context, id string) (*model.Position, error)

	// GetByUserID retrieves positions for a specific user with pagination
	GetByUserID(ctx context.Context, userID string, page, limit int) ([]*model.Position, error)

	// GetOpenPositionsByUserID retrieves all open positions for a specific user
	GetOpenPositionsByUserID(ctx context.Context, userID string) ([]*model.Position, error)

	// GetBySymbol retrieves positions for a specific symbol with pagination
	GetBySymbol(ctx context.Context, symbol string, page, limit int) ([]*model.Position, error)

	// GetBySymbolAndUser retrieves positions for a specific symbol and user with pagination
	GetBySymbolAndUser(ctx context.Context, symbol, userID string, page, limit int) ([]*model.Position, error)

	// GetActiveByUser retrieves all active positions for a specific user
	GetActiveByUser(ctx context.Context, userID string) ([]*model.Position, error)

	// Delete deletes a position
	Delete(ctx context.Context, id string) error

	// Count returns the number of positions matching the provided filters
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
}
