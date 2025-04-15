import React from 'react';
import { Sidebar } from './Sidebar';
import { Header } from './Header';
import { ErrorBoundary } from './ErrorBoundary';
import { toast } from 'sonner';

interface MainLayoutProps {
  children?: React.ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {
  const handleGlobalError = (error: Error) => {
    toast.error('Application Error', {
      description: 'An unexpected error occurred. Some features may be unavailable.',
      duration: 5000,
    });
    console.error('Global application error:', error);
  };

  return (
    <ErrorBoundary onError={handleGlobalError}>
      <div className="flex h-screen bg-brutal-background text-brutal-text">
        <Sidebar />
        <div className="flex flex-col flex-1 overflow-hidden">
          <Header />
          <main className="flex-1 overflow-auto">
            {children}
          </main>
        </div>
      </div>
    </ErrorBoundary>
  );
}
