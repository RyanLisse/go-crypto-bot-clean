#!/bin/bash

# This script retrieves MEXC account data and creates sample files

# Set up colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== MEXC Account Data Retrieval ===${NC}"
echo "This script will retrieve your MEXC account data and create sample files."

# Check if API keys are set
if [ -z "$MEXC_API_KEY" ] || [ -z "$MEXC_API_SECRET" ]; then
    echo -e "${YELLOW}MEXC API credentials not found in environment variables.${NC}"
    echo "Please enter your MEXC API credentials:"

    read -p "MEXC API Key: " MEXC_API_KEY
    read -p "MEXC API Secret: " MEXC_API_SECRET

    # Export for child processes
    export MEXC_API_KEY
    export MEXC_API_SECRET
fi

echo -e "${GREEN}Step 1: Retrieving basic account information...${NC}"
go run scripts/mexc/balance/get_mexc_balance.go
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to retrieve basic account information.${NC}"
    echo "Continuing with detailed retrieval..."
else
    echo -e "${GREEN}Basic account information retrieved successfully.${NC}"
fi

echo -e "${GREEN}Step 2: Retrieving detailed account information...${NC}"
go run scripts/mexc/detailed/get_mexc_balance_detailed.go
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to retrieve detailed account information.${NC}"
    echo "Continuing with sample creation..."
else
    echo -e "${GREEN}Detailed account information retrieved successfully.${NC}"
fi

echo -e "${GREEN}Step 3: Creating sample balance files...${NC}"
go run scripts/mexc/sample/create_sample_balance.go
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to create sample balance files.${NC}"
    exit 1
else
    echo -e "${GREEN}Sample balance files created successfully.${NC}"
fi

echo -e "${GREEN}=== MEXC Account Data Retrieval Complete ===${NC}"
echo "The following files have been created:"
echo "  - mexc_balance.json: Basic account information"
echo "  - mexc_balance_detailed.json: Detailed account information"
echo "  - mexc_account_raw.json: Raw account response from MEXC API"
echo "  - sample_balance.json: Sample balance file for testing"
echo "  - pkg/platform/mexc/sample_balance.go: Go file with sample balance data"

echo -e "${YELLOW}Next steps:${NC}"
echo "1. Use the sample_balance.go file in your code to get real account data"
echo "2. Update the GetAccount method in pkg/platform/mexc/client.go to use the sample data"
echo "3. Run your application to test with real account data"
