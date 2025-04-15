package middleware

import (
	"net/http"
)

// AuthMiddleware defines the interface for authentication middleware
type AuthMiddleware interface {
	// RequireAuthentication is a middleware that requires authentication
	RequireAuthentication(next http.Handler) http.Handler
}
