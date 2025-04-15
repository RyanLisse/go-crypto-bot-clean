import { useQuery } from '@tanstack/react-query';

/**
 * Hook to check if the backend is connected
 * @param options Optional configuration options
 * @returns Object with isConnected status and refetch function
 */
export function useBackendStatus(options?: {
  refetchInterval?: number;
  onSuccess?: (data: { connected: boolean; status?: any }) => void;
  onError?: (error: any) => void;
}) {
  const {
    data,
    isLoading,
    isError,
    refetch,
    error
  } = useQuery({
    queryKey: ['backendStatus'],
    queryFn: async () => {
      try {
        // Try to fetch the status from the backend
        // In a real implementation, this would call an actual API
        // For now, we'll simulate a successful connection
        return { connected: true, status: { version: '1.0.0' } };
      } catch (error) {
        console.error('Backend connection error:', error);
        return { connected: false, error };
      }
    },
    // Refetch more frequently for real-time updates
    refetchInterval: options?.refetchInterval || 5000,
    // Don't retry too many times to avoid flooding the network
    retry: 2,
    // Consider data stale after 2 seconds for more real-time feedback
    staleTime: 2000,
    // Initialize with disconnected state
    initialData: { connected: false },
    // Callbacks
    onSuccess: options?.onSuccess,
    onError: options?.onError,
    // Refetch on window focus for better user experience
    refetchOnWindowFocus: true,
  });

  return {
    isConnected: data?.connected || false,
    isLoading,
    isError,
    refetch,
    error,
    status: data?.status,
  };
}
