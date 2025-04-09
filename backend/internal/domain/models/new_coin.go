package models

import (
	"time"
)

// NewCoin represents a newly listed cryptocurrency
type NewCoin struct {
	ID               int64      `gorm:"primaryKey;autoIncrement" db:"id" json:"id"`
	Symbol           string     `gorm:"uniqueIndex;not null" db:"symbol" json:"symbol"`
	FoundAt          time.Time  `gorm:"not null" db:"found_at" json:"found_at"`
	FirstOpenTime    *time.Time `gorm:"" db:"first_open_time" json:"first_open_time,omitempty"`
	BaseVolume       float64    `gorm:"not null;default:0" db:"base_volume" json:"base_volume"`
	QuoteVolume      float64    `gorm:"not null;default:0" db:"quote_volume" json:"quote_volume"`
	Status           string     `gorm:"not null;default:''" db:"status" json:"status"`                // Trading status (e.g., "1" for tradable)
	BecameTradableAt *time.Time `gorm:"" db:"became_tradable_at" json:"became_tradable_at,omitempty"` // When the coin became tradable
	IsProcessed      bool       `gorm:"not null;default:false" db:"is_processed" json:"is_processed"`
	IsDeleted        bool       `gorm:"not null;default:false" db:"is_deleted" json:"is_deleted"`
	IsUpcoming       bool       `gorm:"not null;default:false" db:"is_upcoming" json:"is_upcoming"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
