import { BaseRepository } from './base';
import { db } from '../../db/client';
import { strategies, strategyParameters, strategyPerformance } from '../../db/schema/strategies';
import { eq } from 'drizzle-orm';

// Strategy interface matching the backend model
export interface Strategy {
  id: string;
  name: string;
  description: string;
  parameters: Record<string, any>;
  isEnabled: boolean;
  userId: string;
  createdAt: Date;
  updatedAt: Date;
}

// Strategy parameter interface
export interface StrategyParameter {
  id: number;
  strategyId: string;
  name: string;
  type: string;
  description: string;
  defaultValue: string;
  min?: string;
  max?: string;
  options?: string[];
  required: boolean;
}

// Strategy performance interface
export interface StrategyPerformance {
  id: number;
  strategyId: string;
  winRate: number;
  profitFactor: number;
  sharpeRatio: number;
  maxDrawdown: number;
  totalTrades: number;
  periodStart: Date;
  periodEnd: Date;
}

export class StrategyRepository extends BaseRepository<Strategy> {
  protected endpoint = 'strategies';
  protected tableName = 'strategies';
  
  // Store strategies in local database
  protected async storeInLocalDb(items: Strategy[]): Promise<void> {
    // Use a transaction to ensure all operations succeed or fail together
    await db.transaction(async (tx) => {
      // For each strategy
      for (const strategy of items) {
        // Insert or update strategy
        await tx
          .insert(strategies)
          .values({
            id: strategy.id,
            userId: strategy.userId,
            name: strategy.name,
            description: strategy.description,
            parameters: strategy.parameters,
            isEnabled: strategy.isEnabled,
            createdAt: new Date(strategy.createdAt),
            updatedAt: new Date(strategy.updatedAt),
          })
          .onConflictDoUpdate({
            target: strategies.id,
            set: {
              name: strategy.name,
              description: strategy.description,
              parameters: strategy.parameters,
              isEnabled: strategy.isEnabled,
              updatedAt: new Date(strategy.updatedAt),
            },
          });
      }
    });
  }
  
  // Get parameters for a strategy
  async getParameters(strategyId: string): Promise<StrategyParameter[]> {
    return this.fetchApi<StrategyParameter[]>(this.getApiUrl(`/${strategyId}/parameters`));
  }
  
  // Store parameters in local database
  async storeParametersInLocalDb(strategyId: string, parameters: StrategyParameter[]): Promise<void> {
    await db.transaction(async (tx) => {
      // Delete existing parameters
      await tx.delete(strategyParameters).where(eq(strategyParameters.strategyId, strategyId));
      
      // Insert new parameters
      for (const param of parameters) {
        await tx.insert(strategyParameters).values({
          strategyId: param.strategyId,
          name: param.name,
          type: param.type,
          description: param.description,
          defaultValue: param.defaultValue,
          min: param.min,
          max: param.max,
          options: param.options,
          required: param.required,
        });
      }
    });
  }
  
  // Get performance metrics for a strategy
  async getPerformance(strategyId: string): Promise<StrategyPerformance[]> {
    return this.fetchApi<StrategyPerformance[]>(this.getApiUrl(`/${strategyId}/performance`));
  }
  
  // Store performance metrics in local database
  async storePerformanceInLocalDb(strategyId: string, performance: StrategyPerformance[]): Promise<void> {
    await db.transaction(async (tx) => {
      // Delete existing performance metrics
      await tx.delete(strategyPerformance).where(eq(strategyPerformance.strategyId, strategyId));
      
      // Insert new performance metrics
      for (const perf of performance) {
        await tx.insert(strategyPerformance).values({
          strategyId: perf.strategyId,
          winRate: perf.winRate,
          profitFactor: perf.profitFactor,
          sharpeRatio: perf.sharpeRatio,
          maxDrawdown: perf.maxDrawdown,
          totalTrades: perf.totalTrades,
          periodStart: new Date(perf.periodStart),
          periodEnd: new Date(perf.periodEnd),
        });
      }
    });
  }
  
  // Sync all strategy data (including parameters and performance)
  async syncFullStrategyData(strategyId: string): Promise<void> {
    try {
      // Get strategy details
      const strategy = await this.getById(strategyId);
      await this.storeInLocalDb([strategy]);
      
      // Get and store parameters
      const parameters = await this.getParameters(strategyId);
      await this.storeParametersInLocalDb(strategyId, parameters);
      
      // Get and store performance metrics
      const performance = await this.getPerformance(strategyId);
      await this.storePerformanceInLocalDb(strategyId, performance);
      
      console.log(`Synced full data for strategy ${strategyId}`);
    } catch (error) {
      console.error(`Failed to sync strategy data: ${error}`);
      throw error;
    }
  }
}
