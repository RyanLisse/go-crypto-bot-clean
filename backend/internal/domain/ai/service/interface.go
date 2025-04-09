package service

import (
	"context"
	"time"
)

// Message represents a single message in a conversation
type Message struct {
	Role      string                 `json:"role"` // "user" or "assistant"
	Content   string                 `json:"content"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"` // For tracking context, actions, etc.
}

// ConversationMemory stores conversation history
type ConversationMemory struct {
	UserID       int       `json:"user_id"`
	SessionID    string    `json:"session_id"`
	Messages     []Message `json:"messages"`
	Summary      string    `json:"summary"`
	LastAccessed time.Time `json:"last_accessed"`
}

// AIService defines the interface for AI interactions
type AIService interface {
	// GenerateResponse generates an AI response based on user message and context
	GenerateResponse(ctx context.Context, userID int, message string) (string, error)
	
	// ExecuteFunction allows the AI to call predefined functions
	ExecuteFunction(ctx context.Context, userID int, functionName string, parameters map[string]interface{}) (interface{}, error)
	
	// StoreConversation saves a conversation to the database
	StoreConversation(ctx context.Context, userID int, sessionID string, messages []Message) error
	
	// RetrieveConversation gets a conversation from the database
	RetrieveConversation(ctx context.Context, userID int, sessionID string) (*ConversationMemory, error)
}
