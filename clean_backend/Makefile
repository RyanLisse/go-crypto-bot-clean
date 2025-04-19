.PHONY: build run-server run-server-safe run-worker test clean

# Default port
PORT ?= 8080

# Build both server and worker
build:
	@echo "Building server and worker..."
	@mkdir -p bin
	@go build -o bin/server cmd/server/main.go
	@go build -o bin/worker cmd/worker/main.go
	@echo "Build complete. Binaries are in the bin/ directory."

# Run the server
run-server:
	@echo "Running server..."
	@go run cmd/server/main.go

# Run the server with automatic port killing
run-server-safe:
	@echo "Running server on port $(PORT) with automatic port killing..."
	@PORT=$(PORT) ./kill_port.sh

# Run the worker
run-worker:
	@echo "Running worker..."
	@go run cmd/worker/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean complete."

# Create database directory
create-db-dir:
	@echo "Creating database directory..."
	@mkdir -p data
	@echo "Database directory created."

# Initialize the project
init: create-db-dir build

# Docker commands
docker-build:
	@echo "Building Docker images..."
	@docker-compose build

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	@docker-compose logs -f

# Help
help:
	@echo "Available targets:"
	@echo "  build       - Build server and worker binaries"
	@echo "  run-server  - Run the server"
	@echo "  run-server-safe - Run the server with automatic port killing"
	@echo "  run-worker  - Run the worker"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  create-db-dir - Create database directory"
	@echo "  init        - Initialize the project"
	@echo "  docker-build - Build Docker images"
	@echo "  docker-up   - Start Docker containers"
	@echo "  docker-down - Stop Docker containers"
	@echo "  docker-logs - Show Docker logs"
	@echo "  help        - Show this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  PORT        - Server port (default: 8080)"
