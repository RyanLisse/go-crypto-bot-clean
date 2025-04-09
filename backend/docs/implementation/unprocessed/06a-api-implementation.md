# RESTful API Implementation

This document details the implementation of the REST API for the Go crypto trading bot, providing endpoints for monitoring portfolio status, viewing trade history, and managing trading parameters.

## 1. API Structure Overview

The API follows a clean, resource-oriented architecture with these components:

```
internal/api/
├── handlers/
│   ├── portfolio.go    # Portfolio status endpoints
│   ├── trade.go        # Trading operation endpoints
│   ├── newcoin.go      # New coin detection endpoints
│   ├── config.go       # Configuration endpoints
│   └── health.go       # Health check endpoints
├── middleware/
│   ├── auth.go         # Authentication middleware
│   ├── logging.go      # Request logging middleware
│   └── recovery.go     # Panic recovery middleware
├── dto/
│   ├── request.go      # Request data models
│   └── response.go     # Response data models
└── router.go           # Route definitions
```

## 2. Router Implementation

The router is implemented using the standard library's `net/http` package with a router like Chi:

```go
// internal/api/router.go
package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	
	"github.com/ryanlisse/cryptobot/internal/api/handlers"
	appMiddleware "github.com/ryanlisse/cryptobot/internal/api/middleware"
)

// Config holds API configuration parameters
type Config struct {
	Port       int
	EnableAuth bool
	AuthToken  string
}

// Router creates and configures the API router
func Router(
	portfolioHandler *handlers.PortfolioHandler,
	tradeHandler *handlers.TradeHandler,
	newCoinHandler *handlers.NewCoinHandler,
	configHandler *handlers.ConfigHandler,
	cfg Config,
) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	
	// Add authentication if enabled
	if cfg.EnableAuth {
		r.Use(appMiddleware.TokenAuth(cfg.AuthToken))
	}
	
	// Health check
	r.Get("/health", handlers.HealthCheck)
	
	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Portfolio endpoints
		r.Route("/portfolio", func(r chi.Router) {
			r.Get("/", portfolioHandler.GetPortfolioSummary)
			r.Get("/active", portfolioHandler.GetActiveTrades)
			r.Get("/performance", portfolioHandler.GetPerformanceMetrics)
			r.Get("/value", portfolioHandler.GetTotalValue)
		})

		// Trading endpoints
		r.Route("/trade", func(r chi.Router) {
			r.Get("/history", tradeHandler.GetTradeHistory)
			r.Post("/buy", tradeHandler.ExecuteTrade)
			r.Post("/sell", tradeHandler.SellCoin)
			r.Get("/status/{id}", tradeHandler.GetTradeStatus)
		})

		// New coin detection endpoints
		r.Route("/newcoins", func(r chi.Router) {
			r.Get("/", newCoinHandler.GetDetectedCoins)
			r.Post("/process", newCoinHandler.ProcessNewCoins)
		})

		// Configuration endpoints
		r.Route("/config", func(r chi.Router) {
			r.Get("/", configHandler.GetCurrentConfig)
			r.Put("/", configHandler.UpdateConfig)
			r.Get("/defaults", configHandler.GetDefaultConfig)
		})
	})

	return r
}
```

## 3. Handler Implementation

Handlers follow a consistent structure, accepting service dependencies via constructor injection and processing requests:

