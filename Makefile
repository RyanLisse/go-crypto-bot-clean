# Makefile for go-crypto-bot-clean

# Variables
DOCKER_COMPOSE = docker-compose
DOCKER_COMPOSE_DEV = docker-compose -f docker-compose.dev.yml

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build            - Build all containers"
	@echo "  make up               - Start all containers in production mode"
	@echo "  make down             - Stop all containers"
	@echo "  make dev              - Start all containers in development mode"
	@echo "  make backend          - Start only the backend container in production mode"
	@echo "  make frontend         - Start only the frontend container in production mode"
	@echo "  make backend-dev      - Start only the backend container in development mode"
	@echo "  make frontend-dev     - Start only the frontend container in development mode"
	@echo "  make logs             - Show logs from all containers"
	@echo "  make clean            - Remove all containers, volumes, and images"
	@echo "  make run-backend      - Run the backend server locally with proper env variables"

# Build all containers
.PHONY: build
build:
	$(DOCKER_COMPOSE) build

# Start all containers in production mode
.PHONY: up
up:
	$(DOCKER_COMPOSE) up -d

# Stop all containers
.PHONY: down
down:
	$(DOCKER_COMPOSE) down

# Start all containers in development mode
.PHONY: dev
dev:
	$(DOCKER_COMPOSE_DEV) up -d

# Start only the backend container in production mode
.PHONY: backend
backend:
	$(DOCKER_COMPOSE) up -d backend

# Start only the frontend container in production mode
.PHONY: frontend
frontend:
	$(DOCKER_COMPOSE) up -d frontend

# Start only the backend container in development mode
.PHONY: backend-dev
backend-dev:
	$(DOCKER_COMPOSE_DEV) up -d backend-dev

# Start only the frontend container in development mode
.PHONY: frontend-dev
frontend-dev:
	$(DOCKER_COMPOSE_DEV) up -d frontend-dev

# Show logs from all containers
.PHONY: logs
logs:
	$(DOCKER_COMPOSE) logs -f

# Remove all containers, volumes, and images
.PHONY: clean
clean:
	$(DOCKER_COMPOSE) down -v --rmi all
	$(DOCKER_COMPOSE_DEV) down -v --rmi all

# Run the backend server locally with proper env variables
.PHONY: run-backend
run-backend:
	@echo "Starting backend server with proper environment variables..."
	cd backend && \
	export MEXC_API_KEY="mx0vglsgdd7flAhfqq" && \
	export MEXC_SECRET_KEY="0351d73e5a444d5ea5de2d527bd2a07a" && \
	export MEXC_BASE_URL="https://api.mexc.com" && \
	export HTTP_PROXY="" && \
	export HTTPS_PROXY="" && \
	go run main.go serve
