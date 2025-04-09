package gemini

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestNewGeminiClient(t *testing.T) {
	client := NewGeminiClient("test-api-key")
	assert.NotNil(t, client)
	assert.Equal(t, "test-api-key", client.APIKey)
	assert.NotNil(t, client.HTTPClient)
}

func TestGeminiClient_AnalyzeMetrics(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-api-key", r.Header.Get("x-goog-api-key"))

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"candidates": [
				{
					"content": {
						"parts": [
							{
								"text": "This is a test analysis"
							}
						]
					}
				}
			]
		}`))
	}))
	defer server.Close()

	// Create a client with the test server URL
	client := NewGeminiClient("test-api-key")
	client.Endpoint = server.URL

	// Create test metrics
	metrics := models.PerformanceReportMetrics{
		Timestamp: time.Now(),
		Metrics: map[string]interface{}{
			"test_metric": 123.45,
		},
		SystemState: models.SystemState{
			CPUUsage:    10.5,
			MemoryUsage: 256.0,
			Latency:     50,
			Goroutines:  10,
			Uptime:      "1h 30m",
		},
	}

	// Test the AnalyzeMetrics method
	analysis, err := client.AnalyzeMetrics(context.Background(), metrics)
	assert.NoError(t, err)
	assert.Equal(t, "This is a test analysis", analysis)
}

func TestGeminiClient_ExtractInsights(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-api-key", r.Header.Get("x-goog-api-key"))

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"candidates": [
				{
					"content": {
						"parts": [
							{
								"text": "Insight 1\nInsight 2\nInsight 3"
							}
						]
					}
				}
			]
		}`))
	}))
	defer server.Close()

	// Create a client with the test server URL
	client := NewGeminiClient("test-api-key")
	client.Endpoint = server.URL

	// Test the ExtractInsights method
	insights, err := client.ExtractInsights(context.Background(), "This is a test analysis")
	assert.NoError(t, err)
	assert.Equal(t, []string{"Insight 1", "Insight 2", "Insight 3"}, insights)
}

func TestGeminiClient_AnalyzeMetrics_Error(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"code": 400,
				"message": "Bad request",
				"status": "INVALID_ARGUMENT"
			}
		}`))
	}))
	defer server.Close()

	// Create a client with the test server URL
	client := NewGeminiClient("test-api-key")
	client.Endpoint = server.URL

	// Create test metrics
	metrics := models.PerformanceReportMetrics{
		Timestamp: time.Now(),
		Metrics: map[string]interface{}{
			"test_metric": 123.45,
		},
	}

	// Test the AnalyzeMetrics method
	_, err := client.AnalyzeMetrics(context.Background(), metrics)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error")
}
