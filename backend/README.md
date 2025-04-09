# Go Crypto Bot

## Project Overview

Go Crypto Bot is a sophisticated cryptocurrency trading bot designed to provide robust, efficient, and flexible trading capabilities across multiple exchanges.

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
- Language: Go (1.21+)
- WebSocket Library: gorilla/websocket
- HTTP Router: chi
- API Documentation: Huma
- Testing: testify
- Rate Limiting: Custom token bucket implementation
- Database: SQLite with sqlx
- Logging: zap
- CLI: cobra

### Frontend
- React 18 with TypeScript
- Material UI for components and styling
- Redux Toolkit for state management
- RTK Query for API calls
- Chart.js for data visualization
- Socket.io for real-time updates

### DevOps
- Docker for containerization
- Docker Compose for multi-container orchestration
- GitHub Actions for CI/CD

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

## Getting Started

### Prerequisites
- Go 1.21+
- Docker and Docker Compose (optional, for containerized deployment)
- Node.js 16+ and npm/yarn (for frontend development)

### Installation

#### Backend
```bash
git clone https://github.com/ryanlisse/go-crypto-bot.git
cd go-crypto-bot
go mod download
```

#### Frontend
```bash
cd frontend
npm install
# or
yarn install
```

### Running with Docker
```bash
# Create a .env file from the example
cp .env.example .env
# Edit the .env file with your API keys and settings

# Start the application
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the application
docker-compose down
```

### Deploying to Railway

The application is configured for easy deployment to Railway:

1. Create a new project in Railway

2. Connect your GitHub repository

3. Railway will automatically detect the Dockerfile and build your application

4. Set up the required environment variables in Railway:
   - `MEXC_API_KEY`
   - `MEXC_SECRET_KEY`
   - `TURSO_URL` (if using Turso database)
   - `TURSO_AUTH_TOKEN` (if using Turso database)
   - `OPENAI_API_KEY` (if using AI features)
   - `GOOGLE_API_KEY` (if using Google services)
   - `ENVIRONMENT=production`

5. Deploy your application

6. Railway will automatically assign a domain to your application

### Running Locally

#### Backend
```bash
# Run the API server
go run main.go serve --port=8080

# Run tests
go test ./...
```

#### Frontend
```bash
cd frontend
npm start
# or
yarn start
```

### Running Backtests
```bash
go run cmd/backtest/main.go backtest --strategy=simple_ma --symbols=BTCUSDT --start=2023-01-01 --end=2023-12-31 --capital=10000 --interval=1h
```

### Available CLI Commands
```bash
# Get help
go run main.go --help

# Run the API server
go run main.go serve --port=8080
```

## Contributing
Contributions are welcome! Please read our contributing guidelines before submitting pull requests.

## License
[Specify your license here]