```go
// internal/api/handlers/portfolio.go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ryanlisse/cryptobot/internal/api/dto"
	"github.com/ryanlisse/cryptobot/internal/domain/service"
)

// PortfolioHandler handles portfolio-related API endpoints
type PortfolioHandler struct {
	portfolioService service.PortfolioService
}

// NewPortfolioHandler creates a new portfolio handler with dependencies
func NewPortfolioHandler(portfolioService service.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{
		portfolioService: portfolioService,
	}
}

// GetPortfolioSummary returns a summary of the current portfolio
func (h *PortfolioHandler) GetPortfolioSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get active trades
	activeTrades, err := h.portfolioService.GetActiveTrades(ctx)
	if err != nil {
		http.Error(w, "Failed to get active trades: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get portfolio value
	totalValue, err := h.portfolioService.GetPortfolioValue(ctx)
	if err != nil {
		http.Error(w, "Failed to get portfolio value: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get performance metrics
	metrics, err := h.portfolioService.GetTradePerformance(ctx, "all")
	if err != nil {
		http.Error(w, "Failed to get performance metrics: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Build response
	response := dto.PortfolioSummaryResponse{
		TotalValue:      totalValue,
		ActiveTradeCount: len(activeTrades),
		ActiveTrades:    mapToTradeResponses(activeTrades),
		Performance:     mapToPerformanceResponse(metrics),
		Timestamp:       time.Now(),
	}
	
	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetActiveTrades returns all active trading positions
func (h *PortfolioHandler) GetActiveTrades(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get active trades
	activeTrades, err := h.portfolioService.GetActiveTrades(ctx)
	if err != nil {
		http.Error(w, "Failed to get active trades: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Map to response DTOs
	response := dto.ActiveTradesResponse{
		Trades:    mapToTradeResponses(activeTrades),
		Count:     len(activeTrades),
		Timestamp: time.Now(),
	}
	
	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Helper functions to map domain models to DTOs
func mapToTradeResponses(trades []*models.BoughtCoin) []dto.TradeResponse {
	response := make([]dto.TradeResponse, len(trades))
	for i, trade := range trades {
		response[i] = dto.TradeResponse{
			ID:              trade.ID,
			Symbol:          trade.Symbol,
			PurchasePrice:   trade.PurchasePrice,
			CurrentPrice:    trade.CurrentPrice,
			Quantity:        trade.Quantity,
			PurchaseTime:    trade.PurchaseTime,
			ProfitPercent:   trade.ProfitPercentage,
			CurrentValue:    trade.CurrentValue,
			StopLossPrice:   trade.StopLossPrice,
			TakeProfitLevels: mapToTakeProfitLevels(trade.TakeProfitLevels),
		}
	}
	return response
}
```

## 4. Data Transfer Objects (DTOs)

DTOs ensure clean separation between API models and domain models:

```go
// internal/api/dto/request.go
package dto

import "time"

// TradeRequest represents a request to execute a trade
type TradeRequest struct {
	Symbol string  `json:"symbol"`
	Amount float64 `json:"amount,omitempty"`
}

// SellRequest represents a request to sell a coin
type SellRequest struct {
	CoinID uint    `json:"coin_id"`
	Amount float64 `json:"amount,omitempty"`
	All    bool    `json:"all,omitempty"`
}

// ConfigUpdateRequest represents a request to update bot configuration
type ConfigUpdateRequest struct {
	USDTPerTrade    *float64   `json:"usdt_per_trade,omitempty"`
	StopLossPercent *float64   `json:"stop_loss_percent,omitempty"`
	TakeProfitLevels []float64 `json:"take_profit_levels,omitempty"`
	SellPercentages  []float64 `json:"sell_percentages,omitempty"`
}

// response.go
package dto

import "time"

// PortfolioSummaryResponse represents the overall portfolio status
type PortfolioSummaryResponse struct {
	TotalValue      float64         `json:"total_value"`
	ActiveTradeCount int            `json:"active_trade_count"`
	ActiveTrades    []TradeResponse `json:"active_trades"`
	Performance     PerformanceResponse `json:"performance"`
	Timestamp       time.Time       `json:"timestamp"`
}

// TradeResponse represents a single trading position
type TradeResponse struct {
	ID              uint      `json:"id"`
	Symbol          string    `json:"symbol"`
	PurchasePrice   float64   `json:"purchase_price"`
	CurrentPrice    float64   `json:"current_price"`
	Quantity        float64   `json:"quantity"`
	PurchaseTime    time.Time `json:"purchase_time"`
	ProfitPercent   float64   `json:"profit_percent"`
	CurrentValue    float64   `json:"current_value"`
	StopLossPrice   float64   `json:"stop_loss_price"`
	TakeProfitLevels []TakeProfitLevelResponse `json:"take_profit_levels"`
}
```

