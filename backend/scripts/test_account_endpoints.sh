#!/bin/bash

# Set base URL
BASE_URL="http://localhost:8080/api/v1"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Function to test an endpoint
test_endpoint() {
    local endpoint=$1
    local method=${2:-GET}
    local data=$3

    echo -e "${YELLOW}Testing ${method} ${endpoint}${NC}"

    if [ "$method" = "GET" ]; then
        response=$(curl -s -X GET "${BASE_URL}${endpoint}")
    elif [ "$method" = "POST" ]; then
        response=$(curl -s -X POST "${BASE_URL}${endpoint}" -H "Content-Type: application/json" -d "${data}")
    elif [ "$method" = "PUT" ]; then
        response=$(curl -s -X PUT "${BASE_URL}${endpoint}" -H "Content-Type: application/json" -d "${data}")
    elif [ "$method" = "DELETE" ]; then
        response=$(curl -s -X DELETE "${BASE_URL}${endpoint}")
    fi

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

# Test account endpoints
echo "=== Testing Account Endpoints ==="
test_endpoint "/account/wallet"
test_endpoint "/account/balance/BTC"
test_endpoint "/account/balance/ETH?days=7"
test_endpoint "/account/balance/USDT?days=30"
test_endpoint "/account/refresh" "POST"

# Test account-test endpoint
test_endpoint "/account-test"

echo "=== Testing Complete ==="
