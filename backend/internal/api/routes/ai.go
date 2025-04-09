package routes

import (
	"github.com/go-chi/chi/v5"
	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/api/middleware"
	"go-crypto-bot-clean/backend/internal/domain/ai/service"
)

// RegisterAIRoutesWithChi registers AI-related routes with Chi router
func RegisterAIRoutesWithChi(
	r chi.Router,
	aiSvc service.AIService,
) {
	// Create AI context middleware
	aiContextMiddleware := middleware.AIContextMiddleware()

	// Register routes
	r.Route("/ai", func(r chi.Router) {
		r.Use(aiContextMiddleware)

		r.Post("/chat", handlers.ChatHandler(aiSvc))
		r.Post("/function", handlers.FunctionHandler(aiSvc))
	})
}
