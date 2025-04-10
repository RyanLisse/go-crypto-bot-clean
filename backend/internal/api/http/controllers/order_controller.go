package controllers

import (
	"encoding/json"
	"net/http"

	"go-crypto-bot-clean/backend/internal/api/http/dto"
	"go-crypto-bot-clean/backend/internal/application/services"
	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/gorilla/mux"
)

// OrderController handles HTTP requests related to orders
type OrderController struct {
	orderService *services.OrderService
}

// NewOrderController creates a new OrderController
func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

// RegisterRoutes registers the order routes
func (c *OrderController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/orders", c.CreateOrder).Methods("POST")
	router.HandleFunc("/api/orders/{id}", c.GetOrder).Methods("GET")
	router.HandleFunc("/api/orders", c.ListOrders).Methods("GET")
	router.HandleFunc("/api/orders/{id}", c.CancelOrder).Methods("DELETE")
}

// CreateOrder handles the creation of a new order
func (c *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest dto.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order := orderRequest.ToModel()
	if err := c.orderService.CreateOrder(r.Context(), order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.OrderResponseFromModel(order)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetOrder handles retrieving an order by ID
func (c *OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	order, err := c.orderService.GetOrderByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := dto.OrderResponseFromModel(order)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListOrders handles retrieving a list of orders
func (c *OrderController) ListOrders(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	status := models.OrderStatus(r.URL.Query().Get("status"))

	orders, err := c.orderService.ListOrders(r.Context(), symbol, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []*dto.OrderResponse
	for _, order := range orders {
		response = append(response, dto.OrderResponseFromModel(order))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CancelOrder handles canceling an order
func (c *OrderController) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := c.orderService.CancelOrder(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
