package huma

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestSetupHuma(t *testing.T) {
	// Create a new router
	router := chi.NewRouter()

	// Setup Huma with the router
	api := SetupHuma(router, DefaultConfig())

	// Verify that the API was created
	assert.NotNil(t, api, "API should not be nil")

	// Create a test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Test the OpenAPI endpoint
	resp, err := http.Get(server.URL + "/openapi.json")
	assert.NoError(t, err, "Should not error when getting OpenAPI spec")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")

	// Decode the OpenAPI spec
	var openAPI map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&openAPI)
	assert.NoError(t, err, "Should not error when decoding OpenAPI spec")

	// Verify the OpenAPI spec
	info, ok := openAPI["info"].(map[string]interface{})
	assert.True(t, ok, "Should have info object")
	assert.Equal(t, "Crypto Trading Bot API", info["title"], "Should have correct title")
	assert.Equal(t, "1.0.0", info["version"], "Should have correct version")

	// Test the docs endpoint
	resp, err = http.Get(server.URL + "/docs")
	assert.NoError(t, err, "Should not error when getting docs")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")
}

func TestHealthEndpoint(t *testing.T) {
	// Create a test API
	_, api := humatest.New(t)

	// Register the health endpoint
	huma.Register(api, huma.Operation{
		OperationID: "health-check",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health check",
		Description: "Returns the health status of the API",
		Tags:        []string{"System"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Status string `json:"status" doc:"Health status" example:"ok"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Status string `json:"status" doc:"Health status" example:"ok"`
			}
		}{}
		resp.Body.Status = "ok"
		return resp, nil
	})

	// Test the health endpoint
	resp := api.Get("/health")
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Decode the response
	var health struct {
		Status string `json:"status"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &health)
	assert.NoError(t, err, "Should not error when decoding response")
	assert.Equal(t, "ok", health.Status, "Should return ok status")
}

func TestPortfolioEndpoints(t *testing.T) {
	// Create a test API
	_, api := humatest.New(t)

	// Register the portfolio endpoints
	registerPortfolioEndpoints(api, "/api/v1")

	// Test the portfolio summary endpoint
	resp := api.Get("/api/v1/portfolio")
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the active trades endpoint
	resp = api.Get("/api/v1/portfolio/active")
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the performance metrics endpoint
	resp = api.Get("/api/v1/portfolio/performance")
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the total value endpoint
	resp = api.Get("/api/v1/portfolio/value")
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")
}

func TestTradeEndpoints(t *testing.T) {
	// Create a test API
	_, api := humatest.New(t)

	// Register the trade endpoints
	registerTradeEndpoints(api, "/api/v1")

	// Test the trade history endpoint
	resp := api.Get("/api/v1/trade/history")
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the execute trade endpoint
	resp = api.Post("/api/v1/trade/buy", map[string]interface{}{
		"symbol": "BTC/USDT",
		"amount": 100.0,
	})
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the sell coin endpoint
	resp = api.Post("/api/v1/trade/sell", map[string]interface{}{
		"coin_id": 123,
		"amount":  0.001,
		"all":     false,
	})
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the trade status endpoint
	resp = api.Get("/api/v1/trade/status/123")
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")
}

func TestNewCoinEndpoints(t *testing.T) {
	// Create a test API
	_, api := humatest.New(t)

	// Register the newcoin endpoints
	registerNewCoinEndpoints(api, "/api/v1")

	// Test the get detected coins endpoint
	resp := api.Get("/api/v1/newcoins")
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the process new coins endpoint
	resp = api.Post("/api/v1/newcoins/process", map[string]interface{}{
		"coin_ids": []int{123, 456},
	})
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the detect new coins endpoint
	resp = api.Post("/api/v1/newcoins/detect", map[string]interface{}{})
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")
}

func TestConfigEndpoints(t *testing.T) {
	// Create a test API
	_, api := humatest.New(t)

	// Register the config endpoints
	registerConfigEndpoints(api, "/api/v1")

	// Test the get current config endpoint
	resp := api.Get("/api/v1/config")
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the update config endpoint
	resp = api.Put("/api/v1/config", map[string]interface{}{
		"usdt_per_trade":     20.0,
		"stop_loss_percent":  10.0,
		"take_profit_levels": []float64{5.0, 10.0, 15.0, 20.0},
		"sell_percentages":   []float64{0.25, 0.25, 0.25, 0.25},
	})
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the get default config endpoint
	resp = api.Get("/api/v1/config/defaults")
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")
}

func TestStatusEndpoints(t *testing.T) {
	// Create a test API
	_, api := humatest.New(t)

	// Register the status endpoints
	registerStatusEndpoints(api, "/api/v1")

	// Test the get status endpoint
	resp := api.Get("/api/v1/status")
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the start processes endpoint
	resp = api.Post("/api/v1/status/start", map[string]interface{}{})
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the stop processes endpoint
	resp = api.Post("/api/v1/status/stop", map[string]interface{}{})
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")
}
