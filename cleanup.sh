#!/bin/bash

# This script will clean up the project structure by removing duplicated files and directories
# It should be run after all necessary files have been copied to their proper locations

echo "Starting cleanup process..."

# Remove duplicated directories
echo "Removing duplicated directories..."
rm -rf api_backup
rm -rf api
rm -rf internal
rm -rf cmd

# Remove duplicated files from root
echo "Removing duplicated files from root..."
rm -f fix_imports.sh
rm -f revert_imports.sh
rm -f run-dev-monorepo.sh
rm -f update_backend_imports.sh
rm -f update_imports.sh

# Remove unnecessary files from root
echo "Removing unnecessary files from root..."
rm -f go.mod
rm -f go.sum
rm -f go.work
rm -f go.work.sum
rm -f bun.lock
rm -f package.json
rm -f package-lock.json

# Remove node_modules from root (if it exists)
if [ -d "node_modules" ]; then
  echo "Removing node_modules from root..."
  rm -rf node_modules
fi

echo "Cleanup complete!"
