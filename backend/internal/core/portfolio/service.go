package portfolio

import (
	"context"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

type PortfolioService interface {
	GetWallet(ctx context.Context) (*models.Wallet, error)
	GetPositions(ctx context.Context) ([]models.Position, error)
	GetTradePerformance(ctx context.Context, timeRange string) (*models.PerformanceMetrics, error)
}
