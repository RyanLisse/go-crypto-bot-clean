#!/bin/bash

BASE_URL="http://localhost:8081/api/v1"

# Function to test an endpoint
test_endpoint() {
  local endpoint=$1
  local method=${2:-"GET"}
  local url="${BASE_URL}${endpoint}"

  echo "Testing ${method} ${url}"

  if [ "$method" == "GET" ]; then
    response=$(curl -s "$url")
  else
    response=$(curl -s -X "$method" "$url")
  fi

  if [[ $response == *"success"* ]]; then
    echo "Success: $response"
  else
    echo "Error: Invalid JSON response"
    echo "$response"
  fi

  echo ""
}

echo "=== Testing API v1 Endpoints ==="

# Test account endpoints
test_endpoint "/account/wallet"
test_endpoint "/account/balance/BTC"
test_endpoint "/account/balance/ETH?days=7"
test_endpoint "/account/balance/USDT?days=30"
test_endpoint "/account/refresh" "POST"

# Test direct test endpoints
test_endpoint "/account-test-simple"
test_endpoint "/account-direct-test"
test_endpoint "/account-test"
test_endpoint "/direct-wallet"
test_endpoint "/account/wallet-direct"
test_endpoint "/wallet-test-direct"

# Test market endpoints
test_endpoint "/market/tickers"
test_endpoint "/market/ticker/BTCUSDT"

echo "=== Testing Complete ==="
