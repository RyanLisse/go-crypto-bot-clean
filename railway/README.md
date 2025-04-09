# Crypto Bot API - Railway Deployment

This directory contains a Railway-optimized version of the Crypto Bot API backend.

## Project Structure

```
railway/
├── cmd/
│   └── api/
│       └── main.go       # Entry point with Railway-specific configuration
├── internal/
│   ├── app/              # Application initialization
│   ├── config/           # Configuration management
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # Middleware components
│   ├── models/           # Data models
│   ├── server/           # HTTP server setup
│   └── services/         # Business logic
├── pkg/                  # Reusable packages
├── .env.example          # Example environment variables
├── Dockerfile            # Optimized for Railway
└── go.mod                # Module definition
```

## Local Development

1. Copy `.env.example` to `.env` and adjust values as needed
2. Run the application:

```bash
cd railway
go run cmd/api/main.go
```

## Deployment to Railway

### Option 1: Deploy via Railway CLI

1. Install the Railway CLI:
```bash
npm i -g @railway/cli
```

2. Login to Railway:
```bash
railway login
```

3. Link to your Railway project:
```bash
railway link
```

4. Deploy the application:
```bash
railway up
```

### Option 2: Deploy via GitHub Integration

1. Push your code to GitHub
2. In the Railway dashboard, create a new project
3. Select "Deploy from GitHub repo"
4. Configure the project:
   - Set the root directory to `railway`
   - Railway will automatically detect the Dockerfile

## Environment Variables

The following environment variables can be configured in Railway:

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | The port the server will listen on | 8080 |
| DATABASE_URL | Database connection string | sqlite3://crypto-bot.db |
| APP_LOG_LEVEL | Logging level (debug, info, warn, error) | info |
| APP_API_BASE_PATH | Base path for API endpoints | /api/v1 |

## Adding Custom Services

To add new services to the API:

1. Create a new service in `internal/services/`
2. Create corresponding handlers in `internal/handlers/`
3. Register routes in `internal/server/server.go`
