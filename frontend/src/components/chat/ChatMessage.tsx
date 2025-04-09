import { Bot, User } from 'lucide-react';
import { Message } from '@/services/conversationService';
import { cn } from '@/lib/utils';
import { Card } from '@/components/ui/card';

interface ChatMessageProps {
  message: Message;
  isLast: boolean;
}

export function ChatMessage({ message, isLast }: ChatMessageProps) {
  const isUser = message.role === 'user';
  const timestamp = message.metadata?.timestamp 
    ? new Date(message.metadata.timestamp).toLocaleTimeString() 
    : new Date().toLocaleTimeString();
  
  // Check if the message contains code blocks
  const hasCodeBlock = message.content.includes('```');
  
  // Process message content to handle markdown-like formatting
  const processContent = (content: string) => {
    if (!hasCodeBlock) return content;
    
    // Split by code blocks
    const parts = content.split(/(```(?:[\w-]+)?\n[\s\S]*?\n```)/g);
    
    return parts.map((part, index) => {
      // Check if this part is a code block
      if (part.startsWith('```') && part.endsWith('```')) {
        // Extract language and code
        const match = part.match(/```([\w-]+)?\n([\s\S]*?)\n```/);
        if (match) {
          const [_, language, code] = match;
          return (
            <div key={index} className="my-2 rounded-md bg-gray-800 text-white p-2 font-mono text-sm overflow-x-auto">
              {code}
            </div>
          );
        }
      }
      
      // Regular text
      return <span key={index}>{part}</span>;
    });
  };
  
  return (
    <div
      className={cn(
        "flex w-full mb-4",
        isUser ? "justify-end" : "justify-start"
      )}
    >
      <div
        className={cn(
          "flex items-start gap-2 max-w-[80%] p-3 rounded-md border-2",
          isUser
            ? "bg-black text-white border-black"
            : "bg-gray-100 text-black border-black"
        )}
      >
        <div className="mt-1 flex-shrink-0">
          {isUser ? (
            <User className="h-5 w-5" />
          ) : (
            <Bot className="h-5 w-5" />
          )}
        </div>
        <div className="flex-1 min-w-0">
          <div className="font-mono whitespace-pre-wrap break-words">
            {processContent(message.content)}
          </div>
          <div className="text-xs opacity-50 mt-1 font-mono">
            {timestamp}
          </div>
        </div>
      </div>
    </div>
  );
}

export default ChatMessage;
