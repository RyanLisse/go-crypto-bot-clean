package models

import (
	"time"

	"gorm.io/gorm"
)

// TradeSide represents the side of a trade (buy/sell)
type TradeSide string

const (
	TradeSideBuy  TradeSide = "BUY"
	TradeSideSell TradeSide = "SELL"
)

// Trade represents a trade in the system
type Trade struct {
	gorm.Model
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Symbol    string    `gorm:"index;not null"`
	Price     float64   `gorm:"type:decimal(20,8);not null"`
	Amount    float64   `gorm:"type:decimal(20,8);not null"`
	Side      string    `gorm:"type:varchar(4);not null"` // buy or sell
	OrderID   string    `gorm:"type:uuid;index"`
	TradeID   string    `gorm:"uniqueIndex;not null"` // Exchange trade ID
	Exchange  string    `gorm:"type:varchar(20);not null"`
	TradeTime time.Time `gorm:"index;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// CalculateValue returns the total value of the trade
func (t *Trade) CalculateValue() float64 {
	return t.Price * t.Amount
}

// IsValid checks if the trade has valid properties
func (t *Trade) IsValid() bool {
	return t.Symbol != "" && t.Price > 0 && t.Amount > 0
}
