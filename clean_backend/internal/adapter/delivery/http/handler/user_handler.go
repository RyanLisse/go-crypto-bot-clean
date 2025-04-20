package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/response"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	logger      *zerolog.Logger
	authService port.AuthServiceInterface
	userService port.UserServiceInterface
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(
	logger *zerolog.Logger,
	authService port.AuthServiceInterface,
	userService port.UserServiceInterface,
) *UserHandler {
	return &UserHandler{
		logger:      logger,
		authService: authService,
		userService: userService,
	}
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		// Current user routes (protected)
		r.Get("/me", h.GetCurrentUser)
		r.Put("/me", h.UpdateCurrentUser)

		// Admin-only routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Get("/", h.ListUsers)
			r.Post("/", h.CreateUser)
			r.Get("/{id}", h.GetUserByID)
			r.Put("/{id}", h.UpdateUser)
			r.Delete("/{id}", h.DeleteUser)
			r.Post("/{id}/roles/{role}", h.AddUserRole)
			r.Delete("/{id}/roles/{role}", h.RemoveUserRole)
		})
	})
}

// GetCurrentUser handles the GET /api/v1/users/me request
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.logger.Error().Msg("User ID not found in context")
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error", h.logger)
		return
	}

	// Get user from user service
	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user")
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to get user", h.logger)
		return
	}

	// Get user roles
	roles, err := h.authService.GetUserRoles(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user roles")
		// Continue with empty roles
		roles = []string{"user"}
	}

	// Create response
	resp := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"roles": roles,
	}

	// Write response
	response.WriteSuccessJSON(w, http.StatusOK, resp, h.logger)
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name string `json:"name"`
}

// UpdateCurrentUser handles the PUT /api/v1/users/me request
func (h *UserHandler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.logger.Error().Msg("User ID not found in context")
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error", h.logger)
		return
	}

	// Parse request body
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid request body", h.logger)
		return
	}

	// Validate request
	if req.Name == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Name is required", h.logger)
		return
	}

	// Get user from user service
	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user")
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to get user", h.logger)
		return
	}

	// Update user
	user.Name = req.Name
	user, err = h.userService.UpdateUser(r.Context(), user)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to update user")
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to update user", h.logger)
		return
	}

	// Get user roles
	roles, err := h.authService.GetUserRoles(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user roles")
		// Continue with empty roles
		roles = []string{"user"}
	}

	// Create response
	resp := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"roles": roles,
	}

	// Write response
	response.WriteSuccessJSON(w, http.StatusOK, resp, h.logger)
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// CreateUser handles the POST /api/v1/users request
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid request body", h.logger)
		return
	}

	// Validate request
	if req.Email == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Email is required", h.logger)
		return
	}

	// Create user
	user := &model.User{
		Email: req.Email,
		Name:  req.Name,
	}

	user, err := h.userService.CreateUser(r.Context(), user)
	if err != nil {
		h.logger.Error().Err(err).Str("email", req.Email).Msg("Failed to create user")
		if apperror.IsConflict(err) {
			response.WriteErrorJSON(w, http.StatusConflict, "User with email already exists", h.logger)
		} else {
			response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to create user", h.logger)
		}
		return
	}

	// Create response
	resp := map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"roles": []string{"user"},
	}

	// Write response
	response.WriteSuccessJSON(w, http.StatusCreated, resp, h.logger)
}

// ListUsers handles the GET /api/v1/users request
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Get users from user service
	users, err := h.userService.ListUsers(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list users")
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to list users", h.logger)
		return
	}

	// Create response
	var resp []map[string]any
	for _, user := range users {
		// Get user roles
		roles, err := h.authService.GetUserRoles(r.Context(), user.ID)
		if err != nil {
			h.logger.Error().Err(err).Str("userID", user.ID).Msg("Failed to get user roles")
			// Continue with empty roles
			roles = []string{"user"}
		}

		resp = append(resp, map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"roles": roles,
		})
	}

	// Write response
	response.WriteSuccessJSON(w, http.StatusOK, resp, h.logger)
}

