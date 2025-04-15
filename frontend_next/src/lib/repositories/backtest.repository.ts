import { BaseRepository } from './interfaces';
import { API_CONFIG } from '@/config';

/**
 * Backtest entity
 */
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

/**
 * Backtest trade interface
 */
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

/**
 * Backtest equity point interface
 */
export interface BacktestEquity {
  id: number;
  backtestId: string;
  timestamp: Date;
  equity: number;
  balance: number;
  drawdown: number;
}

/**
 * Implementation of the backtest repository for API communication
 */
export class BacktestRepository implements BaseRepository<Backtest> {
  private readonly endpoint = 'backtests';
  private readonly apiUrl: string;
  
  constructor() {
    this.apiUrl = API_CONFIG.API_URL;
  }
  
  /**
   * Get auth headers with JWT token
   */
  private getAuthHeaders(): HeadersInit {
    // Get token from local storage in client-side environments
    const token = typeof window !== 'undefined' 
      ? localStorage.getItem('token') 
      : null;
    
    return {
      'Content-Type': 'application/json',
      'Authorization': token ? `Bearer ${token}` : '',
    };
  }
  
  /**
   * Get all backtests
   */
  async getAll(): Promise<Backtest[]> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}`, {
      headers: this.getAuthHeaders(),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch backtests: ${response.status}`);
    }
    
    const data = await response.json();
    return this.mapDatesToBacktests(data);
  }
  
  /**
   * Get backtest by ID
   */
  async getById(id: string): Promise<Backtest> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}/${id}`, {
      headers: this.getAuthHeaders(),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch backtest: ${response.status}`);
    }
    
    const data = await response.json();
    return this.mapDatesToBacktest(data);
  }
  
  /**
   * Create a new backtest
   */
  async create(data: Partial<Backtest>): Promise<Backtest> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(data),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to create backtest: ${response.status}`);
    }
    
    const responseData = await response.json();
    return this.mapDatesToBacktest(responseData);
  }
  
  /**
   * Update an existing backtest
   */
  async update(id: string, data: Partial<Backtest>): Promise<Backtest> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}/${id}`, {
      method: 'PUT',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(data),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to update backtest: ${response.status}`);
    }
    
    const responseData = await response.json();
    return this.mapDatesToBacktest(responseData);
  }
  
  /**
   * Delete a backtest
   */
  async delete(id: string): Promise<void> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}/${id}`, {
      method: 'DELETE',
      headers: this.getAuthHeaders(),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to delete backtest: ${response.status}`);
    }
  }
  
  /**
   * Run a new backtest
   */
  async runBacktest(strategyId: string, params: {
    startDate: Date;
    endDate: Date;
    initialBalance: number;
    parameters: Record<string, any>;
  }): Promise<Backtest> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}/run`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify({
        strategyId,
        ...params,
        startDate: params.startDate.toISOString(),
        endDate: params.endDate.toISOString(),
      }),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to run backtest: ${response.status}`);
    }
    
    const responseData = await response.json();
    return this.mapDatesToBacktest(responseData);
  }
  
  /**
   * Get trades for a backtest
   */
  async getTrades(backtestId: string): Promise<BacktestTrade[]> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}/${backtestId}/trades`, {
      headers: this.getAuthHeaders(),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch backtest trades: ${response.status}`);
    }
    
    const data = await response.json();
    return data.map((item: any) => ({
      ...item,
      entryTime: new Date(item.entryTime),
      exitTime: item.exitTime ? new Date(item.exitTime) : undefined,
    }));
  }
  
  /**
   * Get equity curve for a backtest
   */
  async getEquityCurve(backtestId: string): Promise<BacktestEquity[]> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}/${backtestId}/equity`, {
      headers: this.getAuthHeaders(),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch backtest equity: ${response.status}`);
    }
    
    const data = await response.json();
    return data.map((item: any) => ({
      ...item,
      timestamp: new Date(item.timestamp),
    }));
  }
  
  /**
   * Sync full backtest data including trades and equity curve
   */
  async syncFullBacktestData(backtestId: string): Promise<void> {
    try {
      // Get backtest details
      const backtest = await this.getById(backtestId);
      
      // Get and store trades
      await this.getTrades(backtestId);
      
      // Get and store equity curve
      await this.getEquityCurve(backtestId);
      
      console.log(`Synced full data for backtest ${backtestId}`);
    } catch (error) {
      console.error(`Failed to sync backtest data: ${error}`);
      throw error;
    }
  }
  
  /**
   * Convert API dates to Date objects
   */
  private mapDatesToBacktest(data: any): Backtest {
    return {
      ...data,
      startDate: data.startDate ? new Date(data.startDate) : new Date(),
      endDate: data.endDate ? new Date(data.endDate) : new Date(),
      createdAt: data.createdAt ? new Date(data.createdAt) : new Date(),
      updatedAt: data.updatedAt ? new Date(data.updatedAt) : new Date(),
    };
  }
  
  /**
   * Map dates for an array of backtests
   */
  private mapDatesToBacktests(data: any[]): Backtest[] {
    return data.map(item => this.mapDatesToBacktest(item));
  }
} 