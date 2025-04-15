package repo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// AIConversationEntity represents an AI conversation in the database
type AIConversationEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)"`
	UserID    string    `gorm:"index;type:varchar(50)"`
	Title     string    `gorm:"type:varchar(255)"`
	Tags      string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// AIMessageEntity represents an AI message in the database
type AIMessageEntity struct {
	ID             string `gorm:"primaryKey;type:varchar(50)"`
	ConversationID string `gorm:"index;type:varchar(50)"`
	Role           string `gorm:"type:varchar(20)"`
	Content        string `gorm:"type:text"`
	Timestamp      time.Time
	Metadata       []byte `gorm:"type:json"`
}

// GormAIConversationRepository implements port.ConversationMemoryRepository using GORM
type GormAIConversationRepository struct {
	BaseRepository
}

// NewGormAIConversationRepository creates a new GormAIConversationRepository
func NewGormAIConversationRepository(db *gorm.DB, logger *zerolog.Logger) *GormAIConversationRepository {
	return &GormAIConversationRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// SaveConversation saves a conversation
func (r *GormAIConversationRepository) SaveConversation(ctx context.Context, conversation *model.AIConversation) error {
	// If the conversation doesn't have an ID, generate one
	if conversation.ID == "" {
		conversation.ID = uuid.New().String()
	}

	// Update timestamps
	if conversation.CreatedAt.IsZero() {
		conversation.CreatedAt = time.Now()
	}
	conversation.UpdatedAt = time.Now()

	// Convert tags to JSON string
	tagsJSON, err := json.Marshal(conversation.Tags)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to marshal conversation tags")
		return err
	}

	// Create entity
	entity := &AIConversationEntity{
		ID:        conversation.ID,
		UserID:    conversation.UserID,
		Title:     conversation.Title,
		Tags:      string(tagsJSON),
		CreatedAt: conversation.CreatedAt,
		UpdatedAt: conversation.UpdatedAt,
	}

	// Save entity
	return r.Upsert(ctx, entity, []string{"id"}, []string{
		"user_id", "title", "tags", "updated_at",
	})
}

// GetConversation retrieves a conversation by ID
func (r *GormAIConversationRepository) GetConversation(ctx context.Context, id string) (*model.AIConversation, error) {
	var entity AIConversationEntity
	err := r.FindOne(ctx, &entity, "id = ?", id)
	if err != nil {
		return nil, err
	}

	if entity.ID == "" {
		return nil, nil // Not found
	}

	// Get messages for this conversation
	var messageEntities []AIMessageEntity
	err = r.GetDB(ctx).
		Where("conversation_id = ?", id).
		Order("timestamp ASC").
		Find(&messageEntities).Error
	if err != nil {
		r.logger.Error().Err(err).Str("conversation_id", id).Msg("Failed to get messages for conversation")
		return nil, err
	}

	// Convert to domain model
	conversation := r.toDomain(&entity)
	conversation.Messages = r.messagesToDomain(messageEntities)

	return conversation, nil
}

// ListConversations lists conversations for a user
func (r *GormAIConversationRepository) ListConversations(ctx context.Context, userID string, limit, offset int) ([]*model.AIConversation, error) {
	var entities []AIConversationEntity
	err := r.GetDB(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entities).Error
	if err != nil {
		r.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to list conversations")
		return nil, err
	}

	// Convert to domain models
	conversations := make([]*model.AIConversation, len(entities))
	for i, entity := range entities {
		conversations[i] = r.toDomain(&entity)
	}

	return conversations, nil
}

// SaveMessage saves a message to a conversation
func (r *GormAIConversationRepository) SaveMessage(ctx context.Context, message *model.AIMessage) error {
	// Check if the conversation exists
	var count int64
	err := r.GetDB(ctx).
		Model(&AIConversationEntity{}).
		Where("id = ?", message.ConversationID).
		Count(&count).Error
	if err != nil {
		r.logger.Error().Err(err).Str("conversation_id", message.ConversationID).Msg("Failed to check if conversation exists")
		return err
	}
	if count == 0 {
		return errors.New("conversation not found")
	}

	// If the message doesn't have an ID, generate one
	if message.ID == "" {
		message.ID = uuid.New().String()
	}

	// Set timestamp if not set
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}

	// Convert metadata to JSON
	var metadataJSON []byte
	if message.Metadata != nil {
		var err error
		metadataJSON, err = json.Marshal(message.Metadata)
		if err != nil {
			r.logger.Error().Err(err).Msg("Failed to marshal message metadata")
			return err
		}
	}

	// Create entity
	entity := &AIMessageEntity{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		Role:           message.Role,
		Content:        message.Content,
		Timestamp:      message.Timestamp,
		Metadata:       metadataJSON,
	}

	// Save entity
	return r.Create(ctx, entity)
}

