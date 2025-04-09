package management

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// ScalingStrategy defines how positions should be scaled
type ScalingStrategy struct {
	// MaxScalingSteps is the maximum number of times a position can be scaled
	MaxScalingSteps int

	// ScalingFactor is the multiplier for each scaling step (e.g., 2.0 doubles the position size)
	ScalingFactor float64

	// PriceThresholdPercent is the price movement required to trigger scaling
	PriceThresholdPercent float64

	// MinProfitPercent is the minimum profit required before scaling
	MinProfitPercent float64
}

// DefaultScalingStrategy returns a default scaling strategy
func DefaultScalingStrategy() ScalingStrategy {
	return ScalingStrategy{
		MaxScalingSteps:       3,
		ScalingFactor:         1.5,
		PriceThresholdPercent: 5.0,
		MinProfitPercent:      2.0,
	}
}

// ScalePositionByStrategy scales a position according to the provided strategy
func (pm *PositionManager) ScalePositionByStrategy(
	ctx context.Context,
	positionID string,
	strategy ScalingStrategy,
) (*models.Position, error) {
	// Get the existing position
	position, err := pm.positionRepo.FindByID(ctx, positionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find position: %w", err)
	}

	// Validate position
	if position.Status != "OPEN" {
		return nil, fmt.Errorf("cannot scale closed position")
	}

	// Check if we've reached the maximum scaling steps
	if len(position.Orders) > strategy.MaxScalingSteps {
		return nil, fmt.Errorf("maximum scaling steps reached")
	}

	// Get current price
	currentPrice, err := pm.priceService.GetPrice(ctx, position.Symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get current price: %w", err)
	}

	// Calculate profit percentage
	profitPercent := (currentPrice - position.EntryPrice) / position.EntryPrice * 100

	// Check if profit is sufficient for scaling
	if profitPercent < strategy.MinProfitPercent {
		return nil, fmt.Errorf("insufficient profit for scaling: %.2f%% (required: %.2f%%)",
			profitPercent, strategy.MinProfitPercent)
	}

	// Check if price movement is sufficient for scaling
	priceMovementPercent := (currentPrice - position.EntryPrice) / position.EntryPrice * 100
	if priceMovementPercent < strategy.PriceThresholdPercent {
		return nil, fmt.Errorf("insufficient price movement for scaling: %.2f%% (required: %.2f%%)",
			priceMovementPercent, strategy.PriceThresholdPercent)
	}

	// Calculate scaling quantity
	scalingQuantity := position.Quantity * strategy.ScalingFactor

	// Create buy order for scaling
	order := &models.Order{
		Symbol:    position.Symbol,
		Side:      "BUY",
		Type:      "MARKET",
		Quantity:  scalingQuantity,
		CreatedAt: time.Now(),
	}

	// Execute the scaling
	scaledPosition, err := pm.ScalePosition(ctx, positionID, order)
	if err != nil {
		return nil, fmt.Errorf("failed to scale position: %w", err)
	}

	pm.logger.Info("Position scaled by strategy",
		zap.String("position_id", scaledPosition.ID),
		zap.String("symbol", scaledPosition.Symbol),
		zap.Float64("original_quantity", position.Quantity),
		zap.Float64("scaling_quantity", scalingQuantity),
		zap.Float64("new_quantity", scaledPosition.Quantity),
		zap.Float64("profit_percent", profitPercent),
	)

	return scaledPosition, nil
}
