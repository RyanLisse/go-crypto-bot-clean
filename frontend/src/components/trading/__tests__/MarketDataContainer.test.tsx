import { render, screen, fireEvent } from '@testing-library/react';
import { MarketDataContainer } from '../MarketDataContainer';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useWebSocket } from '@/hooks/useWebSocket';
import type { WebSocketConnectionState } from '@/types/websocket';
import type { MarketData } from '@/types/market';

// Mock the WebSocket hook
vi.mock('@/hooks/useWebSocket', () => {
  const mockSendMessage = vi.fn();
  const mockSubscribeTicker = vi.fn();
  return {
    useWebSocket: () => ({
      isConnected: true,
      connectionState: 'CONNECTED' as WebSocketConnectionState,
      sendMessage: mockSendMessage,
      lastMessage: null,
      subscribeTicker: mockSubscribeTicker,
      connect: vi.fn(),
      disconnect: vi.fn(),
    }),
  };
});

// Mock the MarketDataChart component
vi.mock('../MarketDataChart', () => ({
  MarketDataChart: ({ data }: { data: MarketData[] }) => (
    <div data-testid="market-data-chart">
      {JSON.stringify(data)}
    </div>
  ),
}));

describe('MarketDataContainer', () => {
  const mockSendMessage = vi.fn();
  const mockSubscribeTicker = vi.fn();
  
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useWebSocket).mockReturnValue({
      isConnected: true,
      connectionState: 'CONNECTED' as WebSocketConnectionState,
      sendMessage: mockSendMessage,
      lastMessage: null,
      subscribeTicker: mockSubscribeTicker,
      connect: vi.fn(),
      disconnect: vi.fn(),
    });
  });

  it('should render and subscribe to market data', () => {
    render(<MarketDataContainer />);
    expect(screen.getByTestId('market-data-chart')).toBeInTheDocument();
    expect(mockSubscribeTicker).toHaveBeenCalledWith(['BTC-USD']);
  });

  it('should change trading pair and resubscribe', () => {
    render(<MarketDataContainer />);
    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: 'ETH-USD' } });
    expect(mockSubscribeTicker).toHaveBeenCalledWith(['ETH-USD']);
  });

  it('should update chart with market data', () => {
    const mockMarketData: MarketData = {
      symbol: 'BTC-USD',
      price: '50000',
      timestamp: Date.now(),
    };

    vi.mocked(useWebSocket).mockReturnValue({
      isConnected: true,
      connectionState: 'CONNECTED' as WebSocketConnectionState,
      sendMessage: mockSendMessage,
      lastMessage: { data: JSON.stringify(mockMarketData) },
      subscribeTicker: mockSubscribeTicker,
      connect: vi.fn(),
      disconnect: vi.fn(),
    });

    render(<MarketDataContainer />);
    expect(screen.getByText(/"price":"50000"/)).toBeInTheDocument();
  });

  it('should show disconnected state', () => {
    vi.mocked(useWebSocket).mockReturnValue({
      isConnected: false,
      connectionState: 'DISCONNECTED' as WebSocketConnectionState,
      sendMessage: mockSendMessage,
      lastMessage: null,
      subscribeTicker: mockSubscribeTicker,
      connect: vi.fn(),
      disconnect: vi.fn(),
    });

    render(<MarketDataContainer />);
    expect(screen.getByText(/Disconnected/i)).toBeInTheDocument();
  });
}); 