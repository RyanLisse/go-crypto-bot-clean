package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo/crypto-bot/internal/apperror"
	"github.com/rs/zerolog"
)

// SetupRouter configures and returns the main HTTP router
func SetupRouter(logger zerolog.Logger) *gin.Engine {
	// Set Gin mode based on environment
	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create a new router
	router := gin.New()

	// Add middleware
	router.Use(
		gin.Recovery(),
		corsMiddleware(),
		loggerMiddleware(logger),
		errorMiddleware(),
	)

	// Setup health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Setup API routes
	// API group will be used by handlers to register their routes
	_ = router.Group("/api/v1")

	// Add more route registrations here when handlers are implemented

	return router
}

// corsMiddleware sets up CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// loggerMiddleware sets up request logging
func loggerMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log information after request is processed
		end := time.Now()
		latency := end.Sub(start)

		if raw != "" {
			path = path + "?" + raw
		}

		status := c.Writer.Status()
		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		logEvent := logger.Info()
		if status >= 400 {
			logEvent = logger.Error()
		}

		logEvent.
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Dur("latency", latency).
			Str("ip", ip).
			Str("user_agent", userAgent).
			Msg("HTTP Request")
	}
}

// errorMiddleware handles and standardizes error responses
func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Check if it's an AppError
			var appErr *apperror.AppError
			if apperror.As(err, &appErr) {
				c.JSON(appErr.StatusCode, appErr.ToResponse())
				return
			}

			// If not an AppError, return a generic 500 error
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "An internal server error occurred",
				},
			})
		}
	}
}
