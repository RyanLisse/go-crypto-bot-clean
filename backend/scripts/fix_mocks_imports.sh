#!/usr/bin/env bash
set -e

# This script updates import paths in test files to reference the relocated mocks.
# Run from repository root: bash scripts/fix_mocks_imports.sh

find backend/internal -type f -name "*_test.go" | while read -r file; do
  # Insert mocks import for usecase tests
  if grep -q 'usecase.Mock' "$file"; then
    sed -i '' \
      -e 's|"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"|"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase";\
    mocks "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/usecase"|g' \
      "$file"
    # Update references: usecase.MockX -> mocks.MockX
    sed -i '' \
      -E 's|usecase\.Mock([A-Za-z0-9_]+)|mocks.Mock\1|g' \
      "$file"
  fi

  # Add similar blocks for other packages as needed...

done

echo "Mocks imports updated. Please run 'go test ./...' to verify."
