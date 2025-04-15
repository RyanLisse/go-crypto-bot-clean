#!/bin/bash
# direct_endpoint_test.sh
# Script to test direct MEXC API endpoints which are known to work

# Set the base URL
BASE_URL="http://localhost:8080"
PASSED=0
FAILED=0
TOTAL=0

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Test function that tracks results
test_endpoint() {
    ENDPOINT=$1
    NAME=$2
    EXPECTED_CODE=${3:-200}

    TOTAL=$((TOTAL+1))

    echo -e "${YELLOW}Testing ${ENDPOINT}${NC}"

    RESPONSE=$(curl -s -o response.json -w "%{http_code}" ${BASE_URL}${ENDPOINT})

    if [ "$RESPONSE" -eq "$EXPECTED_CODE" ]; then
        echo -e "${GREEN}✓ Success: ${NAME} (HTTP ${RESPONSE})${NC}"
        PASSED=$((PASSED+1))

        # Show response with better formatting
        echo "Response:"
        cat response.json | jq . | head -n 20
        LINES=$(cat response.json | jq . | wc -l)
        if [ "$LINES" -gt 20 ]; then
            echo "... (truncated, $LINES lines total)"
        fi
    else
        echo -e "${RED}✗ Failed: ${NAME} - Expected HTTP ${EXPECTED_CODE}, got ${RESPONSE}${NC}"
        FAILED=$((FAILED+1))
        echo "Response:"
        cat response.json | jq .
    fi
    echo ""
}

echo "=== Starting Direct Endpoint Tests ==="
echo "Target: ${BASE_URL}"
echo "=================================================================="

# Health Check
test_endpoint "/health" "Health Check"

# Market Data Direct Endpoints
echo "\n=== Testing Direct Market Data Endpoints (Real MEXC Data) ==="

echo "\n--- Testing Direct Ticker (BTC) ---"
test_endpoint "/api/v1/market/direct/ticker/BTCUSDT" "Get Direct Ticker (BTC)"

echo "\n--- Testing Direct Ticker (ETH) ---"
test_endpoint "/api/v1/market/direct/ticker/ETHUSDT" "Get Direct Ticker (ETH)"

echo "\n--- Testing Direct Order Book ---"
test_endpoint "/api/v1/market/orderbook/BTCUSDT" "Get Order Book"

echo "\n--- Testing Direct Symbols ---"
test_endpoint "/api/v1/market/direct/symbols" "Get Direct Symbols"

# Test invalid route to verify error handling
echo "\n=== Testing Error Handling ==="
test_endpoint "/invalid/route" "Invalid Route" 404

# Test Results Summary
echo "=================================================================="
echo "Test Results Summary:"
echo -e "${GREEN}Passed: ${PASSED}/${TOTAL}${NC}"
if [ "$FAILED" -gt 0 ]; then
    echo -e "${RED}Failed: ${FAILED}/${TOTAL}${NC}"
fi

if [ "$FAILED" -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi 