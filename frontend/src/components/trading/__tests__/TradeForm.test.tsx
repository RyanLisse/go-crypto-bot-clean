import { describe, it, expect, beforeEach, vi } from 'vitest';

import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { toast } from 'sonner';
import TradeForm from '../TradeForm';
import * as tradeQueries from '@/hooks/queries/useTradeQueries';
type TradeRequest = { symbol: string; side: 'buy' | 'sell'; amount: number; price: number };


// Mock the toast
vi.mock('sonner', () => ({
  toast: vi.fn(),
}));

// Mock the trade queries
vi.mock('@/hooks/queries/useTradeQueries', () => ({
  useExecuteTradeMutation: vi.fn(),
}));

describe('TradeForm', () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  const mockExecuteTrade = vi.fn();
  const mockMutate = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    
    // Setup the mock implementation
    vi.mocked(tradeQueries.useExecuteTradeMutation).mockReturnValue({
      mutate: mockMutate,
      isPending: false,
      isError: false,
      isIdle: true,
      isPaused: false,
      submittedAt: 0,
      error: null,
      isSuccess: false,
      data: undefined,
      reset: vi.fn(),
      variables: undefined,
      context: undefined,
      status: 'idle',
      failureCount: 0,
      failureReason: null,
      mutateAsync: mockExecuteTrade,
    });
  });

  it('renders the trade form correctly', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <TradeForm />
      </QueryClientProvider>
    );

    // Check if the form elements are rendered
    expect(screen.getByText(/Market/i)).toBeInTheDocument();
    expect(screen.getAllByText(/BTC/i)[0]).toBeInTheDocument();
    expect(screen.getByText(/Buy/i)).toBeInTheDocument();
    expect(screen.getByText(/Sell/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Amount/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Price/i)).toBeInTheDocument();
  });

  it('validates form inputs before submission', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <TradeForm />
      </QueryClientProvider>
    );

    // Try to submit without filling the form
    const submitButton = screen.getByRole('button', { name: /Buy BTC/i });
    fireEvent.click(submitButton);

    // Check if validation error is shown
    await waitFor(() => {
      expect(toast).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'Error',
          description: 'Please enter both amount and price',
        })
      );
    });

    // Verify that the trade execution was not called
    expect(mockMutate).not.toHaveBeenCalled();
  });

  it('executes a buy trade successfully', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <TradeForm />
      </QueryClientProvider>
    );

    // Fill the form
    const amountInput = screen.getByLabelText(/Amount/i);
    const priceInput = screen.getByLabelText(/Price/i);
    
    fireEvent.change(amountInput, { target: { value: '0.1' } });
    fireEvent.change(priceInput, { target: { value: '50000' } });

    // Submit the form
    const submitButton = screen.getByRole('button', { name: /Buy BTC/i });
    fireEvent.click(submitButton);

    // Check if the trade execution was called with correct parameters
    await waitFor(() => {
      expect(mockMutate).toHaveBeenCalledWith(
        {
          symbol: 'BTC',
          side: 'buy',
          amount: 0.1,
          price: 50000,
        },
        expect.anything()
      );
    });
  });

  it('executes a sell trade successfully', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <TradeForm />
      </QueryClientProvider>
    );

    // Select sell order type
    const sellButton = screen.getByRole('button', { name: /Sell/i });
    fireEvent.click(sellButton);

    // Fill the form
    const amountInput = screen.getByLabelText(/Amount/i);
    const priceInput = screen.getByLabelText(/Price/i);
    
    fireEvent.change(amountInput, { target: { value: '0.1' } });
    fireEvent.change(priceInput, { target: { value: '50000' } });

    // Submit the form
    const submitButton = screen.getByRole('button', { name: /Sell BTC/i });
    fireEvent.click(submitButton);

    // Check if the trade execution was called with correct parameters
    await waitFor(() => {
      expect(mockMutate).toHaveBeenCalledWith(
        {
          symbol: 'BTC',
          side: 'sell',
          amount: 0.1,
          price: 50000,
        },
        expect.anything()
      );
    });
  });

  it('shows an error message when trade execution fails', async () => {
    // Mock the mutation to simulate an error
    vi.mocked(tradeQueries.useExecuteTradeMutation).mockReturnValue({
      mutate: vi.fn((_, options) => {
        if (options?.onError) {
          options.onError(new Error('Failed to execute trade'), { symbol: 'BTC', side: 'buy' as const, amount: 0, price: 0 }, undefined);
        }
      }),
      isPending: false,
      isError: true,
      isIdle: false,
      isPaused: false,
      submittedAt: 0,
      error: new Error('Failed to execute trade'),
      isSuccess: false,
      data: undefined,
      reset: vi.fn(),
      variables: undefined,
      context: undefined,
      status: 'error',
      failureCount: 1,
      failureReason: new Error('Failed to execute trade'),
      mutateAsync: vi.fn().mockRejectedValue(new Error('Failed to execute trade')),
    });

    render(
      <QueryClientProvider client={queryClient}>
        <TradeForm />
      </QueryClientProvider>
    );

    // Fill the form
    const amountInput = screen.getByLabelText(/Amount/i);
    const priceInput = screen.getByLabelText(/Price/i);
    
    fireEvent.change(amountInput, { target: { value: '0.1' } });
    fireEvent.change(priceInput, { target: { value: '50000' } });

    // Submit the form
    const submitButton = screen.getByRole('button', { name: /Buy BTC/i });
    fireEvent.click(submitButton);

    // Check if error message is shown
    await waitFor(() => {
      expect(toast).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'Error',
          description: 'Failed to place order: Failed to execute trade',
        })
      );
    });
  });
});
