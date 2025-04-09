package service

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/strategy"
)

// StrategyService provides strategy functionality for the API
type StrategyService struct {
	factory *strategy.Factory
}

// NewStrategyService creates a new strategy service
func NewStrategyService(factory *strategy.Factory) *StrategyService {
	return &StrategyService{
		factory: factory,
	}
}

// StrategyParameter represents a strategy parameter
type StrategyParameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Default     interface{} `json:"default"`
	Min         interface{} `json:"min,omitempty"`
	Max         interface{} `json:"max,omitempty"`
	Options     []string    `json:"options,omitempty"`
	Required    bool        `json:"required"`
}

// StrategyPerformance represents the performance metrics of a strategy
type StrategyPerformance struct {
	WinRate      float64 `json:"winRate"`
	ProfitFactor float64 `json:"profitFactor"`
	SharpeRatio  float64 `json:"sharpeRatio"`
	MaxDrawdown  float64 `json:"maxDrawdown"`
}

// Strategy represents a trading strategy
type Strategy struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Parameters  []StrategyParameter `json:"parameters"`
	Performance StrategyPerformance `json:"performance"`
	IsEnabled   bool                `json:"isEnabled"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

// StrategyListItem represents a strategy in a list
type StrategyListItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	WinRate     float64 `json:"winRate"`
	IsEnabled   bool    `json:"isEnabled"`
}

// ListStrategies lists all available strategies
func (s *StrategyService) ListStrategies(ctx context.Context, enabled *bool) ([]*StrategyListItem, error) {
	// Get available strategies from factory
	availableStrategies := s.factory.GetAvailableStrategies()

	// Convert to API format
	strategies := make([]*StrategyListItem, 0, len(availableStrategies))
	for _, strategyName := range availableStrategies {
		// Get strategy configuration
		config, err := s.factory.GetStrategyConfig(strategyName)
		if err != nil {
			return nil, err
		}

		// Create strategy list item
		item := &StrategyListItem{
			ID:          strategyName,
			Name:        getStrategyDisplayName(strategyName),
			Description: getStrategyDescription(strategyName),
			WinRate:     getStrategyWinRate(strategyName),
			IsEnabled:   config.Enabled,
		}

		// Filter by enabled status if provided
		if enabled != nil && item.IsEnabled != *enabled {
			continue
		}

		strategies = append(strategies, item)
	}

	return strategies, nil
}

// GetStrategy gets a strategy by ID
func (s *StrategyService) GetStrategy(ctx context.Context, id string) (*Strategy, error) {
	// Check if strategy exists
	if !s.factory.StrategyExists(id) {
		return nil, fmt.Errorf("strategy '%s' not found", id)
	}

	// Get strategy configuration
	config, err := s.factory.GetStrategyConfig(id)
	if err != nil {
		return nil, err
	}

	// Create strategy
	strategy := &Strategy{
		ID:          id,
		Name:        getStrategyDisplayName(id),
		Description: getStrategyDescription(id),
		Parameters:  getStrategyParameters(id),
		Performance: StrategyPerformance{
			WinRate:      getStrategyWinRate(id),
			ProfitFactor: getStrategyProfitFactor(id),
			SharpeRatio:  getStrategySharpeRatio(id),
			MaxDrawdown:  getStrategyMaxDrawdown(id),
		},
		IsEnabled: config.Enabled,
		CreatedAt: time.Now().AddDate(0, -1, 0), // Mock creation time
		UpdatedAt: time.Now().AddDate(0, 0, -5), // Mock update time
	}

	return strategy, nil
}

// UpdateStrategyConfig updates a strategy configuration
func (s *StrategyService) UpdateStrategyConfig(ctx context.Context, id string, parameters map[string]interface{}, isEnabled bool) (*Strategy, error) {
	// Check if strategy exists
	if !s.factory.StrategyExists(id) {
		return nil, fmt.Errorf("strategy '%s' not found", id)
	}

	// Get current strategy configuration
	config, err := s.factory.GetStrategyConfig(id)
	if err != nil {
		return nil, err
	}

	// Update parameters
	for name, value := range parameters {
		config.Parameters[name] = value
	}

	// Update enabled status
	config.Enabled = isEnabled

	// Save configuration
	err = s.factory.SaveStrategyConfig(id, config)
	if err != nil {
		return nil, err
	}

	// Get updated strategy
	return s.GetStrategy(ctx, id)
}

// EnableStrategy enables a strategy
func (s *StrategyService) EnableStrategy(ctx context.Context, id string) (*StrategyEnableResult, error) {
	// Check if strategy exists
	if !s.factory.StrategyExists(id) {
		return nil, fmt.Errorf("strategy '%s' not found", id)
	}

	// Get current strategy configuration
	config, err := s.factory.GetStrategyConfig(id)
	if err != nil {
		return nil, err
	}

	// Update enabled status
	config.Enabled = true

	// Save configuration
	err = s.factory.SaveStrategyConfig(id, config)
	if err != nil {
		return nil, err
	}

	// Return result
	return &StrategyEnableResult{
		ID:        id,
		IsEnabled: true,
		UpdatedAt: time.Now(),
	}, nil
}

// DisableStrategy disables a strategy
func (s *StrategyService) DisableStrategy(ctx context.Context, id string) (*StrategyEnableResult, error) {
	// Check if strategy exists
	if !s.factory.StrategyExists(id) {
		return nil, fmt.Errorf("strategy '%s' not found", id)
	}

	// Get current strategy configuration
	config, err := s.factory.GetStrategyConfig(id)
	if err != nil {
		return nil, err
	}

	// Update enabled status
	config.Enabled = false

	// Save configuration
	err = s.factory.SaveStrategyConfig(id, config)
	if err != nil {
		return nil, err
	}

	// Return result
	return &StrategyEnableResult{
		ID:        id,
		IsEnabled: false,
		UpdatedAt: time.Now(),
	}, nil
}

// StrategyEnableResult represents the result of enabling or disabling a strategy
type StrategyEnableResult struct {
	ID        string    `json:"id"`
	IsEnabled bool      `json:"isEnabled"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Helper functions to get strategy information
