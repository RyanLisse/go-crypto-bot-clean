// Package api sets up the HTTP API routing.
package api

import (
	"net/http"

	"go-crypto-bot-clean/backend/internal/api/huma"
	"go-crypto-bot-clean/backend/internal/api/middleware"
	"go-crypto-bot-clean/backend/internal/api/middleware/cors"
	"go-crypto-bot-clean/backend/internal/api/service"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// SetupConsolidatedRouter initializes the Chi router, middleware chain, and API routes.
func SetupConsolidatedRouter(deps *Dependencies) http.Handler {
	r := chi.NewRouter()

	// Add standard middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Middleware())

	// Add custom middleware if needed
	if deps.Config != nil {
		r.Use(RecoveryMiddlewareForChi(deps.logger))
	}

	// Setup Huma for OpenAPI documentation if services are available
	if deps.StrategyService != nil &&
		deps.AuthService != nil && deps.UserService != nil {
		humaConfig := huma.DefaultConfig()
		// Create a service provider with the available services
		serviceProvider := &service.Provider{}

		// Set the services if available
		if deps.StrategyService != nil {
			serviceProvider.StrategyService = deps.StrategyService
		}
		// Pass the internal/auth.AuthProvider directly to NewAuthService
		if deps.AuthService != nil && deps.UserRepository != nil {
			serviceProvider.AuthService = service.NewAuthService(deps.AuthService, deps.UserRepository)
		} else {
			if deps.AuthService == nil {
				deps.logger.Warn("AuthService is nil", zap.String("context", "cannot initialize service.AuthService for Huma"))
			}
			if deps.UserRepository == nil {
				deps.logger.Warn("UserRepository is nil", zap.String("context", "cannot initialize service.AuthService for Huma"))
			}
		}

		if deps.UserService != nil {
			serviceProvider.UserService = deps.UserService
		}
		humaAPI := huma.SetupHuma(r, humaConfig, serviceProvider)
		_ = humaAPI // Use the API to avoid unused variable warning
	}

	// Public endpoints
	r.Get("/health", deps.HealthHandler.HealthCheck)

	// Authentication endpoints
	r.Route("/auth", func(r chi.Router) {
		// Login/Logout routes removed

		// Protected auth endpoints
		if deps.Config != nil && deps.Config.Auth.Enabled {
			r.Group(func(r chi.Router) {
				r.Use(middleware.JWTAuthMiddleware(deps.Config.Auth.JWTSecret))
				r.Get("/me", deps.AuthHandler.GetCurrentUser)
			})
		} else {
			r.Get("/me", deps.AuthHandler.GetCurrentUser)
		}
	})

	// Versioned API group
	r.Route("/api/v1", func(r chi.Router) {
		// Apply authentication middleware to protected endpoints if enabled
		if deps.Config != nil && deps.Config.Auth.Enabled {
			// Skip authentication for public endpoints
			r.Get("/health", deps.HealthHandler.HealthCheck) // Duplicate for backward compatibility

			// Apply JWT authentication to all other endpoints
			r.Group(func(r chi.Router) {
				r.Use(middleware.JWTAuthMiddleware(deps.Config.Auth.JWTSecret))

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

				// Backtest endpoints
				if deps.BacktestHandler != nil {
					r.Post("/backtest/run", deps.BacktestHandler.RunBacktest)
					r.Get("/backtest/results/{id}", deps.BacktestHandler.GetBacktestResults)
					r.Get("/backtest/results", deps.BacktestHandler.ListBacktestResults)
				}

				// NewCoin endpoints
				r.Get("/newcoins", deps.NewCoinHandler.GetDetectedCoins)
				r.Post("/newcoins/process", deps.NewCoinHandler.ProcessNewCoins)
				r.Post("/newcoins/detect", deps.NewCoinHandler.DetectNewCoins)
				r.Post("/newcoins/by-date", deps.NewCoinHandler.GetCoinsByDate)
				r.Post("/newcoins/by-date-range", deps.NewCoinHandler.GetCoinsByDateRange)

				// New endpoints for upcoming coins
				r.Get("/newcoins/upcoming", deps.NewCoinHandler.GetUpcomingCoins)
				r.Post("/newcoins/upcoming/by-date", deps.NewCoinHandler.GetUpcomingCoinsByDate)
				r.Get("/newcoins/upcoming/today-and-tomorrow", deps.NewCoinHandler.GetUpcomingCoinsForTodayAndTomorrow)

				// New endpoints for tradable coins
				if deps.CoinHandler != nil {
					r.Get("/newcoins/tradable", deps.CoinHandler.ListTradableCoins)
					r.Get("/newcoins/tradable/today", deps.CoinHandler.ListTradableCoinsToday)
				}

				// Config endpoints
				r.Get("/config", deps.ConfigHandler.GetCurrentConfig)
				r.Put("/config", deps.ConfigHandler.UpdateConfig)
				r.Get("/config/defaults", deps.ConfigHandler.GetDefaultConfig)

				// Account endpoints
				if deps.EnhancedAccountHandler != nil {
					r.Get("/account/details", deps.EnhancedAccountHandler.GetAccountDetails)
					r.Get("/account/validate-keys", deps.EnhancedAccountHandler.ValidateAPIKeys)
					r.Get("/account/listen-key", deps.EnhancedAccountHandler.GetListenKey)
					r.Put("/account/listen-key/renew", deps.EnhancedAccountHandler.RenewListenKey)
					r.Delete("/account/listen-key/close", deps.EnhancedAccountHandler.CloseListenKey)
				}

				// Analytics endpoints
				if deps.AnalyticsHandler != nil {
					r.Get("/analytics", deps.AnalyticsHandler.GetTradeAnalytics)
					r.Get("/analytics/trades", deps.AnalyticsHandler.GetAllTradePerformance)
					r.Get("/analytics/trades/{id}", deps.AnalyticsHandler.GetTradePerformance)
					r.Get("/analytics/winrate", deps.AnalyticsHandler.GetWinRate)
					r.Get("/analytics/balance-history", deps.AnalyticsHandler.GetBalanceHistory)
					r.Get("/analytics/by-symbol", deps.AnalyticsHandler.GetPerformanceBySymbol)
					r.Get("/analytics/by-reason", deps.AnalyticsHandler.GetPerformanceByReason)
					r.Get("/analytics/by-strategy", deps.AnalyticsHandler.GetPerformanceByStrategy)
				}
			})
		} else {
			// If auth is disabled, register all routes without authentication
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

			// Backtest endpoints
			if deps.BacktestHandler != nil {
				r.Post("/backtest/run", deps.BacktestHandler.RunBacktest)
				r.Get("/backtest/results/{id}", deps.BacktestHandler.GetBacktestResults)
				r.Get("/backtest/results", deps.BacktestHandler.ListBacktestResults)
			}

			// NewCoin endpoints
			r.Get("/newcoins", deps.NewCoinHandler.GetDetectedCoins)
			r.Post("/newcoins/process", deps.NewCoinHandler.ProcessNewCoins)
			r.Post("/newcoins/detect", deps.NewCoinHandler.DetectNewCoins)
			r.Post("/newcoins/by-date", deps.NewCoinHandler.GetCoinsByDate)
			r.Post("/newcoins/by-date-range", deps.NewCoinHandler.GetCoinsByDateRange)

			// New endpoints for upcoming coins
			r.Get("/newcoins/upcoming", deps.NewCoinHandler.GetUpcomingCoins)
			r.Post("/newcoins/upcoming/by-date", deps.NewCoinHandler.GetUpcomingCoinsByDate)
			r.Get("/newcoins/upcoming/today-and-tomorrow", deps.NewCoinHandler.GetUpcomingCoinsForTodayAndTomorrow)

			// New endpoints for tradable coins
			if deps.CoinHandler != nil {
				r.Get("/newcoins/tradable", deps.CoinHandler.ListTradableCoins)
				r.Get("/newcoins/tradable/today", deps.CoinHandler.ListTradableCoinsToday)
			}

			// Config endpoints
			r.Get("/config", deps.ConfigHandler.GetCurrentConfig)
			r.Put("/config", deps.ConfigHandler.UpdateConfig)
			r.Get("/config/defaults", deps.ConfigHandler.GetDefaultConfig)

			// Account endpoints
			if deps.EnhancedAccountHandler != nil {
				r.Get("/account/details", deps.EnhancedAccountHandler.GetAccountDetails)
				r.Get("/account/validate-keys", deps.EnhancedAccountHandler.ValidateAPIKeys)
				r.Get("/account/listen-key", deps.EnhancedAccountHandler.GetListenKey)
				r.Put("/account/listen-key/renew", deps.EnhancedAccountHandler.RenewListenKey)
				r.Delete("/account/listen-key/close", deps.EnhancedAccountHandler.CloseListenKey)
			}

			// Analytics endpoints
			if deps.AnalyticsHandler != nil {
				r.Get("/analytics", deps.AnalyticsHandler.GetTradeAnalytics)
				r.Get("/analytics/trades", deps.AnalyticsHandler.GetAllTradePerformance)
				r.Get("/analytics/trades/{id}", deps.AnalyticsHandler.GetTradePerformance)
				r.Get("/analytics/winrate", deps.AnalyticsHandler.GetWinRate)
				r.Get("/analytics/balance-history", deps.AnalyticsHandler.GetBalanceHistory)
				r.Get("/analytics/by-symbol", deps.AnalyticsHandler.GetPerformanceBySymbol)
				r.Get("/analytics/by-reason", deps.AnalyticsHandler.GetPerformanceByReason)
				r.Get("/analytics/by-strategy", deps.AnalyticsHandler.GetPerformanceByStrategy)
			}
		}
	})

	// WebSocket endpoint
	r.Get("/ws", deps.WebSocketHandler.ServeWS) // Reverted to original correct method name

	return r
}

// RecoveryMiddlewareForChi creates a Chi middleware for panic recovery.
func RecoveryMiddlewareForChi(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered:", zap.Any("error", err))
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
