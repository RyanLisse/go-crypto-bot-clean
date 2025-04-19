package handler

import (
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/response"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/service"
	"github.com/rs/zerolog"
)

// ProtectedHandler handles requests for protected endpoints.
type ProtectedHandler struct {
	logger               *zerolog.Logger
	protectedUserService service.ProtectedUserService
}

// NewProtectedHandler creates a new ProtectedHandler.
func NewProtectedHandler(logger *zerolog.Logger, protectedUserService service.ProtectedUserService) *ProtectedHandler {
	return &ProtectedHandler{logger: logger, protectedUserService: protectedUserService}
}

// GetProtected handles the GET /api/v1/protected request.
func (h *ProtectedHandler) GetProtected(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		// This should ideally be caught by authentication middleware,
		// but as a safeguard, return an internal server error.
		h.logger.Error().Msg("User ID not found in context in protected handler")
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error: user context missing", h.logger)
		return
	}

	// Use the injected use case to get user information
	user, err := h.protectedUserService.GetBasicUserInfo(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user info from use case")
		// Assuming the use case returns an apperror or a standard error that can be handled.
		// In a real scenario, you might check the type of error and return different HTTP status codes.
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to retrieve user information", h.logger)
		return
	}

	h.logger.Info().Str("userID", userID).Msg("Protected endpoint called")
	// Use user data from the use case result
	response.WriteSuccessJSON(w, http.StatusOK, map[string]string{
		"message":   "This is a protected endpoint",
		"userID":    user.ID,
		"userEmail": user.Email, // Assuming model.User has Email field
		"userName":  user.Name,  // Assuming model.User has Name field
	}, h.logger)
}
