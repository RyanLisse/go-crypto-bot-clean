import React from 'react';
import { render, screen, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import OrderBook from '../OrderBook';
import { useWebSocket } from '@/hooks/useWebSocket';

// Mock the useWebSocket hook
vi.mock('@/hooks/useWebSocket', () => ({
  useWebSocket: vi.fn()
}));

describe('OrderBook', () => {
  const mockSendMessage = vi.fn();
  
  beforeEach(() => {
    vi.clearAllMocks();
    (useWebSocket as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      isConnected: true,
      sendMessage: mockSendMessage
    });
  });

  it('renders the order book with loading state initially', () => {
    render(<OrderBook symbol="BTC/USD" />);
    
    expect(screen.getByText('Order Book')).toBeInTheDocument();
    expect(screen.getByText('BTC/USD')).toBeInTheDocument();
    expect(screen.getByRole('status')).toBeInTheDocument(); // Loading spinner
  });

  it('subscribes to order book data when connected', () => {
    render(<OrderBook symbol="BTC/USD" />);
    
    expect(mockSendMessage).toHaveBeenCalledWith({
      type: 'subscribe',
      channel: 'orderbook',
      symbol: 'BTC/USD'
    });
  });

  it('displays order book data when received', () => {
    const mockOrderBookData = {
      type: 'orderbook',
      data: {
        bids: [[50000, 1.5], [49900, 2.0]],
        asks: [[50100, 1.0], [50200, 2.5]]
      }
    };

    let onMessageCallback: (data: string) => void;
    (useWebSocket as unknown as ReturnType<typeof vi.fn>).mockImplementation(({ onMessage }: { onMessage: (data: string) => void }) => {
      onMessageCallback = onMessage;
      return {
        isConnected: true,
        sendMessage: mockSendMessage
      };
    });

    render(<OrderBook symbol="BTC/USD" />);

    act(() => {
      onMessageCallback(JSON.stringify(mockOrderBookData));
    });

    // Check for rendered order book entries
    expect(screen.getByText('50,000.00')).toBeInTheDocument();
    expect(screen.getByText('49,900.00')).toBeInTheDocument();
    expect(screen.getByText('50,100.00')).toBeInTheDocument();
    expect(screen.getByText('50,200.00')).toBeInTheDocument();
  });

  it('handles invalid message data gracefully', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    let onMessageCallback: (data: string) => void;
    
    (useWebSocket as unknown as ReturnType<typeof vi.fn>).mockImplementation(({ onMessage }: { onMessage: (data: string) => void }) => {
      onMessageCallback = onMessage;
      return {
        isConnected: true,
        sendMessage: mockSendMessage
      };
    });

    render(<OrderBook symbol="BTC/USD" />);

    act(() => {
      onMessageCallback('invalid json');
    });

    expect(consoleSpy).toHaveBeenCalledWith(
      'Error processing order book message:',
      expect.any(Error)
    );

    consoleSpy.mockRestore();
  });

  it('resubscribes when symbol changes', () => {
    const { rerender } = render(<OrderBook symbol="BTC/USD" />);
    
    expect(mockSendMessage).toHaveBeenCalledWith({
      type: 'subscribe',
      channel: 'orderbook',
      symbol: 'BTC/USD'
    });

    rerender(<OrderBook symbol="ETH/USD" />);

    expect(mockSendMessage).toHaveBeenCalledWith({
      type: 'subscribe',
      channel: 'orderbook',
      symbol: 'ETH/USD'
    });
  });

  it('displays spread between best bid and ask', () => {
    let onMessageCallback: (data: string) => void;
    (useWebSocket as unknown as ReturnType<typeof vi.fn>).mockImplementation(({ onMessage }: { onMessage: (data: string) => void }) => {
      onMessageCallback = onMessage;
      return {
        isConnected: true,
        sendMessage: mockSendMessage
      };
    });

    render(<OrderBook symbol="BTC/USD" />);

    act(() => {
      onMessageCallback(JSON.stringify({
        type: 'orderbook',
        data: {
          bids: [[50000, 1.0]],
          asks: [[50100, 1.0]]
        }
      }));
    });

    expect(screen.getByText('Spread: 100.00')).toBeInTheDocument();
  });

  it('applies correct CSS classes for bids and asks', () => {
    let onMessageCallback: (data: string) => void;
    (useWebSocket as unknown as ReturnType<typeof vi.fn>).mockImplementation(({ onMessage }: { onMessage: (data: string) => void }) => {
      onMessageCallback = onMessage;
      return {
        isConnected: true,
        sendMessage: mockSendMessage
      };
    });

    render(<OrderBook symbol="BTC/USD" />);

    act(() => {
      onMessageCallback(JSON.stringify({
        type: 'orderbook',
        data: {
          bids: [[50000, 1.0]],
          asks: [[50100, 1.0]]
        }
      }));
    });

    const bidElement = screen.getByText('50,000.00').closest('div');
    const askElement = screen.getByText('50,100.00').closest('div');

    expect(bidElement).toHaveClass('text-brutal-success/90');
    expect(askElement).toHaveClass('text-brutal-error/90');
  });
}); 