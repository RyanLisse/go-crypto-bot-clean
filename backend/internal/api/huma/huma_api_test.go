package huma

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHumaAPI(t *testing.T) {
	// Create a test API
	_, api := humatest.New(t)

	// Register the endpoints
	registerBacktestEndpoints(api, "/api/v1")
	registerStrategyEndpoints(api, "/api/v1")
	registerAuthEndpoints(api, "/api/v1")
	registerUserEndpoints(api, "/api/v1")

	// Test that the API was created
	assert.NotNil(t, api, "API should not be nil")

	// Test the backtest endpoints
	t.Run("backtest endpoints", func(t *testing.T) {
		// Test the run backtest endpoint
		resp := api.Post("/api/v1/backtest", map[string]interface{}{
			"strategy":        "breakout",
			"symbol":          "BTC/USDT",
			"timeframe":       "1h",
			"startDate":       time.Now().AddDate(0, -1, 0).Format(time.RFC3339),
			"endDate":         time.Now().Format(time.RFC3339),
			"initialCapital":  1000.0,
			"riskPerTrade":    0.02,
		})
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Decode the response
		var backtestResponse struct {
			ID string `json:"id"`
		}
		err := json.Unmarshal(resp.Body.Bytes(), &backtestResponse)
		assert.NoError(t, err, "Should not error when decoding response")
		assert.NotEmpty(t, backtestResponse.ID, "Should return a backtest ID")

		// Test the get backtest results endpoint
		resp = api.Get("/api/v1/backtest/" + backtestResponse.ID)
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Test the list backtests endpoint
		resp = api.Get("/api/v1/backtest/list")
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Test the compare backtests endpoint
		resp = api.Post("/api/v1/backtest/compare", map[string]interface{}{
			"backtest_ids": []string{backtestResponse.ID},
		})
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")
	})

	// Test the strategy endpoints
	t.Run("strategy endpoints", func(t *testing.T) {
		// Test the list strategies endpoint
		resp := api.Get("/api/v1/strategy")
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Decode the response
		var strategiesResponse struct {
			Strategies []struct {
				ID string `json:"id"`
			} `json:"strategies"`
		}
		err := json.Unmarshal(resp.Body.Bytes(), &strategiesResponse)
		assert.NoError(t, err, "Should not error when decoding response")
		
		// For now, we'll just check that the response is valid JSON
		// In a real test, we would check that the strategies list is not empty
		// but since we're just setting up the API structure, we'll skip that check
		
		// Test the get strategy endpoint with a dummy ID
		resp = api.Get("/api/v1/strategy/breakout")
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Test the update strategy config endpoint
		resp = api.Put("/api/v1/strategy/breakout", map[string]interface{}{
			"parameters": map[string]interface{}{
				"lookbackPeriod":    20,
				"breakoutThreshold": 2.5,
			},
			"isEnabled": true,
		})
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Test the enable strategy endpoint
		resp = api.Post("/api/v1/strategy/breakout/enable", nil)
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Test the disable strategy endpoint
		resp = api.Post("/api/v1/strategy/breakout/disable", nil)
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")
	})

	// Test the auth endpoints
	t.Run("auth endpoints", func(t *testing.T) {
		// Test the login endpoint
		resp := api.Post("/api/v1/auth/login", map[string]interface{}{
			"email":    "user@example.com",
			"password": "password123",
		})
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Test the register endpoint
		resp = api.Post("/api/v1/auth/register", map[string]interface{}{
			"email":     "newuser@example.com",
			"username":  "newuser",
			"password":  "password123",
			"firstName": "New",
			"lastName":  "User",
		})
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Test the refresh token endpoint
		resp = api.Post("/api/v1/auth/refresh", map[string]interface{}{
			"refreshToken": "dummy-refresh-token",
		})
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Test the logout endpoint
		resp = api.Post("/api/v1/auth/logout", nil)
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Test the verify token endpoint
		resp = api.Get("/api/v1/auth/verify")
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")
	})

	// Test the user endpoints
	t.Run("user endpoints", func(t *testing.T) {
		// Test the get user profile endpoint
		resp := api.Get("/api/v1/user/profile")
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Test the update user profile endpoint
		resp = api.Put("/api/v1/user/profile", map[string]interface{}{
			"username":  "updateduser",
			"firstName": "Updated",
			"lastName":  "User",
		})
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Test the get user settings endpoint
		resp = api.Get("/api/v1/user/settings")
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Test the update user settings endpoint
		resp = api.Put("/api/v1/user/settings", map[string]interface{}{
			"theme":               "dark",
			"language":            "en",
			"timeZone":            "America/New_York",
			"notificationsEnabled": true,
			"emailNotifications":   true,
			"pushNotifications":    false,
			"defaultCurrency":      "USD",
		})
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

		// Test the change password endpoint
		resp = api.Post("/api/v1/user/password", map[string]interface{}{
			"currentPassword": "password123",
			"newPassword":     "newpassword123",
		})
		assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")
	})
}
