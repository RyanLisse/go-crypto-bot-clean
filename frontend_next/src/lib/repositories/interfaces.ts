/**
 * Base repository interface defining common CRUD operations
 */
export interface BaseRepository<T> {
  /**
   * Get all entities
   */
  getAll(): Promise<T[]>;
  
  /**
   * Get entity by ID
   * @param id Entity ID
   */
  getById(id: string): Promise<T>;
  
  /**
   * Create a new entity
   * @param data Entity data to create
   */
  create(data: Partial<T>): Promise<T>;
  
  /**
   * Update an existing entity
   * @param id Entity ID
   * @param data Entity data to update
   */
  update(id: string, data: Partial<T>): Promise<T>;
  
  /**
   * Delete an entity
   * @param id Entity ID
   */
  delete(id: string): Promise<void>;
  
  /**
   * Sync data from remote API to local storage
   */
  syncFromApi?(): Promise<void>;
} 