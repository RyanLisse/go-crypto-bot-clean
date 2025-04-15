#!/bin/bash

# Find all Go files
find . -name "*.go" -type f | while read -r file; do
  # Replace github.com/neo/crypto-bot with github.com/RyanLisse/go-crypto-bot-clean/backend
  sed -i '' 's|github.com/neo/crypto-bot|github.com/RyanLisse/go-crypto-bot-clean/backend|g' "$file"
done

echo "Import paths fixed!"
