# Docker Setup for Go Crypto Bot

This document provides instructions for running the Go Crypto Bot application using Docker containers.

## Prerequisites

- Docker and Docker Compose installed on your system
- Git repository cloned to your local machine

## Environment Variables

Before running the application, you need to set up your environment variables. Copy the example file and update it with your values:

```bash
cp .env.example .env
```

Edit the `.env` file and add your API keys and other configuration values.

## Running the Application

### Production Mode

To run the full application in production mode:

```bash
make up
```

This will start both the backend and frontend containers.

To run only the backend:

```bash
make backend
```

To run only the frontend:

```bash
make frontend
```

### Development Mode

For development with hot reloading:

```bash
make dev
```

This will start both the backend and frontend in development mode with hot reloading.

To run only the backend in development mode:

```bash
make backend-dev
```

To run only the frontend in development mode:

```bash
make frontend-dev
```

### Viewing Logs

To view logs from all containers:

```bash
make logs
```

### Stopping the Application

To stop all containers:

```bash
make down
```

### Cleaning Up

To remove all containers, volumes, and images:

```bash
make clean
```

## Container Structure

### Backend Container

- Built from `backend/Dockerfile`
- Runs on port 8080
- Stores data in a persistent volume
- Uses SQLite for local development
- Can connect to Turso cloud database

### Frontend Container

- Built from `frontend/Dockerfile`
- Runs on port 3000
- Uses Nginx to serve static files
- Proxies API requests to the backend

### Development Containers

- Backend uses Air for hot reloading
- Frontend uses Next.js development server
- Volumes are mounted for real-time code changes

## Customizing the Setup

You can modify the Docker Compose files to add additional services or change configuration:

- `docker-compose.yml` - Production setup
- `docker-compose.dev.yml` - Development setup

## Troubleshooting

### Container Won't Start

Check the logs:

```bash
docker-compose logs [service_name]
```

### API Connection Issues

Make sure the environment variables are set correctly in the `.env` file.

### Volume Permissions

If you encounter permission issues with volumes:

```bash
sudo chown -R $(id -u):$(id -g) ./data
```
