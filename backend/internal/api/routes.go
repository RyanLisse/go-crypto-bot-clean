package api

import (
	"net/http"
	"strings"

	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/api/middleware/cors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// SetupRoutes configures the API routes
func SetupRoutes(
	statusHandler *handlers.StatusHandler,
	reportHandler *handlers.ReportHandler,
	wsHandler *handlers.WebSocketHandler,
	accountHandler *handlers.AccountHandler,
	healthHandler *handlers.HealthHandler,
	logger *zap.Logger,
) http.Handler {
	r := chi.NewRouter()

	// Add core middleware
	setupMiddleware(r, logger)

	// Health check route
	r.Get("/health", healthHandler.HealthCheck)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Status routes
		r.Get("/status", statusHandler.GetStatus)
		r.Post("/status/start", statusHandler.StartProcesses)
		r.Post("/status/stop", statusHandler.StopProcesses)

		// Report routes
		r.Route("/report", func(r chi.Router) {
			// TODO: Review these mappings - are they correct?
			r.Get("/", reportHandler.GetReportsByPeriod)     // Mapped from GetReport
			r.Get("/summary", reportHandler.GetLatestReport) // Mapped from GetSummary
			r.Get("/details", reportHandler.GetReportByID)   // Mapped from GetDetails - Needs ID param?
		})

		// WebSocket endpoint
		r.Get("/ws", wsHandler.HandleWebSocket) // Corrected method name

		// Account routes
		r.Route("/account", func(r chi.Router) {
			r.Get("/", accountHandler.GetAccount)
			r.Get("/balance", accountHandler.GetBalances)
			r.Get("/wallet", accountHandler.GetWallet)
			r.Get("/balance-summary", accountHandler.GetBalanceSummary)
			r.Get("/validate-keys", accountHandler.ValidateAPIKeys)
			r.Post("/sync", accountHandler.SyncWithExchange)
		})
	})

	// Static routes for testing
	FileServer(r, "/test", http.Dir("./static"))

	return r
}

// setupMiddleware configures all necessary middleware for the router
func setupMiddleware(r *chi.Mux, _ *zap.Logger) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Middleware())
}

// FileServer conveniently sets up a http.FileServer handler to serve static files from a http.FileSystem
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
