package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/api/huma"
	"go-crypto-bot-clean/backend/internal/api/middleware/cors"
	"go-crypto-bot-clean/backend/internal/api/websocket"
)

// SetupChiRouter initializes the Chi router with Huma for OpenAPI documentation.
func SetupChiRouter(deps *HumaDependencies) http.Handler {
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Middleware())

	// Setup Huma for OpenAPI documentation
	humaConfig := huma.DefaultConfig()
	humaAPI := huma.SetupHuma(r, humaConfig)
	_ = humaAPI // Use the API to avoid unused variable warning

	// Health check endpoint
	r.Get("/health", adaptGinHandler(deps.HealthHandler.HealthCheck))

	// Versioned API group
	r.Route("/api/v1", func(r chi.Router) {
		// Status endpoints
		r.Get("/status", adaptGinHandler(deps.StatusHandler.GetStatus))
		r.Post("/status/start", adaptGinHandler(deps.StatusHandler.StartProcesses))
		r.Post("/status/stop", adaptGinHandler(deps.StatusHandler.StopProcesses))

		// Portfolio endpoints
		r.Get("/portfolio", adaptGinHandler(deps.PortfolioHandler.GetPortfolioSummary))
		r.Get("/portfolio/active", adaptGinHandler(deps.PortfolioHandler.GetActiveTrades))
		r.Get("/portfolio/performance", adaptGinHandler(deps.PortfolioHandler.GetPerformanceMetrics))
		r.Get("/portfolio/value", adaptGinHandler(deps.PortfolioHandler.GetTotalValue))

		// Trade endpoints
		r.Get("/trade/history", adaptGinHandler(deps.TradeHandler.GetTradeHistory))
		r.Post("/trade/buy", adaptGinHandler(deps.TradeHandler.ExecuteTrade))
		r.Post("/trade/sell", adaptGinHandler(deps.TradeHandler.SellCoin))
		r.Get("/trade/status/{id}", adaptGinHandler(deps.TradeHandler.GetTradeStatus))

		// NewCoin endpoints
		r.Get("/newcoins", adaptGinHandler(deps.NewCoinHandler.GetDetectedCoins))
		r.Post("/newcoins/process", adaptGinHandler(deps.NewCoinHandler.ProcessNewCoins))
		r.Post("/newcoins/detect", adaptGinHandler(deps.NewCoinHandler.DetectNewCoins))
		r.Post("/newcoins/by-date", adaptGinHandler(deps.NewCoinHandler.GetCoinsByDate))
		r.Post("/newcoins/by-date-range", adaptGinHandler(deps.NewCoinHandler.GetCoinsByDateRange))

		// Config endpoints
		r.Get("/config", adaptGinHandler(deps.ConfigHandler.GetCurrentConfig))
		r.Put("/config", adaptGinHandler(deps.ConfigHandler.UpdateConfig))
		r.Get("/config/defaults", adaptGinHandler(deps.ConfigHandler.GetDefaultConfig))

		// Analytics endpoints
		r.Get("/analytics", adaptGinHandler(deps.AnalyticsHandler.GetTradeAnalytics))
		r.Get("/analytics/trades", adaptGinHandler(deps.AnalyticsHandler.GetAllTradePerformance))
		r.Get("/analytics/trades/{id}", adaptGinHandler(deps.AnalyticsHandler.GetTradePerformance))
		r.Get("/analytics/winrate", adaptGinHandler(deps.AnalyticsHandler.GetWinRate))
		r.Get("/analytics/balance-history", adaptGinHandler(deps.AnalyticsHandler.GetBalanceHistory))
		r.Get("/analytics/by-symbol", adaptGinHandler(deps.AnalyticsHandler.GetPerformanceBySymbol))
		r.Get("/analytics/by-reason", adaptGinHandler(deps.AnalyticsHandler.GetPerformanceByReason))
		r.Get("/analytics/by-strategy", adaptGinHandler(deps.AnalyticsHandler.GetPerformanceByStrategy))
	})

	// WebSocket endpoint
	r.Get("/ws", adaptGinHandler(deps.WebSocketHandler.ServeWSGin))

	return r
}

// adaptGinHandler adapts a Gin handler to an http.HandlerFunc.
func adaptGinHandler(ginHandler func(*gin.Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a test context
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = r

		// Call the Gin handler
		ginHandler(ginCtx)
	}
}

// responseWriterAdapter adapts http.ResponseWriter to gin.ResponseWriter.
type responseWriterAdapter struct {
	http.ResponseWriter
	size   int
	status int
}

func (w *responseWriterAdapter) Status() int {
	return w.status
}

func (w *responseWriterAdapter) Size() int {
	return w.size
}

func (w *responseWriterAdapter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriterAdapter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func (w *responseWriterAdapter) WriteString(s string) (int, error) {
	n, err := w.ResponseWriter.Write([]byte(s))
	w.size += n
	return n, err
}

func (w *responseWriterAdapter) Written() bool {
	return w.size > 0
}

// HumaDependencies contains all the dependencies for the Huma API.
type HumaDependencies struct {
	HealthHandler    *handlers.HealthHandler
	StatusHandler    *handlers.StatusHandler
	PortfolioHandler *handlers.PortfolioHandler
	TradeHandler     *handlers.TradeHandler
	NewCoinHandler   *handlers.NewCoinsHandler
	ConfigHandler    *handlers.ConfigHandler
	WebSocketHandler *websocket.Handler
	AnalyticsHandler *handlers.AnalyticsHandler
}
