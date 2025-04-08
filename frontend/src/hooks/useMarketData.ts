import { useQuery } from '@tanstack/react-query';

interface MarketCoin {
  symbol: string;
  price: number;
  priceChange: number;
  priceChangePercent: number;
  volume: number;
}

export interface MarketData {
  topGainers: MarketCoin[];
  topLosers: MarketCoin[];
  mostVolume: MarketCoin[];
  lastUpdated: string;
}

/**
 * Fetch market data from the API
 */
async function fetchMarketData(): Promise<MarketData> {
  const response = await fetch('/api/v1/market/overview');
  
  if (!response.ok) {
    throw new Error('Failed to fetch market data');
  }
  
  return response.json();
}

/**
 * Hook for fetching market data
 */
export function useMarketData() {
  return useQuery({
    queryKey: ['marketData'],
    queryFn: fetchMarketData,
    refetchInterval: 60000, // Refetch every minute
    staleTime: 30000, // Consider data stale after 30 seconds
  });
}
