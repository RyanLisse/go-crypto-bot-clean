package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/util"
)

// WalletHandler handles wallet-related endpoints
//
type WalletHandler struct {
	walletService usecase.WalletService
	logger        *zerolog.Logger
}

// NewWalletHandler creates a new WalletHandler
func NewWalletHandler(walletService usecase.WalletService, logger *zerolog.Logger) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
		logger:        logger,
	}
}

// RegisterRoutes registers wallet-related routes
func (c *WalletHandler) RegisterRoutes(r chi.Router) {
	r.Route("/wallet", func(r chi.Router) {
		r.Get("/real", c.GetRealWallet)
	})
}

// GetRealWallet handles GET /wallet/real
func (c *WalletHandler) GetRealWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wallet, err := c.walletService.GetRealAccountData(ctx)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to fetch real wallet data")
		util.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch real wallet data"})
		return
	}
	util.WriteJSONResponse(w, http.StatusOK, wallet)
}
