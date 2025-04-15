# Migration Guide: Moving from Vite (frontend/) to Next.js (frontend_next/) for Go Crypto Bot Frontend

## Overview

This guide provides a step-by-step process for migrating your Go Crypto Bot Frontend from a Vite-based setup (`frontend/`) to a Next.js-based application (`frontend_next/`), using Bun as the JavaScript runtime and package manager. It covers project setup, file migration, component and routing updates, API integration, environment variables, testing, deployment, and Next.js best practices.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Project Setup](#project-setup)
3. [Migrating Files and Structure](#migrating-files-and-structure)
4. [Component & Layout Migration](#component--layout-migration)
5. [Routing in Next.js](#routing-in-nextjs)
6. [API Integration](#api-integration)
7. [State Management](#state-management)
8. [Environment Variables](#environment-variables)
9. [Testing](#testing)
10. [Deployment](#deployment)
11. [Troubleshooting](#troubleshooting)
12. [Best Practices](#best-practices)
13. [Conclusion](#conclusion)

---

## Prerequisites

- Bun installed (v1.0.0+)
- Node.js (if required by dependencies)
- Git for version control
- Backup of your current `frontend/` project

```bash
# Install Bun if not already installed
curl -fsSL https://bun.sh/install | bash

# Verify Bun installation
bun --version
```

---

## Project Setup

1. **Backup Your Project**

   Always start by backing up your current frontend project.

   ```bash
   cp -r frontend frontend_backup
   ```

2. **Create a New Next.js App with Bun**

   ```bash
   bun create next-app@latest frontend_next --typescript --eslint --tailwind --app --src-dir --import-alias "@/*"
   cd frontend_next
   ```

3. **Install Dependencies**

   Install all required dependencies using Bun:

   ```bash
   # Core dependencies
   bun add react react-dom
   bun add @tanstack/react-query @tanstack/react-query-devtools
   bun add date-fns recharts sonner
   bun add zod @hookform/resolvers react-hook-form
   bun add lucide-react clsx class-variance-authority tailwind-merge

   # Dev dependencies
   bun add -d typescript tailwindcss postcss autoprefixer
   bun add -d eslint @eslint/js
   bun add -d vitest jsdom @testing-library/react
   ```

4. **Initialize UI Component Library (Optional: shadcn/ui)**

   ```bash
   bun add -d shadcn-ui
   bunx shadcn-ui@latest init
   # Follow prompts to match your styling and directory structure
   ```

---

## Migrating Files and Structure

1. **Copy Styling Files**

   - Move your global CSS and any custom styles from `frontend/src/` to `frontend_next/src/app/globals.css` and `frontend_next/src/styles/`.
   - Update `frontend_next/src/app/globals.css` to import your custom styles:

     ```css
     @import '../styles/app.css';

     @tailwind base;
     @tailwind components;
     @tailwind utilities;
     ```

   - Copy your Tailwind config:

     ```bash
     cp frontend/tailwind.config.ts frontend_next/tailwind.config.ts
     ```

2. **Copy Components, Hooks, Lib, and Types**

   ```bash
   cp -r frontend/src/components frontend_next/src/
   cp -r frontend/src/hooks frontend_next/src/
   cp -r frontend/src/lib frontend_next/src/
   cp -r frontend/src/types frontend_next/src/
   ```

3. **Directory Structure Example**

   ```
   frontend_next/src/
   ├── app/
   │   ├── (dashboard)/
   │   │   ├── page.tsx
   │   │   ├── portfolio/
   │   │   ├── trading/
   │   │   ├── new-coins/
   │   │   ├── backtesting/
   │   │   ├── system/
   │   │   ├── config/
   │   │   ├── settings/
   │   │   ├── testing/
   │   │   └── layout.tsx
   │   ├── api/
   │   ├── layout.tsx
   │   └── providers.tsx
   ├── components/
   ├── hooks/
   ├── lib/
   ├── types/
   └── styles/
   ```

---

## Component & Layout Migration

1. **Create Root Layout**

   `frontend_next/src/app/layout.tsx`:

   ```tsx
   import './globals.css';
   import { Providers } from './providers';

   export const metadata = {
     title: 'Go Crypto Bot',
     description: 'Cryptocurrency trading bot for detecting and trading new coin listings',
   };

   export default function RootLayout({ children }: { children: React.ReactNode }) {
     return (
       <html lang="en" suppressHydrationWarning>
         <body>
           <Providers>{children}</Providers>
         </body>
       </html>
     );
   }
   ```

2. **Create Providers Component**

   `frontend_next/src/app/providers.tsx`:

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
           staleTime: 5 * 60 * 1000,
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

3. **Create Dashboard Layout**

   `frontend_next/src/app/(dashboard)/layout.tsx`:

   ```tsx
   'use client';

   import { MainLayout } from '@/components/layout/MainLayout';

   export default function DashboardLayout({ children }: { children: React.ReactNode }) {
     return <MainLayout>{children}</MainLayout>;
   }
   ```

---

## Routing in Next.js

1. **Convert Routes to Pages**

   For each route in your Vite app (`frontend/src/App.tsx`), create a corresponding file in `frontend_next/src/app/(dashboard)/`.

   Example:  
   - `frontend/src/pages/NewCoins.tsx` → `frontend_next/src/app/(dashboard)/new-coins/page.tsx`
   - Repeat for all main routes (portfolio, trading, backtesting, etc.)

2. **Update Navigation Links**

   Replace React Router's `Link` with Next.js `Link`:

   ```tsx
   // Old (Vite/React Router)
   <Link to="/new-coins">New Coins</Link>

   // New (Next.js)
   import Link from 'next/link';
   <Link href="/new-coins">New Coins</Link>
   ```

3. **Remove Vite-Specific Routing Logic**

   Next.js uses file-based routing. Remove any `<Routes>`, `<Route>`, or `react-router` logic.

---

## API Integration

1. **Update API Client**

   - Move your API client to `frontend_next/src/lib/api.ts`.
   - Update environment variable usage:

     ```tsx
     // Next.js
     const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
     ```

2. **(Optional) Create Next.js API Routes**

   Use `frontend_next/src/app/api/` for serverless functions or proxying requests.

   ```tsx
   // Example: src/app/api/status/route.ts
   import { NextResponse } from 'next/server';

   export async function GET() {
     return NextResponse.json({ status: 'healthy' });
   }
   ```

---

## State Management

- React Query, Zustand, and other state libraries work in Next.js as in Vite.
- For hooks/components using browser-only APIs, add `'use client'` at the top of the file.

---

## Environment Variables

1. **Copy and Update Environment Variables**

   ```bash
   cp frontend/.env frontend_next/.env.local
   ```

   - All client-side variables must be prefixed with `NEXT_PUBLIC_` in Next.js.

     ```
     # Old (Vite)
     VITE_API_URL=http://localhost:8080/api/v1

     # New (Next.js)
     NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
     ```

2. **Update References in Code**

   ```tsx
   // Vite
   import.meta.env.VITE_API_URL

   // Next.js
   process.env.NEXT_PUBLIC_API_URL
   ```

---

## Testing

1. **Start the Development Server**

   ```bash
   cd frontend_next
   bun dev
   ```

2. **Test All Features**

   - Compare the Next.js app with the original Vite app.
   - Ensure all components render and function as expected.
   - Test navigation, API integration, forms, and interactivity.

---

## Deployment

1. **Build the Project**

   ```bash
   bun run build
   ```

2. **Preview the Production Build**

   ```bash
   bun run start
   ```

3. **Deploy**

   - Follow your hosting provider's instructions for deploying a Next.js app (Vercel, Netlify, custom server, etc.).

---

## Troubleshooting

- **Styling Issues:**  
  Check for differences in CSS module handling. Adjust imports or module names as needed.

- **Client/Server Rendering Issues:**  
  Add `'use client'` to components using browser APIs.

- **Environment Variables:**  
  Ensure all client-side variables are prefixed with `NEXT_PUBLIC_`.

- **API Integration:**  
  Check for CORS and credentials differences.

- **Routing Issues:**  
  Ensure all navigation uses Next.js `Link` and matches the file-based routing structure.

---

## Best Practices

- Use Next.js file-based routing and layouts for organization.
- Leverage SSR/SSG as needed for performance.
- Keep environment variables secure and properly prefixed.
- Use Next.js API routes for server-side logic or proxying.
- Incrementally adopt Next.js features after migration for further optimization.

---

## Conclusion

By following this guide, you can migrate your Go Crypto Bot Frontend from Vite (`frontend/`) to Next.js (`frontend_next/`), taking advantage of Next.js features while maintaining your app's design and functionality. This approach ensures a smooth transition and sets the foundation for future improvements.
