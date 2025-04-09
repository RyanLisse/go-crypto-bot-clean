package backtest

import (
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCalculateMetrics tests the CalculateMetrics function
func TestCalculateMetrics(t *testing.T) {
	// Create a test backtest result
	result := createTestBacktestResult()

	// Create a performance analyzer
	analyzer := NewPerformanceAnalyzer()

	// Calculate metrics
	metrics, err := analyzer.CalculateMetrics(result)
	require.NoError(t, err)
	require.NotNil(t, metrics)

	// Check basic metrics
	assert.InDelta(t, 20.0, metrics.TotalReturn, 0.01)
	assert.InDelta(t, 365.0, metrics.AnnualizedReturn, 0.01)
	assert.Equal(t, 3, metrics.TotalTrades)
	assert.Equal(t, 2, metrics.WinningTrades)
	assert.Equal(t, 1, metrics.LosingTrades)
	assert.InDelta(t, 66.67, metrics.WinRate, 0.01)
	assert.InDelta(t, 2.0, metrics.ProfitFactor, 0.01)
	assert.InDelta(t, 500.0, metrics.AverageProfitTrade, 0.01)
	assert.InDelta(t, 500.0, metrics.AverageLossTrade, 0.01)
	assert.InDelta(t, 166.67, metrics.ExpectedPayoff, 0.01)
	assert.InDelta(t, 500.0, metrics.MaxDrawdown, 0.01)
	assert.InDelta(t, 5.0, metrics.MaxDrawdownPercent, 0.01)
}

// TestGenerateReport tests the GenerateReport function
func TestGenerateReport(t *testing.T) {
	// Create a test backtest result
	result := createTestBacktestResult()

	// Create a performance analyzer
	analyzer := NewPerformanceAnalyzer()

	// Calculate metrics
	metrics, err := analyzer.CalculateMetrics(result)
	require.NoError(t, err)
	require.NotNil(t, metrics)

	// Generate report
	report, err := analyzer.GenerateReport(result, metrics)
	require.NoError(t, err)
	require.NotNil(t, report)

	// Check report
	assert.Equal(t, metrics, report.Metrics)
	assert.NotEmpty(t, report.MonthlyReturns)
	assert.NotNil(t, report.TradeStats)
	assert.Equal(t, result.EquityCurve, report.EquityCurve)
	assert.Equal(t, result.DrawdownCurve, report.DrawdownCurve)

	// Check advanced trade statistics
	assert.InDelta(t, metrics.CalmarRatio, report.TradeStats.CalmarRatio, 0.01)
	assert.InDelta(t, metrics.OmegaRatio, report.TradeStats.OmegaRatio, 0.01)
	assert.InDelta(t, metrics.InformationRatio, report.TradeStats.InformationRatio, 0.01)
}

// TestGenerateEquityCurve tests the GenerateEquityCurve function
func TestGenerateEquityCurve(t *testing.T) {
	// Create a test backtest result
	result := createTestBacktestResult()

	// Create a performance analyzer
	analyzer := NewPerformanceAnalyzer()

	// Generate equity curve
	equityCurve, err := analyzer.GenerateEquityCurve(result)
	require.NoError(t, err)
	require.NotNil(t, equityCurve)

	// Check equity curve
	assert.Equal(t, result.EquityCurve, equityCurve)
}

// TestGenerateDrawdownCurve tests the GenerateDrawdownCurve function
func TestGenerateDrawdownCurve(t *testing.T) {
	// Create a test backtest result
	result := createTestBacktestResult()

	// Create a performance analyzer
	analyzer := NewPerformanceAnalyzer()

	// Generate drawdown curve
	drawdownCurve, err := analyzer.GenerateDrawdownCurve(result)
	require.NoError(t, err)
	require.NotNil(t, drawdownCurve)

	// Check drawdown curve
	assert.Equal(t, result.DrawdownCurve, drawdownCurve)
}

// createTestBacktestResult creates a test backtest result
func createTestBacktestResult() *BacktestResult {
	// Create test data
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	initialCapital := 10000.0
	finalCapital := 12000.0

	// Create closed positions
	closedPositions := []*models.ClosedPosition{
		{
			ID:         "trade-1",
			Symbol:     "BTCUSDT",
			Side:       models.OrderSideBuy,
			Quantity:   1.0,
			EntryPrice: 10000.0,
			ExitPrice:  10500.0,
			ProfitLoss: 500.0,
			OpenTime:   startTime.Add(1 * time.Hour),
			CloseTime:  startTime.Add(2 * time.Hour),
		},
		{
			ID:         "trade-2",
			Symbol:     "BTCUSDT",
			Side:       models.OrderSideBuy,
			Quantity:   1.0,
			EntryPrice: 10500.0,
			ExitPrice:  10000.0,
			ProfitLoss: -500.0,
			OpenTime:   startTime.Add(3 * time.Hour),
			CloseTime:  startTime.Add(4 * time.Hour),
		},
		{
			ID:         "trade-3",
			Symbol:     "BTCUSDT",
			Side:       models.OrderSideBuy,
			Quantity:   1.0,
			EntryPrice: 10000.0,
			ExitPrice:  10500.0,
			ProfitLoss: 500.0,
			OpenTime:   startTime.Add(5 * time.Hour),
			CloseTime:  startTime.Add(6 * time.Hour),
		},
	}

	// Create equity curve
	equityCurve := []*EquityPoint{
		{Timestamp: startTime, Equity: initialCapital},
		{Timestamp: startTime.Add(1 * time.Hour), Equity: 10100.0},
		{Timestamp: startTime.Add(2 * time.Hour), Equity: 10500.0},
		{Timestamp: startTime.Add(3 * time.Hour), Equity: 10400.0},
		{Timestamp: startTime.Add(4 * time.Hour), Equity: 10000.0},
		{Timestamp: startTime.Add(5 * time.Hour), Equity: 10200.0},
		{Timestamp: startTime.Add(6 * time.Hour), Equity: 10500.0},
		{Timestamp: endTime, Equity: finalCapital},
	}

	// Create drawdown curve
	drawdownCurve := []*DrawdownPoint{
		{Timestamp: startTime, Drawdown: 0.0},
		{Timestamp: startTime.Add(1 * time.Hour), Drawdown: 0.0},
		{Timestamp: startTime.Add(2 * time.Hour), Drawdown: 0.0},
		{Timestamp: startTime.Add(3 * time.Hour), Drawdown: 100.0},
		{Timestamp: startTime.Add(4 * time.Hour), Drawdown: 500.0},
		{Timestamp: startTime.Add(5 * time.Hour), Drawdown: 300.0},
		{Timestamp: startTime.Add(6 * time.Hour), Drawdown: 0.0},
		{Timestamp: endTime, Drawdown: 0.0},
	}

	// Create backtest result
	result := &BacktestResult{
		StartTime:       startTime,
		EndTime:         endTime,
		InitialCapital:  initialCapital,
		FinalCapital:    finalCapital,
		ClosedPositions: closedPositions,
		EquityCurve:     equityCurve,
		DrawdownCurve:   drawdownCurve,
	}

	return result
}
