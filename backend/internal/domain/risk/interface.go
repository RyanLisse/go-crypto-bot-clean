package risk

import (
	"context"
)

// RiskCheck represents the result of a risk check
type RiskCheck struct {
	Allowed   bool    `json:"allowed"`
	Threshold float64 `json:"threshold"`
}

// Service defines the interface for risk management
type Service interface {
	// CheckDailyLossLimit checks if the daily loss limit has been reached
	CheckDailyLossLimit(ctx context.Context, userID int) (*RiskCheck, error)
	
	// CheckMaximumDrawdown checks if the maximum drawdown has been reached
	CheckMaximumDrawdown(ctx context.Context, userID int) (*RiskCheck, error)
	
	// CheckExposureLimit checks if the exposure limit for a symbol has been reached
	CheckExposureLimit(ctx context.Context, userID int, symbol string) (*RiskCheck, error)
}
