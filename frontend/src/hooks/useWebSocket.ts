import { useState, useEffect, useCallback } from 'react';

interface WebSocketConfig {
  url: string;
  onMessage: (message: string) => void;
  onOpen?: () => void;
  onClose?: () => void;
  onError?: (error: Event) => void;
}

interface WebSocketHookResult {
  isConnected: boolean;
  sendMessage: (message: string) => void;
}

export const useWebSocket = (config: WebSocketConfig): WebSocketHookResult => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    const ws = new WebSocket(config.url);

    ws.onopen = () => {
      setIsConnected(true);
      config.onOpen?.();
    };

    ws.onclose = () => {
      setIsConnected(false);
      config.onClose?.();
    };

    ws.onmessage = (event) => {
      config.onMessage(event.data);
    };

    ws.onerror = (error) => {
      config.onError?.(error);
    };

    setSocket(ws);

    return () => {
      ws.close();
    };
  }, [config.url]);

  const sendMessage = useCallback((message: string) => {
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(message);
    }
  }, [socket]);

  return { isConnected, sendMessage };
};
