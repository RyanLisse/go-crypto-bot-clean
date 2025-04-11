# Product Requirements Document: Railway Deployment

## Overview
This document outlines the requirements and implementation strategy for deploying the Go Crypto Bot backend to Railway, a modern cloud platform for deploying applications. Due to initial deployment challenges, we've adopted an incremental approach that gradually adds components to identify and resolve issues.

## Goals
- Deploy the Go Crypto Bot backend to Railway using an incremental approach
- Ensure the application is accessible via a public URL
- Set up proper health checks for monitoring
- Implement a deployment process that can be automated
- Document the deployment process for future reference

## Non-Goals
- Setting up a CI/CD pipeline (will be addressed in a separate task)
- Implementing advanced monitoring and alerting
- Setting up custom domains (will be addressed in a separate task)

## Current Progress

The application is currently deployed at https://piquant-desire-production.up.railway.app with the following components:

- ✅ Basic API structure with Chi router
- ✅ Health check endpoint
- ✅ Configuration management with Viper
- ✅ Structured logging with Zap
- ⏳ SQLite database integration
- ⏳ Turso cloud database integration
- ⏳ Clerk authentication
- ⏳ Google Gemini AI integration
- ⏳ External service integrations

## System Requirements

### Core Dependencies
- Go version: 1.24.2
- SQLite with CGO support
- Turso for cloud database (with local SQLite sync)

### Key Dependencies
- Database:
  - github.com/mattn/go-sqlite3 v1.14.27 (Local SQLite)
  - github.com/jmoiron/sqlx v1.4.0 (SQL extensions)
  - gorm.io/driver/sqlite v1.5.7 (GORM SQLite driver)
  - gorm.io/driver/postgres v1.5.11 (GORM Postgres driver)
  - gorm.io/gorm v1.25.12 (ORM framework)
  - libsql/libsql-client-go (Turso client - to be added)

- API and Routing:
  - github.com/go-chi/chi/v5 v5.2.1
  - github.com/danielgtaylor/huma/v2 v2.32.0
  - github.com/gorilla/websocket v1.5.3

- Authentication:
  - github.com/clerk/clerk-sdk-go/v2 v2.3.0 (Primary authentication provider)
  - github.com/golang-jwt/jwt/v5 v5.2.2 (JWT token handling)

- AI and External Services:
  - github.com/google/generative-ai-go v0.19.0 (Google Gemini AI integration)
  - google.golang.org/api v0.228.0 (Google API client)
  - github.com/sashabaranov/go-openai v1.38.1 (OpenAI integration)
  - github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 (Telegram bot)
  - github.com/slack-go/slack v0.16.0 (Slack integration)

- Configuration and CLI:
  - github.com/spf13/viper v1.20.1
  - github.com/spf13/cobra v1.9.1

- Logging and Monitoring:
  - github.com/sirupsen/logrus v1.9.3
  - go.uber.org/zap v1.27.0

## Deployment Requirements

### Environment Variables
- DATABASE_URL: SQLite database path
- PORT: Application port (default: 8080)
- TURSO_ENABLED: Enable Turso database (true/false)
- TURSO_URL: Turso database URL
- TURSO_AUTH_TOKEN: Turso authentication token
- TURSO_SYNC_ENABLED: Enable local/remote sync (true/false)
- TURSO_SYNC_INTERVAL_SECONDS: Sync interval in seconds
- CLERK_SECRET_KEY: Clerk authentication secret
- OPENAI_API_KEY: OpenAI API key
- GEMINI_API_KEY: Google Gemini API key
- TELEGRAM_BOT_TOKEN: Telegram Bot API token
- SLACK_BOT_TOKEN: Slack Bot token
- GOOGLE_APPLICATION_CREDENTIALS: Path to Google Cloud credentials file

### Resource Requirements
- Minimum Memory: 512MB
- Recommended Memory: 1GB
- Storage: At least 1GB for SQLite database and application files

### Network Requirements
- Exposed Ports: 8080 (HTTP)
- Outbound access needed for:
  - api.openai.com
  - api.telegram.org
  - slack.com
  - clerk.dev
  - googleapis.com

### Health Checks
- HTTP endpoint: /health
- Expected response: 200 OK
- Timeout: 5 seconds
- Interval: 30 seconds

## Implementation Strategy

### Phase 1: Basic Deployment
1. Configure SQLite with CGO support in Docker
2. Set up basic health check endpoint
3. Configure essential environment variables
4. Deploy minimal version with health check

### Phase 2: Database Integration
1. Add SQLite local database support
2. Configure Turso cloud database integration
3. Set up database synchronization
4. Implement basic data models and repositories

### Phase 3: Authentication and Security
1. Integrate Clerk authentication
2. Set up JWT token validation
3. Implement authorization middleware
4. Configure secure API endpoints

### Phase 4: External Services Integration
1. Add Google Gemini AI integration
2. Configure OpenAI integration
3. Set up Telegram bot connection
4. Configure Slack integration

### Phase 5: Monitoring and Optimization
1. Implement detailed logging
2. Add performance monitoring
3. Configure resource scaling
4. Set up backup strategy for SQLite and Turso databases

## Success Criteria
1. Application successfully deploys and starts
2. Health check endpoint responds correctly
3. Local SQLite database operations work as expected
4. Turso cloud database integration functions properly
5. Clerk authentication system is operational
6. JWT token validation works correctly
7. Google Gemini AI integration functions properly
8. OpenAI integration works as expected
9. Telegram and Slack integrations are operational
10. Logs are properly captured and accessible
11. Application can handle expected load

## Rollback Plan
1. Keep previous deployment available
2. Maintain SQLite and Turso database backups
3. Document version-specific configuration
4. Have quick rollback procedure ready
5. Ensure Clerk authentication can be reverted if needed

## Future Considerations
1. Implementation of CI/CD pipeline
2. Setting up custom domain
3. Advanced monitoring and alerting
4. Automated backup system for both SQLite and Turso
5. Load balancing configuration
6. Enhanced security features
7. Improved AI model integration
8. Multi-region deployment
