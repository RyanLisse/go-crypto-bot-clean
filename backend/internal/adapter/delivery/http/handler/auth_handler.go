package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

// AuthHandler handles authentication-related endpoints
type AuthHandler struct {
	cfg    *config.Config
	logger *zerolog.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(cfg *config.Config, logger *zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		cfg:    cfg,
		logger: logger,
	}
}

// RegisterRoutes registers the authentication routes
func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/token", h.CreateToken)
		r.Get("/test-token", h.CreateTestToken)
	})
}

// TokenRequest represents a request to create a token
type TokenRequest struct {
	UserID string   `json:"userId"`
	Roles  []string `json:"roles"`
}

// TokenResponse represents a response with a token
type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// CreateToken creates a JWT token using Clerk's JWT templates
func (h *AuthHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug().Msg("Creating token")

	// Parse request body
	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, err))
		return
	}

	// Validate request
	if req.UserID == "" {
		h.logger.Error().Interface("request", req).Msg("Invalid token request")
		apperror.WriteError(w, apperror.NewInvalid("Missing required fields", nil, nil))
		return
	}

	// Set default roles if not provided
	if len(req.Roles) == 0 {
		req.Roles = []string{"user"}
	}

	// In a real implementation, you would use the Clerk API to create a JWT token
	// using a JWT template. For now, we'll provide instructions on how to do this.

	// Instructions for creating a JWT token using Clerk's API
	instructions := fmt.Sprintf(`
To create a JWT token using Clerk's API, make the following request:

curl -X POST https://api.clerk.com/v1/jwt_templates/%s/jwt \
  -H "Authorization: Bearer %s" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "%s",
    "claims": {
      "roles": %v
    }
  }'

This will return a JWT token that you can use to authenticate with the API.
`, h.cfg.Auth.ClerkJWTTemplate, h.cfg.Auth.ClerkSecretKey, req.UserID, req.Roles)

	// Return instructions
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"instructions": instructions,
	})
}

// CreateTestToken creates a test JWT token for development purposes
func (h *AuthHandler) CreateTestToken(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug().Msg("Creating test token for development")

	// Only allow in development mode
	if h.cfg.ENV != "development" {
		h.logger.Warn().Msg("Test token endpoint is only available in development mode")
		apperror.WriteError(w, apperror.NewForbidden("This endpoint is only available in development mode", nil))
		return
	}

	// Create a token that mimics a Clerk JWT token
	claims := jwt.MapClaims{
		"sub": "user_2NNPBn8mSWz5KXFMDq9UzCVAq1t", // Example Clerk user ID format
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
		"azp": "clerk",
		"iss": "https://clerk.your-site.com",
		"jti": fmt.Sprintf("test_%d", time.Now().UnixNano()),
		"sid": fmt.Sprintf("sess_%d", time.Now().UnixNano()),
		"Custom": map[string]interface{}{
			"roles": []string{"admin", "user"},
		},
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Add a key ID header to mimic Clerk's JWT format
	token.Header["kid"] = "test_key_id"

	// Sign the token
	tokenString, err := token.SignedString([]byte(h.cfg.Auth.JWTSecret))
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to sign test token")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return the token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
		"note":  "This is a test token for development purposes only. It will not work in production.",
	})
}
