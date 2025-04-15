#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Testing Clerk Authentication Endpoints${NC}"
echo "========================================"

# Base URL
BASE_URL="http://localhost:8080/api/v1"

# Test 1: Access protected endpoint without authentication
echo -e "\n${YELLOW}Test 1: Access protected endpoint without authentication${NC}"
RESPONSE=$(curl -s $BASE_URL/alerts)
if [[ "$RESPONSE" == "Unauthorized" ]]; then
    echo -e "${GREEN}✓ Test passed: Unauthorized access correctly rejected${NC}"
else
    echo -e "${RED}✗ Test failed: Expected 'Unauthorized' but got: $RESPONSE${NC}"
fi

# Test 2: Access protected endpoint with invalid token
echo -e "\n${YELLOW}Test 2: Access protected endpoint with invalid token${NC}"
RESPONSE=$(curl -s -H "Authorization: Bearer invalid_token" $BASE_URL/alerts)
if [[ "$RESPONSE" == *"Invalid authentication token"* ]]; then
    echo -e "${GREEN}✓ Test passed: Invalid token correctly rejected${NC}"
else
    echo -e "${RED}✗ Test failed: Expected error message but got: $RESPONSE${NC}"
fi

# Test 3: Get test token
echo -e "\n${YELLOW}Test 3: Get test token${NC}"
TOKEN_RESPONSE=$(curl -s $BASE_URL/auth/test-token)
TOKEN=$(echo $TOKEN_RESPONSE | grep -o '"token":"[^"]*' | sed 's/"token":"//')

if [[ -n "$TOKEN" ]]; then
    echo -e "${GREEN}✓ Test passed: Successfully retrieved test token${NC}"
else
    echo -e "${RED}✗ Test failed: Could not retrieve test token${NC}"
    exit 1
fi

# Test 4: Access protected endpoint with valid token
echo -e "\n${YELLOW}Test 4: Access protected endpoint with valid token${NC}"
RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" $BASE_URL/alerts)
if [[ "$RESPONSE" == "[]" || "$RESPONSE" == *"id"* ]]; then
    echo -e "${GREEN}✓ Test passed: Successfully accessed protected endpoint${NC}"
else
    echo -e "${RED}✗ Test failed: Could not access protected endpoint: $RESPONSE${NC}"
fi

# Test 5: Create an alert with valid token
echo -e "\n${YELLOW}Test 5: Create an alert with valid token${NC}"
ALERT_RESPONSE=$(curl -s -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
    -d '{"symbol":"BTCUSDT","condition":"price_above","threshold":45000,"userId":"user_2NNPBn8mSWz5KXFMDq9UzCVAq1t"}' \
    $BASE_URL/alerts)

ALERT_ID=$(echo $ALERT_RESPONSE | grep -o '"id":"[^"]*' | sed 's/"id":"//')

if [[ -n "$ALERT_ID" ]]; then
    echo -e "${GREEN}✓ Test passed: Successfully created an alert with ID: $ALERT_ID${NC}"
else
    echo -e "${RED}✗ Test failed: Could not create an alert: $ALERT_RESPONSE${NC}"
fi

# Test 6: Retrieve alerts with valid token
echo -e "\n${YELLOW}Test 6: Retrieve alerts with valid token${NC}"
RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" $BASE_URL/alerts)
if [[ "$RESPONSE" == *"BTCUSDT"* && "$RESPONSE" == *"price_above"* ]]; then
    echo -e "${GREEN}✓ Test passed: Successfully retrieved alerts${NC}"
else
    echo -e "${RED}✗ Test failed: Could not retrieve alerts or alert data missing: $RESPONSE${NC}"
fi

echo -e "\n${YELLOW}All tests completed!${NC}"
