package ports

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// TradeRepository defines the interface for trade storage operations
type TradeRepository interface {
	// Store persists a trade in the repository
	Store(ctx context.Context, trade *models.Trade) error

	// GetByID retrieves a trade by its ID
	GetByID(ctx context.Context, id string) (*models.Trade, error)

	// GetBySymbol retrieves all trades for a given symbol
	GetBySymbol(ctx context.Context, symbol string, limit int) ([]*models.Trade, error)

	// GetByTimeRange retrieves trades within a specific time range
	GetByTimeRange(ctx context.Context, symbol string, start, end time.Time, limit int) ([]*models.Trade, error)

	// GetByExchange retrieves trades from a specific exchange
	GetByExchange(ctx context.Context, exchange string, limit int) ([]*models.Trade, error)

	// GetByOrderID retrieves trades associated with a specific order
	GetByOrderID(ctx context.Context, orderID string) ([]*models.Trade, error)

	// Delete removes a trade from the repository
	Delete(ctx context.Context, id string) error

	// DeleteOlderThan removes trades older than the specified time
	DeleteOlderThan(ctx context.Context, before time.Time) error
}
