import { useState, useEffect, useCallback } from 'react';
import websocketClient, { WebSocketConnectionState, WebSocketEventHandlers } from '@/lib/websocket';

/**
 * Hook for using the WebSocket client
 * @param eventHandlers WebSocket event handlers
 * @returns WebSocket client and connection state
 */
export function useWebSocket(eventHandlers: WebSocketEventHandlers = {}) {
  const [connectionState, setConnectionState] = useState<WebSocketConnectionState>(
    websocketClient.getConnectionState()
  );
  const [lastMessage, setLastMessage] = useState<any>(null);

  // Connect to WebSocket
  const connect = useCallback(() => {
    websocketClient.connect();
  }, []);

  // Disconnect from WebSocket
  const disconnect = useCallback(() => {
    websocketClient.disconnect();
  }, []);

  // Send a message to the WebSocket server
  const sendMessage = useCallback((message: any) => {
    websocketClient.send(message);
  }, []);

  // Subscribe to a ticker
  const subscribeTicker = useCallback((symbols: string[]) => {
    websocketClient.subscribeTicker(symbols);
  }, []);

  // Update connection state when it changes
  useEffect(() => {
    const handleOpen = () => {
      setConnectionState(WebSocketConnectionState.OPEN);
      if (eventHandlers.onOpen) {
        eventHandlers.onOpen();
      }
    };

    const handleClose = (event: CloseEvent) => {
      setConnectionState(WebSocketConnectionState.CLOSED);
      if (eventHandlers.onClose) {
        eventHandlers.onClose(event);
      }
    };

    const handleError = (event: Event) => {
      if (eventHandlers.onError) {
        eventHandlers.onError(event);
      }
    };

    const handleMessage = (data: any) => {
      setLastMessage(data);
      if (eventHandlers.onMessage) {
        eventHandlers.onMessage(data);
      }
    };

    // Set event handlers
    websocketClient.setEventHandlers({
      onOpen: handleOpen,
      onClose: handleClose,
      onError: handleError,
      onMessage: handleMessage,
      onMarketData: eventHandlers.onMarketData,
      onTradeNotification: eventHandlers.onTradeNotification,
      onNewCoinAlert: eventHandlers.onNewCoinAlert,
      onSubscriptionSuccess: eventHandlers.onSubscriptionSuccess,
    });

    // Connect to WebSocket on mount
    connect();

    // Disconnect from WebSocket on unmount
    return () => {
      disconnect();
    };
  }, [
    connect,
    disconnect,
    eventHandlers.onOpen,
    eventHandlers.onClose,
    eventHandlers.onError,
    eventHandlers.onMessage,
    eventHandlers.onMarketData,
    eventHandlers.onTradeNotification,
    eventHandlers.onNewCoinAlert,
    eventHandlers.onSubscriptionSuccess,
  ]);

  return {
    isConnected: connectionState === WebSocketConnectionState.OPEN,
    connectionState,
    lastMessage,
    sendMessage,
    subscribeTicker,
    connect,
    disconnect
  };
};
