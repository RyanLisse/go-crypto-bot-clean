# Project Structure

This document outlines the organization of the Go Crypto Bot project.

## Directory Structure

```
go-crypto-bot-clean/
├── backend/               # Backend Go code
│   ├── api/               # API implementation
│   │   ├── database/      # Database access
│   │   ├── huma/          # Huma OpenAPI integration
│   │   ├── middleware/    # API middleware
│   │   ├── models/        # API data models
│   │   ├── repository/    # Data repositories
│   │   ├── service/       # Business logic services
│   │   └── main.go        # API entry point
│   ├── cmd/               # Command-line applications
│   ├── internal/          # Internal packages
│   ├── mocks/             # Mock implementations for testing
│   ├── pkg/               # Reusable packages
│   └── tests/             # Integration and end-to-end tests
│
├── frontend/              # Frontend React application
│   ├── public/            # Static assets
│   ├── src/               # Source code
│   │   ├── components/    # React components
│   │   │   ├── layout/    # Layout components
│   │   │   └── ui/        # UI components
│   │   ├── hooks/         # Custom React hooks
│   │   ├── lib/           # Utility functions
│   │   ├── pages/         # Page components
│   │   └── services/      # API service clients
│   └── ...                # Configuration files
│
├── memory-bank/           # Project documentation and context
├── scripts/               # Utility scripts
└── tasks/                 # Task definitions and documentation
```

## Key Components

### Backend

- **API**: RESTful API implementation using Chi router
- **Services**: Business logic implementation
- **Repositories**: Data access layer
- **Models**: Data structures and domain models

### Frontend

- **Components**: Reusable React components
- **Pages**: Application pages and routes
- **Services**: API client implementations
- **Hooks**: Custom React hooks for state management and data fetching

## Development Workflow

1. Backend development is done in the `backend` directory
2. Frontend development is done in the `frontend` directory
3. Scripts for development and deployment are in the `scripts` directory
4. Project documentation is maintained in the `memory-bank` directory
