import React, { useEffect, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { API_CONFIG } from '@/config';
import { AlertCircle, CheckCircle2, Wifi, WifiOff } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useWebSocket } from '@/hooks/use-websocket';

export function ApiConnectionStatus() {
  const [statusMessage, setStatusMessage] = useState<string>('Checking API connection...');
  const [isApiConnected, setIsApiConnected] = useState<boolean>(false);
  
  // Use WebSocket hook to check WebSocket connection
  const { isConnected: isWsConnected, lastMessageTime } = useWebSocket();
  
  // Query the API status
  const { data: apiStatus, isError, error, refetch } = useQuery({
    queryKey: ['apiStatus'],
    queryFn: async () => {
      const response = await fetch(`${API_CONFIG.API_URL}/status`);
      if (!response.ok) {
        throw new Error(`API status check failed: ${response.status}`);
      }
      return response.json();
    },
    retry: 2,
    refetchInterval: 30000, // Refetch every 30 seconds
  });
  
  // Update status message whenever the API status changes
  useEffect(() => {
    if (isError) {
      setStatusMessage(`Connection error: ${error instanceof Error ? error.message : 'Unknown error'}`);
      setIsApiConnected(false);
    } else if (apiStatus) {
      setStatusMessage(`API connected (v${apiStatus.version || 'unknown'})`);
      setIsApiConnected(true);
    } else {
      setStatusMessage('Waiting for API response...');
      setIsApiConnected(false);
    }
  }, [apiStatus, isError, error]);
  
  // Get the formatted time of the last WebSocket message
  const getLastMessageTime = () => {
    if (!lastMessageTime) return 'Never';
    return lastMessageTime.toLocaleTimeString();
  };
  
  return (
    <Card className="overflow-hidden">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg flex items-center gap-2">
          {isApiConnected ? (
            <CheckCircle2 className="h-5 w-5 text-green-500" />
          ) : (
            <AlertCircle className="h-5 w-5 text-amber-500" />
          )}
          API Connection
        </CardTitle>
        <CardDescription>
          Status of your connection to the backend API
        </CardDescription>
      </CardHeader>
      <CardContent className="text-sm">
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <span className="font-medium">API Status:</span>
            <span className={isApiConnected ? 'text-green-500' : 'text-amber-500'}>
              {isApiConnected ? 'Connected' : 'Disconnected'}
            </span>
          </div>
          
          <div className="flex items-center justify-between">
            <span className="font-medium">WebSocket:</span>
            <span className="flex items-center gap-1">
              {isWsConnected ? (
                <>
                  <Wifi className="h-4 w-4 text-green-500" />
                  <span className="text-green-500">Connected</span>
                </>
              ) : (
                <>
                  <WifiOff className="h-4 w-4 text-amber-500" />
                  <span className="text-amber-500">Disconnected</span>
                </>
              )}
            </span>
          </div>
          
          <div className="flex items-center justify-between">
            <span className="font-medium">API Endpoint:</span>
            <span className="text-xs opacity-70 truncate max-w-[200px]">
              {API_CONFIG.API_URL}
            </span>
          </div>
          
          <div className="flex items-center justify-between">
            <span className="font-medium">Last WS Message:</span>
            <span className="text-xs opacity-70">
              {getLastMessageTime()}
            </span>
          </div>
          
          <div className="mt-4 text-xs text-muted-foreground">
            {statusMessage}
          </div>
        </div>
      </CardContent>
    </Card>
  );
} 