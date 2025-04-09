package risk

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RiskManager implements the RiskService interface
type RiskManager struct {
	riskParams        RiskParameters
	balanceRepo       BalanceHistoryRepository
	positionSizer     PositionSizer
	drawdownMonitor   DrawdownMonitor
	exposureMonitor   ExposureMonitor
	dailyLimitMonitor DailyLimitMonitor
	lock              sync.RWMutex
	logger            Logger
}

// PositionSizer defines the interface for position sizing
type PositionSizer interface {
	CalculatePositionSize(ctx context.Context, symbol string, accountBalance float64, riskPercent float64, stopLossPercent float64) (float64, error)
}

// DrawdownMonitor defines the interface for drawdown monitoring
type DrawdownMonitor interface {
	CalculateDrawdown(ctx context.Context, days int) (float64, error)
	CheckDrawdownLimit(ctx context.Context, maxDrawdownPercent float64) (bool, error)
}

// ExposureMonitor defines the interface for exposure monitoring
type ExposureMonitor interface {
	CalculateTotalExposure(ctx context.Context) (float64, error)
	CheckExposureLimit(ctx context.Context, newOrderValue float64, maxExposurePercent float64) (bool, error)
	GetAccountBalance(ctx context.Context) (float64, error)
}

// DailyLimitMonitor defines the interface for daily limit monitoring
type DailyLimitMonitor interface {
	CalculateDailyPnL(ctx context.Context) (float64, error)
	CheckDailyLossLimit(ctx context.Context, dailyLossLimitPercent float64) (bool, error)
}

// Logger defines the interface for logging
type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// NewRiskManager creates a new RiskManager
func NewRiskManager(
	balanceRepo BalanceHistoryRepository,
	positionSizer PositionSizer,
	drawdownMonitor DrawdownMonitor,
	exposureMonitor ExposureMonitor,
	dailyLimitMonitor DailyLimitMonitor,
	logger Logger,
) *RiskManager {
	return &RiskManager{
		riskParams:        DefaultRiskParameters(),
		balanceRepo:       balanceRepo,
		positionSizer:     positionSizer,
		drawdownMonitor:   drawdownMonitor,
		exposureMonitor:   exposureMonitor,
		dailyLimitMonitor: dailyLimitMonitor,
		lock:              sync.RWMutex{},
		logger:            logger,
	}
}

// CalculatePositionSize determines a safe position size based on risk parameters
func (rm *RiskManager) CalculatePositionSize(
	ctx context.Context,
	symbol string,
	accountBalance float64,
) (float64, error) {
	rm.lock.RLock()
	riskPercent := rm.riskParams.RiskPerTradePercent
	rm.lock.RUnlock()

	// Default stop-loss at 5% below entry
	stopLossPercent := 5.0

	return rm.positionSizer.CalculatePositionSize(
		ctx,
		symbol,
		accountBalance,
		riskPercent,
		stopLossPercent,
	)
}

// CalculateDrawdown computes the maximum peak-to-trough drawdown
func (rm *RiskManager) CalculateDrawdown(ctx context.Context) (float64, error) {
	// Calculate drawdown over the last 90 days
	return rm.drawdownMonitor.CalculateDrawdown(ctx, 90)
}

// CheckExposureLimit verifies if a new order would exceed exposure limits
func (rm *RiskManager) CheckExposureLimit(ctx context.Context, newOrderValue float64) (bool, error) {
	rm.lock.RLock()
	maxExposurePercent := rm.riskParams.MaxExposurePercent
	rm.lock.RUnlock()

	return rm.exposureMonitor.CheckExposureLimit(ctx, newOrderValue, maxExposurePercent)
}

// CheckDailyLossLimit verifies if trading should be allowed based on daily P&L
func (rm *RiskManager) CheckDailyLossLimit(ctx context.Context) (bool, error) {
	rm.lock.RLock()
	dailyLossLimitPercent := rm.riskParams.DailyLossLimitPercent
	rm.lock.RUnlock()

	return rm.dailyLimitMonitor.CheckDailyLossLimit(ctx, dailyLossLimitPercent)
}

