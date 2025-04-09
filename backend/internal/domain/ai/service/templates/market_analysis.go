package templates

import (
	"errors"
	"fmt"
)

// MarketAnalysisData contains data for market analysis prompts
type MarketAnalysisData struct {
	MarketData       map[string]interface{}
	TimeFrame        string
	HistoricalTrends string
	UserQuery        string
}

// Validate validates the market analysis data
func (d *MarketAnalysisData) Validate() error {
	if d.MarketData == nil || len(d.MarketData) == 0 {
		return errors.New("market data is required")
	}
	if d.TimeFrame == "" {
		return errors.New("time frame is required")
	}
	return nil
}

// NewMarketAnalysisTemplate creates a new market analysis template
func NewMarketAnalysisTemplate() *BaseTemplate {
	const templateText = `You are an AI trading assistant for a cryptocurrency bot.

MARKET DATA:
{{range $key, $value := .MarketData}}
- {{$key}}: {{$value}}
{{end}}

TIME FRAME:
{{.TimeFrame}}

HISTORICAL TRENDS:
{{.HistoricalTrends}}

USER QUERY:
{{.UserQuery}}

TASK:
Analyze the current market conditions and provide insights based on the data provided.

OUTPUT FORMAT:
1. Market Summary: Provide a concise summary of the current market conditions.
2. Key Trends: Identify 2-3 key trends in the market data.
3. Risk Assessment: Evaluate the current market risk level (LOW, MEDIUM, HIGH).
4. Opportunities: Identify potential trading opportunities based on the data.
5. Recommendations: Provide actionable recommendations for the user.`

	return NewBaseTemplate(
		"market_analysis",
		"1.0.0",
		"Template for generating market analysis",
		templateText,
	)
}

// RenderMarketAnalysis renders a market analysis prompt
func RenderMarketAnalysis(
	marketData map[string]interface{},
	timeFrame string,
	historicalTrends string,
	userQuery string,
) (string, error) {
	template := NewMarketAnalysisTemplate()
	data := &MarketAnalysisData{
		MarketData:       marketData,
		TimeFrame:        timeFrame,
		HistoricalTrends: historicalTrends,
		UserQuery:        userQuery,
	}

	if err := data.Validate(); err != nil {
		return "", fmt.Errorf("market analysis data validation failed: %w", err)
	}

	return template.Render(data)
}
