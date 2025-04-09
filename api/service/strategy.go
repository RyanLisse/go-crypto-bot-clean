package service

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/pkg/strategy"
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
	// Mock implementation since we can't directly call the factory's GetAvailableStrategies method
	// In a real implementation, we would call the factory to get the available strategies
	
	// Create mock strategies
	strategies := []*StrategyListItem{
		{
			ID:          "breakout",
			Name:        "Breakout Strategy",
			Description: "A strategy that trades breakouts from support and resistance levels",
			WinRate:     65.0,
			IsEnabled:   true,
		},
		{
			ID:          "volume-spike",
			Name:        "Volume Spike Strategy",
			Description: "A strategy that trades based on unusual volume spikes",
			WinRate:     58.0,
			IsEnabled:   false,
		},
		{
			ID:          "new-coin",
			Name:        "New Coin Strategy",
			Description: "A strategy that trades newly listed coins",
			WinRate:     72.0,
			IsEnabled:   true,
		},
	}

	// Filter by enabled status if provided
	if enabled != nil {
		filteredStrategies := make([]*StrategyListItem, 0)
		for _, strategy := range strategies {
			if strategy.IsEnabled == *enabled {
				filteredStrategies = append(filteredStrategies, strategy)
			}
		}
		strategies = filteredStrategies
	}

	return strategies, nil
}

// GetStrategy gets a strategy by ID
func (s *StrategyService) GetStrategy(ctx context.Context, id string) (*Strategy, error) {
	// Mock implementation since we can't directly call the factory's methods
	// In a real implementation, we would call the factory to get the strategy
	
	// Create mock strategy based on ID
	var strategy *Strategy
	switch id {
	case "breakout":
		strategy = &Strategy{
			ID:          "breakout",
			Name:        "Breakout Strategy",
			Description: "A strategy that trades breakouts from support and resistance levels",
			Parameters: []StrategyParameter{
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
			},
			Performance: StrategyPerformance{
				WinRate:      65.0,
				ProfitFactor: 2.1,
				SharpeRatio:  1.5,
				MaxDrawdown:  15.0,
			},
			IsEnabled: true,
			CreatedAt: time.Now().AddDate(0, -3, 0),
			UpdatedAt: time.Now().AddDate(0, 0, -5),
		}
	case "volume-spike":
		strategy = &Strategy{
			ID:          "volume-spike",
			Name:        "Volume Spike Strategy",
			Description: "A strategy that trades based on unusual volume spikes",
			Parameters: []StrategyParameter{
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
			},
			Performance: StrategyPerformance{
				WinRate:      58.0,
				ProfitFactor: 1.8,
				SharpeRatio:  1.2,
				MaxDrawdown:  18.0,
			},
			IsEnabled: false,
			CreatedAt: time.Now().AddDate(0, -2, 0),
			UpdatedAt: time.Now().AddDate(0, 0, -3),
		}
	case "new-coin":
		strategy = &Strategy{
			ID:          "new-coin",
			Name:        "New Coin Strategy",
			Description: "A strategy that trades newly listed coins",
			Parameters: []StrategyParameter{
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
			},
			Performance: StrategyPerformance{
				WinRate:      72.0,
				ProfitFactor: 2.5,
				SharpeRatio:  1.8,
				MaxDrawdown:  12.0,
			},
			IsEnabled: true,
			CreatedAt: time.Now().AddDate(0, -1, 0),
			UpdatedAt: time.Now().AddDate(0, 0, -1),
		}
	default:
		return nil, fmt.Errorf("strategy '%s' not found", id)
	}

	return strategy, nil
}

// StrategyConfig represents a strategy configuration
type StrategyConfig struct {
	Enabled    bool                   `json:"enabled"`
	Parameters map[string]interface{} `json:"parameters"`
}

// UpdateStrategyConfig updates a strategy configuration
func (s *StrategyService) UpdateStrategyConfig(ctx context.Context, id string, parameters map[string]interface{}, isEnabled bool) (*Strategy, error) {
	// Mock implementation since we can't directly call the factory's methods
	// In a real implementation, we would call the factory to update the strategy configuration
	
	// Check if strategy exists
	strategy, err := s.GetStrategy(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update enabled status
	strategy.IsEnabled = isEnabled

	// Update parameters
	// In a real implementation, we would update the parameters in the strategy
	
	// Update timestamp
	strategy.UpdatedAt = time.Now()

	return strategy, nil
}

// StrategyEnableResult represents the result of enabling or disabling a strategy
type StrategyEnableResult struct {
	ID        string    `json:"id"`
	IsEnabled bool      `json:"isEnabled"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// EnableStrategy enables a strategy
func (s *StrategyService) EnableStrategy(ctx context.Context, id string) (*StrategyEnableResult, error) {
	// Mock implementation since we can't directly call the factory's methods
	// In a real implementation, we would call the factory to enable the strategy
	
	// Check if strategy exists
	_, err := s.GetStrategy(ctx, id)
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
	// Mock implementation since we can't directly call the factory's methods
	// In a real implementation, we would call the factory to disable the strategy
	
	// Check if strategy exists
	_, err := s.GetStrategy(ctx, id)
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
