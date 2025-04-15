package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// PositionHandler handles position-related endpoints
type PositionHandler struct {
	useCase usecase.PositionUseCase
	logger  *zerolog.Logger
}

// NewPositionHandler creates a new PositionHandler
func NewPositionHandler(useCase usecase.PositionUseCase, logger *zerolog.Logger) *PositionHandler {
	return &PositionHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// RegisterRoutes registers the position routes
func (h *PositionHandler) RegisterRoutes(r chi.Router) {
	h.logger.Info().Msg("Registering position routes")

	r.Route("/positions", func(r chi.Router) {
		// Create a new position
		r.Post("/", h.CreatePosition)

		// Get all positions
		r.Get("/", h.GetPositions)

		// Get positions by type
		r.Get("/type/{positionType}", h.GetPositionsByType)

		// Get positions by symbol
		r.Get("/symbol/{symbol}", h.GetPositionsBySymbol)

		// Get active positions
		r.Get("/active", h.GetActivePositions)

		// Get closed positions
		r.Get("/closed", h.GetClosedPositions)

		// Position-specific operations
		r.Route("/{positionID}", func(r chi.Router) {
			r.Get("/", h.GetPosition)
			r.Put("/", h.UpdatePosition)
			r.Delete("/", h.DeletePosition)
			r.Put("/close", h.ClosePosition)
			r.Put("/stop-loss", h.SetStopLoss)
			r.Put("/take-profit", h.SetTakeProfit)
			r.Put("/price", h.UpdatePositionPrice)
		})
	})

	h.logger.Info().Msg("Position routes registered")
}

// CreatePosition creates a new position
func (h *PositionHandler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request body
	var req model.PositionCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode position create request")
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, err))
		return
	}

	// Validate required fields
	if req.Symbol == "" {
		apperror.WriteError(w, apperror.NewInvalid("Symbol is required", nil, nil))
		return
	}

	// Create position
	position, err := h.useCase.CreatePosition(ctx, req)
	if err != nil {
		h.logger.Error().Err(err).Str("symbol", req.Symbol).Msg("Failed to create position")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return the created position
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    position,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode position response")
	}
}

// GetPositions returns all positions with pagination
func (h *PositionHandler) GetPositions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// For now, use a test user ID
	userID := "test_user"

	// Get positions
	positions, err := h.useCase.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get positions")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return positions
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    positions,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode positions response")
	}
}

// GetPosition returns a specific position by ID
func (h *PositionHandler) GetPosition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get position ID from URL
	positionID := chi.URLParam(r, "positionID")
	if positionID == "" {
		apperror.WriteError(w, apperror.NewInvalid("Position ID is required", nil, nil))
		return
	}

	// Get position
	position, err := h.useCase.GetPositionByID(ctx, positionID)
	if err != nil {
		h.logger.Error().Err(err).Str("positionID", positionID).Msg("Failed to get position")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return position
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    position,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode position response")
	}
}

// UpdatePosition updates a position
func (h *PositionHandler) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get position ID from URL
	positionID := chi.URLParam(r, "positionID")
	if positionID == "" {
		apperror.WriteError(w, apperror.NewInvalid("Position ID is required", nil, nil))
		return
	}

	// Parse request body
	var req model.PositionUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode position update request")
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, err))
		return
	}

	// Update position
	position, err := h.useCase.UpdatePosition(ctx, positionID, req)
	if err != nil {
		h.logger.Error().Err(err).Str("positionID", positionID).Msg("Failed to update position")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return updated position
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    position,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode position response")
	}
}

// ClosePosition closes a position
func (h *PositionHandler) ClosePosition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get position ID from URL
	positionID := chi.URLParam(r, "positionID")
	if positionID == "" {
		apperror.WriteError(w, apperror.NewInvalid("Position ID is required", nil, nil))
		return
	}

	// Parse request body
	var req struct {
		ExitPrice    float64  `json:"exitPrice"`
		ExitOrderIDs []string `json:"exitOrderIds,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode close position request")
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, err))
		return
	}

	if req.ExitPrice <= 0 {
		apperror.WriteError(w, apperror.NewInvalid("Exit price must be positive", nil, nil))
		return
	}

	// Close position
	position, err := h.useCase.ClosePosition(ctx, positionID, req.ExitPrice, req.ExitOrderIDs)
	if err != nil {
		h.logger.Error().Err(err).Str("positionID", positionID).Msg("Failed to close position")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return closed position
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    position,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode position response")
	}
}

// SetStopLoss sets a stop-loss for a position
func (h *PositionHandler) SetStopLoss(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get position ID from URL
	positionID := chi.URLParam(r, "positionID")
	if positionID == "" {
		apperror.WriteError(w, apperror.NewInvalid("Position ID is required", nil, nil))
		return
	}

	// Parse request body
	var req struct {
		StopLoss float64 `json:"stopLoss"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode set stop-loss request")
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, err))
		return
	}

	if req.StopLoss <= 0 {
		apperror.WriteError(w, apperror.NewInvalid("Stop-loss must be positive", nil, nil))
		return
	}

	// Set stop-loss
	position, err := h.useCase.SetStopLoss(ctx, positionID, req.StopLoss)
	if err != nil {
		h.logger.Error().Err(err).Str("positionID", positionID).Msg("Failed to set stop-loss")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return updated position
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    position,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode position response")
	}
}

