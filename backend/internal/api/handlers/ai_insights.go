package handlers

import (
	"encoding/json"
	"net/http"

	"go-crypto-bot-clean/backend/internal/domain/ai/service"
)

// InsightRequest represents a request for AI-generated insights
type InsightRequest struct {
	UserID      int                      `json:"user_id"`
	Portfolio   map[string]interface{}   `json:"portfolio"`
	TradeHistory []map[string]interface{} `json:"trade_history"`
	InsightTypes []string                 `json:"insight_types"`
}

// InsightResponse represents the response with AI-generated insights
type InsightResponse struct {
	Insights []service.Insight `json:"insights"`
}

// AIInsightsHandler handles requests for AI-generated insights
func AIInsightsHandler(aiSvc service.AIService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var req InsightRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validate request
		if req.Portfolio == nil {
			http.Error(w, "Missing portfolio data", http.StatusBadRequest)
			return
		}

		// Generate insights using the AI service
		insights, err := aiSvc.GenerateInsights(r.Context(), req.UserID, req.Portfolio, req.TradeHistory, req.InsightTypes)
		if err != nil {
			http.Error(w, "Failed to generate insights: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(InsightResponse{
			Insights: insights,
		})
	}
}