## 5. Middleware Implementation

```go
// internal/api/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"
)

// TokenAuth creates middleware for simple API token authentication
func TokenAuth(validToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}
			
			// Check format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}
			
			// Validate token
			if parts[1] != validToken {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			
			// Token is valid, continue
			next.ServeHTTP(w, r)
		})
	}
}
```

## 6. Server Implementation

```go
// cmd/server/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/ryanlisse/cryptobot/internal/api"
	"github.com/ryanlisse/cryptobot/internal/api/handlers"
	"github.com/ryanlisse/cryptobot/internal/core/newcoin"
	"github.com/ryanlisse/cryptobot/internal/core/portfolio"
	"github.com/ryanlisse/cryptobot/internal/core/trade"
	"github.com/ryanlisse/cryptobot/internal/platform/database"
	"github.com/ryanlisse/cryptobot/internal/platform/mexc"
)

func main() {
	// Parse command line flags
	port := flag.Int("port", 8080, "HTTP server port")
	dbPath := flag.String("db", "data/cryptobot.db", "SQLite database path")
	apiKey := flag.String("apikey", "", "MEXC API key")
	secretKey := flag.String("secret", "", "MEXC API secret key")
	flag.Parse()
	
	// Check required flags
	if *apiKey == "" || *secretKey == "" {
		log.Fatal("MEXC API key and secret key are required")
	}
	
	// Create a context that's canceled on Ctrl+C
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v", sig)
		cancel()
	}()
	
	// Connect to database
	dbConfig := database.Config{
		Path: *dbPath,
	}
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)
	log.Println("Connected to database")
	
	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations applied")
	
	// Initialize repositories
	// ...
	
	// Create exchange service
	exchangeService, err := mexc.NewClient(*apiKey, *secretKey, "")
	if err != nil {
		log.Fatalf("Failed to create MEXC client: %v", err)
	}
	if err := exchangeService.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to MEXC: %v", err)
	}
	defer exchangeService.Disconnect()
	log.Println("Connected to MEXC exchange")
	
	// Initialize core services
	tradeService, err := trade.NewService(
		exchangeService,
		boughtCoinRepo,
		logRepo,
		trade.DefaultConfig(),
	)
	if err != nil {
		log.Fatalf("Failed to create trade service: %v", err)
	}
	
	portfolioService := portfolio.NewService(
		exchangeService,
		boughtCoinRepo,
	)
	
	newCoinService := newcoin.NewService(
		exchangeService,
		newCoinRepo,
		purchaseDecisionRepo,
		logRepo,
		newcoin.DefaultConfig(),
	)
	
	// Create API handlers
	portfolioHandler := handlers.NewPortfolioHandler(portfolioService)
	tradeHandler := handlers.NewTradeHandler(tradeService)
	newCoinHandler := handlers.NewNewCoinHandler(newCoinService)
	configHandler := handlers.NewConfigHandler(/* config service */)
	
	// Create router
	router := api.Router(
		portfolioHandler,
		tradeHandler,
		newCoinHandler,
		configHandler,
		api.Config{
			Port:       *port,
			EnableAuth: true,
			AuthToken:  os.Getenv("API_TOKEN"),
		},
	)
	
	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}
	
	// Start server in a goroutine
	go func() {
		log.Printf("API server listening on port %d", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
	
	// Wait for cancellation signal
	<-ctx.Done()
	log.Println("Shutting down API server...")
	
	// Create a shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	
	// Gracefully shutdown the server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	
	log.Println("Server gracefully stopped")
}
```

## 7. API Documentation

Complete API documentation should be created using tools like Swagger or OpenAPI:

