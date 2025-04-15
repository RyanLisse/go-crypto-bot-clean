package entity

import (
	"time"
)

type TransactionEntity struct {
	ID        string    `gorm:"primaryKey"`
	AccountID string    `gorm:"not null;index"`
	Type      string    `gorm:"not null"` // "DEPOSIT", "WITHDRAWAL", "TRADE"
	Asset     string    `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	Status    string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (TransactionEntity) TableName() string { return "transactions" }
