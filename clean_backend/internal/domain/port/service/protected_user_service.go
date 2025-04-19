package service

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// ProtectedUserService defines the interface for retrieving basic protected user information.
type ProtectedUserService interface {
	// GetBasicUserInfo retrieves basic information for a user by their ID.
	GetBasicUserInfo(ctx context.Context, userID string) (*model.User, error)
}
