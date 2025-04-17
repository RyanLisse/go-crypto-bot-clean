package controls

import (
	"context"
	"fmt"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
)

// PositionSizeControl evaluates if a position is too large relative to account size
type PositionSizeControl struct {
	BaseRiskControl
	marketDataService port.MarketDataService
	orderRepo         port.OrderRepository
}

// NewPositionSizeControl creates a new position size risk control
func NewPositionSizeControl(marketDataService port.MarketDataService, orderRepo port.OrderRepository) *PositionSizeControl {
	return &PositionSizeControl{
		BaseRiskControl:   NewBaseRiskControl(model.RiskTypePosition, "Position Size"),
		marketDataService: marketDataService,
		orderRepo:         orderRepo,
	}
}

// Evaluate checks if any orders or positions exceed the maximum position size in the risk profile
func (c *PositionSizeControl) Evaluate(ctx context.Context, userID string, profile *model.RiskProfile) ([]*model.RiskAssessment, error) {
	// Get all orders for the user
	orders, err := c.orderRepo.GetByUserID(ctx, userID, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	var assessments []*model.RiskAssessment

	// Check each order
	for _, order := range orders {
		// Skip filled orders
		if order.Status == model.OrderStatusFilled || order.Status == model.OrderStatusCanceled {
			continue
		}

		// Get current market price
		ticker, err := c.marketDataService.GetTicker(ctx, order.Symbol)
		if err != nil {
			continue // Skip if we can't get market data
		}

		// Calculate position value
		orderValue := order.Quantity * ticker.LastPrice

		// Check if order value exceeds maximum position size
		if orderValue > profile.MaxPositionSize {
			assessment := model.NewRiskAssessment(
				userID,
				model.RiskTypePosition,
				model.RiskLevelHigh,
				fmt.Sprintf("Order value %.2f exceeds maximum position size %.2f", orderValue, profile.MaxPositionSize),
			)
			assessment.Symbol = order.Symbol
			assessment.OrderID = order.ID
			assessment.Recommendation = fmt.Sprintf("Reduce order size to below %.2f or adjust risk profile", profile.MaxPositionSize)
			assessments = append(assessments, assessment)
		}
	}

	return assessments, nil
}
