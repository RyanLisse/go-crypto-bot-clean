// Package strategy provides the strategy endpoints for the Huma API.
package strategy

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"go-crypto-bot-clean/api/service"
)

// RegisterEndpoints registers the strategy endpoints.
func RegisterEndpoints(api huma.API, basePath string, strategyService *service.StrategyService) {
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

		// Convert enabled string to bool pointer
		var enabledPtr *bool
		if input.Enabled != "" {
			enabled := input.Enabled == "true"
			enabledPtr = &enabled
		}

		// Get strategies from service
		strategies, err := strategyService.ListStrategies(ctx, enabledPtr)
		if err != nil {
			return nil, err
		}

		// Convert service strategies to API strategies
		apiStrategies := make([]StrategyListItem, 0, len(strategies))
		for _, strategy := range strategies {
			apiStrategies = append(apiStrategies, StrategyListItem{
				ID:          strategy.ID,
				Name:        strategy.Name,
				Description: strategy.Description,
				WinRate:     strategy.WinRate,
				IsEnabled:   strategy.IsEnabled,
			})
		}

		resp.Body.Strategies = apiStrategies
		resp.Body.Count = len(apiStrategies)
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

		// Get strategy from service
		strategy, err := strategyService.GetStrategy(ctx, input.ID)
		if err != nil {
			return nil, err
		}

		// Convert service strategy to API strategy
		resp.Body = Strategy{
			ID:          strategy.ID,
			Name:        strategy.Name,
			Description: strategy.Description,
			IsEnabled:   strategy.IsEnabled,
			CreatedAt:   strategy.CreatedAt,
			UpdatedAt:   strategy.UpdatedAt,
			Performance: Performance{
				WinRate:      strategy.Performance.WinRate,
				ProfitFactor: strategy.Performance.ProfitFactor,
				SharpeRatio:  strategy.Performance.SharpeRatio,
				MaxDrawdown:  strategy.Performance.MaxDrawdown,
			},
		}

		// Convert service parameters to API parameters
		resp.Body.Parameters = make([]Parameter, 0, len(strategy.Parameters))
		for _, param := range strategy.Parameters {
			resp.Body.Parameters = append(resp.Body.Parameters, Parameter{
				Name:        param.Name,
				Type:        param.Type,
				Description: param.Description,
				Default:     param.Default,
				Min:         param.Min,
				Max:         param.Max,
				Options:     param.Options,
				Required:    param.Required,
			})
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

		// Update strategy configuration
		strategy, err := strategyService.UpdateStrategyConfig(ctx, input.ID, input.Body.Parameters, input.Body.IsEnabled)
		if err != nil {
			return nil, err
		}

		// Set response
		resp.Body.ID = strategy.ID
		resp.Body.Parameters = input.Body.Parameters
		resp.Body.IsEnabled = strategy.IsEnabled
		resp.Body.UpdatedAt = strategy.UpdatedAt

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

		// Enable strategy
		result, err := strategyService.EnableStrategy(ctx, input.ID)
		if err != nil {
			return nil, err
		}

		// Set response
		resp.Body.ID = result.ID
		resp.Body.IsEnabled = result.IsEnabled
		resp.Body.UpdatedAt = result.UpdatedAt

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

		// Disable strategy
		result, err := strategyService.DisableStrategy(ctx, input.ID)
		if err != nil {
			return nil, err
		}

		// Set response
		resp.Body.ID = result.ID
		resp.Body.IsEnabled = result.IsEnabled
		resp.Body.UpdatedAt = result.UpdatedAt

		return resp, nil
	})
}
