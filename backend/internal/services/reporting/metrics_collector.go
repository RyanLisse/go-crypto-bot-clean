package reporting

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/repository"
	"go.uber.org/zap"
)

// DefaultMetricsCollector is the default implementation of MetricsCollector
type DefaultMetricsCollector struct {
	tradeAnalyticsRepo TradeAnalyticsRepository
	balanceHistoryRepo BalanceHistoryRepository
	logger             *zap.Logger
}

// TradeAnalyticsRepository defines the interface for trade analytics
type TradeAnalyticsRepository interface {
	GetWinRate(ctx context.Context, startTime, endTime time.Time) (float64, error)
	GetProfitFactor(ctx context.Context, startTime, endTime time.Time) (float64, error)
	GetSharpeRatio(ctx context.Context, startTime, endTime time.Time) (float64, error)
	GetDrawdown(ctx context.Context, startTime, endTime time.Time) (float64, float64, error)
	GetPerformanceBySymbol(ctx context.Context, startTime, endTime time.Time) (map[string]models.SymbolPerformance, error)
	GetPerformanceByReason(ctx context.Context, startTime, endTime time.Time) (map[string]models.ReasonPerformance, error)
	GetPerformanceByStrategy(ctx context.Context, startTime, endTime time.Time) (map[string]models.StrategyPerformance, error)
}

// BalanceHistoryRepository defines the interface for balance history
type BalanceHistoryRepository interface {
	GetLatestBalance(ctx context.Context) (*repository.BalanceHistory, error)
	GetBalanceHistory(ctx context.Context, startTime, endTime time.Time) ([]*repository.BalanceHistory, error)
}

// NewMetricsCollector creates a new DefaultMetricsCollector
func NewMetricsCollector(
	tradeAnalyticsRepo TradeAnalyticsRepository,
	balanceHistoryRepo BalanceHistoryRepository,
	logger *zap.Logger,
) *DefaultMetricsCollector {
	return &DefaultMetricsCollector{
		tradeAnalyticsRepo: tradeAnalyticsRepo,
		balanceHistoryRepo: balanceHistoryRepo,
		logger:             logger,
	}
}

