package strategy

import (
	"context"
	"testing"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/strategy/advanced"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStrategy is a mock implementation of the Strategy interface
type MockStrategy struct {
	mock.Mock
}

func (m *MockStrategy) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockStrategy) Initialize(ctx context.Context, config map[string]interface{}) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockStrategy) UpdateParameters(ctx context.Context, params map[string]interface{}) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockStrategy) OnPriceUpdate(ctx context.Context, priceUpdate *models.PriceUpdate) (*Signal, error) {
	args := m.Called(ctx, priceUpdate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Signal), args.Error(1)
}

func (m *MockStrategy) OnCandleUpdate(ctx context.Context, candle *models.Candle) (*Signal, error) {
	args := m.Called(ctx, candle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Signal), args.Error(1)
}

func (m *MockStrategy) OnTradeUpdate(ctx context.Context, trade *models.Trade) (*Signal, error) {
	args := m.Called(ctx, trade)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Signal), args.Error(1)
}

func (m *MockStrategy) OnMarketDepthUpdate(ctx context.Context, depth *models.OrderBook) (*Signal, error) {
	args := m.Called(ctx, depth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Signal), args.Error(1)
}

func (m *MockStrategy) OnTimerEvent(ctx context.Context, eventType string) (*Signal, error) {
	args := m.Called(ctx, eventType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Signal), args.Error(1)
}

func (m *MockStrategy) GetTimeframes() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockStrategy) GetRequiredDataTypes() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockStrategy) PerformBacktest(ctx context.Context, historicalData []*models.Candle, params map[string]interface{}) ([]*Signal, *models.BacktestResult, error) {
	args := m.Called(ctx, historicalData, params)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*Signal), args.Get(1).(*models.BacktestResult), args.Error(2)
}

// TestStrategyFactoryCreateStrategy tests the CreateStrategy method
func TestStrategyFactoryCreateStrategy(t *testing.T) {
	// Create a new factory
	factory := NewStrategyFactory()

	// Test creating a strategy that doesn't exist
	ctx := context.Background()
	_, err := factory.CreateStrategy(ctx, "non_existent_strategy", nil)
	assert.Error(t, err, "Creating a non-existent strategy should return an error")

	// Register a mock strategy
	mockStrategy := new(MockStrategy)
	mockStrategy.On("GetName").Return("test_strategy")

	factory.RegisterStrategy("test_strategy", func(ctx context.Context, config map[string]interface{}) (Strategy, error) {
		return mockStrategy, nil
	})

	// Test creating a registered strategy
	strategy, err := factory.CreateStrategy(ctx, "test_strategy", nil)
	assert.NoError(t, err, "Creating a registered strategy should not return an error")
	assert.Equal(t, "test_strategy", strategy.GetName(), "Strategy name should match")

	mockStrategy.AssertExpectations(t)
}

// TestStrategyFactoryListAvailableStrategies tests the ListAvailableStrategies method
func TestStrategyFactoryListAvailableStrategies(t *testing.T) {
	// Create a new factory
	factory := NewStrategyFactory()

	// Register some strategies
	factory.RegisterStrategy("strategy1", nil)
	factory.RegisterStrategy("strategy2", nil)
	factory.RegisterStrategy("strategy3", nil)

	// Test listing available strategies
	strategies := factory.ListAvailableStrategies()
	assert.Len(t, strategies, 3, "Should have 3 registered strategies")
	assert.Contains(t, strategies, "strategy1", "Should contain strategy1")
	assert.Contains(t, strategies, "strategy2", "Should contain strategy2")
	assert.Contains(t, strategies, "strategy3", "Should contain strategy3")
}

// TestStrategyFactoryGetDefaultConfig tests the GetDefaultConfig method
func TestStrategyFactoryGetDefaultConfig(t *testing.T) {
	// Create a new factory
	factory := NewStrategyFactory()

	// Register a strategy with default config
	defaultConfig := map[string]interface{}{
		"param1": 10,
		"param2": "value",
	}

	factory.RegisterStrategyWithConfig("test_strategy", nil, defaultConfig)

	// Test getting default config
	config := factory.GetDefaultConfig("test_strategy")
	assert.Equal(t, defaultConfig, config, "Default config should match")

	// Test getting default config for non-existent strategy
	config = factory.GetDefaultConfig("non_existent_strategy")
	assert.Empty(t, config, "Default config for non-existent strategy should be empty")
}

// TestStrategyFactoryValidateConfig tests the ValidateConfig method
func TestStrategyFactoryValidateConfig(t *testing.T) {
	// Create a new factory
	factory := NewStrategyFactory()

	// Register a strategy with validation
	factory.RegisterStrategyWithValidation("test_strategy", nil, nil, func(config map[string]interface{}) (bool, []string, error) {
		errors := []string{}

		// Check required parameters
		if _, ok := config["required_param"]; !ok {
			errors = append(errors, "required_param is required")
		}

		// Check parameter types
		if val, ok := config["numeric_param"]; ok {
			if _, ok := val.(float64); !ok {
				errors = append(errors, "numeric_param must be a number")
			}
		}

		return len(errors) == 0, errors, nil
	})

	// Test valid config
	valid, errors, err := factory.ValidateConfig("test_strategy", map[string]interface{}{
		"required_param": "value",
		"numeric_param":  10.5,
	})
	assert.NoError(t, err, "Validation should not return an error")
	assert.True(t, valid, "Config should be valid")
	assert.Empty(t, errors, "There should be no validation errors")

	// Test invalid config
	valid, errors, err = factory.ValidateConfig("test_strategy", map[string]interface{}{
		"numeric_param": "not a number",
	})
	assert.NoError(t, err, "Validation should not return an error")
	assert.False(t, valid, "Config should be invalid")
	assert.Len(t, errors, 2, "There should be 2 validation errors")

	// Test non-existent strategy
	valid, errors, err = factory.ValidateConfig("non_existent_strategy", nil)
	assert.Error(t, err, "Validating a non-existent strategy should return an error")
}

