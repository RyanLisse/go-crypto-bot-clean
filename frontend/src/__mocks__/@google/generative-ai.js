// Mock implementation of Google Generative AI
export class GoogleGenerativeAI {
  constructor(apiKey) {
    this.apiKey = apiKey;
  }

  getGenerativeModel({ model, safetySettings }) {
    return {
      startChat: ({ history, generationConfig }) => {
        return {
          sendMessage: async (message) => {
            return {
              response: {
                text: () => "This is a mock response from the AI model."
              }
            };
          }
        };
      }
    };
  }
}

// Mock harm categories and thresholds
export const HarmCategory = {
  HARM_CATEGORY_HARASSMENT: 'HARM_CATEGORY_HARASSMENT',
  HARM_CATEGORY_HATE_SPEECH: 'HARM_CATEGORY_HATE_SPEECH',
  HARM_CATEGORY_SEXUALLY_EXPLICIT: 'HARM_CATEGORY_SEXUALLY_EXPLICIT',
  HARM_CATEGORY_DANGEROUS_CONTENT: 'HARM_CATEGORY_DANGEROUS_CONTENT',
};

export const HarmBlockThreshold = {
  BLOCK_MEDIUM_AND_ABOVE: 'BLOCK_MEDIUM_AND_ABOVE',
};