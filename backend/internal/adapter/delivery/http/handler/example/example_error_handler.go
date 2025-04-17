package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// ExampleErrorInput represents a validation example request
type ExampleErrorInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

// Validate checks if the input is valid
func (i *ExampleErrorInput) Validate() map[string]string {
	errors := make(map[string]string)

	if i.Email == "" {
		errors["email"] = "Email is required"
	} else if !isValidEmail(i.Email) {
		errors["email"] = "Email format is invalid"
	}

	if i.Username == "" {
		errors["username"] = "Username is required"
	} else if len(i.Username) < 3 {
		errors["username"] = "Username must be at least 3 characters"
	}

	if i.Age < 18 {
		errors["age"] = "Age must be at least 18"
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Helper for email validation
func isValidEmail(email string) bool {
	// Simple check for demonstration
	return len(email) > 3 && (email[len(email)-4:] == ".com" || email[len(email)-3:] == ".io")
}

// ErrorExampleController demonstrates different error scenarios
type ErrorExampleController struct {
	logger *zerolog.Logger
}

// NewErrorExampleController creates a new error example controller
func NewErrorExampleController(logger *zerolog.Logger) *ErrorExampleController {
	return &ErrorExampleController{
		logger: logger,
	}
}

// RegisterRoutes registers routes for the error example controller
func (c *ErrorExampleController) RegisterRoutes(r chi.Router) {
	r.Route("/errors", func(r chi.Router) {
		r.Get("/not-found", c.NotFoundExample)
		r.Get("/unauthorized", c.UnauthorizedExample)
		r.Get("/forbidden", c.ForbiddenExample)
		r.Get("/internal", c.InternalErrorExample)
		r.Get("/validation/{type}", c.ValidationErrorExample)
		r.Post("/validation", c.ValidationInputExample)
		r.Get("/wrapped", c.WrappedErrorExample)
		r.Get("/external-api", c.ExternalAPIErrorExample)
		r.Get("/panic", c.PanicExample)
	})
}

// NotFoundExample demonstrates a not found error
func (c *ErrorExampleController) NotFoundExample(w http.ResponseWriter, r *http.Request) {
	// Simulate a not found error
	err := apperror.NewNotFound("user", "123", nil)
	apperror.RespondWithError(w, r, err)
}

// UnauthorizedExample demonstrates an unauthorized error
func (c *ErrorExampleController) UnauthorizedExample(w http.ResponseWriter, r *http.Request) {
	// Simulate an unauthorized error
	err := apperror.NewUnauthorized("Authentication token is invalid or expired", nil)
	apperror.RespondWithError(w, r, err)
}

// ForbiddenExample demonstrates a forbidden error
func (c *ErrorExampleController) ForbiddenExample(w http.ResponseWriter, r *http.Request) {
	// Simulate a forbidden error
	err := apperror.NewForbidden("You do not have permission to access this resource", nil)
	apperror.RespondWithError(w, r, err)
}

// InternalErrorExample demonstrates an internal server error
func (c *ErrorExampleController) InternalErrorExample(w http.ResponseWriter, r *http.Request) {
	// Simulate an internal error
	originalErr := errors.New("database connection failed: connection timeout")
	err := apperror.NewInternal(originalErr)
	apperror.RespondWithError(w, r, err)
}

// ValidationErrorExample demonstrates validation errors
func (c *ErrorExampleController) ValidationErrorExample(w http.ResponseWriter, r *http.Request) {
	errorType := chi.URLParam(r, "type")

	traceID := apperror.GetTraceID(r)

	switch errorType {
	case "single":
		// Single field error
		apperror.WriteValidationError(w, "email", "Email format is invalid", traceID)
	case "multiple":
		// Multiple field errors
		fieldErrors := map[string]string{
			"email":    "Email format is invalid",
			"username": "Username must be at least 3 characters",
			"password": "Password must contain at least one number and one special character",
		}
		apperror.WriteValidationErrors(w, fieldErrors, traceID)
	default:
		// Invalid type
		err := apperror.NewInvalid("Invalid validation error type", nil, nil)
		apperror.RespondWithError(w, r, err)
	}
}

// ValidationInputExample demonstrates input validation
func (c *ErrorExampleController) ValidationInputExample(w http.ResponseWriter, r *http.Request) {
	var input ExampleErrorInput

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperror.RespondWithError(w, r, apperror.NewInvalid("Invalid JSON payload", nil, err))
		return
	}

	// Validate input
	if fieldErrors := input.Validate(); fieldErrors != nil {
		traceID := apperror.GetTraceID(r)
		apperror.WriteValidationErrors(w, fieldErrors, traceID)
		return
	}

	// Successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Validation successful",
		"data":    input,
	})
}

// WrappedErrorExample demonstrates error wrapping
func (c *ErrorExampleController) WrappedErrorExample(w http.ResponseWriter, r *http.Request) {
	// Simulate wrapped errors
	baseErr := errors.New("database query failed")
	repoErr := apperror.WrapError(baseErr, "user repository error")
	serviceErr := apperror.WrapError(repoErr, "user service error")
	finalErr := apperror.WrapError(serviceErr, "get user operation failed")

	apperror.RespondWithError(w, r, finalErr)
}

// ExternalAPIErrorExample demonstrates external API errors
func (c *ErrorExampleController) ExternalAPIErrorExample(w http.ResponseWriter, r *http.Request) {
	// Simulate an external API error
	baseErr := fmt.Errorf("HTTP 503 Service Unavailable")
	err := apperror.NewExternalService("payment-gateway", "Failed to process payment", baseErr)

	apperror.RespondWithError(w, r, err)
}

// PanicExample demonstrates panic recovery
func (c *ErrorExampleController) PanicExample(w http.ResponseWriter, r *http.Request) {
	// Simulate a panic
	panic("This is a simulated panic that will be recovered")
}
