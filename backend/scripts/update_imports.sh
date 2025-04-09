#!/bin/bash

# Replace all occurrences of the GitHub import path with the local module path
find . -type f -name "*.go" -exec sed -i '' 's|github.com/RyanLisse/go-crypto-bot-clean/backend|go-crypto-bot-clean/backend|g' {} +

# Run go mod tidy to update dependencies
go mod tidy
