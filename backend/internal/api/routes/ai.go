package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/ryanlisse/go-crypto-bot/internal/api/handlers"
	"github.com/ryanlisse/go-crypto-bot/internal/api/middleware"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/ai/service"
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
