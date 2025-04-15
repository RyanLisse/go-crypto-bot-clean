#!/bin/bash

# Set environment variables
export MEXC_API_KEY=mx0vglsgdd7flAhfqq
export MEXC_SECRET_KEY=0351d73e5a444d5ea5de2d527bd2a07a
export SERVER_PORT=8082

# Run the server
cd /Users/neo/Developer/experiments/go-crypto-bot-clean/backend
go run ./cmd/server/main.go
