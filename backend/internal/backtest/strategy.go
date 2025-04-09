package backtest

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// BaseStrategy provides a basic implementation of the BacktestStrategy interface
type BaseStrategy struct {
	// Common fields for strategies can be added here
}

// Initialize sets up the strategy
func (s *BaseStrategy) Initialize(ctx context.Context, config interface{}) error {
	// Base implementation does nothing, can be overridden
	return nil
}

// OnTick processes a market tick
func (s *BaseStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*Signal, error) {
	// Base implementation returns no signals, must be overridden
	return nil, nil
}

// OnOrderFilled handles filled orders
func (s *BaseStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
	// Base implementation does nothing, can be overridden
	return nil
}

// ClosePositions is called at the end of the backtest
func (s *BaseStrategy) ClosePositions(ctx context.Context) ([]*Signal, error) {
	// Base implementation returns no signals, can be overridden
	return nil, nil
}

// DefaultStrategy is a simple example strategy
type DefaultStrategy struct {
	BaseStrategy // Embed BaseStrategy
	// Add strategy-specific fields here
}

// NewDefaultStrategy creates a new DefaultStrategy
func NewDefaultStrategy() BacktestStrategy {
	return &DefaultStrategy{}
}

// Initialize overrides the base Initialize method to match the interface
func (s *DefaultStrategy) Initialize(ctx context.Context, config interface{}) error {
	// Call base implementation first if needed (optional)
	if err := s.BaseStrategy.Initialize(ctx, config); err != nil {
		return err
	}
	// Add DefaultStrategy specific initialization here if needed
	return nil
}

// OnTick implements the trading logic for DefaultStrategy
func (s *DefaultStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*Signal, error) {
	// Implement actual trading logic here
	return nil, nil // Placeholder
}

// ClosePositions implements the method required by the interface
// DefaultStrategy might not need to do anything specific on close
func (s *DefaultStrategy) ClosePositions(ctx context.Context) ([]*Signal, error) {
	// You could call the base implementation if it ever does something:
	// return s.BaseStrategy.ClosePositions(ctx)
	return nil, nil // Default implementation returns no closing signals
}

// Ensure DefaultStrategy implements BacktestStrategy
var _ BacktestStrategy = (*DefaultStrategy)(nil)
