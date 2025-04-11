# Railway Deployment Guide

This guide provides instructions for deploying the Go Crypto Bot backend to Railway.

## Prerequisites

- Railway CLI installed (`npm install -g @railway/cli`)
- Railway account with access to the project
- Docker installed locally for testing

## Project Structure

The project is configured for deployment with the following components:

- **Backend**: Go API with SQLite database

## Configuration Files

- `backend/railway.toml`: Backend service configuration

## Environment Variables

### Backend Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| PORT | Port for the HTTP server | Yes | 8080 |
| ENVIRONMENT | Deployment environment | Yes | production |
| LOG_LEVEL | Logging level | Yes | info |
| CONFIG_PATH | Path to configuration files | Yes | /app/configs |
| CONFIG_FILE | Configuration file name | Yes | config.minimal.yaml |
| DB_PATH | Path to SQLite database | Yes | /app/data/minimal.db |
| DATABASE_ENABLED | Enable database | Yes | true |
| AUTH_ENABLED | Enable authentication | No | false |
| CLERK_SECRET_KEY | Clerk API secret key | No | - |
| CLERK_DOMAIN | Clerk domain | No | - |
| TURSO_ENABLED | Enable Turso database | No | false |
| TURSO_URL | Turso database URL | No | - |
| TURSO_AUTH_TOKEN | Turso authentication token | No | - |
| TURSO_SYNC_ENABLED | Enable Turso sync | No | false |
| MEXC_API_KEY | MEXC API key | No | - |
| MEXC_SECRET_KEY | MEXC API secret | No | - |
| OPENAI_API_KEY | OpenAI API key | No | - |
| GOOGLE_API_KEY | Google API key | No | - |

## Deployment Steps

### Initial Setup

1. Login to Railway:
   ```bash
   railway login
   ```

2. Link to the project:
   ```bash
   railway link
   ```

### Deploying the Backend

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Deploy to Railway:
   ```bash
   railway up
   ```

## Monitoring and Maintenance

### Health Checks

- Backend: `/health` and `/health/detailed` endpoints

### Logs

Access logs through the Railway dashboard or CLI:

```bash
railway logs
```

### Resource Monitoring

Monitor resource usage through the Railway dashboard.

## Rollback Procedure

If a deployment fails or causes issues:

1. Identify the last working deployment in the Railway dashboard
2. Roll back to that deployment:
   ```bash
   railway rollback --to <deployment-id>
   ```

## Troubleshooting

### Common Issues

1. **Health Check Failures**:
   - Verify the health check endpoint is correct
   - Check application logs for errors

2. **Database Connection Issues**:
   - Verify database path and permissions
   - Check if the database file exists

3. **Environment Variable Problems**:
   - Ensure all required environment variables are set
   - Check for typos in variable names

### Getting Help

For additional assistance, contact the DevOps team or refer to the Railway documentation at https://docs.railway.app/.
