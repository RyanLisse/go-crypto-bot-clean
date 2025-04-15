package entity

import (
	"time"
)

type StatusEntity struct {
	ID        string    `gorm:"primaryKey"`
	Name      string    `gorm:"not null;uniqueIndex"`
	Value     string    `gorm:"not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (StatusEntity) TableName() string { return "statuses" }
