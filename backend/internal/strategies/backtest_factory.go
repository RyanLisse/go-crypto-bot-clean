package strategies

import (
	"go-crypto-bot-clean/backend/internal/backtest"
)

// RegisterBacktestStrategies registers strategies with the backtest factory
func RegisterBacktestStrategies(factory *backtest.StrategyFactory) {
	// Register the moving average strategy
	factory.RegisterStrategy("moving_average", func() backtest.BacktestStrategy {
		// Create a simple moving average strategy with default parameters
		// This is a placeholder implementation
		return NewSimpleMAStrategy(10, 20, nil)
	})

	// Register other strategies as needed
}

// Initialize the backtest factory with strategies
func init() {
	// This will be called when the strategies package is imported
	// We can't directly call RegisterBacktestStrategies here because
	// we don't have access to the factory instance
}
