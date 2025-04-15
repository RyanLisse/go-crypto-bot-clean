# Build and Deployment Guide

This document outlines the steps to build, run, and deploy the Next.js frontend for the Crypto Bot application.

## Environment Setup

Before building or running the application, ensure you have:

1. **Bun Runtime**: This project uses Bun for improved performance
   ```
   curl -fsSL https://bun.sh/install | bash
   ```

2. **Node.js**: Version 18.17.0 or higher is recommended
   ```
   nvm install 18
   ```

3. **Environment Variables**: Copy `.env.local.example` to `.env.local` and fill in the required values:
   ```
   cp .env.local.example .env.local
   ```
   
   Required environment variables include:
   - Clerk API keys (for authentication)
   - API URL endpoints
   - Any other service-specific credentials

## Development

To run the application in development mode:

```bash
# Install dependencies
bun install

# Run development server
bun dev
```

The application will be available at http://localhost:3000.

## Production Build

To create a production build:

```bash
# Build the application
bun run build

# Start the production server
bun start
```

## Testing

```bash
# Run all tests
bun test

# Run browser-specific tests
bun test:browser

# Type checking
bun type-check
```

## Deployment

### Vercel (Recommended)

The easiest way to deploy this Next.js application is using Vercel:

1. Push your code to GitHub, GitLab, or Bitbucket
2. Import the project in Vercel dashboard
3. Configure environment variables
4. Deploy

### Docker

You can also deploy using Docker:

1. Build the Docker image:
   ```
   docker build -t crypto-bot-frontend .
   ```

2. Run the container:
   ```
   docker run -p 3000:3000 crypto-bot-frontend
   ```

### Custom Server

For custom server deployments:

1. Build the application: `bun run build`
2. Copy the following to your production server:
   - `.next/` directory
   - `public/` directory
   - `package.json`
   - `next.config.js`
   - `.env.local` (with production values)
3. Install dependencies: `bun install --production`
4. Start the server: `bun start`

## Environment Specific Configurations

### Development
- Uses `.env.local` for local environment variables
- Enables React strict mode
- Includes detailed error messages

### Production
- Should use environment variables set on the hosting platform
- Disables console logs (except errors and warnings)
- Optimizes asset loading

## Troubleshooting

If you encounter issues:

1. Verify environment variables are correctly set
2. Check Node.js and Bun versions
3. Clear `.next/` cache: `rm -rf .next/`
4. Reinstall dependencies: `rm -rf node_modules && bun install` 