import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Loader2, Trash, MessageSquare, Plus } from 'lucide-react';
import { useConversationList } from '@/hooks/useConversationList';
import { format } from 'date-fns';
import { toast } from 'sonner';

interface ConversationHistoryProps {
  activeConversationId?: string;
  onSelectConversation: (id: string) => void;
  onNewConversation: () => void;
}

export function ConversationHistory({
  activeConversationId,
  onSelectConversation,
  onNewConversation
}: ConversationHistoryProps) {
  const {
    conversations,
    loading,
    error,
    deleteConversation
  } = useConversationList();

  // Handle conversation deletion
  const handleDeleteConversation = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    
    if (window.confirm('Are you sure you want to delete this conversation?')) {
      const success = await deleteConversation(id);
      if (success) {
        toast.success('Conversation deleted');
      }
    }
  };

  return (
    <Card className="h-full border-2 border-black">
      <CardHeader className="border-b-2 border-black px-4 py-2">
        <div className="flex justify-between items-center">
          <CardTitle className="text-lg font-mono">Conversation History</CardTitle>
          <Button 
            onClick={onNewConversation}
            className="bg-black text-white border-2 border-black hover:bg-gray-800"
          >
            <Plus className="h-4 w-4 mr-2" />
            New
          </Button>
        </div>
      </CardHeader>
      
      <CardContent className="p-0">
        {loading ? (
          <div className="flex justify-center items-center h-40">
            <Loader2 className="h-8 w-8 animate-spin" />
          </div>
        ) : error ? (
          <div className="p-4 text-red-500">
            Error loading conversations: {error}
          </div>
        ) : conversations.length === 0 ? (
          <div className="flex flex-col justify-center items-center h-40 text-center p-4">
            <MessageSquare className="h-12 w-12 mb-2 text-gray-400" />
            <p className="text-gray-500">No conversations yet</p>
            <Button 
              onClick={onNewConversation} 
              variant="outline" 
              className="mt-4 border-2 border-black"
            >
              Start a new conversation
            </Button>
          </div>
        ) : (
          <ScrollArea className="h-[calc(100vh-200px)]">
            <div className="p-2 space-y-2">
              {conversations.map((conversation) => (
                <div
                  key={conversation.id}
                  className={`p-3 rounded-md border-2 flex justify-between items-center cursor-pointer ${
                    conversation.id === activeConversationId
                      ? 'bg-black text-white border-black'
                      : 'bg-white text-black border-gray-200 hover:border-black'
                  }`}
                  onClick={() => onSelectConversation(conversation.id)}
                >
                  <div className="flex-1 truncate">
                    <div className="font-medium truncate">{conversation.title}</div>
                    <div className="text-xs opacity-70">
                      {format(new Date(conversation.updatedAt), 'MMM d, yyyy h:mm a')}
                    </div>
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={(e) => handleDeleteConversation(conversation.id, e)}
                    className={conversation.id === activeConversationId ? 'text-white hover:text-red-300' : 'text-gray-500 hover:text-red-500'}
                  >
                    <Trash className="h-4 w-4" />
                  </Button>
                </div>
              ))}
            </div>
          </ScrollArea>
        )}
      </CardContent>
    </Card>
  );
}

export default ConversationHistory;
