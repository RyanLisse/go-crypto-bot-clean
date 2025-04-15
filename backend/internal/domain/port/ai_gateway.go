package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// AIService defines the interface for interacting with AI services
type AIService interface {
	// Chat sends a message to the AI and returns a response
	Chat(ctx context.Context, message string, conversationID string) (*model.AIMessage, error)
	
	// ChatWithHistory sends a message with conversation history to the AI
	ChatWithHistory(ctx context.Context, messages []model.AIMessage) (*model.AIMessage, error)
	
	// GenerateInsight generates an insight based on provided data
	GenerateInsight(ctx context.Context, insightType string, data map[string]interface{}) (*model.AIInsight, error)
	
	// GenerateTradeRecommendation generates a trade recommendation
	GenerateTradeRecommendation(ctx context.Context, data map[string]interface{}) (*model.AITradeRecommendation, error)
	
	// ExecuteFunction executes a function call from the AI
	ExecuteFunction(ctx context.Context, functionCall model.AIFunctionCall) (*model.AIFunctionResponse, error)
	
	// GenerateEmbedding generates a vector embedding for a text
	GenerateEmbedding(ctx context.Context, text string) (*model.AIEmbedding, error)
}

// ConversationMemoryRepository defines the interface for storing and retrieving conversations
type ConversationMemoryRepository interface {
	// SaveConversation saves a conversation
	SaveConversation(ctx context.Context, conversation *model.AIConversation) error
	
	// GetConversation retrieves a conversation by ID
	GetConversation(ctx context.Context, id string) (*model.AIConversation, error)
	
	// ListConversations lists conversations for a user
	ListConversations(ctx context.Context, userID string, limit, offset int) ([]*model.AIConversation, error)
	
	// SaveMessage saves a message to a conversation
	SaveMessage(ctx context.Context, message *model.AIMessage) error
	
	// GetMessages retrieves messages for a conversation
	GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]*model.AIMessage, error)
	
	// DeleteConversation deletes a conversation
	DeleteConversation(ctx context.Context, id string) error
}

// EmbeddingRepository defines the interface for storing and retrieving embeddings
type EmbeddingRepository interface {
	// SaveEmbedding saves an embedding
	SaveEmbedding(ctx context.Context, embedding *model.AIEmbedding) error
	
	// FindSimilar finds similar embeddings
	FindSimilar(ctx context.Context, vector []float64, limit int) ([]*model.AIEmbedding, error)
	
	// GetEmbedding retrieves an embedding by source ID and type
	GetEmbedding(ctx context.Context, sourceID, sourceType string) (*model.AIEmbedding, error)
	
	// DeleteEmbedding deletes an embedding
	DeleteEmbedding(ctx context.Context, id string) error
}
