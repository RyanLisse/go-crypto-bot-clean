import { useState, useEffect } from 'react';

// Define the shape of the account data
interface AccountData {
  balances: {
    [key: string]: {
      asset: string;
      total: number;
      price: number;
    };
  };
  updatedAt: string;
}

export function useWebSocket() {
  const [isConnected, setIsConnected] = useState(false);
  const [accountData, setAccountData] = useState<AccountData | null>(null);

  // Mock websocket connection for testing
  useEffect(() => {
    // In a real app, this would establish a WebSocket connection

    // Simulate connected state after a delay
    const connectTimeout = setTimeout(() => {
      setIsConnected(true);
    }, 1000);

    // Mock receiving data periodically
    const dataInterval = setInterval(() => {
      if (isConnected) {
        // Simulate receiving account data
        setAccountData({
          balances: {
            BTC: { asset: 'BTC', total: 0.5, price: 60000 + Math.random() * 1000 },
            ETH: { asset: 'ETH', total: 5, price: 3000 + Math.random() * 100 },
            USDT: { asset: 'USDT', total: 1000, price: 1 },
          },
          updatedAt: new Date().toISOString(),
        });
      }
    }, 5000);

    // Clean up
    return () => {
      clearTimeout(connectTimeout);
      clearInterval(dataInterval);
      setIsConnected(false);
    };
  }, [isConnected]);

  return {
    isConnected,
    accountData,
  };
} 