package portfolio

import (
	"context"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// Service defines the interface for portfolio management
type Service interface {
	// GetSummary returns a summary of the portfolio
	GetSummary(ctx context.Context, userID int) (*models.PortfolioSummary, error)
	
	// GetPositions returns all positions in the portfolio
	GetPositions(ctx context.Context, userID int) ([]models.Position, error)
	
	// GetActiveTrades returns all active trades
	GetActiveTrades(ctx context.Context, userID int) ([]*models.BoughtCoin, error)
	
	// GetTradePerformance returns performance metrics for trades
	GetTradePerformance(ctx context.Context, userID int, timeRange string) (*models.PerformanceMetrics, error)
	
	// GetPortfolioValue returns the total value of the portfolio
	GetPortfolioValue(ctx context.Context, userID int) (float64, error)
}
