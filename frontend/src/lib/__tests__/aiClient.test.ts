import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { sendChatMessage, executeTradingFunction, streamChatMessage, getAIMetrics } from '../aiClient';
import { globalRateLimiter, globalCircuitBreaker } from '../security';

// Mock fetch
const mockFetch = vi.fn();
global.fetch = mockFetch;

// Mock localStorage
const mockLocalStorage = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
  length: 0,
  key: vi.fn(),
};
Object.defineProperty(window, 'localStorage', { value: mockLocalStorage });

// Mock performance.now
const mockPerformanceNow = vi.fn();
performance.now = mockPerformanceNow;

describe('AI Client', () => {
  beforeEach(() => {
    vi.resetAllMocks();
    mockLocalStorage.getItem.mockReturnValue('mock-token');
    mockPerformanceNow.mockReturnValueOnce(0).mockReturnValueOnce(100);
    
    // Reset rate limiter and circuit breaker
    vi.spyOn(globalRateLimiter, 'allow').mockReturnValue(true);
    vi.spyOn(globalCircuitBreaker, 'isAllowed').mockReturnValue(true);
    vi.spyOn(globalCircuitBreaker, 'recordSuccess').mockImplementation(() => {});
    vi.spyOn(globalCircuitBreaker, 'recordFailure').mockImplementation(() => {});
  });
  
  afterEach(() => {
    vi.restoreAllMocks();
  });
  
  describe('sendChatMessage', () => {
    it('should send a chat message and return the response', async () => {
      const mockResponse = {
        message: { role: 'assistant', content: 'Hello!' },
        session_id: 'test-session',
      };
      
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      });
      
      const result = await sendChatMessage('Hi there', 'test-session');
      
      expect(mockFetch).toHaveBeenCalledWith('/api/chat', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer mock-token',
        },
        body: JSON.stringify({
          messages: [{ role: 'user', content: 'Hi there' }],
          session_id: 'test-session',
        }),
      });
      
      expect(result).toEqual(mockResponse);
      expect(globalCircuitBreaker.recordSuccess).toHaveBeenCalled();
    });
    
    it('should throw an error if not authenticated', async () => {
      mockLocalStorage.getItem.mockReturnValueOnce(null);
      
      await expect(sendChatMessage('Hi there')).rejects.toThrow('Authentication required');
    });
    
    it('should throw an error if the API returns an error', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
      });
      
      await expect(sendChatMessage('Hi there')).rejects.toThrow('API error: 500');
      expect(globalCircuitBreaker.recordFailure).toHaveBeenCalled();
    });
    
    it('should throw an error if rate limited', async () => {
      vi.spyOn(globalRateLimiter, 'allow').mockReturnValueOnce(false);
      
      await expect(sendChatMessage('Hi there')).rejects.toThrow('Rate limit exceeded');
    });
    
    it('should throw an error if circuit breaker is open', async () => {
      vi.spyOn(globalCircuitBreaker, 'isAllowed').mockReturnValueOnce(false);
      
      await expect(sendChatMessage('Hi there')).rejects.toThrow('Service temporarily unavailable');
    });
  });
  
  describe('executeTradingFunction', () => {
    it('should execute a trading function and return the result', async () => {
      const mockResponse = {
        result: { success: true },
      };
      
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      });
      
      const result = await executeTradingFunction('buyBTC', { amount: 0.1 });
      
      expect(mockFetch).toHaveBeenCalledWith('/api/function', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer mock-token',
        },
        body: JSON.stringify({
          function_name: 'buyBTC',
          parameters: { amount: 0.1 },
        }),
      });
      
      expect(result).toEqual(mockResponse);
      expect(globalCircuitBreaker.recordSuccess).toHaveBeenCalled();
    });
  });
  
  describe('streamChatMessage', () => {
    it('should stream a chat message and return the stream', async () => {
      const mockStream = new ReadableStream();
      
      mockFetch.mockResolvedValueOnce({
        ok: true,
        body: mockStream,
      });
      
      const result = await streamChatMessage('Hi there', 'test-session');
      
      expect(mockFetch).toHaveBeenCalledWith('/api/chat/stream', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer mock-token',
        },
        body: JSON.stringify({
          messages: [{ role: 'user', content: 'Hi there' }],
          session_id: 'test-session',
        }),
      });
      
      expect(result).toBe(mockStream);
      expect(globalCircuitBreaker.recordSuccess).toHaveBeenCalled();
    });
  });
  
  describe('getAIMetrics', () => {
    it('should return AI metrics', () => {
      // Call sendChatMessage to record a request
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({}),
      });
      
      sendChatMessage('Hi there');
      
      const metrics = getAIMetrics();
      
      expect(metrics).toBeDefined();
      expect(typeof metrics).toBe('object');
    });
  });
});
