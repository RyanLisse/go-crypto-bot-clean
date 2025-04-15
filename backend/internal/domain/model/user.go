package model

import (
	"errors"
	"time"
)

// User represents a user in the system
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new User
func NewUser(id, email, name string) *User {
	now := time.Now()
	return &User{
		ID:        id,
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate validates the User
func (u *User) Validate() error {
	if u.ID == "" {
		return ErrInvalidUserID
	}

	if u.Email == "" {
		return errors.New("invalid email")
	}

	return nil
}

// Update updates the User
func (u *User) Update(name string) {
	if name != "" {
		u.Name = name
	}

	u.UpdatedAt = time.Now()
}