// TestGetStrategyForMarketRegime tests selecting a strategy based on market regime
func TestGetStrategyForMarketRegime(t *testing.T) {
	// Create a new factory
	factory := NewStrategyFactory()

	// Register strategies for different regimes
	trendingUpStrategy := new(MockStrategy)
	trendingUpStrategy.On("GetName").Return("trending_up_strategy")

	trendingDownStrategy := new(MockStrategy)
	trendingDownStrategy.On("GetName").Return("trending_down_strategy")

	rangingStrategy := new(MockStrategy)
	rangingStrategy.On("GetName").Return("ranging_strategy")

	volatileStrategy := new(MockStrategy)
	volatileStrategy.On("GetName").Return("volatile_strategy")

	defaultStrategy := new(MockStrategy)
	defaultStrategy.On("GetName").Return("default_strategy")

	factory.RegisterStrategy("trending_up_strategy", func(ctx context.Context, config map[string]interface{}) (Strategy, error) {
		return trendingUpStrategy, nil
	})

	factory.RegisterStrategy("trending_down_strategy", func(ctx context.Context, config map[string]interface{}) (Strategy, error) {
		return trendingDownStrategy, nil
	})

	factory.RegisterStrategy("ranging_strategy", func(ctx context.Context, config map[string]interface{}) (Strategy, error) {
		return rangingStrategy, nil
	})

	factory.RegisterStrategy("volatile_strategy", func(ctx context.Context, config map[string]interface{}) (Strategy, error) {
		return volatileStrategy, nil
	})

	factory.RegisterStrategy("default_strategy", func(ctx context.Context, config map[string]interface{}) (Strategy, error) {
		return defaultStrategy, nil
	})

	// Create market data with different regimes
	trendingUpData := &MarketData{
		Symbol: "BTC/USDT",
		Regime: string(advanced.RegimeTrendingUp),
	}

	trendingDownData := &MarketData{
		Symbol: "BTC/USDT",
		Regime: string(advanced.RegimeTrendingDown),
	}

	rangingData := &MarketData{
		Symbol: "BTC/USDT",
		Regime: string(advanced.RegimeRanging),
	}

	volatileData := &MarketData{
		Symbol: "BTC/USDT",
		Regime: string(advanced.RegimeVolatile),
	}

	unknownData := &MarketData{
		Symbol: "BTC/USDT",
		Regime: "UNKNOWN",
	}

	// Configure factory to use specific strategies for different regimes
	factory.SetRegimeStrategy(string(advanced.RegimeTrendingUp), "trending_up_strategy")
	factory.SetRegimeStrategy(string(advanced.RegimeTrendingDown), "trending_down_strategy")
	factory.SetRegimeStrategy(string(advanced.RegimeRanging), "ranging_strategy")
	factory.SetRegimeStrategy(string(advanced.RegimeVolatile), "volatile_strategy")
	factory.SetDefaultStrategy("default_strategy")

	// Test getting strategy for different regimes
	ctx := context.Background()

	strategy, err := factory.GetStrategyForMarketRegime(ctx, trendingUpData)
	assert.NoError(t, err, "Getting strategy for trending up regime should not return an error")
	assert.Equal(t, "trending_up_strategy", strategy.GetName(), "Should return trending up strategy")

	strategy, err = factory.GetStrategyForMarketRegime(ctx, trendingDownData)
	assert.NoError(t, err, "Getting strategy for trending down regime should not return an error")
	assert.Equal(t, "trending_down_strategy", strategy.GetName(), "Should return trending down strategy")

	strategy, err = factory.GetStrategyForMarketRegime(ctx, rangingData)
	assert.NoError(t, err, "Getting strategy for ranging regime should not return an error")
	assert.Equal(t, "ranging_strategy", strategy.GetName(), "Should return ranging strategy")

	strategy, err = factory.GetStrategyForMarketRegime(ctx, volatileData)
	assert.NoError(t, err, "Getting strategy for volatile regime should not return an error")
	assert.Equal(t, "volatile_strategy", strategy.GetName(), "Should return volatile strategy")

	strategy, err = factory.GetStrategyForMarketRegime(ctx, unknownData)
	assert.NoError(t, err, "Getting strategy for unknown regime should not return an error")
	assert.Equal(t, "default_strategy", strategy.GetName(), "Should return default strategy")

	// Test with no default strategy
	factory.SetDefaultStrategy("")
	_, err = factory.GetStrategyForMarketRegime(ctx, unknownData)
	assert.Error(t, err, "Getting strategy for unknown regime with no default should return an error")
}
