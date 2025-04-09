package backtest

import (
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateMonthlyReturns(t *testing.T) {
	// Create a test backtest result with equity curve spanning multiple months
	result := &BacktestResult{
		InitialCapital: 10000,
		EquityCurve: []*EquityPoint{
			{Timestamp: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Equity: 10000},
			{Timestamp: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), Equity: 10500},
			{Timestamp: time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC), Equity: 11000},
			{Timestamp: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC), Equity: 10800},
			{Timestamp: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC), Equity: 11200},
			{Timestamp: time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC), Equity: 11500},
			{Timestamp: time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC), Equity: 11800},
		},
	}

	// Create a performance analyzer
	analyzer := NewPerformanceAnalyzer()

	// Calculate monthly returns
	monthlyReturns, err := analyzer.CalculateMonthlyReturns(result)
	require.NoError(t, err)
	require.NotNil(t, monthlyReturns)

	// Check that we have returns for each month
	assert.Equal(t, 3, len(monthlyReturns))

	// Check the values (approximate due to floating point)
	assert.InDelta(t, 10.0, monthlyReturns["2023-01"], 0.1) // 10% return in January
	assert.InDelta(t, 1.82, monthlyReturns["2023-02"], 0.1) // ~1.82% return in February
	assert.InDelta(t, 5.36, monthlyReturns["2023-03"], 0.1) // ~5.36% return in March
}

func TestRunMonteCarloSimulation(t *testing.T) {
	// Create a test backtest result with equity curve
	result := &BacktestResult{
		InitialCapital: 10000,
		EquityCurve: []*EquityPoint{
			{Timestamp: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Equity: 10000},
			{Timestamp: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Equity: 10100},
			{Timestamp: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), Equity: 10200},
			{Timestamp: time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC), Equity: 10150},
			{Timestamp: time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC), Equity: 10250},
		},
	}

	// Create a performance analyzer
	analyzer := NewPerformanceAnalyzer()

	// Run Monte Carlo simulation with 10 simulations
	simulations, err := analyzer.RunMonteCarloSimulation(result, 10)
	require.NoError(t, err)
	require.NotNil(t, simulations)

	// Check that we have 10 simulations
	assert.Equal(t, 10, len(simulations))

	// Check that each simulation has the correct length
	for _, sim := range simulations {
		assert.Equal(t, 5, len(sim)) // 4 days of returns + initial capital
		assert.Equal(t, 10000.0, sim[0]) // Initial capital
	}
}

func TestCalculateTradeStats(t *testing.T) {
	// Create a test backtest result with closed positions
	result := &BacktestResult{
		InitialCapital: 10000,
		ClosedPositions: []*models.ClosedPosition{
			{
				OpenTime:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				CloseTime: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				Profit:    100,
			},
			{
				OpenTime:  time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
				CloseTime: time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC),
				Profit:    -50,
			},
			{
				OpenTime:  time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC),
				CloseTime: time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC),
				Profit:    200,
			},
			{
				OpenTime:  time.Date(2023, 1, 7, 0, 0, 0, 0, time.UTC),
				CloseTime: time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC),
				Profit:    -75,
			},
		},
		EquityCurve: []*EquityPoint{
			{Timestamp: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Equity: 10000},
			{Timestamp: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Equity: 10100},
			{Timestamp: time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC), Equity: 10050},
			{Timestamp: time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC), Equity: 10250},
			{Timestamp: time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC), Equity: 10175},
		},
	}

	// Create metrics
	metrics := &PerformanceMetrics{
		TotalTrades:        4,
		WinningTrades:      2,
		LosingTrades:       2,
		AverageProfitTrade: 150,
		AverageLossTrade:   -62.5,
		LargestProfitTrade: 200,
		LargestLossTrade:   -75,
	}

	// Create a performance analyzer
	analyzer := NewPerformanceAnalyzer()

	// Calculate trade statistics
	stats := analyzer.calculateTradeStats(result, metrics)
	require.NotNil(t, stats)

	// Check basic statistics
	assert.Equal(t, 150.0, stats.AverageWin)
	assert.Equal(t, -62.5, stats.AverageLoss)
	assert.Equal(t, 200.0, stats.LargestWin)
	assert.Equal(t, -75.0, stats.LargestLoss)

	// Check consecutive wins/losses
	assert.Equal(t, 1, stats.ConsecutiveWins)
	assert.Equal(t, 1, stats.ConsecutiveLosses)

	// Check holding times
	assert.Equal(t, 24*time.Hour, stats.MedianHoldingTime)
}

func TestAnalyzeRegimes(t *testing.T) {
	// Create a test backtest result with equity curve
	result := &BacktestResult{
		InitialCapital: 10000,
		EquityCurve: []*EquityPoint{
			{Timestamp: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Equity: 10000},
			{Timestamp: time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC), Equity: 10500},
			{Timestamp: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC), Equity: 10200},
			{Timestamp: time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC), Equity: 10800},
		},
	}

	// Create benchmark returns
	benchmarkReturns := map[string]float64{
		"2023-01": 2.0,  // Bull market
		"2023-02": -2.0, // Bear market
		"2023-03": 0.5,  // Sideways market
	}

	// Create a performance analyzer
	analyzer := NewPerformanceAnalyzer()

	// Analyze regimes
	regimes, err := analyzer.AnalyzeRegimes(result, benchmarkReturns)
	require.NoError(t, err)
	require.NotNil(t, regimes)

	// Check that we have regimes for each market type
	assert.Contains(t, regimes, "bull")
	assert.Contains(t, regimes, "bear")
	assert.Contains(t, regimes, "sideways")

	// Check that each regime has the correct months
	assert.Contains(t, regimes["bull"], "2023-01")
	assert.Contains(t, regimes["bear"], "2023-02")
	assert.Contains(t, regimes["sideways"], "2023-03")
}

func TestCalculateCorrelation(t *testing.T) {
	// Create a test backtest result with equity curve
	result := &BacktestResult{
		InitialCapital: 10000,
		EquityCurve: []*EquityPoint{
			{Timestamp: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Equity: 10000},
			{Timestamp: time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC), Equity: 10500},
			{Timestamp: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC), Equity: 10200},
			{Timestamp: time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC), Equity: 10800},
		},
	}

	// Create benchmark returns
	benchmarkReturns := map[string]float64{
		"2023-01": 4.0,  // Positive correlation
		"2023-02": -3.0, // Positive correlation
		"2023-03": 5.0,  // Positive correlation
	}

	// Create a performance analyzer
	analyzer := NewPerformanceAnalyzer()

	// Calculate correlation
	correlation, err := analyzer.CalculateCorrelation(result, benchmarkReturns)
	require.NoError(t, err)

	// Check that correlation is positive (should be close to 1.0)
	assert.Greater(t, correlation, 0.5)
}
