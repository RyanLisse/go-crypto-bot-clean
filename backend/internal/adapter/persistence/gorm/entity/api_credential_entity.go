package entity

import (
	"time"
)

// APICredentialEntity represents the database model for API credentials
type APICredentialEntity struct {
	ID         string    `gorm:"primaryKey;type:varchar(50)"`
	UserID     string    `gorm:"not null;index;type:varchar(50)"`
	Exchange   string    `gorm:"not null;index;type:varchar(20)"`
	APIKey     string    `gorm:"not null;type:varchar(100)"`
	APISecret  []byte    `gorm:"not null;type:blob"`  // Encrypted
	Label      string    `gorm:"type:varchar(50)"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for the APICredentialEntity
func (APICredentialEntity) TableName() string {
	return "api_credentials"
}
