package repository

import (
	"context"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// ClosedPositionRepository defines the interface for closed position data access
type ClosedPositionRepository interface {
	// Create adds a new closed position
	Create(ctx context.Context, position *models.ClosedPosition) (string, error)

	// FindByID returns a specific closed position by ID
	FindByID(ctx context.Context, id string) (*models.ClosedPosition, error)

	// FindAll returns all closed positions matching the filter
	FindAll(ctx context.Context, filter ClosedPositionFilter) ([]*models.ClosedPosition, error)

	// FindBySymbol returns closed positions for a specific symbol
	FindBySymbol(ctx context.Context, symbol string) ([]*models.ClosedPosition, error)

	// FindByTimeRange returns closed positions within a time range
	FindByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*models.ClosedPosition, error)

	// GetTotalPnL returns the total profit/loss for all closed positions
	GetTotalPnL(ctx context.Context, startTime, endTime time.Time) (float64, error)

	// GetWinRate returns the win rate for all closed positions
	GetWinRate(ctx context.Context, startTime, endTime time.Time) (float64, error)

	// GetProfitFactor returns the profit factor for all closed positions
	GetProfitFactor(ctx context.Context, startTime, endTime time.Time) (float64, error)
}

// ClosedPositionFilter used for filtering closed positions in queries
type ClosedPositionFilter struct {
	Symbol   string
	MinPnL   *float64
	MaxPnL   *float64
	FromDate *time.Time
	ToDate   *time.Time
	Strategy string
	Reason   string
}
