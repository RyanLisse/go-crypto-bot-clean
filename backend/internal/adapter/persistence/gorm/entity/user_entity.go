package entity

import (
	"time"
)

// UserEntity represents the database model for users
type UserEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)"`
	Email     string    `gorm:"unique;not null;type:varchar(100)"`
	Name      string    `gorm:"type:varchar(100)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for the UserEntity
func (UserEntity) TableName() string {
	return "users"
}
