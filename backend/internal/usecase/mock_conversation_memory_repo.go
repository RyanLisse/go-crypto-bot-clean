package usecase

import (
	"context"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

type MockConversationMemoryRepository struct{}

func (m *MockConversationMemoryRepository) SaveConversation(ctx context.Context, c *model.AIConversation) error {
	return nil
}
func (m *MockConversationMemoryRepository) SaveMessage(ctx context.Context, msg *model.AIMessage) error {
	return nil
}
func (m *MockConversationMemoryRepository) GetConversation(ctx context.Context, id string) (*model.AIConversation, error) {
	return &model.AIConversation{ID: id, UserID: "test-user"}, nil
}
func (m *MockConversationMemoryRepository) ListConversations(ctx context.Context, userID string, limit, offset int) ([]*model.AIConversation, error) {
	return []*model.AIConversation{{ID: "mock-conv", UserID: userID}}, nil
}
func (m *MockConversationMemoryRepository) SaveConversations(ctx context.Context, cs []*model.AIConversation) error {
	return nil
}
func (m *MockConversationMemoryRepository) SaveMessages(ctx context.Context, msgs []*model.AIMessage) error {
	return nil
}
func (m *MockConversationMemoryRepository) GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]*model.AIMessage, error) {
	return []*model.AIMessage{}, nil
}
func (m *MockConversationMemoryRepository) DeleteConversation(ctx context.Context, id string) error {
	return nil
}

func (m *MockConversationMemoryRepository) DeleteMessage(ctx context.Context, id string) error {
	return nil
}

func (m *MockConversationMemoryRepository) GetMessage(ctx context.Context, id string) (*model.AIMessage, error) {
	return &model.AIMessage{ID: id, ConversationID: "mock-conv", Role: "assistant", Content: "mock message"}, nil
}
