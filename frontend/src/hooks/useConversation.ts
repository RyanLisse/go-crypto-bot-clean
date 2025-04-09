import { useState, useEffect, useCallback } from 'react';
import { ConversationService, Message } from '@/services/conversationService';
import { useAuth } from '@/hooks/useAuth';

const conversationService = new ConversationService();

export function useConversation(conversationId?: string) {
  const { user } = useAuth();
  const [messages, setMessages] = useState<Message[]>([]);
  const [conversation, setConversation] = useState<{ id: string; title: string } | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Load conversation and messages
  useEffect(() => {
    if (!conversationId || !user) return;
    
    const loadConversation = async () => {
      setLoading(true);
      setError(null);
      
      try {
        // Load conversation details
        const conv = await conversationService.getConversation(conversationId);
        if (!conv) {
          setError('Conversation not found');
          setLoading(false);
          return;
        }
        
        setConversation({
          id: conv.id,
          title: conv.title
        });
        
        // Load messages
        const msgs = await conversationService.getMessages(conversationId);
        setMessages(msgs.map(msg => ({
          role: msg.role as 'user' | 'assistant',
          content: msg.content,
          metadata: msg.metadata as Record<string, any> | undefined
        })));
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load conversation');
      } finally {
        setLoading(false);
      }
    };
    
    loadConversation();
  }, [conversationId, user]);
  
  // Create a new conversation
  const createConversation = useCallback(async (title: string) => {
    if (!user) {
      setError('User not authenticated');
      return null;
    }
    
    try {
      const newConversationId = await conversationService.createConversation(user.id, title);
      return newConversationId;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create conversation');
      return null;
    }
  }, [user]);
  
  // Add a message to the conversation
  const addMessage = useCallback(async (message: Message) => {
    if (!conversationId) {
      setError('No active conversation');
      return false;
    }
    
    try {
      await conversationService.addMessage(conversationId, message);
      setMessages(prev => [...prev, message]);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add message');
      return false;
    }
  }, [conversationId]);
  
  // Update conversation title
  const updateTitle = useCallback(async (title: string) => {
    if (!conversationId) {
      setError('No active conversation');
      return false;
    }
    
    try {
      await conversationService.updateConversationTitle(conversationId, title);
      setConversation(prev => prev ? { ...prev, title } : null);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update title');
      return false;
    }
  }, [conversationId]);
  
  // Delete conversation
  const deleteConversation = useCallback(async () => {
    if (!conversationId) {
      setError('No active conversation');
      return false;
    }
    
    try {
      await conversationService.deleteConversation(conversationId);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete conversation');
      return false;
    }
  }, [conversationId]);
  
  return {
    messages,
    conversation,
    loading,
    error,
    createConversation,
    addMessage,
    updateTitle,
    deleteConversation
  };
}
