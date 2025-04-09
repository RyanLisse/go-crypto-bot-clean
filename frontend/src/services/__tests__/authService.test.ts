import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { authService } from '../authService';

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {};
  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => {
      store[key] = value.toString();
    },
    removeItem: (key: string) => {
      delete store[key];
    },
    clear: () => {
      store = {};
    },
  };
})();

// Mock fetch
global.fetch = vi.fn();

// Mock localStorage globally
global.localStorage = localStorageMock;

describe('Authentication Service', () => {
  beforeEach(() => {
    // Clear localStorage before each test
    localStorage.clear();

    // Reset fetch mock
    vi.resetAllMocks?.() || vi.clearAllMocks();
  });

  afterEach(() => {
    vi.clearAllMocks?.() || vi.resetAllMocks?.();
  });

  describe('login', () => {
    it('should store token in localStorage on successful login', async () => {
      // Mock successful login response
      const mockResponse = {
        token: 'test-token',
        expires_at: '2023-12-31T23:59:59Z',
        user_id: 'test-user',
        role: 'user',
      };

      (global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      });

      // Call login method
      const result = await authService.login('testuser', 'password123');

      // Verify fetch was called with correct arguments
      expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: 'testuser',
          password: 'password123',
        }),
      });

      // Verify token was stored in localStorage
      expect(localStorage.getItem('token')).toBe(mockResponse.token);
      expect(localStorage.getItem('user')).toBe(JSON.stringify(mockResponse));

      // Verify function returns user data
      expect(result).toEqual(mockResponse);
    });

    it('should throw an error on failed login', async () => {
      // Mock failed login response
      (global.fetch as any).mockResolvedValueOnce({
        ok: false,
        status: 401,
        statusText: 'Unauthorized',
        json: async () => ({ message: 'Invalid credentials' }),
      });

      // Call login method and expect it to throw
      await expect(authService.login('testuser', 'wrongpassword')).rejects.toThrow('Login failed: 401 Unauthorized');

      // Verify localStorage was not updated
      expect(localStorage.getItem('token')).toBeNull();
      expect(localStorage.getItem('user')).toBeNull();
    });
  });

  describe('logout', () => {
    it('should remove token from localStorage on logout', async () => {
      // Setup localStorage with token
      localStorage.setItem('token', 'test-token');
      localStorage.setItem('user', JSON.stringify({ user_id: 'test-user' }));

      // Mock successful logout response
      (global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ message: 'Successfully logged out' }),
      });

      // Call logout method
      await authService.logout();

      // Verify fetch was called with correct arguments
      expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/auth/logout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer test-token',
        },
      });

      // Verify localStorage was cleared
      expect(localStorage.getItem('token')).toBeNull();
      expect(localStorage.getItem('user')).toBeNull();
    });

    it('should clear localStorage even if logout API call fails', async () => {
      // Setup localStorage with token
      localStorage.setItem('token', 'test-token');
      localStorage.setItem('user', JSON.stringify({ user_id: 'test-user' }));

      // Mock failed logout response
      (global.fetch as any).mockRejectedValueOnce(new Error('Network error'));

      // Call logout method
      await authService.logout();

      // Verify localStorage was still cleared
      expect(localStorage.getItem('token')).toBeNull();
      expect(localStorage.getItem('user')).toBeNull();
    });
  });

  describe('getCurrentUser', () => {
    it('should return user data from localStorage', () => {
      // Setup localStorage with user data
      const userData = { user_id: 'test-user', role: 'user' };
      localStorage.setItem('user', JSON.stringify(userData));

      // Call getCurrentUser method
      const result = authService.getCurrentUser();

      // Verify result
      expect(result).toEqual(userData);
    });

    it('should return null if no user data in localStorage', () => {
      // Call getCurrentUser method with empty localStorage
      const result = authService.getCurrentUser();

      // Verify result
      expect(result).toBeNull();
    });
  });

  describe('isAuthenticated', () => {
    it('should return true if token exists in localStorage', () => {
      // Setup localStorage with token
      localStorage.setItem('token', 'test-token');

      // Call isAuthenticated method
      const result = authService.isAuthenticated();

      // Verify result
      expect(result).toBe(true);
    });

    it('should return false if no token in localStorage', () => {
      // Call isAuthenticated method with empty localStorage
      const result = authService.isAuthenticated();

      // Verify result
      expect(result).toBe(false);
    });
  });

  describe('getAuthHeader', () => {
    it('should return authorization header with token', () => {
      // Setup localStorage with token
      localStorage.setItem('token', 'test-token');

      // Call getAuthHeader method
      const result = authService.getAuthHeader();

      // Verify result
      expect(result).toEqual({ Authorization: 'Bearer test-token' });
    });

    it('should return empty object if no token in localStorage', () => {
      // Call getAuthHeader method with empty localStorage
      const result = authService.getAuthHeader();

      // Verify result
      expect(result).toEqual({});
    });
  });
});
