package repository

import (
	"context"

	"go-crypto-bot-clean/backend/internal/models"
)

// PortfolioRepository defines the interface for portfolio operations
type PortfolioRepository interface {
	GetByUserID(ctx context.Context, userID string) (*models.Portfolio, error)
	AddPosition(ctx context.Context, position *models.Position) error
	UpdatePosition(ctx context.Context, position *models.Position) error
	GetPositionsByPortfolioID(ctx context.Context, portfolioID string) ([]*models.Position, error)
}
