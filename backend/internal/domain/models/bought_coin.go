package models

import "time"

type BoughtCoin struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" db:"id" json:"id"`
	Symbol        string    `gorm:"uniqueIndex;not null" db:"symbol" json:"symbol"`
	PurchasePrice float64   `gorm:"not null" db:"purchase_price" json:"purchase_price"`
	BuyPrice      float64   `gorm:"-" json:"buy_price,omitempty"` // Alias for PurchasePrice for backward compatibility
	Quantity      float64   `gorm:"not null" db:"quantity" json:"quantity"`
	BoughtAt      time.Time `gorm:"not null" db:"bought_at" json:"bought_at"`
	StopLoss      float64   `gorm:"not null" db:"stop_loss" json:"stop_loss"`
	TakeProfit    float64   `gorm:"not null" db:"take_profit" json:"take_profit"`
	CurrentPrice  float64   `gorm:"not null" db:"current_price" json:"current_price"`
	IsDeleted     bool      `gorm:"not null;default:false" db:"is_deleted" json:"is_deleted,omitempty"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" db:"updated_at" json:"updated_at"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}
