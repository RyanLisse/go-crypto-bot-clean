package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"

	"go-crypto-bot-clean/backend/internal/api/http/dto"
)

// RecoveryMiddleware recovers from panics and returns a 500 error
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace
				log.Printf("Panic: %v\n%s", err, debug.Stack())

				// Return a 500 error
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				response := dto.ErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: "Internal server error",
				}

				json.NewEncoder(w).Encode(response)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
