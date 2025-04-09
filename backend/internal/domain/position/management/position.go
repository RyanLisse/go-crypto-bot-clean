package management

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/interfaces"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// PositionManager implements the interfaces.PositionService interface
type PositionManager struct {
	positionRepo interfaces.PositionRepository
	orderService interfaces.OrderService
	priceService interfaces.PriceService
	logger       *zap.Logger
}

// NewPositionManager creates a new position manager
func NewPositionManager(
	positionRepo interfaces.PositionRepository,
	orderService interfaces.OrderService,
	priceService interfaces.PriceService,
	logger *zap.Logger,
) *PositionManager {
	return &PositionManager{
		positionRepo: positionRepo,
		orderService: orderService,
		priceService: priceService,
		logger:       logger,
	}
}

// EnterPosition creates a new position from an order
func (pm *PositionManager) EnterPosition(ctx context.Context, order *models.Order) (*models.Position, error) {
	// Validate order
	if order.Side != "BUY" {
		return nil, fmt.Errorf("can only enter position with buy order")
	}

	// Execute the order if it's not already filled
	if order.Status != "FILLED" {
		executedOrder, err := pm.orderService.ExecuteOrder(ctx, order)
		if err != nil {
			return nil, fmt.Errorf("failed to execute order: %w", err)
		}
		order = executedOrder
	}

	// Create a new position
	position := &models.Position{
		Symbol:       order.Symbol,
		Quantity:     order.Quantity,
		Amount:       order.Quantity, // For backward compatibility
		EntryPrice:   order.Price,
		CurrentPrice: order.Price,
		OpenTime:     time.Now(),
		OpenedAt:     time.Now(), // For backward compatibility
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Status:       "OPEN",
		Orders:       []models.Order{*order},
	}

	// Calculate default stop loss and take profit
	position.StopLoss = order.Price * 0.95  // 5% stop loss
	position.TakeProfit = order.Price * 1.1 // 10% take profit

	// Save the position
	id, err := pm.positionRepo.Create(ctx, position)
	if err != nil {
		return nil, fmt.Errorf("failed to create position: %w", err)
	}
	position.ID = id

	pm.logger.Info("Position entered",
		zap.String("position_id", position.ID),
		zap.String("symbol", position.Symbol),
		zap.Float64("quantity", position.Quantity),
		zap.Float64("entry_price", position.EntryPrice),
	)

	return position, nil
}

// ScalePosition adds to an existing position
func (pm *PositionManager) ScalePosition(ctx context.Context, positionID string, order *models.Order) (*models.Position, error) {
	// Get the existing position
	position, err := pm.positionRepo.FindByID(ctx, positionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find position: %w", err)
	}

	// Validate order
	if order.Symbol != position.Symbol {
		return nil, fmt.Errorf("order symbol %s does not match position symbol %s", order.Symbol, position.Symbol)
	}
	if order.Side != "BUY" {
		return nil, fmt.Errorf("can only scale position with buy order")
	}
	if position.Status != "OPEN" {
		return nil, fmt.Errorf("cannot scale closed position")
	}

	// Execute the order if it's not already filled
	if order.Status != "FILLED" {
		executedOrder, err := pm.orderService.ExecuteOrder(ctx, order)
		if err != nil {
			return nil, fmt.Errorf("failed to execute order: %w", err)
		}
		order = executedOrder
	}

	// Calculate new average entry price
	totalQuantity := position.Quantity + order.Quantity
	totalCost := (position.Quantity * position.EntryPrice) + (order.Quantity * order.Price)
	newEntryPrice := totalCost / totalQuantity

	// Update position
	position.Quantity = totalQuantity
	position.Amount = totalQuantity // For backward compatibility
	position.EntryPrice = newEntryPrice
	position.UpdatedAt = time.Now()

	// Add order to position
	if err := pm.positionRepo.AddOrder(ctx, positionID, order); err != nil {
		return nil, fmt.Errorf("failed to add order to position: %w", err)
	}

	// Update position in repository
	if err := pm.positionRepo.Update(ctx, position); err != nil {
		return nil, fmt.Errorf("failed to update position: %w", err)
	}

	pm.logger.Info("Position scaled",
		zap.String("position_id", position.ID),
		zap.String("symbol", position.Symbol),
		zap.Float64("new_quantity", position.Quantity),
		zap.Float64("new_entry_price", position.EntryPrice),
	)

	return position, nil
}

