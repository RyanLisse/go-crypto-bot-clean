import { useState, useEffect } from 'react';
import { useWalletQuery } from './queries/useAccountQueries';
import { useWebSocket } from './use-websocket';
import { API_CONFIG } from '@/config';

export type WalletBalance = {
  asset: string;
  free: number;
  locked: number;
  total: number;
  price?: number;
  value?: number;
};

export interface WalletData {
  balances: Record<string, WalletBalance>;
  totalValue: number;
  lastUpdated: string;
  dataSource: 'websocket' | 'api' | 'mock';
  isLoading: boolean;
  isError: boolean;
  error: Error | null;
  refetch: () => void;
}

/**
 * Custom hook to fetch and combine real wallet data from multiple sources
 * Prioritizes websocket data, then API data, then falls back to mock data
 */
export function useRealWalletData(): WalletData {
  const [combinedData, setCombinedData] = useState<WalletData>({
    balances: {},
    totalValue: 0,
    lastUpdated: new Date().toISOString(),
    dataSource: 'mock',
    isLoading: true,
    isError: false,
    error: null,
    refetch: () => {}
  });

  // Use WebSocket for real-time updates
  const { isConnected, accountData } = useWebSocket();

  // Use API query as backup and for initial data
  const {
    data: walletData,
    isLoading,
    isError,
    error,
    refetch
  } = useWalletQuery({
    refetchInterval: 30000, // Refetch every 30 seconds
    staleTime: 10000 // Consider data stale after 10 seconds
  });

  // Update combined data whenever sources change
  useEffect(() => {
    console.log('API URL being used:', API_CONFIG.API_URL);
    console.log('WebSocket connected:', isConnected);
    console.log('Account data from WebSocket:', accountData);
    console.log('Wallet data from API:', walletData);

    // Calculate new wallet data
    let newBalances: Record<string, WalletBalance> = {};
    let newTotalValue = 0;
    let newLastUpdated = new Date().toISOString();
    let newDataSource: 'websocket' | 'api' | 'mock' = 'mock';

    // Priority 1: Use WebSocket data if available
    if (isConnected && accountData && Object.keys(accountData.balances || {}).length > 0) {
      newBalances = accountData.balances;
      newLastUpdated = accountData.updatedAt;
      newDataSource = 'websocket';
      
      // Calculate total value
      newTotalValue = Object.values(newBalances).reduce((sum, balance) => {
        const price = balance.price || 0;
        const total = balance.total || 0;
        return sum + (price * total);
      }, 0);
    }
    // Priority 2: Use API data if available
    else if (walletData && Object.keys(walletData.balances || {}).length > 0) {
      newBalances = walletData.balances;
      newLastUpdated = walletData.updatedAt;
      newDataSource = 'api';
      
      // Calculate total value
      newTotalValue = Object.values(newBalances).reduce((sum, balance) => {
        const price = balance.price || 0;
        const total = balance.total || 0;
        return sum + (price * total);
      }, 0);
    }
    // Priority 3: Fall back to mock data if nothing else works
    else if (!isLoading && !isError) {
      // Create mock data if we don't have real data
      newBalances = {
        'BTC': {
          asset: 'BTC',
          free: 0.1,
          locked: 0,
          total: 0.1,
          price: 50000,
          value: 5000
        },
        'ETH': {
          asset: 'ETH',
          free: 1.5,
          locked: 0,
          total: 1.5,
          price: 3000,
          value: 4500
        },
        'USDT': {
          asset: 'USDT',
          free: 10000,
          locked: 0,
          total: 10000,
          price: 1,
          value: 10000
        }
      };
      
      // Calculate total value
      newTotalValue = Object.values(newBalances).reduce((sum, balance) => {
        return sum + (balance.value || 0);
      }, 0);
      
      newDataSource = 'mock';
    }

    // Update state with new data
    setCombinedData({
      balances: newBalances,
      totalValue: newTotalValue,
      lastUpdated: newLastUpdated,
      dataSource: newDataSource,
      isLoading,
      isError,
      error: error instanceof Error ? error : null,
      refetch
    });
  }, [isConnected, accountData, walletData, isLoading, isError, error, refetch]);

  return combinedData;
} 