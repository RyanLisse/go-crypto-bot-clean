import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { SimpleAccountBalance } from '../SimpleAccountBalance';
import { ToastProvider } from '@/components/ui/toast';

// Mock the toast hook
jest.mock('@/hooks/use-toast', () => ({
  useToast: () => ({
    toast: jest.fn(),
  }),
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
});
