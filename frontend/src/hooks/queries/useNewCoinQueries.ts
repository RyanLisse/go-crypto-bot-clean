import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';

// Query keys
export const newCoinKeys = {
  all: ['newCoins'] as const,
  list: () => [...newCoinKeys.all, 'list'] as const,
  upcoming: () => [...newCoinKeys.all, 'upcoming'] as const,
  upcomingTodayAndTomorrow: () => [...newCoinKeys.all, 'upcomingTodayAndTomorrow'] as const,
  byDate: (date: string) => [...newCoinKeys.all, 'byDate', date] as const,
  byDateRange: (startDate: string, endDate: string) =>
    [...newCoinKeys.all, 'byDateRange', startDate, endDate] as const,
};

// Get all new coins
export const useNewCoinsQuery = () => {
  return useQuery({
    queryKey: newCoinKeys.list(),
    queryFn: () => api.getNewCoins(),
    staleTime: 30000, // Consider data stale after 30 seconds
    retry: 2, // Retry failed requests up to 2 times
    refetchOnWindowFocus: true, // Refetch when window regains focus
  });
};

// Get upcoming coins
export const useUpcomingCoinsQuery = () => {
  return useQuery({
    queryKey: newCoinKeys.upcoming(),
    queryFn: () => api.getUpcomingCoins(),
    staleTime: 30000, // Consider data stale after 30 seconds
    retry: 2, // Retry failed requests up to 2 times
    refetchOnWindowFocus: true, // Refetch when window regains focus
  });
};

// Get new coins by date
export const useNewCoinsByDateQuery = (date: string) => {
  return useQuery({
    queryKey: newCoinKeys.byDate(date),
    queryFn: () => api.getNewCoinsByDate(date),
    staleTime: 30000, // Consider data stale after 30 seconds
    enabled: !!date, // Only run query if date is provided
  });
};

// Get new coins by date range
export const useNewCoinsByDateRangeQuery = (startDate: string, endDate: string) => {
  return useQuery({
    queryKey: newCoinKeys.byDateRange(startDate, endDate),
    queryFn: () => api.getNewCoinsByDateRange(startDate, endDate),
    staleTime: 30000, // Consider data stale after 30 seconds
    enabled: !!startDate && !!endDate, // Only run query if both dates are provided
  });
};

// Get upcoming coins for today and tomorrow
export const useUpcomingCoinsForTodayAndTomorrowQuery = () => {
  return useQuery({
    queryKey: newCoinKeys.upcomingTodayAndTomorrow(),
    queryFn: () => api.getUpcomingCoinsForTodayAndTomorrow(),
    staleTime: 30000, // Consider data stale after 30 seconds
    retry: 2, // Retry failed requests up to 2 times
    refetchOnWindowFocus: true, // Refetch when window regains focus
  });
};
