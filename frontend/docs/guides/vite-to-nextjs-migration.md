# Migration Guide: Vite to Next.js with Bun for Go Crypto Bot Frontend

## Overview

This document outlines the step-by-step process for migrating the Go Crypto Bot Frontend from Vite to Next.js while maintaining the exact same design and UI components. We'll use Bun as our JavaScript runtime and package manager throughout the migration process.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Project Setup](#project-setup)
3. [File Structure Migration](#file-structure-migration)
4. [Component Migration](#component-migration)
5. [Routing Migration](#routing-migration)
6. [API Integration](#api-integration)
7. [State Management Migration](#state-management-migration)
8. [Environment Variables](#environment-variables)
9. [Testing the Migration](#testing-the-migration)
10. [Deployment](#deployment)
11. [Troubleshooting](#troubleshooting)

## Prerequisites

Before starting the migration, ensure you have:

- Bun installed (version 1.0.0 or higher)
- Git for version control
- A backup of your current Vite project

```bash
# Install Bun if not already installed
curl -fsSL https://bun.sh/install | bash

# Verify Bun installation
bun --version
```

## Project Setup

### 1. Create a Backup

```bash
cp -r /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend_backup
```

### 2. Create a New Next.js Project with Bun

```bash
cd /Users/neo/Developer/experiments/go-crypto-bot-migration/
bun create next-app@latest new_frontend_next --typescript --eslint --tailwind --app --src-dir --import-alias "@/*"
cd new_frontend_next
```

### 3. Install Existing Dependencies with Bun

```bash
# Install core dependencies
bun add react@18.3.1 react-dom@18.3.1
bun add @tanstack/react-query@5.72.0 @tanstack/react-query-devtools@5.72.0
bun add date-fns@4.1.0 recharts@2.15.2 sonner@1.7.4
bun add zod@3.24.2 @hookform/resolvers@3.10.0 react-hook-form@7.55.0
bun add lucide-react@0.462.0 clsx@2.1.1 class-variance-authority@0.7.1 tailwind-merge@2.6.0

# Install development dependencies
bun add -d typescript@5.8.3 tailwindcss@3.4.17 postcss@8.5.3 autoprefixer@10.4.21
bun add -d eslint@9.24.0 @eslint/js@9.24.0
bun add -d vitest@1.3.1 jsdom@24.0.0 @testing-library/react@14.2.1
```

### 4. Install Shadcn UI Components

```bash
# Install shadcn-ui
bun add -d shadcn-ui

# Initialize shadcn-ui with the exact same configuration as current project
bunx shadcn-ui@latest init
```

When prompted, select the following options to match your current styling:
- Style: Default (or match your current theme)
- Base color: Slate (or match your current base color)
- CSS variables: Yes
- Directory: src/components/ui
- Import alias: @/components/ui

```bash
# Install the specific components you're currently using
bunx shadcn-ui@latest add button card dialog toast input checkbox
bunx shadcn-ui@latest add popover calendar date-picker
bunx shadcn-ui@latest add accordion alert-dialog aspect-ratio avatar
bunx shadcn-ui@latest add collapsible context-menu dropdown-menu
bunx shadcn-ui@latest add hover-card label menubar navigation-menu
bunx shadcn-ui@latest add progress radio-group scroll-area select
bunx shadcn-ui@latest add separator slider slot switch tabs
bunx shadcn-ui@latest add toggle toggle-group tooltip
```

## File Structure Migration

### 1. Copy Styling Files First

To ensure the design remains exactly the same, copy your CSS files:

```bash
# Copy CSS files
cp /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend/src/index.css /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend_next/src/app/globals.css
cp /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend/src/App.css /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend_next/src/styles/app.css

# Copy Tailwind config
cp /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend/tailwind.config.ts /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend_next/tailwind.config.ts
```

Update the `src/app/globals.css` file to import your app.css:

```css
@import '../styles/app.css';

@tailwind base;
@tailwind components;
@tailwind utilities;
```

### 2. Setup App Directory Structure

Create the following directory structure:

```
src/
├── app/
│   ├── (dashboard)/         # Group for authenticated routes
│   │   ├── page.tsx         # Home/dashboard page
│   │   ├── portfolio/       # Portfolio page
│   │   ├── trading/         # Trading page
│   │   ├── new-coins/       # New coins page
│   │   ├── backtesting/     # Backtesting page
│   │   ├── system/          # System status page
│   │   ├── config/          # Bot config page
│   │   ├── settings/        # Settings page
│   │   ├── testing/         # Testing page
│   │   └── layout.tsx       # Dashboard layout
│   ├── api/                 # API routes
│   ├── layout.tsx           # Root layout
│   └── providers.tsx        # App providers
├── components/              # Copy from current project
├── hooks/                   # Copy from current project
├── lib/                     # Copy from current project
├── types/                   # Copy from current project
└── styles/                  # Additional styles
```

## Component Migration

### 1. Copy Components Directory

```bash
cp -r /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend/src/components/* /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend_next/src/components/
```

### 2. Create Root Layout

Create a root layout in `src/app/layout.tsx`:

```tsx
import './globals.css';
import { Providers } from './providers';

export const metadata = {
  title: 'Go Crypto Bot',
  description: 'Cryptocurrency trading bot for detecting and trading new coin listings',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
```

### 3. Create Providers Component

Create a providers component in `src/app/providers.tsx`:

```tsx
'use client';

import { useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { Toaster } from '@/components/ui/toaster';
import { Toaster as Sonner } from '@/components/ui/sonner';
import { TooltipProvider } from '@/components/ui/tooltip';

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(() => new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
        retry: 1,
        staleTime: 5 * 60 * 1000, // 5 minutes
      },
    },
  }));

  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <Toaster />
        <Sonner />
        {children}
      </TooltipProvider>
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}
```

### 4. Create Dashboard Layout

Create a dashboard layout in `src/app/(dashboard)/layout.tsx`:

```tsx
'use client';

import { MainLayout } from '@/components/layout/MainLayout';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <MainLayout>{children}</MainLayout>;
}
```

## Routing Migration

Next.js uses file-based routing, so we need to convert our React Router routes to Next.js pages.

### 1. Copy and Adjust the MainLayout Component

Create or modify `src/components/layout/MainLayout.tsx` to work with Next.js:

```tsx
'use client';

import { ReactNode } from 'react';
// Import your sidebar, navigation components, etc.

export function MainLayout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen flex">
      {/* Your sidebar component */}
      <div className="flex-1 flex flex-col">
        {/* Your main navigation component */}
        <main className="flex-1 overflow-auto">
          {children}
        </main>
      </div>
    </div>
  );
}
```

### 2. Convert Each Route to a Page

For each route in your Vite project's `App.tsx`, create a corresponding page in the Next.js app directory.

Example for the New Coins page (`src/app/(dashboard)/new-coins/page.tsx`):

```tsx
'use client';

// Copy and adapt your existing NewCoins.tsx component
// Make sure to add 'use client' at the top

export default function NewCoinsPage() {
  // Copy the content of your existing NewCoins component
  // Adjust imports to use the new paths
}
```

Repeat this for all your routes:
- `src/app/(dashboard)/page.tsx` (from Index.tsx)
- `src/app/(dashboard)/portfolio/page.tsx` (from Portfolio.tsx)
- `src/app/(dashboard)/trading/page.tsx` (from Trading.tsx)
- `src/app/(dashboard)/backtesting/page.tsx` (from Backtesting.tsx)
- `src/app/(dashboard)/system/page.tsx` (from SystemStatus.tsx)
- `src/app/(dashboard)/config/page.tsx` (from BotConfig.tsx)
- `src/app/(dashboard)/settings/page.tsx` (from Settings.tsx)
- `src/app/(dashboard)/testing/page.tsx` (from Testing.tsx)

### 3. Update Navigation Links

Update any navigation components to use Next.js `Link` component instead of React Router's `Link`:

```tsx
import Link from 'next/link';

// Replace
<Link to="/new-coins">New Coins</Link>

// With
<Link href="/new-coins">New Coins</Link>
```

## API Integration

### 1. Copy API Client Code

```bash
cp -r /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend/src/lib/api.ts /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend_next/src/lib/
```

### 2. Adapt API Client for Next.js

Update your API client to use Next.js environment variables:

```tsx
// src/lib/api.ts
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Rest of your API client remains the same
```

### 3. Copy Custom Hooks

```bash
cp -r /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend/src/hooks/* /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend_next/src/hooks/
```

### 4. Create API Routes (Optional)

You can create Next.js API routes for proxying requests or handling server-side logic:

```tsx
// src/app/api/status/route.ts
import { NextResponse } from 'next/server';

export async function GET() {
  // Fetch from backend or return mock data
  return NextResponse.json({ status: 'healthy' });
}
```

## State Management Migration

### 1. React Query Setup

React Query should work mostly the same between Vite and Next.js, but ensure you're using the `client` directive:

```tsx
'use client';

import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';

export function useNewCoinsQuery() {
  return useQuery({
    queryKey: ['newCoins'],
    queryFn: () => api.getNewCoins(),
  });
}
```

### 2. Local State

Local component state using `useState` and `useReducer` will work the same in Next.js as in your Vite project.

## Environment Variables

### 1. Copy Environment Variables

```bash
cp /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend/.env /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend_next/.env.local
```

### 2. Update Environment Variable Names

In Next.js, client-side environment variables must be prefixed with `NEXT_PUBLIC_`:

```
# Original .env
VITE_ENABLE_RATE_LIMITING=true

# New .env.local
NEXT_PUBLIC_ENABLE_RATE_LIMITING=true
```

Update all environment variable references in your code:

```tsx
// Vite version
const enableRateLimiting = import.meta.env.VITE_ENABLE_RATE_LIMITING === 'true';

// Next.js version
const enableRateLimiting = process.env.NEXT_PUBLIC_ENABLE_RATE_LIMITING === 'true';
```

## Testing the Migration

### 1. Start the Development Server

```bash
cd /Users/neo/Developer/experiments/go-crypto-bot-migration/new_frontend_next
bun dev
```

### 2. Visual Comparison

- Compare the Next.js app side-by-side with the Vite app
- Ensure all components render identically
- Check that all functionality works as expected

### 3. Test All Features

- Test date-based filtering in the New Coins page
- Verify all API integrations work correctly
- Test navigation between pages
- Ensure forms and interactive elements work properly

## Deployment

### 1. Build the Project

```bash
bun run build
```

### 2. Preview the Production Build

```bash
bun run start
```

### 3. Deploy to Your Hosting Service

Follow your hosting provider's instructions for deploying a Next.js app.

## Troubleshooting

### Common Issues and Solutions

1. **Styling Differences**

   **Problem**: Styles don't match exactly between Vite and Next.js.
   
   **Solution**: Check for CSS module differences. Next.js has slightly different CSS module handling. You may need to adjust your CSS imports or module names.

2. **Client-Side Rendering Issues**

   **Problem**: Components that rely on browser APIs fail during server rendering.
   
   **Solution**: Add the `'use client'` directive to components that use browser-only APIs or use dynamic imports with `next/dynamic`.

3. **Environment Variable Access**

   **Problem**: Environment variables not accessible in the browser.
   
   **Solution**: Ensure all client-side environment variables are prefixed with `NEXT_PUBLIC_`.

4. **API Integration Issues**

   **Problem**: API calls fail or return unexpected results.
   
   **Solution**: Check for differences in how fetch requests are handled, especially with CORS and credentials.

5. **Routing Differences**

   **Problem**: Links or navigation don't work as expected.
   
   **Solution**: Ensure all React Router `Link` components are replaced with Next.js `Link` components, and route paths match your file structure.

## Conclusion

This migration guide provides a structured approach to moving your Go Crypto Bot Frontend from Vite to Next.js while preserving the exact same design and functionality. By following these steps, you'll benefit from Next.js features like improved performance, server-side rendering, and API routes while maintaining the user experience of your current application.

The migration approach prioritizes preserving the current design and functionality first, then incrementally leveraging Next.js features as needed. This ensures a smooth transition with minimal disruption to users.
