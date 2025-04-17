#!/bin/bash

# This script builds the Go backend server binary.
# It assumes necessary build tools (Go compiler) are available.

# Create local temp directory for Go build if needed (e.g., for sandbox environments)
# export TMPDIR="$(pwd)/.tmp"
# mkdir -p "$TMPDIR"

echo "Building server binary (./server)..."
go build -o server cmd/server/main.go

if [ $? -ne 0 ]; then
  echo "Error: Failed to build server."
  exit 1
fi

echo "Server built successfully: ./server"
# Clean up temp directory if it was created
# rm -rf "$TMPDIR"
