#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Running tests for the entire project...${NC}"

# Run frontend tests
echo -e "\n${YELLOW}Running frontend tests...${NC}"
cd frontend
if bun run test; then
  echo -e "${GREEN}Frontend tests passed!${NC}"
  FRONTEND_TESTS_PASSED=true
else
  echo -e "${RED}Frontend tests failed!${NC}"
  FRONTEND_TESTS_PASSED=false
fi
cd ..

# Run backend tests
echo -e "\n${YELLOW}Running backend tests...${NC}"
cd backend
if go test ./...; then
  echo -e "${GREEN}Backend tests passed!${NC}"
  BACKEND_TESTS_PASSED=true
else
  echo -e "${RED}Backend tests failed!${NC}"
  BACKEND_TESTS_PASSED=false
fi
cd ..

# Check if all tests passed
if [ "$FRONTEND_TESTS_PASSED" = true ] && [ "$BACKEND_TESTS_PASSED" = true ]; then
  echo -e "\n${GREEN}All tests passed!${NC}"
  exit 0
else
  echo -e "\n${RED}Some tests failed!${NC}"
  exit 1
fi
