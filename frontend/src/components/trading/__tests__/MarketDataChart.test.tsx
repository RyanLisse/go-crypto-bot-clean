import { render, screen } from '@testing-library/react';
import { MarketDataChart } from '../MarketDataChart';
import { describe, it, expect, vi, beforeEach } from 'vitest';

// Mock Recharts components to avoid rendering issues in test environment
vi.mock('recharts', () => ({
  Line: () => null,
  LineChart: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="line-chart">{children}</div>
  ),
  ResponsiveContainer: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="responsive-container">{children}</div>
  ),
  Tooltip: () => null,
  XAxis: () => null,
  YAxis: () => null,
}));

describe('MarketDataChart', () => {
  const mockData = [
    {
      type: 'market_data',
      pair: 'BTC-USD',
      time: '2024-03-20T10:00:00Z',
      value: 50000,
      volume: 100,
    },
    {
      type: 'market_data',
      pair: 'BTC-USD',
      time: '2024-03-20T10:01:00Z',
      value: 51000,
      volume: 150,
    },
    {
      type: 'market_data',
      pair: 'BTC-USD',
      time: '2024-03-20T10:02:00Z',
      value: 50500,
      volume: 120,
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders chart container with correct dimensions', () => {
    render(<MarketDataChart data={mockData} />);
    
    const container = screen.getByTestId('responsive-container');
    expect(container).toBeInTheDocument();
    expect(container.parentElement).toHaveClass('h-[400px]', 'w-full');
  });

  it('renders with empty data array', () => {
    render(<MarketDataChart data={[]} />);
    
    expect(screen.getByTestId('line-chart')).toBeInTheDocument();
  });

  it('renders chart components with provided data', () => {
    const { container } = render(<MarketDataChart data={mockData} />);
    
    expect(screen.getByTestId('line-chart')).toBeInTheDocument();
    expect(container.innerHTML).toContain('data-testid="line-chart"');
  });

  // Note: We're not testing the actual chart rendering or interactions
  // as Recharts components are mocked. The focus is on proper component
  // structure and data passing.
}); 