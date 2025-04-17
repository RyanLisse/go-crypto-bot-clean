#!/bin/bash

# Script to convert controller files to handler files
# This ensures all HTTP handlers are in a consistent location and naming convention

echo "Converting controller files to handler files..."

# Create necessary directories
mkdir -p internal/adapter/delivery/http/handler/auth
mkdir -p internal/adapter/delivery/http/handler/user
mkdir -p internal/adapter/delivery/http/handler/wallet
mkdir -p internal/adapter/delivery/http/handler/example

# Function to convert a controller file to a handler file
convert_file() {
  SOURCE=$1
  TARGET=$2
  CONTROLLER_TYPE=$3
  
  echo "Converting $SOURCE to $TARGET"
  
  # Copy the file
  cp "$SOURCE" "$TARGET"
  
  # Replace package name
  sed -i '' 's/package controller/package handler/g' "$TARGET"
  
  # Replace "Controller" with "Handler" in type names and function names
  sed -i '' "s/${CONTROLLER_TYPE}Controller/${CONTROLLER_TYPE}Handler/g" "$TARGET"
  
  # Replace "New${CONTROLLER_TYPE}Controller" with "New${CONTROLLER_TYPE}Handler"
  sed -i '' "s/New${CONTROLLER_TYPE}Controller/New${CONTROLLER_TYPE}Handler/g" "$TARGET"
}

# Process auth controller files
if [ -f "internal/adapter/http/controller/auth_controller.go" ]; then
  convert_file "internal/adapter/http/controller/auth_controller.go" "internal/adapter/delivery/http/handler/auth/auth_handler.go" "Auth"
fi

if [ -f "internal/adapter/http/controller/auth_controller_test.go" ]; then
  convert_file "internal/adapter/http/controller/auth_controller_test.go" "internal/adapter/delivery/http/handler/auth/auth_handler_test.go" "Auth"
fi

# Process user controller files
if [ -f "internal/adapter/http/controller/user_controller.go" ]; then
  convert_file "internal/adapter/http/controller/user_controller.go" "internal/adapter/delivery/http/handler/user/user_handler.go" "User"
fi

if [ -f "internal/adapter/http/controller/user_controller_test.go" ]; then
  convert_file "internal/adapter/http/controller/user_controller_test.go" "internal/adapter/delivery/http/handler/user/user_handler_test.go" "User"
fi

# Process wallet controller files
if [ -f "internal/adapter/http/controller/wallet_controller.go" ]; then
  convert_file "internal/adapter/http/controller/wallet_controller.go" "internal/adapter/delivery/http/handler/wallet/wallet_handler.go" "Wallet"
fi

# Process example controller files
if [ -f "internal/adapter/http/controller/example_error_controller.go" ]; then
  convert_file "internal/adapter/http/controller/example_error_controller.go" "internal/adapter/delivery/http/handler/example/example_error_handler.go" "ExampleError"
fi

# Create a README.md file in the old location
mkdir -p internal/adapter/http/controller
cat > internal/adapter/http/controller/README.md << EOL
# DEPRECATED: Controller Files Have Moved

The controller files that were previously in this directory have been relocated to:
\`internal/adapter/delivery/http/handler/*\`

Please update your imports to use the new location.
This ensures consistency with other handlers in the project.

All controllers have been renamed to handlers for consistency.
EOL

echo "Controller conversion completed!" 