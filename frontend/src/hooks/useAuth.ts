import { createContext, useContext } from 'react';

// Define a simple auth context type
type AuthContextType = {
  user: {
    id: string;
    email: string;
    name?: string;
  } | null;
};

// Create the auth context with a default value
const AuthContext = createContext<AuthContextType>({
  user: {
    id: 'user123',
    email: 'demo@example.com',
    name: 'Demo User',
  }
});

// Simple Auth provider component
export function AuthProvider({ children }: { children: React.ReactNode }) {
  // For simplicity, we're using a mock user
  const authValue = {
    user: {
      id: 'user123',
      email: 'demo@example.com',
      name: 'Demo User',
    }
  };

  return (
    <AuthContext.Provider value={authValue}>
      {children}
    </AuthContext.Provider>
  );
}

// Hook to use auth context
export function useAuth() {
  return useContext(AuthContext);
}
