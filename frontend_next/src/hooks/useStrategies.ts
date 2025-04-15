'use client';

import { useRepository } from './useRepository';
import { createRepositories } from '@/lib/repositories';
import { Strategy, StrategyParameter, StrategyPerformance } from '@/lib/repositories/strategy.repository';
import { useQuery } from '@tanstack/react-query';

const repositories = createRepositories();

/**
 * Custom hook for working with strategies
 * @returns Strategy-specific repository methods with React Query
 */
export function useStrategies() {
  const baseRepo = useRepository<Strategy>(repositories.strategyRepository, 'strategies');
  
  /**
   * Fetch strategy parameters
   */
  const getParameters = (strategyId: string | null, options = {}) => {
    return useQuery({
      queryKey: ['strategy-parameters', strategyId],
      queryFn: () => {
        if (!strategyId) throw new Error('Strategy ID is required');
        return repositories.strategyRepository.getParameters(strategyId);
      },
      enabled: !!strategyId,
      ...options
    });
  };
  
  /**
   * Fetch strategy performance metrics
   */
  const getPerformance = (strategyId: string | null, options = {}) => {
    return useQuery({
      queryKey: ['strategy-performance', strategyId],
      queryFn: () => {
        if (!strategyId) throw new Error('Strategy ID is required');
        return repositories.strategyRepository.getPerformance(strategyId);
      },
      enabled: !!strategyId,
      ...options
    });
  };
  
  /**
   * Sync full strategy data
   */
  const syncFullStrategyData = (strategyId: string | null, options = {}) => {
    return useQuery({
      queryKey: ['strategy-full', strategyId],
      queryFn: () => {
        if (!strategyId) throw new Error('Strategy ID is required');
        return repositories.strategyRepository.syncFullStrategyData(strategyId);
      },
      enabled: !!strategyId,
      ...options
    });
  };
  
  return {
    ...baseRepo,
    getParameters,
    getPerformance,
    syncFullStrategyData
  };
} 