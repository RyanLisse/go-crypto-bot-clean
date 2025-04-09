# Migration Guide: Implementing Brutalist Design

## Overview

This guide outlines the process of migrating the current frontend at `/Users/neo/Developer/experiments/go-crypto-bot-clean/frontend` to implement the new brutalist design from `/Users/neo/Developer/experiments/go-crypto-bot-clean/crypto-brutal-bot-view`.

## Prerequisites

- Bun installed (for package management)
- Git (for version control)
- Access to both current and new design frontends

## Step-by-Step Migration

### 1. Create Backup

```bash
# Create a backup of the current frontend
cp -r /Users/neo/Developer/experiments/go-crypto-bot-clean/frontend /Users/neo/Developer/experiments/go-crypto-bot-clean/frontend.backup
```

### 2. Update Dependencies

Update `frontend/package.json` with new dependencies from the brutalist design:

```json
{
  "dependencies": {
    "@radix-ui/react-tooltip": "^1.0.7",
    "@tanstack/react-query": "^5.0.0",
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.0.0",
    "lucide-react": "^0.284.0",
    "tailwind-merge": "^1.14.0",
    "tailwindcss-animate": "^1.0.7"
  }
}
```

### 3. Copy New Design Assets

```bash
# Copy the new design's styles and components
cp -r /Users/neo/Developer/experiments/go-crypto-bot-clean/crypto-brutal-bot-view/src/components/ui frontend/src/components/
cp /Users/neo/Developer/experiments/go-crypto-bot-clean/crypto-brutal-bot-view/src/index.css frontend/src/
cp /Users/neo/Developer/experiments/go-crypto-bot-clean/crypto-brutal-bot-view/tailwind.config.ts frontend/
```

### 4. Update Configuration Files

Copy and adapt the following configuration files:

```bash
cp /Users/neo/Developer/experiments/go-crypto-bot-clean/crypto-brutal-bot-view/components.json frontend/
cp /Users/neo/Developer/experiments/go-crypto-bot-clean/crypto-brutal-bot-view/postcss.config.js frontend/
```

### 5. Implement New Layout Components

Copy the core layout components:

```bash
# Copy layout components
cp -r /Users/neo/Developer/experiments/go-crypto-bot-clean/crypto-brutal-bot-view/src/components/layout frontend/src/components/

# Copy pages structure
cp -r /Users/neo/Developer/experiments/go-crypto-bot-clean/crypto-brutal-bot-view/src/pages frontend/src/
```

### 6. Update App Structure

Update `frontend/src/App.tsx`:

```typescript
import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Index from "./pages/Index";
import Dashboard from "./pages/Dashboard";
import NotFound from "./pages/NotFound";
// ... import other pages

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <BrowserRouter>
      <TooltipProvider>
        <Toaster />
        <Sonner />
        <Routes>
          <Route path="/" element={<Index />}>
            <Route index element={<Dashboard />} />
            {/* Add other routes */}
          </Route>
          <Route path="*" element={<NotFound />} />
        </Routes>
      </TooltipProvider>
    </BrowserRouter>
  </QueryClientProvider>
);

export default App;
```

### 7. Implement Brutalist Theme

Update `frontend/src/index.css`:

```css
@import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600;700&display=swap');

@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --background: 0 0% 7%;
    --foreground: 0 0% 97%;
    /* ... copy other CSS variables */
  }

  * {
    @apply border-border;
    box-sizing: border-box;
  }

  html, body {
    @apply font-mono bg-brutal-background text-brutal-text;
    font-feature-settings: "ss01", "ss02", "cv01", "cv02", "cv03";
  }

  /* ... copy other base styles */
}
```

### 8. Update Component Styling

Migrate existing components to use the brutalist design system:

1. Replace Material UI components with brutalist equivalents
2. Update color schemes to use brutal theme colors
3. Implement monospace typography
4. Add high-contrast elements
5. Simplify layouts

### 9. Test and Debug

```bash
# Install dependencies
cd frontend
bun install

# Start development server
bun dev
```

## Verification Checklist

- [ ] All pages render with new brutalist design
- [ ] Typography uses JetBrains Mono
- [ ] High contrast color scheme is applied
- [ ] Components maintain functionality
- [ ] Responsive design works
- [ ] No console errors
- [ ] All features work as before

## Common Issues and Solutions

### 1. Style Conflicts

If you encounter style conflicts:

```bash
# Check for conflicting CSS imports
# Ensure tailwind classes are properly applied
# Verify component class order
```

### 2. Component Integration

For components not matching the new design:

```bash
# Compare with brutalist design components
# Update component props and styling
# Verify accessibility features
```

## Rollback Plan

If issues arise, you can rollback using the backup:

```bash
# Remove modified frontend
rm -rf frontend

# Restore from backup
cp -r /Users/neo/Developer/experiments/go-crypto-bot-clean/frontend.backup frontend
```

## Post-Migration Tasks

1. **Documentation Updates**
   - Update component documentation
   - Document new design system
   - Update README with new styling guidelines

2. **Performance Optimization**
   - Audit bundle size
   - Check for unused styles
   - Optimize images and assets

3. **Accessibility Checks**
   - Verify contrast ratios
   - Test screen reader compatibility
   - Check keyboard navigation

## Conclusion

This migration transforms the existing frontend to implement the new brutalist design while maintaining all functionality. The updated design provides a more distinctive and minimalist user interface with improved contrast and typography.