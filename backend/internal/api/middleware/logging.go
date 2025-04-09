// Package middleware contains API middleware components.
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Logger defines the minimal logger interface for middleware.
type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

// LoggingMiddleware logs incoming requests.
//
//	@summary	Logging middleware
//	@description	Logs method, path, status, latency, request ID, and errors.
func LoggingMiddleware(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		requestID := c.GetHeader("X-Request-ID")

		logFields := []interface{}{
			"method", method,
			"path", path,
			"status", status,
			"latency", latency,
			"request_id", requestID,
		}

		if len(c.Errors) > 0 {
			logger.Error(append([]interface{}{"errors", c.Errors.String()}, logFields...)...)
		} else {
			logger.Info(logFields...)
		}
	}
}
