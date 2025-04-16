package mocks

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

type MockAIService struct{}

func (m *MockAIService) ChatWithHistory(ctx context.Context, messages []model.AIMessage, tradingContext map[string]interface{}) (*model.AIMessage, error) {
	return &model.AIMessage{
		ID: "mock-id",
		ConversationID: "mock-conv-id",
		Role: "assistant",
		Content: "mock response",
	}, nil
}

func (m *MockAIService) Chat(ctx context.Context, message string, conversationID string) (*model.AIMessage, error) {
	return &model.AIMessage{
		ID: "mock-id",
		ConversationID: conversationID,
		Role: "assistant",
		Content: "mock response",
	}, nil
}

func (m *MockAIService) GenerateInsight(ctx context.Context, insightType string, data map[string]interface{}) (*model.AIInsight, error) {
	return &model.AIInsight{}, nil
}

func (m *MockAIService) GenerateTradeRecommendation(ctx context.Context, data map[string]interface{}) (*model.AITradeRecommendation, error) {
	return &model.AITradeRecommendation{}, nil
}

func (m *MockAIService) GenerateEmbedding(ctx context.Context, text string) (*model.AIEmbedding, error) {
	return &model.AIEmbedding{}, nil
}

func (m *MockAIService) ExecuteFunction(ctx context.Context, functionCall model.AIFunctionCall) (*model.AIFunctionResponse, error) {
	return &model.AIFunctionResponse{}, nil
}
