/**
 * AI Client for interacting with the backend AI API
 * Includes security measures like input sanitization, rate limiting, and circuit breaking
 * Provides fallback to local Gemini model when backend is unavailable
 */

import { sanitizeUserInput, globalRateLimiter, globalCircuitBreaker } from './security';
import { sendMessage as sendGeminiMessage, GeminiChatRequest } from './gemini';

export interface ChatMessage {
  role: 'user' | 'assistant';
  content: string;
}

export interface ChatRequest {
  messages: ChatMessage[];
  session_id?: string;
}

export interface ChatResponse {
  message: ChatMessage;
  session_id: string;
}

export interface FunctionRequest {
  function_name: string;
  parameters: Record<string, any>;
}

export interface FunctionResponse {
  result: any;
  error?: string;
}

// Metrics for monitoring AI usage
export class AIMetrics {
  private static instance: AIMetrics;
  private requestCount: number = 0;
  private errorCount: number = 0;
  private totalLatency: number = 0;
  private lastRequestTime: Date | null = null;
  private fallbackCount: number = 0;
  private usingFallback: boolean = false;

  private constructor() {}

  public static getInstance(): AIMetrics {
    if (!AIMetrics.instance) {
      AIMetrics.instance = new AIMetrics();
    }
    return AIMetrics.instance;
  }

  public recordRequest(latency: number, error: boolean, fallback: boolean = false): void {
    this.requestCount++;
    this.totalLatency += latency;
    this.lastRequestTime = new Date();

    if (error) {
      this.errorCount++;
    }

    if (fallback) {
      this.fallbackCount++;
    }
  }

  public setFallbackMode(usingFallback: boolean): void {
    this.usingFallback = usingFallback;
  }

  public isUsingFallback(): boolean {
    return this.usingFallback;
  }

  public getMetrics(): Record<string, any> {
    const avgLatency = this.requestCount > 0 ? this.totalLatency / this.requestCount : 0;

    return {
      requestCount: this.requestCount,
      errorCount: this.errorCount,
      fallbackCount: this.fallbackCount,
      usingFallback: this.usingFallback,
      avgLatencyMs: avgLatency,
      lastRequestTime: this.lastRequestTime,
      errorRate: this.requestCount > 0 ? this.errorCount / this.requestCount : 0,
      fallbackRate: this.requestCount > 0 ? this.fallbackCount / this.requestCount : 0,
    };
  }
}

// Get the metrics instance
const metrics = AIMetrics.getInstance();

/**
 * Send a chat message to the AI
 * @param message The message to send
 * @param sessionId Optional session ID for continuing a conversation
 * @param forceFallback Force using the fallback mode (Gemini) even if backend is available
 * @returns The AI response
 */
export async function sendChatMessage(
  message: string,
  sessionId: string | null = null,
  forceFallback: boolean = false
): Promise<ChatResponse> {
  // Apply security measures
  if (!globalRateLimiter.allow()) {
    throw new Error('Rate limit exceeded. Please try again later.');
  }

  if (!forceFallback && !globalCircuitBreaker.isAllowed()) {
    throw new Error('Service temporarily unavailable due to errors. Please try again later.');
  }

  // Sanitize user input
  const sanitizedMessage = sanitizeUserInput(message);

  // Get authentication token
  const token = localStorage.getItem('token');
  if (!token && !forceFallback) {
    throw new Error('Authentication required');
  }

  const startTime = performance.now();
  let error = false;
  let usingFallback = forceFallback || metrics.isUsingFallback();

  // Try to use the backend API first, unless fallback is forced
  if (!usingFallback) {
    try {
      const response = await fetch('/api/chat', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          messages: [{ role: 'user', content: sanitizedMessage }],
          session_id: sessionId,
        }),
      });

      if (!response.ok) {
        error = true;
        globalCircuitBreaker.recordFailure();
        throw new Error(`API error: ${response.status}`);
      }

      globalCircuitBreaker.recordSuccess();
      metrics.setFallbackMode(false);

      const latency = performance.now() - startTime;
      metrics.recordRequest(latency, error, false);

      return response.json();
    } catch (err) {
      // If backend fails, switch to fallback mode
      console.warn('Backend API failed, switching to fallback mode:', err);
      usingFallback = true;
      metrics.setFallbackMode(true);
      error = false; // Reset error flag for fallback attempt
    }
  }

  // Fallback to Gemini if backend is unavailable or fallback is forced
  if (usingFallback) {
    try {
      // Create a request for Gemini
      const geminiRequest: GeminiChatRequest = {
        message: sanitizedMessage,
      };

      // Send message to Gemini
      const geminiResponse = await sendGeminiMessage(geminiRequest);

      // Format the response to match the ChatResponse interface
      const response: ChatResponse = {
        message: {
          role: 'assistant',
          content: geminiResponse.text,
        },
        session_id: sessionId || 'gemini-fallback-session',
      };

      const latency = performance.now() - startTime;
      metrics.recordRequest(latency, false, true);

      return response;
    } catch (geminiErr) {
      error = true;
      console.error('Gemini fallback also failed:', geminiErr);

      const latency = performance.now() - startTime;
      metrics.recordRequest(latency, true, true);

      throw new Error('Both backend and fallback AI services failed. Please try again later.');
    }
  }

  // This should never happen, but TypeScript requires a return statement
  throw new Error('Unexpected error in AI client');
}