```bash
# Install Swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
swag init -g cmd/server/main.go -o ./docs/swagger
```

Document each endpoint with proper annotations for auto-generating documentation:

```go
// @Summary Get portfolio summary
// @Description Returns a summary of the current portfolio including total value and active trades
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} dto.PortfolioSummaryResponse
// @Failure 500 {string} string "Error message"
// @Router /api/v1/portfolio [get]
func (h *PortfolioHandler) GetPortfolioSummary(w http.ResponseWriter, r *http.Request) {
    // Implementation...
}
```

## 8. Error Handling Strategy

Implement a consistent error handling approach:

```go
// internal/api/errors/errors.go
package errors

import (
	"encoding/json"
	"net/http"
)

// APIError represents a standardized API error response
type APIError struct {
	Status  int    `json:"-"`       // HTTP status code
	Code    string `json:"code"`    // Application-specific error code
	Message string `json:"message"` // Human-readable message
	Details any    `json:"details,omitempty"` // Additional error details
}

// Common error codes
const (
	CodeInvalidRequest = "INVALID_REQUEST"
	CodeNotFound       = "NOT_FOUND"
	CodeServerError    = "SERVER_ERROR"
	CodeUnauthorized   = "UNAUTHORIZED"
)

// Write sends the error as a JSON response
func (e APIError) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	_ = json.NewEncoder(w).Encode(e)
}

// BadRequest creates a 400 Bad Request error
func BadRequest(message string, details any) APIError {
	return APIError{
		Status:  http.StatusBadRequest,
		Code:    CodeInvalidRequest,
		Message: message,
		Details: details,
	}
}

// NotFound creates a 404 Not Found error
func NotFound(message string) APIError {
	return APIError{
		Status:  http.StatusNotFound,
		Code:    CodeNotFound,
		Message: message,
	}
}

// ServerError creates a 500 Internal Server Error
func ServerError(message string) APIError {
	return APIError{
		Status:  http.StatusInternalServerError,
		Code:    CodeServerError,
		Message: message,
	}
}
```

## 9. Testing

Create comprehensive tests for all API endpoints:

```go
// internal/api/handlers/portfolio_test.go
package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	
	"github.com/ryanlisse/cryptobot/internal/api/dto"
	"github.com/ryanlisse/cryptobot/internal/api/handlers"
	"github.com/ryanlisse/cryptobot/internal/domain/models"
)

// MockPortfolioService implements a mock of the PortfolioService for testing
type MockPortfolioService struct {
	activeTrades []*models.BoughtCoin
	portfolioValue float64
	metrics *models.PerformanceMetrics
}

// ... mock method implementations ...

func TestGetPortfolioSummary(t *testing.T) {
	// Arrange
	now := time.Now()
	mockService := &MockPortfolioService{
		activeTrades: []*models.BoughtCoin{
			{
				ID:            1,
				Symbol:        "BTC/USDT",
				PurchasePrice: 30000,
				CurrentPrice:  35000,
				Quantity:      0.1,
				PurchaseTime:  now.Add(-24 * time.Hour),
			},
		},
		portfolioValue: 5000,
		metrics: &models.PerformanceMetrics{
			TotalTrades:   10,
			WinningTrades: 7,
			LosingTrades:  3,
		},
	}
	
	handler := handlers.NewPortfolioHandler(mockService)
	req := httptest.NewRequest("GET", "/api/v1/portfolio", nil)
	rr := httptest.NewRecorder()
	
	// Act
	handler.GetPortfolioSummary(rr, req)
	
	// Assert
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
	
	var response dto.PortfolioSummaryResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if response.TotalValue != 5000 {
		t.Errorf("Expected total value 5000, got %.2f", response.TotalValue)
	}
	
	if response.ActiveTradeCount != 1 {
		t.Errorf("Expected 1 active trade, got %d", response.ActiveTradeCount)
	}
}
```
