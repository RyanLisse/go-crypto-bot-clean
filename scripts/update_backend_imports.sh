#!/bin/bash

# Find all Go files in the backend/api directory
find ./backend/api -name "*.go" -type f | while read -r file; do
  # Replace import paths
  sed -i '' 's|"go-crypto-bot-clean/api/|"go-crypto-bot-clean/backend/api/|g' "$file"
done

echo "Backend import paths updated."
