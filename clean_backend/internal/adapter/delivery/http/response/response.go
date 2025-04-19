package response

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

// Response is the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success creates a successful response with data
func Success(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// Error creates an error response with code and message
func Error(code string, message string) Response {
	return Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}
}

// ErrorWithDetails creates an error response from an error object
func ErrorWithDetails(err interface{}) Response {
	switch e := err.(type) {
	case *ErrorInfo:
		return Response{
			Success: false,
			Error:   e,
		}
	case string:
		return Error("error", e)
	case error:
		return Error("error", e.Error())
	default:
		return Error("unknown_error", "An unknown error occurred")
	}
}

// WriteJSON writes a JSON response to the http.ResponseWriter
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// If encoding fails, write a simple error message
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success":false,"error":{"code":"encoding_error","message":"Failed to encode response"}}`))
		}
	}
}

// WriteSuccessJSON writes a success JSON response
func WriteSuccessJSON(w http.ResponseWriter, statusCode int, data interface{}, logger *zerolog.Logger) {
	response := Success(data)
	WriteJSON(w, statusCode, response)
}

// WriteErrorJSON writes an error JSON response
func WriteErrorJSON(w http.ResponseWriter, statusCode int, err interface{}, logger *zerolog.Logger) {
	response := ErrorWithDetails(err)
	if logger != nil {
		logger.Error().Interface("error", err).Int("status_code", statusCode).Msg("Error response")
	}
	WriteJSON(w, statusCode, response)
}

// ErrorCode converts an error to an error code string
func ErrorCode(code string) string {
	return code
}
