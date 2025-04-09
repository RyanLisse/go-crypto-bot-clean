package templates

import (
	"errors"
	"fmt"
)

// TradeRecommendationData contains data for trade recommendation prompts
type TradeRecommendationData struct {
	PortfolioContext string
	RiskProfile      string
	MarketConditions string
	TargetAsset      string
	UserQuery        string
}

// Validate validates the trade recommendation data
func (d *TradeRecommendationData) Validate() error {
	if d.TargetAsset == "" {
		return errors.New("target asset is required")
	}
	return nil
}

// NewTradeRecommendationTemplate creates a new trade recommendation template
func NewTradeRecommendationTemplate() *BaseTemplate {
	const templateText = `You are an AI trading assistant for a cryptocurrency bot.

CONTEXT:
- Portfolio: {{.PortfolioContext}}
- User Risk Profile: {{.RiskProfile}}
- Current Market Conditions: {{.MarketConditions}}
- User Query: {{.UserQuery}}

TASK:
Analyze whether {{.TargetAsset}} is a good trading opportunity right now.

OUTPUT FORMAT (JSON):
{
  "recommendation": "BUY|SELL|HOLD",
  "confidence": 0.0-1.0,
  "reasoning": "Brief explanation of recommendation",
  "risk_level": "LOW|MEDIUM|HIGH",
  "suggested_position_size": 0.0-1.0,
  "suggested_stop_loss": float,
  "technical_indicators": {
    "rsi": float,
    "macd": "BULLISH|BEARISH|NEUTRAL",
    "support_level": float,
    "resistance_level": float
  }
}`

	return NewBaseTemplate(
		"trade_recommendation",
		"1.0.0",
		"Template for generating trade recommendations",
		templateText,
	)
}

// RenderTradeRecommendation renders a trade recommendation prompt
func RenderTradeRecommendation(
	portfolioContext, riskProfile, marketConditions, targetAsset, userQuery string,
) (string, error) {
	template := NewTradeRecommendationTemplate()
	data := &TradeRecommendationData{
		PortfolioContext: portfolioContext,
		RiskProfile:      riskProfile,
		MarketConditions: marketConditions,
		TargetAsset:      targetAsset,
		UserQuery:        userQuery,
	}

	if err := data.Validate(); err != nil {
		return "", fmt.Errorf("trade recommendation data validation failed: %w", err)
	}

	return template.Render(data)
}
