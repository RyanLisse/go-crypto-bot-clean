package entity

import (
	"time"
)

type MexcApiCredential struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"not null;index"`
	ApiKey    string `gorm:"not null"` // Store encrypted
	ApiSecret string `gorm:"not null"` // Store encrypted
	Label     string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (MexcApiCredential) TableName() string { return "mexc_api_credentials" }
