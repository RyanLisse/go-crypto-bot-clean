package util

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
)

// ParseJSONBody parses a JSON request body into dst and returns a standardized error if decoding fails.
func ParseJSONBody(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return apperror.NewInvalid("Request body is empty", nil, nil)
	}
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		if err == io.EOF {
			return apperror.NewInvalid("Request body is empty", nil, err)
		}
		return apperror.NewInvalid("Invalid JSON request body", nil, err)
	}
	return nil
}

// WriteJSONResponse is a convenience wrapper for writing a JSON response with status code.
func WriteJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"error":{"code":"encoding_error","message":"Failed to encode response"}}`))
	}
}
