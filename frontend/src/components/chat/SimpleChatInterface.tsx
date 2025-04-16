import { useState, useRef, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Loader2, Send, Bot, User } from 'lucide-react';
import { toast } from 'sonner';

// API client for backend communication
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

type Message = {
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
};

export function SimpleChatInterface() {
  const [messages, setMessages] = useState<Message[]>([
    {
      role: 'assistant',
      content: 'Hello! I am BruteBot, your AI trading assistant. How can I help you today?',
      timestamp: new Date(),
    },
  ]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const scrollAreaRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // Mock user ID for testing
  const userId = 'user123';

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

  // Send a message to the AI
  const handleSendMessage = async () => {
    if (!input.trim()) return;

    // Add user message to the conversation
    const userMessage: Message = {
      role: 'user',
      content: input,
      timestamp: new Date(),
    };

    setMessages(prev => [...prev, userMessage]);
    setInput('');
    setIsLoading(true);

    try {
      // Send message to backend API
      const response = await fetch(`${API_URL}/ai/chat`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          user_id: userId,
          message: input,
          session_id: 'simple-chat-session',
        })
      });

      if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
      }

      const data = await response.json();

      // Add assistant message to the conversation
      const assistantMessage: Message = {
        role: 'assistant',
        content: data.data?.response || "I'm sorry, I couldn't process your request.",
        timestamp: new Date(),
      };

      setMessages(prev => [...prev, assistantMessage]);
    } catch (error) {
      console.error('Error sending message:', error);
      toast.error('Failed to get a response from the AI assistant');

      // Add error message
      const errorMessage: Message = {
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

  // Handle key down event
  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !isLoading) {
      handleSendMessage();
    }
  };

  return (
    <Card className="flex flex-col h-full border-2 border-black">
      <CardHeader className="border-b-2 border-black px-4 py-2">
        <CardTitle className="text-lg font-mono">BruteBot AI Trading Assistant</CardTitle>
      </CardHeader>

      <CardContent className="flex-1 p-0">
        <ScrollArea className="h-[500px] p-4" ref={scrollAreaRef}>
          <div className="space-y-4">
            {messages.map((message, index) => (
              <div
                key={index}
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
