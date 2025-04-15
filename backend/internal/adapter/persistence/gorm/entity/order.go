package entity

import (
	"time"
)

type OrderEntity struct {
	ID        string    `gorm:"primaryKey"`
	AccountID string    `gorm:"not null;index"`
	Symbol    string    `gorm:"not null;index"`
	Side      string    `gorm:"not null"` // "BUY" or "SELL"
	Type      string    `gorm:"not null"` // "LIMIT", "MARKET", etc.
	Quantity  float64   `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	Status    string    `gorm:"not null"` // "NEW", "FILLED", etc.
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (OrderEntity) TableName() string { return "orders" }
