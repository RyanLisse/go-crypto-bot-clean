package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo/crypto-bot/internal/adapter/http/middleware"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

// SetupRouter configures and returns the main HTTP router and API group
func SetupRouter(logger zerolog.Logger) (*gin.Engine, *gin.RouterGroup) {
	// Set Gin mode based on environment
	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create a new router
	router := gin.New()

	// Create rate limiters
	ipRateLimiter := middleware.NewIPRateLimiter(rate.Limit(1), 60, &logger) // 60 requests per minute
	dailyRateLimiter := middleware.NewDailyRateLimiter(1000, &logger)        // 1000 requests per day

	// Add middleware
	router.Use(
		gin.Recovery(),
		corsMiddleware(),
		loggerMiddleware(logger),
		middleware.RateLimiterMiddleware(ipRateLimiter),
		middleware.DailyRateLimiterMiddleware(dailyRateLimiter),
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
	apiV1 := router.Group("/api/v1")

	return router, apiV1
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
