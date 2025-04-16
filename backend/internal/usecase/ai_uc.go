package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// AIUsecase handles AI-related operations
type AIUsecase struct {
	aiService              port.AIService
	conversationMemoryRepo port.ConversationMemoryRepository
	embeddingRepo          port.EmbeddingRepository
	logger                 zerolog.Logger
}

// NewAIUsecase creates a new AIUsecase
func NewAIUsecase(
	aiService port.AIService,
	conversationMemoryRepo port.ConversationMemoryRepository,
	embeddingRepo port.EmbeddingRepository,
	logger zerolog.Logger,
) *AIUsecase {
	return &AIUsecase{
		aiService:              aiService,
		conversationMemoryRepo: conversationMemoryRepo,
		embeddingRepo:          embeddingRepo,
		logger:                 logger.With().Str("component", "ai_usecase").Logger(),
	}
}

// Chat sends a message to the AI and returns a response
func (uc *AIUsecase) Chat(ctx context.Context, userID, message, conversationID string, tradingContext map[string]interface{}) (*model.AIMessage, error) {
	// Create a new conversation if conversationID is empty
	if conversationID == "" {
		conversation := &model.AIConversation{
			ID:        uuid.New().String(),
			UserID:    userID,
			Title:     generateTitle(message),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := uc.conversationMemoryRepo.SaveConversation(ctx, conversation); err != nil {
			uc.logger.Error().Err(err).Msg("Failed to save new conversation")
			return nil, err
		}

		conversationID = conversation.ID
	}

	// Create user message
	userMessage := &model.AIMessage{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		Role:           "user",
		Content:        message,
		Timestamp:      time.Now(),
	}

	// Save user message
	if err := uc.conversationMemoryRepo.SaveMessage(ctx, userMessage); err != nil {
		uc.logger.Error().Err(err).Msg("Failed to save user message")
		return nil, err
	}

	// Get conversation history
	messages, err := uc.conversationMemoryRepo.GetMessages(ctx, conversationID, 10, 0)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Failed to get conversation history")
		return nil, err
	}

	// Convert messages to AIMessage slice
	aiMessages := make([]model.AIMessage, len(messages))
	for i, msg := range messages {
		aiMessages[i] = *msg
	}

	// Send message to AI with history and trading context
	response, err := uc.aiService.ChatWithHistory(ctx, aiMessages, tradingContext)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Failed to get AI response")
		return nil, err
	}

	// Save AI response
	if err := uc.conversationMemoryRepo.SaveMessage(ctx, response); err != nil {
		uc.logger.Error().Err(err).Msg("Failed to save AI response")
		// Don't return error here, we still want to return the response to the user
	}

	// Update conversation title if it's the first message
	if len(aiMessages) <= 2 {
		conversation, err := uc.conversationMemoryRepo.GetConversation(ctx, conversationID)
		if err == nil {
			conversation.Title = generateTitle(message)
			_ = uc.conversationMemoryRepo.SaveConversation(ctx, conversation)
		}
	}

	return response, nil
}

// GetConversation retrieves a conversation by ID
func (uc *AIUsecase) GetConversation(ctx context.Context, userID, conversationID string) (*model.AIConversation, error) {
	conversation, err := uc.conversationMemoryRepo.GetConversation(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	// Check if the conversation belongs to the user
	if conversation.UserID != userID {
		return nil, errors.New("unauthorized access to conversation")
	}

	return conversation, nil
}

// ListConversations lists conversations for a user
func (uc *AIUsecase) ListConversations(ctx context.Context, userID string, limit, offset int) ([]*model.AIConversation, error) {
	return uc.conversationMemoryRepo.ListConversations(ctx, userID, limit, offset)
}

// DeleteConversation deletes a conversation
func (uc *AIUsecase) DeleteConversation(ctx context.Context, userID, conversationID string) error {
	// Check if the conversation belongs to the user
	conversation, err := uc.conversationMemoryRepo.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}

	if conversation.UserID != userID {
		return errors.New("unauthorized access to conversation")
	}

	return uc.conversationMemoryRepo.DeleteConversation(ctx, conversationID)
}

// GetMessages retrieves messages for a conversation
func (uc *AIUsecase) GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]*model.AIMessage, error) {
	return uc.conversationMemoryRepo.GetMessages(ctx, conversationID, limit, offset)
}

// GenerateInsight generates an insight based on provided data
func (uc *AIUsecase) GenerateInsight(ctx context.Context, userID, insightType string, data map[string]interface{}) (*model.AIInsight, error) {
	// Add user ID to data
	data["user_id"] = userID

	// Generate insight
	insight, err := uc.aiService.GenerateInsight(ctx, insightType, data)
	if err != nil {
		uc.logger.Error().Err(err).Str("insightType", insightType).Msg("Failed to generate insight")
		return nil, err
	}

	return insight, nil
}

// GenerateTradeRecommendation generates a trade recommendation
func (uc *AIUsecase) GenerateTradeRecommendation(ctx context.Context, userID string, data map[string]interface{}) (*model.AITradeRecommendation, error) {
	// Add user ID to data
	data["user_id"] = userID

	// Generate trade recommendation
	recommendation, err := uc.aiService.GenerateTradeRecommendation(ctx, data)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Failed to generate trade recommendation")
		return nil, err
	}

	return recommendation, nil
}

// ExecuteFunction executes a function call from the AI
func (uc *AIUsecase) ExecuteFunction(ctx context.Context, functionCall model.AIFunctionCall) (*model.AIFunctionResponse, error) {
	return uc.aiService.ExecuteFunction(ctx, functionCall)
}

// Helper functions

// generateTitle generates a title for a conversation based on the first message
func generateTitle(message string) string {
	// Simple implementation: use the first 30 characters of the message
	title := message
	if len(title) > 30 {
		title = title[:30] + "..."
	}
	return title
}
