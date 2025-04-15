'use client';

import { useState } from 'react';
import { 
  useQuery, 
  useMutation, 
  useQueryClient,
  UseMutationResult,
  UseQueryResult 
} from '@tanstack/react-query';
import { BaseRepository } from '@/lib/repositories/interfaces';

/**
 * Custom hook wrapping repository with React Query
 * @param repository Repository instance
 * @param queryKey Base query key for React Query
 * @returns Object with query and mutation hooks
 */
export function useRepository<T>(repository: BaseRepository<T>, queryKey: string) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const queryClient = useQueryClient();
  
  /**
   * Get all entities
   */
  const getAll = (options: { 
    enabled?: boolean, 
    staleTime?: number, 
    refetchInterval?: number 
  } = {}) => {
    return useQuery({
      queryKey: [queryKey],
      queryFn: () => repository.getAll(),
      enabled: options.enabled,
      staleTime: options.staleTime,
      refetchInterval: options.refetchInterval,
    });
  };
  
  /**
   * Get entity by ID
   */
  const getById = (
    id: string | null, 
    options: { 
      enabled?: boolean, 
      staleTime?: number, 
      refetchInterval?: number 
    } = {}
  ): UseQueryResult<T, Error> => {
    return useQuery({
      queryKey: [queryKey, id],
      queryFn: () => {
        if (!id) throw new Error('ID is required');
        return repository.getById(id);
      },
      enabled: !!id && (options.enabled ?? true),
      staleTime: options.staleTime,
      refetchInterval: options.refetchInterval,
    });
  };
  
  /**
   * Create a new entity
   */
  const create = (): UseMutationResult<T, Error, Partial<T>, unknown> => {
    return useMutation({
      mutationFn: (data: Partial<T>) => repository.create(data),
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: [queryKey] });
      },
      onError: (err: Error) => {
        setError(err);
      },
    });
  };
  
  /**
   * Update an existing entity
   */
  const update = (): UseMutationResult<T, Error, { id: string, data: Partial<T> }, unknown> => {
    return useMutation({
      mutationFn: ({ id, data }: { id: string, data: Partial<T> }) => repository.update(id, data),
      onSuccess: (_, variables) => {
        queryClient.invalidateQueries({ queryKey: [queryKey] });
        queryClient.invalidateQueries({ queryKey: [queryKey, variables.id] });
      },
      onError: (err: Error) => {
        setError(err);
      },
    });
  };
  
  /**
   * Delete an entity
   */
  const remove = (): UseMutationResult<void, Error, string, unknown> => {
    return useMutation({
      mutationFn: (id: string) => repository.delete(id),
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: [queryKey] });
      },
      onError: (err: Error) => {
        setError(err);
      },
    });
  };
  
  return {
    loading,
    error,
    getAll,
    getById,
    create,
    update,
    remove,
  };
} 