/**
 * Execute a trading function
 * @param functionName The name of the function to execute
 * @param parameters The parameters for the function
 * @returns The function result
 */
export async function executeTradingFunction(
  functionName: string,
  parameters: Record<string, any>
): Promise<FunctionResponse> {
  // Apply security measures
  if (!globalRateLimiter.allow()) {
    throw new Error('Rate limit exceeded. Please try again later.');
  }

  if (!globalCircuitBreaker.isAllowed()) {
    throw new Error('Service temporarily unavailable due to errors. Please try again later.');
  }

  // Get authentication token
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('Authentication required');
  }

  const startTime = performance.now();
  let error = false;

  try {
    const response = await fetch('/api/function', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({
        function_name: functionName,
        parameters: parameters,
      }),
    });

    if (!response.ok) {
      error = true;
      globalCircuitBreaker.recordFailure();
      throw new Error(`API error: ${response.status}`);
    }

    globalCircuitBreaker.recordSuccess();
    return response.json();
  } catch (err) {
    error = true;
    throw err;
  } finally {
    const latency = performance.now() - startTime;
    metrics.recordRequest(latency, error);
  }
}

/**
 * Stream a chat message to the AI
 * This function returns a ReadableStream for streaming responses
 * @param message The message to send
 * @param sessionId Optional session ID for continuing a conversation
 * @returns A ReadableStream of the AI response
 */
export async function streamChatMessage(
  message: string,
  sessionId: string | null = null
): Promise<ReadableStream<Uint8Array>> {
  // Apply security measures
  if (!globalRateLimiter.allow()) {
    throw new Error('Rate limit exceeded. Please try again later.');
  }

  if (!globalCircuitBreaker.isAllowed()) {
    throw new Error('Service temporarily unavailable due to errors. Please try again later.');
  }

  // Sanitize user input
  const sanitizedMessage = sanitizeUserInput(message);

  // Get authentication token
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('Authentication required');
  }

  const startTime = performance.now();
  let error = false;

  try {
    const response = await fetch('/api/chat/stream', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({
        messages: [{ role: 'user', content: sanitizedMessage }],
        session_id: sessionId,
      }),
    });

    if (!response.ok) {
      error = true;
      globalCircuitBreaker.recordFailure();
      throw new Error(`API error: ${response.status}`);
    }

    globalCircuitBreaker.recordSuccess();
    return response.body as ReadableStream<Uint8Array>;
  } catch (err) {
    error = true;
    throw err;
  } finally {
    const latency = performance.now() - startTime;
    metrics.recordRequest(latency, error);
  }
}

/**
 * Get AI usage metrics
 * @returns Current AI usage metrics
 */
export function getAIMetrics(): Record<string, any> {
  return metrics.getMetrics();
}

export default {
  sendChatMessage,
  executeTradingFunction,
  streamChatMessage,
  getAIMetrics,
};
