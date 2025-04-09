package routes

import (
	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/domain/ai/similarity"

	"github.com/go-chi/chi/v5"
)

// RegisterSimilarityRoutes registers the similarity routes
func RegisterSimilarityRoutes(r chi.Router, similaritySvc *similarity.Service) {
	// Create similarity handler
	handler := handlers.NewSimilarityHandler(similaritySvc)

	// Register routes
	r.Route("/api/similarity", func(r chi.Router) {
		r.Post("/index", handler.IndexMessage)
		r.Post("/search", handler.FindSimilarMessages)
		r.Delete("/conversations/{conversationID}", handler.DeleteConversationMessages)
	})
}
