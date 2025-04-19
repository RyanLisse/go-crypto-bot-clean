package apperror

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog"
)

// ErrorHandler is a function that handles HTTP requests and may return an error
type ErrorHandler func(w http.ResponseWriter, r *http.Request) error

// HandleError wraps an ErrorHandler and handles any returned errors
func HandleError(handler ErrorHandler, logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get request ID from context if available
		requestID := r.Header.Get("X-Request-ID")

		// Call the handler
		err := handler(w, r)
		if err == nil {
			return
		}

		// Log the error
		logEvent := logger.Error().
			Str("method", r.Method).
			Str("path", r.URL.Path)

		if requestID != "" {
			logEvent = logEvent.Str("request_id", requestID)
		}

		logEvent.Err(err).Msg("HTTP handler error")

		// Convert to AppError if needed
		var appErr *AppError
		if !errors.As(err, &appErr) {
			appErr = FromError(err)
		}

		// Respond with error
		RespondWithError(w, r, appErr)
	}
}

// WithLogging adds logging to an HTTP handler
func WithLogging(handler http.Handler, logger *zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get request ID from context if available
		requestID := r.Header.Get("X-Request-ID")

		// Log the request
		logEvent := logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr)

		if requestID != "" {
			logEvent = logEvent.Str("request_id", requestID)
		}

		logEvent.Msg("HTTP request received")

		// Call the handler
		handler.ServeHTTP(w, r)
	})
}

// RespondWithJSON writes a JSON response
func RespondWithJSON(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) error {
	// Set content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Create response
	response := map[string]interface{}{
		"success": true,
		"data":    data,
	}

	// Add request ID if available
	requestID := r.Header.Get("X-Request-ID")
	if requestID != "" {
		response["trace_id"] = requestID
	}

	// Encode response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return NewInternal(err)
	}

	return nil
}

// RespondWithSuccess writes a success response with no data
func RespondWithSuccess(w http.ResponseWriter, r *http.Request, statusCode int) error {
	// Set content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Create response
	response := map[string]interface{}{
		"success": true,
	}

	// Add request ID if available
	requestID := r.Header.Get("X-Request-ID")
	if requestID != "" {
		response["trace_id"] = requestID
	}

	// Encode response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return NewInternal(err)
	}

	return nil
}

// RespondWithCreated writes a 201 Created response with the created resource
func RespondWithCreated(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return RespondWithJSON(w, r, http.StatusCreated, data)
}

// RespondWithOK writes a 200 OK response with the requested resource
func RespondWithOK(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return RespondWithJSON(w, r, http.StatusOK, data)
}

// RespondWithNoContent writes a 204 No Content response
func RespondWithNoContent(w http.ResponseWriter, r *http.Request) error {
	return RespondWithSuccess(w, r, http.StatusNoContent)
}

// RespondWithAccepted writes a 202 Accepted response
func RespondWithAccepted(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return RespondWithJSON(w, r, http.StatusAccepted, data)
}
