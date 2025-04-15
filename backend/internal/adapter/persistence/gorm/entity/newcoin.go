package entity

import (
	"time"
)

type NewCoinEntity struct {
	ID        string    `gorm:"primaryKey"`
	Symbol    string    `gorm:"not null;uniqueIndex"`
	Status    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (NewCoinEntity) TableName() string { return "new_coins" }
