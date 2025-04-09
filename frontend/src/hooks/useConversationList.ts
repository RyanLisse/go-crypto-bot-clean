import { useState, useEffect, useCallback } from 'react';
import { ConversationService } from '@/services/conversationService';
import { useAuth } from '@/hooks/useAuth';

const conversationService = new ConversationService();

export function useConversationList() {
  const { user } = useAuth();
  const [conversations, setConversations] = useState<Array<{ id: string; title: string; updatedAt: Date }>>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Load conversations
  const loadConversations = useCallback(async () => {
    if (!user) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const convs = await conversationService.listConversations(user.id);
      setConversations(convs.map(conv => ({
        id: conv.id,
        title: conv.title,
        updatedAt: conv.updatedAt
      })));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load conversations');
    } finally {
      setLoading(false);
    }
  }, [user]);

  // Load conversations on mount and when user changes
  useEffect(() => {
    loadConversations();
  }, [loadConversations]);
  
  // Delete a conversation
  const deleteConversation = useCallback(async (conversationId: string) => {
    try {
      await conversationService.deleteConversation(conversationId);
      setConversations(prev => prev.filter(conv => conv.id !== conversationId));
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete conversation');
      return false;
    }
  }, []);
  
  return {
    conversations,
    loading,
    error,
    refreshConversations: loadConversations,
    deleteConversation
  };
}
