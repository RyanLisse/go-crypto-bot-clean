#!/bin/bash

# This script starts the server after killing any process using port 8080
# Usage: ./start_server.sh [port]

PORT=${1:-8080}  # Default to port 8080 if not specified
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Kill any process using the port
"$SCRIPT_DIR/kill_port.sh" $PORT

# Build the server
echo "Building server..."
cd "$PROJECT_ROOT" && go build -o server ./cmd/server

# Start the server
echo "Starting server on port $PORT..."
cd "$PROJECT_ROOT" && ./server
