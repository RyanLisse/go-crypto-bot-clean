// Package handlers contains HTTP request handlers.
package handlers

import (
	"encoding/json"
	"net/http"
)

// HTTPHandler is a standard HTTP handler function type
type HTTPHandler func(http.ResponseWriter, *http.Request)

// ResponseData represents the standard response structure
type ResponseData struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// SendJSON sends a JSON response with the given status code and data
func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// SendError sends an error response with the given status code and message
func SendError(w http.ResponseWriter, statusCode int, message string) {
	SendJSON(w, statusCode, ResponseData{
		Status:  "error",
		Message: message,
	})
}

// SendSuccess sends a success response with the given data
func SendSuccess(w http.ResponseWriter, data interface{}) {
	SendJSON(w, http.StatusOK, ResponseData{
		Status: "success",
		Data:   data,
	})
}
