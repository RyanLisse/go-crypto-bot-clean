import React, { useState, useEffect } from 'react';
import { API_CONFIG } from '@/config';
import { useQueryClient } from '@tanstack/react-query';

export function ApiToggle() {
  const [isLocal, setIsLocal] = useState(API_CONFIG.USE_LOCAL_API);
  const queryClient = useQueryClient();

  // Update the API_CONFIG when the toggle changes
  const handleToggle = () => {
    // Toggle the API_CONFIG
    API_CONFIG.USE_LOCAL_API = !API_CONFIG.USE_LOCAL_API;
    
    // Update the state
    setIsLocal(API_CONFIG.USE_LOCAL_API);
    
    // Invalidate all queries to force a refetch with the new API URL
    queryClient.invalidateQueries();
    
    // Log the change
    console.log('API URL changed to:', API_CONFIG.API_URL);
  };

  // Initialize the state from the config on mount
  useEffect(() => {
    setIsLocal(API_CONFIG.USE_LOCAL_API);
  }, []);

  return (
    <div className="flex items-center space-x-2 p-2 bg-brutal-panel rounded-md border border-brutal-border">
      <span className="text-xs text-brutal-text/70">API:</span>
      <button
        onClick={handleToggle}
        className={`px-2 py-1 text-xs rounded-md transition-colors ${
          isLocal
            ? 'bg-green-500/20 text-green-500 hover:bg-green-500/30'
            : 'bg-blue-500/20 text-blue-500 hover:bg-blue-500/30'
        }`}
      >
        {isLocal ? 'Local' : 'Remote'}
      </button>
    </div>
  );
}
