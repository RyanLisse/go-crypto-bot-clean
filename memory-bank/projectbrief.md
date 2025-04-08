# Go Crypto Bot: Project Brief

## Overview
This project involves migrating an existing Python cryptocurrency trading bot to Go. The Go implementation will follow hexagonal architecture principles and support both backend server and CLI functionalities.

## Key Goals
1. Implement all functionality from the existing Python codebase in Go
2. Apply hexagonal architecture for better separation of concerns
3. Create both server and CLI interfaces
4. Integrate with MEXC cryptocurrency exchange API
5. Implement robust data persistence using SQLite
6. Support portfolio management, trading decisions, and trade execution

## Technical Requirements
- Clean architecture with clearly defined domain boundaries
- Separation between business logic and external dependencies
- Comprehensive testing strategy
- Proper error handling and logging
- Documentation of key components

## Success Criteria
- Feature parity with the Python implementation
- Improved maintainability and testability
- Ability to easily swap infrastructure components (database, exchange APIs, etc.)
- Proper adherence to Go best practices
