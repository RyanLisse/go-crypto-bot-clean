# Clean Backend for Go Crypto Bot

This is the clean architecture implementation of the Go Crypto Bot backend, following hexagonal architecture principles.

## Project Structure

```
clean_backend/
├── cmd/                    # Application entry points
│   ├── server/             # HTTP API server
│   └── worker/             # Background worker
├── internal/               # Private application code
│   ├── adapter/            # Adapters (infrastructure & delivery)
│   │   ├── delivery/       # Inbound adapters (HTTP, gRPC, etc.)
│   │   │   └── http/       # HTTP delivery layer
│   │   │       ├── handler/    # HTTP handlers
│   │   │       ├── middleware/ # HTTP middleware
│   │   │       ├── response/   # HTTP response definitions
│   │   │       └── server/     # HTTP server setup
│   │   └── infrastructure/ # Outbound adapters
│   │       ├── gateway/    # External service clients
│   │       └── persistence/ # Data persistence
│   │           └── gorm/   # GORM-based repositories
│   │               ├── entity/  # Database entities
│   │               ├── migrations/ # Database migrations
│   │               └── repo/   # Repository implementations
│   ├── apperror/           # Application error definitions
│   ├── config/             # Configuration management
│   ├── di/                 # Dependency injection
│   ├── domain/             # Domain layer (core business logic)
│   │   ├── model/          # Domain models
│   │   ├── port/           # Ports (interfaces)
│   │   └── service/        # Domain services
│   ├── factory/            # Factory methods for creating components
│   ├── usecase/            # Application use cases
│   │   └── dto/            # Data Transfer Objects
│   └── util/               # Utility functions
│       ├── crypto/         # Cryptography utilities
│       └── validator/      # Validation utilities
├── pkg/                    # Public libraries that can be used by other services
├── scripts/                # Build and deployment scripts
└── test/                   # Integration and end-to-end tests
```

## Hexagonal Architecture

This project follows the Hexagonal Architecture (also known as Ports and Adapters) pattern:

1. **Domain Layer**: Contains the core business logic, domain models, and interfaces (ports).
   - Located in `internal/domain/`
   - Independent of external concerns
   - Defines interfaces that adapters must implement

2. **Application Layer**: Orchestrates the flow of data and coordinates domain operations.
   - Located in `internal/usecase/`
   - Implements use cases that coordinate between domain and adapters
   - Uses DTOs to transfer data between layers

3. **Adapter Layer**: Connects the application to external concerns.
   - **Inbound Adapters**: Handle incoming requests (HTTP, gRPC, etc.)
     - Located in `internal/adapter/delivery/`
   - **Outbound Adapters**: Handle outgoing requests (database, external APIs, etc.)
     - Located in `internal/adapter/infrastructure/`

## Getting Started

### Prerequisites

- Go 1.21 or higher
- SQLite (for local development)

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Create a `.env` file based on `.env.example`:
   ```bash
   cp .env.example .env
   ```
4. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

### Configuration

Configuration is managed through environment variables and/or a `config.yaml` file. See `internal/config/config.go` for available options.

## Development

### Adding a New Feature

1. Define domain models in `internal/domain/model/`
2. Define ports (interfaces) in `internal/domain/port/`
3. Implement domain services in `internal/domain/service/`
4. Implement use cases in `internal/usecase/`
5. Implement adapters in `internal/adapter/`
6. Wire everything together in `internal/di/`

### Running the Server

You can run the server using the Makefile:

```bash
# Run the server normally
make run-server

# Run the server with automatic port killing
make run-server-safe

# Run the server on a different port
PORT=3000 make run-server-safe
```

Alternatively, you can run the server directly:

```bash
go run cmd/server/main.go
```

Or use the provided scripts to automatically kill any process using port 8080 before starting the server:

**On macOS/Linux:**

```bash
./kill_port.sh
```

**On Windows:**

```cmd
kill_port.bat
```

You can also specify a different port using the PORT environment variable:

```bash
# On macOS/Linux
PORT=3000 ./kill_port.sh

# On Windows
set PORT=3000
kill_port.bat
```

### Running the Worker

```bash
go run cmd/worker/main.go
```

### Building

```bash
go build -o bin/server cmd/server/main.go
go build -o bin/worker cmd/worker/main.go
```

### Testing

```bash
go test ./...
```

## Deployment

The application can be deployed using Docker:

```bash
docker build -t crypto-bot-backend .
docker run -p 8080:8080 crypto-bot-backend
```
