#!/bin/bash

# Set colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Running Tests with Test Environment ===${NC}"

# Load test environment variables
export $(grep -v '^#' .env.test | xargs)

# Set additional test-specific environment variables
export ENV=test
export GO_ENV=test

# Run the tests
echo -e "${YELLOW}Running tests...${NC}"
go test ./... -v

# Check the exit code
if [ $? -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed.${NC}"
    exit 1
fi
