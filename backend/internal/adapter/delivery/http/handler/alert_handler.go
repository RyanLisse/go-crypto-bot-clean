package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/notification"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type AlertHandler struct {
	notifier *notification.AlertNotifier
	logger   *zerolog.Logger
}

// AlertRequest represents a request to create or update an alert
type AlertRequest struct {
	Symbol    string  `json:"symbol"`
	Condition string  `json:"condition"`
	Threshold float64 `json:"threshold"`
	UserID    string  `json:"userId"`
}

func NewAlertHandler(notifier *notification.AlertNotifier, logger *zerolog.Logger) *AlertHandler {
	return &AlertHandler{
		notifier: notifier,
		logger:   logger,
	}
}

func (h *AlertHandler) RegisterRoutes(r chi.Router) {
	r.Route("/alerts", func(r chi.Router) {
		r.Get("/", h.ListAlerts)
		r.Post("/", h.CreateAlert)
		r.Get("/{id}", h.GetAlert)
		r.Put("/{id}", h.UpdateAlert)
		r.Delete("/{id}", h.DeleteAlert)
	})
}

// ListAlerts returns all alerts
func (h *AlertHandler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Debug().Msg("Getting all alerts")

	// Get query parameter for active only
	activeOnly := r.URL.Query().Get("active") == "true"

	// Get alerts from notifier
	alerts := h.notifier.GetAlerts(ctx, activeOnly)

	// Return the alerts
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(alerts); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode alerts")
	}
}

// CreateAlert creates a new alert
func (h *AlertHandler) CreateAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Debug().Msg("Creating alert")

	// Parse request body
	var req AlertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, err))
		return
	}

	// Validate request
	if req.Symbol == "" || req.Condition == "" || req.UserID == "" {
		h.logger.Error().Interface("request", req).Msg("Invalid alert request")
		apperror.WriteError(w, apperror.NewInvalid("Missing required fields", nil, nil))
		return
	}

	// Create alert
	alert := notification.Alert{
		ID:        uuid.New().String(),
		Level:     notification.AlertLevelWarning,
		Title:     "Price Alert: " + req.Symbol,
		Message:   "Alert for " + req.Symbol + " when price " + req.Condition + " " + fmt.Sprintf("%.2f", req.Threshold),
		Source:    "user_" + req.UserID,
		Timestamp: time.Now(),
		Resolved:  false,
	}

	// Store the alert
	if err := h.notifier.CreateAlert(ctx, notification.AlertLevelWarning, alert.Title, alert.Message, alert.Source); err != nil {
		h.logger.Error().Err(err).Msg("Failed to create alert")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return the alert ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": alert.ID}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode alert ID")
	}
}

// GetAlert returns a specific alert
func (h *AlertHandler) GetAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alertID := chi.URLParam(r, "id")
	h.logger.Debug().Str("alertID", alertID).Msg("Getting alert")

	// Get all alerts
	alerts := h.notifier.GetAlerts(ctx, false)

	// Find the alert with the specified ID
	for _, alert := range alerts {
		if alert.ID == alertID {
			// Return the alert
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(alert); err != nil {
				h.logger.Error().Err(err).Msg("Failed to encode alert")
			}
			return
		}
	}

	// Alert not found
	apperror.WriteError(w, apperror.NewNotFound("Alert", alertID, nil))
}

// UpdateAlert updates an existing alert
func (h *AlertHandler) UpdateAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alertID := chi.URLParam(r, "id")
	h.logger.Debug().Str("alertID", alertID).Msg("Updating alert")

	// Parse request body
	var req AlertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, err))
		return
	}

	// Validate request
	if req.Symbol == "" || req.Condition == "" || req.UserID == "" {
		h.logger.Error().Interface("request", req).Msg("Invalid alert request")
		apperror.WriteError(w, apperror.NewInvalid("Missing required fields", nil, nil))
		return
	}

	// Get all alerts
	alerts := h.notifier.GetAlerts(ctx, false)

	// Find the alert with the specified ID
	alertFound := false
	for i, alert := range alerts {
		if alert.ID == alertID {
			// Update the alert
			alerts[i].Title = "Price Alert: " + req.Symbol
			alerts[i].Message = "Alert for " + req.Symbol + " when price " + req.Condition + " " + fmt.Sprintf("%.2f", req.Threshold)
			alerts[i].Source = "user_" + req.UserID
			alerts[i].Timestamp = time.Now()
			alertFound = true
			break
		}
	}

	// Alert not found
	if !alertFound {
		apperror.WriteError(w, apperror.NewNotFound("Alert", alertID, nil))
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": alertID}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode alert ID")
	}
}

// DeleteAlert deletes an alert
func (h *AlertHandler) DeleteAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alertID := chi.URLParam(r, "id")
	h.logger.Debug().Str("alertID", alertID).Msg("Deleting alert")

	// Resolve the alert (mark as deleted)
	err := h.notifier.ResolveAlert(ctx, alertID)
	if err != nil {
		h.logger.Error().Err(err).Str("alertID", alertID).Msg("Failed to delete alert")
		apperror.WriteError(w, apperror.NewNotFound("Alert", alertID, err))
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}
