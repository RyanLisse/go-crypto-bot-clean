# Endpoint Testing Plan

This document outlines a comprehensive testing plan to verify all endpoints are working correctly following the Chi router migration and other optimizations.

## Testing Strategy

1. **Automated Tests**: Run existing test suite to verify functionality
2. **Manual API Testing**: Use tools like Postman or curl to test endpoints
3. **Integration Testing**: Test end-to-end workflows across multiple endpoints

## Endpoints to Test

Based on the codebase analysis, the following endpoints need to be verified:

### Core API Endpoints

1. **Health Check**
   - `GET /health`
   - Expected response: Status 200 with JSON containing status, version, and timestamp

2. **Market Data Endpoints**
   - `GET /api/v1/market/tickers`
   - `GET /api/v1/market/ticker/{symbol}`
   - `GET /api/v1/market/orderbook/{symbol}`
   - `GET /api/v1/market/candles/{symbol}/{interval}`
   - `GET /api/v1/market/symbols`

3. **Status Endpoints**
   - `GET /api/v1/status/services`
   - `GET /api/v1/status/exchange`

4. **Alert Endpoints**
   - `GET /api/v1/alerts`
   - `POST /api/v1/alerts`
   - `GET /api/v1/alerts/{id}`
   - `PUT /api/v1/alerts/{id}`
   - `DELETE /api/v1/alerts/{id}`

## Testing Process

### 1. Prepare Testing Environment

```bash
# Start the application in development mode
cd cmd/server
go run main.go
```

### 2. Health Check Test

```bash
# Test health endpoint
curl http://localhost:8080/health | jq
```

Expected output:
```json
{
  "status": "ok",
  "version": "0.1.0",
  "timestamp": "2023-06-10T15:04:05Z"
}
```

### 3. Market Data Tests

```bash
# Test tickers endpoint
curl http://localhost:8080/api/v1/market/tickers | jq

# Test specific ticker
curl http://localhost:8080/api/v1/market/ticker/BTCUSDT | jq

# Test order book
curl http://localhost:8080/api/v1/market/orderbook/BTCUSDT | jq

# Test candles with interval
curl "http://localhost:8080/api/v1/market/candles/BTCUSDT/1h?limit=10" | jq

# Test symbols list
curl http://localhost:8080/api/v1/market/symbols | jq
```

### 4. Status Tests

```bash
# Test service status
curl http://localhost:8080/api/v1/status/services | jq

# Test exchange status
curl http://localhost:8080/api/v1/status/exchange | jq
```

### 5. Alert Tests

```bash
# Get all alerts
curl http://localhost:8080/api/v1/alerts | jq

# Create a new alert
curl -X POST http://localhost:8080/api/v1/alerts \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "condition": "price_above",
    "threshold": 40000,
    "userId": "user123"
  }' | jq

# Get a specific alert
# Replace {id} with an actual alert ID
curl http://localhost:8080/api/v1/alerts/{id} | jq

# Update an alert
# Replace {id} with an actual alert ID
curl -X PUT http://localhost:8080/api/v1/alerts/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "condition": "price_below",
    "threshold": 35000,
    "userId": "user123"
  }' | jq

# Delete an alert
# Replace {id} with an actual alert ID
curl -X DELETE http://localhost:8080/api/v1/alerts/{id} | jq
```

## Integration Testing

Test the following end-to-end workflows:

1. **Market Data Workflow**
   - Get list of symbols
   - Get ticker for a specific symbol
   - Get order book for the same symbol
   - Get candles for the same symbol with different intervals

2. **Alert Management Workflow**
   - Create a new alert
   - Retrieve the created alert to verify its properties
   - Update the alert
   - Verify the alert was updated correctly
   - Delete the alert
   - Verify the alert was removed

## Performance Testing

For each endpoint, verify response times are within acceptable ranges:

```bash
# Example using time with curl
time curl http://localhost:8080/api/v1/market/tickers > /dev/null
```

Expected performance metrics:
- Fast endpoints (health, single ticker): < 100ms
- Medium endpoints (order book, specific candles): < 200ms
- Heavy endpoints (all tickers, multiple candles): < 500ms

## Error Handling Testing

Test error scenarios to ensure proper error responses:

1. **Invalid Routes**
   ```bash
   curl http://localhost:8080/invalid/route
   ```
   Expected: 404 Not Found with error message

2. **Invalid Parameters**
   ```bash
   curl http://localhost:8080/api/v1/market/ticker/INVALID_SYMBOL
   ```
   Expected: 400 Bad Request with specific error message

3. **Rate Limiting**
   - Make multiple rapid requests to test rate limiter
   ```bash
   for i in {1..20}; do curl http://localhost:8080/api/v1/market/tickers; done
   ```
   Expected: 429 Too Many Requests after limit is reached

## Test Reporting

After completing all tests, document results in a report with:

1. Test date and environment
2. List of passed and failed tests
3. Performance metrics
4. Issues discovered
5. Recommendations for fixes or improvements

## Automated Test Script

Create a shell script to automate the process:

```bash
#!/bin/bash
# endpoint_test.sh

# Health check
echo "Testing health endpoint..."
curl -s http://localhost:8080/health | jq

# Market data endpoints
echo "Testing market data endpoints..."
curl -s http://localhost:8080/api/v1/market/tickers | jq
curl -s http://localhost:8080/api/v1/market/ticker/BTCUSDT | jq
curl -s http://localhost:8080/api/v1/market/orderbook/BTCUSDT | jq
curl -s "http://localhost:8080/api/v1/market/candles/BTCUSDT/1h?limit=10" | jq
curl -s http://localhost:8080/api/v1/market/symbols | jq

# Status endpoints
echo "Testing status endpoints..."
curl -s http://localhost:8080/api/v1/status/services | jq
curl -s http://localhost:8080/api/v1/status/exchange | jq

# Alert workflow
echo "Testing alert workflow..."
ALERT_ID=$(curl -s -X POST http://localhost:8080/api/v1/alerts \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "condition": "price_above",
    "threshold": 40000,
    "userId": "user123"
  }' | jq -r '.id')

echo "Created alert with ID: $ALERT_ID"

curl -s http://localhost:8080/api/v1/alerts/$ALERT_ID | jq

curl -s -X PUT http://localhost:8080/api/v1/alerts/$ALERT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "condition": "price_below",
    "threshold": 35000,
    "userId": "user123"
  }' | jq

curl -s http://localhost:8080/api/v1/alerts/$ALERT_ID | jq

curl -s -X DELETE http://localhost:8080/api/v1/alerts/$ALERT_ID | jq

echo "All tests completed!"
```

Make the script executable and run:
```bash
chmod +x endpoint_test.sh
./endpoint_test.sh
``` 