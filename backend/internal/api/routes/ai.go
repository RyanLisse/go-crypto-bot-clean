package routes

import (
	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/api/middleware"
	"go-crypto-bot-clean/backend/internal/domain/ai/service"
	"go-crypto-bot-clean/backend/internal/domain/portfolio"
	"go-crypto-bot-clean/backend/internal/domain/trade"

	"github.com/go-chi/chi/v5"
)

// RegisterAIRoutesWithChi registers AI-related routes with Chi router
func RegisterAIRoutesWithChi(
	r chi.Router,
	aiSvc service.AIService,
	portfolioSvc portfolio.Service,
	tradeSvc trade.Service,
) {
	// Create AI context middleware
	aiContextMiddleware := middleware.AIContextMiddleware(portfolioSvc, tradeSvc)

	// Register routes
	r.Route("/ai", func(r chi.Router) {
		r.Use(aiContextMiddleware)

		// Chat and function endpoints
		r.Post("/chat", handlers.ChatHandler(aiSvc))
		r.Post("/function", handlers.FunctionHandler(aiSvc))
		r.Post("/insights", handlers.AIInsightsHandler(aiSvc))

		// Risk management endpoints
		r.Route("/risk", func(r chi.Router) {
			r.Post("/guardrails", handlers.ApplyRiskGuardrailsHandler(aiSvc))
			r.Post("/confirmation", handlers.CreateTradeConfirmationHandler(aiSvc))
			r.Post("/confirm", handlers.ConfirmTradeHandler(aiSvc))
			r.Post("/pending", handlers.ListPendingTradeConfirmationsHandler(aiSvc))
		})
	})
}
