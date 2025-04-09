package gorm

import (
	"context"
	"fmt"
	"time"

	notification_domain "go-crypto-bot-clean/backend/internal/domain/notification"
	"go-crypto-bot-clean/backend/internal/domain/notification/ports"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GormNotificationPreferenceRepository implements the preference repository using GORM.
type GormNotificationPreferenceRepository struct {
	db *gorm.DB
}

// NewGormNotificationPreferenceRepository creates a new GORM repository instance.
func NewGormNotificationPreferenceRepository(db *gorm.DB) *GormNotificationPreferenceRepository {
	if db == nil {
		panic("database connection cannot be nil for GormNotificationPreferenceRepository")
	}
	// Ensure the table exists (consider moving migrations elsewhere for production)
	db.AutoMigrate(&PreferenceGORM{})
	return &GormNotificationPreferenceRepository{db: db}
}

// GetUserPreferences retrieves enabled preferences for a user from the database.
func (r *GormNotificationPreferenceRepository) GetUserPreferences(ctx context.Context, userID string) ([]notification_domain.Preference, error) {
	var gormPrefs []PreferenceGORM
	result := r.db.WithContext(ctx).Where("user_id = ? AND enabled = ?", userID, true).Find(&gormPrefs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query notification preferences for user %s: %w", userID, result.Error)
	}

	domainPrefs := make([]notification_domain.Preference, len(gormPrefs))
	for i, gp := range gormPrefs {
		domainPrefs[i] = gp.ToDomain()
	}

	return domainPrefs, nil
}

// SaveUserPreference saves/updates a preference using GORM. (For testing)
// Uses Clauses.OnConflict to perform an UPSERT based on UserID and Channel.
func (r *GormNotificationPreferenceRepository) SaveUserPreference(ctx context.Context, pref notification_domain.Preference) error {
	enabled := pref.Enabled // Create a copy of the bool value
	gormPref := PreferenceGORM{
		UserID:    pref.UserID,
		Channel:   pref.Channel,
		Recipient: pref.Recipient,
		Enabled:   &enabled,
		// Let GORM handle ID, CreatedAt, UpdatedAt, DeletedAt
	}

	// Upsert: If conflict on UserID+Channel, update Recipient and Enabled status explicitly.
	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "channel"}},
		// Use Assignments map to explicitly set values, including false for boolean
		DoUpdates: clause.Assignments(map[string]interface{}{
			"recipient":  pref.Recipient,
			"enabled":    &enabled,   // Use pointer to bool
			"updated_at": time.Now(), // Manually update timestamp
		}),
	}).Create(&gormPref)

	if result.Error != nil {
		return fmt.Errorf("failed to save notification preference (UserID: %s, Channel: %s): %w", pref.UserID, pref.Channel, result.Error)
	}

	return nil
}

// Compile-time interface satisfaction checks
var _ ports.NotificationPreferenceRepository = (*GormNotificationPreferenceRepository)(nil)
var _ interface {
	SaveUserPreference(context.Context, notification_domain.Preference) error
	GetUserPreferences(context.Context, string) ([]notification_domain.Preference, error)
} = (*GormNotificationPreferenceRepository)(nil)
