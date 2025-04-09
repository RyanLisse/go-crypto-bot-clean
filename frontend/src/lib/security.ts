/**
 * Security utilities for AI integration
 */

/**
 * Sanitize user input to prevent prompt injection
 * @param input User input to sanitize
 * @returns Sanitized input
 */
export function sanitizeUserInput(input: string): string {
  if (!input) return '';
  
  // Remove potentially harmful characters and sequences
  let sanitized = input
    .replace(/\\n/g, ' ')
    .replace(/\\r/g, ' ')
    .replace(/\\t/g, ' ');
  
  // Limit input length
  const maxInputLength = 1000;
  if (sanitized.length > maxInputLength) {
    sanitized = sanitized.substring(0, maxInputLength);
  }
  
  return sanitized;
}

/**
 * Validate a trade request
 * @param symbol Trading symbol (e.g., BTC)
 * @param amount Trade amount
 * @param orderType Order type (market or limit)
 * @param price Optional price for limit orders
 * @returns Error message if validation fails, null if valid
 */
export function validateTradeRequest(
  symbol: string,
  amount: number,
  orderType: 'market' | 'limit',
  price?: number
): string | null {
  if (!symbol) {
    return 'Symbol is required';
  }
  
  if (!amount || amount <= 0) {
    return 'Amount must be positive';
  }
  
  if (orderType !== 'market' && orderType !== 'limit') {
    return 'Order type must be market or limit';
  }
  
  if (orderType === 'limit' && (!price || price <= 0)) {
    return 'Price must be positive for limit orders';
  }
  
  return null;
}

/**
 * Simple rate limiter for client-side rate limiting
 */
export class RateLimiter {
  private requestCount: number = 0;
  private lastReset: Date = new Date();
  private readonly maxRequests: number;
  private readonly resetPeriodMs: number;
  
  /**
   * Create a new rate limiter
   * @param maxRequests Maximum number of requests allowed in the reset period
   * @param resetPeriodMs Reset period in milliseconds
   */
  constructor(maxRequests: number = 10, resetPeriodMs: number = 60000) {
    this.maxRequests = maxRequests;
    this.resetPeriodMs = resetPeriodMs;
  }
  
  /**
   * Check if a request is allowed
   * @returns True if the request is allowed, false otherwise
   */
  public allow(): boolean {
    const now = new Date();
    if (now.getTime() - this.lastReset.getTime() >= this.resetPeriodMs) {
      this.requestCount = 0;
      this.lastReset = now;
    }
    
    if (this.requestCount >= this.maxRequests) {
      return false;
    }
    
    this.requestCount++;
    return true;
  }
  
  /**
   * Get the current request count
   * @returns Current request count
   */
  public getRequestCount(): number {
    return this.requestCount;
  }
  
  /**
   * Get the time until the next reset in milliseconds
   * @returns Time until the next reset in milliseconds
   */
  public getTimeUntilReset(): number {
    const now = new Date();
    const elapsed = now.getTime() - this.lastReset.getTime();
    return Math.max(0, this.resetPeriodMs - elapsed);
  }
}

/**
 * Circuit breaker for preventing excessive API calls during errors
 */
export class CircuitBreaker {
  private failures: number = 0;
  private lastFailure: Date | null = null;
  private state: 'CLOSED' | 'OPEN' | 'HALF_OPEN' = 'CLOSED';
  private readonly failureThreshold: number;
  private readonly resetTimeoutMs: number;
  
  /**
   * Create a new circuit breaker
   * @param failureThreshold Number of failures before opening the circuit
   * @param resetTimeoutMs Time in milliseconds before trying to close the circuit again
   */
  constructor(failureThreshold: number = 3, resetTimeoutMs: number = 30000) {
    this.failureThreshold = failureThreshold;
    this.resetTimeoutMs = resetTimeoutMs;
  }
  
  /**
   * Check if a request is allowed
   * @returns True if the request is allowed, false otherwise
   */
  public isAllowed(): boolean {
    if (this.state === 'OPEN') {
      // Check if it's time to try again
      if (this.lastFailure && new Date().getTime() - this.lastFailure.getTime() >= this.resetTimeoutMs) {
        this.state = 'HALF_OPEN';
        return true;
      }
      return false;
    }
    
    return true;
  }
  
  /**
   * Record a successful request
   */
  public recordSuccess(): void {
    this.failures = 0;
    this.state = 'CLOSED';
  }
  
  /**
   * Record a failed request
   */
  public recordFailure(): void {
    this.failures++;
    this.lastFailure = new Date();
    
    if (this.state === 'HALF_OPEN' || this.failures >= this.failureThreshold) {
      this.state = 'OPEN';
    }
  }
  
  /**
   * Get the current state of the circuit breaker
   * @returns Current state
   */
  public getState(): 'CLOSED' | 'OPEN' | 'HALF_OPEN' {
    return this.state;
  }
}

// Create singleton instances for global use
export const globalRateLimiter = new RateLimiter(
  Number(import.meta.env.VITE_MAX_AI_REQUESTS_PER_MINUTE || 10),
  60000
);

export const globalCircuitBreaker = new CircuitBreaker(
  Number(import.meta.env.VITE_AI_FAILURE_THRESHOLD || 3),
  Number(import.meta.env.VITE_AI_RESET_TIMEOUT_MS || 30000)
);

export default {
  sanitizeUserInput,
  validateTradeRequest,
  RateLimiter,
  CircuitBreaker,
  globalRateLimiter,
  globalCircuitBreaker,
};
