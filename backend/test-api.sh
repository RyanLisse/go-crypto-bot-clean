#!/bin/bash

# Define colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Base URL
BASE_URL="http://localhost:8080/api/v1"

# Function to test an endpoint
test_endpoint() {
  local endpoint=$1
  local url="${BASE_URL}${endpoint}"
  
  echo -e "Testing ${GREEN}GET ${url}${NC}"
  
  response=$(curl -s "$url")
  
  # Check if response is valid JSON
  if echo "$response" | jq . > /dev/null 2>&1; then
    echo -e "${GREEN}Success! Response:${NC}"
    echo "$response" | jq .
  else
    echo -e "${RED}Error: Invalid JSON response${NC}"
    echo "$response"
  fi
  
  echo ""
}

# Test basic endpoints
echo "=== Testing Basic Endpoints ==="
test_endpoint "/status"
test_endpoint "/health"

# Test account endpoints
echo "=== Testing Account Endpoints ==="
test_endpoint "/account/wallet"
test_endpoint "/portfolio"
test_endpoint "/portfolio/value"
test_endpoint "/wallets"

echo "=== Testing Complete ==="
