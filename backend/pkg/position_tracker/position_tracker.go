package position_tracker

import (
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"go.uber.org/zap"
)

// PositionTracker manages positions for backtesting
type PositionTracker struct {
	openPositions   map[string]*models.Position // Symbol -> Position
	closedPositions []*models.ClosedPosition
	logger          *zap.Logger
}

// NewPositionTracker creates a new position tracker
func NewPositionTracker(logger *zap.Logger) *PositionTracker {
	return &PositionTracker{
		openPositions:   make(map[string]*models.Position),
		closedPositions: make([]*models.ClosedPosition, 0),
		logger:          logger,
	}
}

// OpenPosition opens a new position
func (pt *PositionTracker) OpenPosition(symbol string, side string, price float64, quantity float64, timestamp time.Time) (*models.Position, error) {
	if _, exists := pt.openPositions[symbol]; exists {
		return nil, fmt.Errorf("position already exists for symbol %s", symbol)
	}

	position := &models.Position{
		Symbol:     symbol,
		Side:       models.OrderSide(side),
		Quantity:   quantity,
		EntryPrice: price,
		OpenTime:   timestamp,
	}

	pt.openPositions[symbol] = position
	return position, nil
}

// ClosePositionBySymbol closes a position by symbol
func (pt *PositionTracker) ClosePositionBySymbol(symbol string, price float64, quantity float64, timestamp time.Time) (*models.ClosedPosition, error) {
	position, exists := pt.openPositions[symbol]
	if !exists {
		return nil, fmt.Errorf("no open position found for symbol %s", symbol)
	}

	if quantity > position.Quantity {
		return nil, fmt.Errorf("close quantity %.8f exceeds position quantity %.8f", quantity, position.Quantity)
	}

	// Calculate PnL
	var pnl float64
	if position.Side == models.OrderSideBuy {
		pnl = (price - position.EntryPrice) * quantity
	} else {
		pnl = (position.EntryPrice - price) * quantity
	}
	closedPosition := &models.ClosedPosition{
		Symbol:     symbol,
		Side:       position.Side,
		Quantity:   quantity,
		EntryPrice: position.EntryPrice,
		ExitPrice:  price,
		OpenTime:   position.OpenTime,
		CloseTime:  timestamp,
		ProfitLoss: pnl,
		Profit:     pnl, // Set both ProfitLoss and Profit for backward compatibility
	}

	if quantity == position.Quantity {
		delete(pt.openPositions, symbol)
	} else {
		position.Quantity -= quantity
	}

	pt.closedPositions = append(pt.closedPositions, closedPosition)
	return closedPosition, nil
}

// GetOpenPositions returns all currently open positions
func (pt *PositionTracker) GetOpenPositions() []*models.Position {
	positions := make([]*models.Position, 0, len(pt.openPositions))
	for _, pos := range pt.openPositions {
		positions = append(positions, pos)
	}
	return positions
}

// GetClosedPositions returns all closed positions
func (pt *PositionTracker) GetClosedPositions() []*models.ClosedPosition {
	return pt.closedPositions
}
