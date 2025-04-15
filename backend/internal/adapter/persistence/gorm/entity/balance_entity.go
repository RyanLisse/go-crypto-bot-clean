package entity

import (
	"time"
)

// BalanceEntity represents the database model for balances
type BalanceEntity struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	WalletID  uint      `gorm:"not null;index"`
	Asset     string    `gorm:"size:20;not null"`
	Free      float64   `gorm:"type:decimal(18,8);not null"`
	Locked    float64   `gorm:"type:decimal(18,8);not null"`
	Total     float64   `gorm:"type:decimal(18,8);not null"`
	USDValue  float64   `gorm:"type:decimal(18,8);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for the BalanceEntity
func (BalanceEntity) TableName() string {
	return "balance_entities"
}
