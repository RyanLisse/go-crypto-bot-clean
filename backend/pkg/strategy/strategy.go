// Package strategy provides public interfaces for strategy services
package strategy

import (
	"context"
	"go-crypto-bot-clean/backend/internal/domain/strategy"
)

// Factory is a public interface for the strategy factory
type Factory interface {
	// CreateStrategy creates a new strategy with the given name and parameters
	CreateStrategy(name string, params map[string]interface{}) (Strategy, error)

	// GetAvailableStrategies gets a list of available strategies
	GetAvailableStrategies() []string
}

// Strategy is a public interface for trading strategies
type Strategy interface {
	// GetName gets the name of the strategy
	GetName() string

	// GetDescription gets the description of the strategy
	GetDescription() string

	// GetParameters gets the parameters of the strategy
	GetParameters() map[string]interface{}
}

// factoryAdapter adapts the internal strategy factory to the public interface
type factoryAdapter struct {
	internalFactory *strategy.StrategyFactoryImpl
}

// CreateStrategy adapts the internal CreateStrategy method to the public interface
func (a *factoryAdapter) CreateStrategy(name string, params map[string]interface{}) (Strategy, error) {
	// Create a context for the internal call
	ctx := context.Background()

	// Call the internal factory
	internalStrategy, err := a.internalFactory.CreateStrategy(ctx, name, params)
	if err != nil {
		return nil, err
	}

	// Adapt the internal strategy to the public interface
	return &strategyAdapter{
		internalStrategy: internalStrategy,
	}, nil
}

// GetAvailableStrategies adapts the internal ListAvailableStrategies method to the public interface
func (a *factoryAdapter) GetAvailableStrategies() []string {
	return a.internalFactory.ListAvailableStrategies()
}

// strategyAdapter adapts the internal strategy to the public interface
type strategyAdapter struct {
	internalStrategy strategy.Strategy
}

// GetName adapts the internal GetName method to the public interface
func (a *strategyAdapter) GetName() string {
	return a.internalStrategy.GetName()
}

// GetDescription provides a description of the strategy
func (a *strategyAdapter) GetDescription() string {
	// Internal strategy doesn't have a description method, so we'll return a generic one
	return "Trading strategy: " + a.internalStrategy.GetName()
}

// GetParameters returns the parameters of the strategy
func (a *strategyAdapter) GetParameters() map[string]interface{} {
	// Internal strategy doesn't have a GetParameters method, so we'll return an empty map
	return make(map[string]interface{})
}

// NewFactory creates a new strategy factory
func NewFactory() Factory {
	return &factoryAdapter{
		internalFactory: strategy.NewStrategyFactory(),
	}
}
