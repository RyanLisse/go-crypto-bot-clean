package strategy

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// MarketData contains information about the current market conditions
type MarketData struct {
	Symbol string
	Regime string
}

// StrategyCreator is a function that creates a strategy
type StrategyCreator func(ctx context.Context, config map[string]interface{}) (Strategy, error)

// ConfigValidator is a function that validates a strategy configuration
type ConfigValidator func(config map[string]interface{}) (bool, []string, error)

// StrategyFactoryImpl implements the StrategyFactory interface
type StrategyFactoryImpl struct {
	strategies       map[string]StrategyCreator
	defaultConfigs   map[string]map[string]interface{}
	configValidators map[string]ConfigValidator
	regimeStrategies map[string]string
	defaultStrategy  string
	mu               sync.RWMutex
}

// NewStrategyFactory creates a new strategy factory
func NewStrategyFactory() *StrategyFactoryImpl {
	return &StrategyFactoryImpl{
		strategies:       make(map[string]StrategyCreator),
		defaultConfigs:   make(map[string]map[string]interface{}),
		configValidators: make(map[string]ConfigValidator),
		regimeStrategies: make(map[string]string),
	}
}

// RegisterStrategy registers a strategy creator function
func (f *StrategyFactoryImpl) RegisterStrategy(name string, creator StrategyCreator) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.strategies[name] = creator
}

// RegisterStrategyWithConfig registers a strategy with a default configuration
func (f *StrategyFactoryImpl) RegisterStrategyWithConfig(name string, creator StrategyCreator, defaultConfig map[string]interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.strategies[name] = creator
	f.defaultConfigs[name] = defaultConfig
}

// RegisterStrategyWithValidation registers a strategy with a default configuration and validation function
func (f *StrategyFactoryImpl) RegisterStrategyWithValidation(name string, creator StrategyCreator, defaultConfig map[string]interface{}, validator ConfigValidator) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.strategies[name] = creator
	f.defaultConfigs[name] = defaultConfig
	f.configValidators[name] = validator
}

// CreateStrategy creates a strategy by name
func (f *StrategyFactoryImpl) CreateStrategy(ctx context.Context, name string, config map[string]interface{}) (Strategy, error) {
	f.mu.RLock()
	creator, exists := f.strategies[name]
	f.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("strategy %s not found", name)
	}

	// If no config is provided, use the default config
	if config == nil {
		f.mu.RLock()
		defaultConfig, exists := f.defaultConfigs[name]
		f.mu.RUnlock()

		if exists {
			config = defaultConfig
		} else {
			config = make(map[string]interface{})
		}
	}

	// Create the strategy
	strategy, err := creator(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create strategy %s: %w", name, err)
	}

	return strategy, nil
}

// ListAvailableStrategies returns all available strategy names
func (f *StrategyFactoryImpl) ListAvailableStrategies() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	strategies := make([]string, 0, len(f.strategies))
	for name := range f.strategies {
		strategies = append(strategies, name)
	}

	sort.Strings(strategies)
	return strategies
}

// GetDefaultConfig returns the default configuration for a strategy
func (f *StrategyFactoryImpl) GetDefaultConfig(strategyName string) map[string]interface{} {
	f.mu.RLock()
	defer f.mu.RUnlock()

	config, exists := f.defaultConfigs[strategyName]
	if !exists {
		return make(map[string]interface{})
	}

	return config
}

// ValidateConfig checks if a configuration is valid for a strategy
func (f *StrategyFactoryImpl) ValidateConfig(strategyName string, config map[string]interface{}) (bool, []string, error) {
	f.mu.RLock()
	validator, exists := f.configValidators[strategyName]
	f.mu.RUnlock()

	if !exists {
		return false, nil, fmt.Errorf("strategy %s not found or has no validator", strategyName)
	}

	return validator(config)
}

// SetRegimeStrategy sets the strategy to use for a specific market regime
func (f *StrategyFactoryImpl) SetRegimeStrategy(regime string, strategyName string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.regimeStrategies[regime] = strategyName
}

// SetDefaultStrategy sets the default strategy to use when no regime-specific strategy is found
func (f *StrategyFactoryImpl) SetDefaultStrategy(strategyName string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.defaultStrategy = strategyName
}

// GetStrategyForMarketRegime returns the appropriate strategy for the current market regime
func (f *StrategyFactoryImpl) GetStrategyForMarketRegime(ctx context.Context, marketData *MarketData) (Strategy, error) {
	f.mu.RLock()
	strategyName, exists := f.regimeStrategies[marketData.Regime]
	if !exists {
		strategyName = f.defaultStrategy
	}
	f.mu.RUnlock()

	if strategyName == "" {
		return nil, fmt.Errorf("no strategy found for regime %s and no default strategy set", marketData.Regime)
	}

	return f.CreateStrategy(ctx, strategyName, nil)
}
