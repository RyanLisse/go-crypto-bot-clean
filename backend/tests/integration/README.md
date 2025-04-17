# Integration Tests

This directory contains integration test applications that test the interaction with external services and APIs.

## Test Applications

- `mexc_api`: Tests direct interaction with the MEXC API
- `direct_api`: Tests the MEXC client through the domain port interface
- `mexc_client`: Tests the MEXC client initialization and basic functionality
- `market_data`: Tests market data retrieval from MEXC
- `mexc_api_server`: Runs a test server that exposes MEXC API functionality through HTTP endpoints

## Running Tests

Each test can be run individually:

```bash
cd tests/integration/mexc_api
go run main.go
```

Make sure to set the required environment variables:

```bash
export MEXC_API_KEY=your_api_key
export MEXC_SECRET_KEY=your_secret_key
```

Or create a `.env` file in the project root with these variables.