// ExitPosition closes a position at the specified price
func (pm *PositionManager) ExitPosition(ctx context.Context, positionID string, price float64) error {
	// Get the existing position
	position, err := pm.positionRepo.FindByID(ctx, positionID)
	if err != nil {
		return fmt.Errorf("failed to find position: %w", err)
	}

	// Validate position
	if position.Status != "OPEN" {
		return fmt.Errorf("position is already closed")
	}

	// Create sell order
	sellOrder := &models.Order{
		Symbol:    position.Symbol,
		Side:      "SELL",
		Type:      "MARKET",
		Quantity:  position.Quantity,
		CreatedAt: time.Now(),
	}

	// Execute the sell order
	executedOrder, err := pm.orderService.ExecuteOrder(ctx, sellOrder)
	if err != nil {
		return fmt.Errorf("failed to execute sell order: %w", err)
	}

	// Update position
	position.CurrentPrice = executedOrder.Price
	position.Status = "CLOSED"
	position.UpdatedAt = time.Now()
	position.PnL = (position.CurrentPrice - position.EntryPrice) * position.Quantity
	position.PnLPercentage = (position.CurrentPrice - position.EntryPrice) / position.EntryPrice * 100

	// Add order to position
	if err := pm.positionRepo.AddOrder(ctx, positionID, executedOrder); err != nil {
		return fmt.Errorf("failed to add order to position: %w", err)
	}

	// Update position in repository
	if err := pm.positionRepo.Update(ctx, position); err != nil {
		return fmt.Errorf("failed to update position: %w", err)
	}

	pm.logger.Info("Position exited",
		zap.String("position_id", position.ID),
		zap.String("symbol", position.Symbol),
		zap.Float64("exit_price", position.CurrentPrice),
		zap.Float64("pnl", position.PnL),
		zap.Float64("pnl_percentage", position.PnLPercentage),
	)

	return nil
}

// GetPosition retrieves a position by ID
func (pm *PositionManager) GetPosition(ctx context.Context, positionID string) (*models.Position, error) {
	position, err := pm.positionRepo.FindByID(ctx, positionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find position: %w", err)
	}

	// Update current price and PnL
	if position.Status == "OPEN" {
		currentPrice, err := pm.priceService.GetPrice(ctx, position.Symbol)
		if err != nil {
			pm.logger.Warn("Failed to get current price",
				zap.String("position_id", position.ID),
				zap.String("symbol", position.Symbol),
				zap.Error(err),
			)
		} else {
			position.CurrentPrice = currentPrice
			position.PnL = (position.CurrentPrice - position.EntryPrice) * position.Quantity
			position.PnLPercentage = (position.CurrentPrice - position.EntryPrice) / position.EntryPrice * 100
		}
	}

	return position, nil
}

// GetPositions retrieves positions based on filter criteria
func (pm *PositionManager) GetPositions(ctx context.Context, filter interfaces.PositionFilter) ([]*models.Position, error) {
	positions, err := pm.positionRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find positions: %w", err)
	}

	// Update current prices and PnL for open positions
	for _, pos := range positions {
		if pos.Status == "OPEN" {
			currentPrice, err := pm.priceService.GetPrice(ctx, pos.Symbol)
			if err != nil {
				pm.logger.Warn("Failed to get current price",
					zap.String("position_id", pos.ID),
					zap.String("symbol", pos.Symbol),
					zap.Error(err),
				)
				continue
			}
			pos.CurrentPrice = currentPrice
			pos.PnL = (pos.CurrentPrice - pos.EntryPrice) * pos.Quantity
			pos.PnLPercentage = (pos.CurrentPrice - pos.EntryPrice) / pos.EntryPrice * 100
		}
	}

	return positions, nil
}

// UpdateStopLoss updates the stop loss price for a position
func (pm *PositionManager) UpdateStopLoss(ctx context.Context, positionID string, price float64) error {
	// Get the existing position
	position, err := pm.positionRepo.FindByID(ctx, positionID)
	if err != nil {
		return fmt.Errorf("failed to find position: %w", err)
	}

	// Validate position
	if position.Status != "OPEN" {
		return fmt.Errorf("cannot update stop loss for closed position")
	}

	// Update stop loss
	position.StopLoss = price
	position.UpdatedAt = time.Now()

	// Update position in repository
	if err := pm.positionRepo.Update(ctx, position); err != nil {
		return fmt.Errorf("failed to update position: %w", err)
	}

	pm.logger.Info("Stop loss updated",
		zap.String("position_id", position.ID),
		zap.String("symbol", position.Symbol),
		zap.Float64("stop_loss", position.StopLoss),
	)

	return nil
}

// UpdateTakeProfit updates the take profit price for a position
func (pm *PositionManager) UpdateTakeProfit(ctx context.Context, positionID string, price float64) error {
	// Get the existing position
	position, err := pm.positionRepo.FindByID(ctx, positionID)
	if err != nil {
		return fmt.Errorf("failed to find position: %w", err)
	}

	// Validate position
	if position.Status != "OPEN" {
		return fmt.Errorf("cannot update take profit for closed position")
	}

	// Update take profit
	position.TakeProfit = price
	position.UpdatedAt = time.Now()

	// Update position in repository
	if err := pm.positionRepo.Update(ctx, position); err != nil {
		return fmt.Errorf("failed to update position: %w", err)
	}

	pm.logger.Info("Take profit updated",
		zap.String("position_id", position.ID),
		zap.String("symbol", position.Symbol),
		zap.Float64("take_profit", position.TakeProfit),
	)

	return nil
}

