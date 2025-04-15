package entity

import (
	"time"
)

type WalletEntity struct {
	ID         string    `gorm:"primaryKey"`
	AccountID  string    `gorm:"not null;index"`
	Exchange   string    `gorm:"not null"`
	TotalUSD   float64   `gorm:"not null"`
	LastUpdate time.Time `gorm:"autoUpdateTime"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (WalletEntity) TableName() string { return "wallets" }
