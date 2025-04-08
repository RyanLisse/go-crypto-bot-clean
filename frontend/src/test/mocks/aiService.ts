import { vi } from 'vitest';
import { ChatResponse, FunctionResponse } from '@/lib/aiClient';

/**
 * Mock responses for the AI service
 */
export const mockResponses: Record<string, string> = {
  'portfolio': JSON.stringify({
    portfolioSnapshot: {
      totalValue: 10000,
      change: 2.5,
      assets: [
        { symbol: 'BTC', amount: 0.5, value: 5000, change: 3.2 },
        { symbol: 'ETH', amount: 2.0, value: 3000, change: -1.5 },
        { symbol: 'SOL', amount: 10.0, value: 2000, change: 5.7 },
      ],
    },
  }),
  'chart': JSON.stringify({
    chartData: {
      title: 'BTC Price Chart',
      symbol: 'BTC',
      timeframe: '1D',
      points: [
        { time: '09:00', price: 26000 },
        { time: '10:00', price: 26200 },
        { time: '11:00', price: 26100 },
        { time: '12:00', price: 26400 },
        { time: '13:00', price: 26300 },
      ],
      indicators: [
        { name: 'MA-50', key: 'ma50', color: '#ff0000' },
      ],
    },
  }),
  'default': 'I\'m your AI trading assistant. How can I help you today?',
};

/**
 * Mock AI service for testing
 */
export const mockAIService = {
  sendChatMessage: vi.fn().mockImplementation((message: string, sessionId: string | null = null): Promise<ChatResponse> => {
    // Determine which mock response to return based on the message content
    let responseContent = mockResponses.default;
    
    if (message.toLowerCase().includes('portfolio')) {
      responseContent = mockResponses.portfolio;
    } else if (message.toLowerCase().includes('chart') || message.toLowerCase().includes('price')) {
      responseContent = mockResponses.chart;
    }
    
    return Promise.resolve({
      message: {
        role: 'assistant',
        content: responseContent,
      },
      session_id: sessionId || 'mock-session-id',
    });
  }),
  
  executeTradingFunction: vi.fn().mockImplementation((functionName: string, parameters: Record<string, any>): Promise<FunctionResponse> => {
    return Promise.resolve({
      result: {
        success: true,
        message: `Successfully executed ${functionName}`,
        data: parameters,
      },
    });
  }),
  
  streamChatMessage: vi.fn().mockImplementation((message: string, sessionId: string | null = null): Promise<ReadableStream<Uint8Array>> => {
    // Create a simple readable stream with the mock response
    let responseContent = mockResponses.default;
    
    if (message.toLowerCase().includes('portfolio')) {
      responseContent = mockResponses.portfolio;
    } else if (message.toLowerCase().includes('chart') || message.toLowerCase().includes('price')) {
      responseContent = mockResponses.chart;
    }
    
    const encoder = new TextEncoder();
    const stream = new ReadableStream({
      start(controller) {
        controller.enqueue(encoder.encode(responseContent));
        controller.close();
      },
    });
    
    return Promise.resolve(stream);
  }),
  
  getAIMetrics: vi.fn().mockReturnValue({
    requestCount: 5,
    errorCount: 1,
    avgLatencyMs: 250,
    lastRequestTime: new Date().toISOString(),
    errorRate: 0.2,
  }),
};

// Mock the entire aiClient module
vi.mock('@/lib/aiClient', () => ({
  sendChatMessage: mockAIService.sendChatMessage,
  executeTradingFunction: mockAIService.executeTradingFunction,
  streamChatMessage: mockAIService.streamChatMessage,
  getAIMetrics: mockAIService.getAIMetrics,
}));
