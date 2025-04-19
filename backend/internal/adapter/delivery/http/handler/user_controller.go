package controller

import (
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/util"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// UserController handles user-related HTTP requests
type UserController struct {
	userService service.UserServiceInterface
	authService service.AuthServiceInterface
	logger      *zerolog.Logger
}

// NewUserController creates a new UserController
func NewUserController(userService service.UserServiceInterface, authService service.AuthServiceInterface, logger *zerolog.Logger) *UserController {
	return &UserController{
		userService: userService,
		authService: authService,
		logger:      logger,
	}
}

// RegisterRoutes registers the user routes
func (c *UserController) RegisterRoutes(r chi.Router, authMiddleware middleware.AuthMiddleware) {
	r.Route("/users", func(r chi.Router) {
		// Public routes
		r.Get("/health", c.HealthCheck)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuthentication)

			// Current user routes
			r.Get("/me", c.GetCurrentUser)
			r.Put("/me", c.UpdateCurrentUser)

			// Admin-only routes
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireRole("admin"))

				r.Get("/", c.ListUsers)
				r.Get("/{id}", c.GetUserByID)
				r.Delete("/{id}", c.DeleteUser)
			})
		})
	})
}

// HealthCheck handles the health check endpoint
func (c *UserController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	util.WriteJSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

// GetCurrentUser handles the get current user endpoint
func (c *UserController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get user
	user, err := c.userService.GetUserByID(r.Context(), userID)
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
	util.WriteJSONResponse(w, http.StatusOK, response)
}

// UpdateCurrentUser handles the update current user endpoint
func (c *UserController) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Parse request body
	var request struct {
		Name string `json:"name"`
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

	// Update user
	user, err := c.userService.UpdateUser(r.Context(), userID, request.Name)
	if err != nil {
		c.logger.Error().Err(err).Str("userID", userID).Msg("Failed to update user")
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
	util.WriteJSONResponse(w, http.StatusOK, response)
}

// ListUsers handles the list users endpoint
func (c *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Get users
	users, err := c.userService.ListUsers(r.Context())
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to list users")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Create response
	response := make([]map[string]interface{}, len(users))
	for i, user := range users {
		response[i] = map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		}
	}

	// Write response
	util.WriteJSONResponse(w, http.StatusOK, response)
}

// GetUserByID handles the get user by ID endpoint
func (c *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL
	userID := chi.URLParam(r, "id")
	if userID == "" {
		apperror.WriteError(w, apperror.NewInvalid("User ID is required", nil, nil))
		return
	}

	// Get user
	user, err := c.userService.GetUserByID(r.Context(), userID)
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
	util.WriteJSONResponse(w, http.StatusOK, response)
}

// DeleteUser handles the delete user endpoint
func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL
	userID := chi.URLParam(r, "id")
	if userID == "" {
		apperror.WriteError(w, apperror.NewInvalid("User ID is required", nil, nil))
		return
	}

	// Delete user
	if err := c.userService.DeleteUser(r.Context(), userID); err != nil {
		c.logger.Error().Err(err).Str("userID", userID).Msg("Failed to delete user")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.WriteHeader(http.StatusNoContent)
}
