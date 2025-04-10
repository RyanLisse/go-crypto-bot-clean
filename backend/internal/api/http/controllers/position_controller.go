package controllers

import (
	"encoding/json"
	"net/http"

	"go-crypto-bot-clean/backend/internal/api/http/dto"
	"go-crypto-bot-clean/backend/internal/application/services"
	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/gorilla/mux"
)

// PositionController handles HTTP requests related to positions
type PositionController struct {
	positionService *services.PositionService
}

// NewPositionController creates a new PositionController
func NewPositionController(positionService *services.PositionService) *PositionController {
	return &PositionController{
		positionService: positionService,
	}
}

// RegisterRoutes registers the position routes
func (c *PositionController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/positions", c.ListPositions).Methods("GET")
	router.HandleFunc("/api/positions/{id}", c.GetPosition).Methods("GET")
	router.HandleFunc("/api/positions/{id}/close", c.ClosePosition).Methods("POST")
}

// ListPositions handles retrieving a list of positions
func (c *PositionController) ListPositions(w http.ResponseWriter, r *http.Request) {
	status := models.PositionStatus(r.URL.Query().Get("status"))

	positions, err := c.positionService.ListPositions(r.Context(), status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []*dto.PositionResponse
	for _, position := range positions {
		response = append(response, dto.PositionResponseFromModel(position))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPosition handles retrieving a position by ID
func (c *PositionController) GetPosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	position, err := c.positionService.GetPositionByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := dto.PositionResponseFromModel(position)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ClosePosition handles closing a position
func (c *PositionController) ClosePosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	position, err := c.positionService.ClosePosition(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.PositionResponseFromModel(position)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
