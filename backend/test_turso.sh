#!/bin/bash

# Load environment variables from .env file
set -a
source .env
set +a

# Configure Turso settings
export TURSO_SYNC_INTERVAL_SECONDS="60"  # 1 minute for testing
export TURSO_SYNC_ENABLED="true"

# Create data directory for Turso
mkdir -p ./data/turso

# Run the test with Turso support
go run -tags turso test_turso.go
