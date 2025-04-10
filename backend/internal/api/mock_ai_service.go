package api

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/ai/service"
	"go-crypto-bot-clean/backend/internal/domain/ai/service/function"
	"go-crypto-bot-clean/backend/internal/domain/ai/service/templates"
	"go-crypto-bot-clean/backend/internal/domain/ai/types"
)

// MockAIService is a mock implementation of the AIService interface
type MockAIService struct{}

// GenerateResponse generates an AI response based on user message and context
func (m *MockAIService) GenerateResponse(ctx context.Context, userID int, message string) (string, error) {
	// Return a mock response
	return fmt.Sprintf("This is a mock AI response to: %s", message), nil
}

// ExecuteFunction allows the AI to call predefined functions
func (m *MockAIService) ExecuteFunction(ctx context.Context, userID int, functionName string, parameters map[string]interface{}) (interface{}, error) {
	// Return a mock result
	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Mock function %s executed with parameters: %v", functionName, parameters),
	}, nil
}

// StoreConversation saves a conversation to the database
func (m *MockAIService) StoreConversation(ctx context.Context, userID int, sessionID string, messages []service.Message) error {
	// Do nothing in the mock implementation
	return nil
}

// RetrieveConversation gets a conversation from the database
func (m *MockAIService) RetrieveConversation(ctx context.Context, userID int, sessionID string) (*service.ConversationMemory, error) {
	// Return a mock conversation
	return &service.ConversationMemory{
		UserID:    userID,
		SessionID: sessionID,
		Messages: []service.Message{
			{
				Role:      "user",
				Content:   "Hello, AI!",
				Timestamp: time.Now().Add(-5 * time.Minute),
			},
			{
				Role:      "assistant",
				Content:   "Hello! How can I help you with your trading today?",
				Timestamp: time.Now().Add(-4 * time.Minute),
			},
		},
		Summary:      "Mock conversation summary",
		LastAccessed: time.Now(),
	}, nil
}

// GenerateResponseWithTemplate generates an AI response using a specific template
func (m *MockAIService) GenerateResponseWithTemplate(ctx context.Context, userID int, templateName string, templateData templates.TemplateData) (string, error) {
	// Return a mock response
	return fmt.Sprintf("This is a mock AI response using template %s", templateName), nil
}

// GetAvailableFunctions returns all available functions
func (m *MockAIService) GetAvailableFunctions(ctx context.Context) []function.FunctionDefinition {
	// Return mock functions
	return []function.FunctionDefinition{
		{
			Name:        "get_market_data",
			Description: "Get current market data for a specific cryptocurrency",
			Parameters:  map[string]interface{}{},
			Required:    []string{},
		},
		{
			Name:        "execute_trade",
			Description: "Execute a trade",
			Parameters:  map[string]interface{}{},
			Required:    []string{},
		},
	}
}

// GetAvailableTemplates returns all available templates
func (m *MockAIService) GetAvailableTemplates(ctx context.Context) []string {
	// Return mock templates
	return []string{"trade_recommendation", "market_analysis", "portfolio_optimization"}
}

// ListUserSessions lists all sessions for a user
func (m *MockAIService) ListUserSessions(ctx context.Context, userID int, limit int) ([]string, error) {
	// Return mock sessions
	return []string{"session1", "session2", "session3"}, nil
}

// DeleteSession deletes a session
func (m *MockAIService) DeleteSession(ctx context.Context, userID int, sessionID string) error {
	// Do nothing in the mock implementation
	return nil
}

// FindSimilarMessages finds messages similar to the given query
func (m *MockAIService) FindSimilarMessages(ctx context.Context, query string, limit int) ([]types.SimilarMessage, error) {
	// Return mock similar messages
	return []types.SimilarMessage{
		{
			ConversationID: "mock_conversation_1",
			MessageID:      "mock_message_1",
			Content:        "This is a mock similar message 1",
			Similarity:     0.95,
			Metadata: map[string]interface{}{
				"indexed_at": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			},
		},
		{
			ConversationID: "mock_conversation_2",
			MessageID:      "mock_message_2",
			Content:        "This is a mock similar message 2",
			Similarity:     0.85,
			Metadata: map[string]interface{}{
				"indexed_at": time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			},
		},
	}, nil
}

