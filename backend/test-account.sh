#!/bin/bash

# Change to the backend directory
cd "$(dirname "$0")"

# Build the test program
echo "Building test program..."
go build -o bin/test_account cmd/test_account/main.go

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful"

# Run the test program
echo -e "\nRunning account test..."
./bin/test_account

# Check if test was successful
if [ $? -ne 0 ]; then
    echo "❌ Test failed"
    exit 1
fi

echo -e "\n✅ Test completed successfully"
