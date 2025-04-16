#!/usr/bin/env bash
set -e

# Normalize package declarations for all mocks
find internal/mocks -type f -name '*.go' | while read f; do
  sed -i '' '1s|.*|package mocks|' "$f"
done

echo "All mock files set to package mocks. Now update imports in test files."

# Replace import paths in test files to point to mocks
find internal -type f -name '*_test.go' | while read f; do
  # replace usecase imports
  sed -i '' \
    -E 's|"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"|"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"\n\t"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/usecase"|g' \
    "$f"
  # replace ai gateway mocks import
  sed -i '' \
    -E 's|"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/gateway/ai"|"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/gateway/ai"\n\t"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/adapter/gateway/ai"|g' \
    "$f"
done

echo "Imports updated. Please run 'go test ./...' to verify."
