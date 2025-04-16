
import { SimpleChatInterface } from '@/components/chat/SimpleChatInterface';

export function ChatPage() {
  return (
    <div className="container mx-auto py-6">
      <h1 className="text-3xl font-bold mb-6">BruteBot AI Trading Assistant</h1>

      <div className="grid grid-cols-1 gap-6">
        <div>
          <SimpleChatInterface />
        </div>
      </div>
    </div>
  );
}

export default ChatPage;
