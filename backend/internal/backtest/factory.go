package backtest

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// StrategyCreator is a function that creates a new strategy
type StrategyCreator func(logger *zap.Logger) BacktestStrategy

// StrategyFactory creates and manages trading strategies
type StrategyFactory struct {
	strategies map[string]StrategyCreator
	configs    map[string]map[string]interface{}
	logger     *zap.Logger
	mu         sync.RWMutex
}

// NewStrategyFactory creates a new strategy factory
func NewStrategyFactory() *StrategyFactory {
	logger, _ := zap.NewProduction()
	return &StrategyFactory{
		strategies: make(map[string]StrategyCreator),
		configs:    make(map[string]map[string]interface{}),
		logger:     logger,
	}
}

// RegisterStrategy registers a strategy creator function
func (f *StrategyFactory) RegisterStrategy(name string, creator StrategyCreator) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.strategies[name] = creator
	f.logger.Info("Registered strategy",
		zap.String("name", name),
	)
}

// CreateStrategy creates a new strategy by name
func (f *StrategyFactory) CreateStrategy(ctx context.Context, name string) (BacktestStrategy, error) {
	f.mu.RLock()
	creator, ok := f.strategies[name]
	config := f.getStrategyConfig(name)
	f.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown strategy: %s", name)
	}

	// Create the strategy
	strategy := creator(f.logger)

	// Initialize with config if available
	if config != nil {
		if err := strategy.Initialize(ctx, config); err != nil {
			return nil, fmt.Errorf("failed to initialize strategy %s: %w", name, err)
		}
	}

	return strategy, nil
}

// ListAvailableStrategies returns all registered strategy names
func (f *StrategyFactory) ListAvailableStrategies() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	names := make([]string, 0, len(f.strategies))
	for name := range f.strategies {
		names = append(names, name)
	}

	return names
}

// SaveStrategyConfig saves a configuration for a strategy
func (f *StrategyFactory) SaveStrategyConfig(name string, config map[string]interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.configs[name] = config
}

// GetStrategyConfig returns the configuration for a strategy
func (f *StrategyFactory) GetStrategyConfig(name string) map[string]interface{} {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.getStrategyConfig(name)
}

// getStrategyConfig is an internal method that returns the configuration for a strategy
// without locking (caller must handle locking)
func (f *StrategyFactory) getStrategyConfig(name string) map[string]interface{} {
	config, ok := f.configs[name]
	if !ok {
		return nil
	}
	return config
}

// DeleteStrategyConfig deletes a strategy configuration
func (f *StrategyFactory) DeleteStrategyConfig(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.configs, name)
}
