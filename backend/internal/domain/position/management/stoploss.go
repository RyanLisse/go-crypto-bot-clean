package management

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// StopLossStrategy defines how stop losses should be managed
type StopLossStrategy struct {
	// InitialStopLossPercent is the initial stop loss percentage from entry price
	InitialStopLossPercent float64

	// BreakEvenThresholdPercent is the profit percentage at which to move stop loss to break even
	BreakEvenThresholdPercent float64

	// TrailingStopPercent is the percentage for trailing stop (if enabled)
	TrailingStopPercent float64

	// EnableTrailingStop determines whether to use trailing stops
	EnableTrailingStop bool
}

// DefaultStopLossStrategy returns a default stop loss strategy
func DefaultStopLossStrategy() StopLossStrategy {
	return StopLossStrategy{
		InitialStopLossPercent:     5.0,
		BreakEvenThresholdPercent:  3.0,
		TrailingStopPercent:        2.0,
		EnableTrailingStop:         true,
	}
}

// ApplyStopLossStrategy applies a stop loss strategy to a position
func (pm *PositionManager) ApplyStopLossStrategy(
	ctx context.Context,
	positionID string,
	strategy StopLossStrategy,
) error {
	// Get the existing position
	position, err := pm.positionRepo.FindByID(ctx, positionID)
	if err != nil {
		return fmt.Errorf("failed to find position: %w", err)
	}

	// Validate position
	if position.Status != "OPEN" {
		return fmt.Errorf("cannot apply stop loss strategy to closed position")
	}

	// Get current price
	currentPrice, err := pm.priceService.GetPrice(ctx, position.Symbol)
	if err != nil {
		return fmt.Errorf("failed to get current price: %w", err)
	}

	// Calculate profit percentage
	profitPercent := (currentPrice - position.EntryPrice) / position.EntryPrice * 100

	// Set initial stop loss if not already set
	if position.StopLoss == 0 {
		stopLossPrice := position.EntryPrice * (1 - strategy.InitialStopLossPercent/100)
		if err := pm.UpdateStopLoss(ctx, positionID, stopLossPrice); err != nil {
			return fmt.Errorf("failed to set initial stop loss: %w", err)
		}
		pm.logger.Info("Initial stop loss set",
			zap.String("position_id", position.ID),
			zap.String("symbol", position.Symbol),
			zap.Float64("stop_loss", stopLossPrice),
			zap.Float64("stop_loss_percent", strategy.InitialStopLossPercent),
		)
	}

	// Move stop loss to break even if profit threshold reached
	if profitPercent >= strategy.BreakEvenThresholdPercent && position.StopLoss < position.EntryPrice {
		if err := pm.UpdateStopLoss(ctx, positionID, position.EntryPrice); err != nil {
			return fmt.Errorf("failed to move stop loss to break even: %w", err)
		}
		pm.logger.Info("Stop loss moved to break even",
			zap.String("position_id", position.ID),
			zap.String("symbol", position.Symbol),
			zap.Float64("entry_price", position.EntryPrice),
			zap.Float64("profit_percent", profitPercent),
		)
	}

	// Enable trailing stop if configured
	if strategy.EnableTrailingStop && position.TrailingStop == nil && profitPercent >= strategy.BreakEvenThresholdPercent {
		if err := pm.UpdateTrailingStop(ctx, positionID, strategy.TrailingStopPercent); err != nil {
			return fmt.Errorf("failed to enable trailing stop: %w", err)
		}
		pm.logger.Info("Trailing stop enabled",
			zap.String("position_id", position.ID),
			zap.String("symbol", position.Symbol),
			zap.Float64("trailing_stop_percent", strategy.TrailingStopPercent),
			zap.Float64("profit_percent", profitPercent),
		)
	}

	return nil
}
