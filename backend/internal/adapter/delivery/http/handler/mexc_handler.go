package handler

import (
	"net/http"
	"strconv"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/response"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// MEXCHandler handles MEXC API-related endpoints
type MEXCHandler struct {
	mexcClient port.MEXCClient
	logger     *zerolog.Logger
}

// NewMEXCHandler creates a new MEXCHandler
func NewMEXCHandler(mexcClient port.MEXCClient, logger *zerolog.Logger) *MEXCHandler {
	return &MEXCHandler{
		mexcClient: mexcClient,
		logger:     logger,
	}
}

// RegisterRoutes registers the MEXC API routes
func (h *MEXCHandler) RegisterRoutes(r chi.Router) {
	r.Route("/mexc", func(r chi.Router) {
		// Account endpoints
		r.Get("/account", h.GetAccount)

		// Market data endpoints
		r.Get("/ticker/{symbol}", h.GetTicker)
		r.Get("/orderbook/{symbol}", h.GetOrderBook)
		r.Get("/klines/{symbol}/{interval}", h.GetKlines)
		r.Get("/exchange-info", h.GetExchangeInfo)

		// Symbol endpoints
		r.Get("/symbol/{symbol}", h.GetSymbolInfo)

		// New listings endpoint
		r.Get("/new-listings", h.GetNewListings)
	})
}

// GetAccount returns the user's MEXC account information
func (h *MEXCHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug().Msg("Getting MEXC account information")

	// Get account information from MEXC
	account, err := h.mexcClient.GetAccount(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get MEXC account information")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("internal_error", "Failed to get MEXC account information"))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(account))
}

// GetTicker returns the ticker for a specific symbol
func (h *MEXCHandler) GetTicker(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	if symbol == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.Error("invalid_symbol", "Symbol is required"))
		return
	}

	h.logger.Debug().Str("symbol", symbol).Msg("Getting MEXC ticker")

	// Get ticker from MEXC
	ticker, err := h.mexcClient.GetMarketData(r.Context(), symbol)
	if err != nil {
		h.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get MEXC ticker")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("internal_error", "Failed to get MEXC ticker"))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(ticker))
}

// GetOrderBook returns the order book for a specific symbol
func (h *MEXCHandler) GetOrderBook(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	if symbol == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.Error("invalid_symbol", "Symbol is required"))
		return
	}

	// Parse depth parameter
	depthStr := r.URL.Query().Get("depth")
	depth := 10 // Default depth
	if depthStr != "" {
		parsedDepth, err := strconv.Atoi(depthStr)
		if err != nil || parsedDepth <= 0 {
			response.WriteJSON(w, http.StatusBadRequest, response.Error("invalid_depth", "Depth must be a positive integer"))
			return
		}
		depth = parsedDepth
	}

	h.logger.Debug().Str("symbol", symbol).Int("depth", depth).Msg("Getting MEXC order book")

	// Get order book from MEXC
	orderBook, err := h.mexcClient.GetOrderBook(r.Context(), symbol, depth)
	if err != nil {
		h.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get MEXC order book")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("internal_error", "Failed to get MEXC order book"))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(orderBook))
}

// GetKlines returns the klines for a specific symbol and interval
func (h *MEXCHandler) GetKlines(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	intervalStr := chi.URLParam(r, "interval")

	if symbol == "" || intervalStr == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.Error("invalid_params", "Symbol and interval are required"))
		return
	}

	// Parse limit parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			response.WriteJSON(w, http.StatusBadRequest, response.Error("invalid_limit", "Limit must be a positive integer"))
			return
		}
		limit = parsedLimit
	}

	// Convert interval string to model.KlineInterval
	interval := model.KlineInterval(intervalStr)

	h.logger.Debug().Str("symbol", symbol).Str("interval", string(interval)).Int("limit", limit).Msg("Getting MEXC klines")

	// Get klines from MEXC
	klines, err := h.mexcClient.GetKlines(r.Context(), symbol, interval, limit)
	if err != nil {
		h.logger.Error().Err(err).Str("symbol", symbol).Str("interval", string(interval)).Msg("Failed to get MEXC klines")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("internal_error", "Failed to get MEXC klines"))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(klines))
}

// GetExchangeInfo returns information about all symbols on the exchange
func (h *MEXCHandler) GetExchangeInfo(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug().Msg("Getting MEXC exchange info")

	// Get exchange info from MEXC
	exchangeInfo, err := h.mexcClient.GetExchangeInfo(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get MEXC exchange info")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("internal_error", "Failed to get MEXC exchange info"))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(exchangeInfo))
}

// GetSymbolInfo returns detailed information about a trading symbol
func (h *MEXCHandler) GetSymbolInfo(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	if symbol == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.Error("invalid_symbol", "Symbol is required"))
		return
	}

	h.logger.Debug().Str("symbol", symbol).Msg("Getting MEXC symbol info")

	// Get symbol info from MEXC
	symbolInfo, err := h.mexcClient.GetSymbolInfo(r.Context(), symbol)
	if err != nil {
		h.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get MEXC symbol info")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("internal_error", "Failed to get MEXC symbol info"))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(symbolInfo))
}

// GetNewListings returns information about newly listed coins
func (h *MEXCHandler) GetNewListings(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug().Msg("Getting MEXC new listings")

	// Get new listings from MEXC
	newListings, err := h.mexcClient.GetNewListings(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get MEXC new listings")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("internal_error", "Failed to get MEXC new listings"))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(newListings))
}
