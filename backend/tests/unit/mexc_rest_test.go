package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-crypto-bot-clean/backend/internal/platform/mexc/rest"
)

func createTestClient(handler http.HandlerFunc) *rest.Client {
	server := httptest.NewServer(handler)
	client, _ := rest.NewClient("test-api-key", "test-secret-key",
		rest.WithBaseURL(server.URL),
	)
	return client
}

func TestGetTicker(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"symbol":             "BTCUSDT",
			"lastPrice":          "50000.00",
			"volume":             "1000.5",
			"priceChange":        "500.00",
			"priceChangePercent": "1.01",
			"highPrice":          "51000.00",
			"lowPrice":           "49000.00",
		}
		json.NewEncoder(w).Encode(resp)
	}

	client := createTestClient(handler)
	ctx := context.Background()

	ticker, err := client.GetTicker(ctx, "BTCUSDT")
	assert.NoError(t, err)
	assert.NotNil(t, ticker)
	assert.Equal(t, "BTCUSDT", ticker.Symbol)
	assert.Equal(t, 50000.00, ticker.Price)
	assert.Equal(t, 1000.5, ticker.Volume)
}

func TestGetAllTickers(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		resp := []map[string]string{
			{
				"symbol":             "BTCUSDT",
				"lastPrice":          "50000.00",
				"volume":             "1000.5",
				"priceChange":        "500.00",
				"priceChangePercent": "1.01",
				"highPrice":          "51000.00",
				"lowPrice":           "49000.00",
			},
			{
				"symbol":             "ETHUSDT",
				"lastPrice":          "3000.00",
				"volume":             "5000.0",
				"priceChange":        "50.00",
				"priceChangePercent": "1.67",
				"highPrice":          "3100.00",
				"lowPrice":           "2900.00",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}

	client := createTestClient(handler)
	ctx := context.Background()

	tickers, err := client.GetAllTickers(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tickers)
	assert.Len(t, tickers, 2)
	assert.Contains(t, tickers, "BTCUSDT")
	assert.Contains(t, tickers, "ETHUSDT")
}

func TestGetTickerError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	client := createTestClient(handler)
	ctx := context.Background()

	ticker, err := client.GetTicker(ctx, "BTCUSDT")
	assert.Error(t, err)
	assert.Nil(t, ticker)
}
