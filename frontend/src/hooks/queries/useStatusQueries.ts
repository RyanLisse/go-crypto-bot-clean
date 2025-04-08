import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api, StatusResponse } from '@/lib/api';

// Query keys
export const statusKeys = {
  all: ['status'] as const,
  details: () => [...statusKeys.all, 'details'] as const,
};

// Get system status
export const useStatusQuery = (options?: { enabled?: boolean }) => {
  return useQuery({
    queryKey: statusKeys.details(),
    queryFn: async () => {
      console.log('Fetching status in useStatusQuery...');
      try {
        const result = await api.getStatus();
        console.log('Status query result:', result);
        return result;
      } catch (error) {
        console.error('Error in useStatusQuery:', error);
        throw error;
      }
    },
    refetchInterval: 10000, // Refetch every 10 seconds
    staleTime: 5000, // Consider data stale after 5 seconds
    enabled: options?.enabled !== undefined ? options.enabled : true,
    retry: false, // Don't retry if the backend is down
    onError: (error) => {
      console.error('Status query error handler:', error);
    },
  });
};

// Start processes mutation
export const useStartProcessesMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => api.startProcesses(),
    onSuccess: (data) => {
      // Update the status query data
      queryClient.setQueryData(statusKeys.details(), data);
      // Invalidate the query to refetch
      queryClient.invalidateQueries({ queryKey: statusKeys.details() });
    },
  });
};

// Stop processes mutation
export const useStopProcessesMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => api.stopProcesses(),
    onSuccess: (data) => {
      // Update the status query data
      queryClient.setQueryData(statusKeys.details(), data);
      // Invalidate the query to refetch
      queryClient.invalidateQueries({ queryKey: statusKeys.details() });
    },
  });
};
