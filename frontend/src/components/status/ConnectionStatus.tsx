import React, { useEffect, useState } from 'react';
import { useBackendStatus } from '@/hooks/useBackendStatus';
import { getAIMetrics } from '@/lib/aiClient';
import { toast } from 'sonner';
import { AlertTriangle, Wifi, WifiOff, Server, Clock, RefreshCw } from 'lucide-react';

interface ConnectionStatusProps {
  showIndicator?: boolean;
  className?: string;
  showDetails?: boolean;
  onRefresh?: () => void;
}

export function ConnectionStatus({
  showIndicator = true,
  className = '',
  showDetails = false,
  onRefresh
}: ConnectionStatusProps) {
  const [fallbackNotified, setFallbackNotified] = useState(false);
  const [lastChecked, setLastChecked] = useState<Date>(new Date());
  const [showTooltip, setShowTooltip] = useState(false);

  // Use the enhanced backend status hook with callbacks
  const { isConnected, isLoading, status, refetch } = useBackendStatus({
    refetchInterval: 5000,
    onSuccess: () => setLastChecked(new Date()),
  });

  // Check if we're using fallback mode
  const usingFallback = getAIMetrics().usingFallback;

  // Show a toast notification when fallback mode is activated
  useEffect(() => {
    if (usingFallback && !fallbackNotified) {
      toast.warning(
        'Using AI fallback mode',
        {
          description: 'Backend connection unavailable. Using local AI model with limited capabilities.',
          duration: 5000,
          icon: <AlertTriangle className="h-4 w-4" />,
        }
      );
      setFallbackNotified(true);
    }
  }, [usingFallback, fallbackNotified]);

  // Reset notification state when connection is restored
  useEffect(() => {
    if (isConnected && !usingFallback) {
      setFallbackNotified(false);

      // Show a toast when connection is restored
      if (fallbackNotified) {
        toast.success(
          'Backend connection restored',
          {
            description: 'Full functionality has been restored.',
            duration: 3000,
          }
        );
      }
    }
  }, [isConnected, usingFallback, fallbackNotified]);

  // Handle manual refresh
  const handleRefresh = () => {
    refetch();
    if (onRefresh) onRefresh();
  };

  if (!showIndicator) {
    return null;
  }

  return (
    <div className={`flex items-center ${className}`}>
      <div
        className="relative flex items-center cursor-pointer"
        onMouseEnter={() => setShowTooltip(true)}
        onMouseLeave={() => setShowTooltip(false)}
      >
        {isLoading ? (
          <div className="flex items-center text-brutal-warning">
            <RefreshCw className="h-4 w-4 mr-1 animate-spin" />
            <span className="uppercase text-xs tracking-wider">CONNECTING</span>
          </div>
        ) : isConnected ? (
          <div className="flex items-center text-brutal-success">
            <Wifi className="h-4 w-4 mr-1" />
            <span className="uppercase text-xs tracking-wider">CONNECTED</span>
            {status?.version && (
              <span className="ml-2 text-xs bg-brutal-success/20 text-brutal-success px-1 py-0.5">
                v{status.version}
              </span>
            )}
          </div>
        ) : (
          <div className="flex items-center text-brutal-error">
            <WifiOff className="h-4 w-4 mr-1" />
            <span className="uppercase text-xs tracking-wider">OFFLINE</span>
            {usingFallback && (
              <span className="ml-2 text-xs bg-brutal-warning/20 text-brutal-warning px-1 py-0.5">
                FALLBACK MODE
              </span>
            )}
          </div>
        )}

        {/* Refresh button */}
        <button
          onClick={handleRefresh}
          className="ml-2 p-1 hover:bg-brutal-panel rounded"
          title="Refresh connection status"
        >
          <RefreshCw className="h-3 w-3" />
        </button>

        {/* Tooltip with detailed status */}
        {showTooltip && showDetails && (
          <div className="absolute top-full mt-2 left-0 z-50 bg-brutal-background border-2 border-brutal-border p-2 rounded shadow-md text-xs font-mono w-64">
            <div className="flex justify-between items-center mb-1">
              <span className="font-bold">Backend Status</span>
              <span className={isConnected ? 'text-brutal-success' : 'text-brutal-error'}>
                {isConnected ? 'ONLINE' : 'OFFLINE'}
              </span>
            </div>

            {isConnected && status && (
              <>
                <div className="flex justify-between">
                  <span>Version:</span>
                  <span>{status.version || 'Unknown'}</span>
                </div>
                <div className="flex justify-between">
                  <span>Uptime:</span>
                  <span>{status.uptime || 'Unknown'}</span>
                </div>
                {status.processes && (
                  <div className="mt-1">
                    <span className="font-bold">Services:</span>
                    <div className="ml-2 mt-1">
                      {Object.entries(status.processes).map(([name, proc]: [string, any]) => (
                        <div key={name} className="flex justify-between">
                          <span>{name}:</span>
                          <span className={proc.status === 'running' ? 'text-brutal-success' : 'text-brutal-error'}>
                            {proc.status}
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </>
            )}

            <div className="flex justify-between mt-1">
              <span>AI Mode:</span>
              <span className={usingFallback ? 'text-brutal-warning' : 'text-brutal-success'}>
                {usingFallback ? 'FALLBACK' : 'ONLINE'}
              </span>
            </div>

            <div className="text-brutal-text/50 text-[10px] mt-2">
              Last checked: {lastChecked.toLocaleTimeString()}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default ConnectionStatus;
