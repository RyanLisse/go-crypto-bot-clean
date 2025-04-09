package similarity

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/ai/repository"
	"go-crypto-bot-clean/backend/internal/domain/ai/service/embeddings"
	"go-crypto-bot-clean/backend/internal/domain/ai/types"

	"go.uber.org/zap"
)

// For backward compatibility
type SimilarMessage = types.SimilarMessage

// SimilarityService manages vector similarity search
type SimilarityService struct {
	embeddingsRepo    repository.EmbeddingsRepository
	embeddingsService embeddings.EmbeddingsService
	logger            *zap.Logger
}

// NewSimilarityService creates a new similarity service
func NewSimilarityService(
	embeddingsRepo repository.EmbeddingsRepository,
	embeddingsService embeddings.EmbeddingsService,
	logger *zap.Logger,
) *SimilarityService {
	return &SimilarityService{
		embeddingsRepo:    embeddingsRepo,
		embeddingsService: embeddingsService,
		logger:            logger,
	}
}

// IndexMessage indexes a message for similarity search
func (s *SimilarityService) IndexMessage(
	ctx context.Context,
	conversationID, messageID, content string,
	metadata map[string]interface{},
) error {
	// Generate embedding for the message
	embedding, err := s.embeddingsService.GenerateEmbedding(ctx, content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Add timestamp to metadata
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["indexed_at"] = time.Now().UTC().Format(time.RFC3339)

	// Store embedding in repository
	err = s.embeddingsRepo.StoreEmbedding(ctx, conversationID, messageID, content, embedding, metadata)
	if err != nil {
		return fmt.Errorf("failed to store embedding: %w", err)
	}

	s.logger.Info("Indexed message for similarity search",
		zap.String("conversation_id", conversationID),
		zap.String("message_id", messageID),
		zap.Int("embedding_dimensions", len(embedding)),
	)

	return nil
}

// FindSimilarMessages finds messages similar to the given query
func (s *SimilarityService) FindSimilarMessages(
	ctx context.Context,
	query string,
	limit int,
) ([]types.SimilarMessage, error) {
	// Generate embedding for the query
	embedding, err := s.embeddingsService.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Find similar embeddings
	similarEmbeddings, err := s.embeddingsRepo.FindSimilarEmbeddings(ctx, embedding, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find similar embeddings: %w", err)
	}

	// Convert to SimilarMessage
	similarMessages := make([]types.SimilarMessage, len(similarEmbeddings))
	for i, embedding := range similarEmbeddings {
		similarMessages[i] = types.SimilarMessage{
			ConversationID: embedding.ConversationID,
			MessageID:      embedding.MessageID,
			Content:        embedding.Content,
			Similarity:     embedding.Similarity,
			Metadata:       embedding.Metadata,
		}
	}

	s.logger.Info("Found similar messages",
		zap.Int("count", len(similarMessages)),
		zap.String("query", query),
	)

	return similarMessages, nil
}

// DeleteConversationMessages deletes all indexed messages for a conversation
func (s *SimilarityService) DeleteConversationMessages(
	ctx context.Context,
	conversationID string,
) error {
	err := s.embeddingsRepo.DeleteEmbeddingsByConversation(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("failed to delete conversation messages: %w", err)
	}

	s.logger.Info("Deleted indexed messages for conversation",
		zap.String("conversation_id", conversationID),
	)

	return nil
}
