package reporting

import (
	"context"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

// MockTradeAnalyticsRepository is a mock implementation of the trade analytics repository
type MockTradeAnalyticsRepository struct {
	mock.Mock
}

func (m *MockTradeAnalyticsRepository) GetWinRate(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTradeAnalyticsRepository) GetProfitFactor(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTradeAnalyticsRepository) GetSharpeRatio(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTradeAnalyticsRepository) GetDrawdown(ctx context.Context, startTime, endTime time.Time) (float64, float64, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

func (m *MockTradeAnalyticsRepository) GetPerformanceBySymbol(ctx context.Context, startTime, endTime time.Time) (map[string]models.SymbolPerformance, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).(map[string]models.SymbolPerformance), args.Error(1)
}

func (m *MockTradeAnalyticsRepository) GetPerformanceByReason(ctx context.Context, startTime, endTime time.Time) (map[string]models.ReasonPerformance, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).(map[string]models.ReasonPerformance), args.Error(1)
}

func (m *MockTradeAnalyticsRepository) GetPerformanceByStrategy(ctx context.Context, startTime, endTime time.Time) (map[string]models.StrategyPerformance, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).(map[string]models.StrategyPerformance), args.Error(1)
}

// MockBalanceHistoryRepository is a mock implementation of the balance history repository
type MockBalanceHistoryRepository struct {
	mock.Mock
}

func (m *MockBalanceHistoryRepository) GetLatestBalance(ctx context.Context) (*repository.BalanceHistory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.BalanceHistory), args.Error(1)
}

func (m *MockBalanceHistoryRepository) GetBalanceHistory(ctx context.Context, startTime, endTime time.Time) ([]*repository.BalanceHistory, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).([]*repository.BalanceHistory), args.Error(1)
}

func TestNewMetricsCollector(t *testing.T) {
	// Create dependencies
	tradeAnalyticsRepo := &MockTradeAnalyticsRepository{}
	balanceHistoryRepo := &MockBalanceHistoryRepository{}
	logger := zaptest.NewLogger(t)

	// Create collector
	collector := NewMetricsCollector(tradeAnalyticsRepo, balanceHistoryRepo, logger)

	// Assert
	assert.NotNil(t, collector)
	assert.Equal(t, tradeAnalyticsRepo, collector.tradeAnalyticsRepo)
	assert.Equal(t, balanceHistoryRepo, collector.balanceHistoryRepo)
	assert.Equal(t, logger, collector.logger)
}

