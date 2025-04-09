package service

import (
	"context"
	"fmt"
)

// TradeRequest represents a trade request from the AI
type TradeRequest struct {
	Symbol     string  `json:"symbol"`
	Action     string  `json:"action"`
	Amount     float64 `json:"amount"`
	PriceType  string  `json:"price_type"`
	LimitPrice float64 `json:"limit_price,omitempty"`
	StopLoss   float64 `json:"stop_loss,omitempty"`
	TakeProfit float64 `json:"take_profit,omitempty"`
}

// RiskAssessment represents the result of a risk assessment
type RiskAssessment struct {
	TradeAllowed bool     `json:"trade_allowed"`
	RiskFactors  []string `json:"risk_factors"`
	Explanation  string   `json:"explanation"`
}

// RiskService defines the interface for risk management
type RiskService interface {
	CheckDailyLossLimit(ctx context.Context, userID int) (*RiskCheck, error)
	CheckMaximumDrawdown(ctx context.Context, userID int) (*RiskCheck, error)
	CheckExposureLimit(ctx context.Context, userID int, symbol string) (*RiskCheck, error)
}

// RiskCheck represents the result of a risk check
type RiskCheck struct {
	Allowed   bool    `json:"allowed"`
	Threshold float64 `json:"threshold"`
}

// ValidateAITradeWithRiskSystem validates a trade request against risk parameters
func ValidateAITradeWithRiskSystem(
	ctx context.Context,
	trade *TradeRequest,
	userID int,
	riskSvc RiskService,
) (*RiskAssessment, error) {
	assessment := &RiskAssessment{
		TradeAllowed: false,
		RiskFactors:  []string{},
		Explanation:  "",
	}

	// Check daily loss limit
	dailyLossCheck, err := riskSvc.CheckDailyLossLimit(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check daily loss limit: %w", err)
	}
	if !dailyLossCheck.Allowed {
		assessment.RiskFactors = append(assessment.RiskFactors, "daily_loss_limit_exceeded")
		assessment.Explanation += fmt.Sprintf("Daily loss limit of %.2f%% reached. ", dailyLossCheck.Threshold)
	}

	// Check maximum drawdown
	drawdownCheck, err := riskSvc.CheckMaximumDrawdown(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check maximum drawdown: %w", err)
	}
	if !drawdownCheck.Allowed {
		assessment.RiskFactors = append(assessment.RiskFactors, "max_drawdown_exceeded")
		assessment.Explanation += fmt.Sprintf("Maximum drawdown of %.2f%% reached. ", drawdownCheck.Threshold)
	}

	// Check exposure limit for the specific asset
	exposureCheck, err := riskSvc.CheckExposureLimit(ctx, userID, trade.Symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to check exposure limit: %w", err)
	}
	if !exposureCheck.Allowed {
		assessment.RiskFactors = append(assessment.RiskFactors, "exposure_limit_exceeded")
		assessment.Explanation += fmt.Sprintf("Exposure limit of %.2f%% for %s reached. ",
			exposureCheck.Threshold, trade.Symbol)
	}

	// Trade is allowed if no risk factors were triggered
	assessment.TradeAllowed = len(assessment.RiskFactors) == 0

	return assessment, nil
}
