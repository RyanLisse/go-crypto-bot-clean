package strategy

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
)

func TestStrategyEndpoints(t *testing.T) {
	// Create a test API
	_, api := humatest.New(t)

	// Register the strategy endpoints
	RegisterEndpoints(api, "/api/v1")

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
	assert.NotEmpty(t, strategiesResponse.Strategies, "Should return a list of strategies")

	// Get the first strategy ID
	strategyID := strategiesResponse.Strategies[0].ID

	// Test the get strategy endpoint
	resp = api.Get("/api/v1/strategy/" + strategyID)
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the update strategy config endpoint
	resp = api.Put("/api/v1/strategy/"+strategyID, map[string]interface{}{
		"parameters": map[string]interface{}{
			"lookbackPeriod":    20,
			"breakoutThreshold": 2.5,
		},
		"isEnabled": true,
	})
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Test the enable strategy endpoint
	resp = api.Post("/api/v1/strategy/"+strategyID+"/enable", map[string]interface{}{})
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Decode the response
	var enableResponse struct {
		ID        string `json:"id"`
		IsEnabled bool   `json:"isEnabled"`
	}
	err = json.Unmarshal(resp.Body.Bytes(), &enableResponse)
	assert.NoError(t, err, "Should not error when decoding response")
	assert.Equal(t, strategyID, enableResponse.ID, "Should return the correct strategy ID")
	assert.True(t, enableResponse.IsEnabled, "Strategy should be enabled")

	// Test the disable strategy endpoint
	resp = api.Post("/api/v1/strategy/"+strategyID+"/disable", map[string]interface{}{})
	assert.Equal(t, http.StatusOK, resp.Code, "Should return 200 OK")

	// Decode the response
	var disableResponse struct {
		ID        string `json:"id"`
		IsEnabled bool   `json:"isEnabled"`
	}
	err = json.Unmarshal(resp.Body.Bytes(), &disableResponse)
	assert.NoError(t, err, "Should not error when decoding response")
	assert.Equal(t, strategyID, disableResponse.ID, "Should return the correct strategy ID")
	assert.False(t, disableResponse.IsEnabled, "Strategy should be disabled")
}
