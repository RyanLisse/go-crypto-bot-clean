package handler

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type RiskHandler struct {
	useCase usecase.RiskUseCase
	logger  *zerolog.Logger
}

func NewRiskHandler(useCase usecase.RiskUseCase, logger *zerolog.Logger) *RiskHandler {
	return &RiskHandler{
		useCase: useCase,
		logger:  logger,
	}
}

func (h *RiskHandler) RegisterRoutes(r chi.Router) {
	r.Route("/risk", func(r chi.Router) {
		// Add risk management routes here
		// Example: r.Get("/analysis", h.GetRiskAnalysis)
	})
}
