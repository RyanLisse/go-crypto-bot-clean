package api

import (
	"log"
	"net/http"

	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/api/huma"
	"go-crypto-bot-clean/backend/internal/api/service"
	"go-crypto-bot-clean/backend/internal/api/websocket"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// SetupChiRouter initializes the Chi router with conditional Huma integration for OpenAPI documentation.
func SetupChiRouter(deps *HumaDependencies, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()

	// Add core middleware
	setupMiddleware(r, logger) // Pass logger here

	// Setup Huma conditionally based on service availability
	if hasRequiredServices(deps.ServiceProvider) {
		if err := setupHuma(r, deps.ServiceProvider); err != nil {
			log.Printf("Warning: Huma setup failed: %v. API documentation will not be available.", err)
		}
	}

	// Setup all routes
	setupRoutes(r, deps)

	return r
}

// setupHuma configures Huma for OpenAPI documentation if core services are available
func setupHuma(r *chi.Mux, provider *service.Provider) error {
	humaConfig := huma.DefaultConfig()
	_ = huma.SetupHuma(r, humaConfig, provider)
	return nil // Huma's SetupHuma doesn't return an error, so we'll assume success
}

// hasRequiredServices checks if all required services are available for Huma setup
func hasRequiredServices(provider *service.Provider) bool {
	// Check for required services
	return provider != nil &&
		provider.HasBacktestService() &&
		provider.HasStrategyService() &&
		provider.HasAuthService() &&
		provider.HasUserService()
}

// setupRoutes configures all API routes
func setupRoutes(r *chi.Mux, deps *HumaDependencies) {
	// Health check endpoint
	r.Get("/health", deps.HealthHandler.HealthCheck)

	// Versioned API group
	r.Route("/api/v1", func(r chi.Router) {
		// Status endpoints
		r.Get("/status", deps.StatusHandler.GetStatus)
		r.Post("/status/start", deps.StatusHandler.StartProcesses)
		r.Post("/status/stop", deps.StatusHandler.StopProcesses)

		// Portfolio endpoints
		r.Get("/portfolio", deps.PortfolioHandler.GetPortfolioSummary)
		r.Get("/portfolio/active", deps.PortfolioHandler.GetActiveTrades)
		r.Get("/portfolio/performance", deps.PortfolioHandler.GetPerformanceMetrics)
		r.Get("/portfolio/value", deps.PortfolioHandler.GetTotalValue)

		// Trade endpoints
		r.Get("/trade/history", deps.TradeHandler.GetTradeHistory)
		r.Post("/trade/buy", deps.TradeHandler.ExecuteTrade)
		r.Post("/trade/sell", deps.TradeHandler.SellCoin)
		r.Get("/trade/status/{id}", deps.TradeHandler.GetTradeStatus)

		// NewCoin endpoints
		r.Get("/newcoins", deps.NewCoinHandler.GetDetectedCoins)
		r.Post("/newcoins/process", deps.NewCoinHandler.ProcessNewCoins)
		r.Post("/newcoins/detect", deps.NewCoinHandler.DetectNewCoins)
		r.Post("/newcoins/by-date", deps.NewCoinHandler.GetCoinsByDate)
		r.Post("/newcoins/by-date-range", deps.NewCoinHandler.GetCoinsByDateRange)

		// Config endpoints
		r.Get("/config", deps.ConfigHandler.GetCurrentConfig)
		r.Put("/config", deps.ConfigHandler.UpdateConfig)
		r.Get("/config/defaults", deps.ConfigHandler.GetDefaultConfig)

		// Analytics endpoints
		r.Get("/analytics", deps.AnalyticsHandler.GetTradeAnalytics)
		r.Get("/analytics/trades", deps.AnalyticsHandler.GetAllTradePerformance)
		r.Get("/analytics/trades/{id}", deps.AnalyticsHandler.GetTradePerformance)
		r.Get("/analytics/winrate", deps.AnalyticsHandler.GetWinRate)
		r.Get("/analytics/balance-history", deps.AnalyticsHandler.GetBalanceHistory)
		r.Get("/analytics/by-symbol", deps.AnalyticsHandler.GetPerformanceBySymbol)
		r.Get("/analytics/by-reason", deps.AnalyticsHandler.GetPerformanceByReason)
		r.Get("/analytics/by-strategy", deps.AnalyticsHandler.GetPerformanceByStrategy)
	})

	// WebSocket endpoint
	r.Get("/ws", deps.WebSocketHandler.ServeWS)
}

// HumaDependencies contains all the dependencies for the API.
type HumaDependencies struct {
	HealthHandler    *handlers.HealthHandler
	StatusHandler    *handlers.StatusHandler
	PortfolioHandler *handlers.PortfolioHandler
	TradeHandler     *handlers.TradeHandler
	NewCoinHandler   *handlers.NewCoinsHandler
	ConfigHandler    *handlers.ConfigHandler
	WebSocketHandler *websocket.Handler
	AnalyticsHandler *handlers.AnalyticsHandler
	ServiceProvider  *service.Provider
}
