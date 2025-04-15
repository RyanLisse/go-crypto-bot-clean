package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/response"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

// TestHandler handles test-related endpoints
type TestHandler struct {
	cfg    *config.Config
	logger *zerolog.Logger
}

// NewTestHandler creates a new TestHandler
func NewTestHandler(cfg *config.Config, logger *zerolog.Logger) *TestHandler {
	return &TestHandler{
		cfg:    cfg,
		logger: logger,
	}
}

// RegisterRoutes registers the test routes
func (h *TestHandler) RegisterRoutes(r chi.Router) {
	// Simple test endpoint
	r.Get("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info().Msg("Test endpoints list requested")
		response.WriteJSON(w, http.StatusOK, response.Success(map[string]interface{}{
			"endpoints": []string{
				"/api/v1/market/tickers",
				"/api/v1/market/ticker/{symbol}",
				"/api/v1/market/orderbook/{symbol}",
				"/api/v1/market/candles/{symbol}/{interval}",
				"/api/v1/market/symbols",
				"/api/v1/account/wallet",
				"/api/v1/account/balance/{asset}",
				"/api/v1/account/refresh",
				"/api/v1/status/services",
				"/api/v1/status/exchanges",
			},
		}))
	})
}

// GenerateTestToken generates a test JWT token for testing purposes
func (h *TestHandler) GenerateTestToken(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug().Msg("Generating test token")

	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    "test_user_123",
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
		"iat":    time.Now().Unix(),
		"jti":    "test_session_123",
		"sid":    "test_session_123",
		"roles":  []string{"admin", "user"},
		"org_id": "test_org_123",
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(h.cfg.Auth.JWTSecret))
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to sign token")
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
