package backtest

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/risk"
	"go.uber.org/zap"
)

// RiskManager defines the interface for risk management in backtesting
type RiskManager interface {
	// CalculatePositionSize calculates the position size based on risk parameters
	CalculatePositionSize(ctx context.Context, symbol string, price float64, accountBalance float64) (float64, error)

	// IsTradeAllowed checks if a trade is allowed based on risk parameters
	IsTradeAllowed(ctx context.Context, symbol string, orderValue float64) (bool, string, error)

	// UpdateRiskParameters updates the risk parameters
	UpdateRiskParameters(ctx context.Context, params risk.RiskParameters) error

	// GetRiskStatus returns the current risk status
	GetRiskStatus(ctx context.Context) (*risk.RiskStatus, error)

	// UpdateAccountBalance updates the account balance and recalculates risk metrics
	UpdateAccountBalance(balance float64, timestamp time.Time)

	// GetBalanceHistory returns the balance history
	GetBalanceHistory() []BalanceRecord

	// GetDrawdownHistory returns the drawdown history
	GetDrawdownHistory() []DrawdownRecord
}

// BacktestRiskManager implements the RiskManager interface for backtesting
type BacktestRiskManager struct {
	riskParams      risk.RiskParameters
	accountBalance  float64
	totalExposure   float64
	drawdown        float64
	dailyPnL        float64
	balanceHistory  []BalanceRecord
	tradingEnabled  bool
	disabledReason  string
	logger          *zap.Logger
	positionTracker PositionTracker
}

// BalanceRecord represents a historical balance record
type BalanceRecord struct {
	Timestamp time.Time
	Balance   float64
}

// DrawdownRecord represents a historical drawdown record
type DrawdownRecord struct {
	Timestamp time.Time
	Drawdown  float64
}

// NewBacktestRiskManager creates a new BacktestRiskManager
func NewBacktestRiskManager(initialBalance float64, params risk.RiskParameters, positionTracker PositionTracker, logger *zap.Logger) *BacktestRiskManager {
	if logger == nil {
		logger, _ = zap.NewDevelopment()
	}

	return &BacktestRiskManager{
		riskParams:      params,
		accountBalance:  initialBalance,
		totalExposure:   0,
		drawdown:        0,
		dailyPnL:        0,
		balanceHistory:  []BalanceRecord{{Timestamp: time.Now(), Balance: initialBalance}},
		tradingEnabled:  true,
		disabledReason:  "",
		logger:          logger,
		positionTracker: positionTracker,
	}
}

// CalculatePositionSize calculates the position size based on risk parameters
func (r *BacktestRiskManager) CalculatePositionSize(ctx context.Context, symbol string, price float64, accountBalance float64) (float64, error) {
	// Calculate position size based on risk per trade
	riskAmount := accountBalance * r.riskParams.RiskPerTradePercent / 100
	stopLossPercent := 5.0 // Default stop loss percent

	// Calculate position size
	positionSize := riskAmount / (price * stopLossPercent / 100)

	r.logger.Debug("Calculated position size",
		zap.String("symbol", symbol),
		zap.Float64("price", price),
		zap.Float64("account_balance", accountBalance),
		zap.Float64("risk_percent", r.riskParams.RiskPerTradePercent),
		zap.Float64("risk_amount", riskAmount),
		zap.Float64("stop_loss_percent", stopLossPercent),
		zap.Float64("position_size", positionSize),
	)

	return positionSize, nil
}

// IsTradeAllowed checks if a trade is allowed based on risk parameters
func (r *BacktestRiskManager) IsTradeAllowed(ctx context.Context, symbol string, orderValue float64) (bool, string, error) {
	// Check if trading is enabled
	if !r.tradingEnabled {
		return false, r.disabledReason, nil
	}

	// Check minimum account balance
	if r.accountBalance < r.riskParams.MinAccountBalance {
		r.disabledReason = "Account balance below minimum"
		r.tradingEnabled = false
		return false, r.disabledReason, nil
	}

	// Check maximum drawdown
	if r.drawdown > r.riskParams.MaxDrawdownPercent {
		r.disabledReason = "Maximum drawdown exceeded"
		r.tradingEnabled = false
		return false, r.disabledReason, nil
	}

	// Check daily loss limit
	if r.dailyPnL < 0 && -r.dailyPnL/r.accountBalance*100 > r.riskParams.DailyLossLimitPercent {
		r.disabledReason = "Daily loss limit exceeded"
		r.tradingEnabled = false
		return false, r.disabledReason, nil
	}

	// Check maximum exposure
	newExposure := r.totalExposure + orderValue
	if newExposure/r.accountBalance*100 > r.riskParams.MaxExposurePercent {
		return false, "Maximum exposure exceeded", nil
	}

	return true, "", nil
}