// GetUserByID handles the GET /api/v1/users/{id} request
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL
	userID := chi.URLParam(r, "id")
	if userID == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "User ID is required", h.logger)
		return
	}

	// Get user from user service
	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user")
		if apperror.IsNotFound(err) {
			response.WriteErrorJSON(w, http.StatusNotFound, "User not found", h.logger)
		} else {
			response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to get user", h.logger)
		}
		return
	}

	// Get user roles
	roles, err := h.authService.GetUserRoles(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user roles")
		// Continue with empty roles
		roles = []string{"user"}
	}

	// Create response
	resp := map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"roles": roles,
	}

	// Write response
	response.WriteSuccessJSON(w, http.StatusOK, resp, h.logger)
}

// UpdateUser handles the PUT /api/v1/users/{id} request
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL
	userID := chi.URLParam(r, "id")
	if userID == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "User ID is required", h.logger)
		return
	}

	// Parse request body
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid request body", h.logger)
		return
	}

	// Validate request
	if req.Name == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Name is required", h.logger)
		return
	}

	// Get user from user service
	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user")
		if apperror.IsNotFound(err) {
			response.WriteErrorJSON(w, http.StatusNotFound, "User not found", h.logger)
		} else {
			response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to get user", h.logger)
		}
		return
	}

	// Update user
	user.Name = req.Name
	user, err = h.userService.UpdateUser(r.Context(), user)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to update user")
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to update user", h.logger)
		return
	}

	// Get user roles
	roles, err := h.authService.GetUserRoles(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user roles")
		// Continue with empty roles
		roles = []string{"user"}
	}

	// Create response
	resp := map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"roles": roles,
	}

	// Write response
	response.WriteSuccessJSON(w, http.StatusOK, resp, h.logger)
}

// DeleteUser handles the DELETE /api/v1/users/{id} request
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL
	userID := chi.URLParam(r, "id")
	if userID == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "User ID is required", h.logger)
		return
	}

	// Delete user
	err := h.userService.DeleteUser(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to delete user")
		if apperror.IsNotFound(err) {
			response.WriteErrorJSON(w, http.StatusNotFound, "User not found", h.logger)
		} else {
			response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to delete user", h.logger)
		}
		return
	}

	// Write response
	response.WriteSuccessJSON(w, http.StatusOK, map[string]any{"message": "User deleted successfully"}, h.logger)
}

// AddUserRole handles the POST /api/v1/users/{id}/roles/{role} request
func (h *UserHandler) AddUserRole(w http.ResponseWriter, r *http.Request) {
	// Get user ID and role from URL
	userID := chi.URLParam(r, "id")
	if userID == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "User ID is required", h.logger)
		return
	}

	role := chi.URLParam(r, "role")
	if role == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Role is required", h.logger)
		return
	}

	// Check if user exists
	_, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user")
		if apperror.IsNotFound(err) {
			response.WriteErrorJSON(w, http.StatusNotFound, "User not found", h.logger)
		} else {
			response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to get user", h.logger)
		}
		return
	}

	// Set user role
	err = h.userService.SetUserRole(r.Context(), userID, role)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Str("role", role).Msg("Failed to set user role")
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to set user role", h.logger)
		return
	}

	// Write response
	response.WriteSuccessJSON(w, http.StatusOK, map[string]any{"message": "Role added successfully"}, h.logger)
}

// RemoveUserRole handles the DELETE /api/v1/users/{id}/roles/{role} request
func (h *UserHandler) RemoveUserRole(w http.ResponseWriter, r *http.Request) {
	// Get user ID and role from URL
	userID := chi.URLParam(r, "id")
	if userID == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "User ID is required", h.logger)
		return
	}

	role := chi.URLParam(r, "role")
	if role == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Role is required", h.logger)
		return
	}

	// Check if user exists
	_, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user")
		if apperror.IsNotFound(err) {
			response.WriteErrorJSON(w, http.StatusNotFound, "User not found", h.logger)
		} else {
			response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to get user", h.logger)
		}
		return
	}

	// Remove user role
	err = h.userService.RemoveUserRole(r.Context(), userID, role)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Str("role", role).Msg("Failed to remove user role")
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to remove user role", h.logger)
		return
	}

	// Write response
	response.WriteSuccessJSON(w, http.StatusOK, map[string]any{"message": "Role removed successfully"}, h.logger)
}
