package service

import (
	"context"

	"go-crypto-bot-clean/backend/internal/domain/risk"
)

// RiskServiceAdapter adapts the risk.Service interface to the RiskService interface
type RiskServiceAdapter struct {
	riskSvc risk.Service
}

// NewRiskServiceAdapter creates a new RiskServiceAdapter
func NewRiskServiceAdapter(riskSvc risk.Service) *RiskServiceAdapter {
	return &RiskServiceAdapter{
		riskSvc: riskSvc,
	}
}

// CheckDailyLossLimit checks if the daily loss limit has been reached
func (a *RiskServiceAdapter) CheckDailyLossLimit(ctx context.Context, userID int) (*RiskCheck, error) {
	result, err := a.riskSvc.CheckDailyLossLimit(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &RiskCheck{
		Allowed:   result.Allowed,
		Threshold: result.Threshold,
	}, nil
}

// CheckMaximumDrawdown checks if the maximum drawdown has been reached
func (a *RiskServiceAdapter) CheckMaximumDrawdown(ctx context.Context, userID int) (*RiskCheck, error) {
	result, err := a.riskSvc.CheckMaximumDrawdown(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &RiskCheck{
		Allowed:   result.Allowed,
		Threshold: result.Threshold,
	}, nil
}

// CheckExposureLimit checks if the exposure limit for a symbol has been reached
func (a *RiskServiceAdapter) CheckExposureLimit(ctx context.Context, userID int, symbol string) (*RiskCheck, error) {
	result, err := a.riskSvc.CheckExposureLimit(ctx, userID, symbol)
	if err != nil {
		return nil, err
	}
	return &RiskCheck{
		Allowed:   result.Allowed,
		Threshold: result.Threshold,
	}, nil
}
