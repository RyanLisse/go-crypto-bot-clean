package controls

import (
	"context"
	"fmt"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
)

// DrawdownControl evaluates if positions have experienced excessive drawdown
type DrawdownControl struct {
	BaseRiskControl
	marketDataService port.MarketDataService
	positionRepo      port.PositionRepository
}

// NewDrawdownControl creates a new drawdown risk control
func NewDrawdownControl(marketDataService port.MarketDataService, positionRepo port.PositionRepository) *DrawdownControl {
	return &DrawdownControl{
		BaseRiskControl:   NewBaseRiskControl(model.RiskTypeDrawdown, "Position Drawdown"),
		marketDataService: marketDataService,
		positionRepo:      positionRepo,
	}
}

// Evaluate checks if any positions have experienced a drawdown beyond the threshold
func (c *DrawdownControl) Evaluate(ctx context.Context, userID string, profile *model.RiskProfile) ([]*model.RiskAssessment, error) {
	// Get all active positions for the user
	positions, err := c.positionRepo.GetByUserID(ctx, userID, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	// If there are no positions, no drawdown risk
	if len(positions) == 0 {
		return nil, nil
	}

	var assessments []*model.RiskAssessment

	// Check each position for drawdown
	for _, position := range positions {
		// Get current market price
		ticker, err := c.marketDataService.GetTicker(ctx, position.Symbol)
		if err != nil {
			continue // Skip if we can't get market data
		}

		// Calculate current position value
		currentValue := position.Quantity * ticker.Price

		// Calculate entry position value
		entryValue := position.Quantity * position.EntryPrice

		// Skip positions that are in profit
		if currentValue >= entryValue {
			continue
		}

		// Calculate drawdown percentage
		drawdownPct := ((entryValue - currentValue) / entryValue) * 100

		// Check if drawdown exceeds threshold - using MaxDrawdown field which is stored as a decimal (0-1)
		// So multiply by 100 to get percentage for comparison
		maxDrawdownPct := profile.MaxDrawdown * 100
		if drawdownPct > maxDrawdownPct {
			riskLevel := determineDrawdownRiskLevel(drawdownPct, maxDrawdownPct)

			assessment := createRiskAssessment(
				userID,
				model.RiskTypeDrawdown,
				riskLevel,
				fmt.Sprintf("Position %s has experienced %.2f%% drawdown, exceeding threshold of %.2f%%",
					position.Symbol, drawdownPct, maxDrawdownPct),
			)
			assessment.Symbol = position.Symbol
			assessment.PositionID = position.ID

			// Provide recommendations based on severity
			if riskLevel == model.RiskLevelCritical {
				assessment.Recommendation = "Immediate action required: Consider closing position to limit further losses"
			} else if riskLevel == model.RiskLevelHigh {
				assessment.Recommendation = "Consider setting a tighter stop loss or reducing position size"
			} else {
				assessment.Recommendation = "Monitor position closely and consider implementing a stop loss"
			}

			assessments = append(assessments, assessment)
		}
	}

	return assessments, nil
}

// determineDrawdownRiskLevel calculates the appropriate risk level based on how much
// the drawdown exceeds the threshold
func determineDrawdownRiskLevel(drawdown, threshold float64) model.RiskLevel {
	// Calculate the ratio of drawdown to threshold
	ratio := drawdown / threshold

	if ratio >= 2.0 {
		return model.RiskLevelCritical
	} else if ratio >= 1.5 {
		return model.RiskLevelHigh
	} else {
		return model.RiskLevelMedium
	}
}