func TestMetricsCollector_CollectMetrics(t *testing.T) {
	// Create dependencies
	tradeAnalyticsRepo := &MockTradeAnalyticsRepository{}
	balanceHistoryRepo := &MockBalanceHistoryRepository{}
	logger := zaptest.NewLogger(t)

	// Create collector
	collector := NewMetricsCollector(tradeAnalyticsRepo, balanceHistoryRepo, logger)

	// Create test data
	ctx := context.Background()
	// Use fixed times to avoid issues with time precision in tests
	now := time.Date(2023, 4, 7, 12, 0, 0, 0, time.UTC)
	dayStart := time.Date(2023, 4, 6, 12, 0, 0, 0, time.UTC)
	weekStart := time.Date(2023, 3, 31, 12, 0, 0, 0, time.UTC)
	monthStart := time.Date(2023, 3, 8, 12, 0, 0, 0, time.UTC)

	// Set up expectations
	tradeAnalyticsRepo.On("GetWinRate", ctx, dayStart, now).Return(0.75, nil)
	tradeAnalyticsRepo.On("GetProfitFactor", ctx, dayStart, now).Return(2.5, nil)
	tradeAnalyticsRepo.On("GetDrawdown", ctx, dayStart, now).Return(100.0, 5.0, nil)

	tradeAnalyticsRepo.On("GetWinRate", ctx, weekStart, now).Return(0.65, nil)
	tradeAnalyticsRepo.On("GetProfitFactor", ctx, weekStart, now).Return(2.0, nil)
	tradeAnalyticsRepo.On("GetSharpeRatio", ctx, weekStart, now).Return(1.5, nil)
	tradeAnalyticsRepo.On("GetDrawdown", ctx, weekStart, now).Return(200.0, 10.0, nil)

	tradeAnalyticsRepo.On("GetWinRate", ctx, monthStart, now).Return(0.60, nil)
	tradeAnalyticsRepo.On("GetProfitFactor", ctx, monthStart, now).Return(1.8, nil)
	tradeAnalyticsRepo.On("GetSharpeRatio", ctx, monthStart, now).Return(1.2, nil)
	tradeAnalyticsRepo.On("GetDrawdown", ctx, monthStart, now).Return(300.0, 15.0, nil)

	symbolPerformance := map[string]models.SymbolPerformance{
		"BTC": {
			Symbol:        "BTC",
			TotalTrades:   10,
			WinningTrades: 7,
			LosingTrades:  3,
			WinRate:       0.7,
			TotalProfit:   1000.0,
			AverageProfit: 100.0,
			ProfitFactor:  2.0,
		},
	}
	tradeAnalyticsRepo.On("GetPerformanceBySymbol", ctx, monthStart, now).Return(symbolPerformance, nil)

	reasonPerformance := map[string]models.ReasonPerformance{
		"breakout": {
			Reason:        "breakout",
			TotalTrades:   5,
			WinningTrades: 4,
			LosingTrades:  1,
			WinRate:       0.8,
			TotalProfit:   500.0,
			AverageProfit: 100.0,
			ProfitFactor:  4.0,
		},
	}
	tradeAnalyticsRepo.On("GetPerformanceByReason", ctx, monthStart, now).Return(reasonPerformance, nil)

	strategyPerformance := map[string]models.StrategyPerformance{
		"trend_following": {
			Strategy:      "trend_following",
			TotalTrades:   8,
			WinningTrades: 5,
			LosingTrades:  3,
			WinRate:       0.625,
			TotalProfit:   800.0,
			AverageProfit: 100.0,
			ProfitFactor:  2.0,
		},
	}
	tradeAnalyticsRepo.On("GetPerformanceByStrategy", ctx, monthStart, now).Return(strategyPerformance, nil)

	latestBalance := &repository.BalanceHistory{
		ID:            1,
		Timestamp:     now,
		Balance:       10000.0,
		Equity:        10500.0,
		FreeBalance:   9000.0,
		LockedBalance: 1000.0,
		UnrealizedPnL: 500.0,
	}
	balanceHistoryRepo.On("GetLatestBalance", ctx).Return(latestBalance, nil)

	balanceHistory := []*repository.BalanceHistory{
		{
			ID:        1,
			Timestamp: monthStart,
			Balance:   9000.0,
		},
		{
			ID:        2,
			Timestamp: now,
			Balance:   10000.0,
		},
	}
	balanceHistoryRepo.On("GetBalanceHistory", ctx, monthStart, now).Return(balanceHistory, nil)

	// Call the method with fixed time ranges
	metrics, err := collector.CollectMetrics(ctx, now, dayStart, weekStart, monthStart)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, metrics)

	// Check metrics
	assert.Equal(t, 0.75, metrics["daily_win_rate"])
	assert.Equal(t, 2.5, metrics["daily_profit_factor"])
	assert.Equal(t, 100.0, metrics["daily_drawdown"])
	assert.Equal(t, 5.0, metrics["daily_drawdown_percent"])

	assert.Equal(t, 0.65, metrics["weekly_win_rate"])
	assert.Equal(t, 2.0, metrics["weekly_profit_factor"])
	assert.Equal(t, 1.5, metrics["weekly_sharpe_ratio"])
	assert.Equal(t, 200.0, metrics["weekly_drawdown"])
	assert.Equal(t, 10.0, metrics["weekly_drawdown_percent"])

	assert.Equal(t, 0.60, metrics["monthly_win_rate"])
	assert.Equal(t, 1.8, metrics["monthly_profit_factor"])
	assert.Equal(t, 1.2, metrics["monthly_sharpe_ratio"])
	assert.Equal(t, 300.0, metrics["monthly_drawdown"])
	assert.Equal(t, 15.0, metrics["monthly_drawdown_percent"])

	assert.Equal(t, symbolPerformance, metrics["symbol_performance"])
	assert.Equal(t, reasonPerformance, metrics["reason_performance"])
	assert.Equal(t, strategyPerformance, metrics["strategy_performance"])

	assert.Equal(t, 10000.0, metrics["current_balance"])
	assert.Equal(t, 10500.0, metrics["current_equity"])
	assert.Equal(t, 9000.0, metrics["free_balance"])
	assert.Equal(t, 1000.0, metrics["locked_balance"])
	assert.Equal(t, 500.0, metrics["unrealized_pnl"])

	assert.Equal(t, 11.11111111111111, metrics["monthly_balance_growth_percent"])

	// Verify all expectations were met
	tradeAnalyticsRepo.AssertExpectations(t)
	balanceHistoryRepo.AssertExpectations(t)
}
