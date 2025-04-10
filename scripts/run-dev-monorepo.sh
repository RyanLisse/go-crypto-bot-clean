#!/bin/bash

# run-dev-monorepo.sh - Script to run backend and frontend services for the crypto bot
# Author: Ryan Lisse
# Date: 2024-04-08

# Text colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Print header
echo -e "${BLUE}=======================================${NC}"
echo -e "${BLUE}   Crypto Bot Development Environment  ${NC}"
echo -e "${BLUE}=======================================${NC}"

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Check for required tools
if ! command_exists go; then
  echo -e "${RED}Error: Go is not installed. Please install Go to run the backend.${NC}"
  exit 1
fi

if ! command_exists bun; then
  echo -e "${YELLOW}Warning: Bun is not installed. Falling back to npm for frontend.${NC}"
  USE_NPM=true
else
  USE_NPM=false
fi

# Function to handle cleanup on exit
cleanup() {
  echo -e "\n${YELLOW}Shutting down services...${NC}"

  # Kill all background processes
  if [ -n "$BACKEND_PID" ]; then
    echo -e "${YELLOW}Stopping backend (PID: $BACKEND_PID)${NC}"
    kill $BACKEND_PID 2>/dev/null
  fi

  if [ -n "$FRONTEND_PID" ]; then
    echo -e "${YELLOW}Stopping frontend (PID: $FRONTEND_PID)${NC}"
    kill $FRONTEND_PID 2>/dev/null
  fi

  echo -e "${GREEN}All services stopped. Goodbye!${NC}"
  exit 0
}

# Set up trap to catch Ctrl+C and other termination signals
trap cleanup SIGINT SIGTERM

# Check if config file exists
CONFIG_FILE="backend/configs/config.yaml"
if [ ! -f "$CONFIG_FILE" ]; then
  echo -e "${YELLOW}Warning: Config file not found at $CONFIG_FILE${NC}"
  echo -e "${YELLOW}The backend may fail to start if the config file is missing.${NC}"
fi

# Define ports
BACKEND_PORT=8080
FRONTEND_PORT=3000

# Start backend server
echo -e "\n${GREEN}Starting backend API server on port $BACKEND_PORT...${NC}"
cd backend
go run cmd/api/main.go --port=$BACKEND_PORT &
BACKEND_PID=$!
cd ..

# Check if backend started successfully
sleep 2
if ! ps -p $BACKEND_PID > /dev/null; then
  echo -e "${RED}Error: Backend failed to start.${NC}"
  exit 1
fi

echo -e "${GREEN}Backend API server running with PID: $BACKEND_PID${NC}"
echo -e "${GREEN}API server available at: http://localhost:$BACKEND_PORT${NC}"

# Start frontend
echo -e "\n${GREEN}Starting frontend development server on port $FRONTEND_PORT...${NC}"
cd frontend

# Set environment variables for the frontend
export NEXT_PUBLIC_API_URL="http://localhost:$BACKEND_PORT/api/v1"
export NEXT_PUBLIC_WS_URL="ws://localhost:$BACKEND_PORT/ws"

if [ "$USE_NPM" = true ]; then
  echo -e "${YELLOW}Using npm to start frontend...${NC}"
  npm run dev -- --port=$FRONTEND_PORT &
else
  echo -e "${GREEN}Using bun to start frontend...${NC}"
  bun run dev -- --port=$FRONTEND_PORT &
fi

FRONTEND_PID=$!
cd ..

# Check if frontend started successfully
sleep 3
if ! ps -p $FRONTEND_PID > /dev/null; then
  echo -e "${RED}Error: Frontend failed to start.${NC}"
  kill $BACKEND_PID
  exit 1
fi

echo -e "${GREEN}Frontend development server running with PID: $FRONTEND_PID${NC}"
echo -e "${GREEN}Frontend available at: http://localhost:$FRONTEND_PORT${NC}"

# Print instructions
echo -e "\n${BLUE}=======================================${NC}"
echo -e "${GREEN}All services are running!${NC}"
echo -e "${BLUE}=======================================${NC}"
echo -e "Backend API:   ${YELLOW}http://localhost:$BACKEND_PORT${NC}"
echo -e "Frontend:      ${YELLOW}http://localhost:$FRONTEND_PORT${NC}"
echo -e "\nPress ${RED}Ctrl+C${NC} to stop all services.\n"
echo -e "${BLUE}Environment Variables:${NC}"
echo -e "NEXT_PUBLIC_API_URL=${YELLOW}$NEXT_PUBLIC_API_URL${NC}"
echo -e "NEXT_PUBLIC_WS_URL=${YELLOW}$NEXT_PUBLIC_WS_URL${NC}"

# Wait for user to press Ctrl+C
wait
