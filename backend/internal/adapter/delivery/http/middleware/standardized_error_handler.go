package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/rs/zerolog"
)

// responseWriter is a wrapper around http.ResponseWriter that tracks status code
// Used in StandardizedErrorHandler.Middleware and LoggingMiddleware
//
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}


// StandardizedErrorHandler is a middleware that handles errors in a standardized way
type StandardizedErrorHandler struct {
	logger *zerolog.Logger
}

// NewStandardizedErrorHandler creates a new StandardizedErrorHandler
func NewStandardizedErrorHandler(logger *zerolog.Logger) *StandardizedErrorHandler {
	return &StandardizedErrorHandler{
		logger: logger,
	}
}

// Middleware returns a middleware function that handles errors in a standardized way
func (h *StandardizedErrorHandler) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a response writer that can capture the status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			// Call the next handler
			next.ServeHTTP(rw, r)
			
			// Log the request
			h.logger.Debug().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", rw.statusCode).
				Msg("Request processed")
			
			// If there was an error, log it
			if rw.statusCode >= 400 {
				h.logger.Error().
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Int("status", rw.statusCode).
					Msg("Request error")
			}
		})
	}
}

// RecoverMiddleware creates a middleware that recovers from panics with standardized error responses
func (h *StandardizedErrorHandler) RecoverMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic with stack trace
					stack := debug.Stack()
					h.logger.Error().
						Interface("error", err).
						Str("stack", string(stack)).
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Msg("Panic recovered")
					
					// Convert the panic to an error response
					var appErr *apperror.AppError
					switch e := err.(type) {
					case *apperror.AppError:
						appErr = e
					case error:
						appErr = apperror.NewInternal(e)
					case string:
						appErr = apperror.NewInternal(errors.New(e))
					default:
						appErr = apperror.NewInternal(errors.New("unknown panic"))
					}
					
					apperror.WriteError(w, appErr)
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware creates a middleware that logs requests and responses
func (h *StandardizedErrorHandler) LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a response writer that can capture the status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			// Log the request
			h.logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Msg("Request received")
			
			// Call the next handler
			next.ServeHTTP(rw, r)
			
			// Log the response
			h.logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", rw.statusCode).
				Str("remote_addr", r.RemoteAddr).
				Msg("Response sent")
		})
	}
}

// ErrorResponseMiddleware creates a middleware that converts errors to standardized responses
func (h *StandardizedErrorHandler) ErrorResponseMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a response writer that can capture errors
			ew := &errorResponseWriter{
				ResponseWriter: w,
				logger:         h.logger,
				request:        r,
			}
			
			// Call the next handler
			next.ServeHTTP(ew, r)
		})
	}
}

// errorResponseWriter is a wrapper around http.ResponseWriter that captures errors
type errorResponseWriter struct {
	http.ResponseWriter
	logger  *zerolog.Logger
	request *http.Request
}

// WriteHeader captures the status code and writes a standardized error response for error status codes
func (ew *errorResponseWriter) WriteHeader(code int) {
	// If it's an error status code, convert it to a standardized error response
	if code >= 400 {
		var appErr *apperror.AppError
		
		switch code {
		case http.StatusBadRequest:
			appErr = apperror.NewInvalid("Bad request", nil, nil)
		case http.StatusUnauthorized:
			appErr = apperror.NewUnauthorized("Unauthorized", nil)
		case http.StatusForbidden:
			appErr = apperror.NewForbidden("Forbidden", nil)
		case http.StatusNotFound:
			appErr = apperror.NewNotFound("resource", ew.request.URL.Path, nil)
		case http.StatusMethodNotAllowed:
			appErr = apperror.NewInvalid("Method not allowed", nil, nil)
		case http.StatusConflict:
			appErr = &apperror.AppError{
				StatusCode: http.StatusConflict,
				Code:       "CONFLICT",
				Message:    "Resource conflict",
			}
		case http.StatusTooManyRequests:
			appErr = apperror.NewRateLimit("rate_limit_exceeded", nil)
		default:
			appErr = apperror.NewInternal(errors.New("internal server error"))
		}
		
		// Log the error
		ew.logger.Error().
			Int("status", code).
			Str("method", ew.request.Method).
			Str("path", ew.request.URL.Path).
			Msg("Error response")
		
		// Write the error response
		apperror.WriteError(ew.ResponseWriter, appErr)
		return
	}
	
	// Otherwise, just write the header
	ew.ResponseWriter.WriteHeader(code)
}
