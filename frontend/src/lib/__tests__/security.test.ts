import { describe, it, expect, beforeEach, vi } from 'vitest';
import { 
  sanitizeUserInput, 
  validateTradeRequest, 
  RateLimiter, 
  CircuitBreaker 
} from '../security';

describe('sanitizeUserInput', () => {
  it('should return empty string for null or undefined input', () => {
    expect(sanitizeUserInput('')).toBe('');
    expect(sanitizeUserInput(undefined as any)).toBe('');
    expect(sanitizeUserInput(null as any)).toBe('');
  });
  
  it('should remove potentially harmful characters', () => {
    expect(sanitizeUserInput('Hello\\nWorld')).toBe('Hello World');
    expect(sanitizeUserInput('Hello\\rWorld')).toBe('Hello World');
    expect(sanitizeUserInput('Hello\\tWorld')).toBe('Hello World');
  });
  
  it('should limit input length', () => {
    const longInput = 'a'.repeat(2000);
    const sanitized = sanitizeUserInput(longInput);
    expect(sanitized.length).toBeLessThanOrEqual(1000);
  });
});

describe('validateTradeRequest', () => {
  it('should validate a valid market order', () => {
    const result = validateTradeRequest('BTC', 0.1, 'market');
    expect(result).toBeNull();
  });
  
  it('should validate a valid limit order', () => {
    const result = validateTradeRequest('ETH', 1.0, 'limit', 2000);
    expect(result).toBeNull();
  });
  
  it('should reject missing symbol', () => {
    const result = validateTradeRequest('', 0.1, 'market');
    expect(result).toBe('Symbol is required');
  });
  
  it('should reject invalid amount', () => {
    const result = validateTradeRequest('BTC', 0, 'market');
    expect(result).toBe('Amount must be positive');
    
    const result2 = validateTradeRequest('BTC', -1, 'market');
    expect(result2).toBe('Amount must be positive');
  });
  
  it('should reject invalid order type', () => {
    const result = validateTradeRequest('BTC', 0.1, 'invalid' as any);
    expect(result).toBe('Order type must be market or limit');
  });
  
  it('should reject limit order without price', () => {
    const result = validateTradeRequest('BTC', 0.1, 'limit');
    expect(result).toBe('Price must be positive for limit orders');
    
    const result2 = validateTradeRequest('BTC', 0.1, 'limit', 0);
    expect(result2).toBe('Price must be positive for limit orders');
    
    const result3 = validateTradeRequest('BTC', 0.1, 'limit', -1);
    expect(result3).toBe('Price must be positive for limit orders');
  });
});

describe('RateLimiter', () => {
  let rateLimiter: RateLimiter;
  
  beforeEach(() => {
    rateLimiter = new RateLimiter(3, 1000);
    vi.useFakeTimers();
  });
  
  it('should allow requests within the limit', () => {
    expect(rateLimiter.allow()).toBe(true);
    expect(rateLimiter.allow()).toBe(true);
    expect(rateLimiter.allow()).toBe(true);
  });
  
  it('should reject requests over the limit', () => {
    expect(rateLimiter.allow()).toBe(true);
    expect(rateLimiter.allow()).toBe(true);
    expect(rateLimiter.allow()).toBe(true);
    expect(rateLimiter.allow()).toBe(false);
  });
  
  it('should reset after the reset period', () => {
    expect(rateLimiter.allow()).toBe(true);
    expect(rateLimiter.allow()).toBe(true);
    expect(rateLimiter.allow()).toBe(true);
    expect(rateLimiter.allow()).toBe(false);
    
    // Advance time by reset period
    vi.advanceTimersByTime(1000);
    
    // Should be allowed again
    expect(rateLimiter.allow()).toBe(true);
  });
  
  it('should return correct request count', () => {
    expect(rateLimiter.getRequestCount()).toBe(0);
    rateLimiter.allow();
    expect(rateLimiter.getRequestCount()).toBe(1);
    rateLimiter.allow();
    expect(rateLimiter.getRequestCount()).toBe(2);
  });
  
  it('should return correct time until reset', () => {
    rateLimiter.allow();
    
    // Advance time by 500ms
    vi.advanceTimersByTime(500);
    
    // Should have 500ms remaining
    expect(rateLimiter.getTimeUntilReset()).toBe(500);
    
    // Advance time by another 600ms (past reset)
    vi.advanceTimersByTime(600);
    
    // Should be 0 (reset already happened)
    expect(rateLimiter.getTimeUntilReset()).toBe(0);
  });
});

describe('CircuitBreaker', () => {
  let circuitBreaker: CircuitBreaker;
  
  beforeEach(() => {
    circuitBreaker = new CircuitBreaker(3, 1000);
    vi.useFakeTimers();
  });
  
  it('should start in closed state', () => {
    expect(circuitBreaker.getState()).toBe('CLOSED');
    expect(circuitBreaker.isAllowed()).toBe(true);
  });
  
  it('should open after threshold failures', () => {
    expect(circuitBreaker.isAllowed()).toBe(true);
    
    circuitBreaker.recordFailure();
    expect(circuitBreaker.getState()).toBe('CLOSED');
    
    circuitBreaker.recordFailure();
    expect(circuitBreaker.getState()).toBe('CLOSED');
    
    circuitBreaker.recordFailure();
    expect(circuitBreaker.getState()).toBe('OPEN');
    
    // Should reject requests in open state
    expect(circuitBreaker.isAllowed()).toBe(false);
  });
  
  it('should transition to half-open after reset timeout', () => {
    circuitBreaker.recordFailure();
    circuitBreaker.recordFailure();
    circuitBreaker.recordFailure();
    
    expect(circuitBreaker.getState()).toBe('OPEN');
    expect(circuitBreaker.isAllowed()).toBe(false);
    
    // Advance time by reset timeout
    vi.advanceTimersByTime(1000);
    
    // Should be in half-open state and allow one request
    expect(circuitBreaker.isAllowed()).toBe(true);
    expect(circuitBreaker.getState()).toBe('HALF_OPEN');
  });
  
  it('should close after successful request in half-open state', () => {
    circuitBreaker.recordFailure();
    circuitBreaker.recordFailure();
    circuitBreaker.recordFailure();
    
    expect(circuitBreaker.getState()).toBe('OPEN');
    
    // Advance time by reset timeout
    vi.advanceTimersByTime(1000);
    
    // Should be in half-open state
    expect(circuitBreaker.isAllowed()).toBe(true);
    expect(circuitBreaker.getState()).toBe('HALF_OPEN');
    
    // Record success
    circuitBreaker.recordSuccess();
    
    // Should be closed again
    expect(circuitBreaker.getState()).toBe('CLOSED');
    expect(circuitBreaker.isAllowed()).toBe(true);
  });
  
  it('should open immediately after failure in half-open state', () => {
    circuitBreaker.recordFailure();
    circuitBreaker.recordFailure();
    circuitBreaker.recordFailure();
    
    // Advance time by reset timeout
    vi.advanceTimersByTime(1000);
    
    // Should be in half-open state
    expect(circuitBreaker.isAllowed()).toBe(true);
    expect(circuitBreaker.getState()).toBe('HALF_OPEN');
    
    // Record failure
    circuitBreaker.recordFailure();
    
    // Should be open again
    expect(circuitBreaker.getState()).toBe('OPEN');
    expect(circuitBreaker.isAllowed()).toBe(false);
  });
});
