package api

import (
	"net/http"

	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/api/middleware"

	"github.com/gin-gonic/gin"
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
	router.Use(gin.WrapH(middleware.RecoveryMiddleware(middleware.RecoveryOptions{
		Logger:           logger,
		EnableStackTrace: true,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))))

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
