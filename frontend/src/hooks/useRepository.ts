import { useState, useCallback } from 'react';
import { StrategyRepository, BacktestRepository } from '../api/repository';

// Repository factory to get the appropriate repository instance
const getRepository = (type: 'strategy' | 'backtest') => {
  switch (type) {
    case 'strategy':
      return new StrategyRepository();
    case 'backtest':
      return new BacktestRepository();
    default:
      throw new Error(`Unknown repository type: ${type}`);
  }
};

// Hook to use repositories with loading and error states
export function useRepository<T>(type: 'strategy' | 'backtest') {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  
  // Get repository instance
  const repository = getRepository(type);
  
  // Wrap repository methods with loading and error handling
  const withLoadingAndError = useCallback(
    <R>(fn: () => Promise<R>): Promise<R> => {
      setLoading(true);
      setError(null);
      
      return fn()
        .catch((err) => {
          setError(err);
          throw err;
        })
        .finally(() => {
          setLoading(false);
        });
    },
    []
  );
  
  // Wrapped repository methods
  const getAll = useCallback(
    () => withLoadingAndError(() => repository.getAll()),
    [repository, withLoadingAndError]
  );
  
  const getById = useCallback(
    (id: string) => withLoadingAndError(() => repository.getById(id)),
    [repository, withLoadingAndError]
  );
  
  const create = useCallback(
    (data: Partial<T>) => withLoadingAndError(() => repository.create(data)),
    [repository, withLoadingAndError]
  );
  
  const update = useCallback(
    (id: string, data: Partial<T>) => withLoadingAndError(() => repository.update(id, data)),
    [repository, withLoadingAndError]
  );
  
  const remove = useCallback(
    (id: string) => withLoadingAndError(() => repository.delete(id)),
    [repository, withLoadingAndError]
  );
  
  const syncFromApi = useCallback(
    () => withLoadingAndError(() => repository.syncFromApi()),
    [repository, withLoadingAndError]
  );
  
  // Return wrapped repository methods and state
  return {
    repository,
    loading,
    error,
    getAll,
    getById,
    create,
    update,
    remove,
    syncFromApi,
  };
}
