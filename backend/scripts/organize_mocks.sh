#!/bin/bash

# Script to organize all mock files into internal/mocks directory structure
# This ensures all mocks are in a consistent location

echo "Organizing mock files..."

# Create necessary directories if they don't exist
mkdir -p internal/mocks/domain/port
mkdir -p internal/mocks/domain/service
mkdir -p internal/mocks/usecase
mkdir -p internal/mocks/infrastructure
mkdir -p internal/mocks/adapter

# Check if internal/adapter/persistence/mock directory exists and move its contents
if [ -d "internal/adapter/persistence/mock" ]; then
  echo "Moving mocks from internal/adapter/persistence/mock to internal/mocks/adapter/persistence"
  mkdir -p internal/mocks/adapter/persistence
  cp -r internal/adapter/persistence/mock/* internal/mocks/adapter/persistence/
  echo "Creating a README.md file in the old location"
  mkdir -p internal/adapter/persistence/mock
  cat > internal/adapter/persistence/mock/README.md << EOL
# DEPRECATED: Mock Files Have Moved

The mock files that were previously in this directory have been relocated to:
\`internal/mocks/adapter/persistence\`

Please update your imports to use the new location.
This ensures consistency with other mocks in the project.
EOL
fi

# Create a simple README for the mocks directory
cat > internal/mocks/README.md << EOL
# Mocks Directory

This directory contains all mock implementations used for testing.

## Directory Structure

- \`domain/port\`: Mocks for domain ports (repositories, clients, etc.)
- \`domain/service\`: Mocks for domain services
- \`usecase\`: Mocks for use cases
- \`infrastructure\`: Mocks for infrastructure components (DB, cache, etc.)
- \`adapter\`: Mocks for adapters (HTTP, persistence, etc.)

## Best Practices

1. Use a tool like mockery to generate mocks
2. Keep mocks in sync with their interfaces
3. Don't modify generated mocks directly
4. Include a generation command in a comment
EOL

echo "Mocks organization completed!" 