package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/ai/service"
	"go-crypto-bot-clean/backend/internal/domain/audit"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GenerateInsights generates AI insights based on portfolio and trade history
func (s *GeminiAIService) GenerateInsights(
	ctx context.Context,
	userID int,
	portfolio map[string]interface{},
	tradeHistory []map[string]interface{},
	insightTypes []string,
) ([]service.Insight, error) {
	// Create audit event
	if s.AuditService != nil {
		event, err := audit.CreateAuditEvent(
			userID,
			audit.EventTypeAI,
			audit.EventSeverityInfo,
			"GENERATE_INSIGHTS",
			"User requested AI insights",
			map[string]interface{}{
				"insight_types": insightTypes,
			},
			"", // IP will be added by middleware
			"", // User agent will be added by middleware
			"", // Request ID will be added by middleware
		)
		if err == nil {
			s.AuditService.LogEvent(ctx, event)
		}
	}

	// Convert portfolio and trade history to JSON strings for the prompt
	portfolioJSON, err := json.Marshal(portfolio)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal portfolio: %w", err)
	}

	tradeHistoryJSON, err := json.Marshal(tradeHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal trade history: %w", err)
	}

	// Create prompt for Gemini
	prompt := fmt.Sprintf(`You are an AI trading assistant for a cryptocurrency bot.

PORTFOLIO DATA:
%s

TRADE HISTORY:
%s

TASK:
Generate %d insights based on the portfolio and trade history data. Focus on the following insight types: %v.

For each insight, provide:
1. A concise title
2. A detailed description
3. The type of insight (one of: portfolio, market, opportunity)
4. The importance level (high, medium, low)
5. Relevant metrics with values and percentage changes where applicable
6. A specific recommendation when appropriate

Format your response as a JSON array of insight objects with the following structure:
[
  {
    "title": "Insight title",
    "description": "Detailed description",
    "type": "portfolio|market|opportunity",
    "importance": "high|medium|low",
    "metrics": [
      {
        "name": "Metric name",
        "value": "Metric value",
        "change": 10.5 // percentage change, can be positive or negative
      }
    ],
    "recommendation": "Specific recommendation" // optional
  }
]

Ensure your insights are data-driven, actionable, and relevant to the user's portfolio and trading history.`,
		string(portfolioJSON),
		string(tradeHistoryJSON),
		len(insightTypes)*2, // Generate 2 insights per type
		insightTypes,
	)

	// Call Gemini API
	model := s.Client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		s.Logger.Error("Failed to generate insights",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		return nil, fmt.Errorf("failed to generate insights: %w", err)
	}

	// Extract response text
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	responseText, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	// Parse JSON response
	var rawInsights []map[string]interface{}
	if err := json.Unmarshal([]byte(responseText), &rawInsights); err != nil {
		s.Logger.Error("Failed to parse insights JSON",
			zap.Error(err),
			zap.String("response", string(responseText)),
		)
		return nil, fmt.Errorf("failed to parse insights: %w", err)
	}

	// Convert to Insight objects
	insights := make([]service.Insight, 0, len(rawInsights))
	for _, raw := range rawInsights {
		insight := service.Insight{
			ID:             uuid.New().String(),
			Title:          getStringValue(raw, "title"),
			Description:    getStringValue(raw, "description"),
			Type:           getStringValue(raw, "type"),
			Importance:     getStringValue(raw, "importance"),
			Timestamp:      time.Now(),
			Recommendation: getStringValue(raw, "recommendation"),
		}

		// Parse metrics
		if metricsRaw, ok := raw["metrics"].([]interface{}); ok {
			insight.Metrics = make([]service.Metric, 0, len(metricsRaw))
			for _, metricRaw := range metricsRaw {
				if metricMap, ok := metricRaw.(map[string]interface{}); ok {
					metric := service.Metric{
						Name:  getStringValue(metricMap, "name"),
						Value: getStringValue(metricMap, "value"),
					}

					// Get change value if present
					if change, ok := metricMap["change"].(float64); ok {
						metric.Change = change
					}

					insight.Metrics = append(insight.Metrics, metric)
				}
			}
		}

		insights = append(insights, insight)
	}

	// Validate and sanitize insights
	if s.SecurityService != nil {
		for i, insight := range insights {
			// Sanitize title
			sanitizedTitle, err := s.SecurityService.SanitizeInput(ctx, insight.Title)
			if err == nil {
				insights[i].Title = sanitizedTitle
			}

			// Sanitize description
			sanitizedDesc, err := s.SecurityService.SanitizeInput(ctx, insight.Description)
			if err == nil {
				insights[i].Description = sanitizedDesc
			}

			// Sanitize recommendation
			if insight.Recommendation != "" {
				sanitizedRec, err := s.SecurityService.SanitizeInput(ctx, insight.Recommendation)
				if err == nil {
					insights[i].Recommendation = sanitizedRec
				}
			}
		}
	}

	return insights, nil
}

// Helper function to get string value from map
func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}
