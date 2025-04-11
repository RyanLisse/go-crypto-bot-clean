// Package models contains the data models for the application
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common fields for all models
type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate is a GORM hook that sets the ID before creating a record
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// SystemInfo represents system information
type SystemInfo struct {
	Base
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	StartTime   time.Time `json:"start_time"`
	Uptime      int64 `json:"uptime"` // in seconds
}

// HealthCheck represents a health check record
type HealthCheck struct {
	Base
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details"`
	Component string    `json:"component"`
}

// LogEntry represents a log entry
type LogEntry struct {
	Base
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
	Details   string    `json:"details"`
}
