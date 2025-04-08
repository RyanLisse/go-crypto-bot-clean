import { GoogleGenerativeAI, HarmCategory, HarmBlockThreshold } from '@google/generative-ai';

// Define types for the Gemini API
export interface GeminiMessage {
  role: 'user' | 'model';
  parts: string[];
}

export interface GeminiChatRequest {
  message: string;
  history?: GeminiMessage[];
  tradingContext?: {
    portfolio?: {
      totalValue: number;
      holdings: Array<{
        symbol: string;
        quantity: number;
        value: number;
      }>;
    };
    marketData?: {
      topGainers: Array<{
        symbol: string;
        priceChange: number;
      }>;
      topLosers: Array<{
        symbol: string;
        priceChange: number;
      }>;
    };
  };
}

export interface GeminiChatResponse {
  text: string;
  references?: Array<{
    text: string;
    source: string;
  }>;
}

// Initialize the Gemini API client
const API_KEY = import.meta.env.VITE_GEMINI_API_KEY;

// Create a client with the API key
const genAI = new GoogleGenerativeAI(API_KEY || '');

// Get the Gemini model (1.5 Flash)
const model = genAI.getGenerativeModel({
  model: 'gemini-1.5-flash',
  safetySettings: [
    {
      category: HarmCategory.HARM_CATEGORY_HARASSMENT,
      threshold: HarmBlockThreshold.BLOCK_MEDIUM_AND_ABOVE,
    },
    {
      category: HarmCategory.HARM_CATEGORY_HATE_SPEECH,
      threshold: HarmBlockThreshold.BLOCK_MEDIUM_AND_ABOVE,
    },
    {
      category: HarmCategory.HARM_CATEGORY_SEXUALLY_EXPLICIT,
      threshold: HarmBlockThreshold.BLOCK_MEDIUM_AND_ABOVE,
    },
    {
      category: HarmCategory.HARM_CATEGORY_DANGEROUS_CONTENT,
      threshold: HarmBlockThreshold.BLOCK_MEDIUM_AND_ABOVE,
    },
  ],
});

// Create a chat session
const chat = model.startChat({
  history: [
    {
      role: 'user',
      parts: [
        'You are a crypto trading assistant. You help users understand their portfolio, market trends, and provide trading advice. Keep responses concise and focused on crypto trading.',
      ],
    },
    {
      role: 'model',
      parts: [
        'I\'m your crypto trading assistant. I\'ll help you understand your portfolio, analyze market trends, and provide trading insights. How can I assist with your crypto investments today?',
      ],
    },
  ],
  generationConfig: {
    temperature: 0.7,
    topK: 40,
    topP: 0.95,
    maxOutputTokens: 1024,
  },
});

/**
 * Send a message to the Gemini AI model
 * @param request The chat request containing the message and optional context
 * @returns The AI response
 */
export async function sendMessage(request: GeminiChatRequest): Promise<GeminiChatResponse> {
  try {
    // Prepare the message with trading context if provided
    let messageText = request.message;
    
    if (request.tradingContext) {
      messageText += '\n\nTrading Context:';
      
      if (request.tradingContext.portfolio) {
        messageText += `\nPortfolio Total Value: $${request.tradingContext.portfolio.totalValue.toFixed(2)}`;
        messageText += '\nHoldings:';
        request.tradingContext.portfolio.holdings.forEach(holding => {
          messageText += `\n- ${holding.symbol}: ${holding.quantity} ($${holding.value.toFixed(2)})`;
        });
      }
      
      if (request.tradingContext.marketData) {
        messageText += '\nMarket Data:';
        messageText += '\nTop Gainers:';
        request.tradingContext.marketData.topGainers.forEach(coin => {
          messageText += `\n- ${coin.symbol}: ${coin.priceChange > 0 ? '+' : ''}${coin.priceChange.toFixed(2)}%`;
        });
        
        messageText += '\nTop Losers:';
        request.tradingContext.marketData.topLosers.forEach(coin => {
          messageText += `\n- ${coin.symbol}: ${coin.priceChange.toFixed(2)}%`;
        });
      }
    }
    
    // Send the message to the model
    const result = await chat.sendMessage(messageText);
    const response = await result.response;
    const text = response.text();
    
    return {
      text,
    };
  } catch (error) {
    console.error('Error sending message to Gemini:', error);
    throw error;
  }
}

/**
 * Send a message to the Gemini AI model via the API route
 * This is an alternative approach that uses a custom API route
 * @param request The chat request
 * @returns The AI response
 */
export async function sendMessageViaAPI(request: GeminiChatRequest): Promise<GeminiChatResponse> {
  try {
    const response = await fetch('/api/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error sending message via API:', error);
    throw error;
  }
}

export default {
  sendMessage,
  sendMessageViaAPI,
};
