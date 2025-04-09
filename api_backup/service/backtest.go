// Package service provides the service layer for the API
package service

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/pkg/backtest"
)

// BacktestService provides backtest functionality for the API
type BacktestService struct {
	backtestService backtest.Service
}

// NewBacktestService creates a new backtest service
func NewBacktestService(backtestService *backtest.Service) *BacktestService {
	return &BacktestService{
		backtestService: *backtestService,
	}
}

// BacktestRequest represents a request to run a backtest
type BacktestRequest struct {
	Strategy       string    `json:"strategy"`
	Symbol         string    `json:"symbol"`
	Timeframe      string    `json:"timeframe"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	InitialCapital float64   `json:"initialCapital"`
	RiskPerTrade   float64   `json:"riskPerTrade"`
}

// BacktestResult represents the result of a backtest
type BacktestResult struct {
	ID                 string    `json:"id"`
	Strategy           string    `json:"strategy"`
	Symbol             string    `json:"symbol"`
	Timeframe          string    `json:"timeframe"`
	StartDate          time.Time `json:"startDate"`
	EndDate            time.Time `json:"endDate"`
	InitialCapital     float64   `json:"initialCapital"`
	FinalCapital       float64   `json:"finalCapital"`
	TotalReturn        float64   `json:"totalReturn"`
	AnnualizedReturn   float64   `json:"annualizedReturn"`
	MaxDrawdown        float64   `json:"maxDrawdown"`
	SharpeRatio        float64   `json:"sharpeRatio"`
	WinRate            float64   `json:"winRate"`
	ProfitFactor       float64   `json:"profitFactor"`
	TotalTrades        int       `json:"totalTrades"`
	WinningTrades      int       `json:"winningTrades"`
	LosingTrades       int       `json:"losingTrades"`
	AverageProfitTrade float64   `json:"averageProfitTrade"`
	AverageLossTrade   float64   `json:"averageLossTrade"`
	MaxConsecutiveWins int       `json:"maxConsecutiveWins"`
	MaxConsecutiveLoss int       `json:"maxConsecutiveLoss"`
	EquityCurve        []struct {
		Timestamp time.Time `json:"timestamp"`
		Equity    float64   `json:"equity"`
	} `json:"equityCurve"`
	DrawdownCurve []struct {
		Timestamp time.Time `json:"timestamp"`
		Drawdown  float64   `json:"drawdown"`
	} `json:"drawdownCurve"`
	Trades []struct {
		ID        string    `json:"id"`
		Timestamp time.Time `json:"timestamp"`
		Side      string    `json:"side"`
		Price     float64   `json:"price"`
		Quantity  float64   `json:"quantity"`
		Profit    float64   `json:"profit"`
	} `json:"trades"`
	CreatedAt time.Time `json:"createdAt"`
}

// RunBacktest runs a backtest with the given configuration
func (s *BacktestService) RunBacktest(ctx context.Context, req *BacktestRequest) (*BacktestResult, error) {
	// Create a simplified implementation that doesn't depend on internal types
	params := make(map[string]interface{})
	params["riskPerTrade"] = req.RiskPerTrade

	// Create a mock result - we don't actually call the backtest service
	// This is a temporary solution until we have a proper implementation
	// that doesn't rely on internal packages

	// Create a mock result for now
	result := &BacktestResult{
		ID:                 fmt.Sprintf("bt-%d", time.Now().Unix()),
		Strategy:           req.Strategy,
		Symbol:             req.Symbol,
		Timeframe:          req.Timeframe,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		InitialCapital:     req.InitialCapital,
		FinalCapital:       req.InitialCapital * 1.15, // Mock 15% profit
		TotalReturn:        15.0,
		AnnualizedReturn:   20.0,
		MaxDrawdown:        5.0,
		SharpeRatio:        1.5,
		WinRate:            65.0,
		ProfitFactor:       2.1,
		TotalTrades:        25,
		WinningTrades:      16,
		LosingTrades:       9,
		AverageProfitTrade: 2.5,
		AverageLossTrade:   -1.2,
		MaxConsecutiveWins: 5,
		MaxConsecutiveLoss: 2,
		EquityCurve: []struct {
			Timestamp time.Time `json:"timestamp"`
			Equity    float64   `json:"equity"`
		}{
			{Timestamp: req.StartDate, Equity: req.InitialCapital},
			{Timestamp: req.EndDate, Equity: req.InitialCapital * 1.15},
		},
		DrawdownCurve: []struct {
			Timestamp time.Time `json:"timestamp"`
			Drawdown  float64   `json:"drawdown"`
		}{
			{Timestamp: req.StartDate, Drawdown: 0},
			{Timestamp: req.EndDate, Drawdown: 0},
		},
		Trades: []struct {
			ID        string    `json:"id"`
			Timestamp time.Time `json:"timestamp"`
			Side      string    `json:"side"`
			Price     float64   `json:"price"`
			Quantity  float64   `json:"quantity"`
			Profit    float64   `json:"profit"`
		}{
			{
				ID:        "trade-1",
				Timestamp: req.StartDate.Add(24 * time.Hour),
				Side:      "BUY",
				Price:     100.0,
				Quantity:  1.0,
				Profit:    0.0,
			},
			{
				ID:        "trade-2",
				Timestamp: req.StartDate.Add(48 * time.Hour),
				Side:      "SELL",
				Price:     110.0,
				Quantity:  1.0,
				Profit:    10.0,
			},
		},
		CreatedAt: time.Now(),
	}

	return result, nil
}

// GetBacktestResult gets a backtest result by ID
func (s *BacktestService) GetBacktestResult(ctx context.Context, id string) (*BacktestResult, error) {
	// This is a simplified mock implementation
	// In a real implementation, we would call the backtest service to get the result

	// Create a mock result
	result := &BacktestResult{
		ID:                 id,
		Strategy:           "moving_average",
		Symbol:             "BTC/USDT",
		Timeframe:          "1h",
		StartDate:          time.Now().AddDate(0, -1, 0),
		EndDate:            time.Now(),
		InitialCapital:     10000.0,
		FinalCapital:       11500.0,
		TotalReturn:        15.0,
		AnnualizedReturn:   20.0,
		MaxDrawdown:        5.0,
		SharpeRatio:        1.5,
		WinRate:            65.0,
		ProfitFactor:       2.1,
		TotalTrades:        25,
		WinningTrades:      16,
		LosingTrades:       9,
		AverageProfitTrade: 2.5,
		AverageLossTrade:   -1.2,
		MaxConsecutiveWins: 5,
		MaxConsecutiveLoss: 2,
		CreatedAt:          time.Now().AddDate(0, 0, -1),
	}

	return result, nil
}

// ListBacktestResults lists all backtest results
func (s *BacktestService) ListBacktestResults(ctx context.Context) ([]*BacktestResult, error) {
	// This is a simplified mock implementation
	// In a real implementation, we would call the backtest service to get the results

	// Create mock results
	results := []*BacktestResult{
		{
			ID:                 "bt-1",
			Strategy:           "moving_average",
			Symbol:             "BTC/USDT",
			Timeframe:          "1h",
			StartDate:          time.Now().AddDate(0, -1, 0),
			EndDate:            time.Now(),
			InitialCapital:     10000.0,
			FinalCapital:       11500.0,
			TotalReturn:        15.0,
			AnnualizedReturn:   20.0,
			MaxDrawdown:        5.0,
			SharpeRatio:        1.5,
			WinRate:            65.0,
			ProfitFactor:       2.1,
			TotalTrades:        25,
			WinningTrades:      16,
			LosingTrades:       9,
			AverageProfitTrade: 2.5,
			AverageLossTrade:   -1.2,
			MaxConsecutiveWins: 5,
			MaxConsecutiveLoss: 2,
			CreatedAt:          time.Now().AddDate(0, 0, -1),
		},
		{
			ID:                 "bt-2",
			Strategy:           "breakout",
			Symbol:             "ETH/USDT",
			Timeframe:          "4h",
			StartDate:          time.Now().AddDate(0, -2, 0),
			EndDate:            time.Now(),
			InitialCapital:     10000.0,
			FinalCapital:       12000.0,
			TotalReturn:        20.0,
			AnnualizedReturn:   25.0,
			MaxDrawdown:        8.0,
			SharpeRatio:        1.8,
			WinRate:            70.0,
			ProfitFactor:       2.5,
			TotalTrades:        30,
			WinningTrades:      21,
			LosingTrades:       9,
			AverageProfitTrade: 3.0,
			AverageLossTrade:   -1.5,
			MaxConsecutiveWins: 7,
			MaxConsecutiveLoss: 3,
			CreatedAt:          time.Now().AddDate(0, 0, -2),
		},
	}

	return results, nil
}

// BacktestComparisonResult represents the result of comparing multiple backtests
type BacktestComparisonResult struct {
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
	} `json:"backtests"`
	BestPerformer struct {
		ID               string  `json:"id"`
		Strategy         string  `json:"strategy"`
		TotalReturn      float64 `json:"totalReturn"`
		AnnualizedReturn float64 `json:"annualizedReturn"`
		SharpeRatio      float64 `json:"sharpeRatio"`
	} `json:"bestPerformer"`
	WorstPerformer struct {
		ID               string  `json:"id"`
		Strategy         string  `json:"strategy"`
		TotalReturn      float64 `json:"totalReturn"`
		AnnualizedReturn float64 `json:"annualizedReturn"`
		SharpeRatio      float64 `json:"sharpeRatio"`
	} `json:"worstPerformer"`
	ComparisonDate time.Time `json:"comparisonDate"`
}

// CompareBacktests compares multiple backtests
func (s *BacktestService) CompareBacktests(ctx context.Context, ids []string) (*BacktestComparisonResult, error) {
	// Get backtest results
	results := make([]*BacktestResult, 0, len(ids))
	for _, id := range ids {
		result, err := s.GetBacktestResult(ctx, id)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	// Create comparison result
	comparison := &BacktestComparisonResult{
		Backtests: make([]struct {
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
		}, 0, len(results)),
		ComparisonDate: time.Now(),
	}

	// Find best and worst performers
	var bestPerformer, worstPerformer *BacktestResult
	for i, result := range results {
		// Add to backtests
		comparison.Backtests = append(comparison.Backtests, struct {
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
		}{
			ID:               result.ID,
			Strategy:         result.Strategy,
			Symbol:           result.Symbol,
			Timeframe:        result.Timeframe,
			TotalReturn:      result.TotalReturn,
			AnnualizedReturn: result.AnnualizedReturn,
			MaxDrawdown:      result.MaxDrawdown,
			SharpeRatio:      result.SharpeRatio,
			WinRate:          result.WinRate,
			ProfitFactor:     result.ProfitFactor,
		})

		// Update best and worst performers
		if i == 0 || result.SharpeRatio > bestPerformer.SharpeRatio {
			bestPerformer = result
		}
		if i == 0 || result.SharpeRatio < worstPerformer.SharpeRatio {
			worstPerformer = result
		}
	}

	// Set best and worst performers
	comparison.BestPerformer = struct {
		ID               string  `json:"id"`
		Strategy         string  `json:"strategy"`
		TotalReturn      float64 `json:"totalReturn"`
		AnnualizedReturn float64 `json:"annualizedReturn"`
		SharpeRatio      float64 `json:"sharpeRatio"`
	}{
		ID:               bestPerformer.ID,
		Strategy:         bestPerformer.Strategy,
		TotalReturn:      bestPerformer.TotalReturn,
		AnnualizedReturn: bestPerformer.AnnualizedReturn,
		SharpeRatio:      bestPerformer.SharpeRatio,
	}

	comparison.WorstPerformer = struct {
		ID               string  `json:"id"`
		Strategy         string  `json:"strategy"`
		TotalReturn      float64 `json:"totalReturn"`
		AnnualizedReturn float64 `json:"annualizedReturn"`
		SharpeRatio      float64 `json:"sharpeRatio"`
	}{
		ID:               worstPerformer.ID,
		Strategy:         worstPerformer.Strategy,
		TotalReturn:      worstPerformer.TotalReturn,
		AnnualizedReturn: worstPerformer.AnnualizedReturn,
		SharpeRatio:      worstPerformer.SharpeRatio,
	}

	return comparison, nil
}
