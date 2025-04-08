import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import TradingAssistant from '../TradingAssistant';
import { mockAIService } from '@/test/mocks/aiService';
import { globalRateLimiter } from '@/lib/security';

// Import the mock
import '@/test/mocks/aiService';

// Mock toast
vi.mock('sonner', () => ({
  toast: {
    error: vi.fn(),
    success: vi.fn(),
  },
}));

describe('TradingAssistant', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.spyOn(globalRateLimiter, 'allow').mockReturnValue(true);
  });
  
  it('should render the component', () => {
    render(<TradingAssistant />);
    
    expect(screen.getByText(/Trading Assistant/i)).toBeInTheDocument();
    expect(screen.getByText(/Hello! I'm your crypto trading assistant/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/Ask about trading strategies/i)).toBeInTheDocument();
  });
  
  it('should send a message and display the response', async () => {
    const user = userEvent.setup();
    render(<TradingAssistant />);
    
    const input = screen.getByPlaceholderText(/Ask about trading strategies/i);
    const sendButton = screen.getByRole('button');
    
    await user.type(input, 'Tell me about Bitcoin');
    await user.click(sendButton);
    
    // Check that the user message is displayed
    expect(screen.getByText('Tell me about Bitcoin')).toBeInTheDocument();
    
    // Check that the AI service was called
    expect(mockAIService.sendChatMessage).toHaveBeenCalledWith(
      'Tell me about Bitcoin',
      null
    );
    
    // Wait for the response to be displayed
    await waitFor(() => {
      expect(screen.getByText(/I'm your AI trading assistant/i)).toBeInTheDocument();
    });
  });
  
  it('should display a chart when requesting price data', async () => {
    const user = userEvent.setup();
    render(<TradingAssistant />);
    
    const input = screen.getByPlaceholderText(/Ask about trading strategies/i);
    const sendButton = screen.getByRole('button');
    
    await user.type(input, 'Show me the Bitcoin price chart');
    await user.click(sendButton);
    
    // Wait for the chart to be displayed
    await waitFor(() => {
      expect(screen.getByText('BTC Price Chart')).toBeInTheDocument();
    });
  });
  
  it('should display a portfolio when requesting portfolio data', async () => {
    const user = userEvent.setup();
    render(<TradingAssistant />);
    
    const input = screen.getByPlaceholderText(/Ask about trading strategies/i);
    const sendButton = screen.getByRole('button');
    
    await user.type(input, 'Show me my portfolio');
    await user.click(sendButton);
    
    // Wait for the portfolio to be displayed
    await waitFor(() => {
      expect(screen.getByText('Portfolio Summary')).toBeInTheDocument();
      expect(screen.getByText('$10000.00')).toBeInTheDocument();
    });
  });
  
  it('should show metrics when clicking the Show Metrics button', async () => {
    const user = userEvent.setup();
    render(<TradingAssistant />);
    
    const metricsButton = screen.getByText('Show Metrics');
    await user.click(metricsButton);
    
    expect(screen.getByText('Requests:')).toBeInTheDocument();
    expect(screen.getByText('Errors:')).toBeInTheDocument();
    expect(screen.getByText('Avg Latency:')).toBeInTheDocument();
  });
  
  it('should disable input when rate limited', async () => {
    vi.spyOn(globalRateLimiter, 'allow').mockReturnValue(false);
    
    render(<TradingAssistant />);
    
    const input = screen.getByPlaceholderText(/Ask about trading strategies/i);
    const sendButton = screen.getByRole('button');
    
    expect(input).toBeDisabled();
    expect(sendButton).toBeDisabled();
    expect(screen.getByText(/Rate limit reached/i)).toBeInTheDocument();
  });
});
