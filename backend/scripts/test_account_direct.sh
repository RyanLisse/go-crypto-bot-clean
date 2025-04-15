#!/bin/bash

# Set base URL (default to localhost:8081)
BASE_URL=${1:-"http://localhost:8081"}

# Function to test an endpoint
function test_endpoint() {
    ENDPOINT=$1
    METHOD=${2:-"GET"}
    
    echo "Testing $METHOD $ENDPOINT"
    echo "URL: $BASE_URL$ENDPOINT"
    
    if [ "$METHOD" = "GET" ]; then
        curl -s -X GET "$BASE_URL$ENDPOINT" | jq
    elif [ "$METHOD" = "POST" ]; then
        curl -s -X POST "$BASE_URL$ENDPOINT" | jq
    fi
    
    echo -e "\n------------------------\n"
}

# Test endpoints
echo "Testing account endpoints..."

# Update BASE_URL to use port 8081
BASE_URL="http://localhost:8081"

# Test the wallet endpoint which is working correctly
test_endpoint "/api/v1/account/wallet"

# Also test the direct test endpoints which should work
test_endpoint "/api/v1/account-test"
test_endpoint "/api/v1/account-wallet-test"

echo "Testing complete."
