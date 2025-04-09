// Package controls provides risk control implementations for the trading bot.
package controls

import (
	"context"
	"errors"
	"fmt"
)

// PositionSizer calculates appropriate position sizes based on risk parameters
type PositionSizer struct {
	priceService PriceService
	logger       Logger
}

// PriceService defines the interface for getting price information
type PriceService interface {
	GetPrice(ctx context.Context, symbol string) (float64, error)
}

// Logger defines the interface for logging
type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// NewPositionSizer creates a new PositionSizer
func NewPositionSizer(priceService PriceService, logger Logger) *PositionSizer {
	return &PositionSizer{
		priceService: priceService,
		logger:       logger,
	}
}

// CalculatePositionSize determines a safe position size based on risk parameters
func (ps *PositionSizer) CalculatePositionSize(
	ctx context.Context,
	symbol string,
	accountBalance float64,
	riskPercent float64,
	stopLossPercent float64,
) (float64, error) {
	// Get current price
	currentPrice, err := ps.priceService.GetPrice(ctx, symbol)
	if err != nil {
		return 0, fmt.Errorf("failed to get price for %s: %w", symbol, err)
	}

	// Calculate stop-loss price based on percentage
	stopLossPrice := currentPrice * (1 - stopLossPercent/100)

	// Calculate risk amount based on account balance
	riskAmount := accountBalance * (riskPercent / 100)

	// Calculate position size
	priceDifference := currentPrice - stopLossPrice
	riskPerUnit := priceDifference

	if riskPerUnit <= 0 {
		return 0, errors.New("invalid stop-loss placement, risk per unit is zero or negative")
	}

	// Position size = risk amount / risk per unit
	positionSize := riskAmount / riskPerUnit

	// Convert to coin quantity based on price
	quantity := positionSize / currentPrice

	ps.logger.Info("Calculated position size",
		"symbol", symbol,
		"account_balance", accountBalance,
		"risk_percent", riskPercent,
		"risk_amount", riskAmount,
		"stop_loss_percent", stopLossPercent,
		"stop_loss_price", stopLossPrice,
		"quantity", quantity)

	return quantity, nil
}

// CalculateOrderValue calculates the total order value for a given quantity and price
func (ps *PositionSizer) CalculateOrderValue(
	ctx context.Context,
	symbol string,
	quantity float64,
) (float64, error) {
	// Get current price
	currentPrice, err := ps.priceService.GetPrice(ctx, symbol)
	if err != nil {
		return 0, fmt.Errorf("failed to get price for %s: %w", symbol, err)
	}

	// Calculate order value
	orderValue := quantity * currentPrice

	ps.logger.Info("Calculated order value",
		"symbol", symbol,
		"quantity", quantity,
		"price", currentPrice,
		"order_value", orderValue)

	return orderValue, nil
}

// CalculateMaxQuantity calculates the maximum quantity that can be purchased
// based on the maximum order value and current price
func (ps *PositionSizer) CalculateMaxQuantity(
	ctx context.Context,
	symbol string,
	maxOrderValue float64,
) (float64, error) {
	// Get current price
	currentPrice, err := ps.priceService.GetPrice(ctx, symbol)
	if err != nil {
		return 0, fmt.Errorf("failed to get price for %s: %w", symbol, err)
	}

	// Calculate maximum quantity
	maxQuantity := maxOrderValue / currentPrice

	ps.logger.Info("Calculated maximum quantity",
		"symbol", symbol,
		"max_order_value", maxOrderValue,
		"price", currentPrice,
		"max_quantity", maxQuantity)

	return maxQuantity, nil
}
