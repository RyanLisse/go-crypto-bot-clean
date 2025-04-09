# Crypto Trading Bot Frontend Documentation

This directory contains comprehensive documentation for the Crypto Trading Bot frontend application, which is built with Vite/React and follows a brutalist design approach.

## Documentation Overview

- [Product Requirements Document](./product-requirements-document.md): Detailed requirements for the frontend application
- [Implementation Guide](./implementation-guide.md): Step-by-step guide for implementing the frontend
- [API Documentation](./api-documentation.md): Comprehensive documentation of the API endpoints and WebSocket interface
- [Data Models](./data-models.ts): TypeScript interfaces for all API data structures

## Current Implementation

The frontend is implemented using Vite and React with React Router, following a brutalist design approach. Key features include:

1. **Brutalist Design System**: Minimalist, high-contrast UI with monospace fonts (JetBrains Mono)
2. **Real-time Data**: Dashboard with portfolio charts and real-time account balance
3. **TanStack Query**: For efficient data fetching and caching
4. **AI Integration**: Google Gemini 1.5 Flash for AI-powered chat assistant
5. **Responsive Layout**: Works on both desktop and mobile devices

## Project Structure

```
new_frontend/
├── src/
│   ├── pages/           # React Router pages
│   │   ├── Dashboard.tsx
│   │   ├── Portfolio.tsx
│   │   ├── Trading.tsx
│   │   └── ...
│   ├── components/      # Reusable React components
│   │   ├── dashboard/   # Dashboard-specific components
│   │   ├── layout/      # Layout components (Header, Sidebar)
│   │   └── ui/          # UI components (buttons, cards, etc.)
│   ├── hooks/           # Custom React hooks
│   ├── lib/             # Utility functions and API clients
│   │   └── gemini.ts    # Google Gemini AI integration
│   └── index.css        # Global styles and theme
├── public/              # Static assets
└── index.html           # Entry HTML file
```

## Technology Stack

- **Framework**: Vite + React
- **Routing**: React Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: React Context and TanStack Query
- **Data Fetching**: TanStack Query
- **Charts**: Recharts
- **AI Integration**: Google Gemini 1.5 Flash
- **UI Components**: Radix UI primitives with shadcn/ui
- **Local Database**: Drizzle ORM with SQLite
- **Package Manager**: Bun

## Development Workflow

1. Start the development server with `bun dev`
2. Make changes to components and pages
3. Preview the application at http://localhost:8080
4. Lint your code with `bun lint`
5. Build for production with `bun build`

## Best Practices

- Use TanStack Query for all API data fetching
- Follow the brutalist design system for consistent UI
- Organize code by feature rather than by type
- Implement responsive design for both desktop and mobile
- Use TypeScript for better type safety and developer experience
- Write unit tests for critical components and services
- Use Sonner for toast notifications instead of the deprecated toast component

## Contributing

If you'd like to contribute to the frontend documentation or example implementation, please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## Support

If you have questions or need help with the frontend implementation, please:

1. Check the existing documentation
2. Look for similar issues in the issue tracker
3. Create a new issue if your question hasn't been addressed

## License

This documentation and example implementation are provided under the same license as the Crypto Trading Bot project.
