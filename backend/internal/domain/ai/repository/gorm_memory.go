package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/ai/service"
	"gorm.io/gorm"
)

// ConversationMemoryModel is the GORM model for conversation memories
type ConversationMemoryModel struct {
	UserID       int       `gorm:"primaryKey;column:user_id"`
	SessionID    string    `gorm:"primaryKey;column:session_id"`
	MessagesJSON string    `gorm:"column:messages_json"`
	Summary      string    `gorm:"column:summary"`
	LastAccessed time.Time `gorm:"column:last_accessed;index"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (ConversationMemoryModel) TableName() string {
	return "conversation_memories"
}

// GormConversationMemoryRepository implements ConversationMemoryRepository using GORM
type GormConversationMemoryRepository struct {
	db *gorm.DB
}

// NewGormConversationMemoryRepository creates a new GormConversationMemoryRepository
func NewGormConversationMemoryRepository(db *gorm.DB) (*GormConversationMemoryRepository, error) {
	// Auto-migrate the schema
	err := db.AutoMigrate(&ConversationMemoryModel{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate conversation_memories table: %w", err)
	}

	return &GormConversationMemoryRepository{db: db}, nil
}

// StoreConversation saves a conversation to the database
func (r *GormConversationMemoryRepository) StoreConversation(
	ctx context.Context,
	userID int,
	sessionID string,
	messages []service.Message,
) error {
	// Convert messages to JSON
	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to marshal messages: %w", err)
	}

	// Create model
	model := ConversationMemoryModel{
		UserID:       userID,
		SessionID:    sessionID,
		MessagesJSON: string(messagesJSON),
		LastAccessed: time.Now().UTC(),
	}

	// Insert or update conversation
	result := r.db.WithContext(ctx).Save(&model)
	if result.Error != nil {
		return fmt.Errorf("failed to store conversation: %w", result.Error)
	}

	return nil
}

// RetrieveConversation gets a conversation from the database
func (r *GormConversationMemoryRepository) RetrieveConversation(
	ctx context.Context,
	userID int,
	sessionID string,
) (*service.ConversationMemory, error) {
	var model ConversationMemoryModel

	result := r.db.WithContext(ctx).Where("user_id = ? AND session_id = ?", userID, sessionID).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // No conversation found
		}
		return nil, fmt.Errorf("failed to retrieve conversation: %w", result.Error)
	}

	// Parse messages JSON
	var messages []service.Message
	err := json.Unmarshal([]byte(model.MessagesJSON), &messages)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal messages: %w", err)
	}

	// Update last accessed time
	model.LastAccessed = time.Now().UTC()
	result = r.db.WithContext(ctx).Save(&model)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update last accessed time: %w", result.Error)
	}

	return &service.ConversationMemory{
		UserID:       userID,
		SessionID:    sessionID,
		Messages:     messages,
		Summary:      model.Summary,
		LastAccessed: model.LastAccessed,
	}, nil
}

// ListUserSessions lists all sessions for a user
func (r *GormConversationMemoryRepository) ListUserSessions(
	ctx context.Context,
	userID int,
	limit int,
) ([]string, error) {
	var models []ConversationMemoryModel

	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("last_accessed DESC").
		Limit(limit).
		Find(&models)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list user sessions: %w", result.Error)
	}

	sessions := make([]string, len(models))
	for i, model := range models {
		sessions[i] = model.SessionID
	}

	return sessions, nil
}

// DeleteSession deletes a session
func (r *GormConversationMemoryRepository) DeleteSession(
	ctx context.Context,
	userID int,
	sessionID string,
) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		Delete(&ConversationMemoryModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete session: %w", result.Error)
	}

	return nil
}
