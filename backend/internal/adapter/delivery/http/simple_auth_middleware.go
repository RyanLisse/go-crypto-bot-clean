package http

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/rs/zerolog"
)

// NewSimpleAuthMiddleware creates a new test auth middleware
func NewSimpleAuthMiddleware(logger *zerolog.Logger) middleware.AuthMiddleware {
	return middleware.NewTestAuthMiddleware(logger)
}