// IndexMessage indexes a message for similarity search
func (m *MockAIService) IndexMessage(ctx context.Context, conversationID, messageID, content string, metadata map[string]interface{}) error {
	// Do nothing in the mock implementation
	return nil
}

// GenerateInsights generates AI insights based on portfolio and trade history
func (m *MockAIService) GenerateInsights(ctx context.Context, userID int, portfolio map[string]interface{}, tradeHistory []map[string]interface{}, insightTypes []string) ([]service.Insight, error) {
	// Return mock insights
	return []service.Insight{
		{
			ID:             "insight_1",
			Type:           "portfolio_allocation",
			Title:          "Portfolio Allocation Insight",
			Description:    "Your portfolio is well-diversified across major cryptocurrencies.",
			Importance:     "medium",
			Timestamp:      time.Now(),
			Recommendation: "Continue with current allocation strategy.",
		},
		{
			ID:             "insight_2",
			Type:           "risk_assessment",
			Title:          "Risk Assessment Insight",
			Description:    "Your current risk level is moderate based on your trading history.",
			Importance:     "high",
			Timestamp:      time.Now(),
			Recommendation: "Consider reducing exposure to volatile assets.",
		},
	}, nil
}

// ApplyRiskGuardrails applies risk guardrails to a trade recommendation
func (m *MockAIService) ApplyRiskGuardrails(ctx context.Context, userID int, recommendation *service.TradeRecommendation) (*service.GuardrailsResult, error) {
	// Return a mock result
	return &service.GuardrailsResult{
		OriginalRecommendation: recommendation,
		ModifiedRecommendation: recommendation,
		Modifications:          []string{},
		Timestamp:              time.Now(),
	}, nil
}

// CreateTradeConfirmation creates a trade confirmation
func (m *MockAIService) CreateTradeConfirmation(ctx context.Context, userID int, trade *service.TradeRequest, recommendation *service.TradeRecommendation) (*service.TradeConfirmation, error) {
	// Return a mock confirmation
	return &service.TradeConfirmation{
		ID:                 "mock_confirmation_id",
		UserID:             userID,
		TradeRequest:       trade,
		Recommendation:     recommendation,
		Status:             service.ConfirmationPending,
		ConfirmationReason: "Mock confirmation reason",
		CreatedAt:          time.Now(),
		ExpiresAt:          time.Now().Add(24 * time.Hour),
	}, nil
}

// ConfirmTrade confirms a trade
func (m *MockAIService) ConfirmTrade(ctx context.Context, confirmationID string, approve bool) (*service.TradeConfirmation, error) {
	// Return a mock confirmation
	status := service.ConfirmationApproved
	if !approve {
		status = service.ConfirmationRejected
	}

	now := time.Now()
	return &service.TradeConfirmation{
		ID:                 confirmationID,
		UserID:             1,
		Status:             status,
		ConfirmationReason: "Mock confirmation reason",
		CreatedAt:          now.Add(-1 * time.Hour),
		ExpiresAt:          now.Add(23 * time.Hour),
		ConfirmedAt:        &now,
	}, nil
}

// ListPendingTradeConfirmations lists all pending trade confirmations for a user
func (m *MockAIService) ListPendingTradeConfirmations(ctx context.Context, userID int) ([]*service.TradeConfirmation, error) {
	// Return mock confirmations
	now := time.Now()
	return []*service.TradeConfirmation{
		{
			ID:                 "mock_confirmation_id_1",
			UserID:             userID,
			Status:             service.ConfirmationPending,
			ConfirmationReason: "Mock confirmation reason 1",
			CreatedAt:          now.Add(-2 * time.Hour),
			ExpiresAt:          now.Add(22 * time.Hour),
		},
		{
			ID:                 "mock_confirmation_id_2",
			UserID:             userID,
			Status:             service.ConfirmationPending,
			ConfirmationReason: "Mock confirmation reason 2",
			CreatedAt:          now.Add(-1 * time.Hour),
			ExpiresAt:          now.Add(23 * time.Hour),
		},
	}, nil
}
