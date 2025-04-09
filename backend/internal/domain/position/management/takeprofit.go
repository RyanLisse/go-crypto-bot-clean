package management

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// TakeProfitStrategy defines how take profits should be managed
type TakeProfitStrategy struct {
	// InitialTakeProfitPercent is the initial take profit percentage from entry price
	InitialTakeProfitPercent float64

	// PartialTakeProfitLevels defines levels at which to take partial profits
	// Each level is a percentage of profit and the percentage of position to close
	PartialTakeProfitLevels []PartialTakeProfitLevel
}

// PartialTakeProfitLevel defines a level at which to take partial profits
type PartialTakeProfitLevel struct {
	// ProfitPercent is the profit percentage at which to take partial profits
	ProfitPercent float64

	// PositionClosePercent is the percentage of the position to close
	PositionClosePercent float64
}

// DefaultTakeProfitStrategy returns a default take profit strategy
func DefaultTakeProfitStrategy() TakeProfitStrategy {
	return TakeProfitStrategy{
		InitialTakeProfitPercent: 15.0,
		PartialTakeProfitLevels: []PartialTakeProfitLevel{
			{ProfitPercent: 5.0, PositionClosePercent: 25.0},
			{ProfitPercent: 10.0, PositionClosePercent: 25.0},
			{ProfitPercent: 20.0, PositionClosePercent: 25.0},
		},
	}
}

// ApplyTakeProfitStrategy applies a take profit strategy to a position
func (pm *PositionManager) ApplyTakeProfitStrategy(
	ctx context.Context,
	positionID string,
	strategy TakeProfitStrategy,
) error {
	// Get the existing position
	position, err := pm.positionRepo.FindByID(ctx, positionID)
	if err != nil {
		return fmt.Errorf("failed to find position: %w", err)
	}

	// Validate position
	if position.Status != "OPEN" {
		return fmt.Errorf("cannot apply take profit strategy to closed position")
	}

	// Get current price
	currentPrice, err := pm.priceService.GetPrice(ctx, position.Symbol)
	if err != nil {
		return fmt.Errorf("failed to get current price: %w", err)
	}

	// Calculate profit percentage
	profitPercent := (currentPrice - position.EntryPrice) / position.EntryPrice * 100

	// Set initial take profit if not already set
	if position.TakeProfit == 0 {
		takeProfitPrice := position.EntryPrice * (1 + strategy.InitialTakeProfitPercent/100)
		if err := pm.UpdateTakeProfit(ctx, positionID, takeProfitPrice); err != nil {
			return fmt.Errorf("failed to set initial take profit: %w", err)
		}
		pm.logger.Info("Initial take profit set",
			zap.String("position_id", position.ID),
			zap.String("symbol", position.Symbol),
			zap.Float64("take_profit", takeProfitPrice),
			zap.Float64("take_profit_percent", strategy.InitialTakeProfitPercent),
		)
	}

	// Check for partial take profit levels
	for _, level := range strategy.PartialTakeProfitLevels {
		if profitPercent >= level.ProfitPercent {
			// Calculate quantity to close
			quantityToClose := position.Quantity * (level.PositionClosePercent / 100)

			// Create sell order for partial take profit
			sellOrder := &models.Order{
				Symbol:   position.Symbol,
				Side:     "SELL",
				Type:     "MARKET",
				Quantity: quantityToClose,
			}

			// Execute the sell order
			executedOrder, err := pm.orderService.ExecuteOrder(ctx, sellOrder)
			if err != nil {
				return fmt.Errorf("failed to execute partial take profit: %w", err)
			}

			// Update position
			position.Quantity -= quantityToClose
			position.Amount = position.Quantity // For backward compatibility

			// Add order to position
			if err := pm.positionRepo.AddOrder(ctx, positionID, executedOrder); err != nil {
				return fmt.Errorf("failed to add order to position: %w", err)
			}

			// Update position in repository
			if err := pm.positionRepo.Update(ctx, position); err != nil {
				return fmt.Errorf("failed to update position: %w", err)
			}

			pm.logger.Info("Partial take profit executed",
				zap.String("position_id", position.ID),
				zap.String("symbol", position.Symbol),
				zap.Float64("profit_percent", profitPercent),
				zap.Float64("position_close_percent", level.PositionClosePercent),
				zap.Float64("quantity_closed", quantityToClose),
				zap.Float64("remaining_quantity", position.Quantity),
			)

			// If position is fully closed, mark it as closed
			if position.Quantity <= 0 {
				position.Status = "CLOSED"
				if err := pm.positionRepo.Update(ctx, position); err != nil {
					return fmt.Errorf("failed to update position status: %w", err)
				}
				pm.logger.Info("Position fully closed by partial take profits",
					zap.String("position_id", position.ID),
					zap.String("symbol", position.Symbol),
				)
				break
			}
		}
	}

	return nil
}
