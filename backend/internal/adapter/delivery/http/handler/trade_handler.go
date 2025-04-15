package handler

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type TradeHandler struct {
	useCase usecase.TradeUseCase
	logger  *zerolog.Logger
}

func NewTradeHandler(useCase usecase.TradeUseCase, logger *zerolog.Logger) *TradeHandler {
	return &TradeHandler{
		useCase: useCase,
		logger:  logger,
	}
}

func (h *TradeHandler) RegisterRoutes(r chi.Router) {
	r.Route("/trade", func(r chi.Router) {
		// Add trading routes here
		// Example: r.Post("/orders", h.PlaceOrder)
	})
}
