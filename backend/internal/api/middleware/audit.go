package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/audit"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// AuditMiddleware creates a middleware that logs audit events
func AuditMiddleware(auditSvc audit.Service, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get request ID from context or header
			requestID := r.Context().Value("request_id")
			if requestID == nil {
				requestID = r.Header.Get("X-Request-ID")
				if requestID == "" {
					requestID = generateRequestID()
				}
			}

			// Get user ID from context
			userID := 0
			if user := r.Context().Value("user"); user != nil {
				if userMap, ok := user.(map[string]interface{}); ok {
					if id, ok := userMap["id"].(int); ok {
						userID = id
					}
				}
			}

			// Create response writer wrapper to capture status code
			rw := newResponseWriter(w)

			// Record start time
			startTime := time.Now()

			// Create context with request ID
			ctx := context.WithValue(r.Context(), "request_id", requestID)
			r = r.WithContext(ctx)

			// Call the next handler
			next.ServeHTTP(rw, r)

			// Record end time
			duration := time.Since(startTime)

			// Determine event type and severity based on path and status code
			eventType := determineEventType(r.URL.Path)
			severity := determineEventSeverity(rw.statusCode)

			// Create metadata
			metadata := map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"query":       r.URL.RawQuery,
				"status_code": rw.statusCode,
				"duration_ms": duration.Milliseconds(),
				"host":        r.Host,
				"referer":     r.Referer(),
			}

			// Create audit event
			event, err := audit.CreateAuditEvent(
				userID,
				eventType,
				severity,
				r.Method+" "+r.URL.Path,
				fmt.Sprintf("%s %s - %d", r.Method, r.URL.Path, rw.statusCode),
				metadata,
				r.RemoteAddr,
				r.UserAgent(),
				requestID.(string),
			)
			if err != nil {
				logger.Error("Failed to create audit event", zap.Error(err))
				return
			}

			// Log audit event
			if err := auditSvc.LogEvent(r.Context(), event); err != nil {
				logger.Error("Failed to log audit event", zap.Error(err))
			}
		})
	}
}

// responseWriter is a wrapper for http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// newResponseWriter creates a new responseWriter
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// determineEventType determines the event type based on the path
func determineEventType(path string) audit.EventType {
	switch {
	case strings.HasPrefix(path, "/api/auth"):
		return audit.EventTypeAuth
	case strings.HasPrefix(path, "/api/ai"):
		return audit.EventTypeAI
	case strings.HasPrefix(path, "/api/admin"):
		return audit.EventTypeAdmin
	case strings.HasPrefix(path, "/api/trading"):
		return audit.EventTypeTrading
	default:
		return audit.EventTypeAI
	}
}

// determineEventSeverity determines the event severity based on the status code
func determineEventSeverity(statusCode int) audit.EventSeverity {
	switch {
	case statusCode >= 500:
		return audit.EventSeverityError
	case statusCode >= 400:
		return audit.EventSeverityWarning
	default:
		return audit.EventSeverityInfo
	}
}

// RegisterAuditMiddleware registers audit middleware with a Chi router
func RegisterAuditMiddleware(r chi.Router, auditSvc audit.Service, logger *zap.Logger) {
	r.Use(AuditMiddleware(auditSvc, logger))
}
