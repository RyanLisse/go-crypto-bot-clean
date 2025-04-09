package templates

import (
	"errors"
	"fmt"
)

// PortfolioOptimizationData contains data for portfolio optimization prompts
type PortfolioOptimizationData struct {
	CurrentPortfolio map[string]float64
	RiskProfile      string
	MarketOutlook    string
	UserGoals        string
	UserQuery        string
}

// Validate validates the portfolio optimization data
func (d *PortfolioOptimizationData) Validate() error {
	if d.CurrentPortfolio == nil || len(d.CurrentPortfolio) == 0 {
		return errors.New("current portfolio is required")
	}
	if d.RiskProfile == "" {
		return errors.New("risk profile is required")
	}
	return nil
}

// NewPortfolioOptimizationTemplate creates a new portfolio optimization template
func NewPortfolioOptimizationTemplate() *BaseTemplate {
	const templateText = `You are an AI trading assistant for a cryptocurrency bot.

CURRENT PORTFOLIO:
{{range $asset, $allocation := .CurrentPortfolio}}
- {{$asset}}: {{$allocation}}%
{{end}}

RISK PROFILE:
{{.RiskProfile}}

MARKET OUTLOOK:
{{.MarketOutlook}}

USER GOALS:
{{.UserGoals}}

USER QUERY:
{{.UserQuery}}

TASK:
Analyze the current portfolio and suggest optimizations based on the user's risk profile, market outlook, and goals.

OUTPUT FORMAT (JSON):
{
  "analysis": {
    "current_diversification": "LOW|MEDIUM|HIGH",
    "risk_assessment": "LOW|MEDIUM|HIGH",
    "performance_outlook": "Brief assessment of current portfolio performance outlook"
  },
  "recommendations": [
    {
      "action": "BUY|SELL|HOLD",
      "asset": "Asset symbol",
      "allocation_change": float,
      "reasoning": "Brief explanation"
    }
  ],
  "suggested_portfolio": {
    "asset1": float,
    "asset2": float
  },
  "expected_improvement": {
    "risk_reduction": float,
    "potential_return_increase": float
  }
}`

	return NewBaseTemplate(
		"portfolio_optimization",
		"1.0.0",
		"Template for generating portfolio optimization recommendations",
		templateText,
	)
}

// RenderPortfolioOptimization renders a portfolio optimization prompt
func RenderPortfolioOptimization(
	currentPortfolio map[string]float64,
	riskProfile string,
	marketOutlook string,
	userGoals string,
	userQuery string,
) (string, error) {
	template := NewPortfolioOptimizationTemplate()
	data := &PortfolioOptimizationData{
		CurrentPortfolio: currentPortfolio,
		RiskProfile:      riskProfile,
		MarketOutlook:    marketOutlook,
		UserGoals:        userGoals,
		UserQuery:        userQuery,
	}

	if err := data.Validate(); err != nil {
		return "", fmt.Errorf("portfolio optimization data validation failed: %w", err)
	}

	return template.Render(data)
}
