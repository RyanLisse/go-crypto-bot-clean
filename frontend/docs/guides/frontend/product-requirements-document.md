# Crypto Trading Bot Frontend - Product Requirements Document

## 1. Introduction

### 1.1 Purpose
This document outlines the requirements for the Crypto Trading Bot Frontend, a Vite/React web application with a brutalist design that provides a user-friendly interface for interacting with the Crypto Trading Bot backend API. The frontend enables users to monitor their portfolio, execute trades, configure the bot, view real-time market data, and interact with an AI assistant.

### 1.2 Scope
The frontend application provides a complete user interface for all functionality exposed by the Crypto Trading Bot API, including authentication, portfolio management, trading, new coin detection, system configuration, and AI-powered assistance.

### 1.3 Definitions and Acronyms
- **JWT**: JSON Web Token, used for authentication
- **API**: Application Programming Interface
- **WebSocket**: Protocol for real-time communication
- **SPA**: Single Page Application
- **UI**: User Interface
- **UX**: User Experience
- **TanStack Query**: Data fetching and caching library (formerly React Query)
- **Gemini**: Google's AI model used for the chat assistant

## 2. Product Overview

### 2.1 Product Perspective
The Crypto Trading Bot Frontend is a Vite/React web application that communicates with the Crypto Trading Bot backend API. It provides a visual representation of the data and functionality offered by the API, making it accessible to users without technical knowledge of API interactions. The application follows a brutalist design approach with a focus on minimalism, high contrast, and monospace typography.

### 2.2 User Classes and Characteristics
1. **Traders**: Users who actively manage their portfolio and execute trades
2. **Investors**: Users who monitor their portfolio but rely on the bot for automated trading
3. **Administrators**: Users who configure and maintain the system
4. **Analysts**: Users who analyze trading performance and market trends

### 2.3 Operating Environment
- Modern web browsers (Chrome, Firefox, Safari, Edge)
- Desktop and mobile devices
- Internet connection required
- Node.js v18+ for development

### 2.4 Design and Implementation Constraints
- Must follow brutalist design principles with monospace typography
- Must be responsive and work on both desktop and mobile devices
- Must implement secure authentication using JWT
- Must handle real-time updates via WebSocket connections
- Must integrate with Google Gemini 1.5 Flash for AI assistance
- Must use TanStack Query for data fetching and caching
- Must be accessible and follow WCAG 2.1 AA standards

### 2.5 Assumptions and Dependencies
- Requires the Crypto Trading Bot backend API to be operational
- Assumes the API follows the documented endpoints and response formats
- Depends on modern browser features like WebSockets and localStorage
- Requires Google Gemini API key for AI assistant functionality
- Uses Bun as the preferred package manager

## 3. System Features

### 3.1 Authentication
#### 3.1.1 Description
The system must provide secure authentication using JWT tokens.

#### 3.1.2 Requirements
- User login form with username and password fields
- JWT token storage in localStorage or cookies
- Automatic token refresh mechanism
- Logout functionality
- Protected routes for authenticated users
- Role-based access control (admin vs. regular user)

### 3.2 Dashboard
#### 3.2.1 Description
The main dashboard provides an overview of the user's portfolio, recent trades, and market data.

#### 3.2.2 Requirements
- Portfolio summary with total value and performance metrics
- Active trades list with current profit/loss
- Market overview with key cryptocurrency prices
- Real-time updates via WebSocket
- Performance charts (daily, weekly, monthly)
- System status indicators

### 3.3 Portfolio Management
#### 3.3.1 Description
Allows users to view and manage their cryptocurrency portfolio.

#### 3.3.2 Requirements
- Detailed list of all holdings with current values
- Historical performance charts
- Position details (entry price, current price, profit/loss)
- Filtering and sorting options
- Export functionality (CSV, PDF)

### 3.4 Trading Interface
#### 3.4.1 Description
Enables users to execute manual trades and view trade history.

#### 3.4.2 Requirements
- Buy/sell form with symbol, quantity, and order type
- Market data display (price, volume, 24h change)
- Order book visualization
- Trade history with filtering options
- Real-time price updates
- Confirmation dialogs for trade execution

### 3.5 New Coin Detection
#### 3.5.1 Description
Displays newly detected coins and allows users to configure detection parameters.

#### 3.5.2 Requirements
- List of newly detected coins with key metrics
- Real-time alerts for new coin listings
- Configuration options for detection parameters
- Manual trigger for coin detection
- Historical data on previously detected coins

### 3.6 Bot Configuration
#### 3.6.1 Description
Allows users to configure the trading bot's behavior.

#### 3.6.2 Requirements
- Strategy selection and configuration
- Risk management settings
- Trading pair selection
- Schedule configuration
- Parameter validation
- Configuration history and versioning

### 3.7 System Status and Logs
#### 3.7.1 Description
Provides information about the system status and operation logs.

#### 3.7.2 Requirements
- System health indicators
- Process status (running/stopped)
- Error logs and warnings
- Performance metrics
- Start/stop controls for system processes

### 3.8 AI Assistant
#### 3.8.1 Description
Provides an AI-powered chat interface that helps users understand their portfolio, analyze market trends, and get insights on trading strategies.

#### 3.8.2 Requirements
- Chat interface for interacting with the AI assistant
- Integration with Google Gemini 1.5 Flash model
- Context-aware responses based on user's portfolio and trading data
- Ability to ask questions about portfolio performance
- Ability to get explanations about trading strategies
- Ability to analyze market trends and provide insights
- Conversation history preservation during the session
- Clear indication of AI processing status
- Error handling for API failures

