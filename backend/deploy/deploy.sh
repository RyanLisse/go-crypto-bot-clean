#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Print header
echo -e "${GREEN}=== Crypto Bot Backend Deployment Script ===${NC}"
echo -e "${YELLOW}This script will deploy the backend to production${NC}"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

# Check if .env file exists
if [ ! -f "../.env" ]; then
    echo -e "${RED}Error: .env file not found. Please create it from .env.template${NC}"
    exit 1
fi

# Validate environment variables
echo -e "${YELLOW}Validating environment variables...${NC}"
source ../.env

# Check required environment variables
if [ -z "$MEXC_API_KEY" ] || [ -z "$MEXC_SECRET_KEY" ]; then
    echo -e "${RED}Error: MEXC API credentials are missing in .env file${NC}"
    exit 1
fi

if [ -z "$TURSO_URL" ] || [ -z "$TURSO_AUTH_TOKEN" ]; then
    echo -e "${RED}Error: TursoDB credentials are missing in .env file${NC}"
    exit 1
fi

if [ -z "$JWT_SECRET" ]; then
    echo -e "${YELLOW}Warning: JWT_SECRET is not set. Generating a random one...${NC}"
    JWT_SECRET=$(openssl rand -base64 32)
    echo "JWT_SECRET=$JWT_SECRET" >> ../.env
fi

echo -e "${GREEN}Environment validation passed${NC}"

# Build the Docker image
echo -e "${YELLOW}Building Docker image...${NC}"
cd ..
docker build -t crypto-bot-backend:latest .

# Tag the image for deployment
echo -e "${YELLOW}Tagging image for deployment...${NC}"
TIMESTAMP=$(date +%Y%m%d%H%M%S)
docker tag crypto-bot-backend:latest crypto-bot-backend:$TIMESTAMP

# Create a secure .env file for production
echo -e "${YELLOW}Creating secure production environment file...${NC}"
grep -v "^#" .env > .env.production
chmod 600 .env.production

echo -e "${GREEN}Deployment preparation complete!${NC}"
echo -e "${YELLOW}To deploy to your production server, run:${NC}"
echo -e "  docker save crypto-bot-backend:$TIMESTAMP | ssh user@your-server 'docker load'"
echo -e "  scp .env.production user@your-server:/path/to/deployment/"
echo -e "  ssh user@your-server 'docker stop crypto-bot || true && docker rm crypto-bot || true'"
echo -e "  ssh user@your-server 'docker run -d --name crypto-bot --restart always --env-file /path/to/deployment/.env.production -p 8080:8080 crypto-bot-backend:$TIMESTAMP'"

echo -e "${GREEN}Done!${NC}"
