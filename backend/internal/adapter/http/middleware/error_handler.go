package middleware

import (
	"errors"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/rs/zerolog"
)

// ErrorHandlerMiddleware creates a middleware that handles errors
func ErrorHandlerMiddleware(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a response writer that can capture the status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			// Call the next handler
			next.ServeHTTP(rw, r)
			
			// Log the request
			logger.Debug().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", rw.statusCode).
				Msg("Request processed")
			
			// If there was an error, log it
			if rw.statusCode >= 400 {
				logger.Error().
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Int("status", rw.statusCode).
					Msg("Request error")
			}
		})
	}
}

// RecoverMiddleware creates a middleware that recovers from panics
func RecoverMiddleware(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error().
						Interface("error", err).
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

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and calls the underlying WriteHeader
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the status code if not already set and calls the underlying Write
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}
