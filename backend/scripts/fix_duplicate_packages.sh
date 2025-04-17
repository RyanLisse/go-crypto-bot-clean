#!/usr/bin/env bash
set -e

# This script fixes duplicate package declarations in mock files
# Run from the backend directory: bash scripts/fix_duplicate_packages.sh

find internal/mocks -type f -name "*.go" | while read -r file; do
  # Check if the file has duplicate package declarations
  if grep -q "package mocks\s*\n\s*package mocks" "$file"; then
    echo "Fixing duplicate package declaration in $file"
    # Replace the duplicate package declaration with a single one
    sed -i '' -e 's/package mocks\s*\n\s*package mocks/package mocks/g' "$file"
  fi
done

echo "Fixed duplicate package declarations in mock files."
