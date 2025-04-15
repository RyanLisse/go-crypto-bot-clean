#!/bin/bash
# endpoint_test.sh
# Script to verify all endpoints are working after the Chi router migration

# Set the base URL (modify as needed)
BASE_URL="http://localhost:8080"
PASSED=0
FAILED=0
TOTAL=0

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Get auth token for protected endpoints
get_auth_token() {
    TOKEN_RESPONSE=$(curl -s ${BASE_URL}/api/v1/auth/test-token)
    echo $(echo $TOKEN_RESPONSE | grep -o '"token":"[^"]*' | sed 's/"token":"//')
}

# Test function that tracks results
test_endpoint() {
    ENDPOINT=$1
    METHOD=${2:-GET}
    DATA=$3
    EXPECTED_CODE=${4:-200}
    NAME=$5
    PROTECTED=${6:-false}

    TOTAL=$((TOTAL+1))

    echo -e "${YELLOW}Testing ${METHOD} ${ENDPOINT}${NC}"

    # Add authentication header for protected endpoints
    if [ "$PROTECTED" = "true" ]; then
        if [ -z "$AUTH_TOKEN" ]; then
            AUTH_TOKEN=$(get_auth_token)
            if [ -z "$AUTH_TOKEN" ]; then
                echo -e "${RED}Failed to get auth token, skipping protected endpoint test${NC}"
                FAILED=$((FAILED+1))
                return
            fi
        fi
        AUTH_HEADER="-H \"Authorization: Bearer $AUTH_TOKEN\""
    else
        AUTH_HEADER=""
    fi

    if [ "$METHOD" = "GET" ]; then
        if [ "$PROTECTED" = "true" ]; then
            RESPONSE=$(curl -s -o response.json -w "%{http_code}" -X ${METHOD} -H "Authorization: Bearer $AUTH_TOKEN" ${BASE_URL}${ENDPOINT})
        else
            RESPONSE=$(curl -s -o response.json -w "%{http_code}" -X ${METHOD} ${BASE_URL}${ENDPOINT})
        fi
    else
        if [ "$PROTECTED" = "true" ]; then
            RESPONSE=$(curl -s -o response.json -w "%{http_code}" -X ${METHOD} -H "Authorization: Bearer $AUTH_TOKEN" -H "Content-Type: application/json" -d "${DATA}" ${BASE_URL}${ENDPOINT})
        else
            RESPONSE=$(curl -s -o response.json -w "%{http_code}" -X ${METHOD} -H "Content-Type: application/json" -d "${DATA}" ${BASE_URL}${ENDPOINT})
        fi
    fi

    if [ "$RESPONSE" -eq "$EXPECTED_CODE" ]; then
        echo -e "${GREEN}✓ Success: ${NAME} (HTTP ${RESPONSE})${NC}"
        PASSED=$((PASSED+1))

        # Show response with better formatting
        echo "Response:"

        # For market data endpoints, show more detailed information
        if [[ "$ENDPOINT" == *"/api/v1/market/"* ]]; then
            echo "=== REAL MEXC DATA ==="

            # For tickers endpoint, show count and first few tickers
            if [[ "$ENDPOINT" == *"/tickers"* ]]; then
                COUNT=$(cat response.json | jq '.data | length')
                echo "Total tickers: $COUNT"
                echo "First 3 tickers:"
                cat response.json | jq '.data | .[0:3]'

            # For single ticker endpoint, show all data
            elif [[ "$ENDPOINT" == *"/ticker/"* ]]; then
                echo "Ticker data:"
                cat response.json | jq '.data'

            # For orderbook endpoint, show summary
            elif [[ "$ENDPOINT" == *"/orderbook/"* ]]; then
                BID_COUNT=$(cat response.json | jq '.data.bids | length')
                ASK_COUNT=$(cat response.json | jq '.data.asks | length')
                echo "Order book summary:"
                echo "Total bids: $BID_COUNT"
                echo "Total asks: $ASK_COUNT"
                echo "Top 3 bids:"
                cat response.json | jq '.data.bids | .[0:3]'
                echo "Top 3 asks:"
                cat response.json | jq '.data.asks | .[0:3]'

            # For candles endpoint, show summary
            elif [[ "$ENDPOINT" == *"/candles/"* ]]; then
                CANDLE_COUNT=$(cat response.json | jq '.data | length')
                echo "Total candles: $CANDLE_COUNT"
                echo "First 2 candles:"
                cat response.json | jq '.data | .[0:2]'

            # For symbols endpoint, show count and first few symbols
            elif [[ "$ENDPOINT" == *"/symbols"* ]]; then
                SYMBOL_COUNT=$(cat response.json | jq '.data | length')
                echo "Total symbols: $SYMBOL_COUNT"
                echo "First 3 symbols:"
                cat response.json | jq '.data | .[0:3]'

            # Default case
            else
                cat response.json | jq . | head -n 15
                LINES=$(cat response.json | jq . | wc -l)
                if [ "$LINES" -gt 15 ]; then
                    echo "... (truncated, $LINES lines total)"
                fi
            fi
        else
            # Default behavior for non-market endpoints
            cat response.json | jq . | head -n 10
            LINES=$(cat response.json | jq . | wc -l)
            if [ "$LINES" -gt 10 ]; then
                echo "... (truncated, $LINES lines total)"
            fi
        fi
    else
        echo -e "${RED}✗ Failed: ${NAME} - Expected HTTP ${EXPECTED_CODE}, got ${RESPONSE}${NC}"
        FAILED=$((FAILED+1))
        echo "Response:"
        cat response.json | jq .
    fi
    echo ""

    # Return the ID if this is a POST request creating a resource
    if [ "$METHOD" = "POST" ]; then
        ID=$(cat response.json | jq -r '.id' 2>/dev/null)
        if [ "$?" -eq 0 ] && [ -n "$ID" ]; then
            echo "$ID"
        fi
    fi
}

