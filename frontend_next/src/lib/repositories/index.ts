import { StrategyRepository } from './strategy.repository';
import { BacktestRepository } from './backtest.repository';

/**
 * Factory function that creates repository instances
 * @returns Object containing repository instances
 */
export function createRepositories() {
  return {
    strategy: new StrategyRepository(),
    backtest: new BacktestRepository(),
  };
}

// Export repository interfaces and types
export * from './interfaces';
export * from './strategy.repository';
export * from './backtest.repository'; 