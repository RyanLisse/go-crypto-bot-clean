package gemini

import (
	"context"
	"fmt"

	"go-crypto-bot-clean/backend/internal/domain/ai/service"
)

// SetRiskGuardrails sets the risk guardrails for the AI service
func (s *GeminiAIService) SetRiskGuardrails(guardrails *service.AIRiskGuardrails) {
	s.RiskGuardrails = guardrails
}

// SetConfirmationFlow sets the confirmation flow for the AI service
func (s *GeminiAIService) SetConfirmationFlow(flow *service.ConfirmationFlow) {
	s.ConfirmationFlow = flow
}

// ApplyRiskGuardrails applies risk guardrails to a trade recommendation
func (s *GeminiAIService) ApplyRiskGuardrails(
	ctx context.Context,
	userID int,
	recommendation *service.TradeRecommendation,
) (*service.GuardrailsResult, error) {
	if s.RiskGuardrails == nil {
		return nil, fmt.Errorf("risk guardrails not set")
	}

	// Get account balance from portfolio service
	portfolio, err := s.PortfolioSvc.GetPortfolio(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}

	// Apply guardrails
	return s.RiskGuardrails.ApplyGuardrails(ctx, userID, recommendation, portfolio.TotalValue)
}

// CreateTradeConfirmation creates a trade confirmation
func (s *GeminiAIService) CreateTradeConfirmation(
	ctx context.Context,
	userID int,
	trade *service.TradeRequest,
	recommendation *service.TradeRecommendation,
) (*service.TradeConfirmation, error) {
	if s.ConfirmationFlow == nil {
		return nil, fmt.Errorf("confirmation flow not set")
	}

	// Validate trade with risk system
	assessment, err := service.ValidateAITradeWithRiskSystem(ctx, trade, userID, s.RiskSvc)
	if err != nil {
		return nil, fmt.Errorf("failed to validate trade: %w", err)
	}

	// Check if trade requires confirmation
	requiresConfirmation, reason := s.ConfirmationFlow.RequiresConfirmation(
		ctx, trade, recommendation, assessment)
	
	if !requiresConfirmation {
		return nil, nil // No confirmation required
	}

	// Create confirmation
	return s.ConfirmationFlow.CreateConfirmation(
		ctx, userID, trade, recommendation, assessment, reason)
}

// ConfirmTrade confirms a trade
func (s *GeminiAIService) ConfirmTrade(
	ctx context.Context,
	confirmationID string,
	approve bool,
) (*service.TradeConfirmation, error) {
	if s.ConfirmationFlow == nil {
		return nil, fmt.Errorf("confirmation flow not set")
	}

	return s.ConfirmationFlow.ConfirmTrade(ctx, confirmationID, approve)
}

// ListPendingTradeConfirmations lists all pending trade confirmations for a user
func (s *GeminiAIService) ListPendingTradeConfirmations(
	ctx context.Context,
	userID int,
) ([]*service.TradeConfirmation, error) {
	if s.ConfirmationFlow == nil {
		return nil, fmt.Errorf("confirmation flow not set")
	}

	return s.ConfirmationFlow.ListPendingConfirmations(ctx, userID)
}
