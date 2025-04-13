package ai

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/neo/crypto-bot/internal/domain/model"
)

// MockGeminiAIService is a mock implementation of the GeminiAIService for testing
type MockGeminiAIService struct{}

// NewMockGeminiAIService creates a new MockGeminiAIService
func NewMockGeminiAIService() *MockGeminiAIService {
	return &MockGeminiAIService{}
}

// Chat sends a message to the AI and returns a mock response
func (s *MockGeminiAIService) Chat(ctx context.Context, message string, conversationID string) (*model.AIMessage, error) {
	return &model.AIMessage{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        "This is a mock response from the AI assistant.",
		Timestamp:      time.Now(),
		Metadata: map[string]interface{}{
			"model": "mock-gemini",
		},
	}, nil
}

// ChatWithHistory sends a message with conversation history to the AI and returns a mock response
func (s *MockGeminiAIService) ChatWithHistory(ctx context.Context, messages []model.AIMessage) (*model.AIMessage, error) {
	conversationID := ""
	if len(messages) > 0 {
		conversationID = messages[0].ConversationID
	}

	return &model.AIMessage{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        "This is a mock response from the AI assistant with history context.",
		Timestamp:      time.Now(),
		Metadata: map[string]interface{}{
			"model": "mock-gemini",
		},
	}, nil
}

// GenerateInsight generates a mock insight
func (s *MockGeminiAIService) GenerateInsight(ctx context.Context, insightType string, data map[string]interface{}) (*model.AIInsight, error) {
	return &model.AIInsight{
		ID:          uuid.New().String(),
		UserID:      data["user_id"].(string),
		Type:        insightType,
		Title:       "Mock Insight Title",
		Description: "This is a mock insight description for testing purposes.",
		CreatedAt:   time.Now(),
		Confidence:  0.85,
		Metadata: map[string]interface{}{
			"model": "mock-gemini",
			"data":  data,
		},
	}, nil
}

// GenerateTradeRecommendation generates a mock trade recommendation
func (s *MockGeminiAIService) GenerateTradeRecommendation(ctx context.Context, data map[string]interface{}) (*model.AITradeRecommendation, error) {
	return &model.AITradeRecommendation{
		ID:         uuid.New().String(),
		UserID:     data["user_id"].(string),
		Symbol:     "BTC/USDT",
		Action:     "buy",
		Quantity:   0.01,
		Reasoning:  "This is a mock trade recommendation for testing purposes.",
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		Confidence: 0.75,
		Status:     "pending",
	}, nil
}

// ExecuteFunction executes a mock function call
func (s *MockGeminiAIService) ExecuteFunction(ctx context.Context, functionCall model.AIFunctionCall) (*model.AIFunctionResponse, error) {
	return &model.AIFunctionResponse{
		Name:   functionCall.Name,
		Result: map[string]interface{}{"message": "Mock function execution result"},
	}, nil
}

// GenerateEmbedding generates a mock vector embedding
func (s *MockGeminiAIService) GenerateEmbedding(ctx context.Context, text string) (*model.AIEmbedding, error) {
	return &model.AIEmbedding{
		ID:         uuid.New().String(),
		SourceID:   uuid.New().String(),
		SourceType: "message",
		Vector:     make([]float64, 128), // Mock 128-dimensional vector
		CreatedAt:  time.Now(),
	}, nil
}

// RegisterFunction registers a mock function handler
func (s *MockGeminiAIService) RegisterFunction(name string, handler FunctionHandler) {
	// Do nothing in the mock implementation
}
