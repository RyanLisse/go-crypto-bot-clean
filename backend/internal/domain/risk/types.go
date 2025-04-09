// Package risk provides risk management functionality for the trading bot.
package risk

import (
	"context"
	"time"
)

// RiskParameters contains risk management configuration
type RiskParameters struct {
	MaxDrawdownPercent    float64 `json:"max_drawdown_percent"`
	RiskPerTradePercent   float64 `json:"risk_per_trade_percent"`
	MaxExposurePercent    float64 `json:"max_exposure_percent"`
	DailyLossLimitPercent float64 `json:"daily_loss_limit_percent"`
	MinAccountBalance     float64 `json:"min_account_balance"`
}

// RiskService defines the interface for risk management
type RiskService interface {
	// Calculation methods
	CalculatePositionSize(ctx context.Context, symbol string, accountBalance float64) (float64, error)
	CalculateDrawdown(ctx context.Context) (float64, error)

	// Limit checking methods
	CheckExposureLimit(ctx context.Context, newOrderValue float64) (bool, error)
	CheckDailyLossLimit(ctx context.Context) (bool, error)

	// Risk status methods
	GetRiskStatus(ctx context.Context) (*RiskStatus, error)
	UpdateRiskParameters(ctx context.Context, params RiskParameters) error

	// Monitoring methods
	IsTradeAllowed(ctx context.Context, symbol string, orderValue float64) (bool, string, error)
}

// RiskStatus represents the current risk metrics
type RiskStatus struct {
	CurrentDrawdown float64   `json:"current_drawdown"`
	TotalExposure   float64   `json:"total_exposure"`
	TodayPnL        float64   `json:"today_pnl"`
	AccountBalance  float64   `json:"account_balance"`
	TradingEnabled  bool      `json:"trading_enabled"`
	DisabledReason  string    `json:"disabled_reason,omitempty"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Import the BalanceHistory type from the controls package
// This is defined in controls/drawdown.go

// DefaultRiskParameters returns the default risk parameters
func DefaultRiskParameters() RiskParameters {
	return RiskParameters{
		MaxDrawdownPercent:    20.0,  // Maximum 20% drawdown before halting trading
		RiskPerTradePercent:   1.0,   // Risk 1% of account per trade
		MaxExposurePercent:    50.0,  // Maximum 50% of account in open positions
		DailyLossLimitPercent: 5.0,   // Stop trading if daily losses exceed 5%
		MinAccountBalance:     100.0, // Minimum account balance required
	}
}
