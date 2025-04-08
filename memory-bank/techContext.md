# Technology Context

## Technology Stack

### Primary Language
- Go (Golang) 1.21+

### Database
- SQLite3 for data persistence
- Migrations handled via embedded SQL files

### API
- RESTful API using Gin framework
- WebSocket for real-time updates

### External Integration
- MEXC Exchange API
- WebSocket API for market data

### Development Tools
- Go Modules for dependency management
- Make for build automation
- Docker for containerization
- GitHub Actions for CI/CD

## Implementation Documentation

The implementation is thoroughly documented with step-by-step guidelines:

### Core Documentation
- Architecture overview and implementation strategy
- Domain models and repository interfaces
- Service interfaces and business logic implementation

### Database Layer
- Database layer overview and repository pattern
- SQLite setup and connection management
- Migration system and schema management
- Repository implementation examples

### API Layer
- API layer overview and structure
- Middleware components for authentication, logging, etc.
- REST API endpoint handlers
- WebSocket implementation for real-time data

### Advanced Trading Features
- Advanced trading strategies using multi-indicator analysis
- Position management with lifecycle controls
- Risk management and capital protection systems

## Key Libraries and Dependencies

### Core
- github.com/gin-gonic/gin - Web framework
- github.com/jmoiron/sqlx - Enhanced database access
- github.com/mattn/go-sqlite3 - SQLite driver
- github.com/gorilla/websocket - WebSocket implementation

### Testing
- **testing**: Go standard library testing
- **testify**: Enhanced testing assertions and mocks

## Build & Deployment
- **Makefile**: For build automation
- **Docker**: For containerization (optional)
- **go.mod/go.sum**: For dependency management
