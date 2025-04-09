// Package backtest provides public interfaces for backtesting services
package backtest

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/backtest"
	"go-crypto-bot-clean/backend/internal/domain/models"
)

// Service is a public interface for the backtesting service
type Service interface {
	// RunBacktest runs a backtest with the given parameters
	RunBacktest(ctx context.Context, strategyName string, symbol string, startTime, endTime time.Time, initialBalance float64, params map[string]interface{}) (*models.BacktestResult, error)

	// GetBacktestResults gets the results of a backtest
	GetBacktestResults(ctx context.Context, backtestID string) (*models.BacktestResult, error)
}

// serviceAdapter adapts the internal backtest service to the public interface
type serviceAdapter struct {
	internalService *backtest.Service
}

// RunBacktest adapts the internal RunBacktest method to the public interface
func (a *serviceAdapter) RunBacktest(ctx context.Context, strategyName string, symbol string, startTime, endTime time.Time, initialBalance float64, params map[string]interface{}) (*models.BacktestResult, error) {
	// Convert parameters to internal request config
	reqConfig := &backtest.BacktestRequestConfig{
		Strategy:       strategyName,
		Symbol:         symbol,
		Timeframe:      "1h", // Default timeframe
		StartTime:      startTime,
		EndTime:        endTime,
		InitialCapital: initialBalance,
	}

	// Add any additional parameters
	if riskPerTrade, ok := params["riskPerTrade"].(float64); ok {
		reqConfig.RiskPerTrade = riskPerTrade
	}

	// Run the backtest
	internalResult, err := a.internalService.RunBacktest(ctx, reqConfig)
	if err != nil {
		return nil, err
	}

	// Convert internal result to public result
	return convertToPublicBacktestResult(internalResult)
}

// GetBacktestResults adapts the internal GetBacktestResult method to the public interface
func (a *serviceAdapter) GetBacktestResults(ctx context.Context, backtestID string) (*models.BacktestResult, error) {
	// Get the internal result
	internalResult, err := a.internalService.GetBacktestResult(ctx, backtestID)
	if err != nil {
		return nil, err
	}

	// Convert internal result to public result
	return convertToPublicBacktestResult(internalResult)
}

// convertToPublicBacktestResult converts an internal backtest result to a public one
func convertToPublicBacktestResult(internalResult *backtest.BacktestResult) (*models.BacktestResult, error) {
	// Create a simplified public result with the essential fields
	result := &models.BacktestResult{
		ID:             "bt-" + time.Now().Format("20060102150405"), // Generate a simple ID
		Strategy:       "unknown",                                   // Will be filled if available
		Symbol:         "unknown",                                   // Will be filled if available
		Timeframe:      "1h",                                        // Default
		StartTime:      internalResult.StartTime,
		EndTime:        internalResult.EndTime,
		InitialCapital: internalResult.InitialCapital,
		FinalCapital:   internalResult.FinalCapital,
		TotalTrades:    len(internalResult.Trades),
		WinningTrades:  internalResult.PerformanceMetrics.WinningTrades,
		LosingTrades:   internalResult.PerformanceMetrics.LosingTrades,
		WinRate:        internalResult.PerformanceMetrics.WinRate,
		ProfitFactor:   internalResult.PerformanceMetrics.ProfitFactor,
		MaxDrawdown:    internalResult.PerformanceMetrics.MaxDrawdown,
		SharpeRatio:    internalResult.PerformanceMetrics.SharpeRatio,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Fill in strategy and symbol if available
	if internalResult.Config != nil {
		if len(internalResult.Config.Symbols) > 0 {
			result.Symbol = internalResult.Config.Symbols[0]
		}
		result.Timeframe = internalResult.Config.Interval
	}

	return result, nil
}

// NewService creates a new backtesting service
func NewService() Service {
	return &serviceAdapter{
		internalService: backtest.NewService(),
	}
}
