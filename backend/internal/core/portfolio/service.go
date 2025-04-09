package portfolio

import (
	"context"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

type PortfolioService interface {
	GetWallet(ctx context.Context) (*models.Wallet, error)
	GetPositions(ctx context.Context) ([]models.Position, error)
	GetTradePerformance(ctx context.Context, timeRange string) (*models.PerformanceMetrics, error)
}