// UpdateRiskParameters updates the risk parameters
func (r *BacktestRiskManager) UpdateRiskParameters(ctx context.Context, params risk.RiskParameters) error {
	r.riskParams = params
	return nil
}

// GetRiskStatus returns the current risk status
func (r *BacktestRiskManager) GetRiskStatus(ctx context.Context) (*risk.RiskStatus, error) {
	return &risk.RiskStatus{
		AccountBalance:  r.accountBalance,
		TotalExposure:   r.totalExposure,
		CurrentDrawdown: r.drawdown,
		TodayPnL:        r.dailyPnL,
		TradingEnabled:  r.tradingEnabled,
		DisabledReason:  r.disabledReason,
		UpdatedAt:       time.Now(),
	}, nil
}

// UpdateAccountBalance updates the account balance and recalculates risk metrics
func (r *BacktestRiskManager) UpdateAccountBalance(balance float64, timestamp time.Time) {
	// Update account balance
	r.accountBalance = balance

	// Add balance record
	r.balanceHistory = append(r.balanceHistory, BalanceRecord{
		Timestamp: timestamp,
		Balance:   balance,
	})

	// Calculate drawdown
	highWaterMark := r.balanceHistory[0].Balance
	for _, record := range r.balanceHistory {
		if record.Balance > highWaterMark {
			highWaterMark = record.Balance
		}
	}
	r.drawdown = (highWaterMark - balance) / highWaterMark * 100

	// Calculate daily P&L
	startOfDay := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
	var startBalance float64
	for i := len(r.balanceHistory) - 1; i >= 0; i-- {
		if r.balanceHistory[i].Timestamp.Before(startOfDay) {
			startBalance = r.balanceHistory[i].Balance
			break
		}
	}
	if startBalance == 0 {
		startBalance = r.balanceHistory[0].Balance
	}
	r.dailyPnL = balance - startBalance

	// Update total exposure
	r.totalExposure = 0
	for _, position := range r.positionTracker.GetOpenPositions() {
		r.totalExposure += position.Quantity * position.EntryPrice
	}

	// Check if trading should be enabled/disabled
	if r.accountBalance < r.riskParams.MinAccountBalance {
		r.tradingEnabled = false
		r.disabledReason = "Account balance below minimum"
	} else if r.drawdown > r.riskParams.MaxDrawdownPercent {
		r.tradingEnabled = false
		r.disabledReason = "Maximum drawdown exceeded"
	} else if r.dailyPnL < 0 && -r.dailyPnL/r.accountBalance*100 > r.riskParams.DailyLossLimitPercent {
		r.tradingEnabled = false
		r.disabledReason = "Daily loss limit exceeded"
	} else {
		r.tradingEnabled = true
		r.disabledReason = ""
	}
}

// GetBalanceHistory returns the balance history
func (r *BacktestRiskManager) GetBalanceHistory() []BalanceRecord {
	return r.balanceHistory
}

// GetDrawdownHistory returns the drawdown history
func (r *BacktestRiskManager) GetDrawdownHistory() []DrawdownRecord {
	drawdownHistory := make([]DrawdownRecord, 0, len(r.balanceHistory))

	if len(r.balanceHistory) == 0 {
		return drawdownHistory
	}

	highWaterMark := r.balanceHistory[0].Balance

	for _, record := range r.balanceHistory {
		if record.Balance > highWaterMark {
			highWaterMark = record.Balance
		}

		drawdown := (highWaterMark - record.Balance) / highWaterMark * 100
		drawdownHistory = append(drawdownHistory, DrawdownRecord{
			Timestamp: record.Timestamp,
			Drawdown:  drawdown,
		})
	}

	return drawdownHistory
}
