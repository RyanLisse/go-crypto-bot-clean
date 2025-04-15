#!/bin/bash

# Load environment variables from .env file
set -a
source .env
set +a

# Disable Turso
export TURSO_URL=""
export TURSO_AUTH_TOKEN=""

# Create data directory for SQLite
mkdir -p ./data

# Run the server with SQLite
go run cmd/server/main.go
