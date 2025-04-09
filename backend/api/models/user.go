// Package models contains the database models for the API
package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID            string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Email         string         `gorm:"uniqueIndex;not null;type:varchar(255)" json:"email"`
	Username      string         `gorm:"uniqueIndex;not null;type:varchar(50)" json:"username"`
	PasswordHash  string         `gorm:"not null;type:varchar(255)" json:"-"`
	FirstName     string         `gorm:"type:varchar(50)" json:"firstName,omitempty"`
	LastName      string         `gorm:"type:varchar(50)" json:"lastName,omitempty"`
	Roles         []UserRole     `gorm:"foreignKey:UserID" json:"roles"`
	Settings      *UserSettings  `gorm:"foreignKey:UserID" json:"settings,omitempty"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	LastLoginAt   *time.Time     `json:"lastLoginAt,omitempty"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// UserRole represents a role assigned to a user
type UserRole struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string         `gorm:"index;not null;type:varchar(36)" json:"userId"`
	Role      string         `gorm:"not null;type:varchar(50)" json:"role"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the UserRole model
func (UserRole) TableName() string {
	return "user_roles"
}

// UserSettings represents user settings
type UserSettings struct {
	ID                  uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID              string         `gorm:"uniqueIndex;not null;type:varchar(36)" json:"userId"`
	Theme               string         `gorm:"not null;type:varchar(20);default:'light'" json:"theme"`
	Language            string         `gorm:"not null;type:varchar(10);default:'en'" json:"language"`
	TimeZone            string         `gorm:"not null;type:varchar(50);default:'UTC'" json:"timeZone"`
	NotificationsEnabled bool           `gorm:"not null;default:true" json:"notificationsEnabled"`
	EmailNotifications   bool           `gorm:"not null;default:true" json:"emailNotifications"`
	PushNotifications    bool           `gorm:"not null;default:false" json:"pushNotifications"`
	DefaultCurrency      string         `gorm:"not null;type:varchar(10);default:'USD'" json:"defaultCurrency"`
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the UserSettings model
func (UserSettings) TableName() string {
	return "user_settings"
}

// RefreshToken represents a refresh token for a user
type RefreshToken struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string         `gorm:"index;not null;type:varchar(36)" json:"userId"`
	Token     string         `gorm:"uniqueIndex;not null;type:varchar(255)" json:"-"`
	ExpiresAt time.Time      `gorm:"not null" json:"expiresAt"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	RevokedAt *time.Time     `json:"revokedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the RefreshToken model
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
