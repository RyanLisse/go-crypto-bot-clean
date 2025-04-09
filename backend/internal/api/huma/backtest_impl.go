package huma

import (
	"context"
	"net/http"
	"time"

	"go-crypto-bot-clean/backend/internal/api/service"

	"github.com/danielgtaylor/huma/v2"
)

// registerBacktestEndpointsWithService registers the backtest endpoints with service implementation.
func registerBacktestEndpointsWithService(api huma.API, basePath string, services *service.Provider) {
	// POST /backtest
	huma.Register(api, huma.Operation{
		OperationID: "run-backtest",
		Method:      http.MethodPost,
		Path:        basePath + "/backtest",
		Summary:     "Run backtest",
		Description: "Runs a backtest with the given configuration",
		Tags:        []string{"Backtest"},
	}, func(ctx context.Context, input *struct {
		Strategy       string    `json:"strategy"`
		Symbol         string    `json:"symbol"`
		Timeframe      string    `json:"timeframe"`
		StartDate      time.Time `json:"startDate"`
		EndDate        time.Time `json:"endDate"`
		InitialCapital float64   `json:"initialCapital"`
		RiskPerTrade   float64   `json:"riskPerTrade"`
	}) (*struct {
		Body *service.BacktestResult
	}, error) {
		// Convert API request to service request
		req := &service.BacktestRequest{
			Strategy:       input.Strategy,
			Symbol:         input.Symbol,
			Timeframe:      input.Timeframe,
			StartDate:      input.StartDate,
			EndDate:        input.EndDate,
			InitialCapital: input.InitialCapital,
			RiskPerTrade:   input.RiskPerTrade,
		}

		// Run backtest
		result, err := services.BacktestService.RunBacktest(ctx, req)
		if err != nil {
			return nil, err
		}

		// Return result
		return &struct {
			Body *service.BacktestResult
		}{
			Body: result,
		}, nil
	})

	// GET /backtest/{id}
	huma.Register(api, huma.Operation{
		OperationID: "get-backtest",
		Method:      http.MethodGet,
		Path:        basePath + "/backtest/{id}",
		Summary:     "Get backtest",
		Description: "Gets a backtest result by ID",
		Tags:        []string{"Backtest"},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id"`
	}) (*struct {
		Body *service.BacktestResult
	}, error) {
		// Get backtest result
		result, err := services.BacktestService.GetBacktestResult(ctx, input.ID)
		if err != nil {
			return nil, err
		}

		// Return result
		return &struct {
			Body *service.BacktestResult
		}{
			Body: result,
		}, nil
	})

	// GET /backtest/list
	huma.Register(api, huma.Operation{
		OperationID: "list-backtests",
		Method:      http.MethodGet,
		Path:        basePath + "/backtest/list",
		Summary:     "List backtests",
		Description: "Lists all backtest results",
		Tags:        []string{"Backtest"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body []*service.BacktestResult
	}, error) {
		// List backtest results
		results, err := services.BacktestService.ListBacktestResults(ctx)
		if err != nil {
			return nil, err
		}

		// Return results
		return &struct {
			Body []*service.BacktestResult
		}{
			Body: results,
		}, nil
	})

	// POST /backtest/compare
	huma.Register(api, huma.Operation{
		OperationID: "compare-backtests",
		Method:      http.MethodPost,
		Path:        basePath + "/backtest/compare",
		Summary:     "Compare backtests",
		Description: "Compares multiple backtests",
		Tags:        []string{"Backtest"},
	}, func(ctx context.Context, input *struct {
		IDs []string `json:"ids"`
	}) (*struct {
		Body *service.BacktestComparisonResult
	}, error) {
		// Compare backtests
		result, err := services.BacktestService.CompareBacktests(ctx, input.IDs)
		if err != nil {
			return nil, err
		}

		// Return result
		return &struct {
			Body *service.BacktestComparisonResult
		}{
			Body: result,
		}, nil
	})
}
