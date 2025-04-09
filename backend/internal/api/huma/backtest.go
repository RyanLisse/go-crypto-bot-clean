package huma

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// registerBacktestEndpoints registers the backtest endpoints.
// Deprecated: Use registerBacktestEndpointsWithService instead.
func registerBacktestEndpoints(api huma.API, basePath string) {
	// POST /backtest
	huma.Register(api, huma.Operation{
		OperationID: "run-backtest",
		Method:      http.MethodPost,
		Path:        basePath + "/backtest",
		Summary:     "Run a backtest",
		Description: "Runs a backtest with the specified parameters",
		Tags:        []string{"Backtesting"},
	}, func(ctx context.Context, input *struct {
		Strategy       string    `json:"strategy"`
		Symbol         string    `json:"symbol"`
		Timeframe      string    `json:"timeframe"`
		StartDate      time.Time `json:"startDate"`
		EndDate        time.Time `json:"endDate"`
		InitialCapital float64   `json:"initialCapital"`
		RiskPerTrade   float64   `json:"riskPerTrade"`
	}) (*struct {
		Body struct {
			ID        string    `json:"id"`
			Status    string    `json:"status"`
			Message   string    `json:"message,omitempty"`
			CreatedAt time.Time `json:"createdAt"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				ID        string    `json:"id"`
				Status    string    `json:"status"`
				Message   string    `json:"message,omitempty"`
				CreatedAt time.Time `json:"createdAt"`
			}
		}{}
		resp.Body.ID = "bt-" + uuid.New().String()[:6]
		resp.Body.Status = "completed"
		resp.Body.Message = "Backtest completed successfully"
		resp.Body.CreatedAt = time.Now()
		return resp, nil
	})

	// GET /backtest/{id}
	huma.Register(api, huma.Operation{
		OperationID: "get-backtest-result",
		Method:      http.MethodGet,
		Path:        basePath + "/backtest/{id}",
		Summary:     "Get backtest result",
		Description: "Returns the result of a specific backtest",
		Tags:        []string{"Backtesting"},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id"`
	}) (*struct {
		Body struct {
			ID               string    `json:"id"`
			Strategy         string    `json:"strategy"`
			Symbol           string    `json:"symbol"`
			Timeframe        string    `json:"timeframe"`
			StartDate        time.Time `json:"startDate"`
			EndDate          time.Time `json:"endDate"`
			InitialCapital   float64   `json:"initialCapital"`
			FinalCapital     float64   `json:"finalCapital"`
			TotalReturn      float64   `json:"totalReturn"`
			AnnualizedReturn float64   `json:"annualizedReturn"`
			MaxDrawdown      float64   `json:"maxDrawdown"`
			SharpeRatio      float64   `json:"sharpeRatio"`
			WinRate          float64   `json:"winRate"`
			ProfitFactor     float64   `json:"profitFactor"`
			TotalTrades      int       `json:"totalTrades"`
			WinningTrades    int       `json:"winningTrades"`
			LosingTrades     int       `json:"losingTrades"`
			CreatedAt        time.Time `json:"createdAt"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				ID               string    `json:"id"`
				Strategy         string    `json:"strategy"`
				Symbol           string    `json:"symbol"`
				Timeframe        string    `json:"timeframe"`
				StartDate        time.Time `json:"startDate"`
				EndDate          time.Time `json:"endDate"`
				InitialCapital   float64   `json:"initialCapital"`
				FinalCapital     float64   `json:"finalCapital"`
				TotalReturn      float64   `json:"totalReturn"`
				AnnualizedReturn float64   `json:"annualizedReturn"`
				MaxDrawdown      float64   `json:"maxDrawdown"`
				SharpeRatio      float64   `json:"sharpeRatio"`
				WinRate          float64   `json:"winRate"`
				ProfitFactor     float64   `json:"profitFactor"`
				TotalTrades      int       `json:"totalTrades"`
				WinningTrades    int       `json:"winningTrades"`
				LosingTrades     int       `json:"losingTrades"`
				CreatedAt        time.Time `json:"createdAt"`
			}
		}{}
		resp.Body.ID = input.ID
		resp.Body.Strategy = "breakout"
		resp.Body.Symbol = "BTC/USDT"
		resp.Body.Timeframe = "1h"
		resp.Body.StartDate = time.Now().AddDate(0, -1, 0)
		resp.Body.EndDate = time.Now()
		resp.Body.InitialCapital = 1000.0
		resp.Body.FinalCapital = 1250.0
		resp.Body.TotalReturn = 25.0
		resp.Body.AnnualizedReturn = 300.0
		resp.Body.MaxDrawdown = 15.0
		resp.Body.SharpeRatio = 1.5
		resp.Body.WinRate = 65.0
		resp.Body.ProfitFactor = 2.1
		resp.Body.TotalTrades = 50
		resp.Body.WinningTrades = 32
		resp.Body.LosingTrades = 18
		resp.Body.CreatedAt = time.Now()
		return resp, nil
	})

	// GET /backtest/list
	huma.Register(api, huma.Operation{
		OperationID: "list-backtests",
		Method:      http.MethodGet,
		Path:        basePath + "/backtest/list",
		Summary:     "List backtests",
		Description: "Returns a list of backtests",
		Tags:        []string{"Backtesting"},
	}, func(ctx context.Context, input *struct {
		Limit    int    `query:"limit"`
		Offset   int    `query:"offset"`
		Strategy string `query:"strategy"`
		Symbol   string `query:"symbol"`
	}) (*struct {
		Body struct {
			Backtests []struct {
				ID          string    `json:"id"`
				Strategy    string    `json:"strategy"`
				Symbol      string    `json:"symbol"`
				Timeframe   string    `json:"timeframe"`
				StartDate   time.Time `json:"startDate"`
				EndDate     time.Time `json:"endDate"`
				TotalReturn float64   `json:"totalReturn"`
				MaxDrawdown float64   `json:"maxDrawdown"`
				SharpeRatio float64   `json:"sharpeRatio"`
				WinRate     float64   `json:"winRate"`
				TotalTrades int       `json:"totalTrades"`
				CreatedAt   time.Time `json:"createdAt"`
			} `json:"backtests"`
			Count     int       `json:"count"`
			Timestamp time.Time `json:"timestamp"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Backtests []struct {
					ID          string    `json:"id"`
					Strategy    string    `json:"strategy"`
					Symbol      string    `json:"symbol"`
					Timeframe   string    `json:"timeframe"`
					StartDate   time.Time `json:"startDate"`
					EndDate     time.Time `json:"endDate"`
					TotalReturn float64   `json:"totalReturn"`
					MaxDrawdown float64   `json:"maxDrawdown"`
					SharpeRatio float64   `json:"sharpeRatio"`
					WinRate     float64   `json:"winRate"`
					TotalTrades int       `json:"totalTrades"`
					CreatedAt   time.Time `json:"createdAt"`
				} `json:"backtests"`
				Count     int       `json:"count"`
				Timestamp time.Time `json:"timestamp"`
			}
		}{}
		resp.Body.Backtests = []struct {
			ID          string    `json:"id"`
			Strategy    string    `json:"strategy"`
			Symbol      string    `json:"symbol"`
			Timeframe   string    `json:"timeframe"`
			StartDate   time.Time `json:"startDate"`
			EndDate     time.Time `json:"endDate"`
			TotalReturn float64   `json:"totalReturn"`
			MaxDrawdown float64   `json:"maxDrawdown"`
			SharpeRatio float64   `json:"sharpeRatio"`
			WinRate     float64   `json:"winRate"`
			TotalTrades int       `json:"totalTrades"`
			CreatedAt   time.Time `json:"createdAt"`
		}{
			{
				ID:          "bt-123456",
				Strategy:    "breakout",
				Symbol:      "BTC/USDT",
				Timeframe:   "1h",
				StartDate:   time.Now().AddDate(0, -1, 0),
				EndDate:     time.Now(),
				TotalReturn: 25.0,
				MaxDrawdown: 15.0,
				SharpeRatio: 1.5,
				WinRate:     65.0,
				TotalTrades: 50,
				CreatedAt:   time.Now(),
			},
		}
		resp.Body.Count = len(resp.Body.Backtests)
		resp.Body.Timestamp = time.Now()
		return resp, nil
	})

	// POST /backtest/compare
	huma.Register(api, huma.Operation{
		OperationID: "compare-backtests",
		Method:      http.MethodPost,
		Path:        basePath + "/backtest/compare",
		Summary:     "Compare backtests",
		Description: "Compares multiple backtests",
		Tags:        []string{"Backtesting"},
	}, func(ctx context.Context, input *struct {
		BacktestIDs []string `json:"backtest_ids"`
	}) (*struct {
		Body struct {
			Backtests []struct {
				ID               string  `json:"id"`
				Strategy         string  `json:"strategy"`
				Symbol           string  `json:"symbol"`
				Timeframe        string  `json:"timeframe"`
				TotalReturn      float64 `json:"totalReturn"`
				AnnualizedReturn float64 `json:"annualizedReturn"`
				MaxDrawdown      float64 `json:"maxDrawdown"`
				SharpeRatio      float64 `json:"sharpeRatio"`
				WinRate          float64 `json:"winRate"`
				ProfitFactor     float64 `json:"profitFactor"`
				TotalTrades      int     `json:"totalTrades"`
			} `json:"backtests"`
			Comparison struct {
				BestTotalReturn       string  `json:"bestTotalReturn"`
				BestSharpeRatio       string  `json:"bestSharpeRatio"`
				BestDrawdown          string  `json:"bestDrawdown"`
				BestWinRate           string  `json:"bestWinRate"`
				BestProfitFactor      string  `json:"bestProfitFactor"`
				ReturnDifference      float64 `json:"returnDifference"`
				DrawdownDifference    float64 `json:"drawdownDifference"`
				SharpeRatioDifference float64 `json:"sharpeRatioDifference"`
			} `json:"comparison"`
			Timestamp time.Time `json:"timestamp"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Backtests []struct {
					ID               string  `json:"id"`
					Strategy         string  `json:"strategy"`
					Symbol           string  `json:"symbol"`
					Timeframe        string  `json:"timeframe"`
					TotalReturn      float64 `json:"totalReturn"`
					AnnualizedReturn float64 `json:"annualizedReturn"`
					MaxDrawdown      float64 `json:"maxDrawdown"`
					SharpeRatio      float64 `json:"sharpeRatio"`
					WinRate          float64 `json:"winRate"`
					ProfitFactor     float64 `json:"profitFactor"`
					TotalTrades      int     `json:"totalTrades"`
				} `json:"backtests"`
				Comparison struct {
					BestTotalReturn       string  `json:"bestTotalReturn"`
					BestSharpeRatio       string  `json:"bestSharpeRatio"`
					BestDrawdown          string  `json:"bestDrawdown"`
					BestWinRate           string  `json:"bestWinRate"`
					BestProfitFactor      string  `json:"bestProfitFactor"`
					ReturnDifference      float64 `json:"returnDifference"`
					DrawdownDifference    float64 `json:"drawdownDifference"`
					SharpeRatioDifference float64 `json:"sharpeRatioDifference"`
				} `json:"comparison"`
				Timestamp time.Time `json:"timestamp"`
			}
		}{}
		resp.Body.Backtests = []struct {
			ID               string  `json:"id"`
			Strategy         string  `json:"strategy"`
			Symbol           string  `json:"symbol"`
			Timeframe        string  `json:"timeframe"`
			TotalReturn      float64 `json:"totalReturn"`
			AnnualizedReturn float64 `json:"annualizedReturn"`
			MaxDrawdown      float64 `json:"maxDrawdown"`
			SharpeRatio      float64 `json:"sharpeRatio"`
			WinRate          float64 `json:"winRate"`
			ProfitFactor     float64 `json:"profitFactor"`
			TotalTrades      int     `json:"totalTrades"`
		}{
			{
				ID:               input.BacktestIDs[0],
				Strategy:         "breakout",
				Symbol:           "BTC/USDT",
				Timeframe:        "1h",
				TotalReturn:      25.0,
				AnnualizedReturn: 300.0,
				MaxDrawdown:      15.0,
				SharpeRatio:      1.5,
				WinRate:          65.0,
				ProfitFactor:     2.1,
				TotalTrades:      50,
			},
		}
		resp.Body.Comparison.BestTotalReturn = input.BacktestIDs[0]
		resp.Body.Comparison.BestSharpeRatio = input.BacktestIDs[0]
		resp.Body.Comparison.BestDrawdown = input.BacktestIDs[0]
		resp.Body.Comparison.BestWinRate = input.BacktestIDs[0]
		resp.Body.Comparison.BestProfitFactor = input.BacktestIDs[0]
		resp.Body.Comparison.ReturnDifference = 0.0
		resp.Body.Comparison.DrawdownDifference = 0.0
		resp.Body.Comparison.SharpeRatioDifference = 0.0
		resp.Body.Timestamp = time.Now()
		return resp, nil
	})
}
