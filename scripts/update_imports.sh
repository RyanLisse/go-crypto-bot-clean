#!/bin/bash

# Update import paths in all Go files in the backend/internal/api directory
find backend/internal/api -name "*.go" -type f -exec sed -i '' 's|go-crypto-bot-clean/api|go-crypto-bot-clean/backend/internal/api|g' {} \;

echo "Import paths updated successfully!"
