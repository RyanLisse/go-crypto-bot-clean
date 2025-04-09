package service

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/risk"
	"go-crypto-bot-clean/backend/internal/domain/risk/controls"
	"go.uber.org/zap"
)

// RiskManager implements the risk.RiskService interface
type RiskManager struct {
	riskParams        risk.RiskParameters
	balanceRepo       risk.BalanceHistoryRepository
	positionSizer     *controls.PositionSizer
	drawdownMonitor   *controls.DrawdownMonitor
	exposureMonitor   *controls.ExposureMonitor
	dailyLimitMonitor *controls.DailyLimitMonitor
	lock              sync.RWMutex
	logger            *zap.Logger
}

// NewRiskManager creates a new RiskManager
func NewRiskManager(
	balanceRepo risk.BalanceHistoryRepository,
	positionSizer *controls.PositionSizer,
	drawdownMonitor *controls.DrawdownMonitor,
	exposureMonitor *controls.ExposureMonitor,
	dailyLimitMonitor *controls.DailyLimitMonitor,
	logger *zap.Logger,
) *RiskManager {
	return &RiskManager{
		riskParams:        risk.DefaultRiskParameters(),
		balanceRepo:       balanceRepo,
		positionSizer:     positionSizer,
		drawdownMonitor:   drawdownMonitor,
		exposureMonitor:   exposureMonitor,
		dailyLimitMonitor: dailyLimitMonitor,
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

// GetRiskStatus returns the current risk status
func (rm *RiskManager) GetRiskStatus(ctx context.Context) (*risk.RiskStatus, error) {
	// Get account balance
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

	// Calculate daily P&L
	dailyPnL, err := rm.dailyLimitMonitor.CalculateDailyPnL(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate daily P&L: %w", err)
	}

	// Get risk parameters
	rm.lock.RLock()
	params := rm.riskParams
	rm.lock.RUnlock()

	// Calculate exposure percentage
	exposurePercent := 0.0
	if accountBalance > 0 {
		exposurePercent = (totalExposure / accountBalance) * 100
	}

	// Calculate daily P&L percentage
	dailyPnLPercent := 0.0
	if accountBalance > 0 {
		dailyPnLPercent = (dailyPnL / accountBalance) * 100
	}

	// Determine trading allowed status
	tradingAllowed := true
	var tradingDisabledReason string

	// Check minimum account balance
	if accountBalance < params.MinAccountBalance {
		tradingAllowed = false
		tradingDisabledReason = fmt.Sprintf("Account balance below minimum: %.2f < %.2f",
			accountBalance, params.MinAccountBalance)
	}

	// Check maximum drawdown
	if drawdown > params.MaxDrawdownPercent {
		tradingAllowed = false
		tradingDisabledReason = fmt.Sprintf("Drawdown exceeds maximum: %.2f%% > %.2f%%",
			drawdown, params.MaxDrawdownPercent)
	}

	// Check exposure limit
	if exposurePercent > params.MaxExposurePercent {
		tradingAllowed = false
		tradingDisabledReason = fmt.Sprintf("Exposure exceeds maximum: %.2f%% > %.2f%%",
			exposurePercent, params.MaxExposurePercent)
	}

	// Check daily loss limit
	if dailyPnLPercent < -params.DailyLossLimitPercent {
		tradingAllowed = false
		tradingDisabledReason = fmt.Sprintf("Daily loss exceeds limit: %.2f%% < -%.2f%%",
			dailyPnLPercent, params.DailyLossLimitPercent)
	}

	// Create risk status
	status := &risk.RiskStatus{
		AccountBalance:  accountBalance,
		TotalExposure:   totalExposure,
		CurrentDrawdown: drawdown,
		TodayPnL:        dailyPnL,
		TradingEnabled:  tradingAllowed,
		DisabledReason:  tradingDisabledReason,
		UpdatedAt:       time.Now(),
	}

	return status, nil
}

// UpdateRiskParameters updates the risk parameters
func (rm *RiskManager) UpdateRiskParameters(ctx context.Context, params risk.RiskParameters) error {
	// Validate parameters
	if params.MaxDrawdownPercent <= 0 || params.MaxDrawdownPercent > 100 {
		return fmt.Errorf("invalid max drawdown percent: must be between 0 and 100")
	}

	if params.RiskPerTradePercent <= 0 || params.RiskPerTradePercent > 100 {
		return fmt.Errorf("invalid risk per trade percent: must be between 0 and 100")
	}

	if params.MaxExposurePercent <= 0 || params.MaxExposurePercent > 100 {
		return fmt.Errorf("invalid max exposure percent: must be between 0 and 100")
	}

	if params.DailyLossLimitPercent <= 0 || params.DailyLossLimitPercent > 100 {
		return fmt.Errorf("invalid daily loss limit percent: must be between 0 and 100")
	}

	if params.MinAccountBalance <= 0 {
		return fmt.Errorf("invalid min account balance: must be greater than 0")
	}

	// Update parameters
	rm.lock.Lock()
	rm.riskParams = params
	rm.lock.Unlock()

	rm.logger.Info("Risk parameters updated",
		zap.Float64("max_drawdown_percent", params.MaxDrawdownPercent),
		zap.Float64("risk_per_trade_percent", params.RiskPerTradePercent),
		zap.Float64("max_exposure_percent", params.MaxExposurePercent),
		zap.Float64("daily_loss_limit_percent", params.DailyLossLimitPercent),
		zap.Float64("min_account_balance", params.MinAccountBalance))

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
	drawdown, err := rm.CalculateDrawdown(ctx)
	if err != nil {
		return false, "", fmt.Errorf("failed to calculate drawdown: %w", err)
	}

	if drawdown > maxDrawdownPercent {
		reason := fmt.Sprintf("Drawdown exceeds maximum: %.2f%% > %.2f%%", drawdown, maxDrawdownPercent)
		return false, reason, nil
	}

	// Check daily loss limit
	withinDailyLimit, err := rm.CheckDailyLossLimit(ctx)
	if err != nil {
		return false, "", fmt.Errorf("failed to check daily loss limit: %w", err)
	}

	if !withinDailyLimit {
		rm.lock.RLock()
		dailyLossLimitPercent := rm.riskParams.DailyLossLimitPercent
		rm.lock.RUnlock()

		dailyPnL, _ := rm.dailyLimitMonitor.CalculateDailyPnL(ctx)
		dailyPnLPercent := (dailyPnL / accountBalance) * 100

		reason := fmt.Sprintf("Daily loss exceeds limit: %.2f%% < -%.2f%%",
			dailyPnLPercent, dailyLossLimitPercent)
		return false, reason, nil
	}

	// Check exposure limit
	withinExposureLimit, err := rm.CheckExposureLimit(ctx, orderValue)
	if err != nil {
		return false, "", fmt.Errorf("failed to check exposure limit: %w", err)
	}

	if !withinExposureLimit {
		rm.lock.RLock()
		maxExposurePercent := rm.riskParams.MaxExposurePercent
		rm.lock.RUnlock()

		totalExposure, _ := rm.exposureMonitor.CalculateTotalExposure(ctx)
		newTotalExposure := totalExposure + orderValue
		exposurePercent := (newTotalExposure / accountBalance) * 100

		reason := fmt.Sprintf("Trade would exceed exposure limit: %.2f%% > %.2f%%",
			exposurePercent, maxExposurePercent)
		return false, reason, nil
	}

	return true, "", nil
}

// CalculateMaxDrawdown calculates the maximum drawdown from a series of balances
func (rm *RiskManager) CalculateMaxDrawdown(balances []float64) float64 {
	if len(balances) == 0 {
		return 0
	}

	maxDrawdown := 0.0
	peak := balances[0]

	for _, balance := range balances {
		if balance > peak {
			peak = balance
		}

		drawdown := (peak - balance) / peak * 100
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	return math.Round(maxDrawdown*100) / 100 // Round to 2 decimal places
}
