#!/bin/bash

# Test script to verify real wallet data is being returned

echo "=== Testing Real Wallet Data ==="

# Test wallet endpoint
echo "Testing GET /api/v1/account/wallet"
response=$(curl -s http://localhost:8082/api/v1/account/wallet)
echo "Response: $response"
echo ""

# Test refresh endpoint
echo "Testing POST /api/v1/account/refresh"
response=$(curl -s -X POST http://localhost:8082/api/v1/account/refresh)
echo "Response: $response"
echo ""

# Test wallet endpoint again after refresh
echo "Testing GET /api/v1/account/wallet (after refresh)"
response=$(curl -s http://localhost:8082/api/v1/account/wallet)
echo "Response: $response"
echo ""

echo "=== Testing Complete ==="
