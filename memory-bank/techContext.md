# Technology Context

## Technology Stack

### Primary Language
- Go (Golang) 1.21+

### Database
- GORM (Go Object Relational Mapper) for data persistence
- Auto-migrations handled via GORM

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
- GORM setup and connection management
- Auto-migration system and schema management
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
  - Position sizing based on account balance and risk parameters
  - Drawdown monitoring with historical balance tracking
  - Exposure limits for trading positions
  - Daily loss limits to prevent excessive losses

## Key Libraries and Dependencies

### Core
- github.com/gin-gonic/gin - Web framework
- gorm.io/gorm - Go Object Relational Mapper
- gorm.io/driver/sqlite - SQLite driver for GORM
- github.com/gorilla/websocket - WebSocket implementation
- github.com/go-telegram-bot-api/telegram-bot-api/v5 - Telegram Bot API client
- github.com/sirupsen/logrus - Structured logger

### Testing
- **testing**: Go standard library testing
- **testify**: Enhanced testing assertions and mocks

## Build & Deployment
- **Makefile**: For build automation
- **Docker**: For containerization (optional)
- **go.mod/go.sum**: For dependency management
