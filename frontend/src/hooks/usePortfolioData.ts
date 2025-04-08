import { useQuery } from '@tanstack/react-query';
import { PortfolioData } from '@/types/portfolio';

/**
 * Fetch portfolio data from the API
 */
async function fetchPortfolioData(): Promise<PortfolioData> {
  const response = await fetch('/api/v1/portfolio');
  
  if (!response.ok) {
    throw new Error('Failed to fetch portfolio data');
  }
  
  return response.json();
}

/**
 * Hook for fetching portfolio data
 */
export function usePortfolioData() {
  return useQuery({
    queryKey: ['portfolioData'],
    queryFn: fetchPortfolioData,
    refetchInterval: 60000, // Refetch every minute
    staleTime: 30000, // Consider data stale after 30 seconds
  });
}
