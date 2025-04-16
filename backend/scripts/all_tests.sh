#!/bin/bash
# all_tests.sh
# Streamlined API testing script focused on MEXC data verification

# Set the base URL
BASE_URL="http://localhost:8080"
PASSED=0
FAILED=0
TOTAL=0
MOCK_DATA="NO" # Default assumption

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Server management functions
kill_existing_server() {
    echo -e "${YELLOW}Checking for existing server on port 8080...${NC}"
    local PID=$(lsof -ti:8080)
    if [ -n "$PID" ]; then
        echo -e "${YELLOW}Found server running on port 8080 (PID: $PID). Killing it...${NC}"
        kill -9 $PID
        sleep 2
        echo -e "${GREEN}Server stopped.${NC}"
    else
        echo -e "${GREEN}No server running on port 8080.${NC}"
    fi
}

# Check if MEXC credentials are configured
check_mexc_credentials() {
    echo -e "${YELLOW}Checking MEXC API credentials...${NC}"
    
    # Check if .env file exists
    if [ ! -f ".env" ]; then
        echo -e "${RED}Error: .env file not found. MEXC credentials may not be configured.${NC}"
        echo -e "${YELLOW}Creating sample .env file with placeholder MEXC credentials...${NC}"
        
        # Create a sample .env file with MEXC placeholders
        cat > .env << EOL
# MEXC API Credentials - REPLACE THESE WITH REAL VALUES
MEXC_API_KEY=your_mexc_api_key_here
MEXC_API_SECRET=your_mexc_api_secret_here
# Other environment variables
LOG_LEVEL=debug
EOL
        echo -e "${RED}WARNING: You need to edit .env and add real MEXC credentials${NC}"
        return 1
    fi
    
    # Check if MEXC credentials are in .env
    if ! grep -q "MEXC_API_KEY" .env || ! grep -q "MEXC_API_SECRET" .env; then
        echo -e "${RED}Error: MEXC credentials not found in .env file${NC}"
        return 1
    fi
    
    # Check if MEXC credentials are placeholder values
    MEXC_KEY=$(grep "MEXC_API_KEY" .env | cut -d'=' -f2)
    MEXC_SECRET=$(grep "MEXC_API_SECRET" .env | cut -d'=' -f2)
    
    if [[ "$MEXC_KEY" == *"your_mexc_api_key"* || "$MEXC_SECRET" == *"your_mexc_api_secret"* ]]; then
        echo -e "${RED}Error: MEXC credentials appear to be placeholder values${NC}"
        return 1
    fi
    
    echo -e "${GREEN}MEXC credentials found in .env file${NC}"
    return 0
}

start_server() {
    echo -e "${YELLOW}Starting backend server...${NC}"
    
    # Fixed path to server main.go
    SERVER_PATH="./cmd/server/main.go"
    
    if [ -f "$SERVER_PATH" ]; then
        # First check MEXC credentials
        check_mexc_credentials
        CREDS_STATUS=$?
        
        if [ $CREDS_STATUS -eq 1 ]; then
            echo -e "${YELLOW}Warning: MEXC credentials issue detected. May use mock data.${NC}"
        fi
        
        # Generate a valid base64-encoded 32-byte key for MEXC_CRED_ENCRYPTION_KEY
        ENCRYPTION_KEY="ZeDN1nbevBjwqlr6Zgu+JUebsgeicW6e+zqv8R0GegE="
        echo -e "${YELLOW}Setting MEXC_CRED_ENCRYPTION_KEY environment variable...${NC}"
        
        echo -e "${YELLOW}Starting server from $SERVER_PATH...${NC}"
        MEXC_CRED_ENCRYPTION_KEY="$ENCRYPTION_KEY" go run $SERVER_PATH > /tmp/server.log 2>&1 &
        SERVER_PID=$!
        echo -e "${GREEN}Server started with PID: $SERVER_PID${NC}"
    else
        echo -e "${RED}Error: Could not find server at $SERVER_PATH${NC}"
        # Attempt to find the main.go file
        MAIN_FILE=$(find ./cmd -name "main.go" | head -1)
        if [ -n "$MAIN_FILE" ]; then
            echo -e "${YELLOW}Found potential server at $MAIN_FILE. Trying to start...${NC}"
            MEXC_CRED_ENCRYPTION_KEY="ZeDN1nbevBjwqlr6Zgu+JUebsgeicW6e+zqv8R0GegE=" go run $MAIN_FILE > /tmp/server.log 2>&1 &
            SERVER_PID=$!
            echo -e "${GREEN}Server started with PID: $SERVER_PID${NC}"
        else
            echo -e "${RED}No server executable found. Please ensure the server code exists.${NC}"
            exit 1
        fi
    fi
    
    # Wait for server to initialize
    echo -e "${YELLOW}Waiting for server to initialize...${NC}"
    for i in {1..15}; do  # Increased timeout to 15 seconds
        sleep 1
        if curl -s "$BASE_URL/health" > /dev/null; then
            echo -e "${GREEN}Server is up and running!${NC}"
            return 0
        fi
        echo -n "."
    done
    
    echo -e "\n${RED}Server failed to start or health check failed. Check logs:${NC}"
    cat /tmp/server.log | tail -n 20
    exit 1
}

