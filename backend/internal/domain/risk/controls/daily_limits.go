package controls

import (
	"context"
	"fmt"
	"time"
)

// DailyLimitMonitor tracks and enforces daily trading limits
type DailyLimitMonitor struct {
	tradeRepo    TradeRepository
	accountSvc   AccountService
	logger       Logger
}

// TradeRepository defines the interface for accessing trade data
type TradeRepository interface {
	GetTradesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]Trade, error)
}

// Trade represents a completed trade
type Trade struct {
	ID        uint
	Symbol    string
	Quantity  float64
	BuyPrice  float64
	SellPrice float64
	PnL       float64
	BuyTime   time.Time
	SellTime  time.Time
}

// NewDailyLimitMonitor creates a new DailyLimitMonitor
func NewDailyLimitMonitor(
	tradeRepo TradeRepository,
	accountSvc AccountService,
	logger Logger,
) *DailyLimitMonitor {
	return &DailyLimitMonitor{
		tradeRepo:  tradeRepo,
		accountSvc: accountSvc,
		logger:     logger,
	}
}

// CalculateDailyPnL calculates the profit/loss for the current day
func (dm *DailyLimitMonitor) CalculateDailyPnL(ctx context.Context) (float64, error) {
	// Get today's date range
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Get today's trades
	trades, err := dm.tradeRepo.GetTradesByDateRange(ctx, startOfDay, endOfDay)
	if err != nil {
		return 0, fmt.Errorf("failed to get today's trades: %w", err)
	}

	// Calculate today's P&L
	var todayPnL float64
	for _, trade := range trades {
		todayPnL += trade.PnL
	}

	dm.logger.Info("Calculated daily P&L",
		"date", startOfDay.Format("2006-01-02"),
		"trade_count", len(trades),
		"pnl", todayPnL)

	return todayPnL, nil
}

// CheckDailyLossLimit verifies if trading should be allowed based on daily P&L
func (dm *DailyLimitMonitor) CheckDailyLossLimit(
	ctx context.Context,
	dailyLossLimitPercent float64,
) (bool, error) {
	// Get today's P&L
	todayPnL, err := dm.CalculateDailyPnL(ctx)
	if err != nil {
		return false, err
	}

	// Get account balance
	accountBalance, err := dm.accountSvc.GetBalance(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get account balance: %w", err)
	}

	// Calculate maximum allowed daily loss
	maxDailyLoss := accountBalance * (dailyLossLimitPercent / 100)

	// Check if today's losses exceed the limit
	allowed := todayPnL >= -maxDailyLoss

	if !allowed {
		dm.logger.Warn("Trading disabled due to daily loss limit",
			"today_pnl", todayPnL,
			"max_daily_loss", maxDailyLoss)
	}

	return allowed, nil
}
