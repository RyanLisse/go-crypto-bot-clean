package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ryanlisse/go-crypto-bot/internal/api/handlers"
	"github.com/ryanlisse/go-crypto-bot/internal/api/middleware"
	"go.uber.org/zap"
)

// ZapLoggerAdapter adapts zap.Logger to middleware.Logger
type ZapLoggerAdapter struct {
	logger *zap.Logger
}

// Info implements middleware.Logger.Info
func (a *ZapLoggerAdapter) Info(args ...interface{}) {
	a.logger.Sugar().Info(args...)
}

// Error implements middleware.Logger.Error
func (a *ZapLoggerAdapter) Error(args ...interface{}) {
	a.logger.Sugar().Error(args...)
}

// SetupRoutes configures the API routes
func SetupRoutes(
	router *gin.Engine,
	statusHandler *handlers.StatusHandler,
	reportHandler *handlers.ReportHandler,
	wsHandler *handlers.WebSocketHandler,
	accountHandler *handlers.AccountHandler,
	logger *zap.Logger,
) {
	// Create logger adapter
	loggerAdapter := &ZapLoggerAdapter{logger: logger}

	// Middleware
	router.Use(middleware.LoggingMiddleware(loggerAdapter))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RecoveryMiddleware(loggerAdapter))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Status routes
		v1.GET("/status", statusHandler.GetStatus)
		v1.POST("/status/start", statusHandler.StartProcesses)
		v1.POST("/status/stop", statusHandler.StopProcesses)

		// Report routes
		reportHandler.RegisterRoutes(v1)

		// WebSocket routes
		wsHandler.RegisterRoutes(v1)

		// Account routes
		account := v1.Group("/account")
		{
			account.GET("", accountHandler.GetAccount)
			account.GET("/balance", accountHandler.GetBalances)
			account.GET("/wallet", accountHandler.GetWallet)
			account.GET("/balance-summary", accountHandler.GetBalanceSummary)
			account.GET("/validate-keys", accountHandler.ValidateAPIKeys)
			account.POST("/sync", accountHandler.SyncWithExchange)
		}
	}

	// Static routes for testing
	router.Static("/test", "./static")
}
