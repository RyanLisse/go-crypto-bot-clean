#!/usr/bin/env bash
set -e

# This script fixes all mock files with duplicate package declarations and other issues
# Run from the backend directory: bash scripts/fix_all_mocks.sh

echo "Fixing mock files..."

# Process each mock file
find internal/mocks -type f -name "*.go" | while read -r file; do
  echo "Processing $file"
  
  # Create a temporary file
  tmp_file=$(mktemp)
  
  # Write the package declaration
  echo "package mocks" > "$tmp_file"
  echo "" >> "$tmp_file"
  
  # Extract the imports and content after the second package declaration
  awk 'BEGIN{found=0} /^import/{found=1; print; next} found{print}' "$file" >> "$tmp_file"
  
  # Replace the original file with the fixed one
  mv "$tmp_file" "$file"
done

echo "Fixed all mock files."

# Now fix the references in the factory package
echo "Fixing references in factory package..."

# Update usecase_factory.go to use mocks package
sed -i 's/usecase\.Mock\([A-Za-z]*\)UseCase/mocks.Mock\1UseCase/g' internal/factory/usecase_factory.go

# Add import for mocks package if needed
if ! grep -q '"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/usecase"' internal/factory/usecase_factory.go; then
  sed -i '/^import (/a\
    "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/usecase"' internal/factory/usecase_factory.go
fi

echo "Fixed references in factory package."

echo "All fixes completed. Please run 'go build ./...' to verify."
