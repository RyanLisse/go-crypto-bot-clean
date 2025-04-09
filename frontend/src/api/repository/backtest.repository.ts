import { BaseRepository } from './base';
import { db } from '../../db/client';
import { backtests, backtestTrades, backtestEquity } from '../../db/schema/backtests';
import { eq } from 'drizzle-orm';

// Backtest interface matching the backend model
export interface Backtest {
  id: string;
  userId: string;
  strategyId: string;
  name: string;
  description: string;
  startDate: Date;
  endDate: Date;
  initialBalance: number;
  finalBalance: number;
  totalTrades: number;
  winningTrades: number;
  losingTrades: number;
  winRate: number;
  profitFactor: number;
  sharpeRatio: number;
  maxDrawdown: number;
  parameters: Record<string, any>;
  status: string;
  createdAt: Date;
  updatedAt: Date;
}

// Backtest trade interface
export interface BacktestTrade {
  id: number;
  backtestId: string;
  symbol: string;
  entryTime: Date;
  entryPrice: number;
  exitTime?: Date;
  exitPrice?: number;
  quantity: number;
  direction: string;
  profitLoss?: number;
  profitLossPct?: number;
  exitReason?: string;
}

// Backtest equity point interface
export interface BacktestEquity {
  id: number;
  backtestId: string;
  timestamp: Date;
  equity: number;
  balance: number;
  drawdown: number;
}

export class BacktestRepository extends BaseRepository<Backtest> {
  protected endpoint = 'backtests';
  protected tableName = 'backtests';
  
  // Store backtests in local database
  protected async storeInLocalDb(items: Backtest[]): Promise<void> {
    // Use a transaction to ensure all operations succeed or fail together
    await db.transaction(async (tx) => {
      // For each backtest
      for (const backtest of items) {
        // Insert or update backtest
        await tx
          .insert(backtests)
          .values({
            id: backtest.id,
            userId: backtest.userId,
            strategyId: backtest.strategyId,
            name: backtest.name,
            description: backtest.description,
            startDate: new Date(backtest.startDate),
            endDate: new Date(backtest.endDate),
            initialBalance: backtest.initialBalance,
            finalBalance: backtest.finalBalance,
            totalTrades: backtest.totalTrades,
            winningTrades: backtest.winningTrades,
            losingTrades: backtest.losingTrades,
            winRate: backtest.winRate,
            profitFactor: backtest.profitFactor,
            sharpeRatio: backtest.sharpeRatio,
            maxDrawdown: backtest.maxDrawdown,
            parameters: backtest.parameters,
            status: backtest.status,
            createdAt: new Date(backtest.createdAt),
            updatedAt: new Date(backtest.updatedAt),
          })
          .onConflictDoUpdate({
            target: backtests.id,
            set: {
              name: backtest.name,
              description: backtest.description,
              finalBalance: backtest.finalBalance,
              totalTrades: backtest.totalTrades,
              winningTrades: backtest.winningTrades,
              losingTrades: backtest.losingTrades,
              winRate: backtest.winRate,
              profitFactor: backtest.profitFactor,
              sharpeRatio: backtest.sharpeRatio,
              maxDrawdown: backtest.maxDrawdown,
              status: backtest.status,
              updatedAt: new Date(backtest.updatedAt),
            },
          });
      }
    });
  }
  
  // Run a new backtest
  async runBacktest(strategyId: string, params: {
    startDate: Date;
    endDate: Date;
    initialBalance: number;
    parameters: Record<string, any>;
  }): Promise<Backtest> {
    return this.fetchApi<Backtest>(this.getApiUrl('/run'), {
      method: 'POST',
      body: JSON.stringify({
        strategyId,
        ...params,
      }),
    });
  }
  
  // Get trades for a backtest
  async getTrades(backtestId: string): Promise<BacktestTrade[]> {
    return this.fetchApi<BacktestTrade[]>(this.getApiUrl(`/${backtestId}/trades`));
  }
  
  // Store trades in local database
  async storeTradesInLocalDb(backtestId: string, trades: BacktestTrade[]): Promise<void> {
    await db.transaction(async (tx) => {
      // Delete existing trades
      await tx.delete(backtestTrades).where(eq(backtestTrades.backtestId, backtestId));
      
      // Insert new trades
      for (const trade of trades) {
        await tx.insert(backtestTrades).values({
          backtestId: trade.backtestId,
          symbol: trade.symbol,
          entryTime: new Date(trade.entryTime),
          entryPrice: trade.entryPrice,
          exitTime: trade.exitTime ? new Date(trade.exitTime) : undefined,
          exitPrice: trade.exitPrice,
          quantity: trade.quantity,
          direction: trade.direction,
          profitLoss: trade.profitLoss,
          profitLossPct: trade.profitLossPct,
          exitReason: trade.exitReason,
        });
      }
    });
  }
  
  // Get equity curve for a backtest
  async getEquityCurve(backtestId: string): Promise<BacktestEquity[]> {
    return this.fetchApi<BacktestEquity[]>(this.getApiUrl(`/${backtestId}/equity`));
  }
  
  // Store equity curve in local database
  async storeEquityCurveInLocalDb(backtestId: string, equity: BacktestEquity[]): Promise<void> {
    await db.transaction(async (tx) => {
      // Delete existing equity points
      await tx.delete(backtestEquity).where(eq(backtestEquity.backtestId, backtestId));
      
      // Insert new equity points
      for (const point of equity) {
        await tx.insert(backtestEquity).values({
          backtestId: point.backtestId,
          timestamp: new Date(point.timestamp),
          equity: point.equity,
          balance: point.balance,
          drawdown: point.drawdown,
        });
      }
    });
  }
  
  // Sync all backtest data (including trades and equity curve)
  async syncFullBacktestData(backtestId: string): Promise<void> {
    try {
      // Get backtest details
      const backtest = await this.getById(backtestId);
      await this.storeInLocalDb([backtest]);
      
      // Get and store trades
      const trades = await this.getTrades(backtestId);
      await this.storeTradesInLocalDb(backtestId, trades);
      
      // Get and store equity curve
      const equity = await this.getEquityCurve(backtestId);
      await this.storeEquityCurveInLocalDb(backtestId, equity);
      
      console.log(`Synced full data for backtest ${backtestId}`);
    } catch (error) {
      console.error(`Failed to sync backtest data: ${error}`);
      throw error;
    }
  }
}