// GetRiskStatus returns the current risk metrics
func (rm *RiskManager) GetRiskStatus(ctx context.Context) (*RiskStatus, error) {
	// Get account balance from exposure monitor
	accountBalance, err := rm.exposureMonitor.GetAccountBalance(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get account balance: %w", err)
	}

	// Calculate drawdown
	drawdown, err := rm.CalculateDrawdown(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate drawdown: %w", err)
	}

	// Calculate total exposure
	totalExposure, err := rm.exposureMonitor.CalculateTotalExposure(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate total exposure: %w", err)
	}

	// Calculate today's P&L
	todayPnL, err := rm.dailyLimitMonitor.CalculateDailyPnL(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate today's P&L: %w", err)
	}

	// Check if trading is allowed
	allowed, reason, err := rm.IsTradeAllowed(ctx, "", 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check if trading is allowed: %w", err)
	}

	return &RiskStatus{
		CurrentDrawdown: drawdown,
		TotalExposure:   totalExposure,
		TodayPnL:        todayPnL,
		AccountBalance:  accountBalance,
		TradingEnabled:  allowed,
		DisabledReason:  reason,
		UpdatedAt:       time.Now(),
	}, nil
}

// UpdateRiskParameters updates the risk parameters
func (rm *RiskManager) UpdateRiskParameters(ctx context.Context, params RiskParameters) error {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	// Validate parameters
	if params.MaxDrawdownPercent <= 0 || params.MaxDrawdownPercent > 100 {
		return fmt.Errorf("invalid max drawdown percent: %f", params.MaxDrawdownPercent)
	}
	if params.RiskPerTradePercent <= 0 || params.RiskPerTradePercent > 100 {
		return fmt.Errorf("invalid risk per trade percent: %f", params.RiskPerTradePercent)
	}
	if params.MaxExposurePercent <= 0 || params.MaxExposurePercent > 100 {
		return fmt.Errorf("invalid max exposure percent: %f", params.MaxExposurePercent)
	}
	if params.DailyLossLimitPercent <= 0 || params.DailyLossLimitPercent > 100 {
		return fmt.Errorf("invalid daily loss limit percent: %f", params.DailyLossLimitPercent)
	}
	if params.MinAccountBalance < 0 {
		return fmt.Errorf("invalid min account balance: %f", params.MinAccountBalance)
	}

	// Update parameters
	rm.riskParams = params

	rm.logger.Info("Updated risk parameters",
		"max_drawdown_percent", params.MaxDrawdownPercent,
		"risk_per_trade_percent", params.RiskPerTradePercent,
		"max_exposure_percent", params.MaxExposurePercent,
		"daily_loss_limit_percent", params.DailyLossLimitPercent,
		"min_account_balance", params.MinAccountBalance)

	return nil
}

// IsTradeAllowed performs a comprehensive check of all risk controls
func (rm *RiskManager) IsTradeAllowed(
	ctx context.Context,
	symbol string,
	orderValue float64,
) (bool, string, error) {
	// Check account minimum balance
	accountBalance, err := rm.exposureMonitor.GetAccountBalance(ctx)
	if err != nil {
		return false, "", fmt.Errorf("failed to get account balance: %w", err)
	}

	rm.lock.RLock()
	minAccountBalance := rm.riskParams.MinAccountBalance
	maxDrawdownPercent := rm.riskParams.MaxDrawdownPercent
	rm.lock.RUnlock()

	if accountBalance < minAccountBalance {
		reason := fmt.Sprintf("Account balance below minimum: %.2f < %.2f", accountBalance, minAccountBalance)
		return false, reason, nil
	}

	// Check drawdown limit
	drawdownAllowed, err := rm.drawdownMonitor.CheckDrawdownLimit(ctx, maxDrawdownPercent)
	if err != nil {
		return false, "", err
	}

	if !drawdownAllowed {
		reason := "Maximum drawdown limit reached"
		return false, reason, nil
	}

	// Only check exposure if we're actually placing an order
	if orderValue > 0 {
		// Check exposure limit
		exposureAllowed, err := rm.CheckExposureLimit(ctx, orderValue)
		if err != nil {
			return false, "", err
		}

		if !exposureAllowed {
			reason := "Maximum exposure limit would be exceeded"
			return false, reason, nil
		}
	}

	// Check daily loss limit
	dailyLossAllowed, err := rm.CheckDailyLossLimit(ctx)
	if err != nil {
		return false, "", err
	}

	if !dailyLossAllowed {
		reason := "Daily loss limit reached"
		return false, reason, nil
	}

	return true, "", nil
}
