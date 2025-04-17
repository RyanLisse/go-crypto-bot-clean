#!/bin/bash

# This script runs the server in development mode using SQLite.
# It loads variables from .env but explicitly disables Turso
# to force the application's fallback to SQLite.
# It uses 'go run' for convenience, compiling and running in one step.

echo "Starting server in development mode (SQLite)..."

# Load environment variables from .env file if it exists
if [ -f ".env" ]; then
  echo "Loading environment variables from .env file..."
  set -a
  source .env
  set +a
else
  echo "Warning: .env file not found. Proceeding without it."
fi

# Explicitly disable Turso to force SQLite usage
echo "Disabling Turso (setting TURSO_URL and TURSO_AUTH_TOKEN to empty)..."
export TURSO_URL=""
export TURSO_AUTH_TOKEN=""

# Ensure data directory exists (SQLite database file might be stored here)
mkdir -p ./data

# Run the server using 'go run' (compiles and runs)
echo "Running server with 'go run'..."
go run cmd/server/main.go

if [ $? -ne 0 ]; then
  echo "Error: Failed to run server."
  exit 1
fi
