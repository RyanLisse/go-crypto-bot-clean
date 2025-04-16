#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== MEXC API Configuration Check ===${NC}"

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo -e "${RED}Error: .env file not found!${NC}"
    echo "Please create a .env file with MEXC_API_KEY, MEXC_SECRET_KEY, and MEXC_CRED_ENCRYPTION_KEY."
    exit 1
fi

# Check if MEXC credentials are in .env
if ! grep -q "MEXC_API_KEY" .env; then
    echo -e "${RED}Error: MEXC_API_KEY not found in .env file!${NC}"
    echo "Please add MEXC_API_KEY to your .env file."
    exit 1
fi

if ! grep -q "MEXC_SECRET_KEY" .env; then
    echo -e "${RED}Error: MEXC_SECRET_KEY not found in .env file!${NC}"
    echo "Please add MEXC_SECRET_KEY to your .env file."
    exit 1
fi

if ! grep -q "MEXC_CRED_ENCRYPTION_KEY" .env; then
    echo -e "${RED}Error: MEXC_CRED_ENCRYPTION_KEY not found in .env file!${NC}"
    echo "Please add MEXC_CRED_ENCRYPTION_KEY to your .env file."
    exit 1
fi

# Extract values from .env
MEXC_API_KEY=$(grep "MEXC_API_KEY" .env | cut -d'=' -f2)
MEXC_SECRET_KEY=$(grep "MEXC_SECRET_KEY" .env | cut -d'=' -f2)
MEXC_CRED_ENCRYPTION_KEY=$(grep "MEXC_CRED_ENCRYPTION_KEY" .env | cut -d'=' -f2)

# Check if values are empty or placeholder
if [ -z "$MEXC_API_KEY" ] || [[ "$MEXC_API_KEY" == *"your_mexc_api_key"* ]]; then
    echo -e "${RED}Error: MEXC_API_KEY is empty or appears to be a placeholder value!${NC}"
    echo "Please set a valid MEXC_API_KEY in your .env file."
    exit 1
fi

if [ -z "$MEXC_SECRET_KEY" ] || [[ "$MEXC_SECRET_KEY" == *"your_mexc_secret"* ]]; then
    echo -e "${RED}Error: MEXC_SECRET_KEY is empty or appears to be a placeholder value!${NC}"
    echo "Please set a valid MEXC_SECRET_KEY in your .env file."
    exit 1
fi

if [ -z "$MEXC_CRED_ENCRYPTION_KEY" ]; then
    echo -e "${RED}Error: MEXC_CRED_ENCRYPTION_KEY is empty!${NC}"
    echo "Please set a valid MEXC_CRED_ENCRYPTION_KEY in your .env file."
    exit 1
fi

# Check config.yaml
if [ ! -f "configs/config.yaml" ]; then
    echo -e "${RED}Error: configs/config.yaml file not found!${NC}"
    echo "Please make sure the config file exists."
    exit 1
fi

# Export environment variables so they are available to the application
export MEXC_API_KEY
export MEXC_SECRET_KEY
export MEXC_CRED_ENCRYPTION_KEY

echo -e "${GREEN}MEXC configuration check completed successfully!${NC}"
echo "API Key: ${MEXC_API_KEY:0:5}...${MEXC_API_KEY: -4}"
echo "Secret Key: ${MEXC_SECRET_KEY:0:5}...${MEXC_SECRET_KEY: -4}"
echo "Encryption Key: ${MEXC_CRED_ENCRYPTION_KEY:0:10}...${MEXC_CRED_ENCRYPTION_KEY: -4}"

echo -e "${YELLOW}Testing MEXC client creation...${NC}"
go run cmd/test_mexc_client/main.go
exit_code=$?

if [ $exit_code -eq 0 ]; then
    echo -e "${GREEN}MEXC client test successful!${NC}"
else
    echo -e "${RED}MEXC client test failed with exit code $exit_code${NC}"
    echo "Please check your MEXC credentials and encryption key."
fi

exit $exit_code 