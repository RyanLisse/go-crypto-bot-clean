import { BaseRepository } from './interfaces';
import { API_CONFIG } from '@/config';

/**
 * Strategy entity
 */
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

/**
 * Strategy parameter interface
 */
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

/**
 * Strategy performance interface
 */
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

/**
 * Implementation of the strategy repository for API communication
 */
export class StrategyRepository implements BaseRepository<Strategy> {
  private readonly endpoint = 'strategies';
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
   * Get all strategies
   */
  async getAll(): Promise<Strategy[]> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}`, {
      headers: this.getAuthHeaders(),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch strategies: ${response.status}`);
    }
    
    const data = await response.json();
    return this.mapDatesToStrategies(data);
  }
  
  /**
   * Get strategy by ID
   */
  async getById(id: string): Promise<Strategy> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}/${id}`, {
      headers: this.getAuthHeaders(),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch strategy: ${response.status}`);
    }
    
    const data = await response.json();
    return this.mapDatesToStrategy(data);
  }
  
  /**
   * Create a new strategy
   */
  async create(data: Partial<Strategy>): Promise<Strategy> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(data),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to create strategy: ${response.status}`);
    }
    
    const responseData = await response.json();
    return this.mapDatesToStrategy(responseData);
  }
  
  /**
   * Update an existing strategy
   */
  async update(id: string, data: Partial<Strategy>): Promise<Strategy> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}/${id}`, {
      method: 'PUT',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(data),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to update strategy: ${response.status}`);
    }
    
    const responseData = await response.json();
    return this.mapDatesToStrategy(responseData);
  }
  
  /**
   * Delete a strategy
   */
  async delete(id: string): Promise<void> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}/${id}`, {
      method: 'DELETE',
      headers: this.getAuthHeaders(),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to delete strategy: ${response.status}`);
    }
  }
  
  /**
   * Get parameters for a strategy
   */
  async getParameters(strategyId: string): Promise<StrategyParameter[]> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}/${strategyId}/parameters`, {
      headers: this.getAuthHeaders(),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch strategy parameters: ${response.status}`);
    }
    
    return response.json();
  }
  
  /**
   * Get performance metrics for a strategy
   */
  async getPerformance(strategyId: string): Promise<StrategyPerformance[]> {
    const response = await fetch(`${this.apiUrl}/${this.endpoint}/${strategyId}/performance`, {
      headers: this.getAuthHeaders(),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch strategy performance: ${response.status}`);
    }
    
    const data = await response.json();
    return data.map((item: any) => ({
      ...item,
      periodStart: new Date(item.periodStart),
      periodEnd: new Date(item.periodEnd),
    }));
  }
  
  /**
   * Sync full strategy data including parameters and performance
   */
  async syncFullStrategyData(strategyId: string): Promise<void> {
    try {
      // Get strategy details
      const strategy = await this.getById(strategyId);
      
      // Get and store parameters
      await this.getParameters(strategyId);
      
      // Get and store performance metrics
      await this.getPerformance(strategyId);
      
      console.log(`Synced full data for strategy ${strategyId}`);
    } catch (error) {
      console.error(`Failed to sync strategy data: ${error}`);
      throw error;
    }
  }
  
  /**
   * Convert API dates to Date objects
   */
  private mapDatesToStrategy(data: any): Strategy {
    return {
      ...data,
      createdAt: data.createdAt ? new Date(data.createdAt) : new Date(),
      updatedAt: data.updatedAt ? new Date(data.updatedAt) : new Date(),
    };
  }
  
  /**
   * Map dates for an array of strategies
   */
  private mapDatesToStrategies(data: any[]): Strategy[] {
    return data.map(item => this.mapDatesToStrategy(item));
  }
} 