# Test Scripts

This directory contains individual test scripts for manually testing the MEXC API integration. 
These scripts were moved from the original `scripts/test` directory to separate packages to avoid package conflicts.

Each script is in its own package with a dedicated `main.go` file:

- `account_direct` - Direct test of MEXC account API
- `exchange_info` - Test of MEXC exchange information API
- `rest` - Test of MEXC REST API
- `rest_account` - Test of MEXC REST account API
- `public` - Test of MEXC public API
- `api` - General MEXC API test
- `api_alt` - Alternative MEXC API test
- `api_key` - MEXC API key functionality test

## Usage

To run any of these test scripts, navigate to its directory and run:

```bash
go run main.go
```

For example:

```bash
cd scripts/test_scripts/account_direct
go run main.go
```

Note: Most tests require valid MEXC API credentials to be set in the .env file. 