cleanup() {
    echo -e "\n${YELLOW}Cleaning up...${NC}"
    if [ -n "$SERVER_PID" ]; then
        echo -e "${YELLOW}Stopping server (PID: $SERVER_PID)...${NC}"
        kill -9 $SERVER_PID 2>/dev/null
        echo -e "${GREEN}Server stopped.${NC}"
    fi
    
    # Clean up temporary files
    rm -f direct_ticker.json cached_ticker.json response.json
}

# Set up trap to clean up server on script exit
trap cleanup EXIT INT TERM

# Initialize server
kill_existing_server
start_server

# Get auth token for protected endpoints
get_auth_token() {
    TOKEN_RESPONSE=$(curl -s ${BASE_URL}/api/v1/auth/test-token)
    echo $(echo $TOKEN_RESPONSE | grep -o '"token":"[^"]*' | sed 's/"token":"//')
}

# Test function that tracks results
test_endpoint() {
    ENDPOINT=$1
    NAME=$2
    METHOD=${3:-GET}
    DATA=$4
    EXPECTED_CODE=${5:-200}
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

        # For direct market endpoints, check if this is actually real MEXC data
        if [[ "$ENDPOINT" == *"/api/v1/market/direct/"* ]] && [[ "$NAME" == *"Direct"* ]]; then
            # For single ticker endpoint, look for real price data
            if [[ "$ENDPOINT" == *"/ticker/"* ]]; then
                # Try different potential JSON paths since the structure might vary
                PRICE=$(cat response.json | jq -r '.data.lastPrice // .data.price // .lastPrice // .price // "0"')
                
                if [[ -n "$PRICE" && "$PRICE" != "null" && "$PRICE" != "0" ]]; then
                    echo -e "${GREEN}✓ REAL MEXC DATA DETECTED (Price: $PRICE)${NC}"
                    # We have verified we have real data from MEXC
                    REAL_MEXC_DATA=true
                else
                    echo -e "${RED}✗ Possible mock data (no valid price found)${NC}"
                    # Mark that we might be using mock data
                    MOCK_DATA="YES"
                fi
            fi
            
            # For symbols endpoint, count symbols to determine if real data
            if [[ "$ENDPOINT" == *"/symbols"* ]]; then
                # Try different JSON paths for the symbols array
                SYMBOL_COUNT=$(cat response.json | jq -r '.data | length // 0')
                if [[ "$SYMBOL_COUNT" -gt 50 ]]; then
                    echo -e "${GREEN}✓ REAL MEXC DATA DETECTED ($SYMBOL_COUNT symbols found)${NC}"
                    REAL_MEXC_DATA=true
                else
                    echo -e "${RED}✗ Possible mock data (only $SYMBOL_COUNT symbols found)${NC}"
                    MOCK_DATA="YES"
                fi
            fi
            
            # For market data endpoints, show more detailed information
            if [[ "$ENDPOINT" == *"/ticker/"* ]]; then
                echo "Ticker data:"
                cat response.json | jq '.data'
            elif [[ "$ENDPOINT" == *"/symbols"* ]]; then
                SYMBOL_COUNT=$(cat response.json | jq '.data | length')
                echo "Total symbols: $SYMBOL_COUNT"
                echo "First 3 symbols:"
                cat response.json | jq '.data | .[0:3]'
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
        
        # If direct endpoint fails, we might be using mock data
        if [[ "$ENDPOINT" == *"/api/v1/market/direct/"* ]]; then
            MOCK_DATA="YES" 
        fi
    fi
    echo ""
}

print_section_header() {
    local title=$1
    local color=$2
    
    echo -e "\n${color}=================================================================="
    echo -e "=== ${title} ==="
    echo -e "==================================================================${NC}\n"
}

print_subsection_header() {
    local title=$1
    local color=$2
    
    echo -e "\n${color}--- ${title} ---${NC}\n"
}

# Start script
echo "=================================================================="
echo "=== STREAMLINED API TESTING SUITE ==="
echo "=================================================================="
echo "Target: ${BASE_URL}"
echo "Starting tests at $(date)"
echo "=================================================================="

# Health Check - Most basic test
print_section_header "HEALTH CHECK" "$CYAN"
test_endpoint "/health" "Health Check"

# Check the API routes that actually exist
print_section_header "API ROUTES VERIFICATION" "$CYAN"
test_endpoint "/api/v1/status/services" "Service Status"
test_endpoint "/api/v1/market/tickers" "Get All Tickers"

