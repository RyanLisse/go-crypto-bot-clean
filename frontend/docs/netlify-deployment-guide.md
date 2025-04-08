# Netlify Deployment Guide for Go Crypto Bot Frontend

This guide outlines the steps to deploy the Go Crypto Bot frontend to Netlify for production use.

## Prerequisites

- A Netlify account
- Access to the Go Crypto Bot repository
- Production API endpoints configured and accessible

## Environment Variables

The following environment variables must be set in Netlify:

| Variable | Description | Example |
|----------|-------------|---------|
| `API_URL` | Production API endpoint URL | `https://api.crypto-bot.com/api` |
| `WS_URL` | Production WebSocket endpoint URL | `wss://api.crypto-bot.com/ws` |

## Deployment Steps

### 1. Connect to Git Repository

1. Log in to your Netlify account
2. Click "Add new site" > "Import an existing project"
3. Connect to your Git provider and select the Go Crypto Bot repository
4. Select the branch to deploy (typically `main` or `master`)

### 2. Configure Build Settings

Configure the following build settings:

- **Base directory**: `new_frontend`
- **Build command**: `npm run build`
- **Publish directory**: `dist`

### 3. Set Environment Variables

1. Go to Site settings > Build & deploy > Environment variables
2. Add the required environment variables:
   - `API_URL`: Your production API endpoint
   - `WS_URL`: Your production WebSocket endpoint

### 4. Deploy the Site

1. Click "Deploy site"
2. Wait for the build and deployment to complete
3. Once deployed, Netlify will provide a URL for your site

## Post-Deployment Verification

After deployment, verify the following:

1. The site loads correctly at the provided Netlify URL
2. API connections are working (check network requests)
3. WebSocket connections are established properly
4. All features are functioning as expected

## Troubleshooting

### Common Issues

#### API Connection Errors

If you see API connection errors in the console:
- Verify the `API_URL` environment variable is set correctly
- Check CORS settings on the API server
- Ensure the API is accessible from the Netlify domain

#### WebSocket Connection Failures

If WebSocket connections fail:
- Verify the `WS_URL` environment variable is set correctly
- Check that the WebSocket server accepts connections from the Netlify domain
- Inspect browser console for specific error messages

#### Build Failures

If the build fails:
- Check the Netlify build logs for specific errors
- Ensure all dependencies are properly installed
- Verify that the build command and directory settings are correct

## Custom Domain Setup

To use a custom domain:

1. Go to Site settings > Domain management
2. Click "Add custom domain"
3. Follow the instructions to configure DNS settings
4. Enable HTTPS for your custom domain

## Continuous Deployment

Netlify automatically rebuilds and deploys when changes are pushed to the connected branch. To disable this:

1. Go to Site settings > Build & deploy > Continuous Deployment
2. Toggle "Stop builds" if you want to pause automatic deployments

## Performance Optimization

The current configuration includes:
- Code splitting for vendor libraries
- Disabled source maps in production
- Security headers for better protection
- SPA routing configuration

## Security Considerations

- Environment variables are securely stored in Netlify
- Security headers are configured in `netlify.toml`
- API keys should never be exposed in the frontend code
