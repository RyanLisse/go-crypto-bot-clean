import { useState, useEffect, useCallback, useRef } from 'react';
import { createWebSocketClient, WebSocketMessageType, AccountUpdatePayload, PortfolioUpdatePayload } from '@/lib/api';
import { API_CONFIG } from '@/config';

// Flag to enable/disable WebSocket connection
const WEBSOCKET_ENABLED = true;

// Define the return type for the hook
interface UseWebSocketReturn {
  isConnected: boolean;
  accountData: AccountUpdatePayload | null;
  portfolioData: PortfolioUpdatePayload | null;
  connect: () => void;
  disconnect: () => void;
  lastMessageTime: Date | null;
}

export function useWebSocket(): UseWebSocketReturn {
  const [isConnected, setIsConnected] = useState(false);
  const [accountData, setAccountData] = useState<AccountUpdatePayload | null>(null);
  const [portfolioData, setPortfolioData] = useState<PortfolioUpdatePayload | null>(null);
  const [lastMessageTime, setLastMessageTime] = useState<Date | null>(null);

  // Use a ref to keep the WebSocket client instance
  const wsClientRef = useRef(createWebSocketClient());

  // Connect to WebSocket
  const connect = useCallback(() => {
    if (!WEBSOCKET_ENABLED) {
      console.log('WebSocket connection disabled');
      return null;
    }

    try {
      console.log('Connecting to WebSocket at:', API_CONFIG.WS_URL);
      const socket = wsClientRef.current.connect();
      if (socket) {
        console.log('WebSocket connected successfully');
        setIsConnected(true);
      } else {
        console.error('WebSocket connection failed - null socket returned');
      }
      return socket;
    } catch (error) {
      console.error('Failed to connect to WebSocket:', error);
      return null;
    }
  }, []);

  // Disconnect from WebSocket
  const disconnect = useCallback(() => {
    console.log('Disconnecting WebSocket');
    wsClientRef.current.disconnect();
    setIsConnected(false);
  }, []);

  // Set up event listeners when the component mounts
  useEffect(() => {
    if (!WEBSOCKET_ENABLED) {
      // Don't attempt to connect if disabled
      console.log('WebSocket disabled, not connecting');
      return;
    }

    const wsClient = wsClientRef.current;
    console.log('Setting up WebSocket event listeners');

    // Handle account updates
    const handleAccountUpdate = (data: AccountUpdatePayload) => {
      console.log('Account update received:', data);
      setAccountData(data);
      setLastMessageTime(new Date());
    };

    // Handle portfolio updates
    const handlePortfolioUpdate = (data: PortfolioUpdatePayload) => {
      console.log('Portfolio update received:', data);
      setPortfolioData(data);
      setLastMessageTime(new Date());
    };

    // Handle authentication success
    const handleAuthSuccess = (data: unknown) => {
      console.log('WebSocket authenticated successfully:', data);
      setIsConnected(true);
    };

    // Handle authentication failure
    const handleAuthFailure = (data: unknown) => {
      console.error('WebSocket authentication failed:', data);
      setIsConnected(false);
    };

    // Handle connection errors
    const handleError = (data: unknown) => {
      console.error('WebSocket error:', data);
      setIsConnected(false);
    };

    try {
      // Add event listeners
      wsClient.addEventListener(WebSocketMessageType.ACCOUNT_UPDATE, handleAccountUpdate);
      wsClient.addEventListener(WebSocketMessageType.PORTFOLIO_UPDATE, handlePortfolioUpdate);
      wsClient.addEventListener(WebSocketMessageType.AUTH_SUCCESS, handleAuthSuccess);
      wsClient.addEventListener(WebSocketMessageType.AUTH_FAILURE, handleAuthFailure);
      wsClient.addEventListener(WebSocketMessageType.ERROR, handleError);

      // Connect to WebSocket
      console.log('Initiating WebSocket connection');
      connect();

      // Set up a reconnection timer
      const reconnectInterval = setInterval(() => {
        if (!wsClient.isConnected) {
          console.log('WebSocket disconnected, attempting to reconnect...');
          connect();
        }
      }, 10000); // Try to reconnect every 10 seconds if disconnected

      // Clean up when the component unmounts
      return () => {
        try {
          clearInterval(reconnectInterval);
          wsClient.removeEventListener(WebSocketMessageType.ACCOUNT_UPDATE, handleAccountUpdate);
          wsClient.removeEventListener(WebSocketMessageType.PORTFOLIO_UPDATE, handlePortfolioUpdate);
          wsClient.removeEventListener(WebSocketMessageType.AUTH_SUCCESS, handleAuthSuccess);
          wsClient.removeEventListener(WebSocketMessageType.AUTH_FAILURE, handleAuthFailure);
          wsClient.removeEventListener(WebSocketMessageType.ERROR, handleError);
          disconnect();
        } catch (error) {
          console.error('Error during WebSocket cleanup:', error);
        }
      };
    } catch (error) {
      console.error('Error setting up WebSocket:', error);
      return () => {}; // Empty cleanup function
    }
  }, [connect, disconnect]);

  return {
    isConnected,
    accountData,
    portfolioData,
    connect,
    disconnect,
    lastMessageTime
  };
}
