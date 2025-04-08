import { useState, useRef, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Loader2, Send, Bot, User } from 'lucide-react';
import { sendMessage, GeminiChatRequest, GeminiMessage } from '@/lib/gemini';
import { usePortfolioData } from '@/hooks/usePortfolioData';
import { useMarketData } from '@/hooks/useMarketData';
import { toast } from 'sonner';

interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
}

export function ChatAssistant() {
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
  const scrollAreaRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  
  // Get portfolio and market data to provide context to the AI
  const { data: portfolioData } = usePortfolioData();
  const { data: marketData } = useMarketData();
  
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
  
  // Convert chat messages to Gemini format
  const getGeminiHistory = (): GeminiMessage[] => {
    return messages.map(msg => ({
      role: msg.role === 'user' ? 'user' : 'model',
      parts: [msg.content],
    }));
  };
  
  const handleSendMessage = async () => {
    if (!input.trim()) return;
    
    // Add user message
    const userMessage: ChatMessage = {
      id: Date.now().toString(),
      role: 'user',
      content: input,
      timestamp: new Date(),
    };
    
    setMessages(prev => [...prev, userMessage]);
    setInput('');
    setIsLoading(true);
    
    try {
      // Prepare trading context
      const tradingContext = {
        portfolio: portfolioData ? {
          totalValue: portfolioData.totalValue,
          holdings: portfolioData.holdings.map(h => ({
            symbol: h.symbol,
            quantity: parseFloat(h.amount),
            value: h.value,
          })),
        } : undefined,
        marketData: marketData ? {
          topGainers: marketData.topGainers.map(coin => ({
            symbol: coin.symbol,
            priceChange: coin.priceChange,
          })),
          topLosers: marketData.topLosers.map(coin => ({
            symbol: coin.symbol,
            priceChange: coin.priceChange,
          })),
        } : undefined,
      };
      
      // Create request
      const request: GeminiChatRequest = {
        message: input,
        history: getGeminiHistory(),
        tradingContext,
      };
      
      // Send message to Gemini
      const response = await sendMessage(request);
      
      // Add assistant message
      const assistantMessage: ChatMessage = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: response.text,
        timestamp: new Date(),
      };
      
      setMessages(prev => [...prev, assistantMessage]);
    } catch (error) {
      console.error('Error sending message:', error);
      toast.error('Failed to get a response from the AI assistant');
      
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
        <CardTitle className="text-lg font-mono">AI Trading Assistant</CardTitle>
      </CardHeader>
      <CardContent className="flex-1 p-0">
        <ScrollArea className="h-[400px] p-4" ref={scrollAreaRef}>
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
                  <div>
                    <div className="font-mono whitespace-pre-wrap">{message.content}</div>
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
            placeholder="Ask about your portfolio, market trends, or trading advice..."
            disabled={isLoading}
            className="flex-1 border-2 border-black font-mono"
          />
          <Button
            onClick={handleSendMessage}
            disabled={isLoading || !input.trim()}
            className="bg-black text-white border-2 border-black hover:bg-gray-800"
          >
            {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
          </Button>
        </div>
      </CardFooter>
    </Card>
  );
}

export default ChatAssistant;
