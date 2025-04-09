package repository

import (
	"context"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// BoughtCoinRepository defines operations for managing bought coins
type BoughtCoinRepository interface {
	// FindAll returns all bought coins that haven't been deleted
	FindAll(ctx context.Context) ([]models.BoughtCoin, error)

	// FindByID returns a specific bought coin by ID
	FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error)

	// FindBySymbol returns a specific bought coin by symbol
	FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error)

	// Create adds a new bought coin
	Create(ctx context.Context, coin *models.BoughtCoin) (int64, error)

	// Update modifies an existing bought coin
	Update(ctx context.Context, coin *models.BoughtCoin) error

	// Delete marks a bought coin as deleted
	Delete(ctx context.Context, id int64) error
}
