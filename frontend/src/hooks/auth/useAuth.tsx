import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { authService } from '../../services/authService';
import { User } from '@/types';

/**
 * Hook for authentication state and actions
 */
export function useAuth() {
  const queryClient = useQueryClient();

  // Query for authentication state
  const { data: user, isLoading } = useQuery<User | null>({
    queryKey: ['auth', 'user'],
    queryFn: () => authService.getCurrentUser(),
    initialData: authService.getCurrentUser(),
    staleTime: Infinity, // Don't refetch automatically
  });

  // Mutation for login
  const loginMutation = useMutation({
    mutationFn: ({ username, password }: { username: string; password: string }) => 
      authService.login(username, password),
    onSuccess: () => {
      // Invalidate auth queries to trigger a refetch
      queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });

  // Mutation for logout
  const logoutMutation = useMutation({
    mutationFn: () => authService.logout(),
    onSuccess: () => {
      // Invalidate auth queries to trigger a refetch
      queryClient.invalidateQueries({ queryKey: ['auth'] });
      // Clear all queries to prevent showing stale data after logout
      queryClient.clear();
    },
  });

  // Check if user is authenticated
  const isAuthenticated = authService.isAuthenticated();

  // Login function
  const login = (username: string, password: string) => {
    return loginMutation.mutate({ username, password });
  };

  // Logout function
  const logout = () => {
    return logoutMutation.mutate();
  };

  return {
    user,
    isAuthenticated,
    isLoading,
    login,
    logout,
    loginError: loginMutation.error,
    isLoggingIn: loginMutation.isPending,
    isLoggingOut: logoutMutation.isPending,
  };
}

export default useAuth;
