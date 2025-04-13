# Product Requirements Document: Real MEXC API Integration

## Overview
This document outlines the requirements and implementation plan for integrating real MEXC API data into the backend of the crypto trading bot application. Currently, the application is using mock data for account, wallet, and portfolio information. This PRD defines the approach to replace mock implementations with real API calls.

## Current Backend Structure
The backend follows a hexagonal architecture with the following key components:

### API Layer
- Located in `backend/internal/api`
- Handles HTTP requests and responses
- Contains handlers for various endpoints including account, portfolio, and wallet

### Core Layer
- Located in `backend/internal/core`
- Contains business logic and service implementations
- Key services:
  - `account` - Manages account-related operations
  - `portfolio` - Manages portfolio-related operations
  - `trade` - Manages trading operations

### Platform Layer
- Located in `backend/internal/platform`
- Contains external integrations including MEXC API client
- Key components:
  - `mexc/rest` - REST API client for MEXC
  - `mexc/websocket` - WebSocket client for MEXC
  - `database` - Database access and repositories

### Domain Layer
- Located in `backend/internal/domain`
- Contains domain models and interfaces
- Defines contracts between layers

## Problem Statement
Currently, the application falls back to mock implementations when the real MEXC API client initialization fails. This results in the frontend displaying mock data instead of real account information from the MEXC exchange.

## Requirements

### Functional Requirements
1. The backend must use real MEXC API data for the following endpoints:
   - `/api/v1/account` - Get account information
   - `/api/v1/account/balance` - Get account balances
   - `/api/v1/account/wallet` - Get wallet information
   - `/api/v1/portfolio` - Get portfolio summary
   - `/api/v1/portfolio/value` - Get total portfolio value

2. The application must properly handle API errors and provide meaningful error messages to the frontend.

3. The application must implement proper caching to minimize API calls to the MEXC exchange.

4. The application must implement proper rate limiting to avoid exceeding MEXC API limits.

### Non-Functional Requirements
1. Performance: API responses should be returned within 500ms.
2. Reliability: The application should handle MEXC API outages gracefully.
3. Security: API keys must be securely stored and transmitted.

## Implementation Approach

### 1. Fix Real Account Service Initialization
- Modify the `initializeRealAccountService` method in `dependencies.go` to properly initialize the real account service
- Ensure proper error handling and logging
- Remove fallback to mock services unless explicitly configured

### 2. Update API Handlers
- Ensure all handlers are using the real account service
- Remove any fallback to mock data in the handlers
- Implement proper error handling

### 3. Implement Caching
- Use the existing caching mechanism in the real account service
- Configure appropriate TTL values for different types of data

### 4. Implement Rate Limiting
- Use the existing rate limiting mechanism in the MEXC client
- Configure appropriate rate limits based on MEXC API documentation

## Testing Plan
1. Unit tests for the real account service
2. Integration tests for the API endpoints
3. End-to-end tests with the frontend

## Success Criteria
1. All specified API endpoints return real data from the MEXC exchange
2. The application handles API errors gracefully
3. The application respects MEXC API rate limits
4. The frontend displays real account information

## Future Enhancements
1. Implement WebSocket for real-time updates
2. Add support for additional exchanges
3. Implement more sophisticated caching strategies
