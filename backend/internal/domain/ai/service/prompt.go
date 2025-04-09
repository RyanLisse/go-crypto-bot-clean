package service

import "fmt"

// PromptTemplate defines a template for AI prompts
type PromptTemplate interface {
	GeneratePrompt() string
}

// TradeRecommendationPrompt is a template for trade recommendations
type TradeRecommendationPrompt struct {
	PortfolioContext string
	RiskProfile      string
	MarketConditions string
	TargetAsset      string
}

// GeneratePrompt generates a structured prompt for trade recommendations
func (t *TradeRecommendationPrompt) GeneratePrompt() string {
	return fmt.Sprintf(`You are an AI trading assistant for a cryptocurrency bot.
	
CONTEXT:
- Portfolio: %s
- User Risk Profile: %s
- Current Market Conditions: %s

TASK:
Analyze whether %s is a good trading opportunity right now.

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
}`, t.PortfolioContext, t.RiskProfile, t.MarketConditions, t.TargetAsset)
}
