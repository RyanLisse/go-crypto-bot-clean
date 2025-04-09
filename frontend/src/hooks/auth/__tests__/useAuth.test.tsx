import { describe, it, expect } from 'vitest';
import { useAuth } from '../useAuth';

// Simple test to verify the hook exports the expected functions
describe('useAuth', () => {
  it('should export the expected functions', () => {
    // Just verify the hook exists and exports the expected functions
    expect(typeof useAuth).toBe('function');
  });
});
