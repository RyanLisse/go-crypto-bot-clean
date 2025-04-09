import { db } from '@/db/client';
import { conversations, conversationMessages, type NewConversation, type NewConversationMessage } from '@/db/schema';
import { eq, desc } from 'drizzle-orm';
import { v4 as uuidv4 } from 'uuid';

export type Message = {
  role: 'user' | 'assistant';
  content: string;
  metadata?: Record<string, any>;
};

export class ConversationService {
  /**
   * Create a new conversation
   */
  async createConversation(userId: string, title: string): Promise<string> {
    const conversationId = uuidv4();
    
    await db.insert(conversations).values({
      id: conversationId,
      userId,
      title,
      createdAt: new Date(),
      updatedAt: new Date(),
    });
    
    return conversationId;
  }
  
  /**
   * Get a conversation by ID
   */
  async getConversation(conversationId: string) {
    return await db.query.conversations.findFirst({
      where: eq(conversations.id, conversationId),
    });
  }
  
  /**
   * List all conversations for a user
   */
  async listConversations(userId: string, limit = 20) {
    return await db.query.conversations.findMany({
      where: eq(conversations.userId, userId),
      orderBy: [desc(conversations.updatedAt)],
      limit,
    });
  }
  
  /**
   * Add a message to a conversation
   */
  async addMessage(conversationId: string, message: Message): Promise<string> {
    const messageId = uuidv4();
    
    await db.insert(conversationMessages).values({
      id: messageId,
      conversationId,
      role: message.role,
      content: message.content,
      timestamp: new Date(),
      metadata: message.metadata || null,
    });
    
    // Update conversation's updatedAt timestamp
    await db.update(conversations)
      .set({ updatedAt: new Date() })
      .where(eq(conversations.id, conversationId));
    
    return messageId;
  }
  
  /**
   * Get all messages for a conversation
   */
  async getMessages(conversationId: string) {
    return await db.query.conversationMessages.findMany({
      where: eq(conversationMessages.conversationId, conversationId),
      orderBy: [conversationMessages.timestamp],
    });
  }
  
  /**
   * Delete a conversation and all its messages
   */
  async deleteConversation(conversationId: string) {
    // Delete all messages first (due to foreign key constraint)
    await db.delete(conversationMessages)
      .where(eq(conversationMessages.conversationId, conversationId));
    
    // Then delete the conversation
    await db.delete(conversations)
      .where(eq(conversations.id, conversationId));
  }
  
  /**
   * Update a conversation's title
   */
  async updateConversationTitle(conversationId: string, title: string) {
    await db.update(conversations)
      .set({ 
        title,
        updatedAt: new Date()
      })
      .where(eq(conversations.id, conversationId));
  }
}
