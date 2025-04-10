package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-crypto-bot-clean/backend/internal/api/http/dto"
	"go-crypto-bot-clean/backend/internal/application/services"

	"github.com/gorilla/mux"
)

// TradeController handles HTTP requests related to trades
type TradeController struct {
	tradeService *services.TradeService
}

// NewTradeController creates a new TradeController
func NewTradeController(tradeService *services.TradeService) *TradeController {
	return &TradeController{
		tradeService: tradeService,
	}
}

// RegisterRoutes registers the trade routes
func (c *TradeController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/trades", c.ListTrades).Methods("GET")
	router.HandleFunc("/api/trades/{id}", c.GetTrade).Methods("GET")
	router.HandleFunc("/api/orders/{orderId}/trades", c.GetTradesByOrderID).Methods("GET")
}

// ListTrades handles retrieving a list of trades
func (c *TradeController) ListTrades(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	limitStr := r.URL.Query().Get("limit")
	
	limit := 50 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	trades, err := c.tradeService.GetTradesBySymbol(r.Context(), symbol, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []*dto.TradeResponse
	for _, trade := range trades {
		response = append(response, dto.TradeResponseFromModel(trade))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTrade handles retrieving a trade by ID
func (c *TradeController) GetTrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	trade, err := c.tradeService.GetTradeByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := dto.TradeResponseFromModel(trade)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTradesByOrderID handles retrieving trades by order ID
func (c *TradeController) GetTradesByOrderID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderId"]

	trades, err := c.tradeService.GetTradesByOrderID(r.Context(), orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []*dto.TradeResponse
	for _, trade := range trades {
		response = append(response, dto.TradeResponseFromModel(trade))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
