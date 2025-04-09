#!/bin/bash

# Find all Go files
find . -name "*.go" -type f | while read -r file; do
  # Replace import paths
  sed -i '' 's|"go-crypto-bot-clean/|"github.com/RyanLisse/go-crypto-bot-clean/|g' "$file"
done

echo "Import paths updated."
