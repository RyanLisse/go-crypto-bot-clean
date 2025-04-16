#!/bin/bash

# Simple test script for MEXC API endpoints
# This script tests the MEXC API endpoints exposed by the backend

# Set the base URL
BASE_URL="http://localhost:8080/api/v1/mexc"

# Test ticker endpoint
echo "Testing ticker endpoint..."
curl -s "${BASE_URL}/ticker/BTCUSDT" | jq .

# Test order book endpoint
echo -e "\nTesting order book endpoint..."
curl -s "${BASE_URL}/orderbook/BTCUSDT" | jq .

# Test klines endpoint
echo -e "\nTesting klines endpoint..."
curl -s "${BASE_URL}/klines/BTCUSDT/60m" | jq .

# Test exchange info endpoint
echo -e "\nTesting exchange info endpoint..."
curl -s "${BASE_URL}/exchange-info" | jq .

# Test account endpoint
echo -e "\nTesting account endpoint..."
curl -s "${BASE_URL}/account" | jq .
