// API configuration
export const API_CONFIG = {
  USE_LOCAL_API: process.env.NEXT_PUBLIC_USE_LOCAL_API === 'true' || true,
  LOCAL_API_URL: process.env.NEXT_PUBLIC_LOCAL_API_URL || 'http://localhost:8080',
  REMOTE_API_URL: process.env.NEXT_PUBLIC_REMOTE_API_URL || 'https://api.crypto-bot.example.com',
  
  // Websocket configuration
  WS_LOCAL_URL: process.env.NEXT_PUBLIC_WS_LOCAL_URL || 'ws://localhost:8080/ws',
  WS_REMOTE_URL: process.env.NEXT_PUBLIC_WS_REMOTE_URL || 'wss://api.crypto-bot.example.com/ws',
  
  // Other configuration
  REFRESH_INTERVAL: parseInt(process.env.NEXT_PUBLIC_REFRESH_INTERVAL || '30000', 10),
  DASHBOARD_REFRESH_INTERVAL: parseInt(process.env.NEXT_PUBLIC_DASHBOARD_REFRESH_INTERVAL || '60000', 10),
  
  // Return the appropriate API URL based on configuration
  get API_URL() {
    return this.USE_LOCAL_API ? this.LOCAL_API_URL : this.REMOTE_API_URL;
  },
  
  // Return the appropriate Websocket URL based on configuration
  get WS_URL() {
    return this.USE_LOCAL_API ? this.WS_LOCAL_URL : this.WS_REMOTE_URL;
  }
}; 