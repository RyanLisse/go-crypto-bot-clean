#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting crypto bot backend server...${NC}"

# Load environment variables from .env file
if [ -f ".env" ]; then
  echo -e "${GREEN}Loading environment variables from .env file...${NC}"
  set -a
  source .env
  set +a
else
  echo -e "${RED}Error: .env file not found. Please create an .env file with required configuration.${NC}"
  exit 1
fi

# Check for MEXC credentials
if [ -z "$MEXC_API_KEY" ] || [ -z "$MEXC_SECRET_KEY" ]; then
  echo -e "${RED}Error: MEXC API credentials are not properly configured.${NC}"
  echo "Please set MEXC_API_KEY and MEXC_SECRET_KEY in your .env file."
  exit 1
fi

if [ -z "$MEXC_CRED_ENCRYPTION_KEY" ]; then
  echo -e "${RED}Error: MEXC_CRED_ENCRYPTION_KEY is not set in your .env file.${NC}"
  echo "This key is required for encrypting and decrypting API credentials."
  exit 1
fi

echo -e "${GREEN}MEXC credentials found:${NC}"
echo "  API Key: ${MEXC_API_KEY:0:5}...${MEXC_API_KEY: -4}"
echo "  API Secret: ${MEXC_SECRET_KEY:0:5}...${MEXC_SECRET_KEY: -4}"

# Check if Turso is configured
if [ -z "$TURSO_URL" ] || [ -z "$TURSO_AUTH_TOKEN" ]; then
  echo -e "${RED}Error: Turso is not properly configured.${NC}"
  echo "Please set TURSO_URL and TURSO_AUTH_TOKEN in your .env file."
  exit 1
fi

# Configure Turso settings
export TURSO_SYNC_INTERVAL_SECONDS="300"  # 5 minutes
export TURSO_SYNC_ENABLED="true"

# Create data directory for Turso
mkdir -p ./data/turso

# Build the server if not already built
if [ ! -f "./server" ] || [ "$1" == "--rebuild" ]; then
  echo -e "${YELLOW}Building server...${NC}"
  go build -o server cmd/server/main.go
  if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to build server.${NC}"
    exit 1
  fi
  echo -e "${GREEN}Server built successfully.${NC}"
fi

# Run the server
echo -e "${GREEN}Starting server...${NC}"
./server
