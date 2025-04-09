package huma

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// registerStrategyEndpoints registers the strategy endpoints.
func registerStrategyEndpoints(api huma.API, basePath string) {
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
			Strategies []struct {
				ID          string  `json:"id"`
				Name        string  `json:"name"`
				Description string  `json:"description"`
				WinRate     float64 `json:"winRate"`
				IsEnabled   bool    `json:"isEnabled"`
			} `json:"strategies"`
			Count     int       `json:"count"`
			Timestamp time.Time `json:"timestamp"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Strategies []struct {
					ID          string  `json:"id"`
					Name        string  `json:"name"`
					Description string  `json:"description"`
					WinRate     float64 `json:"winRate"`
					IsEnabled   bool    `json:"isEnabled"`
				} `json:"strategies"`
				Count     int       `json:"count"`
				Timestamp time.Time `json:"timestamp"`
			}
		}{}

		// Add some sample strategies
		resp.Body.Strategies = []struct {
			ID          string  `json:"id"`
			Name        string  `json:"name"`
			Description string  `json:"description"`
			WinRate     float64 `json:"winRate"`
			IsEnabled   bool    `json:"isEnabled"`
		}{
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
			filteredStrategies := []struct {
				ID          string  `json:"id"`
				Name        string  `json:"name"`
				Description string  `json:"description"`
				WinRate     float64 `json:"winRate"`
				IsEnabled   bool    `json:"isEnabled"`
			}{}
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
		Body struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Parameters  []struct {
				Name        string      `json:"name"`
				Type        string      `json:"type"`
				Description string      `json:"description"`
				Default     interface{} `json:"default"`
				Min         interface{} `json:"min,omitempty"`
				Max         interface{} `json:"max,omitempty"`
				Options     []string    `json:"options,omitempty"`
				Required    bool        `json:"required"`
			} `json:"parameters"`
			Performance struct {
				WinRate      float64 `json:"winRate"`
				ProfitFactor float64 `json:"profitFactor"`
				SharpeRatio  float64 `json:"sharpeRatio"`
				MaxDrawdown  float64 `json:"maxDrawdown"`
			} `json:"performance"`
			IsEnabled bool      `json:"isEnabled"`
			CreatedAt time.Time `json:"createdAt"`
			UpdatedAt time.Time `json:"updatedAt"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Parameters  []struct {
					Name        string      `json:"name"`
					Type        string      `json:"type"`
					Description string      `json:"description"`
					Default     interface{} `json:"default"`
					Min         interface{} `json:"min,omitempty"`
					Max         interface{} `json:"max,omitempty"`
					Options     []string    `json:"options,omitempty"`
					Required    bool        `json:"required"`
				} `json:"parameters"`
				Performance struct {
					WinRate      float64 `json:"winRate"`
					ProfitFactor float64 `json:"profitFactor"`
					SharpeRatio  float64 `json:"sharpeRatio"`
					MaxDrawdown  float64 `json:"maxDrawdown"`
				} `json:"performance"`
				IsEnabled bool      `json:"isEnabled"`
				CreatedAt time.Time `json:"createdAt"`
				UpdatedAt time.Time `json:"updatedAt"`
			}
		}{}

		// Set strategy details based on ID
		switch input.ID {
		case "breakout":
			resp.Body.ID = "breakout"
			resp.Body.Name = "Breakout Strategy"
			resp.Body.Description = "A strategy that trades breakouts from support and resistance levels"
			resp.Body.Parameters = []struct {
				Name        string      `json:"name"`
				Type        string      `json:"type"`
				Description string      `json:"description"`
				Default     interface{} `json:"default"`
				Min         interface{} `json:"min,omitempty"`
				Max         interface{} `json:"max,omitempty"`
				Options     []string    `json:"options,omitempty"`
				Required    bool        `json:"required"`
			}{
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
			resp.Body.Performance.WinRate = 65.0
			resp.Body.Performance.ProfitFactor = 2.1
			resp.Body.Performance.SharpeRatio = 1.5
			resp.Body.Performance.MaxDrawdown = 15.0
			resp.Body.IsEnabled = true
		case "volume-spike":
			resp.Body.ID = "volume-spike"
			resp.Body.Name = "Volume Spike Strategy"
			resp.Body.Description = "A strategy that trades based on unusual volume spikes"
			resp.Body.Parameters = []struct {
				Name        string      `json:"name"`
				Type        string      `json:"type"`
				Description string      `json:"description"`
				Default     interface{} `json:"default"`
				Min         interface{} `json:"min,omitempty"`
				Max         interface{} `json:"max,omitempty"`
				Options     []string    `json:"options,omitempty"`
				Required    bool        `json:"required"`
			}{
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
			resp.Body.Performance.WinRate = 58.0
			resp.Body.Performance.ProfitFactor = 1.8
			resp.Body.Performance.SharpeRatio = 1.2
			resp.Body.Performance.MaxDrawdown = 18.0
			resp.Body.IsEnabled = false
		case "new-coin":
			resp.Body.ID = "new-coin"
			resp.Body.Name = "New Coin Strategy"
			resp.Body.Description = "A strategy that trades newly listed coins"
			resp.Body.Parameters = []struct {
				Name        string      `json:"name"`
				Type        string      `json:"type"`
				Description string      `json:"description"`
				Default     interface{} `json:"default"`
				Min         interface{} `json:"min,omitempty"`
				Max         interface{} `json:"max,omitempty"`
				Options     []string    `json:"options,omitempty"`
				Required    bool        `json:"required"`
			}{
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
			resp.Body.Performance.WinRate = 72.0
			resp.Body.Performance.ProfitFactor = 2.5
			resp.Body.Performance.SharpeRatio = 1.8
			resp.Body.Performance.MaxDrawdown = 12.0
			resp.Body.IsEnabled = true
		default:
			// Return a default strategy if ID not found
			resp.Body.ID = input.ID
			resp.Body.Name = "Unknown Strategy"
			resp.Body.Description = "Strategy not found"
			resp.Body.IsEnabled = false
		}

		resp.Body.CreatedAt = time.Now().AddDate(0, -3, 0)
		resp.Body.UpdatedAt = time.Now().AddDate(0, 0, -5)
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
