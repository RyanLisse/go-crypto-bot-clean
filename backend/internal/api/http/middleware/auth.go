package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-crypto-bot-clean/backend/internal/api/http/dto"
)

// AuthMiddleware authenticates requests
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for certain paths
		if isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			unauthorized(w, "Missing Authorization header")
			return
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			unauthorized(w, "Invalid Authorization header format")
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			unauthorized(w, "Empty token")
			return
		}

		// Validate the token (in a real implementation, this would verify the JWT)
		// For now, we'll just check if it's not empty
		if token == "" {
			unauthorized(w, "Invalid token")
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// isPublicPath checks if a path is public (doesn't require authentication)
func isPublicPath(path string) bool {
	publicPaths := []string{
		"/health",
		"/api/auth/login",
		"/api/auth/register",
	}

	for _, publicPath := range publicPaths {
		if path == publicPath {
			return true
		}
	}

	return false
}

// unauthorized returns a 401 Unauthorized response
func unauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	response := dto.ErrorResponse{
		Status:  http.StatusUnauthorized,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}
