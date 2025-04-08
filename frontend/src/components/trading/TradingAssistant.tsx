import { useState, useRef, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Loader2, Send, Bot, User, AlertTriangle, WifiOff } from 'lucide-react';
import { toast } from 'sonner';
import CryptoChart, { ChartData } from '../charts/CryptoChart';
import PortfolioCard, { PortfolioData } from '../portfolio/PortfolioCard';
import { sendChatMessage, getAIMetrics } from '@/lib/aiClient';
import { sanitizeUserInput, globalRateLimiter } from '@/lib/security';
import { useBackendStatus } from '@/hooks/useBackendStatus';

interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
  visualData?: {
    type: 'chart' | 'portfolio';
    data: ChartData | PortfolioData;
  };
}

export function TradingAssistant() {
  const [messages, setMessages] = useState<ChatMessage[]>([
    {
      id: '1',
      role: 'assistant',
      content: 'Hello! I\'m your crypto trading assistant. How can I help you today?',
      timestamp: new Date(),
    },
  ]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [showMetrics, setShowMetrics] = useState(false);
  const [metrics, setMetrics] = useState<Record<string, any>>({});
  const [fallbackNotified, setFallbackNotified] = useState(false);
  const { isConnected } = useBackendStatus();
  const scrollAreaRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // Update metrics periodically
  useEffect(() => {
    const updateMetrics = () => {
      setMetrics(getAIMetrics());
    };

    // Update metrics immediately
    updateMetrics();

    // Update metrics every 5 seconds
    const intervalId = setInterval(updateMetrics, 5000);

    return () => {
      clearInterval(intervalId);
    };
  }, []);

  // Check for fallback mode and show notification
  useEffect(() => {
    const usingFallback = getAIMetrics().usingFallback;

    if (usingFallback && !fallbackNotified && !isConnected) {
      // Add a message about fallback mode
      const fallbackMessage: ChatMessage = {
        id: Date.now().toString(),
        role: 'assistant',
        content: 'I\'m currently operating in offline mode with limited capabilities. Some advanced features like market data analysis and portfolio visualization may not be available until connection to the backend is restored.',
        timestamp: new Date(),
      };

      setMessages(prev => [...prev, fallbackMessage]);
      setFallbackNotified(true);

      // Show toast notification
      toast.warning(
        'Using AI fallback mode',
        {
          description: 'Backend connection unavailable. Using local AI model with limited capabilities.',
          duration: 5000,
          icon: <AlertTriangle className="h-4 w-4" />,
        }
      );
    }

    // Reset notification state when connection is restored
    if (isConnected && !usingFallback) {
      setFallbackNotified(false);
    }
  }, [isConnected, fallbackNotified]);

  // Scroll to bottom when messages change
  useEffect(() => {
    if (scrollAreaRef.current) {
      scrollAreaRef.current.scrollTop = scrollAreaRef.current.scrollHeight;
    }
  }, [messages]);

  // Focus input on mount
  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.focus();
    }
  }, []);

  const handleSendMessage = async () => {
    if (!input.trim()) return;

    // Check rate limiting
    if (!globalRateLimiter.allow()) {
      toast.error('You\'re sending messages too quickly. Please wait a moment before trying again.');
      return;
    }

    // Sanitize user input
    const sanitizedInput = sanitizeUserInput(input);

    // Add user message
    const userMessage: ChatMessage = {
      id: Date.now().toString(),
      role: 'user',
      content: input, // Show original input to user
      timestamp: new Date(),
    };

    setMessages(prev => [...prev, userMessage]);
    setInput('');
    setIsLoading(true);

    try {
      // Send sanitized message to API, force fallback if not connected
      const response = await sendChatMessage(sanitizedInput, sessionId, !isConnected);

      // Save session ID for continuing the conversation
      if (response.session_id) {
        setSessionId(response.session_id);
      }

      // Check if response contains visualization data
      let visualData: ChatMessage['visualData'] = undefined;

      try {
        const parsedContent = JSON.parse(response.message.content);

        if (parsedContent.chartData) {
          visualData = {
            type: 'chart',
            data: parsedContent.chartData as ChartData,
          };
        } else if (parsedContent.portfolioSnapshot) {
          visualData = {
            type: 'portfolio',
            data: parsedContent.portfolioSnapshot as PortfolioData,
          };
        }
      } catch (e) {
        // Not JSON, regular chat message
      }

      // Add assistant message
      const assistantMessage: ChatMessage = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: response.message.content,
        timestamp: new Date(),
        visualData,
      };

      setMessages(prev => [...prev, assistantMessage]);

      // Log metrics for monitoring
      const metrics = getAIMetrics();
      console.debug('AI Metrics:', metrics);
    } catch (error: any) {
      console.error('Error sending message:', error);

      // Show appropriate error message based on the error
      if (error.message.includes('Rate limit')) {
        toast.error('Rate limit exceeded. Please try again later.');
      } else if (error.message.includes('temporarily unavailable')) {
        toast.error('AI service is temporarily unavailable. Please try again later.');
      } else {
        toast.error('Failed to get a response from the AI assistant');
      }

      // Add error message
      const errorMessage: ChatMessage = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: 'Sorry, I encountered an error while processing your request. Please try again later.',
        timestamp: new Date(),
      };

      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);

      // Focus input after sending
      if (inputRef.current) {
        inputRef.current.focus();
      }
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !isLoading) {
      handleSendMessage();
    }
  };

  return (
    <Card className="flex flex-col h-full border-2 border-black">
      <CardHeader className="border-b-2 border-black px-4 py-2">
        <div className="flex justify-between items-center">
          <div className="flex items-center">
            <CardTitle className="text-lg font-mono">Trading Assistant</CardTitle>
            {getAIMetrics().usingFallback && (
              <div className="ml-2 flex items-center text-xs bg-brutal-warning/20 text-brutal-warning px-1 py-0.5">
                <WifiOff className="h-3 w-3 mr-1" />
                OFFLINE MODE
              </div>
            )}
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setShowMetrics(!showMetrics)}
            className="font-mono text-xs"
          >
            {showMetrics ? 'Hide Metrics' : 'Show Metrics'}
          </Button>
        </div>

        {/* Metrics Display */}
        {showMetrics && (
          <div className="mt-2 p-2 bg-gray-100 border border-black text-xs font-mono">
            <div className="grid grid-cols-2 gap-2">
              <div>Requests: {metrics.requestCount || 0}</div>
              <div>Errors: {metrics.errorCount || 0}</div>
              <div>Avg Latency: {metrics.avgLatencyMs ? `${Math.round(metrics.avgLatencyMs)}ms` : 'N/A'}</div>
              <div>Error Rate: {metrics.errorRate ? `${(metrics.errorRate * 100).toFixed(1)}%` : '0%'}</div>
              <div>Fallback: {metrics.fallbackCount || 0}</div>
              <div className={metrics.usingFallback ? 'text-brutal-warning' : ''}>
                Mode: {metrics.usingFallback ? 'OFFLINE' : 'ONLINE'}
              </div>
            </div>
            <div className="mt-1 text-gray-500">
              Last Request: {metrics.lastRequestTime ? new Date(metrics.lastRequestTime).toLocaleTimeString() : 'N/A'}
            </div>
          </div>
        )}
      </CardHeader>
      <CardContent className="flex-1 p-0">
        <ScrollArea className="h-[500px] p-4" ref={scrollAreaRef}>
          <div className="space-y-4">
            {messages.map((message) => (
              <div
                key={message.id}
                className={`flex ${
                  message.role === 'user' ? 'justify-end' : 'justify-start'
                }`}
              >
                <div
                  className={`flex items-start gap-2 max-w-[80%] ${
                    message.role === 'user'
                      ? 'bg-black text-white'
                      : 'bg-gray-100 text-black'
                  } p-3 rounded-md border-2 border-black`}
                >
                  <div className="mt-1">
                    {message.role === 'user' ? (
                      <User className="h-5 w-5" />
                    ) : (
                      <Bot className="h-5 w-5" />
                    )}
                  </div>
                  <div className="w-full">
                    <div className="font-mono whitespace-pre-wrap">{message.content}</div>

                    {/* Render visualization if available */}
                    {message.visualData && (
                      <div className="mt-4">
                        {message.visualData.type === 'chart' && (
                          <CryptoChart data={message.visualData.data as ChartData} />
                        )}
                        {message.visualData.type === 'portfolio' && (
                          <PortfolioCard data={message.visualData.data as PortfolioData} />
                        )}
                      </div>
                    )}

                    <div className="text-xs opacity-50 mt-1 font-mono">
                      {message.timestamp.toLocaleTimeString()}
                    </div>
                  </div>
                </div>
              </div>
            ))}
            {isLoading && (
              <div className="flex justify-start">
                <div className="flex items-center gap-2 bg-gray-100 p-3 rounded-md border-2 border-black">
                  <Bot className="h-5 w-5" />
                  <Loader2 className="h-4 w-4 animate-spin" />
                  <span className="font-mono text-sm">Thinking...</span>
                </div>
              </div>
            )}

            {/* Rate limit warning */}
            {!globalRateLimiter.allow() && (
              <div className="flex justify-center">
                <div className="flex items-center gap-2 bg-yellow-50 p-3 rounded-md border-2 border-yellow-500 text-yellow-700 max-w-[80%]">
                  <AlertTriangle className="h-5 w-5" />
                  <span className="font-mono text-sm">Rate limit reached. Please wait a moment before sending more messages.</span>
                </div>
              </div>
            )}
          </div>
        </ScrollArea>
      </CardContent>
      <CardFooter className="border-t-2 border-black p-4">
        <div className="flex w-full items-center space-x-2">
          <Input
            ref={inputRef}
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Ask about trading strategies, market analysis, or portfolio advice..."
            disabled={isLoading || !globalRateLimiter.allow()}
            className="flex-1 border-2 border-black font-mono"
          />
          <Button
            onClick={handleSendMessage}
            disabled={isLoading || !input.trim() || !globalRateLimiter.allow()}
            className="bg-black text-white border-2 border-black hover:bg-gray-800"
          >
            {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
          </Button>
        </div>
      </CardFooter>
    </Card>
  );
}

export default TradingAssistant;
