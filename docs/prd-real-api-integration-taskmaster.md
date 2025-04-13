<context>
# Overview  
The crypto trading bot application currently uses mock data for account, wallet, and portfolio information. This project aims to replace these mock implementations with real API calls to the MEXC exchange, ensuring that users see accurate, real-time data in the frontend.

# Core Features  
1. **Real MEXC API Integration**
   - What it does: Connects to the MEXC exchange API to fetch real account data
   - Why it's important: Provides users with accurate information about their accounts and trades
   - How it works: Uses the MEXC REST API client to fetch data and properly handle errors

2. **Error Handling and Resilience**
   - What it does: Gracefully handles API errors and outages
   - Why it's important: Ensures the application remains usable even when the exchange API is unavailable
   - How it works: Implements proper error handling, logging, and fallback mechanisms

3. **Performance Optimization**
   - What it does: Implements caching and rate limiting
   - Why it's important: Minimizes API calls to the exchange and ensures compliance with rate limits
   - How it works: Uses caching mechanisms with appropriate TTL values and configures rate limiters

# User Experience  
- **User Personas**: Crypto traders who need accurate, real-time information about their accounts and trades
- **Key User Flows**: Viewing account balances, portfolio summary, and trade performance
- **UI/UX Considerations**: Ensuring the frontend displays real data with appropriate loading states and error messages
</context>
<PRD>
# Technical Architecture  
## System Components
1. **Account Service**
   - Responsible for fetching account information from the MEXC API
   - Implements caching to minimize API calls
   - Handles API errors gracefully

2. **Portfolio Service**
   - Calculates portfolio value and performance metrics using real data
   - Integrates with the Account Service to get wallet information
   - Updates trade information with current prices

3. **API Layer**
   - Exposes RESTful endpoints for the frontend
   - Handles authentication and request validation
   - Returns appropriate error responses

## Data Models
1. **Wallet**
   - Contains balances for all assets
   - Includes free, locked, and total amounts
   - Stores the last update timestamp

2. **Balance**
   - Represents the account balance
   - Includes fiat and crypto assets
   - Tracks available and locked amounts

3. **Trade/Position**
   - Represents an active trade
   - Includes entry price, current price, and profit/loss
   - Stores stop-loss and take-profit levels

## APIs and Integrations
1. **MEXC REST API**
   - Used for fetching account information, wallet data, and market prices
   - Requires API key and secret for authentication
   - Subject to rate limits

2. **MEXC WebSocket API**
   - Used for real-time updates (future enhancement)
   - Requires authentication for account-specific data
   - Provides efficient real-time data streaming

## Infrastructure Requirements
1. **Configuration Management**
   - Secure storage for API keys
   - Environment-specific configuration
   - Feature flags for gradual rollout

2. **Logging and Monitoring**
   - Detailed logging of API interactions
   - Monitoring of API response times and error rates
   - Alerts for API outages or rate limit issues

# Development Roadmap  
## Phase 1: Foundation
1. **Fix Real Account Service Initialization**
   - Modify the `initializeRealAccountService` method in `dependencies.go`
   - Implement proper error handling and logging
   - Add validation of API keys during initialization

2. **Update API Handlers**
   - Remove fallback to mock data in account handlers
   - Implement proper error handling
   - Add appropriate HTTP status codes for different error scenarios

## Phase 2: Core Functionality
1. **Implement Real Portfolio Service**
   - Create a real implementation of the portfolio service
   - Integrate with the account service to get wallet data
   - Calculate performance metrics using real data

2. **Optimize Caching**
   - Configure appropriate TTL values for different types of data
   - Implement cache invalidation strategies
   - Add cache warming for frequently accessed data

## Phase 3: Reliability and Performance
1. **Enhance Error Handling**
   - Implement circuit breaker pattern for API calls
   - Add retry mechanisms with exponential backoff
   - Create meaningful error messages for the frontend

2. **Optimize Rate Limiting**
   - Configure rate limiters based on MEXC API documentation
   - Implement request prioritization
   - Add monitoring for rate limit usage

## Phase 4: Testing and Validation
1. **Implement Comprehensive Testing**
   - Unit tests for services
   - Integration tests for API endpoints
   - End-to-end tests with the frontend

2. **Validation and Monitoring**
   - Validate data accuracy against the MEXC web interface
   - Monitor API response times and error rates
   - Set up alerts for API outages or rate limit issues

# Logical Dependency Chain
1. **Foundation First**
   - Fix the real account service initialization
   - Update API handlers to use real data
   - These changes provide the foundation for all other enhancements

2. **Core Functionality Next**
   - Implement the real portfolio service
   - Optimize caching
   - These changes ensure that all endpoints return real data

3. **Reliability and Performance**
   - Enhance error handling
   - Optimize rate limiting
   - These changes ensure the application remains reliable and performant

4. **Testing and Validation Last**
   - Implement comprehensive testing
   - Validate data accuracy
   - These changes ensure the application works correctly in all scenarios

# Risks and Mitigations  
## Technical Challenges
1. **API Rate Limits**
   - Risk: Exceeding MEXC API rate limits
   - Mitigation: Implement proper rate limiting and caching

2. **API Changes**
   - Risk: MEXC API changes breaking the integration
   - Mitigation: Implement versioned API clients and monitor for changes

3. **API Outages**
   - Risk: MEXC API outages affecting the application
   - Mitigation: Implement circuit breaker pattern and graceful degradation

## MVP Considerations
1. **Scope Management**
   - Risk: Scope creep delaying the implementation
   - Mitigation: Focus on the core functionality first, then add enhancements

2. **Testing Complexity**
   - Risk: Difficulty testing with real API calls
   - Mitigation: Implement proper mocking for tests and use a staging environment

## Resource Constraints
1. **API Key Management**
   - Risk: Secure storage and transmission of API keys
   - Mitigation: Use environment variables and secure storage solutions

2. **Performance Impact**
   - Risk: Real API calls affecting performance
   - Mitigation: Implement proper caching and optimization

# Appendix  
## Research Findings
1. **MEXC API Documentation**
   - The MEXC API provides endpoints for account information, wallet data, and market prices
   - The API has rate limits that need to be respected
   - The API requires authentication for account-specific data

2. **Current Implementation Analysis**
   - The application currently falls back to mock data when the real API client initialization fails
   - The mock data is not representative of the real account state
   - The fallback mechanism is not configurable

## Technical Specifications
1. **API Endpoints**
   - `/api/v1/account` - Get account information
   - `/api/v1/account/balance` - Get account balances
   - `/api/v1/account/wallet` - Get wallet information
   - `/api/v1/portfolio` - Get portfolio summary
   - `/api/v1/portfolio/value` - Get total portfolio value

2. **Error Handling**
   - HTTP 500 for server errors
   - HTTP 400 for client errors
   - HTTP 401 for authentication errors
   - HTTP 429 for rate limit errors
</PRD>