// UpdateTrailingStop updates the trailing stop offset for a position
func (pm *PositionManager) UpdateTrailingStop(ctx context.Context, positionID string, offset float64) error {
	// Get the existing position
	position, err := pm.positionRepo.FindByID(ctx, positionID)
	if err != nil {
		return fmt.Errorf("failed to find position: %w", err)
	}

	// Validate position
	if position.Status != "OPEN" {
		return fmt.Errorf("cannot update trailing stop for closed position")
	}

	// Update trailing stop
	trailingStop := offset
	position.TrailingStop = &trailingStop
	position.UpdatedAt = time.Now()

	// Update position in repository
	if err := pm.positionRepo.Update(ctx, position); err != nil {
		return fmt.Errorf("failed to update position: %w", err)
	}

	pm.logger.Info("Trailing stop updated",
		zap.String("position_id", position.ID),
		zap.String("symbol", position.Symbol),
		zap.Float64("trailing_stop", *position.TrailingStop),
	)

	return nil
}

// CheckPositions monitors all open positions and executes stop loss/take profit orders
func (pm *PositionManager) CheckPositions(ctx context.Context) error {
	// Get all open positions
	filter := interfaces.PositionFilter{
		Status: "OPEN",
	}
	positions, err := pm.positionRepo.FindAll(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to find positions: %w", err)
	}

	for _, pos := range positions {
		// Get current price
		currentPrice, err := pm.priceService.GetPrice(ctx, pos.Symbol)
		if err != nil {
			pm.logger.Warn("Failed to get current price",
				zap.String("position_id", pos.ID),
				zap.String("symbol", pos.Symbol),
				zap.Error(err),
			)
			continue
		}

		// Update position with current price
		pos.CurrentPrice = currentPrice
		pos.PnL = (pos.CurrentPrice - pos.EntryPrice) * pos.Quantity
		pos.PnLPercentage = (pos.CurrentPrice - pos.EntryPrice) / pos.EntryPrice * 100

		// Check stop loss
		if pos.CurrentPrice <= pos.StopLoss {
			pm.logger.Info("Stop loss triggered",
				zap.String("position_id", pos.ID),
				zap.String("symbol", pos.Symbol),
				zap.Float64("current_price", pos.CurrentPrice),
				zap.Float64("stop_loss", pos.StopLoss),
			)
			if err := pm.ExitPosition(ctx, pos.ID, pos.CurrentPrice); err != nil {
				pm.logger.Error("Failed to exit position on stop loss",
					zap.String("position_id", pos.ID),
					zap.String("symbol", pos.Symbol),
					zap.Error(err),
				)
			}
			continue
		}

		// Check take profit
		if pos.CurrentPrice >= pos.TakeProfit {
			pm.logger.Info("Take profit triggered",
				zap.String("position_id", pos.ID),
				zap.String("symbol", pos.Symbol),
				zap.Float64("current_price", pos.CurrentPrice),
				zap.Float64("take_profit", pos.TakeProfit),
			)
			if err := pm.ExitPosition(ctx, pos.ID, pos.CurrentPrice); err != nil {
				pm.logger.Error("Failed to exit position on take profit",
					zap.String("position_id", pos.ID),
					zap.String("symbol", pos.Symbol),
					zap.Error(err),
				)
			}
			continue
		}

		// Check trailing stop
		if pos.TrailingStop != nil {
			// Calculate trailing stop price
			trailingStopPrice := pos.CurrentPrice * (1 - *pos.TrailingStop/100)

			// If trailing stop price is higher than current stop loss, update it
			if trailingStopPrice > pos.StopLoss {
				pm.logger.Info("Trailing stop updated",
					zap.String("position_id", pos.ID),
					zap.String("symbol", pos.Symbol),
					zap.Float64("old_stop_loss", pos.StopLoss),
					zap.Float64("new_stop_loss", trailingStopPrice),
				)
				if err := pm.UpdateStopLoss(ctx, pos.ID, trailingStopPrice); err != nil {
					pm.logger.Error("Failed to update trailing stop",
						zap.String("position_id", pos.ID),
						zap.String("symbol", pos.Symbol),
						zap.Error(err),
					)
				}
			}
		}

		// Update position in repository
		if err := pm.positionRepo.Update(ctx, pos); err != nil {
			pm.logger.Error("Failed to update position",
				zap.String("position_id", pos.ID),
				zap.String("symbol", pos.Symbol),
				zap.Error(err),
			)
		}
	}

	return nil
}
