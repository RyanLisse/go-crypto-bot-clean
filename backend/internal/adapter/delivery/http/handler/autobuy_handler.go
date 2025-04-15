package handler

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type AutoBuyHandler struct {
	useCase usecase.AutoBuyUseCase
	logger  *zerolog.Logger
}

func NewAutoBuyHandler(useCase usecase.AutoBuyUseCase, logger *zerolog.Logger) *AutoBuyHandler {
	return &AutoBuyHandler{
		useCase: useCase,
		logger:  logger,
	}
}

func (h *AutoBuyHandler) RegisterRoutes(r chi.Router) {
	r.Route("/autobuy", func(r chi.Router) {
		// Add auto-buy related routes here
		// Example: r.Post("/rules", h.CreateRule)
	})
}
