package backtest

import (
	"fmt"
	// "go-crypto-bot-clean/backend/internal/domain/models" // Not needed if local DefaultStrategy removed
)

// StrategyFactory creates strategy instances from strategy names
type StrategyFactory struct {
	strategies map[string]func() BacktestStrategy
}

// NewStrategyFactory creates a new strategy factory
func NewStrategyFactory() *StrategyFactory {
	return &StrategyFactory{
		strategies: make(map[string]func() BacktestStrategy),
	}
}

// RegisterStrategy registers a strategy with the factory
func (f *StrategyFactory) RegisterStrategy(name string, creator func() BacktestStrategy) {
	f.strategies[name] = creator
}

// CreateStrategy creates a strategy instance from a strategy name
func (f *StrategyFactory) CreateStrategy(name string) (BacktestStrategy, error) {
	creator, ok := f.strategies[name]
	if !ok {
		return nil, fmt.Errorf("unknown strategy: %s", name)
	}
	return creator(), nil
}
