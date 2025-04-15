package entity

import (
	"time"
)

type AccountEntity struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"uniqueIndex;not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (AccountEntity) TableName() string { return "accounts" }
