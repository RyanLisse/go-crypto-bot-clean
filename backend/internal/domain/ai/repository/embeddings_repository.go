package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"go-crypto-bot-clean/backend/internal/domain/ai/repository/schema"

	"gorm.io/gorm"
)

// EmbeddingsRepository defines the interface for managing embeddings
type EmbeddingsRepository interface {
	// StoreEmbedding stores a new embedding
	StoreEmbedding(ctx context.Context, conversationID, messageID, content string, embedding []float32, metadata map[string]interface{}) error

	// FindSimilarEmbeddings finds similar embeddings to the given embedding
	FindSimilarEmbeddings(ctx context.Context, embedding []float32, limit int) ([]SimilarEmbedding, error)

	// GetEmbeddingsByConversation gets all embeddings for a conversation
	GetEmbeddingsByConversation(ctx context.Context, conversationID string) ([]schema.ConversationEmbedding, error)

	// DeleteEmbeddingsByConversation deletes all embeddings for a conversation
	DeleteEmbeddingsByConversation(ctx context.Context, conversationID string) error
}

// SimilarEmbedding represents a similar embedding with its similarity score
type SimilarEmbedding struct {
	ConversationID string
	MessageID      string
	Content        string
	Similarity     float64
	Metadata       map[string]interface{}
}

// embeddingsRepositoryImpl implements the EmbeddingsRepository interface
type embeddingsRepositoryImpl struct {
	db *gorm.DB
}

// NewEmbeddingsRepository creates a new embeddings repository
func NewEmbeddingsRepository(db *gorm.DB) (EmbeddingsRepository, error) {
	// Create the embeddings table if it doesn't exist
	if err := schema.CreateEmbeddingsTable(db); err != nil {
		return nil, err
	}

	return &embeddingsRepositoryImpl{
		db: db,
	}, nil
}

// StoreEmbedding stores a new embedding
func (r *embeddingsRepositoryImpl) StoreEmbedding(
	ctx context.Context,
	conversationID, messageID, content string,
	embedding []float32,
	metadata map[string]interface{},
) error {
	// Convert embedding to binary data
	embeddingBytes := make([]byte, len(embedding)*4)
	for i, v := range embedding {
		// Convert float32 to bytes (little-endian)
		bits := math.Float32bits(v)
		embeddingBytes[i*4] = byte(bits)
		embeddingBytes[i*4+1] = byte(bits >> 8)
		embeddingBytes[i*4+2] = byte(bits >> 16)
		embeddingBytes[i*4+3] = byte(bits >> 24)
	}

	// Convert metadata to JSON
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Create embedding record
	embeddingRecord := schema.ConversationEmbedding{
		ConversationID:  conversationID,
		MessageID:       messageID,
		Content:         content,
		EmbeddingVector: embeddingBytes,
		Dimensions:      len(embedding),
		Metadata:        string(metadataJSON),
	}

	// Store embedding in database
	if err := r.db.WithContext(ctx).Create(&embeddingRecord).Error; err != nil {
		return fmt.Errorf("failed to store embedding: %w", err)
	}

	return nil
}

// FindSimilarEmbeddings finds similar embeddings to the given embedding
func (r *embeddingsRepositoryImpl) FindSimilarEmbeddings(
	ctx context.Context,
	embedding []float32,
	limit int,
) ([]SimilarEmbedding, error) {
	// Convert embedding to binary data
	embeddingBytes := make([]byte, len(embedding)*4)
	for i, v := range embedding {
		// Convert float32 to bytes (little-endian)
		bits := math.Float32bits(v)
		embeddingBytes[i*4] = byte(bits)
		embeddingBytes[i*4+1] = byte(bits >> 8)
		embeddingBytes[i*4+2] = byte(bits >> 16)
		embeddingBytes[i*4+3] = byte(bits >> 24)
	}

	// Use vector_top_k to find similar embeddings
	// This is specific to Turso/libSQL and won't work with other databases
	var results []struct {
		ConversationID string
		MessageID      string
		Content        string
		Metadata       string
		Distance       float64
	}

	// Execute the query
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			e.conversation_id,
			e.message_id,
			e.content,
			e.metadata,
			vector_distance_cos(e.embedding_vector, ?) as distance
		FROM
			vector_top_k('idx_conversation_embeddings_vector', ?, ?) as v
		JOIN
			conversation_embeddings e ON e.id = v.id
		ORDER BY
			distance ASC
	`, embeddingBytes, embeddingBytes, limit).Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find similar embeddings: %w", err)
	}

	// Convert results to SimilarEmbedding
	similarEmbeddings := make([]SimilarEmbedding, 0, len(results))
	for _, result := range results {
		// Parse metadata
		var metadata map[string]interface{}
		if result.Metadata != "" {
			if err := json.Unmarshal([]byte(result.Metadata), &metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		// Convert distance to similarity (1 - distance)
		similarity := 1.0 - result.Distance

		similarEmbeddings = append(similarEmbeddings, SimilarEmbedding{
			ConversationID: result.ConversationID,
			MessageID:      result.MessageID,
			Content:        result.Content,
			Similarity:     similarity,
			Metadata:       metadata,
		})
	}

	return similarEmbeddings, nil
}

// GetEmbeddingsByConversation gets all embeddings for a conversation
func (r *embeddingsRepositoryImpl) GetEmbeddingsByConversation(
	ctx context.Context,
	conversationID string,
) ([]schema.ConversationEmbedding, error) {
	var embeddings []schema.ConversationEmbedding

	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Find(&embeddings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get embeddings for conversation: %w", err)
	}

	return embeddings, nil
}

// DeleteEmbeddingsByConversation deletes all embeddings for a conversation
func (r *embeddingsRepositoryImpl) DeleteEmbeddingsByConversation(
	ctx context.Context,
	conversationID string,
) error {
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Delete(&schema.ConversationEmbedding{}).Error

	if err != nil {
		return fmt.Errorf("failed to delete embeddings for conversation: %w", err)
	}

	return nil
}
