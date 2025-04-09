// Package service provides the service layer for the API
package service

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/backtest"
	"go-crypto-bot-clean/backend/internal/domain/models"
)

// BacktestService provides backtest functionality for the API
type BacktestService struct {
	backtestService *backtest.Service
}

// NewBacktestService creates a new backtest service
func NewBacktestService(backtestService *backtest.Service) *BacktestService {
	return &BacktestService{
		backtestService: backtestService,
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
	ID               string                 `json:"id"`
	Strategy         string                 `json:"strategy"`
	Symbol           string                 `json:"symbol"`
	Timeframe        string                 `json:"timeframe"`
	StartDate        time.Time              `json:"startDate"`
	EndDate          time.Time              `json:"endDate"`
	InitialCapital   float64                `json:"initialCapital"`
	FinalCapital     float64                `json:"finalCapital"`
	TotalReturn      float64                `json:"totalReturn"`
	AnnualizedReturn float64                `json:"annualizedReturn"`
	MaxDrawdown      float64                `json:"maxDrawdown"`
	SharpeRatio      float64                `json:"sharpeRatio"`
	WinRate          float64                `json:"winRate"`
	ProfitFactor     float64                `json:"profitFactor"`
	TotalTrades      int                    `json:"totalTrades"`
	WinningTrades    int                    `json:"winningTrades"`
	LosingTrades     int                    `json:"losingTrades"`
	AverageProfitTrade float64              `json:"averageProfitTrade"`
	AverageLossTrade   float64              `json:"averageLossTrade"`
	MaxConsecutiveWins int                  `json:"maxConsecutiveWins"`
	MaxConsecutiveLoss int                  `json:"maxConsecutiveLoss"`
	EquityCurve        []*backtest.EquityPoint   `json:"equityCurve"`
	DrawdownCurve      []*backtest.DrawdownPoint `json:"drawdownCurve"`
	Trades             []*models.Order           `json:"trades"`
	CreatedAt          time.Time                 `json:"createdAt"`
}

// RunBacktest runs a backtest with the given configuration
func (s *BacktestService) RunBacktest(ctx context.Context, req *BacktestRequest) (*BacktestResult, error) {
	// Convert API request to backtest service request
	backtestReq := &backtest.BacktestRequestConfig{
		Strategy:       req.Strategy,
		Symbol:         req.Symbol,
		Timeframe:      req.Timeframe,
		StartTime:      req.StartDate,
		EndTime:        req.EndDate,
		InitialCapital: req.InitialCapital,
		RiskPerTrade:   req.RiskPerTrade,
	}

	// Run backtest
	result, err := s.backtestService.RunBacktest(ctx, backtestReq)
	if err != nil {
		return nil, err
	}

	// Convert backtest result to API result
	apiResult := &BacktestResult{
		ID:               "bt-" + result.ID, // Add prefix for API
		Strategy:         req.Strategy,
		Symbol:           req.Symbol,
		Timeframe:        req.Timeframe,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		InitialCapital:   req.InitialCapital,
		FinalCapital:     result.FinalCapital,
		TotalReturn:      result.PerformanceMetrics.TotalReturn,
		AnnualizedReturn: result.PerformanceMetrics.AnnualizedReturn,
		MaxDrawdown:      result.PerformanceMetrics.MaxDrawdown,
		SharpeRatio:      result.PerformanceMetrics.SharpeRatio,
		WinRate:          result.PerformanceMetrics.WinRate,
		ProfitFactor:     result.PerformanceMetrics.ProfitFactor,
		TotalTrades:      result.PerformanceMetrics.TotalTrades,
		WinningTrades:    result.PerformanceMetrics.WinningTrades,
		LosingTrades:     result.PerformanceMetrics.LosingTrades,
		AverageProfitTrade: result.PerformanceMetrics.AverageProfitTrade,
		AverageLossTrade:   result.PerformanceMetrics.AverageLossTrade,
		MaxConsecutiveWins: result.PerformanceMetrics.MaxConsecutiveWins,
		MaxConsecutiveLoss: result.PerformanceMetrics.MaxConsecutiveLoss,
		EquityCurve:        result.EquityCurve,
		DrawdownCurve:      result.DrawdownCurve,
		Trades:             result.Trades,
		CreatedAt:          time.Now(),
	}

	return apiResult, nil
}

// GetBacktestResult gets a backtest result by ID
func (s *BacktestService) GetBacktestResult(ctx context.Context, id string) (*BacktestResult, error) {
	// Remove prefix for internal service
	internalID := id
	if len(id) > 3 && id[:3] == "bt-" {
		internalID = id[3:]
	}

	// Get backtest result
	result, err := s.backtestService.GetBacktestResult(ctx, internalID)
	if err != nil {
		return nil, err
	}

	// Convert backtest result to API result
	apiResult := &BacktestResult{
		ID:               id, // Keep original ID
		Strategy:         result.Config.Strategy.Name(),
		Symbol:           result.Config.Symbols[0], // Assuming single symbol
		Timeframe:        result.Config.Interval,
		StartDate:        result.StartTime,
		EndDate:          result.EndTime,
		InitialCapital:   result.InitialCapital,
		FinalCapital:     result.FinalCapital,
		TotalReturn:      result.PerformanceMetrics.TotalReturn,
		AnnualizedReturn: result.PerformanceMetrics.AnnualizedReturn,
		MaxDrawdown:      result.PerformanceMetrics.MaxDrawdown,
		SharpeRatio:      result.PerformanceMetrics.SharpeRatio,
		WinRate:          result.PerformanceMetrics.WinRate,
		ProfitFactor:     result.PerformanceMetrics.ProfitFactor,
		TotalTrades:      result.PerformanceMetrics.TotalTrades,
		WinningTrades:    result.PerformanceMetrics.WinningTrades,
		LosingTrades:     result.PerformanceMetrics.LosingTrades,
		AverageProfitTrade: result.PerformanceMetrics.AverageProfitTrade,
		AverageLossTrade:   result.PerformanceMetrics.AverageLossTrade,
		MaxConsecutiveWins: result.PerformanceMetrics.MaxConsecutiveWins,
		MaxConsecutiveLoss: result.PerformanceMetrics.MaxConsecutiveLoss,
		EquityCurve:        result.EquityCurve,
		DrawdownCurve:      result.DrawdownCurve,
		Trades:             result.Trades,
		CreatedAt:          time.Now(),
	}

	return apiResult, nil
}

// ListBacktestResults lists all backtest results
func (s *BacktestService) ListBacktestResults(ctx context.Context) ([]*BacktestResult, error) {
	// Get backtest results
	results, err := s.backtestService.ListBacktestResults(ctx)
	if err != nil {
		return nil, err
	}

	// Convert backtest results to API results
	apiResults := make([]*BacktestResult, 0, len(results))
	for _, result := range results {
		apiResult := &BacktestResult{
			ID:               "bt-" + result.ID, // Add prefix for API
			Strategy:         result.Config.Strategy.Name(),
			Symbol:           result.Config.Symbols[0], // Assuming single symbol
			Timeframe:        result.Config.Interval,
			StartDate:        result.StartTime,
			EndDate:          result.EndTime,
			InitialCapital:   result.InitialCapital,
			FinalCapital:     result.FinalCapital,
			TotalReturn:      result.PerformanceMetrics.TotalReturn,
			AnnualizedReturn: result.PerformanceMetrics.AnnualizedReturn,
			MaxDrawdown:      result.PerformanceMetrics.MaxDrawdown,
			SharpeRatio:      result.PerformanceMetrics.SharpeRatio,
			WinRate:          result.PerformanceMetrics.WinRate,
			ProfitFactor:     result.PerformanceMetrics.ProfitFactor,
			TotalTrades:      result.PerformanceMetrics.TotalTrades,
			WinningTrades:    result.PerformanceMetrics.WinningTrades,
			LosingTrades:     result.PerformanceMetrics.LosingTrades,
			AverageProfitTrade: result.PerformanceMetrics.AverageProfitTrade,
			AverageLossTrade:   result.PerformanceMetrics.AverageLossTrade,
			MaxConsecutiveWins: result.PerformanceMetrics.MaxConsecutiveWins,
			MaxConsecutiveLoss: result.PerformanceMetrics.MaxConsecutiveLoss,
			EquityCurve:        result.EquityCurve,
			DrawdownCurve:      result.DrawdownCurve,
			Trades:             result.Trades,
			CreatedAt:          time.Now(),
		}
		apiResults = append(apiResults, apiResult)
	}

	return apiResults, nil
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

	// Compare backtests
	comparison := &BacktestComparisonResult{
		Backtests: results,
		Comparison: &BacktestComparison{
			BestTotalReturn:       results[0].ID,
			BestSharpeRatio:       results[0].ID,
			BestDrawdown:          results[0].ID,
			BestWinRate:           results[0].ID,
			BestProfitFactor:      results[0].ID,
			ReturnDifference:      0.0,
			DrawdownDifference:    0.0,
			SharpeRatioDifference: 0.0,
		},
		Timestamp: time.Now(),
	}

	// Find best metrics
	for _, result := range results {
		if result.TotalReturn > results[0].TotalReturn {
			comparison.Comparison.BestTotalReturn = result.ID
		}
		if result.SharpeRatio > results[0].SharpeRatio {
			comparison.Comparison.BestSharpeRatio = result.ID
		}
		if result.MaxDrawdown < results[0].MaxDrawdown {
			comparison.Comparison.BestDrawdown = result.ID
		}
		if result.WinRate > results[0].WinRate {
			comparison.Comparison.BestWinRate = result.ID
		}
		if result.ProfitFactor > results[0].ProfitFactor {
			comparison.Comparison.BestProfitFactor = result.ID
		}
	}

	// Calculate differences
	maxReturn := results[0].TotalReturn
	minReturn := results[0].TotalReturn
	maxDrawdown := results[0].MaxDrawdown
	minDrawdown := results[0].MaxDrawdown
	maxSharpeRatio := results[0].SharpeRatio
	minSharpeRatio := results[0].SharpeRatio

	for _, result := range results {
		if result.TotalReturn > maxReturn {
			maxReturn = result.TotalReturn
		}
		if result.TotalReturn < minReturn {
			minReturn = result.TotalReturn
		}
		if result.MaxDrawdown > maxDrawdown {
			maxDrawdown = result.MaxDrawdown
		}
		if result.MaxDrawdown < minDrawdown {
			minDrawdown = result.MaxDrawdown
		}
		if result.SharpeRatio > maxSharpeRatio {
			maxSharpeRatio = result.SharpeRatio
		}
		if result.SharpeRatio < minSharpeRatio {
			minSharpeRatio = result.SharpeRatio
		}
	}

	comparison.Comparison.ReturnDifference = maxReturn - minReturn
	comparison.Comparison.DrawdownDifference = maxDrawdown - minDrawdown
	comparison.Comparison.SharpeRatioDifference = maxSharpeRatio - minSharpeRatio

	return comparison, nil
}

// BacktestComparisonResult represents the result of comparing multiple backtests
type BacktestComparisonResult struct {
	Backtests  []*BacktestResult   `json:"backtests"`
	Comparison *BacktestComparison `json:"comparison"`
	Timestamp  time.Time           `json:"timestamp"`
}

// BacktestComparison represents the comparison metrics between backtests
type BacktestComparison struct {
	BestTotalReturn       string  `json:"bestTotalReturn"`
	BestSharpeRatio       string  `json:"bestSharpeRatio"`
	BestDrawdown          string  `json:"bestDrawdown"`
	BestWinRate           string  `json:"bestWinRate"`
	BestProfitFactor      string  `json:"bestProfitFactor"`
	ReturnDifference      float64 `json:"returnDifference"`
	DrawdownDifference    float64 `json:"drawdownDifference"`
	SharpeRatioDifference float64 `json:"sharpeRatioDifference"`
}
