import { useState } from 'react';
import { ChatInterface } from '@/components/chat/ChatInterface';
import { ConversationHistory } from '@/components/chat/ConversationHistory';
import { AuthProvider } from '@/hooks/auth';

export function ChatPage() {
  const [activeConversationId, setActiveConversationId] = useState<string | undefined>(undefined);
  
  const handleNewConversation = () => {
    setActiveConversationId(undefined);
  };
  
  return (
    <AuthProvider>
      <div className="container mx-auto py-6">
        <h1 className="text-3xl font-bold mb-6">AI Trading Assistant</h1>
        
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div className="md:col-span-1 hidden md:block">
            <ConversationHistory 
              activeConversationId={activeConversationId}
              onSelectConversation={setActiveConversationId}
              onNewConversation={handleNewConversation}
            />
          </div>
          
          <div className="md:col-span-3">
            <ChatInterface initialConversationId={activeConversationId} />
          </div>
        </div>
      </div>
    </AuthProvider>
  );
}

export default ChatPage;
