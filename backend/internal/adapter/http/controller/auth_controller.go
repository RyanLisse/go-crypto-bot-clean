package controller

import (
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/util"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	authService service.AuthServiceInterface
	logger      *zerolog.Logger
}

// NewAuthController creates a new AuthController
func NewAuthController(authService service.AuthServiceInterface, logger *zerolog.Logger) *AuthController {
	return &AuthController{
		authService: authService,
		logger:      logger,
	}
}

// RegisterRoutes registers the authentication routes
func (c *AuthController) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/verify", c.VerifyToken)
	})
}

// VerifyToken handles the verify token endpoint
func (c *AuthController) VerifyToken(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request struct {
		Token string `json:"token"`
	}
	// Use standardized JSON body parsing utility for better error handling
	if err := util.ParseJSONBody(r, &request); err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			apperror.WriteError(w, appErr)
		} else {
			apperror.WriteError(w, apperror.NewInternal(err))
		}
		return
	}

	// Verify token
	userID, err := c.authService.VerifyToken(r.Context(), request.Token)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to verify token")
		apperror.WriteError(w, apperror.NewUnauthorized("Invalid token", err))
		return
	}

	// Get user
	user, err := c.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		c.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Get user roles
	roles, err := c.authService.GetUserRoles(r.Context(), userID)
	if err != nil {
		c.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user roles")
		// Continue with empty roles
		roles = []string{"user"}
	}

	// Create response
	response := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"roles": roles,
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	util.WriteJSONResponse(w, http.StatusOK, response)
}
