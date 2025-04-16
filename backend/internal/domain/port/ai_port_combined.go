package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// AI ports consolidated from ai_port.go and ai_gateway.go
// This file should be used instead of those two files

// AIService defines the interface for AI services
type AIService interface {
	// Chat sends a message to the AI and returns a response
	Chat(ctx context.Context, message string, conversationID string) (*model.AIMessage, error)

	// ChatWithHistory sends a message with conversation history and trading context to the AI
	ChatWithHistory(ctx context.Context, messages []model.AIMessage, tradingContext map[string]interface{}) (*model.AIMessage, error)

	// GenerateInsight generates an insight based on provided data
	GenerateInsight(ctx context.Context, insightType string, data map[string]interface{}) (*model.AIInsight, error)

	// GenerateTradeRecommendation generates a trade recommendation
	GenerateTradeRecommendation(ctx context.Context, data map[string]interface{}) (*model.AITradeRecommendation, error)

	// ExecuteFunction executes a function call from the AI
	ExecuteFunction(ctx context.Context, functionCall model.AIFunctionCall) (*model.AIFunctionResponse, error)

	// GenerateEmbedding generates a vector embedding for a text
	GenerateEmbedding(ctx context.Context, text string) (*model.AIEmbedding, error)
}

// ConversationMemoryRepository defines the interface for conversation memory repositories
type ConversationMemoryRepository interface {
	// SaveConversation saves a conversation
	SaveConversation(ctx context.Context, conversation *model.AIConversation) error

	// GetConversation gets a conversation by ID
	GetConversation(ctx context.Context, conversationID string) (*model.AIConversation, error)

	// ListConversations lists conversations for a user
	ListConversations(ctx context.Context, userID string, limit, offset int) ([]*model.AIConversation, error)

	// DeleteConversation deletes a conversation
	DeleteConversation(ctx context.Context, conversationID string) error

	// SaveMessage saves a message
	SaveMessage(ctx context.Context, message *model.AIMessage) error

	// GetMessages gets messages for a conversation
	GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]*model.AIMessage, error)

	// GetMessage gets a message by ID
	GetMessage(ctx context.Context, messageID string) (*model.AIMessage, error)

	// DeleteMessage deletes a message
	DeleteMessage(ctx context.Context, messageID string) error
}

// EmbeddingRepository defines the interface for embedding repositories
type EmbeddingRepository interface {
	// SaveEmbedding saves an embedding
	SaveEmbedding(ctx context.Context, embedding *model.AIEmbedding) error

	// GetEmbedding gets an embedding by ID
	GetEmbedding(ctx context.Context, embeddingID string) (*model.AIEmbedding, error)

	// SearchEmbeddings searches for embeddings similar to a query embedding
	SearchEmbeddings(ctx context.Context, queryVector []float64, limit int) ([]*model.AIEmbedding, error)

	// FindSimilar finds similar embeddings (alias for SearchEmbeddings)
	FindSimilar(ctx context.Context, vector []float64, limit int) ([]*model.AIEmbedding, error)

	// DeleteEmbedding deletes an embedding
	DeleteEmbedding(ctx context.Context, embeddingID string) error
}
