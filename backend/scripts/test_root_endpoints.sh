#!/bin/bash

# Function to test an endpoint
test_endpoint() {
  local endpoint=$1
  local method=${2:-"GET"}
  local url="http://localhost:8080${endpoint}"
  
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

echo "=== Testing Root Endpoints ==="

# Test root endpoints
test_endpoint "/health"
test_endpoint "/root-test"
test_endpoint "/wallet-root"
test_endpoint "/wallet-test"
test_endpoint "/test"
test_endpoint "/account-test-root"

echo "=== Testing Complete ==="
