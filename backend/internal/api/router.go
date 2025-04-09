// Package api sets up the HTTP API routing.
package api

import (
	"go-crypto-bot-clean/backend/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin engine, middleware chain, and API routes.
func SetupRouter(deps *Dependencies) *gin.Engine {
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Skip other middleware for now to simplify testing
	// router.Use(
	// 	middleware.RecoveryMiddleware(deps.Logger),
	// 	middleware.LoggingMiddleware(deps.Logger),
	// 	middleware.AuthMiddleware(deps.ValidAPIKeys),
	// 	middleware.RateLimiterMiddleware(
	// 		deps.RateLimit.Rate,
	// 		deps.RateLimit.Capacity,
	// 		func(c *gin.Context) string { return c.ClientIP() }, // Use client IP as rate limit key
	// 		deps.Logger,
	// 	),
	// )

	// Public endpoints
	router.GET("/health", deps.HealthHandler.HealthCheck)

	// Authentication endpoints
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", deps.AuthHandler.Login)
		authGroup.POST("/logout", deps.AuthHandler.Logout)

		// Protected auth endpoints
		authProtected := authGroup.Group("")
		if deps.Config.Auth.Enabled {
			authProtected.Use(middleware.JWTAuthMiddleware(deps.Config.Auth.JWTSecret))
		}
		authProtected.GET("/me", deps.AuthHandler.GetCurrentUser)
	}

	// Versioned API group
	apiV1 := router.Group("/api/v1")

	// Apply authentication middleware to protected endpoints if enabled
	if deps.Config.Auth.Enabled {
		// Skip authentication for public endpoints
		apiV1.GET("/health", deps.HealthHandler.HealthCheck) // Duplicate for backward compatibility

		// Apply JWT authentication to all other endpoints
		apiV1.Use(middleware.JWTAuthMiddleware(deps.Config.Auth.JWTSecret))
	}

	// API endpoints
	{
		// Status endpoints
		apiV1.GET("/status", deps.StatusHandler.GetStatus)
		apiV1.POST("/status/start", deps.StatusHandler.StartProcesses)
		apiV1.POST("/status/stop", deps.StatusHandler.StopProcesses)

		// Portfolio endpoints
		apiV1.GET("/portfolio", deps.PortfolioHandler.GetPortfolioSummary)
		apiV1.GET("/portfolio/active", deps.PortfolioHandler.GetActiveTrades)
		apiV1.GET("/portfolio/performance", deps.PortfolioHandler.GetPerformanceMetrics)
		apiV1.GET("/portfolio/value", deps.PortfolioHandler.GetTotalValue)

		// Trade endpoints
		apiV1.GET("/trade/history", deps.TradeHandler.GetTradeHistory)
		apiV1.POST("/trade/buy", deps.TradeHandler.ExecuteTrade)
		apiV1.POST("/trade/sell", deps.TradeHandler.SellCoin)
		apiV1.GET("/trade/status/:id", deps.TradeHandler.GetTradeStatus)

		// Backtest endpoints
		apiV1.POST("/backtest/run", deps.BacktestHandler.RunBacktest)
		apiV1.GET("/backtest/results/:id", deps.BacktestHandler.GetBacktestResults)
		apiV1.GET("/backtest/results", deps.BacktestHandler.ListBacktestResults)

		// NewCoin endpoints
		apiV1.GET("/newcoins", deps.NewCoinHandler.GetDetectedCoins)
		apiV1.POST("/newcoins/process", deps.NewCoinHandler.ProcessNewCoins)
		apiV1.POST("/newcoins/detect", deps.NewCoinHandler.DetectNewCoins)
		apiV1.POST("/newcoins/by-date", deps.NewCoinHandler.GetCoinsByDate)
		apiV1.POST("/newcoins/by-date-range", deps.NewCoinHandler.GetCoinsByDateRange)
		// New endpoints for upcoming coins
		apiV1.GET("/newcoins/upcoming", deps.NewCoinHandler.GetUpcomingCoins)
		apiV1.POST("/newcoins/upcoming/by-date", deps.NewCoinHandler.GetUpcomingCoinsByDate)
		apiV1.GET("/newcoins/upcoming/today-and-tomorrow", deps.NewCoinHandler.GetUpcomingCoinsForTodayAndTomorrow)
		// New endpoints for tradable coins
		apiV1.GET("/newcoins/tradable", deps.CoinHandler.ListTradableCoins)
		apiV1.GET("/newcoins/tradable/today", deps.CoinHandler.ListTradableCoinsToday)

		// Config endpoints
		apiV1.GET("/config", deps.ConfigHandler.GetCurrentConfig)
		apiV1.PUT("/config", deps.ConfigHandler.UpdateConfig)
		apiV1.GET("/config/defaults", deps.ConfigHandler.GetDefaultConfig)

		// Account endpoints
		apiV1.GET("/account/details", deps.EnhancedAccountHandler.GetAccountDetails)
		apiV1.GET("/account/validate-keys", deps.EnhancedAccountHandler.ValidateAPIKeys)
		apiV1.GET("/account/listen-key", deps.EnhancedAccountHandler.GetListenKey)
		apiV1.PUT("/account/listen-key/renew", deps.EnhancedAccountHandler.RenewListenKey)
		apiV1.DELETE("/account/listen-key/close", deps.EnhancedAccountHandler.CloseListenKey)

		// Analytics endpoints
		apiV1.GET("/analytics", deps.AnalyticsHandler.GetTradeAnalytics)
		apiV1.GET("/analytics/trades", deps.AnalyticsHandler.GetAllTradePerformance)
		apiV1.GET("/analytics/trades/:id", deps.AnalyticsHandler.GetTradePerformance)
		apiV1.GET("/analytics/winrate", deps.AnalyticsHandler.GetWinRate)
		apiV1.GET("/analytics/balance-history", deps.AnalyticsHandler.GetBalanceHistory)
		apiV1.GET("/analytics/by-symbol", deps.AnalyticsHandler.GetPerformanceBySymbol)
		apiV1.GET("/analytics/by-reason", deps.AnalyticsHandler.GetPerformanceByReason)
		apiV1.GET("/analytics/by-strategy", deps.AnalyticsHandler.GetPerformanceByStrategy)
	}

	// WebSocket endpoint
	router.GET("/ws", deps.WebSocketHandler.ServeWSGin)

	return router
}

// Placeholder middleware functions to be implemented later

// RecoveryMiddleware recovers from panics.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement recovery logic
		c.Next()
	}
}

// LoggingMiddleware logs requests.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement logging logic
		c.Next()
	}
}

// AuthMiddleware authenticates requests.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement authentication logic
		c.Next()
	}
}

// RateLimiterMiddleware limits request rate.
func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement rate limiting logic
		c.Next()
	}
}
