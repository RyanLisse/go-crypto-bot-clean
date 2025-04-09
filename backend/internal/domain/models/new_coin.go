package models

import (
	"time"

	"gorm.io/gorm"
)

// NewCoin represents a newly listed cryptocurrency on an exchange.
type NewCoin struct {
	ID               int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Symbol           string         `gorm:"uniqueIndex;not null" json:"symbol"`
	FoundAt          time.Time      `gorm:"index;not null" json:"found_at"` // Indexed
	FirstOpenTime    *time.Time     `gorm:"" json:"first_open_time,omitempty"`
	BaseVolume       float64        `gorm:"not null;default:0" json:"base_volume"`
	QuoteVolume      float64        `gorm:"not null;default:0" json:"quote_volume"`
	Status           string         `gorm:"index;not null;default:''" json:"status"`   // Trading status (e.g., "1" for tradable)
	BecameTradableAt *time.Time     `gorm:"index" json:"became_tradable_at,omitempty"` // When the coin became tradable
	IsProcessed      bool           `gorm:"index;not null;default:false" json:"is_processed"`
	IsUpcoming       bool           `gorm:"index;not null;default:false" json:"is_upcoming"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"` // Use GORM's soft delete
	IsDeleted        bool           `gorm:"index;not null;default:false" json:"is_deleted"`
}
