package handler

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type AIHandler struct {
	useCase *usecase.AIUsecase
	logger  *zerolog.Logger
}

func NewAIHandler(useCase *usecase.AIUsecase, logger *zerolog.Logger) *AIHandler {
	return &AIHandler{
		useCase: useCase,
		logger:  logger,
	}
}

func (h *AIHandler) RegisterRoutes(r chi.Router) {
	r.Route("/ai", func(r chi.Router) {
		// Add AI-related routes here
		// Example: r.Post("/analyze", h.AnalyzeMarket)
	})
}
