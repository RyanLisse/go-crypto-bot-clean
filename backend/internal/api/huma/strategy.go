package huma

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// StrategyResponse represents a trading strategy
type StrategyResponse struct {
	Body struct {
		ID          string    `json:"id" doc:"Unique identifier for the strategy" example:"breakout"`
		Name        string    `json:"name" doc:"Name of the strategy" example:"Breakout Strategy"`
		Description string    `json:"description" doc:"Description of the strategy" example:"A strategy that trades breakouts from support and resistance levels"`
		Parameters  []struct {
			Name        string      `json:"name" doc:"Name of the parameter" example:"lookbackPeriod"`
			Type        string      `json:"type" doc:"Type of the parameter" example:"integer" enum:"integer,float,boolean,string,array"`
			Description string      `json:"description" doc:"Description of the parameter" example:"Number of periods to look back for support/resistance"`
			Default     interface{} `json:"default" doc:"Default value of the parameter" example:"20"`
			Min         interface{} `json:"min,omitempty" doc:"Minimum value of the parameter" example:"5"`
			Max         interface{} `json:"max,omitempty" doc:"Maximum value of the parameter" example:"100"`
			Options     []string    `json:"options,omitempty" doc:"Available options for the parameter" example:"[\"option1\", \"option2\"]"`
			Required    bool        `json:"required" doc:"Whether the parameter is required" example:"true"`
		} `json:"parameters" doc:"Parameters of the strategy"`
		Performance struct {
			WinRate      float64 `json:"winRate" doc:"Win rate percentage" example:"65.0"`
			ProfitFactor float64 `json:"profitFactor" doc:"Profit factor" example:"2.1"`
			SharpeRatio  float64 `json:"sharpeRatio" doc:"Sharpe ratio" example:"1.5"`
			MaxDrawdown  float64 `json:"maxDrawdown" doc:"Maximum drawdown percentage" example:"15.0"`
		} `json:"performance" doc:"Performance metrics of the strategy"`
		IsEnabled  bool      `json:"isEnabled" doc:"Whether the strategy is enabled" example:"true"`
		CreatedAt  time.Time `json:"createdAt" doc:"Time when the strategy was created" example:"2023-01-01T00:00:00Z"`
		UpdatedAt  time.Time `json:"updatedAt" doc:"Time when the strategy was last updated" example:"2023-01-02T00:00:00Z"`
	}
}

// StrategyListResponse represents a list of trading strategies
type StrategyListResponse struct {
	Body struct {
		Strategies []struct {
			ID          string  `json:"id" doc:"Unique identifier for the strategy" example:"breakout"`
			Name        string  `json:"name" doc:"Name of the strategy" example:"Breakout Strategy"`
			Description string  `json:"description" doc:"Description of the strategy" example:"A strategy that trades breakouts from support and resistance levels"`
			WinRate     float64 `json:"winRate" doc:"Win rate percentage" example:"65.0"`
			IsEnabled   bool    `json:"isEnabled" doc:"Whether the strategy is enabled" example:"true"`
		} `json:"strategies" doc:"List of strategies"`
		Count     int       `json:"count" doc:"Number of strategies" example:"3"`
		Timestamp time.Time `json:"timestamp" doc:"Timestamp of the response" example:"2023-02-02T10:00:00Z"`
	}
}

// StrategyConfigRequest represents a request to update strategy configuration
type StrategyConfigRequest struct {
	Body struct {
		Parameters map[string]interface{} `json:"parameters" doc:"Parameters of the strategy" example:"{\"lookbackPeriod\": 20, \"breakoutThreshold\": 2.5}" binding:"required"`
		IsEnabled  bool                   `json:"isEnabled" doc:"Whether the strategy is enabled" example:"true"`
	}
}

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
		Enabled *bool `query:"enabled" doc:"Filter by enabled status" example:"true"`
	}) (*StrategyListResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
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
		ID string `path:"id" doc:"Strategy ID" example:"breakout"`
	}) (*StrategyResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
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
		ID string `path:"id" doc:"Strategy ID" example:"breakout"`
		StrategyConfigRequest
	}) (*StrategyResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
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
		ID string `path:"id" doc:"Strategy ID" example:"breakout"`
	}) (*struct {
		Body struct {
			ID        string    `json:"id" doc:"Unique identifier for the strategy" example:"breakout"`
			IsEnabled bool      `json:"isEnabled" doc:"Whether the strategy is enabled" example:"true"`
			UpdatedAt time.Time `json:"updatedAt" doc:"Time when the strategy was updated" example:"2023-02-02T10:00:00Z"`
		}
	}, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
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
		ID string `path:"id" doc:"Strategy ID" example:"breakout"`
	}) (*struct {
		Body struct {
			ID        string    `json:"id" doc:"Unique identifier for the strategy" example:"breakout"`
			IsEnabled bool      `json:"isEnabled" doc:"Whether the strategy is enabled" example:"false"`
			UpdatedAt time.Time `json:"updatedAt" doc:"Time when the strategy was updated" example:"2023-02-02T10:00:00Z"`
		}
	}, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})
}