// These would be replaced with actual implementations that get data from the strategy

func getStrategyDisplayName(id string) string {
	switch id {
	case "breakout":
		return "Breakout Strategy"
	case "volume-spike":
		return "Volume Spike Strategy"
	case "new-coin":
		return "New Coin Strategy"
	default:
		return id
	}
}

func getStrategyDescription(id string) string {
	switch id {
	case "breakout":
		return "A strategy that trades breakouts from support and resistance levels"
	case "volume-spike":
		return "A strategy that trades based on unusual volume spikes"
	case "new-coin":
		return "A strategy that trades newly listed coins"
	default:
		return "Unknown strategy"
	}
}

func getStrategyParameters(id string) []StrategyParameter {
	switch id {
	case "breakout":
		return []StrategyParameter{
			{
				Name:        "lookbackPeriod",
				Type:        "integer",
				Description: "Number of periods to look back for support/resistance",
				Default:     20,
				Min:         5,
				Max:         100,
				Required:    true,
			},
			{
				Name:        "breakoutThreshold",
				Type:        "float",
				Description: "Threshold for breakout detection",
				Default:     2.0,
				Min:         0.5,
				Max:         5.0,
				Required:    true,
			},
			{
				Name:        "confirmationPeriods",
				Type:        "integer",
				Description: "Number of periods to confirm breakout",
				Default:     3,
				Min:         1,
				Max:         10,
				Required:    false,
			},
		}
	case "volume-spike":
		return []StrategyParameter{
			{
				Name:        "volumeMultiplier",
				Type:        "float",
				Description: "Multiplier for average volume to detect spike",
				Default:     3.0,
				Min:         1.5,
				Max:         10.0,
				Required:    true,
			},
			{
				Name:        "averagePeriod",
				Type:        "integer",
				Description: "Number of periods to calculate average volume",
				Default:     24,
				Min:         6,
				Max:         72,
				Required:    true,
			},
			{
				Name:        "priceChangeThreshold",
				Type:        "float",
				Description: "Minimum price change percentage to consider",
				Default:     1.0,
				Min:         0.1,
				Max:         5.0,
				Required:    false,
			},
		}
	case "new-coin":
		return []StrategyParameter{
			{
				Name:        "maxAgeHours",
				Type:        "integer",
				Description: "Maximum age of coin in hours",
				Default:     24,
				Min:         1,
				Max:         72,
				Required:    true,
			},
			{
				Name:        "minVolume",
				Type:        "float",
				Description: "Minimum volume in USDT",
				Default:     100000.0,
				Min:         10000.0,
				Max:         1000000.0,
				Required:    true,
			},
			{
				Name:        "entryType",
				Type:        "string",
				Description: "Type of entry strategy",
				Default:     "immediate",
				Options:     []string{"immediate", "pullback", "breakout"},
				Required:    true,
			},
		}
	default:
		return []StrategyParameter{}
	}
}

func getStrategyWinRate(id string) float64 {
	switch id {
	case "breakout":
		return 65.0
	case "volume-spike":
		return 58.0
	case "new-coin":
		return 72.0
	default:
		return 50.0
	}
}

func getStrategyProfitFactor(id string) float64 {
	switch id {
	case "breakout":
		return 2.1
	case "volume-spike":
		return 1.8
	case "new-coin":
		return 2.5
	default:
		return 1.0
	}
}

func getStrategySharpeRatio(id string) float64 {
	switch id {
	case "breakout":
		return 1.5
	case "volume-spike":
		return 1.2
	case "new-coin":
		return 1.8
	default:
		return 1.0
	}
}

func getStrategyMaxDrawdown(id string) float64 {
	switch id {
	case "breakout":
		return 15.0
	case "volume-spike":
		return 18.0
	case "new-coin":
		return 12.0
	default:
		return 20.0
	}
}
