package ai

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// StubAIService is a stub implementation of the AIService interface
type StubAIService struct {
	logger zerolog.Logger
}

// NewStubAIService creates a new StubAIService
func NewStubAIService(cfg *config.Config, logger zerolog.Logger) (*StubAIService, error) {
	return &StubAIService{
		logger: logger.With().Str("component", "stub_ai_service").Logger(),
	}, nil
}

// Chat sends a message to the AI and returns a response
func (s *StubAIService) Chat(ctx context.Context, message string, conversationID string) (*model.AIMessage, error) {
	return &model.AIMessage{
		ID:        uuid.NewString(),
		Role:      "assistant",
		Content:   "[Stub AI response]",
		Timestamp: time.Now(),
		Metadata:  map[string]interface{}{"stub": true},
	}, nil
}

// ChatWithHistory continues a conversation with a user based on message history
func (s *StubAIService) ChatWithHistory(ctx context.Context, messages []model.AIMessage, tradingContext map[string]interface{}) (*model.AIMessage, error) {
	return &model.AIMessage{
		ID:        uuid.NewString(),
		Role:      "assistant",
		Content:   "[Stub AI response with history]",
		Timestamp: time.Now(),
		Metadata:  map[string]interface{}{"stub": true},
	}, nil
}

// GenerateInsight generates an insight based on provided data
func (s *StubAIService) GenerateInsight(ctx context.Context, insightType string, data map[string]interface{}) (*model.AIInsight, error) {
	return &model.AIInsight{
		ID:          uuid.NewString(),
		UserID:      "stub-user",
		Type:        insightType,
		Title:       "Stub Insight",
		Description: "[Stub insight description]",
		CreatedAt:   time.Now(),
		Confidence:  0.9,
		Metadata:    map[string]interface{}{"stub": true},
	}, nil
}

// GenerateTradeRecommendation generates a trade recommendation
func (s *StubAIService) GenerateTradeRecommendation(ctx context.Context, data map[string]interface{}) (*model.AITradeRecommendation, error) {
	return &model.AITradeRecommendation{
		ID:         uuid.NewString(),
		UserID:     "stub-user",
		Symbol:     "BTC/USDT",
		Action:     "buy",
		Quantity:   0.1,
		Reasoning:  "[Stub trade recommendation reasoning]",
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(1 * time.Hour),
		Confidence: 0.8,
		Status:     "pending",
		Metadata:   map[string]interface{}{"stub": true},
	}, nil
}

// ExecuteFunction executes a function call from the AI
func (s *StubAIService) ExecuteFunction(ctx context.Context, functionCall model.AIFunctionCall) (*model.AIFunctionResponse, error) {
	return &model.AIFunctionResponse{
		Name:   functionCall.Name,
		Result: "stub-result",
	}, nil
}

// GenerateEmbedding generates a vector embedding for a text
func (s *StubAIService) GenerateEmbedding(ctx context.Context, text string) (*model.AIEmbedding, error) {
	return nil, errors.New("embedding not supported in stub service")
}
