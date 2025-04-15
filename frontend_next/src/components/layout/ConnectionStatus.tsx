import React from 'react';
import { useWebSocket } from '@/hooks/use-websocket';
import { useStatusQuery } from '@/hooks/queries';
import { Wifi, WifiOff, Server } from 'lucide-react';
import { cn } from '@/lib/utils';

// Flag to disable backend connection checks during development
const DISABLE_BACKEND_CHECK = true;

export function ConnectionStatus() {
  const { isConnected } = useWebSocket();

  // Skip the backend status query if disabled
  const {
    data: statusData,
    isLoading,
    isError
  } = useStatusQuery({
    enabled: !DISABLE_BACKEND_CHECK
  });

  // Always show as connected in development mode if backend check is disabled
  const isBackendConnected = DISABLE_BACKEND_CHECK ? true : (!isError && !isLoading && statusData?.status === 'ok');

  return (
    <div className="flex items-center space-x-3">
      <div className="flex items-center">
        <div
          className={cn(
            "w-2 h-2 rounded-full mr-2",
            isBackendConnected ? "bg-brutal-success" : "bg-brutal-error"
          )}
        />
        <div className="flex items-center text-xs">
          <Server className="h-3 w-3 mr-1" />
          <span>Backend</span>
        </div>
      </div>

      <div className="flex items-center">
        <div
          className={cn(
            "w-2 h-2 rounded-full mr-2",
            isConnected ? "bg-brutal-success" : "bg-brutal-error"
          )}
        />
        <div className="flex items-center text-xs">
          {isConnected ? (
            <>
              <Wifi className="h-3 w-3 mr-1" />
              <span>Live</span>
            </>
          ) : (
            <>
              <WifiOff className="h-3 w-3 mr-1" />
              <span>Offline</span>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