// CollectMetrics collects performance metrics
func (c *DefaultMetricsCollector) CollectMetrics(ctx context.Context, timeRanges ...time.Time) (map[string]interface{}, error) {
	// Define time ranges
	var now, dayStart, weekStart, monthStart time.Time

	// If time ranges are provided, use them
	if len(timeRanges) >= 4 {
		now = timeRanges[0]
		dayStart = timeRanges[1]
		weekStart = timeRanges[2]
		monthStart = timeRanges[3]
	} else {
		// Otherwise use default time ranges
		now = time.Now()
		dayStart = now.Add(-24 * time.Hour)
		weekStart = now.Add(-7 * 24 * time.Hour)
		monthStart = now.Add(-30 * 24 * time.Hour)
	}

	// Collect metrics
	metrics := make(map[string]interface{})

	// Get daily metrics
	dailyWinRate, err := c.tradeAnalyticsRepo.GetWinRate(ctx, dayStart, now)
	if err != nil {
		c.logger.Error("Failed to get daily win rate", zap.Error(err))
	} else {
		metrics["daily_win_rate"] = dailyWinRate
	}

	dailyProfitFactor, err := c.tradeAnalyticsRepo.GetProfitFactor(ctx, dayStart, now)
	if err != nil {
		c.logger.Error("Failed to get daily profit factor", zap.Error(err))
	} else {
		metrics["daily_profit_factor"] = dailyProfitFactor
	}

	dailyDrawdown, dailyDrawdownPercent, err := c.tradeAnalyticsRepo.GetDrawdown(ctx, dayStart, now)
	if err != nil {
		c.logger.Error("Failed to get daily drawdown", zap.Error(err))
	} else {
		metrics["daily_drawdown"] = dailyDrawdown
		metrics["daily_drawdown_percent"] = dailyDrawdownPercent
	}

	// Get weekly metrics
	weeklyWinRate, err := c.tradeAnalyticsRepo.GetWinRate(ctx, weekStart, now)
	if err != nil {
		c.logger.Error("Failed to get weekly win rate", zap.Error(err))
	} else {
		metrics["weekly_win_rate"] = weeklyWinRate
	}

	weeklyProfitFactor, err := c.tradeAnalyticsRepo.GetProfitFactor(ctx, weekStart, now)
	if err != nil {
		c.logger.Error("Failed to get weekly profit factor", zap.Error(err))
	} else {
		metrics["weekly_profit_factor"] = weeklyProfitFactor
	}

	weeklySharpeRatio, err := c.tradeAnalyticsRepo.GetSharpeRatio(ctx, weekStart, now)
	if err != nil {
		c.logger.Error("Failed to get weekly Sharpe ratio", zap.Error(err))
	} else {
		metrics["weekly_sharpe_ratio"] = weeklySharpeRatio
	}

	weeklyDrawdown, weeklyDrawdownPercent, err := c.tradeAnalyticsRepo.GetDrawdown(ctx, weekStart, now)
	if err != nil {
		c.logger.Error("Failed to get weekly drawdown", zap.Error(err))
	} else {
		metrics["weekly_drawdown"] = weeklyDrawdown
		metrics["weekly_drawdown_percent"] = weeklyDrawdownPercent
	}

	// Get monthly metrics
	monthlyWinRate, err := c.tradeAnalyticsRepo.GetWinRate(ctx, monthStart, now)
	if err != nil {
		c.logger.Error("Failed to get monthly win rate", zap.Error(err))
	} else {
		metrics["monthly_win_rate"] = monthlyWinRate
	}

	monthlyProfitFactor, err := c.tradeAnalyticsRepo.GetProfitFactor(ctx, monthStart, now)
	if err != nil {
		c.logger.Error("Failed to get monthly profit factor", zap.Error(err))
	} else {
		metrics["monthly_profit_factor"] = monthlyProfitFactor
	}

	monthlySharpeRatio, err := c.tradeAnalyticsRepo.GetSharpeRatio(ctx, monthStart, now)
	if err != nil {
		c.logger.Error("Failed to get monthly Sharpe ratio", zap.Error(err))
	} else {
		metrics["monthly_sharpe_ratio"] = monthlySharpeRatio
	}

	monthlyDrawdown, monthlyDrawdownPercent, err := c.tradeAnalyticsRepo.GetDrawdown(ctx, monthStart, now)
	if err != nil {
		c.logger.Error("Failed to get monthly drawdown", zap.Error(err))
	} else {
		metrics["monthly_drawdown"] = monthlyDrawdown
		metrics["monthly_drawdown_percent"] = monthlyDrawdownPercent
	}

	// Get performance by symbol
	symbolPerformance, err := c.tradeAnalyticsRepo.GetPerformanceBySymbol(ctx, monthStart, now)
	if err != nil {
		c.logger.Error("Failed to get symbol performance", zap.Error(err))
	} else {
		metrics["symbol_performance"] = symbolPerformance
	}

	// Get performance by reason
	reasonPerformance, err := c.tradeAnalyticsRepo.GetPerformanceByReason(ctx, monthStart, now)
	if err != nil {
		c.logger.Error("Failed to get reason performance", zap.Error(err))
	} else {
		metrics["reason_performance"] = reasonPerformance
	}

	// Get performance by strategy
	strategyPerformance, err := c.tradeAnalyticsRepo.GetPerformanceByStrategy(ctx, monthStart, now)
	if err != nil {
		c.logger.Error("Failed to get strategy performance", zap.Error(err))
	} else {
		metrics["strategy_performance"] = strategyPerformance
	}

	// Get latest balance
	latestBalance, err := c.balanceHistoryRepo.GetLatestBalance(ctx)
	if err != nil {
		c.logger.Error("Failed to get latest balance", zap.Error(err))
	} else {
		metrics["current_balance"] = latestBalance.Balance
		metrics["current_equity"] = latestBalance.Equity
		metrics["free_balance"] = latestBalance.FreeBalance
		metrics["locked_balance"] = latestBalance.LockedBalance
		metrics["unrealized_pnl"] = latestBalance.UnrealizedPnL
	}

	// Get balance history for the last month
	balanceHistory, err := c.balanceHistoryRepo.GetBalanceHistory(ctx, monthStart, now)
	if err != nil {
		c.logger.Error("Failed to get balance history", zap.Error(err))
	} else {
		// Calculate balance growth
		if len(balanceHistory) > 0 {
			firstBalance := balanceHistory[0].Balance
			lastBalance := balanceHistory[len(balanceHistory)-1].Balance

			if firstBalance > 0 {
				balanceGrowth := (lastBalance - firstBalance) / firstBalance * 100
				metrics["monthly_balance_growth_percent"] = balanceGrowth
			}
		}
	}

	return metrics, nil
}
