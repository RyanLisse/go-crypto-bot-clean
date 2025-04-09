package handlers

import (
	"encoding/json"
	"net/http"

	"go-crypto-bot-clean/backend/internal/domain/ai/service"
)

// TradeRecommendationRequest represents a request to apply risk guardrails to a trade recommendation
type TradeRecommendationRequest struct {
	UserID         int                         `json:"user_id" doc:"User ID"`
	Recommendation service.TradeRecommendation `json:"recommendation" doc:"Trade recommendation"`
}

// TradeConfirmationRequest represents a request to create a trade confirmation
type TradeConfirmationRequest struct {
	UserID         int                         `json:"user_id" doc:"User ID"`
	Trade          service.TradeRequest        `json:"trade" doc:"Trade request"`
	Recommendation service.TradeRecommendation `json:"recommendation" doc:"Trade recommendation"`
}

// ConfirmTradeRequest represents a request to confirm a trade
type ConfirmTradeRequest struct {
	ConfirmationID string `json:"confirmation_id" doc:"Confirmation ID"`
	Approve        bool   `json:"approve" doc:"Whether to approve the trade"`
}

// ListPendingConfirmationsRequest represents a request to list pending trade confirmations
type ListPendingConfirmationsRequest struct {
	UserID int `json:"user_id" doc:"User ID"`
}

// ApplyRiskGuardrailsHandler handles requests to apply risk guardrails to a trade recommendation
func ApplyRiskGuardrailsHandler(aiSvc service.AIService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var input TradeRecommendationRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Apply risk guardrails
		result, err := aiSvc.ApplyRiskGuardrails(r.Context(), input.UserID, &input.Recommendation)
		if err != nil {
			http.Error(w, "Failed to apply risk guardrails: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

// CreateTradeConfirmationHandler handles requests to create a trade confirmation
func CreateTradeConfirmationHandler(aiSvc service.AIService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var input TradeConfirmationRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Create trade confirmation
		confirmation, err := aiSvc.CreateTradeConfirmation(r.Context(), input.UserID, &input.Trade, &input.Recommendation)
		if err != nil {
			http.Error(w, "Failed to create trade confirmation: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if confirmation == nil {
			// No confirmation required
			confirmation = &service.TradeConfirmation{
				ID:                 "",
				Status:             service.ConfirmationApproved,
				ConfirmationReason: "No confirmation required",
			}
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(confirmation)
	}
}

// ConfirmTradeHandler handles requests to confirm a trade
func ConfirmTradeHandler(aiSvc service.AIService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var input ConfirmTradeRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Confirm trade
		confirmation, err := aiSvc.ConfirmTrade(r.Context(), input.ConfirmationID, input.Approve)
		if err != nil {
			http.Error(w, "Failed to confirm trade: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(confirmation)
	}
}

// ListPendingTradeConfirmationsHandler handles requests to list pending trade confirmations
func ListPendingTradeConfirmationsHandler(aiSvc service.AIService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var input ListPendingConfirmationsRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// List pending trade confirmations
		confirmations, err := aiSvc.ListPendingTradeConfirmations(r.Context(), input.UserID)
		if err != nil {
			http.Error(w, "Failed to list pending trade confirmations: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(confirmations)
	}
}
