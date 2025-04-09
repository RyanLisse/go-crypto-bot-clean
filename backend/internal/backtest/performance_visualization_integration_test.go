package backtest

import (
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPerformanceVisualizationIntegration tests the integration of performance metrics and visualization components
func TestPerformanceVisualizationIntegration(t *testing.T) {
	// Create a test backtest result with equity curve, drawdown curve, and closed positions
	result := &BacktestResult{
		InitialCapital: 10000,
		FinalCapital:   11000,
		StartTime:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:        time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
		EquityCurve: []*EquityPoint{
			{Timestamp: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Equity: 10000},
			{Timestamp: time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC), Equity: 10200},
			{Timestamp: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), Equity: 10100},
			{Timestamp: time.Date(2023, 1, 22, 0, 0, 0, 0, time.UTC), Equity: 10500},
			{Timestamp: time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC), Equity: 11000},
		},
		DrawdownCurve: []*DrawdownPoint{
			{Timestamp: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Drawdown: 0},
			{Timestamp: time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC), Drawdown: 0},
			{Timestamp: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), Drawdown: 100},
			{Timestamp: time.Date(2023, 1, 22, 0, 0, 0, 0, time.UTC), Drawdown: 0},
			{Timestamp: time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC), Drawdown: 0},
		},
		ClosedPositions: []*models.ClosedPosition{
			{
				OpenTime:  time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				CloseTime: time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC),
				Profit:    200,
			},
			{
				OpenTime:  time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC),
				CloseTime: time.Date(2023, 1, 12, 0, 0, 0, 0, time.UTC),
				Profit:    -100,
			},
			{
				OpenTime:  time.Date(2023, 1, 18, 0, 0, 0, 0, time.UTC),
				CloseTime: time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC),
				Profit:    400,
			},
			{
				OpenTime:  time.Date(2023, 1, 25, 0, 0, 0, 0, time.UTC),
				CloseTime: time.Date(2023, 1, 28, 0, 0, 0, 0, time.UTC),
				Profit:    500,
			},
		},
	}

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

	// Test monthly returns
	monthlyReturns, err := analyzer.CalculateMonthlyReturns(result)
	require.NoError(t, err)
	require.NotNil(t, monthlyReturns)
	assert.Equal(t, 1, len(monthlyReturns)) // Only one month in the test data
	assert.InDelta(t, 10.0, monthlyReturns["2023-01"], 0.1) // 10% return in January

	// Test Monte Carlo simulation
	simulations, err := analyzer.RunMonteCarloSimulation(result, 10)
	require.NoError(t, err)
	require.NotNil(t, simulations)
	assert.Equal(t, 10, len(simulations)) // 10 simulations
	assert.Equal(t, 5, len(simulations[0])) // 5 data points per simulation

	// Test equity curve
	equityCurve, err := analyzer.GenerateEquityCurve(result)
	require.NoError(t, err)
	require.NotNil(t, equityCurve)
	assert.Equal(t, result.EquityCurve, equityCurve)

	// Test drawdown curve
	drawdownCurve, err := analyzer.GenerateDrawdownCurve(result)
	require.NoError(t, err)
	require.NotNil(t, drawdownCurve)
	assert.Equal(t, result.DrawdownCurve, drawdownCurve)

	// Test trade statistics
	tradeStats := report.TradeStats
	require.NotNil(t, tradeStats)
	assert.Equal(t, 3, tradeStats.WinningTrades)
	assert.Equal(t, 1, tradeStats.LosingTrades)
	assert.InDelta(t, 366.67, tradeStats.AverageWin, 0.1)
	assert.Equal(t, -100.0, tradeStats.AverageLoss)

	// Test performance metrics
	assert.InDelta(t, 10.0, metrics.TotalReturn, 0.1) // 10% return
	assert.Greater(t, metrics.AnnualizedReturn, 100.0) // Annualized return should be high for a 1-month period
	assert.InDelta(t, 75.0, metrics.WinRate, 0.1) // 3 out of 4 trades were profitable
	assert.InDelta(t, 11.0, metrics.ProfitFactor, 0.1) // (200+400+500)/100 = 11
}