# MEXC Connectivity Check - Verify credentials work
print_section_header "MEXC CREDENTIAL VERIFICATION" "$BLUE"
test_endpoint "/api/v1/market/direct/ticker/BTCUSDT" "Direct BTC Ticker (MEXC API)"
test_endpoint "/api/v1/market/direct/symbols" "Direct Symbols List (MEXC API)"

# Essential Cached Data Tests - Check if caching works
print_section_header "CACHED DATA VERIFICATION" "$BLUE"
test_endpoint "/api/v1/market/ticker/BTCUSDT" "Cached BTC Ticker"
test_endpoint "/api/v1/market/symbols" "Cached Symbols"

# Direct vs Cached Comparison - Detect if using mock data
print_section_header "DIRECT VS CACHED COMPARISON" "$PURPLE"
echo -e "${YELLOW}Comparing direct vs cached data to detect mock data usage...${NC}"

# Cache a copy of direct ticker response
curl -s ${BASE_URL}/api/v1/market/direct/ticker/BTCUSDT -o direct_ticker.json
DIRECT_PRICE=$(cat direct_ticker.json | jq -r '.data.lastPrice // .data.price // .lastPrice // .price // "0"')

# Cache a copy of cached ticker response
curl -s ${BASE_URL}/api/v1/market/ticker/BTCUSDT -o cached_ticker.json
CACHED_PRICE=$(cat cached_ticker.json | jq -r '.data.Price // .data.price // .Price // .price // "0"')

# Compare prices to detect mock data
echo "Direct API price: $DIRECT_PRICE"
echo "Cached data price: $CACHED_PRICE"

if [[ "$DIRECT_PRICE" != "$CACHED_PRICE" && "$CACHED_PRICE" == "40000" ]]; then
    echo -e "${RED}✗ CACHED DATA APPEARS TO BE MOCK DATA (hardcoded 40000)${NC}"
    MOCK_DATA="YES"
elif [[ "$DIRECT_PRICE" != "0" && "$CACHED_PRICE" != "0" && "$DIRECT_PRICE" != "$CACHED_PRICE" ]]; then
    echo -e "${YELLOW}⚠ DIFFERENT PRICES BUT NOT HARDCODED (might be slightly delayed cache)${NC}"
elif [[ "$DIRECT_PRICE" == "0" || "$CACHED_PRICE" == "0" ]]; then
    echo -e "${RED}✗ ONE OR BOTH PRICES ARE MISSING OR ZERO${NC}"
    MOCK_DATA="YES"
else
    echo -e "${GREEN}✓ DIRECT AND CACHED PRICES MATCH (real data)${NC}"
fi

# Perform a simple candle test
print_section_header "CANDLE DATA TEST" "$BLUE"
test_endpoint "/api/v1/market/direct/candles/BTCUSDT/1h?limit=5" "Direct BTC Candles"

# Error Handling Test
print_section_header "ERROR HANDLING" "$RED"
test_endpoint "/invalid/route" "Invalid Route" "GET" "" 404
test_endpoint "/api/v1/market/direct/ticker/INVALID_SYMBOL" "Invalid Symbol Test"

# Performance Test
print_section_header "PERFORMANCE TEST" "$CYAN"
echo "Testing BTC ticker endpoint response time..."
time curl -s ${BASE_URL}/api/v1/market/direct/ticker/BTCUSDT > /dev/null

# Test Results Summary
echo "=================================================================="
echo "TEST RESULTS SUMMARY:"
echo -e "${GREEN}Passed: ${PASSED}/${TOTAL}${NC}"
if [ "$FAILED" -gt 0 ]; then
    echo -e "${RED}Failed: ${FAILED}/${TOTAL}${NC}"
fi

# Mock Data Assessment
echo "=================================================================="
echo "MOCK DATA ASSESSMENT:"
if [ "$MOCK_DATA" == "YES" ]; then
    echo -e "${RED}Using mock data: YES${NC}"
    echo -e "${YELLOW}Possible reasons:${NC}"
    echo "1. MEXC API credentials are missing or invalid (check your .env file)"
    echo "2. MEXC API might be down or rate-limited"
    echo "3. Network connectivity issues to MEXC API"
    echo "4. The direct API endpoints are not properly implemented"
    echo -e "${YELLOW}How to fix:${NC}"
    echo "1. Ensure valid MEXC API credentials in your .env file"
    echo "2. Check network connectivity to api.mexc.com"
    echo "3. Review the market service implementation"
else
    echo -e "${GREEN}Using mock data: NO${NC}"
    echo "Evidence:"
    echo "- Direct API tests succeeded with real-looking data"
    echo "- Price data appears to be current market values"
    echo "- Symbol count and data patterns match real MEXC API"
fi
echo "=================================================================="

if [ "$FAILED" -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi 