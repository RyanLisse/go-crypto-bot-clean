package schema

import (
	"fmt"

	"gorm.io/gorm"
)

// ConversationEmbedding represents a vector embedding for a conversation message
type ConversationEmbedding struct {
	ID              uint    `gorm:"primaryKey"`
	ConversationID  string  `gorm:"index;not null"`
	MessageID       string  `gorm:"index;not null"`
	Content         string  `gorm:"type:text;not null"`
	EmbeddingVector []byte  `gorm:"type:blob;not null"` // Store embedding as binary data
	Dimensions      int     `gorm:"not null"`           // Number of dimensions in the embedding
	Metadata        string  `gorm:"type:text"`          // JSON metadata
}

// TableName returns the table name for the conversation embedding
func (ConversationEmbedding) TableName() string {
	return "conversation_embeddings"
}

// CreateEmbeddingsTable creates the conversation embeddings table
func CreateEmbeddingsTable(db *gorm.DB) error {
	// Create the table
	if err := db.AutoMigrate(&ConversationEmbedding{}); err != nil {
		return fmt.Errorf("failed to create conversation embeddings table: %w", err)
	}

	// Check if the vector index exists
	var count int64
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='idx_conversation_embeddings_vector'").Count(&count)
	
	// Create the vector index if it doesn't exist
	if count == 0 {
		// Create a vector index for similarity search
		// Note: This is specific to Turso/libSQL and won't work with other databases
		err := db.Exec(`
			CREATE INDEX IF NOT EXISTS idx_conversation_embeddings_vector 
			ON conversation_embeddings (
				libsql_vector_idx(
					embedding_vector,
					'type=diskann',
					'metric=cosine'
				)
			)
		`).Error
		
		if err != nil {
			return fmt.Errorf("failed to create vector index: %w", err)
		}
	}

	return nil
}
