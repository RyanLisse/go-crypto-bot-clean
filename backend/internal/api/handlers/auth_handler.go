package handlers

import (
	"net/http"

	"go-crypto-bot-clean/backend/internal/auth" // Corrected path
)

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthHandler handles authentication-related HTTP requests (primarily getting user info from context)
type AuthHandler struct {
	// No service needed now, user data comes from context via Clerk middleware
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	// No service injection needed
	return &AuthHandler{}
}

// Login/Logout endpoints likely handled by Clerk frontend/middleware now, removed from backend handler.

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user data directly from context populated by Clerk middleware
	// Note: We import the auth package containing GetUserFromContext
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		// This implies the Clerk middleware didn't run or failed
		SendError(w, http.StatusUnauthorized, "User data not found in context")
		return
	}

	SendSuccess(w, user)
}
