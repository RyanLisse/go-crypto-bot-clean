package entity

import (
	"time"
)

type PositionEntity struct {
	ID         string    `gorm:"primaryKey"`
	AccountID  string    `gorm:"not null;index"`
	Symbol     string    `gorm:"not null;index"`
	Side       string    `gorm:"not null"` // "LONG" or "SHORT"
	Quantity   float64   `gorm:"not null"`
	EntryPrice float64   `gorm:"not null"`
	Status     string    `gorm:"not null"` // "OPEN", "CLOSED"
	OpenedAt   time.Time `gorm:"not null"`
	ClosedAt   *time.Time
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (PositionEntity) TableName() string { return "positions" }
