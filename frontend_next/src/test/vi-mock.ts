import { vi } from 'vitest';

// Jest-style mock utility for Vitest
// This helps tests that use jest.mock() syntax to work with Vitest

/**
 * Enhanced jest mock compatibility layer
 * Provides a closer approximation to Jest's module mocking behavior
 */

// Store original mock implementations
const originalMockFn = vi.fn;
const originalMock = vi.mock;
const originalSpyOn = vi.spyOn;

// Enhanced mock function with Jest compatibility
const enhancedMockFn = (...args: any[]) => {
  const mockFn = originalMockFn(...args);
  
  // Add Jest-specific mock properties
  if (!mockFn.mockName) {
    mockFn.mockName = (name: string) => {
      (mockFn as any)._mockName = name;
      return mockFn;
    };
  }
  
  return mockFn;
};

// Enhanced mock with better Jest compatibility
const enhancedMock = (moduleName: string, factory?: () => any) => {
  try {
    return originalMock(moduleName, factory);
  } catch (error) {
    console.warn(`Mock for "${moduleName}" failed, using fallback`, error);
    // Fallback to vi.mock without hoisting
    // @ts-ignore - Working around type constraints for compatibility
    return originalMock.call(vi, moduleName, () => {
      return factory ? factory() : {};
    });
  }
};

// Export a jest-like object that maps to vitest functions
const jest = {
  fn: enhancedMockFn,
  mock: enhancedMock,
  spyOn: originalSpyOn,
  clearAllMocks: vi.clearAllMocks,
  resetAllMocks: vi.resetAllMocks,
  
  // Add commonly used Jest functions
  resetModules: vi.resetModules,
  useFakeTimers: vi.useFakeTimers,
  useRealTimers: vi.useRealTimers,
  runAllTimers: vi.runAllTimers,
  advanceTimersByTime: vi.advanceTimersByTime,
};

// Add jest to global to mimic Jest's environment
(global as any).jest = jest;

// Log success message
console.log('Jest compatibility layer initialized');

export { jest }; 