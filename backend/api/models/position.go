package models

import (
	"time"

	"gorm.io/gorm"
)

type Position struct {
	ID           string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID       string     `gorm:"index;not null;type:varchar(36)" json:"userId"`
	Symbol       string     `gorm:"index;not null;type:varchar(20)" json:"symbol"`
	Quantity     float64    `gorm:"not null" json:"quantity"`
	EntryPrice   float64    `gorm:"not null" json:"entryPrice"`
	CurrentPrice float64    `gorm:"not null" json:"currentPrice"`
	OpenTime     *time.Time `json:"openTime,omitempty"`
	CloseTime    *time.Time `json:"closeTime,omitempty"`
	Status       string     `gorm:"index;type:varchar(20);not null" json:"status"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
