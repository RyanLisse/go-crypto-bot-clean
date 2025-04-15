'use client';

import React, { useState, useEffect } from 'react';
import { API_CONFIG } from '@/config';

export const ApiToggle: React.FC = () => {
  const [useLocalApi, setUseLocalApi] = useState<boolean>(API_CONFIG.USE_LOCAL_API);
  
  // Update localStorage when toggle changes
  useEffect(() => {
    localStorage.setItem('useLocalApi', useLocalApi.toString());
    // Implement reload logic if needed to apply change
  }, [useLocalApi]);
  
  const handleToggle = () => {
    setUseLocalApi(prev => !prev);
  };
  
  return (
    <div className="flex items-center space-x-2 py-2">
      <span className={`text-sm ${!useLocalApi ? 'font-medium' : ''}`}>
        Remote API
      </span>
      <button
        onClick={handleToggle}
        className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 ${
          useLocalApi ? 'bg-indigo-600' : 'bg-gray-200'
        }`}
      >
        <span 
          className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
            useLocalApi ? 'translate-x-6' : 'translate-x-1'
          }`} 
        />
      </button>
      <span className={`text-sm ${useLocalApi ? 'font-medium' : ''}`}>
        Local API
      </span>
    </div>
  );
};

export default ApiToggle; 