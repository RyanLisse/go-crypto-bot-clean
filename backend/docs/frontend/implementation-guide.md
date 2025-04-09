# Crypto Trading Bot Frontend Implementation Guide

This guide provides instructions for implementing the Crypto Trading Bot frontend application using Vite and React with React Router. It covers project setup, authentication, API integration, WebSocket communication, AI integration, and best practices for frontend development.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Project Structure](#project-structure)
3. [Authentication](#authentication)
4. [API Integration](#api-integration)
5. [WebSocket Integration](#websocket-integration)
6. [AI Integration](#ai-integration)
7. [Brutalist Design System](#brutalist-design-system)
8. [Error Handling](#error-handling)
9. [Testing](#testing)
10. [Best Practices](#best-practices)

## Getting Started

### Prerequisites

- Node.js (v18 or later)
- Bun (preferred package manager)
- Basic knowledge of TypeScript, React, and React Router
- Access to the Crypto Trading Bot API

### Project Setup

The frontend is built with Vite and React. Here's how to set up the project:

```bash
# Create a new Vite application with React and TypeScript
bun create vite new_frontend -- --template react-ts
cd new_frontend

# Install dependencies
bun install

# Install required packages
bun add @tanstack/react-query @tanstack/react-query-devtools sonner recharts react-router-dom
bun add @radix-ui/react-dialog @radix-ui/react-toast @radix-ui/react-tooltip
bun add @google/generative-ai

# Install dev dependencies
bun add -D tailwindcss postcss autoprefixer
bun x tailwindcss init -p
```

## Project Structure

The project follows this structure:

```
new_frontend/
├── src/
│   ├── pages/           # React Router pages
│   │   ├── Dashboard.tsx  # Dashboard page
│   │   ├── Portfolio.tsx  # Portfolio page
│   │   ├── Trading.tsx    # Trading page
│   │   ├── NewCoins.tsx   # New coins page
│   │   ├── Backtesting.tsx # Backtesting page
│   │   ├── SystemStatus.tsx # System status page
│   │   ├── BotConfig.tsx  # Bot configuration page
│   │   ├── Settings.tsx   # Settings page
│   │   └── NotFound.tsx   # 404 page
│   ├── components/      # Reusable React components
│   │   ├── dashboard/     # Dashboard-specific components
│   │   ├── layout/        # Layout components (Header, Sidebar)
│   │   └── ui/            # UI components (buttons, cards, etc.)
│   ├── hooks/           # Custom React hooks
│   ├── lib/             # Utility functions and API clients
│   │   └── gemini.ts      # Google Gemini AI integration
│   ├── App.tsx          # Main application component with routing
│   ├── main.tsx         # Application entry point
│   └── index.css        # Global styles
├── public/              # Static assets
└── index.html           # HTML entry point
```

### Using TypeScript Data Models

We provide TypeScript interfaces for all API data structures in the [data-models.ts](./data-models.ts) file. Copy this file to your project's `src/types` directory to use these interfaces in your application:

```typescript
import { PortfolioSummary, BoughtCoin, Order } from '../types/data-models';

// Use in component props
interface PortfolioProps {
  portfolio: PortfolioSummary;
  activeTrades: BoughtCoin[];
  isLoading: boolean;
}

// Use in API service functions
async function getPortfolio(): Promise<PortfolioSummary> {
  const response = await apiClient.get('/api/v1/portfolio');
  return response.data;
}
```

## Authentication

The Crypto Trading Bot API uses JWT (JSON Web Token) for authentication. Here's how to implement authentication in your frontend:

### Login Flow

1. Create a login form that collects username and password
2. Send a POST request to `/auth/login` with the credentials
3. Store the received JWT token in localStorage or a secure cookie
4. Include the token in the Authorization header for subsequent API requests

### Example Authentication Service

```javascript
// src/services/authService.js
import axios from 'axios';

const API_URL = 'http://localhost:8080';

export const authService = {
  async login(username, password) {
    const response = await axios.post(`${API_URL}/auth/login`, {
      username,
      password
    });

    if (response.data.token) {
      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data));
    }

    return response.data;
  },

  logout() {
    // Call the logout endpoint
    axios.post(`${API_URL}/auth/logout`).catch(error => console.error('Logout error:', error));

    // Clear local storage
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  getCurrentUser() {
    return JSON.parse(localStorage.getItem('user'));
  },

  isAuthenticated() {
    return !!localStorage.getItem('token');
  },

  getAuthHeader() {
    const token = localStorage.getItem('token');
    return token ? { Authorization: `Bearer ${token}` } : {};
  }
};
```

### Protected Routes

Create a component to handle protected routes:

```jsx
// src/components/ProtectedRoute.jsx
import React from 'react';
import { Navigate } from 'react-router-dom';
import { authService } from '../services/authService';

const ProtectedRoute = ({ children }) => {
  if (!authService.isAuthenticated()) {
    return <Navigate to="/login" />;
  }

  return children;
};

export default ProtectedRoute;
```

Use it in your router:

```jsx
// src/App.jsx
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import ProtectedRoute from './components/ProtectedRoute';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        } />
      </Routes>
    </BrowserRouter>
  );
}
```

## API Integration

### API Client Setup

Create a centralized API client to handle all API requests:

```javascript
// src/services/apiClient.js
import axios from 'axios';
import { authService } from './authService';

const API_URL = 'http://localhost:8080';

// Create axios instance
const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json'
  }
});

// Add request interceptor to add auth token
apiClient.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => Promise.reject(error)
);

// Add response interceptor to handle token expiration
apiClient.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      // Token expired or invalid
      authService.logout();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default apiClient;
```

### Service Modules

Create service modules for different API endpoints:

```javascript
// src/services/portfolioService.js
import apiClient from './apiClient';

export const portfolioService = {
  async getSummary() {
    const response = await apiClient.get('/api/v1/portfolio');
    return response.data;
  },

  async getActiveTrades() {
    const response = await apiClient.get('/api/v1/portfolio/active');
    return response.data;
  },

  async getPerformance(timeRange = '30d') {
    const response = await apiClient.get(`/api/v1/portfolio/performance?time_range=${timeRange}`);
    return response.data;
  },

  async getTotalValue() {
    const response = await apiClient.get('/api/v1/portfolio/value');
    return response.data;
  }
};

// src/services/tradeService.js
import apiClient from './apiClient';

export const tradeService = {
  async getTradeHistory() {
    const response = await apiClient.get('/api/v1/trade/history');
    return response.data;
  },

  async executeTrade(symbol, quantity, orderType = 'MARKET') {
    const response = await apiClient.post('/api/v1/trade/buy', {
      symbol,
      quantity,
      order_type: orderType
    });
    return response.data;
  },

  async sellCoin(symbol, quantity) {
    const response = await apiClient.post('/api/v1/trade/sell', {
      symbol,
      quantity
    });
    return response.data;
  },

  async getTradeStatus(id) {
    const response = await apiClient.get(`/api/v1/trade/status/${id}`);
    return response.data;
  }
};

// Add more service modules for other API endpoints
```

## WebSocket Integration

The Crypto Trading Bot API provides real-time updates via WebSocket. Here's how to integrate it:

### WebSocket Service

```javascript
// src/services/websocketService.js
export class WebSocketService {
  constructor(url) {
    this.url = url;
    this.socket = null;
    this.listeners = {
      market_data: [],
      trade_notification: [],
      new_coin_alert: [],
      error: []
    };
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 3000; // 3 seconds
  }

  connect() {
    if (this.socket) {
      this.disconnect();
    }

    this.socket = new WebSocket(this.url);

    this.socket.onopen = () => {
      console.log('WebSocket connected');
      this.reconnectAttempts = 0;

      // Subscribe to ticker updates
      this.send({
        type: 'subscribe_ticker',
        payload: {
          symbols: ['BTCUSDT', 'ETHUSDT']
        }
      });
    };

    this.socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        const { type } = data;

        if (this.listeners[type]) {
          this.listeners[type].forEach(callback => callback(data));
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    this.socket.onclose = () => {
      console.log('WebSocket disconnected');

      // Attempt to reconnect
      if (this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++;
        setTimeout(() => this.connect(), this.reconnectDelay);
      }
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }

  disconnect() {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }

  send(data) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(data));
    } else {
      console.error('WebSocket not connected');
    }
  }

  subscribe(type, callback) {
    if (this.listeners[type]) {
      this.listeners[type].push(callback);
    }

    return () => {
      if (this.listeners[type]) {
        this.listeners[type] = this.listeners[type].filter(cb => cb !== callback);
      }
    };
  }
}

// Create and export a singleton instance
export const websocketService = new WebSocketService('ws://localhost:8080/ws');
```

### Using WebSocket in Components

```jsx
// src/components/MarketData.jsx
import React, { useEffect, useState } from 'react';
import { websocketService } from '../services/websocketService';

const MarketData = () => {
  const [marketData, setMarketData] = useState({});

  useEffect(() => {
    // Connect to WebSocket when component mounts
    websocketService.connect();

    // Subscribe to market data updates
    const unsubscribe = websocketService.subscribe('market_data', (data) => {
      const { payload } = data;
      setMarketData(prevData => ({
        ...prevData,
        [payload.symbol]: payload
      }));
    });

    // Cleanup when component unmounts
    return () => {
      unsubscribe();
    };
  }, []);

  return (
    <div>
      <h2>Market Data</h2>
      <div className="market-data-grid">
        {Object.values(marketData).map(data => (
          <div key={data.symbol} className="market-data-card">
            <h3>{data.symbol}</h3>
            <p>Price: ${data.price.toFixed(2)}</p>
            <p>Volume: {data.volume.toFixed(2)}</p>
            <p>Updated: {new Date(data.timestamp * 1000).toLocaleTimeString()}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default MarketData;
```

## AI Integration

The frontend integrates with Google's Gemini 1.5 Flash model to provide an AI-powered chat assistant. This allows users to ask questions about their portfolio, trading strategies, and market data.

### Setting Up Gemini Integration

```typescript
// src/lib/gemini.ts
import { GoogleGenerativeAI, HarmCategory, HarmBlockThreshold } from '@google/generative-ai';

// Initialize the Google Generative AI with your API key
const genAI = new GoogleGenerativeAI(process.env.NEXT_PUBLIC_GEMINI_API_KEY || '');

// System prompt to guide the AI's behavior
const getSystemPrompt = () => {
  return `You are an AI assistant for a cryptocurrency trading bot.
  You help users understand their portfolio, analyze market trends, and provide insights on trading strategies.

  Current portfolio data and trading context will be provided to you in each conversation.
  Base your responses on this data when answering questions about the user's portfolio or trades.

  Keep responses concise, technical, and focused on cryptocurrency trading.`;
};

// Get the Gemini model
export const getGeminiModel = () => {
  return genAI.getGenerativeModel({
    model: 'gemini-1.5-flash',
    safetySettings: [
      {
        category: HarmCategory.HARM_CATEGORY_HARASSMENT,
        threshold: HarmBlockThreshold.BLOCK_MEDIUM_AND_ABOVE,
      },
      {
        category: HarmCategory.HARM_CATEGORY_HATE_SPEECH,
        threshold: HarmBlockThreshold.BLOCK_MEDIUM_AND_ABOVE,
      },
    ],
  });
};

// Generate a response from Gemini
export const generateGeminiResponse = async (
  prompt: string,
  history: { role: string; content: string }[]
) => {
  try {
    const model = getGeminiModel();

    // Format conversation history for Gemini
    let formattedPrompt = getSystemPrompt() + '\n\nConversation history:\n';

    // Add conversation history to the prompt
    for (const msg of history) {
      const role = msg.role === 'user' ? 'User' : 'Assistant';
      formattedPrompt += `${role}: ${msg.content}\n`;
    }

    // Add the current user question
    formattedPrompt += `\nUser: ${prompt}\n\nAssistant: `;

    // Generate content
    const result = await model.generateContent(formattedPrompt);
    const response = result.response;
    return response.text();
  } catch (error) {
    console.error('Error generating Gemini response:', error);
    return 'Sorry, I encountered an error. Please try again later.';
  }
};
```

### Creating the Chat Component

```tsx
// src/components/chat/ChatInterface.tsx
import { useState } from 'react';
import { generateGeminiResponse } from '@/lib/gemini';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useToast } from '@/hooks/use-toast';

type Message = {
  role: 'user' | 'assistant';
  content: string;
};

export function ChatInterface() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim()) return;

    // Add user message
    const userMessage = { role: 'user' as const, content: input };
    setMessages((prev) => [...prev, userMessage]);
    setInput('');
    setIsLoading(true);

    try {
      // Get AI response
      const aiResponse = await generateGeminiResponse(
        input,
        messages
      );

      // Add AI message
      setMessages((prev) => [
        ...prev,
        { role: 'assistant', content: aiResponse },
      ]);
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to get response from AI',
        variant: 'destructive',
      });
      console.error('Chat error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="brutal-card h-full flex flex-col">
      <div className="brutal-card-header">AI ASSISTANT</div>
      <div className="flex-1 overflow-auto p-4 space-y-4">
        {messages.length === 0 ? (
          <div className="text-brutal-text/50 text-center py-8">
            Ask me anything about your portfolio or trading strategies
          </div>
        ) : (
          messages.map((msg, i) => (
            <div
              key={i}
              className={`${msg.role === 'user' ? 'text-right' : 'text-left'}`}
            >
              <div
                className={`inline-block p-3 rounded ${msg.role === 'user' ? 'bg-brutal-info text-white' : 'bg-brutal-panel border border-brutal-border'}`}
              >
                {msg.content}
              </div>
            </div>
          ))
        )}
      </div>
      <form onSubmit={handleSubmit} className="p-4 border-t border-brutal-border">
        <div className="flex gap-2">
          <Input
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Type your message..."
            disabled={isLoading}
            className="flex-1"
          />
          <Button type="submit" disabled={isLoading || !input.trim()}>
            {isLoading ? 'Thinking...' : 'Send'}
          </Button>
        </div>
      </form>
    </div>
  );
}
```

## Brutalist Design System

The frontend follows a brutalist design approach with a focus on minimalism, high contrast, and monospace typography.

### Setting Up the Theme

```css
/* src/styles/globals.css */
@import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600;700&display=swap');

@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --background: 0 0% 7%;
    --foreground: 0 0% 97%;

    --card: 0 0% 12%;
    --card-foreground: 0 0% 97%;

    --popover: 0 0% 12%;
    --popover-foreground: 0 0% 97%;

    --primary: 210 100% 47%;
    --primary-foreground: 0 0% 100%;

    --secondary: 240 4.8% 15.9%;
    --secondary-foreground: 0 0% 97%;

    --muted: 0 0% 15%;
    --muted-foreground: 0 0% 70%;

    --accent: 0 0% 15%;
    --accent-foreground: 0 0% 97%;

    --destructive: 0 84.2% 60.2%;
    --destructive-foreground: 0 0% 97%;

    --border: 0 0% 20%;
    --input: 0 0% 15%;
    --ring: 0 0% 20%;

    --radius: 0;
  }

  * {
    @apply border-border;
    box-sizing: border-box;
  }

  html, body {
    @apply font-mono bg-brutal-background text-brutal-text;
    font-feature-settings: "ss01", "ss02", "cv01", "cv02", "cv03";
  }

  /* Brutalist card style */
  .brutal-card {
    @apply bg-brutal-panel border border-brutal-border p-4 flex flex-col;
  }

  .brutal-card-header {
    @apply text-xs uppercase tracking-widest mb-2 text-brutal-text/70;
  }
}
```

### Tailwind Configuration

```typescript
// tailwind.config.ts
import type { Config } from "tailwindcss";

export default {
  darkMode: ["class"],
  content: [
    "./pages/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./app/**/*.{ts,tsx}",
    "./src/**/*.{ts,tsx}",
  ],
  prefix: "",
  theme: {
    container: {
      center: true,
      padding: '2rem',
      screens: {
        'sm': '640px',
        'md': '768px',
        'lg': '1024px',
        'xl': '1280px',
        '2xl': '1400px'
      }
    },
    fontFamily: {
      'mono': ['JetBrains Mono', 'monospace'],
      'sans': ['JetBrains Mono', 'Inter', 'sans-serif']
    },
    extend: {
      colors: {
        // Brutalist theme specific colors
        brutal: {
          background: '#121212',
          panel: '#1e1e1e',
          border: '#333333',
          text: '#f7f7f7',
          error: '#ff4d4d',
          info: '#3a86ff',
          success: '#00b894',
          warning: '#fdcb6e',
        }
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
} satisfies Config;
```

## Error Handling

Implement consistent error handling throughout your application using Sonner for toast notifications:

```typescript
// src/lib/error-handler.ts
import { toast } from 'sonner';

export const handleApiError = (error: any) => {
  let errorMessage = 'An unexpected error occurred';

  if (error.response) {
    // The request was made and the server responded with an error status
    const { data } = error.response;
    errorMessage = data.message || `Error: ${error.response.status}`;
  } else if (error.request) {
    // The request was made but no response was received
    errorMessage = 'No response from server. Please check your connection.';
  } else {
    // Something happened in setting up the request
    errorMessage = error.message;
  }

  // Display error to user
  toast.error(errorMessage);

  // Log error for debugging
  console.error('API Error:', error);

  return errorMessage;
};
```

Use it with TanStack Query:

```typescript
// src/hooks/queries.ts
import { useQuery } from '@tanstack/react-query';
import { handleApiError } from '@/lib/error-handler';
import { portfolioService } from '@/services/portfolio-service';

export const usePortfolioValueQuery = () => {
  return useQuery({
    queryKey: ['portfolioValue'],
    queryFn: async () => {
      try {
        return await portfolioService.getTotalValue();
      } catch (error) {
        handleApiError(error);
        throw error;
      }
    },
    refetchInterval: 60000, // Refetch every minute
  });
};
```

## Testing

The frontend can use Vitest for unit testing and Playwright for end-to-end testing, which integrate well with Vite.

### Setting Up Vitest

First, install Vitest and testing libraries:

```bash
bun add -D vitest happy-dom @testing-library/react @testing-library/jest-dom
```

Then create a Vitest configuration file:

```typescript
// vitest.config.ts
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react-swc'
import { resolve } from 'path'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'happy-dom',
    globals: true,
    setupFiles: ['./tests/setup.ts'],
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
    },
  },
})
```

### Writing Component Tests

```typescript
// src/components/dashboard/StatsCard.test.tsx
import { render, screen } from '@testing-library/react'
import { StatsCard } from './StatsCard'

describe('StatsCard', () => {
  it('renders the title and value correctly', () => {
    render(
      <StatsCard
        title="Total Balance"
        value="$1,234.56"
        change={5.2}
        isLoading={false}
      />
    )

    expect(screen.getByText('TOTAL BALANCE')).toBeInTheDocument()
    expect(screen.getByText('$1,234.56')).toBeInTheDocument()
    expect(screen.getByText('+5.2%')).toBeInTheDocument()
  })

  it('shows loading state when isLoading is true', () => {
    render(
      <StatsCard
        title="Total Balance"
        value="$0.00"
        change={0}
        isLoading={true}
      />
    )

    expect(screen.getByTestId('loading-skeleton')).toBeInTheDocument()
  })
})
```

### Setting Up Playwright

Install Playwright for end-to-end testing:

```bash
bun x playwright install
```

Then configure it to work with your Vite application:

```typescript
// playwright.config.ts
import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:8080',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 12'] },
    },
  ],
  webServer: {
    command: 'bun dev',
    port: 8080,
    reuseExistingServer: !process.env.CI,
  },
})
```

### Writing E2E Tests

```typescript
// tests/e2e/dashboard.spec.ts
import { test, expect } from '@playwright/test'

test('dashboard page loads correctly', async ({ page }) => {
  await page.goto('/')

  // Check that the page title is correct
  await expect(page).toHaveTitle(/Crypto Trading Bot/)

  // Check that the main components are visible
  await expect(page.getByText('DASHBOARD')).toBeVisible()
  await expect(page.getByText('OVERVIEW')).toBeVisible()

  // Check that the stats cards are visible
  await expect(page.getByText('TOTAL BALANCE')).toBeVisible()
  await expect(page.getByText('24H CHANGE')).toBeVisible()

  // Check that the chart is visible
  await expect(page.getByTestId('performance-chart')).toBeVisible()
})
```

Add the test script to your package.json:

```json
{
  "scripts": {
    "test": "bun test",
    "test:watch": "bun test --watch",
    "test:e2e": "bunx playwright test",
    "test:e2e:ui": "bunx playwright test --ui"
  }
}
```

## Local Data Persistence

### Using Drizzle ORM with SQLite

For offline capabilities and local data persistence, you can use Drizzle ORM with SQLite. This is particularly useful for storing user preferences, caching API responses, and enabling offline functionality.

#### Setting Up Drizzle

First, install Drizzle and its dependencies:

```bash
bun add drizzle-orm @libsql/client
bun add -D drizzle-kit
```

#### Creating Schema

```typescript
// src/lib/db/schema.ts
import { sqliteTable, text, integer, real } from 'drizzle-orm/sqlite-core';

export const userPreferences = sqliteTable('user_preferences', {
  id: text('id').primaryKey(),
  theme: text('theme').default('dark'),
  currency: text('currency').default('USD'),
  language: text('language').default('en'),
  lastUpdated: integer('last_updated'),
});

export const cachedPortfolio = sqliteTable('cached_portfolio', {
  id: text('id').primaryKey(),
  totalValue: real('total_value'),
  dailyChange: real('daily_change'),
  weeklyChange: real('weekly_change'),
  data: text('data', { mode: 'json' }),
  timestamp: integer('timestamp'),
});
```

#### Setting Up the Database

```typescript
// src/lib/db/index.ts
import { createClient } from '@libsql/client';
import { drizzle } from 'drizzle-orm/libsql';
import * as schema from './schema';

// Create a local SQLite database
const client = createClient({
  url: 'file:local.db',
});

export const db = drizzle(client, { schema });
```

#### Using the Database

```typescript
// src/hooks/useLocalPortfolio.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { db } from '@/lib/db';
import { cachedPortfolio } from '@/lib/db/schema';
import { eq } from 'drizzle-orm';

export function useLocalPortfolio() {
  const queryClient = useQueryClient();

  const { data: portfolio, isLoading } = useQuery({
    queryKey: ['localPortfolio'],
    queryFn: async () => {
      const result = await db.select().from(cachedPortfolio).where(eq(cachedPortfolio.id, 'main'));
      return result[0] || null;
    },
  });

  const updatePortfolio = useMutation({
    mutationFn: async (newData) => {
      await db
        .insert(cachedPortfolio)
        .values({
          id: 'main',
          ...newData,
          timestamp: Date.now(),
        })
        .onConflictDoUpdate({
          target: cachedPortfolio.id,
          set: {
            ...newData,
            timestamp: Date.now(),
          },
        });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['localPortfolio'] });
    },
  });

  return { portfolio, isLoading, updatePortfolio };
}
```

### When to Use Drizzle on the Frontend

Drizzle ORM with SQLite is beneficial in the following scenarios:

1. **Offline Support**: Store critical data locally to enable app functionality without an internet connection
2. **Performance Optimization**: Cache API responses to reduce network requests
3. **User Preferences**: Store user settings and preferences locally
4. **Data Synchronization**: Implement a sync mechanism between local and remote data
5. **Reduced API Calls**: Store non-critical or slowly changing data locally

However, for most frontend applications, you may not need a full ORM. Consider these alternatives:

- **localStorage/sessionStorage**: For simple key-value storage needs
- **IndexedDB**: For more complex data without needing an ORM
- **TanStack Query**: For caching API responses with its built-in caching mechanism

## Best Practices

### State Management

For larger applications, consider using a state management library like Redux or Zustand:

```javascript
// Example using Zustand
import create from 'zustand';
import { portfolioService } from '../services/portfolioService';
import { handleApiError } from '../utils/errorHandler';

export const usePortfolioStore = create((set) => ({
  portfolio: null,
  activeTrades: [],
  loading: false,
  error: null,

  fetchPortfolio: async () => {
    set({ loading: true, error: null });
    try {
      const portfolio = await portfolioService.getSummary();
      set({ portfolio, loading: false });
    } catch (error) {
      const errorMessage = handleApiError(error);
      set({ error: errorMessage, loading: false });
    }
  },

  fetchActiveTrades: async () => {
    set({ loading: true, error: null });
    try {
      const activeTrades = await portfolioService.getActiveTrades();
      set({ activeTrades, loading: false });
    } catch (error) {
      const errorMessage = handleApiError(error);
      set({ error: errorMessage, loading: false });
    }
  }
}));
```

### Code Organization

Organize your code by feature rather than by type:

```
src/
├── features/
│   ├── auth/
│   │   ├── components/
│   │   ├── services/
│   │   └── pages/
│   ├── portfolio/
│   │   ├── components/
│   │   ├── services/
│   │   └── pages/
│   ├── trading/
│   │   ├── components/
│   │   ├── services/
│   │   └── pages/
│   └── newcoins/
│       ├── components/
│       ├── services/
│       └── pages/
├── shared/
│   ├── components/
│   ├── services/
│   └── utils/
└── App.jsx
```

### Responsive Design

Ensure your application works well on both desktop and mobile devices:

```jsx
// Example responsive component
import { useMediaQuery } from '@mui/material';

const Dashboard = () => {
  const isMobile = useMediaQuery('(max-width:600px)');

  return (
    <div className={`dashboard ${isMobile ? 'mobile' : 'desktop'}`}>
      {isMobile ? (
        <MobileDashboard />
      ) : (
        <DesktopDashboard />
      )}
    </div>
  );
};
```

## Example Implementation

A complete example implementation is available in the `web/auth_example.html` file in the repository. This example demonstrates:

1. Authentication with JWT
2. Making authenticated API requests
3. Handling errors
4. Basic UI for testing

To use the example:

1. Start the Crypto Trading Bot API server
2. Open the `web/auth_example.html` file in a browser
3. Log in with the credentials defined in the AuthHandler (admin/admin123 or user/user123)
4. Test making API requests to different endpoints

For a more comprehensive implementation, refer to the [Product Requirements Document](./product-requirements-document.md) which outlines the full scope of the frontend application.

## Next Steps

1. Set up a proper frontend project using a framework like React, Vue, or Angular
2. Implement the authentication flow
3. Create components for each feature (portfolio, trading, new coins, etc.)
4. Integrate WebSocket for real-time updates
5. Add charts and visualizations for data
6. Implement responsive design for mobile devices
7. Add comprehensive error handling and loading states

By following this guide, you should be able to create a fully functional frontend for the Crypto Trading Bot that provides a great user experience and leverages all the features of the API.
