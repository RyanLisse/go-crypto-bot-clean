package backtest

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// PositionTracker defines the interface for tracking positions during a backtest
type PositionTracker interface {
	// OpenPosition opens a new position
	OpenPosition(symbol string, side string, entryPrice float64, quantity float64, timestamp time.Time) (*models.Position, error)

	// ClosePosition closes an existing position
	ClosePosition(positionID string, exitPrice float64, timestamp time.Time) (*models.ClosedPosition, error)

	// UpdatePosition updates an existing position (e.g., for partial closes)
	UpdatePosition(positionID string, newQuantity float64, timestamp time.Time) (*models.Position, error)

	// GetOpenPositions returns all currently open positions
	GetOpenPositions() []*models.Position

	// GetClosedPositions returns all closed positions
	GetClosedPositions() []*models.ClosedPosition

	// CalculateUnrealizedPnL calculates the unrealized P&L for all open positions
	CalculateUnrealizedPnL(currentPrices map[string]float64) (float64, error)
}

// DefaultPositionTracker implements the PositionTracker interface
type DefaultPositionTracker struct {
	openPositions   map[string]*models.Position
	closedPositions []*models.ClosedPosition
}

// NewPositionTracker creates a new DefaultPositionTracker
func NewPositionTracker() *DefaultPositionTracker {
	return &DefaultPositionTracker{
		openPositions:   make(map[string]*models.Position),
		closedPositions: make([]*models.ClosedPosition, 0),
	}
}

// OpenPosition opens a new position
func (t *DefaultPositionTracker) OpenPosition(symbol string, side string, entryPrice float64, quantity float64, timestamp time.Time) (*models.Position, error) {
	positionID := uuid.New().String()

	// Convert string side to OrderSide type
	var orderSide models.OrderSide
	if side == "BUY" || side == "buy" {
		orderSide = models.OrderSideBuy
	} else if side == "SELL" || side == "sell" {
		orderSide = models.OrderSideSell
	} else {
		// Default to buy if invalid
		orderSide = models.OrderSideBuy
	}

	position := &models.Position{
		ID:         positionID,
		Symbol:     symbol,
		Side:       orderSide,
		EntryPrice: entryPrice,
		Quantity:   quantity,
		OpenTime:   timestamp,
	}

	t.openPositions[positionID] = position
	return position, nil
}

// ClosePosition closes an existing position
func (t *DefaultPositionTracker) ClosePosition(positionID string, exitPrice float64, timestamp time.Time) (*models.ClosedPosition, error) {
	position, ok := t.openPositions[positionID]
	if !ok {
		return nil, fmt.Errorf("position not found: %s", positionID)
	}

	// Calculate profit/loss
	var profit float64
	if position.Side == models.OrderSideBuy {
		profit = (exitPrice - position.EntryPrice) * position.Quantity
	} else {
		profit = (position.EntryPrice - exitPrice) * position.Quantity
	}

	// Create closed position
	closedPosition := &models.ClosedPosition{
		ID:                   position.ID,
		Symbol:               position.Symbol,
		Side:                 position.Side,
		Quantity:             position.Quantity,
		Amount:               position.Quantity, // Use Quantity for Amount
		EntryPrice:           position.EntryPrice,
		ExitPrice:            exitPrice,
		OpenTime:             position.OpenTime,
		CloseTime:            timestamp,
		ProfitLoss:           profit,
		Profit:               profit, // Set both ProfitLoss and Profit
		ProfitLossPercentage: (profit / (position.EntryPrice * position.Quantity)) * 100,
		ExitReason:           "backtest",
	}

	// Remove from open positions and add to closed positions
	delete(t.openPositions, positionID)
	t.closedPositions = append(t.closedPositions, closedPosition)

	return closedPosition, nil
}

// UpdatePosition updates an existing position (e.g., for partial closes)
func (t *DefaultPositionTracker) UpdatePosition(positionID string, newQuantity float64, timestamp time.Time) (*models.Position, error) {
	position, ok := t.openPositions[positionID]
	if !ok {
		return nil, fmt.Errorf("position not found: %s", positionID)
	}

	// Update quantity
	position.Quantity = newQuantity

	return position, nil
}

// GetOpenPositions returns all currently open positions
func (t *DefaultPositionTracker) GetOpenPositions() []*models.Position {
	positions := make([]*models.Position, 0, len(t.openPositions))
	for _, position := range t.openPositions {
		positions = append(positions, position)
	}
	return positions
}

// GetClosedPositions returns all closed positions
func (t *DefaultPositionTracker) GetClosedPositions() []*models.ClosedPosition {
	return t.closedPositions
}

// CalculateUnrealizedPnL calculates the unrealized P&L for all open positions
func (t *DefaultPositionTracker) CalculateUnrealizedPnL(currentPrices map[string]float64) (float64, error) {
	var totalPnL float64

	for _, position := range t.openPositions {
		currentPrice, ok := currentPrices[position.Symbol]
		if !ok {
			return 0, fmt.Errorf("current price not found for symbol: %s", position.Symbol)
		}

		var positionPnL float64
		if position.Side == "BUY" {
			positionPnL = (currentPrice - position.EntryPrice) * position.Quantity
		} else {
			positionPnL = (position.EntryPrice - currentPrice) * position.Quantity
		}

		totalPnL += positionPnL
	}

	return totalPnL, nil
}
