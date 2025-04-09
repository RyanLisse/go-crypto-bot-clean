package repository

import (
	"context"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// TradeAnalyticsRepository defines the interface for trade analytics persistence and retrieval
type TradeAnalyticsRepository interface {
	// GetAnalytics retrieves analytics for a specific time frame
	GetAnalytics(ctx context.Context, timeFrame models.TimeFrame, startTime, endTime time.Time) (*models.TradeAnalytics, error)
	
	// GetPerformanceBySymbol retrieves performance metrics grouped by symbol
	GetPerformanceBySymbol(ctx context.Context, startTime, endTime time.Time) (map[string]models.SymbolPerformance, error)
	
	// GetPerformanceByReason retrieves performance metrics grouped by decision reason
	GetPerformanceByReason(ctx context.Context, startTime, endTime time.Time) (map[string]models.ReasonPerformance, error)
	
	// GetPerformanceByStrategy retrieves performance metrics grouped by strategy
	GetPerformanceByStrategy(ctx context.Context, startTime, endTime time.Time) (map[string]models.StrategyPerformance, error)
	
	// GetTradePerformance retrieves performance metrics for a specific trade
	GetTradePerformance(ctx context.Context, tradeID string) (*models.TradePerformance, error)
	
	// GetAllTradePerformance retrieves performance metrics for all trades in a time range
	GetAllTradePerformance(ctx context.Context, startTime, endTime time.Time) ([]*models.TradePerformance, error)
	
	// GetBalanceHistory retrieves the balance history for a time range
	GetBalanceHistory(ctx context.Context, startTime, endTime time.Time, interval time.Duration) ([]models.BalancePoint, error)
	
	// GetEquityCurve retrieves the equity curve for a time range
	GetEquityCurve(ctx context.Context, startTime, endTime time.Time, interval time.Duration) ([]models.EquityPoint, error)
	
	// GetDrawdown calculates the maximum drawdown for a time range
	GetDrawdown(ctx context.Context, startTime, endTime time.Time) (float64, float64, error)
	
	// GetWinRate calculates the win rate for a time range
	GetWinRate(ctx context.Context, startTime, endTime time.Time) (float64, error)
	
	// GetProfitFactor calculates the profit factor for a time range
	GetProfitFactor(ctx context.Context, startTime, endTime time.Time) (float64, error)
	
	// GetSharpeRatio calculates the Sharpe ratio for a time range
	GetSharpeRatio(ctx context.Context, startTime, endTime time.Time) (float64, error)
}
