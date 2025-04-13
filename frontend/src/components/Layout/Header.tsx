
import React, { useState, useEffect } from 'react';
import { useBackendStatus } from '../../hooks/useBackendStatus';

export function Header() {
  const { isConnected, isLoading } = useBackendStatus({
    refetchInterval: 5000, // Check every 5 seconds
  });

  const [currentDateTime, setCurrentDateTime] = useState(new Date());

  // Update the date and time every second
  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentDateTime(new Date());
    }, 1000);

    return () => {
      clearInterval(timer);
    };
  }, []);

  const currentDate = currentDateTime.toLocaleDateString('en-US', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  });

  const currentTime = currentDateTime.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  });

  return (
    <header className="h-12 md:h-14 border-b border-brutal-border flex items-center justify-between px-3 md:px-6">
      <div className="flex items-center">
        <div className="uppercase text-xs tracking-wider">
          DASHBOARD<span className="text-brutal-text/30 mx-2">/</span>OVERVIEW
        </div>
      </div>

      <div className="flex items-center space-x-3 md:space-x-6 text-xs md:text-sm">
        <div className="hidden sm:flex items-center">
          {isConnected ? (
            <>
              <div className="h-2 w-2 rounded-full bg-green-500 mr-1"></div>
              <span className="uppercase text-xs tracking-wider mr-2">Backend</span>
              <div className="h-2 w-2 rounded-full bg-green-500 mr-1"></div>
              <span className="uppercase text-xs tracking-wider">Online</span>
            </>
          ) : (
            <>
              <div className="h-2 w-2 rounded-full bg-green-500 mr-1"></div>
              <span className="uppercase text-xs tracking-wider mr-2">Backend</span>
              <div className="h-2 w-2 rounded-full bg-red-500 mr-1"></div>
              <span className="uppercase text-xs tracking-wider">{isLoading ? 'Connecting...' : 'Offline'}</span>
            </>
          )}
        </div>
        <div className="text-brutal-text/70 font-mono">
          {currentDate.replace(/\//g, '/')} {currentTime}
        </div>
      </div>
    </header>
  );
}
