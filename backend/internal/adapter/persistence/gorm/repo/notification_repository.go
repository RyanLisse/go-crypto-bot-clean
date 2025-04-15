package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// NotificationPreferenceEntity represents notification preferences in the database
type NotificationPreferenceEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)"`
	UserID    string    `gorm:"uniqueIndex;type:varchar(50)"`
	Settings  []byte    `gorm:"type:json"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// NotificationEntity represents a notification in the database
type NotificationEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)"`
	UserID    string    `gorm:"index;type:varchar(50)"`
	Type      string    `gorm:"type:varchar(50)"`
	Title     string    `gorm:"type:varchar(255)"`
	Message   string    `gorm:"type:text"`
	Data      []byte    `gorm:"type:json"`
	Read      bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ReadAt    *time.Time
}

// GormNotificationRepository implements port.NotificationRepository using GORM
type GormNotificationRepository struct {
	BaseRepository
}

// NewGormNotificationRepository creates a new GormNotificationRepository
func NewGormNotificationRepository(db *gorm.DB, logger *zerolog.Logger) *GormNotificationRepository {
	return &GormNotificationRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// SavePreferences saves notification preferences for a user
func (r *GormNotificationRepository) SavePreferences(ctx context.Context, userID string, preferences map[string]interface{}) error {
	// Convert preferences to JSON
	preferencesJSON, err := json.Marshal(preferences)
	if err != nil {
		r.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to marshal notification preferences")
		return err
	}

	// Create entity
	entity := &NotificationPreferenceEntity{
		ID:        userID, // Use userID as the ID for simplicity
		UserID:    userID,
		Settings:  preferencesJSON,
		UpdatedAt: time.Now(),
	}

	// Save entity
	return r.Upsert(ctx, entity, []string{"user_id"}, []string{
		"settings", "updated_at",
	})
}

// GetPreferences retrieves notification preferences for a user
func (r *GormNotificationRepository) GetPreferences(ctx context.Context, userID string) (map[string]interface{}, error) {
	var entity NotificationPreferenceEntity
	err := r.FindOne(ctx, &entity, "user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	if entity.ID == "" {
		// Return default preferences if none are found
		return map[string]interface{}{
			"email":    true,
			"push":     true,
			"in_app":   true,
			"trade":    true,
			"system":   true,
			"security": true,
		}, nil
	}

	// Parse preferences
	var preferences map[string]interface{}
	if len(entity.Settings) > 0 {
		if err := json.Unmarshal(entity.Settings, &preferences); err != nil {
			r.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to unmarshal notification preferences")
			return nil, err
		}
	}

	return preferences, nil
}

// SaveNotification saves a notification
func (r *GormNotificationRepository) SaveNotification(ctx context.Context, notification map[string]interface{}) error {
	// Extract required fields
	userID, ok := notification["user_id"].(string)
	if !ok {
		r.logger.Error().Interface("notification", notification).Msg("Missing user_id in notification")
		return fmt.Errorf("missing user_id in notification")
	}

	notificationType, ok := notification["type"].(string)
	if !ok {
		r.logger.Error().Interface("notification", notification).Msg("Missing type in notification")
		return fmt.Errorf("missing type in notification")
	}

	title, ok := notification["title"].(string)
	if !ok {
		r.logger.Error().Interface("notification", notification).Msg("Missing title in notification")
		return fmt.Errorf("missing title in notification")
	}

	message, ok := notification["message"].(string)
	if !ok {
		r.logger.Error().Interface("notification", notification).Msg("Missing message in notification")
		return fmt.Errorf("missing message in notification")
	}

	// Convert data to JSON
	var dataJSON []byte
	if data, ok := notification["data"]; ok && data != nil {
		var err error
		dataJSON, err = json.Marshal(data)
		if err != nil {
			r.logger.Error().Err(err).Interface("notification", notification).Msg("Failed to marshal notification data")
			return err
		}
	}

	// Create entity
	entity := &NotificationEntity{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      notificationType,
		Title:     title,
		Message:   message,
		Data:      dataJSON,
		Read:      false,
		CreatedAt: time.Now(),
	}

	// Save entity
	return r.Create(ctx, entity)
}

// GetNotifications retrieves notifications for a user
func (r *GormNotificationRepository) GetNotifications(ctx context.Context, userID string, limit, offset int) ([]map[string]interface{}, error) {
	var entities []NotificationEntity
	err := r.GetDB(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entities).Error
	if err != nil {
		r.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to get notifications")
		return nil, err
	}

	// Convert to maps
	notifications := make([]map[string]interface{}, len(entities))
	for i, entity := range entities {
		notification := map[string]interface{}{
			"id":         entity.ID,
			"user_id":    entity.UserID,
			"type":       entity.Type,
			"title":      entity.Title,
			"message":    entity.Message,
			"read":       entity.Read,
			"created_at": entity.CreatedAt,
		}

		if entity.ReadAt != nil {
			notification["read_at"] = *entity.ReadAt
		}

		// Parse data
		if len(entity.Data) > 0 {
			var data interface{}
			if err := json.Unmarshal(entity.Data, &data); err != nil {
				r.logger.Error().Err(err).Str("notification_id", entity.ID).Msg("Failed to unmarshal notification data")
			} else {
				notification["data"] = data
			}
		}

		notifications[i] = notification
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read
func (r *GormNotificationRepository) MarkAsRead(ctx context.Context, notificationID string) error {
	now := time.Now()
	return r.Update(ctx, &NotificationEntity{ID: notificationID}, map[string]interface{}{
		"read":    true,
		"read_at": now,
	})
}

// MarkAllAsRead marks all notifications for a user as read
func (r *GormNotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	now := time.Now()
	return r.GetDB(ctx).
		Model(&NotificationEntity{}).
		Where("user_id = ? AND read = ?", userID, false).
		Updates(map[string]interface{}{
			"read":    true,
			"read_at": now,
		}).Error
}

// DeleteNotification deletes a notification
func (r *GormNotificationRepository) DeleteNotification(ctx context.Context, notificationID string) error {
	return r.DeleteByID(ctx, &NotificationEntity{}, notificationID)
}

// Ensure GormNotificationRepository implements port.NotificationRepository
var _ port.NotificationRepository = (*GormNotificationRepository)(nil)
