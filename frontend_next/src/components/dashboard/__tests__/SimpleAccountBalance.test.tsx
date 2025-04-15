import { vi, describe, it, expect } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { SimpleAccountBalance } from '../SimpleAccountBalance';
import { ToastProvider } from '@/components/ui/toast';

// Mock the toast hook using vi.mock instead of jest.mock
vi.mock('@/hooks/use-toast', () => ({
  useToast: () => ({
    toast: vi.fn(),
    success: vi.fn(),
    error: vi.fn(),
  })
}));

// Mock the queries hook
vi.mock('@/hooks/queries', () => ({
  useWalletData: () => ({
    data: {
      balance: 1234.56,
      assets: [
        { symbol: 'BTC', amount: 0.5, valueUsd: 15000 },
        { symbol: 'ETH', amount: 10, valueUsd: 20000 }
      ]
    },
    isLoading: false,
    isError: false,
    error: null
  })
}));

describe('SimpleAccountBalance', () => {
  it('renders loading state initially', () => {
    render(
      <ToastProvider>
        <SimpleAccountBalance />
      </ToastProvider>
    );
    
    // Check for loading indicator
    expect(screen.getByRole('status')).toBeInTheDocument();
  });

  it('renders account balance data after loading', async () => {
    render(
      <ToastProvider>
        <SimpleAccountBalance />
      </ToastProvider>
    );
    
    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.queryByRole('status')).not.toBeInTheDocument();
    });
    
    // Check for account balance heading
    expect(screen.getByText('ACCOUNT BALANCE')).toBeInTheDocument();
    
    // Check for some balance data
    expect(screen.getByText('BTC')).toBeInTheDocument();
    expect(screen.getByText('ETH')).toBeInTheDocument();
    expect(screen.getByText('USDT')).toBeInTheDocument();
    
    // Check for last updated text
    expect(screen.getByText(/Last updated:/)).toBeInTheDocument();
  });

  it('renders the account balance correctly', () => {
    render(<SimpleAccountBalance />);
    
    // We're mocking the data, so we can expect specific values
    expect(screen.getByText(/\$1,234.56/)).toBeInTheDocument();
  });
});
