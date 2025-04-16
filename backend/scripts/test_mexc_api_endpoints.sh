#!/bin/bash

# Test script for MEXC API endpoints
# This script tests the MEXC API endpoints exposed by the backend

# Set the base URL
BASE_URL="http://localhost:8080/api/v1/mexc"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Function to make a request and check the response
test_endpoint() {
    local endpoint=$1
    local description=$2

    echo -e "${YELLOW}Testing: ${description} (${endpoint})${NC}"

    # Make the request
    response=$(curl -s "${BASE_URL}${endpoint}")

    # Check if the response contains "success":true
    if echo "$response" | grep -q '"success":true'; then
        echo -e "${GREEN}✓ Success: ${description}${NC}"
        # Print a sample of the response (first 200 characters)
        echo -e "Response sample: ${response:0:200}...\n"
        return 0
    else
        echo -e "${RED}✗ Failed: ${description}${NC}"
        echo -e "Response: ${response}\n"
        return 1
    fi
}

# Main function to run all tests
run_tests() {
    echo -e "${YELLOW}Starting MEXC API endpoint tests...${NC}\n"

    # Test account endpoint
    test_endpoint "/account" "Get account information"

    # Test market data endpoints
    test_endpoint "/ticker/BTCUSDT" "Get ticker for BTCUSDT"
    test_endpoint "/ticker/ETHUSDT" "Get ticker for ETHUSDT"
    test_endpoint "/ticker/SOLUSDT" "Get ticker for SOLUSDT"

    # Test order book endpoints
    test_endpoint "/orderbook/BTCUSDT" "Get order book for BTCUSDT"
    test_endpoint "/orderbook/ETHUSDT?depth=5" "Get order book for ETHUSDT with depth=5"

    # Test klines endpoints
    test_endpoint "/klines/BTCUSDT/60m" "Get 1h klines for BTCUSDT"
    test_endpoint "/klines/ETHUSDT/4h?limit=5" "Get 4h klines for ETHUSDT with limit=5"

    # Test exchange info endpoint
    test_endpoint "/exchange-info" "Get exchange information"

    # Test symbol info endpoint
    test_endpoint "/symbol/BTCUSDT" "Get symbol information for BTCUSDT"

    # New listings endpoint is not available in the MEXC API
    # test_endpoint "/new-listings" "Get new listings"

    echo -e "${GREEN}All tests completed!${NC}"
}

# Run the tests
run_tests
