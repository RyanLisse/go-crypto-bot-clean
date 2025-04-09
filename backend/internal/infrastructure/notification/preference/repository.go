package preference

import (
	"context"
	"fmt"
	"sync"

	"go-crypto-bot-clean/backend/internal/domain/notification"
	"go-crypto-bot-clean/backend/internal/domain/notification/ports"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PreferenceModel is the GORM model for notification preferences
type PreferenceModel struct {
	gorm.Model
	UserID    string `gorm:"index:idx_user_channel,priority:1"`
	Channel   string `gorm:"index:idx_user_channel,priority:2"`
	Recipient string
	Enabled   bool
}

// TableName returns the table name for the model
func (PreferenceModel) TableName() string {
	return "notification_preferences"
}

// ToEntity converts the model to a domain entity
func (m *PreferenceModel) ToEntity() notification.Preference {
	return notification.Preference{
		UserID:    m.UserID,
		Channel:   m.Channel,
		Recipient: m.Recipient,
		Enabled:   m.Enabled,
	}
}

// FromEntity converts a domain entity to a model
func (m *PreferenceModel) FromEntity(entity notification.Preference) {
	m.UserID = entity.UserID
	m.Channel = entity.Channel
	m.Recipient = entity.Recipient
	m.Enabled = entity.Enabled
}

// Repository implements the NotificationPreferenceRepository interface
type Repository struct {
	db     *gorm.DB
	logger *zap.Logger
	cache  map[string][]notification.Preference
	mu     sync.RWMutex
}

// NewRepository creates a new notification preference repository
func NewRepository(db *gorm.DB, logger *zap.Logger) ports.NotificationPreferenceRepository {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	// Ensure the table exists
	err := db.AutoMigrate(&PreferenceModel{})
	if err != nil {
		logger.Error("Failed to migrate notification preferences table", zap.Error(err))
	}

	return &Repository{
		db:     db,
		logger: logger,
		cache:  make(map[string][]notification.Preference),
	}
}

// GetUserPreferences retrieves all active notification preferences for a given user
func (r *Repository) GetUserPreferences(ctx context.Context, userID string) ([]notification.Preference, error) {
	// Check cache first
	r.mu.RLock()
	if prefs, ok := r.cache[userID]; ok {
		r.mu.RUnlock()
		return prefs, nil
	}
	r.mu.RUnlock()

	// Query database
	var models []PreferenceModel
	result := r.db.WithContext(ctx).Where("user_id = ? AND enabled = ?", userID, true).Find(&models)
	if result.Error != nil {
		r.logger.Error("Failed to get user notification preferences",
			zap.String("user_id", userID),
			zap.Error(result.Error),
		)
		return nil, fmt.Errorf("failed to get user preferences: %w", result.Error)
	}

	// Convert to domain entities
	preferences := make([]notification.Preference, len(models))
	for i, model := range models {
		preferences[i] = model.ToEntity()
	}

	// Update cache
	r.mu.Lock()
	r.cache[userID] = preferences
	r.mu.Unlock()

	return preferences, nil
}

// SaveUserPreference saves a user notification preference
func (r *Repository) SaveUserPreference(ctx context.Context, pref notification.Preference) error {
	// Check if preference already exists
	var model PreferenceModel
	result := r.db.WithContext(ctx).Where("user_id = ? AND channel = ?", pref.UserID, pref.Channel).First(&model)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new preference
			model = PreferenceModel{}
			model.FromEntity(pref)
			result = r.db.WithContext(ctx).Create(&model)
			if result.Error != nil {
				r.logger.Error("Failed to create user notification preference",
					zap.String("user_id", pref.UserID),
					zap.String("channel", pref.Channel),
					zap.Error(result.Error),
				)
				return fmt.Errorf("failed to create user preference: %w", result.Error)
			}
		} else {
			r.logger.Error("Failed to query user notification preference",
				zap.String("user_id", pref.UserID),
				zap.String("channel", pref.Channel),
				zap.Error(result.Error),
			)
			return fmt.Errorf("failed to query user preference: %w", result.Error)
		}
	} else {
		// Update existing preference
		model.FromEntity(pref)
		result = r.db.WithContext(ctx).Save(&model)
		if result.Error != nil {
			r.logger.Error("Failed to update user notification preference",
				zap.String("user_id", pref.UserID),
				zap.String("channel", pref.Channel),
				zap.Error(result.Error),
			)
			return fmt.Errorf("failed to update user preference: %w", result.Error)
		}
	}

	// Invalidate cache
	r.mu.Lock()
	delete(r.cache, pref.UserID)
	r.mu.Unlock()

	return nil
}

// DeleteUserPreference deletes a user notification preference
func (r *Repository) DeleteUserPreference(ctx context.Context, userID, channel string) error {
	result := r.db.WithContext(ctx).Where("user_id = ? AND channel = ?", userID, channel).Delete(&PreferenceModel{})
	if result.Error != nil {
		r.logger.Error("Failed to delete user notification preference",
			zap.String("user_id", userID),
			zap.String("channel", channel),
			zap.Error(result.Error),
		)
		return fmt.Errorf("failed to delete user preference: %w", result.Error)
	}

	// Invalidate cache
	r.mu.Lock()
	delete(r.cache, userID)
	r.mu.Unlock()

	return nil
}

// GetAllUserPreferences retrieves all notification preferences (including disabled ones) for a given user
func (r *Repository) GetAllUserPreferences(ctx context.Context, userID string) ([]notification.Preference, error) {
	var models []PreferenceModel
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models)
	if result.Error != nil {
		r.logger.Error("Failed to get all user notification preferences",
			zap.String("user_id", userID),
			zap.Error(result.Error),
		)
		return nil, fmt.Errorf("failed to get all user preferences: %w", result.Error)
	}

	// Convert to domain entities
	preferences := make([]notification.Preference, len(models))
	for i, model := range models {
		preferences[i] = model.ToEntity()
	}

	return preferences, nil
}
