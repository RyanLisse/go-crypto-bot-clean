package service

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/ai/service/function"
	"go-crypto-bot-clean/backend/internal/domain/ai/service/templates"
	"go-crypto-bot-clean/backend/internal/domain/ai/types"
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

	// GenerateResponseWithTemplate generates an AI response using a specific template
	GenerateResponseWithTemplate(ctx context.Context, userID int, templateName string, templateData templates.TemplateData) (string, error)

	// ExecuteFunction allows the AI to call predefined functions
	ExecuteFunction(ctx context.Context, userID int, functionName string, parameters map[string]interface{}) (interface{}, error)

	// GetAvailableFunctions returns all available functions
	GetAvailableFunctions(ctx context.Context) []function.FunctionDefinition

	// GetAvailableTemplates returns all available templates
	GetAvailableTemplates(ctx context.Context) []string

	// StoreConversation saves a conversation to the database
	StoreConversation(ctx context.Context, userID int, sessionID string, messages []Message) error

	// RetrieveConversation gets a conversation from the database
	RetrieveConversation(ctx context.Context, userID int, sessionID string) (*ConversationMemory, error)

	// ListUserSessions lists all sessions for a user
	ListUserSessions(ctx context.Context, userID int, limit int) ([]string, error)

	// DeleteSession deletes a session
	DeleteSession(ctx context.Context, userID int, sessionID string) error

	// Risk Management Methods

	// ApplyRiskGuardrails applies risk guardrails to a trade recommendation
	ApplyRiskGuardrails(ctx context.Context, userID int, recommendation *TradeRecommendation) (*GuardrailsResult, error)

	// CreateTradeConfirmation creates a trade confirmation
	CreateTradeConfirmation(ctx context.Context, userID int, trade *TradeRequest, recommendation *TradeRecommendation) (*TradeConfirmation, error)

	// ConfirmTrade confirms a trade
	ConfirmTrade(ctx context.Context, confirmationID string, approve bool) (*TradeConfirmation, error)

	// ListPendingTradeConfirmations lists all pending trade confirmations for a user
	ListPendingTradeConfirmations(ctx context.Context, userID int) ([]*TradeConfirmation, error)

	// FindSimilarMessages finds messages similar to the given query
	FindSimilarMessages(ctx context.Context, query string, limit int) ([]types.SimilarMessage, error)

	// IndexMessage indexes a message for similarity search
	IndexMessage(ctx context.Context, conversationID, messageID, content string, metadata map[string]interface{}) error
}
