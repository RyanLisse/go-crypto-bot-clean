package interfaces

import (
	"context"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// PositionService defines the interface for position management
type PositionService interface {
	// Entry methods
	EnterPosition(ctx context.Context, order *models.Order) (*models.Position, error)
	ScalePosition(ctx context.Context, positionID string, order *models.Order) (*models.Position, error)

	// Exit methods
	ExitPosition(ctx context.Context, positionID string, price float64) error

	// Management methods
	GetPosition(ctx context.Context, positionID string) (*models.Position, error)
	GetPositions(ctx context.Context, filter PositionFilter) ([]*models.Position, error)
	UpdateStopLoss(ctx context.Context, positionID string, price float64) error
	UpdateTakeProfit(ctx context.Context, positionID string, price float64) error
	UpdateTrailingStop(ctx context.Context, positionID string, offset float64) error

	// Monitoring methods
	CheckPositions(ctx context.Context) error
}
