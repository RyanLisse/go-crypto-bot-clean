package controls

import (
	"context"
	"fmt"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
)

// LiquidityControl evaluates if a market has sufficient liquidity for safe trading
type LiquidityControl struct {
	BaseRiskControl
	marketDataService port.MarketDataService
}

// NewLiquidityControl creates a new liquidity risk control
func NewLiquidityControl(marketDataService port.MarketDataService) *LiquidityControl {
	return &LiquidityControl{
		BaseRiskControl:   NewBaseRiskControl(model.RiskTypeLiquidity, "Market Liquidity"),
		marketDataService: marketDataService,
	}
}

// Evaluate checks if the 24-hour trading volume is below the minimum liquidity threshold
func (c *LiquidityControl) Evaluate(ctx context.Context, userID string, profile *model.RiskProfile) ([]*model.RiskAssessment, error) {
	// Get all available trading symbols
	symbols, err := c.marketDataService.GetAllSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get symbols: %w", err)
	}

	var assessments []*model.RiskAssessment

	// Check liquidity for each symbol
	for _, symbol := range symbols {
		// Skip non-trading symbols
		if symbol.Status != "TRADING" {
			continue
		}

		// Get market data
		ticker, err := c.marketDataService.GetTicker(ctx, symbol.Symbol)
		if err != nil {
			continue // Skip if we can't get market data
		}

		// Calculate USD volume
		usdVolume := ticker.Volume * ticker.Price

		// Check if volume is below minimum liquidity threshold
		if usdVolume < profile.MinLiquidity {
			assessment := model.NewRiskAssessment(
				userID,
				model.RiskTypeLiquidity,
				model.RiskLevelMedium,
				fmt.Sprintf("Low liquidity for %s: 24h volume $%.2f is below minimum threshold $%.2f",
					symbol.Symbol, usdVolume, profile.MinLiquidity),
			)
			assessment.Symbol = symbol.Symbol
			assessment.Recommendation = "Consider trading assets with higher liquidity or reducing position size"
			assessments = append(assessments, assessment)
		}
	}

	return assessments, nil
}
