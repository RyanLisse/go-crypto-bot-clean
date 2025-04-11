package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	// Define a simple status handler
	http.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		response := map[string]interface{}{
			"status":  "ok",
			"version": "1.0.0",
			"services": map[string]string{
				"database":  "connected",
				"exchange":  "connected",
				"websocket": "connected",
			},
		}

		json.NewEncoder(w).Encode(response)
	})

	// Define a health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Portfolio endpoints
	http.HandleFunc("/api/v1/portfolio", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"totalValue": 25432.98,
			"assets": []map[string]interface{}{
				{
					"symbol": "BTC",
					"amount": 0.5,
					"value":  15432.15,
				},
				{
					"symbol": "ETH",
					"amount": 5.2,
					"value":  8920.43,
				},
				{
					"symbol": "USDT",
					"amount": 1080.40,
					"value":  1080.40,
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	})

	// Portfolio value endpoint
	http.HandleFunc("/api/v1/portfolio/value", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"totalValue": 25432.98,
			"currency":   "USD",
			"timestamp":  time.Now().Unix(),
		}

		json.NewEncoder(w).Encode(response)
	})

	// Portfolio performance endpoint
	http.HandleFunc("/api/v1/portfolio/performance", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"dailyChange":        1.25,
			"weeklyChange":       3.75,
			"monthlyChange":      12.5,
			"yearlyChange":       45.30,
			"totalChange":        112.45,
			"dailyChangeValue":   315.22,
			"weeklyChangeValue":  945.67,
			"monthlyChangeValue": 3120.45,
			"yearlyChangeValue":  8950.67,
			"performance": []map[string]interface{}{
				{
					"date":  "2025-04-10",
					"value": 25120.75,
				},
				{
					"date":  "2025-04-09",
					"value": 24950.35,
				},
				{
					"date":  "2025-04-08",
					"value": 24800.12,
				},
				{
					"date":  "2025-04-07",
					"value": 24750.88,
				},
				{
					"date":  "2025-04-06",
					"value": 24600.45,
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	})

	// Account details endpoint
	http.HandleFunc("/api/v1/account/details", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"userId":   "user123",
			"email":    "user@example.com",
			"verified": true,
			"balances": []map[string]interface{}{
				{
					"symbol": "BTC",
					"free":   0.5,
					"locked": 0.0,
					"total":  0.5,
				},
				{
					"symbol": "ETH",
					"free":   5.2,
					"locked": 0.0,
					"total":  5.2,
				},
				{
					"symbol": "USDT",
					"free":   1080.40,
					"locked": 0.0,
					"total":  1080.40,
				},
			},
			"accountLevel":   "standard",
			"tradingEnabled": true,
			"createdAt":      "2023-01-15T00:00:00Z",
		}

		json.NewEncoder(w).Encode(response)
	})

	// Analytics endpoint
	http.HandleFunc("/api/v1/analytics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"totalTrades":        142,
			"successfulTrades":   98,
			"winRate":            69.01,
			"totalProfitLoss":    2347.89,
			"averageTradeProfit": 23.95,
			"largestWin":         245.67,
			"largestLoss":        87.32,
			"profitFactor":       2.75,
			"sharpeRatio":        1.85,
			"maxDrawdown":        12.5,
		}

		json.NewEncoder(w).Encode(response)
	})

	// Trades endpoints
	http.HandleFunc("/api/v1/trades", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"trades": []map[string]interface{}{
				{
					"id":        "t1",
					"symbol":    "BTC-USDT",
					"side":      "BUY",
					"price":     30864.30,
					"amount":    0.1,
					"timestamp": time.Now().Add(-24 * time.Hour).Unix(),
					"status":    "FILLED",
				},
				{
					"id":        "t2",
					"symbol":    "ETH-USDT",
					"side":      "BUY",
					"price":     1715.85,
					"amount":    2.0,
					"timestamp": time.Now().Add(-12 * time.Hour).Unix(),
					"status":    "FILLED",
				},
				{
					"id":        "t3",
					"symbol":    "ETH-USDT",
					"side":      "SELL",
					"price":     1725.30,
					"amount":    1.0,
					"timestamp": time.Now().Add(-2 * time.Hour).Unix(),
					"status":    "FILLED",
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	})

	// Config endpoint
	http.HandleFunc("/api/v1/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"tradingEnabled": true,
			"riskManagement": map[string]interface{}{
				"maxPositionSize": 0.1,
				"maxDrawdown":     0.2,
				"dailyLossLimit":  0.05,
			},
			"notifications": map[string]interface{}{
				"email":    true,
				"telegram": false,
				"slack":    false,
			},
			"strategies": []map[string]interface{}{
				{
					"id":      "newcoin",
					"name":    "New Coin Strategy",
					"enabled": true,
					"params": map[string]interface{}{
						"lookbackPeriod": 24,
						"minVolume":      100000,
					},
				},
				{
					"id":      "breakout",
					"name":    "Breakout Strategy",
					"enabled": false,
					"params": map[string]interface{}{
						"timePeriod":        14,
						"breakoutThreshold": 0.05,
					},
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	})

	// New coins endpoint
	http.HandleFunc("/api/v1/newcoins", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"newCoins": []map[string]interface{}{
				{
					"symbol":    "NEW-USDT",
					"name":      "New Coin",
					"price":     0.052,
					"listedAt":  time.Now().Add(-48 * time.Hour).Unix(),
					"volume24h": 1250000.0,
					"change24h": 15.75,
				},
				{
					"symbol":    "HYPR-USDT",
					"name":      "Hyper Protocol",
					"price":     0.38,
					"listedAt":  time.Now().Add(-36 * time.Hour).Unix(),
					"volume24h": 3500000.0,
					"change24h": 23.4,
				},
				{
					"symbol":    "META-USDT",
					"name":      "Metaverse Token",
					"price":     1.25,
					"listedAt":  time.Now().Add(-24 * time.Hour).Unix(),
					"volume24h": 5750000.0,
					"change24h": 9.8,
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	})

	// Start the server
	log.Println("Starting simple API server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
