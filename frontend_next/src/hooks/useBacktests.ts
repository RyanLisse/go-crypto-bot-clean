'use client';

import { useRepository } from './useRepository';
import { createRepositories } from '@/lib/repositories';
import { Backtest, BacktestResults } from '@/lib/repositories/backtest.repository';
import { useQuery } from '@tanstack/react-query';

const repositories = createRepositories();

/**
 * Custom hook for working with backtests
 * @returns Backtest-specific repository methods with React Query
 */
export function useBacktests() {
  const baseRepo = useRepository<Backtest>(repositories.backtestRepository, 'backtests');
  
  /**
   * Fetch backtest results
   */
  const getResults = (backtestId: string | null, options = {}) => {
    return useQuery({
      queryKey: ['backtest-results', backtestId],
      queryFn: () => {
        if (!backtestId) throw new Error('Backtest ID is required');
        return repositories.backtestRepository.getResults(backtestId);
      },
      enabled: !!backtestId,
      ...options
    });
  };
  
  /**
   * Fetch backtests for a specific strategy
   */
  const getByStrategy = (strategyId: string | null, options = {}) => {
    return useQuery({
      queryKey: ['backtests-by-strategy', strategyId],
      queryFn: () => {
        if (!strategyId) throw new Error('Strategy ID is required');
        return repositories.backtestRepository.getByStrategy(strategyId);
      },
      enabled: !!strategyId,
      ...options
    });
  };
  
  /**
   * Start a new backtest
   */
  const startBacktest = async (data: { strategyId: string, parameters: Record<string, any> }) => {
    return repositories.backtestRepository.startBacktest(data.strategyId, data.parameters);
  };
  
  /**
   * Cancel a running backtest
   */
  const cancelBacktest = async (backtestId: string) => {
    return repositories.backtestRepository.cancelBacktest(backtestId);
  };
  
  return {
    ...baseRepo,
    getResults,
    getByStrategy,
    startBacktest,
    cancelBacktest
  };
} 