// SetTakeProfit sets a take-profit for a position
func (h *PositionHandler) SetTakeProfit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get position ID from URL
	positionID := chi.URLParam(r, "positionID")
	if positionID == "" {
		apperror.WriteError(w, apperror.NewInvalid("Position ID is required", nil, nil))
		return
	}

	// Parse request body
	var req struct {
		TakeProfit float64 `json:"takeProfit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode set take-profit request")
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, err))
		return
	}

	if req.TakeProfit <= 0 {
		apperror.WriteError(w, apperror.NewInvalid("Take-profit must be positive", nil, nil))
		return
	}

	// Set take-profit
	position, err := h.useCase.SetTakeProfit(ctx, positionID, req.TakeProfit)
	if err != nil {
		h.logger.Error().Err(err).Str("positionID", positionID).Msg("Failed to set take-profit")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return updated position
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    position,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode position response")
	}
}

// UpdatePositionPrice updates the current price of a position
func (h *PositionHandler) UpdatePositionPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get position ID from URL
	positionID := chi.URLParam(r, "positionID")
	if positionID == "" {
		apperror.WriteError(w, apperror.NewInvalid("Position ID is required", nil, nil))
		return
	}

	// Parse request body
	var req struct {
		CurrentPrice float64 `json:"currentPrice"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode update price request")
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, err))
		return
	}

	if req.CurrentPrice <= 0 {
		apperror.WriteError(w, apperror.NewInvalid("Current price must be positive", nil, nil))
		return
	}

	// Update price
	position, err := h.useCase.UpdatePositionPrice(ctx, positionID, req.CurrentPrice)
	if err != nil {
		h.logger.Error().Err(err).Str("positionID", positionID).Msg("Failed to update position price")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return updated position
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    position,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode position response")
	}
}

// DeletePosition deletes a position
func (h *PositionHandler) DeletePosition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get position ID from URL
	positionID := chi.URLParam(r, "positionID")
	if positionID == "" {
		apperror.WriteError(w, apperror.NewInvalid("Position ID is required", nil, nil))
		return
	}

	// Delete position
	err := h.useCase.DeletePosition(ctx, positionID)
	if err != nil {
		h.logger.Error().Err(err).Str("positionID", positionID).Msg("Failed to delete position")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Position deleted successfully",
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
	}
}

// GetPositionsByType returns positions filtered by type
func (h *PositionHandler) GetPositionsByType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get position type from URL
	positionTypeStr := chi.URLParam(r, "positionType")
	if positionTypeStr == "" {
		apperror.WriteError(w, apperror.NewInvalid("Position type is required", nil, nil))
		return
	}

	// Convert to position type
	positionType := model.PositionType(positionTypeStr)

	// Get positions by type
	positions, err := h.useCase.GetOpenPositionsByType(ctx, positionType)
	if err != nil {
		h.logger.Error().Err(err).Str("type", positionTypeStr).Msg("Failed to get positions by type")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return positions
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    positions,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode positions response")
	}
}

// GetPositionsBySymbol returns positions filtered by symbol
func (h *PositionHandler) GetPositionsBySymbol(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get symbol from URL
	symbol := chi.URLParam(r, "symbol")
	if symbol == "" {
		apperror.WriteError(w, apperror.NewInvalid("Symbol is required", nil, nil))
		return
	}

	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get positions by symbol
	positions, err := h.useCase.GetPositionsBySymbol(ctx, symbol, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get positions by symbol")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return positions
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    positions,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode positions response")
	}
}

// GetActivePositions returns all active positions
func (h *PositionHandler) GetActivePositions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// For now, use a test user ID
	userID := "test_user"

	// Get active positions
	positions, err := h.useCase.GetActiveByUser(ctx, userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get active positions")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return positions
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    positions,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode positions response")
	}
}

// GetClosedPositions returns closed positions within a time range
func (h *PositionHandler) GetClosedPositions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse time range parameters
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var from, to time.Time
	var err error

	// Parse from time
	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			apperror.WriteError(w, apperror.NewInvalid("Invalid 'from' parameter format. Use RFC3339 format (e.g., 2023-01-01T00:00:00Z)", nil, err))
			return
		}
	} else {
		// Default to 30 days ago
		from = time.Now().AddDate(0, 0, -30)
	}

	// Parse to time
	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			apperror.WriteError(w, apperror.NewInvalid("Invalid 'to' parameter format. Use RFC3339 format (e.g., 2023-01-01T00:00:00Z)", nil, err))
			return
		}
	} else {
		// Default to now
		to = time.Now()
	}

	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get closed positions
	positions, err := h.useCase.GetClosedPositions(ctx, from, to, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get closed positions")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return positions
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    positions,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode positions response")
	}
}

// Helper function to get pagination parameters from request
func getPaginationParams(r *http.Request) (int, int) {
	// Default values
	limit := 10
	offset := 0

	// Parse limit parameter
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Parse offset parameter
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	return limit, offset
}
