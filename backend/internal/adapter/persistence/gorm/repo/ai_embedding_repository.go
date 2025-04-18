package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// AIEmbeddingEntity represents an AI embedding in the database
type AIEmbeddingEntity struct {
	ID         string    `gorm:"primaryKey;type:varchar(50)"`
	SourceID   string    `gorm:"index;type:varchar(50)"`
	SourceType string    `gorm:"index;type:varchar(20)"`
	Vector     []byte    `gorm:"type:blob"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

// GormAIEmbeddingRepository implements port.EmbeddingRepository using GORM
type GormAIEmbeddingRepository struct {
	BaseRepository
}

// NewGormAIEmbeddingRepository creates a new GormAIEmbeddingRepository
func NewGormAIEmbeddingRepository(db *gorm.DB, logger *zerolog.Logger) *GormAIEmbeddingRepository {
	return &GormAIEmbeddingRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// SaveEmbedding saves an embedding
func (r *GormAIEmbeddingRepository) SaveEmbedding(ctx context.Context, embedding *model.AIEmbedding) error {
	// If the embedding doesn't have an ID, generate one
	if embedding.ID == "" {
		embedding.ID = uuid.New().String()
	}

	// Convert vector to JSON
	vectorJSON, err := json.Marshal(embedding.Vector)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to marshal embedding vector")
		return err
	}

	// Create entity
	entity := &AIEmbeddingEntity{
		ID:         embedding.ID,
		SourceID:   embedding.SourceID,
		SourceType: embedding.SourceType,
		Vector:     vectorJSON,
		CreatedAt:  embedding.CreatedAt,
	}

	// Save entity
	return r.Upsert(ctx, entity, []string{"id"}, []string{
		"source_id", "source_type", "vector",
	})
}

// SearchEmbeddings searches for embeddings similar to a query vector
// This is a placeholder implementation - in a real system you'd use a vector database like Pinecone, Weaviate, etc.
func (r *GormAIEmbeddingRepository) SearchEmbeddings(ctx context.Context, queryVector []float64, limit int) ([]*model.AIEmbedding, error) {
	// In a real implementation, this would use cosine similarity search or other vector similarity
	// For this simple implementation, just return the most recent embeddings

	var entities []AIEmbeddingEntity
	if err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&entities).Error; err != nil {
		return nil, err
	}

	// Convert to domain models
	embeddings := make([]*model.AIEmbedding, len(entities))
	for i, entity := range entities {
		embeddings[i] = r.toDomain(&entity)
	}

	return embeddings, nil
}

// GetEmbedding retrieves an embedding by ID
func (r *GormAIEmbeddingRepository) GetEmbedding(ctx context.Context, embeddingID string) (*model.AIEmbedding, error) {
	var entity AIEmbeddingEntity
	err := r.FindOne(ctx, &entity, "id = ?", embeddingID)
	if err != nil {
		return nil, err
	}

	if entity.ID == "" {
		return nil, nil // Not found
	}

	return r.toDomain(&entity), nil
}

// GetEmbeddingBySource retrieves an embedding by source ID and type
func (r *GormAIEmbeddingRepository) GetEmbeddingBySource(ctx context.Context, sourceID, sourceType string) (*model.AIEmbedding, error) {
	var entity AIEmbeddingEntity
	err := r.FindOne(ctx, &entity, "source_id = ? AND source_type = ?", sourceID, sourceType)
	if err != nil {
		return nil, err
	}

	if entity.ID == "" {
		return nil, nil // Not found
	}

	return r.toDomain(&entity), nil
}

// DeleteEmbedding deletes an embedding
func (r *GormAIEmbeddingRepository) DeleteEmbedding(ctx context.Context, id string) error {
	return r.DeleteByID(ctx, &AIEmbeddingEntity{}, id)
}

// FindSimilar is an alias for SearchEmbeddings for backward compatibility
func (r *GormAIEmbeddingRepository) FindSimilar(ctx context.Context, vector []float64, limit int) ([]*model.AIEmbedding, error) {
	return r.SearchEmbeddings(ctx, vector, limit)
}

// Helper methods for entity conversion

// toDomain converts a database entity to a domain model
func (r *GormAIEmbeddingRepository) toDomain(entity *AIEmbeddingEntity) *model.AIEmbedding {
	if entity == nil {
		return nil
	}

	// Parse vector
	var vector []float64
	if len(entity.Vector) > 0 {
		if err := json.Unmarshal(entity.Vector, &vector); err != nil {
			r.logger.Error().Err(err).Msg("Failed to unmarshal embedding vector")
			return nil
		}
	}

	return &model.AIEmbedding{
		ID:         entity.ID,
		SourceID:   entity.SourceID,
		SourceType: entity.SourceType,
		Vector:     vector,
		CreatedAt:  entity.CreatedAt,
	}
}

// Ensure GormAIEmbeddingRepository implements port.EmbeddingRepository
var _ port.EmbeddingRepository = (*GormAIEmbeddingRepository)(nil)
