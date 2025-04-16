# Crypto Trading Bot (MEXC Focused)

A Go-based automated cryptocurrency trading bot designed primarily for the MEXC exchange. This bot provides automated execution based on predefined strategies, real-time market data analysis, and robust risk management.

## Current Status

**June 2024**: Major architectural refactoring has been completed! The codebase now fully adheres to Clean Architecture principles. Current development focus is on implementing the Risk Management System.

## Features

- MEXC Exchange Integration via REST API and WebSocket
- Account & Portfolio Management
- Market Data Handling
- Trade Execution & Order Management
- Position Management
- Strategy Engine
- Risk Management (In Progress)
- New Coin Detection & AutoBuy
- AI Assistant (Gemini Integration)
- Notifications
- Analytics & Reporting
- System Status & Monitoring

## Project Structure

```
├── cmd/
│   └── server/               # Application entry point
├── internal/
│   ├── config/               # Configuration loading
│   ├── domain/               # Core business logic and entities
│   │   ├── model/            # Domain entities
│   │   └── port/             # Interfaces for external dependencies
│   ├── usecase/              # Application logic / features
│   ├── adapter/              # Implementations of ports
│   │   ├── delivery/         # HTTP API, CLI
│   │   ├── persistence/      # Database implementations
│   │   ├── gateway/          # External service integrations
│   │   └── cache/            # Cache implementations
│   ├── platform/             # Low-level platform components
│   └── apperror/             # Application-specific error types
├── pkg/                      # Shared libraries
└── migrations/               # Database migration files
```

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL (or compatible database)
- MEXC API credentials

### Installation

1. Clone the repository
2. Install dependencies:
   ```
   go mod tidy
   ```
3. Configure environment variables (see `.env.example`)
4. Run the application:
   ```
   go run cmd/server/main.go
   ```

## Configuration

Configuration is loaded from environment variables and/or a config file. See `.env.example` for required variables.

## Development

This project follows Clean Architecture principles with a clear separation of concerns between domain, use cases, and adapters. See `docs/` directory for architecture documentation and guides.

## License

[MIT License](LICENSE)