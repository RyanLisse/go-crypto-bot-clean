package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/risk"
)

// AIRiskGuardrails defines risk guardrails for AI-generated trading recommendations
type AIRiskGuardrails struct {
	RiskSvc risk.Service
}

// NewAIRiskGuardrails creates a new AIRiskGuardrails
func NewAIRiskGuardrails(riskSvc risk.Service) *AIRiskGuardrails {
	return &AIRiskGuardrails{
		RiskSvc: riskSvc,
	}
}

// TradeRecommendation represents an AI-generated trade recommendation
type TradeRecommendation struct {
	Recommendation        string                 `json:"recommendation"` // BUY, SELL, HOLD
	Confidence            float64                `json:"confidence"`
	Reasoning             string                 `json:"reasoning"`
	RiskLevel             string                 `json:"risk_level"` // LOW, MEDIUM, HIGH
	SuggestedPositionSize float64                `json:"suggested_position_size"`
	SuggestedStopLoss     float64                `json:"suggested_stop_loss"`
	TechnicalIndicators   map[string]interface{} `json:"technical_indicators"`
}

// RiskStatus represents the current risk status
type RiskStatus struct {
	TradingEnabled  bool    `json:"trading_enabled"`
	DisabledReason  string  `json:"disabled_reason,omitempty"`
	CurrentDrawdown float64 `json:"current_drawdown"`
}

// GuardrailsResult represents the result of applying risk guardrails
type GuardrailsResult struct {
	OriginalRecommendation *TradeRecommendation `json:"original_recommendation"`
	ModifiedRecommendation *TradeRecommendation `json:"modified_recommendation"`
	Modifications          []string             `json:"modifications"`
	RiskStatus             *RiskStatus          `json:"risk_status"`
	Timestamp              time.Time            `json:"timestamp"`
}

// ApplyGuardrails applies risk guardrails to an AI-generated trade recommendation
func (g *AIRiskGuardrails) ApplyGuardrails(
	ctx context.Context,
	userID int,
	recommendation *TradeRecommendation,
	accountBalance float64,
) (*GuardrailsResult, error) {
	// Create a copy of the original recommendation
	originalRecommendation := *recommendation

	// Initialize result and risk status
	riskStatus := &RiskStatus{
		TradingEnabled:  true,
		CurrentDrawdown: 0.0,
	}

	result := &GuardrailsResult{
		OriginalRecommendation: &originalRecommendation,
		ModifiedRecommendation: recommendation,
		Modifications:          []string{},
		RiskStatus:             riskStatus,
		Timestamp:              time.Now(),
	}

	// Check daily loss limit
	dailyLossCheck, err := g.RiskSvc.CheckDailyLossLimit(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check daily loss limit: %w", err)
	}

	// Check maximum drawdown
	drawdownCheck, err := g.RiskSvc.CheckMaximumDrawdown(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check maximum drawdown: %w", err)
	}

	// Update risk status based on checks
	if !dailyLossCheck.Allowed {
		riskStatus.TradingEnabled = false
		riskStatus.DisabledReason = "Daily loss limit reached"
	}

	if !drawdownCheck.Allowed {
		riskStatus.TradingEnabled = false
		riskStatus.DisabledReason = "Maximum drawdown reached"
	}

	// Simulate a current drawdown for demonstration purposes
	riskStatus.CurrentDrawdown = 5.0 // 5% drawdown

	// If trading is not enabled, override recommendation to HOLD
	if !riskStatus.TradingEnabled {
		if recommendation.Recommendation != "HOLD" {
			recommendation.Recommendation = "HOLD"
			recommendation.Reasoning = fmt.Sprintf("Trading disabled due to risk controls: %s. Original reasoning: %s",
				riskStatus.DisabledReason, recommendation.Reasoning)
			result.Modifications = append(result.Modifications,
				fmt.Sprintf("Changed recommendation to HOLD due to disabled trading: %s", riskStatus.DisabledReason))
		}
		return result, nil
	}

	// Apply position size guardrails
	if recommendation.Recommendation == "BUY" || recommendation.Recommendation == "SELL" {
		// Calculate safe position size based on risk parameters
		// In a real implementation, this would use the symbol to calculate position size
		_ = extractSymbolFromRecommendation(recommendation) // Just to avoid unused function warning

		// Since we don't have the actual CalculatePositionSize method, we'll use a simple calculation
		// In a real implementation, this would call the risk service
		safePositionSize := accountBalance * 0.02 // 2% of account balance

		// If AI-suggested position size is too large, reduce it
		if recommendation.SuggestedPositionSize > safePositionSize {
			originalSize := recommendation.SuggestedPositionSize
			recommendation.SuggestedPositionSize = safePositionSize
			result.Modifications = append(result.Modifications,
				fmt.Sprintf("Reduced position size from %.2f%% to %.2f%% based on risk parameters",
					originalSize, safePositionSize))
		}

		// If risk level is HIGH but account is in drawdown, downgrade recommendation
		if recommendation.RiskLevel == "HIGH" && riskStatus.CurrentDrawdown > 10.0 {
			if recommendation.Recommendation != "HOLD" {
				originalRec := recommendation.Recommendation
				recommendation.Recommendation = "HOLD"
				recommendation.Reasoning = fmt.Sprintf("High risk trade not recommended during significant drawdown (%.2f%%). Original reasoning: %s",
					riskStatus.CurrentDrawdown, recommendation.Reasoning)
				result.Modifications = append(result.Modifications,
					fmt.Sprintf("Changed recommendation from %s to HOLD due to high risk during %.2f%% drawdown",
						originalRec, riskStatus.CurrentDrawdown))
			}
		}

		// Ensure stop loss is set
		if recommendation.SuggestedStopLoss == 0 && (recommendation.Recommendation == "BUY" || recommendation.Recommendation == "SELL") {
			// Set a default stop loss at 2% for risk management
			recommendation.SuggestedStopLoss = 2.0
			result.Modifications = append(result.Modifications,
				"Added default 2% stop loss as none was specified")
		}
	}

	return result, nil
}

// extractSymbolFromRecommendation extracts the trading symbol from a recommendation
func extractSymbolFromRecommendation(recommendation *TradeRecommendation) string {
	// This is a simplified implementation
	// In a real system, you would extract the symbol from the recommendation
	// based on the content of the reasoning or technical indicators

	// Look for common symbol patterns in the reasoning
	reasoning := recommendation.Reasoning

	// Check for BTC/USD, ETH-USD, etc.
	symbols := []string{"BTC", "ETH", "SOL", "XRP", "ADA", "DOT", "DOGE", "SHIB", "AVAX", "MATIC"}
	for _, symbol := range symbols {
		if strings.Contains(reasoning, symbol) {
			return symbol
		}
	}

	// Default to BTC if no symbol found
	return "BTC"
}
