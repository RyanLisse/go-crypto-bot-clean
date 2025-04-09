package gemini

import (
	"context"
	"fmt"

	"go-crypto-bot-clean/backend/internal/domain/ai/similarity"
	"go-crypto-bot-clean/backend/internal/domain/ai/types"

	"go.uber.org/zap"
)

// FindSimilarMessages finds messages similar to the given query
func (s *GeminiAIService) FindSimilarMessages(
	ctx context.Context,
	query string,
	limit int,
) ([]types.SimilarMessage, error) {
	// Check if similarity service is available
	if s.SimilarityService == nil {
		return nil, fmt.Errorf("similarity service not available")
	}

	// Log the request
	s.Logger.Info("Finding similar messages",
		zap.String("query", query),
		zap.Int("limit", limit),
	)

	// Find similar messages
	similarMessages, err := s.SimilarityService.FindSimilarMessages(ctx, query, limit)
	if err != nil {
		s.Logger.Error("Failed to find similar messages",
			zap.Error(err),
			zap.String("query", query),
		)
		return nil, fmt.Errorf("failed to find similar messages: %w", err)
	}

	return similarMessages, nil
}

// IndexMessage indexes a message for similarity search
func (s *GeminiAIService) IndexMessage(
	ctx context.Context,
	conversationID, messageID, content string,
	metadata map[string]interface{},
) error {
	// Check if similarity service is available
	if s.SimilarityService == nil {
		return fmt.Errorf("similarity service not available")
	}

	// Log the request
	s.Logger.Info("Indexing message",
		zap.String("conversation_id", conversationID),
		zap.String("message_id", messageID),
		zap.Int("content_length", len(content)),
	)

	// Index message
	err := s.SimilarityService.IndexMessage(ctx, conversationID, messageID, content, metadata)
	if err != nil {
		s.Logger.Error("Failed to index message",
			zap.Error(err),
			zap.String("conversation_id", conversationID),
			zap.String("message_id", messageID),
		)
		return fmt.Errorf("failed to index message: %w", err)
	}

	return nil
}

// SetSimilarityService sets the similarity service
func (s *GeminiAIService) SetSimilarityService(similaritySvc *similarity.Service) {
	s.SimilarityService = similaritySvc
}
