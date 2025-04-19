package entity

import (
	"time"
)

// APICredentialEntity represents an API credential in the database
type APICredentialEntity struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)"`
	UserID       string    `gorm:"not null;index;type:varchar(36)"`
	Exchange     string    `gorm:"not null;index;type:varchar(20)"`
	APIKey       string    `gorm:"not null;type:varchar(255)"`
	APISecret    string    `gorm:"not null;type:text"`
	Label        string    `gorm:"type:varchar(100)"`
	Status       string    `gorm:"not null;type:varchar(20);default:'active'"`
	FailureCount int       `gorm:"not null;default:0"`
	LastUsed     time.Time
	LastVerified time.Time
	ExpiresAt    time.Time
	RotationDue  time.Time
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for APICredentialEntity
func (APICredentialEntity) TableName() string {
	return "api_credentials"
}

// MexcApiCredential represents a MEXC API credential in the database
type MexcApiCredential struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"not null;index"`
	ApiKey    string `gorm:"not null"` // Store encrypted
	ApiSecret string `gorm:"not null"` // Store encrypted
	Label     string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for MexcApiCredential
func (MexcApiCredential) TableName() string { 
	return "mexc_api_credentials" 
}
