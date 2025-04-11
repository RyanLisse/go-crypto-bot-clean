package huma

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// POST /backtest/compare
func registerBacktestCompare(api huma.API, basePath string) {
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
