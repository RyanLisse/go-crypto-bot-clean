#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

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
  
  echo "Response: $response"
  echo ""
}

echo "=== Testing Account Endpoints ==="

# Test account endpoints
test_endpoint "/account-test"
test_endpoint "/account-wallet-test"
test_endpoint "/account/wallet"
test_endpoint "/account/balance/BTC"
test_endpoint "/account/balance/ETH?days=7"
test_endpoint "/account/balance/USDT?days=30"
test_endpoint "/account/refresh" "POST"

echo "=== Testing Complete ==="
