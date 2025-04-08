import { useState, useEffect, useCallback, useRef } from 'react';
import { createWebSocketClient, WebSocketMessageType, AccountUpdatePayload, PortfolioUpdatePayload } from '@/lib/api';

// Flag to disable WebSocket connection during development
const DISABLE_WEBSOCKET = true;

// Define the return type for the hook
interface UseWebSocketReturn {
  isConnected: boolean;
  accountData: AccountUpdatePayload | null;
  portfolioData: PortfolioUpdatePayload | null;
  connect: () => void;
  disconnect: () => void;
}

export function useWebSocket(): UseWebSocketReturn {
  const [isConnected, setIsConnected] = useState(false);
  const [accountData, setAccountData] = useState<AccountUpdatePayload | null>(null);
  const [portfolioData, setPortfolioData] = useState<PortfolioUpdatePayload | null>(null);

  // Use a ref to keep the WebSocket client instance
  const wsClientRef = useRef(createWebSocketClient());

  // Connect to WebSocket
  const connect = useCallback(() => {
    if (DISABLE_WEBSOCKET) {
      console.log('WebSocket connection disabled in development mode');
      return null;
    }

    try {
      const socket = wsClientRef.current.connect();
      if (socket) {
        setIsConnected(true);
      }
      return socket;
    } catch (error) {
      console.error('Failed to connect to WebSocket:', error);
      return null;
    }
  }, []);

  // Disconnect from WebSocket
  const disconnect = useCallback(() => {
    wsClientRef.current.disconnect();
    setIsConnected(false);
  }, []);

  // Set up event listeners when the component mounts
  useEffect(() => {
    if (DISABLE_WEBSOCKET) {
      // Don't attempt to connect in development mode
      return;
    }

    const wsClient = wsClientRef.current;

    // Handle account updates
    const handleAccountUpdate = (data: AccountUpdatePayload) => {
      console.log('Account update received:', data);
      setAccountData(data);
    };

    // Handle portfolio updates
    const handlePortfolioUpdate = (data: PortfolioUpdatePayload) => {
      console.log('Portfolio update received:', data);
      setPortfolioData(data);
    };

    try {
      // Add event listeners
      wsClient.addEventListener(WebSocketMessageType.ACCOUNT_UPDATE, handleAccountUpdate);
      wsClient.addEventListener(WebSocketMessageType.PORTFOLIO_UPDATE, handlePortfolioUpdate);

      // Connect to WebSocket
      connect();

      // Clean up when the component unmounts
      return () => {
        try {
          wsClient.removeEventListener(WebSocketMessageType.ACCOUNT_UPDATE, handleAccountUpdate);
          wsClient.removeEventListener(WebSocketMessageType.PORTFOLIO_UPDATE, handlePortfolioUpdate);
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
    disconnect
  };
}
