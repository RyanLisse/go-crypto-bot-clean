// Package strategy provides the strategy endpoints for the Huma API.
package strategy

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// Parameter represents a strategy parameter
type Parameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Default     interface{} `json:"default"`
	Min         interface{} `json:"min,omitempty"`
	Max         interface{} `json:"max,omitempty"`
	Options     []string    `json:"options,omitempty"`
	Required    bool        `json:"required"`
}

// Performance represents the performance metrics of a strategy
type Performance struct {
	WinRate      float64 `json:"winRate"`
	ProfitFactor float64 `json:"profitFactor"`
	SharpeRatio  float64 `json:"sharpeRatio"`
	MaxDrawdown  float64 `json:"maxDrawdown"`
}

// Strategy represents a trading strategy
type Strategy struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Parameters  []Parameter  `json:"parameters"`
	Performance Performance  `json:"performance"`
	IsEnabled   bool         `json:"isEnabled"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

// StrategyListItem represents a strategy in a list
type StrategyListItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	WinRate     float64 `json:"winRate"`
	IsEnabled   bool    `json:"isEnabled"`
}

// RegisterEndpoints registers the strategy endpoints.
func RegisterEndpoints(api huma.API, basePath string) {
	// GET /strategy
	huma.Register(api, huma.Operation{
		OperationID: "list-strategies",
		Method:      http.MethodGet,
		Path:        basePath + "/strategy",
		Summary:     "List strategies",
		Description: "Returns a list of available trading strategies",
		Tags:        []string{"Strategy"},
	}, func(ctx context.Context, input *struct {
		Enabled string `query:"enabled"`
	}) (*struct {
		Body struct {
			Strategies []StrategyListItem `json:"strategies"`
			Count      int                `json:"count"`
			Timestamp  time.Time          `json:"timestamp"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Strategies []StrategyListItem `json:"strategies"`
				Count      int                `json:"count"`
				Timestamp  time.Time          `json:"timestamp"`
			}
		}{}

		// Add some sample strategies
		resp.Body.Strategies = []StrategyListItem{
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
		if input.Enabled != "" {
			filteredStrategies := []StrategyListItem{}
			for _, strategy := range resp.Body.Strategies {
				isEnabled := input.Enabled == "true"
				if strategy.IsEnabled == isEnabled {
					filteredStrategies = append(filteredStrategies, strategy)
				}
			}
			resp.Body.Strategies = filteredStrategies
		}

		resp.Body.Count = len(resp.Body.Strategies)
		resp.Body.Timestamp = time.Now()
		return resp, nil
	})

	// GET /strategy/{id}
	huma.Register(api, huma.Operation{
		OperationID: "get-strategy",
		Method:      http.MethodGet,
		Path:        basePath + "/strategy/{id}",
		Summary:     "Get strategy",
		Description: "Returns details of a specific trading strategy",
		Tags:        []string{"Strategy"},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id"`
	}) (*struct {
		Body Strategy
	}, error) {
		resp := &struct {
			Body Strategy
		}{}

		// Set strategy details based on ID
		switch input.ID {
		case "breakout":
			resp.Body = Strategy{
				ID:          "breakout",
				Name:        "Breakout Strategy",
				Description: "A strategy that trades breakouts from support and resistance levels",
				Parameters: []Parameter{
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
				Performance: Performance{
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
			resp.Body = Strategy{
				ID:          "volume-spike",
				Name:        "Volume Spike Strategy",
				Description: "A strategy that trades based on unusual volume spikes",
				Parameters: []Parameter{
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
				Performance: Performance{
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
			resp.Body = Strategy{
				ID:          "new-coin",
				Name:        "New Coin Strategy",
				Description: "A strategy that trades newly listed coins",
				Parameters: []Parameter{
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
				Performance: Performance{
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
			// Return a default strategy if ID not found
			resp.Body = Strategy{
				ID:          input.ID,
				Name:        "Unknown Strategy",
				Description: "Strategy not found",
				IsEnabled:   false,
				CreatedAt:   time.Now().AddDate(0, -1, 0),
				UpdatedAt:   time.Now().AddDate(0, 0, -1),
			}
		}

		return resp, nil
	})

	// PUT /strategy/{id}
	huma.Register(api, huma.Operation{
		OperationID: "update-strategy-config",
		Method:      http.MethodPut,
		Path:        basePath + "/strategy/{id}",
		Summary:     "Update strategy configuration",
		Description: "Updates the configuration of a specific trading strategy",
		Tags:        []string{"Strategy"},
	}, func(ctx context.Context, input *struct {
		ID   string `path:"id"`
		Body struct {
			Parameters map[string]interface{} `json:"parameters"`
			IsEnabled  bool                   `json:"isEnabled"`
		}
	}) (*struct {
		Body struct {
			ID         string                 `json:"id"`
			Parameters map[string]interface{} `json:"parameters"`
			IsEnabled  bool                   `json:"isEnabled"`
			UpdatedAt  time.Time              `json:"updatedAt"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				ID         string                 `json:"id"`
				Parameters map[string]interface{} `json:"parameters"`
				IsEnabled  bool                   `json:"isEnabled"`
				UpdatedAt  time.Time              `json:"updatedAt"`
			}
		}{}
		resp.Body.ID = input.ID
		resp.Body.Parameters = input.Body.Parameters
		resp.Body.IsEnabled = input.Body.IsEnabled
		resp.Body.UpdatedAt = time.Now()
		return resp, nil
	})

	// POST /strategy/{id}/enable
	huma.Register(api, huma.Operation{
		OperationID: "enable-strategy",
		Method:      http.MethodPost,
		Path:        basePath + "/strategy/{id}/enable",
		Summary:     "Enable strategy",
		Description: "Enables a specific trading strategy",
		Tags:        []string{"Strategy"},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id"`
	}) (*struct {
		Body struct {
			ID        string    `json:"id"`
			IsEnabled bool      `json:"isEnabled"`
			UpdatedAt time.Time `json:"updatedAt"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				ID        string    `json:"id"`
				IsEnabled bool      `json:"isEnabled"`
				UpdatedAt time.Time `json:"updatedAt"`
			}
		}{}
		resp.Body.ID = input.ID
		resp.Body.IsEnabled = true
		resp.Body.UpdatedAt = time.Now()
		return resp, nil
	})

	// POST /strategy/{id}/disable
	huma.Register(api, huma.Operation{
		OperationID: "disable-strategy",
		Method:      http.MethodPost,
		Path:        basePath + "/strategy/{id}/disable",
		Summary:     "Disable strategy",
		Description: "Disables a specific trading strategy",
		Tags:        []string{"Strategy"},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id"`
	}) (*struct {
		Body struct {
			ID        string    `json:"id"`
			IsEnabled bool      `json:"isEnabled"`
			UpdatedAt time.Time `json:"updatedAt"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				ID        string    `json:"id"`
				IsEnabled bool      `json:"isEnabled"`
				UpdatedAt time.Time `json:"updatedAt"`
			}
		}{}
		resp.Body.ID = input.ID
		resp.Body.IsEnabled = false
		resp.Body.UpdatedAt = time.Now()
		return resp, nil
	})
}
