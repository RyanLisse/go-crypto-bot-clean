import React, { useEffect } from 'react';
import { AlertTriangle, RefreshCw, Info } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';

interface DashboardErrorFallbackProps {
  error?: Error;
  resetErrorBoundary?: () => void;
}

/**
 * A specialized error fallback component for the Dashboard
 * Displays a user-friendly error message with a retry button
 */
export function DashboardErrorFallback({
  error,
  resetErrorBoundary
}: DashboardErrorFallbackProps) {
  // Log error details to console
  useEffect(() => {
    console.error('Dashboard error details:', {
      message: error?.message,
      stack: error?.stack,
      name: error?.name,
    });

    // Try to fetch status directly to debug
    fetch('http://localhost:8080/api/v1/status')
      .then(response => {
        console.log('Direct status fetch response:', response);
        return response.json();
      })
      .then(data => console.log('Direct status fetch data:', data))
      .catch(err => console.error('Direct status fetch error:', err));
  }, [error]);

  return (
    <div className="grid gap-6 p-6">
      <Card className="p-6 border-destructive/30 bg-destructive/5">
        <div className="flex flex-col items-center text-center">
          <AlertTriangle className="w-12 h-12 text-destructive mb-4" />
          <h2 className="text-xl font-bold mb-2">Dashboard Error</h2>
          <p className="text-brutal-text/70 mb-6 max-w-md">
            {error?.message || 'There was a problem loading the dashboard data.'}
          </p>
          <p className="text-brutal-text/50 mb-6 text-sm max-w-md">
            The backend API may be unavailable or experiencing issues.
            Fallback data is being displayed where possible.
          </p>
          <Button
            onClick={resetErrorBoundary}
            variant="outline"
            className="flex items-center gap-2 mb-4"
          >
            <RefreshCw className="w-4 h-4" />
            Retry
          </Button>

          {error && (
            <Collapsible className="w-full max-w-md">
              <CollapsibleTrigger asChild>
                <Button variant="ghost" size="sm" className="flex items-center gap-1">
                  <Info className="w-4 h-4" />
                  Show Error Details
                </Button>
              </CollapsibleTrigger>
              <CollapsibleContent>
                <div className="mt-2 p-2 bg-brutal-border/10 rounded text-left text-xs font-mono overflow-auto max-h-40">
                  <p><strong>Message:</strong> {error.message}</p>
                  {error.stack && (
                    <p className="mt-1"><strong>Stack:</strong> {error.stack.split('\n').slice(0, 3).join('\n')}</p>
                  )}
                  <p className="mt-1"><strong>API URL:</strong> {import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'}</p>
                </div>
              </CollapsibleContent>
            </Collapsible>
          )}
        </div>
      </Card>

      {/* Display fallback UI for dashboard components */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <Card className="p-4 animate-pulse">
          <div className="h-24 bg-brutal-border/20 rounded-md"></div>
        </Card>
        <Card className="p-4 animate-pulse">
          <div className="h-24 bg-brutal-border/20 rounded-md"></div>
        </Card>
        <Card className="p-4 animate-pulse">
          <div className="h-24 bg-brutal-border/20 rounded-md"></div>
        </Card>
      </div>
    </div>
  );
}
