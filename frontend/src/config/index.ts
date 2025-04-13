// API Configuration
export const API_CONFIG = {
  // Set this to true to use the local API, false to use the remote API
  USE_LOCAL_API: true,
  
  // Local API URL
  LOCAL_API_URL: 'http://localhost:8080/api/v1',
  LOCAL_WS_URL: 'ws://localhost:8080/ws',
  
  // Remote API URL (from environment variables)
  REMOTE_API_URL: import.meta.env.VITE_API_URL || 'https://full-backend-production-cd25.up.railway.app/api/v1',
  REMOTE_WS_URL: import.meta.env.VITE_WS_URL || 'wss://full-backend-production-cd25.up.railway.app/ws',
  
  // Get the active API URL based on the USE_LOCAL_API flag
  get API_URL() {
    return this.USE_LOCAL_API ? this.LOCAL_API_URL : this.REMOTE_API_URL;
  },
  
  // Get the active WebSocket URL based on the USE_LOCAL_API flag
  get WS_URL() {
    return this.USE_LOCAL_API ? this.LOCAL_WS_URL : this.REMOTE_WS_URL;
  }
};

// Other configuration settings
export const APP_CONFIG = {
  // App name
  APP_NAME: 'Crypto Bot',
  
  // Version
  VERSION: '1.0.0',
  
  // Debug mode
  DEBUG: true,
  
  // TanStack Query configuration
  QUERY_CONFIG: {
    // Default stale time (30 seconds)
    DEFAULT_STALE_TIME: 30000,
    
    // Default cache time (5 minutes)
    DEFAULT_CACHE_TIME: 300000,
    
    // Default retry count
    DEFAULT_RETRY_COUNT: 2,
    
    // Default refetch interval (30 seconds)
    DEFAULT_REFETCH_INTERVAL: 30000
  }
};