### 3.9 WebSocket Integration
#### 3.9.1 Description
Enables real-time updates throughout the application.

#### 3.9.2 Requirements
- Connection status indicator
- Automatic reconnection
- Market data updates
- Trade notifications
- New coin alerts
- Portfolio value updates

## 4. External Interface Requirements

### 4.1 User Interfaces
- Modern, clean design following material design principles
- Responsive layout for desktop and mobile devices
- Dark/light theme options
- Customizable dashboard layouts
- Accessibility compliance (WCAG 2.1 AA)

### 4.2 API Interfaces
- RESTful API communication for data retrieval and actions
- WebSocket connection for real-time updates
- JWT authentication for all API requests
- Error handling and retry mechanisms

### 4.3 Hardware Interfaces
- Support for standard input devices (mouse, keyboard, touch)
- Responsive design for various screen sizes

## 5. Non-Functional Requirements

### 5.1 Performance
- Initial load time under 3 seconds on broadband connections
- Real-time updates with less than 500ms latency
- Smooth animations and transitions (60fps)
- Efficient data caching to minimize API calls

### 5.2 Security
- Secure authentication using JWT
- HTTPS for all communications
- Protection against common web vulnerabilities (XSS, CSRF)
- Secure storage of sensitive information
- Session timeout after period of inactivity

### 5.3 Reliability
- Graceful handling of API errors
- Offline mode for basic functionality
- Data persistence across page refreshes
- Automatic recovery from WebSocket disconnections

### 5.4 Usability
- Intuitive navigation and controls
- Consistent design language throughout the application
- Helpful error messages and guidance
- Tooltips and help documentation
- Progressive disclosure of complex features

### 5.5 Compatibility
- Support for modern browsers (last 2 versions)
- Mobile-friendly design
- Responsive layouts for various screen sizes

## 6. Technical Stack

### 6.1 Frontend Framework
- Vite for fast development and optimized builds
- React 18+ for component-based UI development
- React Router for client-side routing
- React Context API and TanStack Query for state management

### 6.2 UI Components
- Tailwind CSS for styling with custom brutalist theme
- JetBrains Mono as the primary font
- Recharts for data visualization
- Sonner for toast notifications
- Radix UI primitives with shadcn/ui for accessible components

### 6.3 API Communication and Data Persistence
- TanStack Query for data fetching, caching, and synchronization
- Native Fetch API for REST requests
- Native WebSocket for real-time communication
- Google Generative AI SDK for Gemini integration
- Drizzle ORM with SQLite for local data persistence and offline support

### 6.4 Build Tools and Testing
- Bun as the package manager and runtime
- Vite for bundling and optimization
- ESLint for code quality
- TypeScript for type safety
- Vitest for unit testing (optional)
- Playwright for end-to-end testing (optional)

## 7. Implementation Phases

### 7.1 Phase 1: Core Infrastructure
- Authentication system
- Basic dashboard
- API integration layer
- WebSocket connection

### 7.2 Phase 2: Portfolio and Trading
- Portfolio management
- Trading interface
- Trade history
- Real-time market data

### 7.3 Phase 3: Bot Configuration
- Strategy configuration
- Risk management settings
- System status and controls

### 7.4 Phase 4: Advanced Features
- New coin detection interface
- Performance analytics
- Customizable dashboards
- Mobile optimization

## 8. Appendices

### 8.1 Mockups and Wireframes
[To be developed based on the requirements]

### 8.2 API Documentation
The frontend will integrate with the following API endpoints:

#### Authentication
- POST /auth/login
- POST /auth/logout
- GET /auth/me

#### Portfolio
- GET /api/v1/portfolio
- GET /api/v1/portfolio/active
- GET /api/v1/portfolio/performance
- GET /api/v1/portfolio/value

#### Trading
- GET /api/v1/trade/history
- POST /api/v1/trade/buy
- POST /api/v1/trade/sell
- GET /api/v1/trade/status/:id

#### New Coins
- GET /api/v1/newcoins
- POST /api/v1/newcoins/process
- POST /api/v1/newcoins/detect

#### Configuration
- GET /api/v1/config
- PUT /api/v1/config
- GET /api/v1/config/defaults

#### System Status
- GET /api/v1/status
- POST /api/v1/status/start
- POST /api/v1/status/stop

#### WebSocket
- WS /ws

### 8.3 User Stories
1. As a trader, I want to see my portfolio value in real-time so I can monitor my investments.
2. As an investor, I want to configure the bot's trading strategy so it aligns with my risk tolerance.
3. As an administrator, I want to view system logs so I can troubleshoot issues.
4. As a trader, I want to receive notifications about new coin listings so I can make timely investment decisions.
5. As a user, I want to securely log in to the system so my trading data remains private.
6. As an investor, I want to see historical performance charts so I can evaluate the bot's effectiveness.
7. As a trader, I want to manually execute trades when I see opportunities the bot might miss.
8. As an administrator, I want to start and stop the bot processes so I can perform maintenance.
9. As a user, I want the interface to work on my mobile device so I can monitor my portfolio on the go.
10. As an investor, I want to set risk management parameters so I can protect my capital.

---

This PRD will guide the development of the Crypto Trading Bot Frontend, ensuring that all stakeholders have a clear understanding of the requirements and functionality to be implemented.
