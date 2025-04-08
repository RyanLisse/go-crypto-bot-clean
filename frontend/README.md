# Crypto Trading Bot Frontend

This is the frontend application for the Crypto Trading Bot. It provides a user interface for managing trading strategies, monitoring portfolio performance, and configuring the trading bot.

## Features

- Dashboard with key metrics and performance indicators
- Portfolio management and tracking
- Trading interface for executing trades
- Backtesting interface for strategy evaluation
- Settings for configuring the trading bot
- Real-time updates via WebSocket

## Technology Stack

- React 18
- TypeScript
- Material UI for components and styling
- Redux Toolkit for state management
- RTK Query for API calls
- React Router for navigation
- Chart.js for data visualization
- Socket.io for real-time updates

## Project Structure

```
frontend/
├── public/              # Static files
├── src/                 # Source code
│   ├── assets/          # Images, fonts, etc.
│   ├── components/      # Reusable components
│   ├── hooks/           # Custom React hooks
│   ├── pages/           # Page components
│   ├── services/        # API services
│   ├── store/           # Redux store
│   ├── styles/          # Global styles
│   ├── types/           # TypeScript type definitions
│   ├── utils/           # Utility functions
│   ├── App.tsx          # Main App component
│   └── index.tsx        # Entry point
├── .dockerignore        # Docker ignore file
├── .env.example         # Example environment variables
├── Dockerfile           # Docker configuration
├── nginx.conf           # Nginx configuration for production
├── package.json         # Dependencies and scripts
└── tsconfig.json        # TypeScript configuration
```

## Getting Started

### Prerequisites

- Node.js 16+
- npm or yarn

### Installation

1. Clone the repository
2. Navigate to the frontend directory
3. Install dependencies:

```bash
npm install
# or
yarn install
```

4. Create a `.env` file based on `.env.example`
5. Start the development server:

```bash
npm start
# or
yarn start
```

### Building for Production

```bash
npm run build
# or
yarn build
```

### Docker

To build and run the frontend using Docker:

```bash
docker build -t crypto-bot-frontend .
docker run -p 3000:80 crypto-bot-frontend
```

## Development Guidelines

- Follow the established folder structure
- Use TypeScript for type safety
- Create reusable components in the components directory
- Use Material UI components and styling
- Follow the Redux Toolkit pattern for state management
- Use RTK Query for API calls
- Write unit tests for components and utilities

## Available Scripts

- `npm start`: Start the development server
- `npm run build`: Build the application for production
- `npm test`: Run tests
- `npm run lint`: Run ESLint
- `npm run lint:fix`: Fix ESLint issues
- `npm run format`: Format code with Prettier
