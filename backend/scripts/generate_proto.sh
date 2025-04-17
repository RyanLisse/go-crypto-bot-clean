#!/bin/bash

# Exit on error
set -e

# Directory containing .proto files
PROTO_DIR="./pkg/platform/mexc/websocket/proto"

# Output directory for generated Go files
OUTPUT_DIR="./pkg/platform/mexc/websocket/proto"

# Create output directory if it doesn't exist
mkdir -p $OUTPUT_DIR

# Generate Go code from Protocol Buffers
echo "Generating Go code from Protocol Buffers..."
protoc --go_out=. --go_opt=paths=source_relative \
    $PROTO_DIR/mexc.proto

echo "Protocol Buffer code generation completed successfully!"
