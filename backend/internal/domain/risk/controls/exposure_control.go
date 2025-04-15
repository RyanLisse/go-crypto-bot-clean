package controls

import (
	"context"
	"fmt"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
)

// ExposureControl evaluates if a user's total market exposure is too high
type ExposureControl struct {
	BaseRiskControl
	marketDataService port.MarketDataService
	positionRepo      port.PositionRepository
}

// NewExposureControl creates a new exposure risk control
func NewExposureControl(marketDataService port.MarketDataService, positionRepo port.PositionRepository) *ExposureControl {
	return &ExposureControl{
		BaseRiskControl:   NewBaseRiskControl(model.RiskTypeExposure, "Total Exposure"),
		marketDataService: marketDataService,
		positionRepo:      positionRepo,
	}
}

// Evaluate checks if the user's total market exposure exceeds their risk profile limits
func (c *ExposureControl) Evaluate(ctx context.Context, userID string, profile *model.RiskProfile) ([]*model.RiskAssessment, error) {
	// Get all active positions for the user
	positions, err := c.positionRepo.GetByUserID(ctx, userID, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	// If there are no positions, no exposure risk
	if len(positions) == 0 {
		return nil, nil
	}

	var totalExposure float64

	// Calculate total exposure across all positions
	for _, position := range positions {
		// Get current market price
		ticker, err := c.marketDataService.GetTicker(ctx, position.Symbol)
		if err != nil {
			continue // Skip if we can't get market data
		}

		// Add position's current value to total exposure
		positionValue := position.Quantity * ticker.Price
		totalExposure += positionValue
	}

	// Determine if total exposure exceeds the maximum allowed by the risk profile
	if totalExposure > profile.MaxTotalExposure {
		assessment := createRiskAssessment(
			userID,
			model.RiskTypeExposure,
			determineExposureRiskLevel(totalExposure, profile.MaxTotalExposure),
			fmt.Sprintf("Total market exposure of $%.2f exceeds maximum of $%.2f",
				totalExposure, profile.MaxTotalExposure),
		)
		assessment.Recommendation = "Reduce total position sizes or close some positions to bring exposure within limits"
		return []*model.RiskAssessment{assessment}, nil
	}

	return nil, nil
}

// determineExposureRiskLevel calculates the appropriate risk level based on how much
// the total exposure exceeds the threshold
func determineExposureRiskLevel(exposure, threshold float64) model.RiskLevel {
	// Calculate what percentage of the threshold the exposure represents
	percentage := (exposure / threshold) * 100

	if percentage >= 150 {
		return model.RiskLevelCritical
	} else if percentage >= 125 {
		return model.RiskLevelHigh
	} else {
		return model.RiskLevelMedium
	}
}
