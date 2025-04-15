#!/bin/bash

echo "=== Testing Direct Endpoints ==="

# Test health endpoint
echo "Testing GET /health"
response=$(curl -s http://localhost:8080/health)
if [[ $response == *"ok"* ]]; then
  echo "Success: Health endpoint works"
  echo "$response"
else
  echo "Error: Health endpoint failed"
  echo "$response"
fi

# Test account test root endpoint
echo -e "\nTesting GET /account-test-root"
response=$(curl -s http://localhost:8080/account-test-root)
if [[ $response == *"success"* ]]; then
  echo "Success: Account test root endpoint works"
  echo "$response"
else
  echo "Error: Account test root endpoint failed"
  echo "$response"
fi

# Test direct wallet endpoint
echo -e "\nTesting GET /api/v1/account/wallet-direct"
response=$(curl -s http://localhost:8080/api/v1/account/wallet-direct)
if [[ $response == *"success"* ]]; then
  echo "Success: Direct wallet endpoint works"
  echo "$response"
else
  echo "Error: Direct wallet endpoint failed"
  echo "$response"
fi

# Test market tickers endpoint
echo -e "\nTesting GET /api/v1/market/tickers"
response=$(curl -s http://localhost:8080/api/v1/market/tickers)
if [[ $response == *"success"* ]]; then
  echo "Success: Market tickers endpoint works"
  echo "$response"
else
  echo "Error: Market tickers endpoint failed"
  echo "$response"
fi

echo -e "\n=== Testing Complete ==="
