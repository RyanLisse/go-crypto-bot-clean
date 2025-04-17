package mocks

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
)

// MockAIUsecase is a mock implementation of the AIUsecase interface for testing
type MockAIUsecase struct{}

// Message represents a chat message
type Message struct {
	Role      string                 `json:"role"`
	Content   string                 `json:"content"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Chat implements the AIUsecase interface
func (m *MockAIUsecase) Chat(ctx context.Context, userID, message, sessionID string, tradingContext map[string]interface{}) (*usecase.AIMessage, error) {
	return &usecase.AIMessage{
		Content: "Hello from AI!",
		Metadata: map[string]interface{}{
			"function_calls": map[string]interface{}{},
		},
	}, nil
}

// GetConversationHistory implements the AIUsecase interface
func (m *MockAIUsecase) GetConversationHistory(ctx context.Context, userID, sessionID string, page, pageSize int) ([]Message, error) {
	return []Message{}, nil
}

// ListConversations implements the AIUsecase interface
func (m *MockAIUsecase) ListConversations(ctx context.Context, userID string, limit, offset int) ([]interface{}, error) {
	return []interface{}{}, nil
}

// GetConversation implements the AIUsecase interface
func (m *MockAIUsecase) GetConversation(ctx context.Context, userID, conversationID string) (interface{}, error) {
	return map[string]interface{}{
		"id":      conversationID,
		"user_id": userID,
		"title":   "Test Conversation",
	}, nil
}

// GetMessages implements the AIUsecase interface
func (m *MockAIUsecase) GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]interface{}, error) {
	return []interface{}{}, nil
}

// DeleteConversation implements the AIUsecase interface
func (m *MockAIUsecase) DeleteConversation(ctx context.Context, userID, conversationID string) error {
	return nil
}
