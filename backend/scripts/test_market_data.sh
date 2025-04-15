#!/bin/bash

# Define base URL
BASE_URL="http://localhost:8081"

# Function to test an endpoint
test_endpoint() {
    local endpoint=$1
    local description=$2
    
    echo "Testing $description..."
    
    # Use curl to make the request and display response
    response=$(curl -s "$BASE_URL$endpoint")
    
    # Check if the response contains "success":true
    if echo "$response" | grep -q '"success":true'; then
        echo "✅ SUCCESS: $description"
    else
        echo "❌ FAILURE: $description"
        echo "Response: $response"
    fi
    
    echo ""
}

# Test market data endpoints
echo "===== Testing Market Data Endpoints ====="

test_endpoint "/api/v1/market/tickers" "Get all tickers"
test_endpoint "/api/v1/market/ticker/BTCUSDT" "Get ticker for BTC"
test_endpoint "/api/v1/market/ticker?symbol=ETHUSDT" "Get ticker for ETH (query param)"
test_endpoint "/api/v1/market/orderbook/BTCUSDT" "Get order book for BTC"
test_endpoint "/api/v1/market/candles/BTCUSDT/1d?limit=5" "Get candles for BTC (daily, 5 days)"
test_endpoint "/api/v1/market/symbols" "Get all symbols"

# Test direct API endpoints
echo "===== Testing Direct API Endpoints ====="

test_endpoint "/api/v1/market/direct/ticker/BTCUSDT" "Get ticker directly from API"
test_endpoint "/api/v1/market/direct/orderbook/BTCUSDT" "Get order book directly from API"
test_endpoint "/api/v1/market/direct/symbols" "Get symbols directly from API"
test_endpoint "/api/v1/market/direct/candles/BTCUSDT/1d?limit=5" "Get candles directly from API"

echo "===== All tests completed =====" 