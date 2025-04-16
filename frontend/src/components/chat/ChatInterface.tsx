import { useState, useRef, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Loader2, Send, Bot, User, Save, Trash, Plus, List } from 'lucide-react';
import { useConversation } from '@/hooks/useConversation';
import { useConversationList } from '@/hooks/useConversationList';
import { useAuth } from '@/hooks/auth';
import { usePortfolioData } from '@/hooks/usePortfolioData';
import { useMarketData } from '@/hooks/useMarketData';
import { toast } from 'sonner';
import { Message } from '@/services/conversationService';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Drawer, DrawerContent, DrawerTrigger } from '@/components/ui/drawer';
import { format } from 'date-fns';

// API client for backend communication
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

interface ChatInterfaceProps {
  initialConversationId?: string;
}

export function ChatInterface({ initialConversationId }: ChatInterfaceProps) {
  const [activeConversationId, setActiveConversationId] = useState<string | undefined>(initialConversationId);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const scrollAreaRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // Get auth context
  const { user } = useAuth();

  // Get conversation hooks
  const {
    messages,
    conversation,
    loading: conversationLoading,
    error: conversationError,
    createConversation,
    addMessage,
    updateTitle,
    deleteConversation
  } = useConversation(activeConversationId);

  // Get conversation list hook
  const {
    conversations,
    loading: conversationsLoading,
    error: conversationsError,
    refreshConversations,
    deleteConversation: deleteConversationFromList
  } = useConversationList();

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
  }, [activeConversationId]);

  // Show error toast if there's an error
  useEffect(() => {
    if (conversationError) {
      toast.error(conversationError);
    }
    if (conversationsError) {
      toast.error(conversationsError);
    }
  }, [conversationError, conversationsError]);

  // Create a new conversation
  const handleNewConversation = async () => {
    if (!user) {
      toast.error('You must be logged in to create a conversation');
      return;
    }

    const title = `New Conversation ${new Date().toLocaleString()}`;
    const conversationId = await createConversation(title);

    if (conversationId) {
      setActiveConversationId(conversationId);
      await refreshConversations();
      toast.success('New conversation created');
    }
  };

  // Send a message to the AI
  const handleSendMessage = async () => {
    if (!input.trim()) return;
    if (!user) {
      toast.error('You must be logged in to send messages');
      return;
    }

    // Create or get conversation ID
    let conversationId = activeConversationId;
    if (!conversationId) {
      conversationId = await createConversation(`Conversation ${new Date().toLocaleString()}`);
      if (!conversationId) {
        toast.error('Failed to create conversation');
        return;
      }
      setActiveConversationId(conversationId);
      await refreshConversations();
    }

    // Add user message to the conversation
    const userMessage: Message = {
      role: 'user',
      content: input,
      metadata: {
        timestamp: new Date().toISOString()
      }
    };

    await addMessage(userMessage);
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

      // Send message to backend API
      const response = await fetch(`${API_URL}/api/ai/chat`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          user_id: user.id,
          message: input,
          session_id: conversationId,
          trading_context: tradingContext
        })
      });

      if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
      }

      const data = await response.json();

      // Add assistant message to the conversation
      const assistantMessage: Message = {
        role: 'assistant',
        content: data.data.response,
        metadata: {
          timestamp: new Date().toISOString(),
          function_calls: data.data.function_calls
        }
      };

      await addMessage(assistantMessage);

      // Update conversation title if it's the first message
      if (messages.length <= 1 && conversation?.title.startsWith('Conversation ')) {
        const newTitle = input.length > 30 ? `${input.substring(0, 30)}...` : input;
        await updateTitle(newTitle);
        await refreshConversations();
      }
    } catch (error) {
      console.error('Error sending message:', error);
      toast.error('Failed to get a response from the AI assistant');

      // Add error message
      const errorMessage: Message = {
        role: 'assistant',
        content: 'Sorry, I encountered an error while processing your request. Please try again later.',
        metadata: {
          timestamp: new Date().toISOString(),
          error: true
        }
      };

      await addMessage(errorMessage);
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

  // Handle conversation deletion
  const handleDeleteConversation = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this conversation?')) {
      const success = await deleteConversationFromList(id);
      if (success) {
        toast.success('Conversation deleted');
        if (id === activeConversationId) {
          setActiveConversationId(undefined);
        }
      }
    }
  };

  // Switch to a different conversation
  const handleSwitchConversation = (id: string) => {
    setActiveConversationId(id);
    setIsDrawerOpen(false);
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-2xl font-bold">AI Trading Assistant</h2>
        <div className="flex gap-2">
          <Drawer open={isDrawerOpen} onOpenChange={setIsDrawerOpen}>
            <DrawerTrigger asChild>
              <Button variant="outline" className="border-2 border-black">
                <List className="h-4 w-4 mr-2" />
                Conversations
              </Button>
            </DrawerTrigger>
            <DrawerContent className="p-4 border-t-2 border-black">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-bold">Your Conversations</h3>
                <Button onClick={handleNewConversation} className="bg-black text-white">
                  <Plus className="h-4 w-4 mr-2" />
                  New Chat
                </Button>
              </div>

              {conversationsLoading ? (
                <div className="flex justify-center p-4">
                  <Loader2 className="h-6 w-6 animate-spin" />
                </div>
              ) : conversations.length === 0 ? (
                <div className="text-center p-4 text-gray-500">
                  No conversations yet. Start a new chat!
                </div>
              ) : (
                <div className="space-y-2 max-h-[400px] overflow-y-auto">
                  {conversations.map((conv) => (
                    <div
                      key={conv.id}
                      className={`p-3 rounded-md border-2 flex justify-between items-center cursor-pointer ${
                        conv.id === activeConversationId
                          ? 'bg-black text-white border-black'
                          : 'bg-white text-black border-gray-200 hover:border-black'
                      }`}
                      onClick={() => handleSwitchConversation(conv.id)}
                    >
                      <div className="flex-1 truncate">
                        <div className="font-medium truncate">{conv.title}</div>
                        <div className="text-xs opacity-70">
                          {format(new Date(conv.updatedAt), 'MMM d, yyyy h:mm a')}
                        </div>
                      </div>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={(e) => {
                          e.stopPropagation();
                          handleDeleteConversation(conv.id);
                        }}
                        className={conv.id === activeConversationId ? 'text-white hover:text-red-300' : 'text-gray-500 hover:text-red-500'}
                      >
                        <Trash className="h-4 w-4" />
                      </Button>
                    </div>
                  ))}
                </div>
              )}
            </DrawerContent>
          </Drawer>

          <Button onClick={handleNewConversation} className="bg-black text-white border-2 border-black">
            <Plus className="h-4 w-4 mr-2" />
            New Chat
          </Button>
        </div>
      </div>

      <Card className="flex flex-col h-full border-2 border-black">
        <CardHeader className="border-b-2 border-black px-4 py-2">
          <div className="flex justify-between items-center">
            <CardTitle className="text-lg font-mono">
              {conversation ? conversation.title : 'New Conversation'}
            </CardTitle>
            {activeConversationId && (
              <Button
                variant="ghost"
                size="icon"
                onClick={() => handleDeleteConversation(activeConversationId)}
                className="text-gray-500 hover:text-red-500"
              >
                <Trash className="h-4 w-4" />
              </Button>
            )}
          </div>
        </CardHeader>

        <CardContent className="flex-1 p-0">
          {conversationLoading ? (
            <div className="flex justify-center items-center h-full">
              <Loader2 className="h-8 w-8 animate-spin" />
            </div>
          ) : messages.length === 0 ? (
            <div className="flex flex-col justify-center items-center h-full text-center p-4">
              <Bot className="h-16 w-16 mb-4 text-gray-400" />
              <h3 className="text-xl font-bold mb-2">Welcome to AI Trading Assistant</h3>
              <p className="text-gray-500 max-w-md">
                Ask me about your portfolio, market trends, or get trading advice. I'm here to help you make informed decisions.
              </p>
            </div>
          ) : (
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
                          {message.metadata?.timestamp
                            ? new Date(message.metadata.timestamp).toLocaleTimeString()
                            : new Date().toLocaleTimeString()}
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
          )}
        </CardContent>

        <CardFooter className="border-t-2 border-black p-4">
          <div className="flex w-full items-center space-x-2">
            <Input
              ref={inputRef}
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Ask about your portfolio, market trends, or trading advice..."
              disabled={isLoading || conversationLoading}
              className="flex-1 border-2 border-black font-mono"
            />
            <Button
              onClick={handleSendMessage}
              disabled={isLoading || conversationLoading || !input.trim()}
              className="bg-black text-white border-2 border-black hover:bg-gray-800"
            >
              {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
            </Button>
          </div>
        </CardFooter>
      </Card>
    </div>
  );
}

export default ChatInterface;
