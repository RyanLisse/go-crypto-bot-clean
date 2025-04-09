Remembering...

I'll create a comprehensive Product Requirements Document (PRD) for the Go Crypto Trading Bot based on the architecture and context provided in the documents.

# Product Requirements Document: Go Crypto Trading Bot

## 1. Executive Summary

### Product Vision
A robust, modular cryptocurrency trading bot built in Go, designed to autonomously detect and trade new coins on the MEXC exchange while providing flexible deployment options (backend server and CLI) and advanced trading strategies.

### Key Objectives
- Automate cryptocurrency trading with intelligent decision-making
- Provide a scalable and maintainable software architecture
- Support multiple deployment modes (server and CLI)
- Implement advanced risk management strategies

## 2. Target Users
- Cryptocurrency traders
- Algorithmic trading enthusiasts
- Technical users comfortable with command-line interfaces
- Small to medium-scale individual investors

## 3. Functional Requirements

### 3.1 Core Trading Functionality
- Automatic detection of new coins on MEXC exchange
- Dynamic trading strategy with configurable parameters:
  - Stop-loss at 15% 
  - Take-profit levels: [5%, 10%, 15%, 20%]
- Real-time wallet balance monitoring
- Soft delete mechanism for trading positions
- Comprehensive logging of trading decisions

### 3.2 Deployment Modes
#### Backend Server Mode
- RESTful API endpoints for:
  - Configuration management
  - Trade monitoring
  - Portfolio status
  - Historical trade data retrieval

#### CLI Mode
- Standalone executable with commands:
  - `start`: Launch trading bot
  - `stop`: Gracefully stop trading
  - `status`: Display current trading status
  - `logs`: View and filter trading logs
  - `config`: Manage bot configuration

### 3.3 Exchange Integration
- MEXC Exchange API integration
  - REST API for order placement and account information
  - WebSocket support for real-time price updates
- Fallback mechanisms for API failures
- Rate limiting and API request optimization

### 3.4 Data Management
- SQLite database for:
  - Storing trading positions
  - Logging trade events
  - Tracking new coin discoveries
- Soft delete mechanism for trading records
- Configurable data retention policies

## 4. Non-Functional Requirements

### 4.1 Performance
- Low-latency trading execution
- Efficient concurrent processing using Go goroutines
- Minimal resource consumption
- Caching mechanism for API responses

### 4.2 Reliability
- Automatic recovery from temporary network failures
- Exponential backoff for API retries
- Graceful shutdown preserving trading state
- Comprehensive error handling and logging

### 4.3 Security
- Secure storage of API credentials
- Environment variable or configuration file support
- Prevention of SQL injection
- Secure WebSocket communication

### 4.4 Scalability
- Modular architecture supporting easy extension
- Potential future support for multiple exchanges
- Configurable trading strategies
- Efficient memory management

## 5. Technical Specifications

### 5.1 Technology Stack
- **Language**: Go 1.18+
- **Database**: SQLite
- **Libraries**:
  - Gin/Echo (Web Framework)
  - Gorilla WebSocket
  - Viper (Configuration)
  - Logrus/Zap (Logging)
  - GORM/SQLx (Database ORM)

### 5.2 Project Structure
- Layered architecture with:
  - Domain layer (business logic)
  - Platform layer (infrastructure)
  - Application layer (API/CLI)

## 6. Configuration Options

### 6.1 Trading Configuration
- Stop-loss percentage
- Take-profit levels
- Trading amount per transaction
- Coin selection criteria
- Polling intervals

### 6.2 Exchange Configuration
- API keys
- Rate limit settings
- Supported trading pairs

### 6.3 System Configuration
- Logging level
- Database path
- Caching settings

## 7. Risks and Mitigations

### 7.1 Market Risks
- Cryptocurrency price volatility
- Mitigation: Robust stop-loss and take-profit mechanisms

### 7.2 Technical Risks
- API changes by MEXC
- Mitigation: Flexible API client, comprehensive error handling

### 7.3 Compliance Risks
- Regulatory changes in cryptocurrency trading
- Mitigation: Configurable trading limits, logging for auditing

## 8. Future Roadmap
- Machine learning trading strategy optimization
- Multi-exchange support
- Advanced risk management algorithms
- Enhanced visualization and reporting
- Mobile app companion

## 9. Success Metrics
- Trading performance metrics
- Bot uptime
- Number of successful trades
- Risk-adjusted return
- User adoption and feedback

## 10. Acceptance Criteria
- Autonomous trading execution
- Less than 5% drawdown in volatile markets
- 99.9% uptime
- Configurable and extensible architecture
- Comprehensive logging and error tracking

---

This PRD provides a comprehensive overview of the Go Crypto Trading Bot, outlining its purpose, features, technical requirements, and future potential. The document serves as a strategic guide for development, ensuring a robust, flexible, and user-centric trading solution.