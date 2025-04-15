#!/bin/bash

# Create a migration script for moving Go files from /internal to /backend/internal
# This script will:
# 1. Copy all files, preserving directory structure
# 2. Update import paths to use the backend module
# 3. Keep backup of any existing files in the destination

# Set directories
SOURCE_DIR="internal"
DEST_DIR="backend/internal"
BACKUP_DIR="backend/internal_backup_$(date +%Y%m%d_%H%M%S)"

# Create backup of any existing files in destination
echo "Creating backup of existing files in $DEST_DIR to $BACKUP_DIR"
mkdir -p "$BACKUP_DIR"
cp -rf "$DEST_DIR"/* "$BACKUP_DIR"

# Copy files from source to destination, preserving directory structure
echo "Copying files from $SOURCE_DIR to $DEST_DIR"
find "$SOURCE_DIR" -type f -name "*.go" | while read file; do
    rel_path="${file#$SOURCE_DIR/}"
    dest_file="$DEST_DIR/$rel_path"
    dest_dir=$(dirname "$dest_file")
    
    # Create destination directory if it doesn't exist
    mkdir -p "$dest_dir"
    
    # Copy the file
    cp "$file" "$dest_file"
    echo "Copied: $file -> $dest_file"
    
    # Update import paths in the copied file
    sed -i '' 's#github.com/RyanLisse/go-crypto-bot-clean/internal#github.com/RyanLisse/go-crypto-bot-clean/backend/internal#g' "$dest_file"
    echo "Updated imports in: $dest_file"
done

echo "Migration completed. Files from $SOURCE_DIR have been copied to $DEST_DIR"
echo "Original files in $DEST_DIR were backed up to $BACKUP_DIR"
echo ""
echo "Next steps:"
echo "1. Review the migrated files and resolve any conflicts"
echo "2. Run tests to ensure everything works correctly"
echo "3. When satisfied, remove the original files from $SOURCE_DIR" 