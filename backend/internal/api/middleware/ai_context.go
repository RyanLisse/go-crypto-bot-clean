package middleware

import (
	"context"
	"net/http"
	"time"
)

// AIContextMiddleware adds trading context to AI requests
func AIContextMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract user ID from context (set by auth middleware)
			userID, ok := r.Context().Value("userID").(int)
			if !ok {
				http.Error(w, "User ID not found in context", http.StatusUnauthorized)
				return
			}

			// Add context to request
			ctx := context.WithValue(r.Context(), "aiContext", map[string]interface{}{
				"timestamp": time.Now().UTC(),
				"userID":    userID,
			})

			// Pass to next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
