#!/bin/bash

# Test script for MEXC API klines endpoint with different interval formats
# This script tests the MEXC API klines endpoint with different interval formats

# Set the base URL
BASE_URL="http://localhost:8080/api/v1/mexc"

# Test different interval formats
echo "Testing klines endpoint with interval=1h..."
curl -s "${BASE_URL}/klines/BTCUSDT/1h" | jq .

echo -e "\nTesting klines endpoint with interval=1hour..."
curl -s "${BASE_URL}/klines/BTCUSDT/1hour" | jq .

echo -e "\nTesting klines endpoint with interval=60m..."
curl -s "${BASE_URL}/klines/BTCUSDT/60m" | jq .

echo -e "\nTesting klines endpoint with interval=60..."
curl -s "${BASE_URL}/klines/BTCUSDT/60" | jq .

echo -e "\nTesting klines endpoint with interval=1d..."
curl -s "${BASE_URL}/klines/BTCUSDT/1d" | jq .
