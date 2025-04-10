package services

import (
	"context"
	"time"
)

// RiskParameters represents risk management settings
type RiskParameters struct {
	MaxPositionSize   float64 `json:"max_position_size"`
	MaxDrawdown       float64 `json:"max_drawdown"`
	MaxDailyLoss      float64 `json:"max_daily_loss"`
	MaxLeverage       float64 `json:"max_leverage"`
	StopLossPercent   float64 `json:"stop_loss_percent"`
	TakeProfitPercent float64 `json:"take_profit_percent"`
}

// PositionRisk represents risk metrics for a position
type PositionRisk struct {
	PositionID     string  `json:"position_id"`
	Symbol         string  `json:"symbol"`
	CurrentRisk    float64 `json:"current_risk"`
	MaxLoss        float64 `json:"max_loss"`
	RiskReward     float64 `json:"risk_reward"`
	BreakEvenPrice float64 `json:"break_even_price"`
}

// RiskAlertType represents the type of risk alert
type RiskAlertType string

const (
	RiskAlertTypeStopLoss   RiskAlertType = "STOP_LOSS"
	RiskAlertTypeTakeProfit RiskAlertType = "TAKE_PROFIT"
	RiskAlertTypeDrawdown   RiskAlertType = "DRAWDOWN"
	RiskAlertTypeDailyLoss  RiskAlertType = "DAILY_LOSS"
)

// RiskAlert represents a risk management alert
type RiskAlert struct {
	ID          string        `json:"id"`
	Symbol      string        `json:"symbol"`
	Price       float64       `json:"price"`
	Type        RiskAlertType `json:"type"`
	Active      bool          `json:"active"`
	CreatedAt   time.Time     `json:"created_at"`
	TriggeredAt *time.Time    `json:"triggered_at,omitempty"`
}

// RiskService defines the interface for risk management operations
type RiskService interface {
	// Risk parameters
	GetRiskParameters(ctx context.Context) (*RiskParameters, error)
	UpdateRiskParameters(ctx context.Context, params *RiskParameters) error

	// Position risk assessment
	ValidateNewPosition(ctx context.Context, symbol string, size float64, price float64) error
	CalculatePositionRisk(ctx context.Context, positionID string) (*PositionRisk, error)
	GetAllPositionsRisk(ctx context.Context) ([]PositionRisk, error)

	// Portfolio risk metrics
	GetPortfolioValue(ctx context.Context) (float64, error)
	GetPortfolioRisk(ctx context.Context) (float64, error)
	GetDailyPnL(ctx context.Context) (float64, error)
	GetDrawdown(ctx context.Context) (float64, error)

	// Risk alerts
	SetRiskAlert(ctx context.Context, symbol string, price float64, alertType RiskAlertType) error
	GetActiveAlerts(ctx context.Context) ([]RiskAlert, error)
}
