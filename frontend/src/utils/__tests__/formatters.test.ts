import { formatNumber } from '../formatters';

describe('formatNumber', () => {
  it('formats numbers with default 2 decimal places', () => {
    expect(formatNumber(123.456)).toBe('123.46');
    expect(formatNumber(123)).toBe('123.00');
  });

  it('formats numbers with specified decimal places', () => {
    expect(formatNumber(123.456, 3)).toBe('123.456');
    expect(formatNumber(123, 4)).toBe('123.0000');
  });

  it('handles zero correctly', () => {
    expect(formatNumber(0)).toBe('0.00');
    expect(formatNumber(0, 3)).toBe('0.000');
  });

  it('handles negative numbers correctly', () => {
    expect(formatNumber(-123.456)).toBe('-123.46');
    expect(formatNumber(-123.456, 3)).toBe('-123.456');
  });

  it('formats large numbers with thousands separators', () => {
    expect(formatNumber(1234567.89)).toBe('1,234,567.89');
    expect(formatNumber(1000000)).toBe('1,000,000.00');
  });

  it('handles very small decimal numbers', () => {
    expect(formatNumber(0.0001, 4)).toBe('0.0001');
    expect(formatNumber(0.00001, 5)).toBe('0.00001');
  });

  it('rounds numbers according to specified decimal places', () => {
    expect(formatNumber(123.456, 1)).toBe('123.5');
    expect(formatNumber(123.444, 1)).toBe('123.4');
  });

  it('uses en-US locale formatting', () => {
    // Verifies that we use period as decimal separator and comma as thousands separator
    expect(formatNumber(1234.56)).toBe('1,234.56');
    expect(formatNumber(9876543.21)).toBe('9,876,543.21');
  });
}); 