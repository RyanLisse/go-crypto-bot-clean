package gorm

import (
	"time"

	notification_domain "go-crypto-bot-clean/backend/internal/domain/notification"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PreferenceGORM represents the GORM model for notification preferences
type PreferenceGORM struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    string    `gorm:"index:idx_user_channel,unique;not null"`
	Channel   string    `gorm:"index:idx_user_channel,unique;not null"`
	Recipient string    `gorm:"not null"`
	Enabled   *bool     `gorm:""`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUID rather than relying on default value generation.
func (p *PreferenceGORM) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

// TableName specifies the table name for the GORM model.
func (PreferenceGORM) TableName() string {
	return "notification_preferences"
}

// FromDomain converts a domain Preference to a GORM model
func (p *PreferenceGORM) FromDomain(pref notification_domain.Preference) {
	p.UserID = pref.UserID
	p.Channel = pref.Channel
	p.Recipient = pref.Recipient
	enabled := pref.Enabled // Create a copy of the bool value
	p.Enabled = &enabled
}

// ToDomain converts a GORM model to a domain Preference
func (p *PreferenceGORM) ToDomain() notification_domain.Preference {
	enabled := false
	if p.Enabled != nil {
		enabled = *p.Enabled
	}
	return notification_domain.Preference{
		UserID:    p.UserID,
		Channel:   p.Channel,
		Recipient: p.Recipient,
		Enabled:   enabled,
	}
}

// FromDomain converts the domain entity to the GORM model.
// Note: This might not be needed if we only save/update via specific fields or maps,
// but can be useful for direct creation.
func FromDomain(d notification_domain.Preference) PreferenceGORM {
	enabled := d.Enabled // Create a copy of the bool value
	return PreferenceGORM{
		UserID:    d.UserID,
		Channel:   d.Channel,
		Recipient: d.Recipient,
		Enabled:   &enabled, // Use address of the copy
		// ID, CreatedAt, UpdatedAt are usually handled by GORM or the DB
	}
}
