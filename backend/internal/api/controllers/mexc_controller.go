package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"go-crypto-bot-clean/backend/internal/domain/service"
)

// MexcController handles MEXC API endpoints
type MexcController struct {
	exchangeService service.ExchangeService
	logger          *zap.Logger
}

// NewMexcController creates a new MEXC controller
func NewMexcController(exchangeService service.ExchangeService, logger *zap.Logger) *MexcController {
	return &MexcController{
		exchangeService: exchangeService,
		logger:          logger,
	}
}

// RegisterRoutes registers the MEXC API routes
func (c *MexcController) RegisterRoutes(r chi.Router) {
	r.Route("/mexc", func(r chi.Router) {
		r.Get("/ticker/{symbol}", c.GetTicker)
		r.Get("/tickers", c.GetTickers)
		r.Get("/klines/{symbol}", c.GetKlines)
		r.Get("/depth/{symbol}", c.GetOrderBook)
		r.Get("/trades/{symbol}", c.GetRecentTrades)
		r.Get("/account", c.GetAccountInfo)
		r.Get("/orders", c.GetOpenOrders)
		r.Post("/order", c.CreateOrder)
		r.Delete("/order/{orderId}", c.CancelOrder)
	})
}

// GetTicker returns the ticker for a specific symbol
func (c *MexcController) GetTicker(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	if symbol == "" {
		http.Error(w, "Symbol is required", http.StatusBadRequest)
		return
	}

	ticker, err := c.exchangeService.GetTicker(r.Context(), symbol)
	if err != nil {
		c.logger.Error("Failed to get ticker", zap.String("symbol", symbol), zap.Error(err))
		http.Error(w, "Failed to get ticker: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticker)
}

// GetTickers returns tickers for all symbols
func (c *MexcController) GetTickers(w http.ResponseWriter, r *http.Request) {
	// tickers, err := c.exchangeService.GetTickers(r.Context())
	// if err != nil {
	// 	c.logger.Error("Failed to get tickers", zap.Error(err))
	// 	http.Error(w, "Failed to get tickers: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	var tickers interface{} // Placeholder

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickers)
}

// GetKlines returns klines (candlestick data) for a specific symbol
func (c *MexcController) GetKlines(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	if symbol == "" {
		http.Error(w, "Symbol is required", http.StatusBadRequest)
		return
	}

	interval := r.URL.Query().Get("interval")
	if interval == "" {
		interval = "1m" // Default interval
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 500 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	klines, err := c.exchangeService.GetKlines(r.Context(), symbol, interval, limit)
	if err != nil {
		c.logger.Error("Failed to get klines",
			zap.String("symbol", symbol),
			zap.String("interval", interval),
			zap.Error(err))
		http.Error(w, "Failed to get klines: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(klines)
}

// GetOrderBook returns the order book for a specific symbol
func (c *MexcController) GetOrderBook(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	if symbol == "" {
		http.Error(w, "Symbol is required", http.StatusBadRequest)
		return
	}
	// limitStr := r.URL.Query().Get("limit")
	// limit := 100 // Default limit
	// if limitStr != "" {
	// 	parsedLimit, err := strconv.Atoi(limitStr)
	// 	if err == nil && parsedLimit > 0 {
	// 		limit = parsedLimit
	// 	}
	// }

	// orderBook, err := c.exchangeService.GetOrderBook(r.Context(), symbol, limit)
	// if err != nil {
	// 	c.logger.Error("Failed to get order book",
	// 		zap.String("symbol", symbol),
	// 		zap.Int("limit", limit),
	// 		zap.Error(err))
	// 	http.Error(w, "Failed to get order book: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	var orderBook interface{} // Placeholder

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderBook)
}

// GetRecentTrades returns recent trades for a specific symbol
func (c *MexcController) GetRecentTrades(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	if symbol == "" {
		http.Error(w, "Symbol is required", http.StatusBadRequest)
		return
	}

	// limitStr := r.URL.Query().Get("limit")
	// limit := 500 // Default limit
	// if limitStr != "" {
	// 	parsedLimit, err := strconv.Atoi(limitStr)
	// 	if err == nil && parsedLimit > 0 {
	// 		limit = parsedLimit
	// 	}
	// }

	// trades, err := c.exchangeService.GetRecentTrades(r.Context(), symbol, limit)
	// if err != nil {
	// 	c.logger.Error("Failed to get recent trades",
	// 		zap.String("symbol", symbol),
	// 		zap.Int("limit", limit),
	// 		zap.Error(err))
	// 	http.Error(w, "Failed to get recent trades: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	var trades interface{} // Placeholder

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trades)
}

// GetAccountInfo returns account information
func (c *MexcController) GetAccountInfo(w http.ResponseWriter, r *http.Request) {
	// accountInfo, err := c.exchangeService.GetAccountInfo(r.Context())
	// if err != nil {
	// 	c.logger.Error("Failed to get account info", zap.Error(err))
	// 	http.Error(w, "Failed to get account info: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	var accountInfo interface{} // Placeholder

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accountInfo)
}

// GetOpenOrders returns open orders
func (c *MexcController) GetOpenOrders(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")

	orders, err := c.exchangeService.GetOpenOrders(r.Context(), symbol)
	if err != nil {
		c.logger.Error("Failed to get open orders",
			zap.String("symbol", symbol),
			zap.Error(err))
		http.Error(w, "Failed to get open orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// CreateOrder creates a new order
func (c *MexcController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest struct {
		Symbol   string  `json:"symbol"`
		Side     string  `json:"side"`
		Type     string  `json:"type"`
		Quantity float64 `json:"quantity"`
		Price    float64 `json:"price,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if orderRequest.Symbol == "" {
		http.Error(w, "Symbol is required", http.StatusBadRequest)
		return
	}
	if orderRequest.Side == "" {
		http.Error(w, "Side is required", http.StatusBadRequest)
		return
	}
	if orderRequest.Type == "" {
		http.Error(w, "Type is required", http.StatusBadRequest)
		return
	}
	if orderRequest.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than 0", http.StatusBadRequest)
		return
	}
	if orderRequest.Type == "LIMIT" && orderRequest.Price <= 0 {
		http.Error(w, "Price must be greater than 0 for LIMIT orders", http.StatusBadRequest)
		return
	}

	// Create order
	// order, err := c.exchangeService.CreateOrder(
	// 	r.Context(),
	// 	orderRequest.Symbol,
	// 	orderRequest.Side,
	// 	orderRequest.Type,
	// 	orderRequest.Quantity,
	// 	orderRequest.Price,
	// )
	// if err != nil {
	// 	c.logger.Error("Failed to create order",
	// 		zap.String("symbol", orderRequest.Symbol),
	// 		zap.String("side", orderRequest.Side),
	// 		zap.String("type", orderRequest.Type),
	// 		zap.Float64("quantity", orderRequest.Quantity),
	// 		zap.Float64("price", orderRequest.Price),
	// 		zap.Error(err))
	// 	http.Error(w, "Failed to create order: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	var order interface{} // Placeholder

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// CancelOrder cancels an order
func (c *MexcController) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderId := chi.URLParam(r, "orderId")
	if orderId == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Symbol is required", http.StatusBadRequest)
		return
	}

	err := c.exchangeService.CancelOrder(r.Context(), symbol, orderId)
	if err != nil {
		c.logger.Error("Failed to cancel order",
			zap.String("orderId", orderId),
			zap.String("symbol", symbol),
			zap.Error(err))
		http.Error(w, "Failed to cancel order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Order cancelled successfully",
	})
}
