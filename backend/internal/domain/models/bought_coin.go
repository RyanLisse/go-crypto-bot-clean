package models

import (
	"time"

	"gorm.io/gorm"
)

// BoughtCoin represents a coin purchase transaction recorded by the bot.
type BoughtCoin struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Symbol        string         `gorm:"uniqueIndex;not null" json:"symbol"`
	PurchasePrice float64        `gorm:"not null" json:"purchase_price"`
	BuyPrice      float64        `gorm:"-" json:"buy_price,omitempty"` // Alias for PurchasePrice for backward compatibility
	Quantity      float64        `gorm:"not null" json:"quantity"`
	BoughtAt      time.Time      `gorm:"index;not null" json:"bought_at"` // Indexed for querying by time
	StopLoss      float64        `gorm:"not null" json:"stop_loss"`
	TakeProfit    float64        `gorm:"not null" json:"take_profit"`
	CurrentPrice  float64        `gorm:"not null" json:"current_price"` // This might be better managed externally
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"` // Use GORM's soft delete
	IsDeleted     bool           `gorm:"not null" json:"is_deleted"`
}
