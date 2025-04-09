import { db } from '../../db/client';

// Base API URL from environment variables
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

// Base repository class with common methods
export abstract class BaseRepository<T> {
  protected abstract endpoint: string;
  protected abstract tableName: string;
  
  // Get the full API URL for the endpoint
  protected getApiUrl(path: string = ''): string {
    return `${API_BASE_URL}/${this.endpoint}${path}`;
  }
  
  // Get auth headers with JWT token
  protected getAuthHeaders(): HeadersInit {
    // Get token from local storage or auth provider
    const token = localStorage.getItem('token');
    
    return {
      'Content-Type': 'application/json',
      'Authorization': token ? `Bearer ${token}` : '',
    };
  }
  
  // Generic fetch method with error handling
  protected async fetchApi<R>(
    url: string, 
    options: RequestInit = {}
  ): Promise<R> {
    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          ...this.getAuthHeaders(),
          ...options.headers,
        },
      });
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `API error: ${response.status}`);
      }
      
      return await response.json();
    } catch (error) {
      console.error(`API request failed: ${error}`);
      throw error;
    }
  }
  
  // Generic CRUD methods
  
  // Get all items
  async getAll(): Promise<T[]> {
    return this.fetchApi<T[]>(this.getApiUrl());
  }
  
  // Get item by ID
  async getById(id: string): Promise<T> {
    return this.fetchApi<T>(this.getApiUrl(`/${id}`));
  }
  
  // Create new item
  async create(data: Partial<T>): Promise<T> {
    return this.fetchApi<T>(this.getApiUrl(), {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }
  
  // Update item
  async update(id: string, data: Partial<T>): Promise<T> {
    return this.fetchApi<T>(this.getApiUrl(`/${id}`), {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }
  
  // Delete item
  async delete(id: string): Promise<void> {
    return this.fetchApi<void>(this.getApiUrl(`/${id}`), {
      method: 'DELETE',
    });
  }
  
  // Sync data from API to local database
  async syncFromApi(): Promise<void> {
    try {
      // Get data from API
      const items = await this.getAll();
      
      // Store in local database
      // This is a simplified example - in a real app, you'd need to handle
      // more complex syncing logic with conflict resolution
      await this.storeInLocalDb(items);
      
      console.log(`Synced ${items.length} items from API to local database`);
    } catch (error) {
      console.error(`Failed to sync from API: ${error}`);
      throw error;
    }
  }
  
  // Abstract method to be implemented by subclasses
  protected abstract storeInLocalDb(items: T[]): Promise<void>;
}
