#!/bin/bash

# Find all Go files
find . -name "*.go" -type f | while read -r file; do
  # Replace import paths
  sed -i '' 's|"github.com/RyanLisse/go-crypto-bot-clean/|"go-crypto-bot-clean/|g' "$file"
done

echo "Import paths reverted."