// GetMessages retrieves messages for a conversation
func (r *GormAIConversationRepository) GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]*model.AIMessage, error) {
	// Check if the conversation exists
	var count int64
	err := r.GetDB(ctx).
		Model(&AIConversationEntity{}).
		Where("id = ?", conversationID).
		Count(&count).Error
	if err != nil {
		r.logger.Error().Err(err).Str("conversation_id", conversationID).Msg("Failed to check if conversation exists")
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("conversation not found")
	}

	// Get messages
	var entities []AIMessageEntity
	err = r.GetDB(ctx).
		Where("conversation_id = ?", conversationID).
		Order("timestamp ASC").
		Limit(limit).
		Offset(offset).
		Find(&entities).Error
	if err != nil {
		r.logger.Error().Err(err).Str("conversation_id", conversationID).Msg("Failed to get messages")
		return nil, err
	}

	// Convert to domain models
	messages := r.messagesToDomain(entities)

	// Convert to pointer slice
	result := make([]*model.AIMessage, len(messages))
	for i := range messages {
		result[i] = &messages[i]
	}

	return result, nil
}

// DeleteConversation deletes a conversation
func (r *GormAIConversationRepository) DeleteConversation(ctx context.Context, id string) error {
	// Use a transaction to delete the conversation and its messages
	return r.Transaction(ctx, func(tx *gorm.DB) error {
		// Delete messages first (due to foreign key constraint)
		if err := tx.Where("conversation_id = ?", id).Delete(&AIMessageEntity{}).Error; err != nil {
			return err
		}

		// Delete the conversation
		if err := tx.Delete(&AIConversationEntity{}, "id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}

// Helper methods for entity conversion

// toDomain converts a database entity to a domain model
func (r *GormAIConversationRepository) toDomain(entity *AIConversationEntity) *model.AIConversation {
	if entity == nil {
		return nil
	}

	// Parse tags
	var tags []string
	if entity.Tags != "" {
		if err := json.Unmarshal([]byte(entity.Tags), &tags); err != nil {
			r.logger.Error().Err(err).Msg("Failed to unmarshal conversation tags")
		}
	}

	return &model.AIConversation{
		ID:        entity.ID,
		UserID:    entity.UserID,
		Title:     entity.Title,
		Tags:      tags,
		Messages:  []model.AIMessage{}, // Will be populated separately
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// messageToDomain converts a message entity to a domain model
func (r *GormAIConversationRepository) messageToDomain(entity *AIMessageEntity) *model.AIMessage {
	if entity == nil {
		return nil
	}

	// Parse metadata
	var metadata map[string]interface{}
	if len(entity.Metadata) > 0 {
		if err := json.Unmarshal(entity.Metadata, &metadata); err != nil {
			r.logger.Error().Err(err).Msg("Failed to unmarshal message metadata")
		}
	}

	return &model.AIMessage{
		ID:             entity.ID,
		ConversationID: entity.ConversationID,
		Role:           entity.Role,
		Content:        entity.Content,
		Timestamp:      entity.Timestamp,
		Metadata:       metadata,
	}
}

// messagesToDomain converts message entities to domain models
func (r *GormAIConversationRepository) messagesToDomain(entities []AIMessageEntity) []model.AIMessage {
	messages := make([]model.AIMessage, len(entities))
	for i, entity := range entities {
		msg := r.messageToDomain(&entity)
		if msg != nil {
			messages[i] = *msg
		}
	}
	return messages
}

// Ensure GormAIConversationRepository implements port.ConversationMemoryRepository
var _ port.ConversationMemoryRepository = (*GormAIConversationRepository)(nil)
