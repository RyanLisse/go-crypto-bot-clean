package models

import (
	"time"

	"gorm.io/gorm"
)

// BalanceHistory represents a snapshot of account balance and equity at a specific time.
type BalanceHistory struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Timestamp     time.Time      `gorm:"index;not null" json:"timestamp"` // Indexed for time-series queries
	Balance       float64        `gorm:"not null" json:"balance"`
	Equity        float64        `gorm:"not null" json:"equity"`
	FreeBalance   float64        `gorm:"not null" json:"free_balance"`
	LockedBalance float64        `gorm:"not null" json:"locked_balance"`
	UnrealizedPnL float64        `gorm:"not null" json:"unrealized_pnl"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"` // Added soft delete
}

// TableName specifies the table name for the BalanceHistory model
func (BalanceHistory) TableName() string {
	return "balance_history"
}

// BalanceHistoryEntry represents a simplified balance history entry for API responses
type BalanceHistoryEntry struct {
	Date    time.Time `json:"date"`
	Balance float64   `json:"balance"`
}
