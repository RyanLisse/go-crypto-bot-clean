package entity

import (
	"time"
)

// APICredential represents the database model for API credentials
type APICredential struct {
	ID            string    `gorm:"primaryKey;type:varchar(50)"`
	UserID        string    `gorm:"not null;index;type:varchar(50)"`
	Exchange      string    `gorm:"not null;index;type:varchar(20)"`
	APIKey        string    `gorm:"not null;type:varchar(100)"`
	APISecret     []byte    `gorm:"not null;type:blob"`  // Encrypted
	Label         string    `gorm:"type:varchar(50)"`
	Status        string    `gorm:"not null;type:varchar(20);default:'active'"`
	LastUsed      *time.Time `gorm:"type:timestamp"`
	LastVerified  *time.Time `gorm:"type:timestamp"`
	ExpiresAt     *time.Time `gorm:"type:timestamp"`
	RotationDue   *time.Time `gorm:"type:timestamp"`
	FailureCount  int       `gorm:"not null;default:0"`
	Metadata      []byte    `gorm:"type:json"`  // JSON metadata
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for the APICredential
func (APICredential) TableName() string {
	return "api_credentials"
}
