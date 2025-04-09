# Go Crypto Bot - Monorepo

## Project Overview

Go Crypto Bot is a sophisticated cryptocurrency trading bot designed to provide robust, efficient, and flexible trading capabilities across multiple exchanges. This repository is structured as a monorepo containing multiple deployable components that work together as a complete system.

## Current Status

The project is currently under active development with a focus on building a reliable and scalable architecture. The backend API and WebSocket services are operational, and the frontend React application is in development. The system currently supports the MEXC exchange with plans to add more exchanges in the future.

## Features

### WebSocket Client
- Real-time market data streaming
- Thread-safe connection management
- Automatic reconnection
- Ticker subscription mechanism
- Comprehensive error handling

### REST API Client
- Comprehensive MEXC Exchange API integration
- Rate limiting with token bucket algorithm
- Efficient caching for market data
- Thread-safe operations
- Robust error handling

### Backtesting Framework
- Historical data loading from various sources
- Position tracking and P&L calculation
- Performance metrics (Sharpe ratio, drawdown, etc.)
- Slippage models for realistic simulation
- Strategy interface for testing different strategies
- CLI command for running backtests

### Current Capabilities
- MEXC Exchange WebSocket and REST API integration
- Ticker data processing and caching
- Exponential backoff reconnection strategy
- Strategy implementation and backtesting
- Risk management and position sizing
- Real-time notifications via Telegram and Slack
- Docker containerization for easy deployment
- React frontend for monitoring and control

## Technology Stack

### Backend
- Language: Go (1.24.2)
- Web Frameworks: chi, gin-gonic/gin
- WebSocket Library: gorilla/websocket
- API Documentation: Huma v2
- Testing: testify
- Rate Limiting: Custom token bucket implementation
- Database: SQLite with sqlx and GORM
- Logging: zap
- CLI: cobra with viper for configuration
- JWT Authentication: golang-jwt/jwt/v5
- AI Integration: google/generative-ai-go

### Frontend
- React 18 with TypeScript
- Material UI for components and styling
- Redux Toolkit for state management
- Axios for API calls
- Chart.js with react-chartjs-2 for data visualization
- Socket.io for real-time updates
- Formik and Yup for form validation
- React Router v6 for routing
- React Toastify for notifications
- Date-fns for date manipulation
- Build tools: Create React App with Vite configuration available

### DevOps
- Docker for containerization
- Docker Compose for multi-container orchestration
- GitHub Actions for CI/CD
- Bun as an optional JavaScript runtime (with npm fallback)

## Development Approach
- Test-Driven Development (TDD)
- Modular and extensible design
- Performance-focused implementation

## Upcoming Milestones
- Parameter optimization for backtesting
- Monte Carlo simulation for strategy robustness testing
- Walk-forward analysis
- Visualization tools for equity curves and drawdowns
- Frontend integration with React
- Advanced trading strategy implementation
- Machine learning integration

## Monorepo Structure

This repository is organized as a monorepo containing multiple deployable components:

- `backend/`: Go backend services for core trading functionality
  - `cmd/`: Command-line tools and entry points
  - `internal/`: Shared internal packages
  - `pkg/`: Shared public packages
  - `configs/`: Configuration templates for different environments
  - `tests/`: Test suites for backend components
  - `docs/`: API documentation
  - `data/`: Data storage directory
- `frontend/`: React frontend application for monitoring and control
  - `src/`: Source code for the React application
  - `public/`: Static assets
  - `docs/`: Frontend documentation
- `memory-bank/`: Storage for AI-assisted trading strategies and historical data
- `.env.example`: Example environment variables configuration

## Getting Started

### Prerequisites
- Go 1.24+
- Docker and Docker Compose (for containerized deployment)
- Node.js 18+ and npm (for frontend development)
- Bun (optional, for faster frontend development)

### Installation

#### Clone Repository
```bash
git clone https://github.com/RyanLisse/go-crypto-bot-clean.git
cd go-crypto-bot-clean
```

#### Backend
```bash
go mod download
```

#### Frontend
```bash
cd frontend
npm install
# or
yarn install
```

### Development Workflow

The repository includes a development script that starts both the backend and frontend services:

```bash
# Make the script executable if needed
chmod +x run-dev-monorepo.sh

# Run the development environment
./run-dev-monorepo.sh
```

This script will:
- Start the Go backend API server on port 8080
- Start the React frontend development server on port 3000 (using npm or Bun)
- Set up the necessary environment variables for API and WebSocket connections
- Provide a clean shutdown mechanism with Ctrl+C

## Deployment

This monorepo is designed for deployment to multiple platforms. Each component can be deployed independently or as part of the complete system.

### Component Deployment

#### Backend API (Go)
- Can be deployed to any cloud provider that supports Go or containers
- Supports containerized deployment with Docker
- Can be deployed to Kubernetes, ECS, or standalone servers

#### Frontend Application (React)
- Can be deployed to Netlify, Vercel, or any static hosting provider
- Supports containerized deployment for environments like Kubernetes
- Configuration via environment variables for different environments

#### Dashboard (React)
- Can be deployed to Netlify, Vercel, or any static hosting provider
- Separate deployment process from the main frontend
- Shares authentication with the main application

### Running with Docker
```bash
# Create a .env file from the example
cp .env.example .env
# Edit the .env file with your API keys and settings

# Start all services
docker-compose up -d

# Start only specific services
docker-compose up -d backend frontend

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### Running Locally

#### Backend
```bash
# Run the API server
cd backend
go run cmd/api/main.go --port=8080

# Run tests
go test ./...
```

#### Frontend
```bash
cd frontend
# Using npm
npm start

# Or using Bun (if installed)
bun start
```

### Available CLI Commands
```bash
# Get help for the API server
cd backend
go run cmd/api/main.go --help
```

## CI/CD Integration

This monorepo uses GitHub Actions for continuous integration and deployment:

- Separate workflows for each deployable component
- Shared linting and testing steps
- Deployment to different environments based on branch
- Comprehensive test coverage across all components

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting pull requests.

When contributing to this monorepo:
1. Clearly indicate which component(s) your changes affect
2. Run the appropriate tests for the affected components
3. Ensure your changes don't break other components that may depend on shared code

## License
MIT License