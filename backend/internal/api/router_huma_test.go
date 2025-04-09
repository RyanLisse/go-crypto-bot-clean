package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/api/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSetupChiRouter(t *testing.T) {
	// Create dependencies
	deps := &HumaDependencies{
		HealthHandler:    handlers.NewHealthHandler(),
		StatusHandler:    &handlers.StatusHandler{},
		PortfolioHandler: &handlers.PortfolioHandler{},
		TradeHandler:     &handlers.TradeHandler{},
		NewCoinHandler:   &handlers.NewCoinsHandler{},
		ConfigHandler:    &handlers.ConfigHandler{},
		WebSocketHandler: &websocket.Handler{},
	}

	// Setup router
	router := SetupChiRouter(deps)

	// Create a test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Test the health endpoint
	resp, err := http.Get(server.URL + "/health")
	assert.NoError(t, err, "Should not error when getting health")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")

	// Test the OpenAPI endpoint
	resp, err = http.Get(server.URL + "/openapi.json")
	assert.NoError(t, err, "Should not error when getting OpenAPI spec")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")

	// Test the docs endpoint
	resp, err = http.Get(server.URL + "/docs")
	assert.NoError(t, err, "Should not error when getting docs")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")
}

func TestCORSMiddleware(t *testing.T) {
	// Create dependencies
	deps := &HumaDependencies{
		HealthHandler:    handlers.NewHealthHandler(),
		StatusHandler:    &handlers.StatusHandler{},
		PortfolioHandler: &handlers.PortfolioHandler{},
		TradeHandler:     &handlers.TradeHandler{},
		NewCoinHandler:   &handlers.NewCoinsHandler{},
		ConfigHandler:    &handlers.ConfigHandler{},
		WebSocketHandler: &websocket.Handler{},
	}

	// Setup router
	router := SetupChiRouter(deps)

	// Create a test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Create a request with OPTIONS method
	req, err := http.NewRequest(http.MethodOptions, server.URL+"/health", nil)
	assert.NoError(t, err, "Should not error when creating request")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Should not error when sending request")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")

	// Check CORS headers
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"), "Should have correct CORS header")
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"), "Should have correct CORS header")
	assert.Equal(t, "Content-Type, Authorization", resp.Header.Get("Access-Control-Allow-Headers"), "Should have correct CORS header")
}
