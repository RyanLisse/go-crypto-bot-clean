#!/bin/bash

# Load environment variables from .env file
set -a
source .env
set +a

# Check if Turso is configured
if [ -z "$TURSO_URL" ] || [ -z "$TURSO_AUTH_TOKEN" ]; then
  echo "Error: Turso is not properly configured. Please set TURSO_URL and TURSO_AUTH_TOKEN in your .env file."
  exit 1
fi

# Configure Turso settings
export TURSO_SYNC_INTERVAL_SECONDS="300"  # 5 minutes
export TURSO_SYNC_ENABLED="true"

# Create data directory for Turso
mkdir -p ./data/turso

# Run the server
./server
