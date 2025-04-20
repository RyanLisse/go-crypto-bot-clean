package middleware

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// SetUserIDInContext sets the user ID in the context
func SetUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// SetUserInContext sets the user in the context
func SetUserInContext(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

// SetRolesInContext sets the user roles in the context
func SetRolesInContext(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, RolesKey, roles)
}

// HasRole checks if the user has the specified role
func HasRole(ctx context.Context, role string) bool {
	roles, ok := GetRolesFromContext(ctx)
	if !ok {
		return false
	}

	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}