echo "=== Starting Endpoint Tests ==="
echo "Target: ${BASE_URL}"
echo "=================================================================="

# Health Check
test_endpoint "/health" "GET" "" 200 "Health Check"

# Market Data Endpoints
echo "\n=== Testing Market Data Endpoints (Real MEXC Data) ==="

echo "\n--- Testing Get All Tickers ---"
test_endpoint "/api/v1/market/tickers" "GET" "" 200 "Get All Tickers"

echo "\n--- Testing Get Single Ticker (BTC) ---"
test_endpoint "/api/v1/market/ticker/BTCUSDT" "GET" "" 200 "Get Single Ticker (BTC)"

echo "\n--- Testing Get Single Ticker (ETH) ---"
test_endpoint "/api/v1/market/ticker/ETHUSDT" "GET" "" 200 "Get Single Ticker (ETH)"

echo "\n--- Testing Get Order Book ---"
test_endpoint "/api/v1/market/orderbook/BTCUSDT" "GET" "" 200 "Get Order Book"

echo "\n--- Testing Get Candles (1h) ---"
test_endpoint "/api/v1/market/candles/BTCUSDT/1h?limit=10" "GET" "" 200 "Get Candles (1h)"

echo "\n--- Testing Get Candles (15m) ---"
test_endpoint "/api/v1/market/candles/ETHUSDT/15m?limit=5" "GET" "" 200 "Get Candles (15m)"

echo "\n--- Testing Get Symbols ---"
test_endpoint "/api/v1/market/symbols" "GET" "" 200 "Get Symbols"

# Status Endpoints
test_endpoint "/api/v1/status/services" "GET" "" 200 "Get Services Status"
test_endpoint "/api/v1/status/exchange" "GET" "" 200 "Get Exchange Status"

# Alert Workflow
echo "=== Testing Alert Workflow ==="

# Get auth token for protected endpoints
AUTH_TOKEN=$(get_auth_token)
if [ -z "$AUTH_TOKEN" ]; then
    echo -e "${RED}Failed to get auth token, skipping alert workflow tests${NC}"
    FAILED=$((FAILED+6))
    TOTAL=$((TOTAL+6))
    echo ""
else
    echo -e "${GREEN}Successfully obtained auth token for protected endpoints${NC}"

    # Test getting all alerts
    test_endpoint "/api/v1/alerts" "GET" "" 200 "Get All Alerts" "true"

    # Create alert
    ALERT_DATA='{
      "symbol": "BTCUSDT",
      "condition": "price_above",
      "threshold": 40000,
      "userId": "user_2NNPBn8mSWz5KXFMDq9UzCVAq1t"
    }'

    # Create alert and extract ID directly
    RESPONSE=$(curl -s -X POST -H "Authorization: Bearer $AUTH_TOKEN" -H "Content-Type: application/json" -d "$ALERT_DATA" ${BASE_URL}/api/v1/alerts)
    echo "Alert creation response: $RESPONSE"

    # Get all alerts to find the one we just created
    ALL_ALERTS=$(curl -s -H "Authorization: Bearer $AUTH_TOKEN" ${BASE_URL}/api/v1/alerts)
    echo "All alerts: $ALL_ALERTS"

    # Count alerts to verify creation
    ALERT_COUNT=$(echo $ALL_ALERTS | jq '. | length')
    if [ "$ALERT_COUNT" -gt 0 ]; then
        echo -e "${GREEN}✓ Success: Alert creation verified (found $ALERT_COUNT alerts)${NC}"
        PASSED=$((PASSED+1))
    else
        echo -e "${RED}✗ Failed: Alert creation could not be verified${NC}"
        FAILED=$((FAILED+1))
    fi
    TOTAL=$((TOTAL+1))
fi

# Error Handling Tests
echo "=== Testing Error Handling ==="

# Invalid route
test_endpoint "/invalid/route" "GET" "" 404 "Invalid Route"

# Invalid parameter
test_endpoint "/api/v1/market/ticker/INVALID_SYMBOL" "GET" "" 400 "Invalid Symbol"

# Performance Tests
echo "=== Testing Performance ==="

echo "Testing health endpoint performance..."
time curl -s ${BASE_URL}/health > /dev/null

echo "Testing tickers endpoint performance..."
time curl -s ${BASE_URL}/api/v1/market/tickers > /dev/null

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