import { LoginResponse, User } from '@/types';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export const authService = {
  /**
   * Login with username and password
   * @param username User's username
   * @param password User's password
   * @returns User data including token
   */
  async login(username: string, password: string): Promise<LoginResponse> {
    const response = await fetch(`${API_URL}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username,
        password
      }),
    });

    if (!response.ok) {
      throw new Error(`Login failed: ${response.status} ${response.statusText}`);
    }

    const data = await response.json();

    // Store token and user data in localStorage
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data));

    return data;
  },

  /**
   * Logout the current user
   */
  async logout(): Promise<void> {
    try {
      // Call the logout endpoint
      await fetch(`${API_URL}/auth/logout`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...this.getAuthHeader(),
        },
      });
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      // Always clear local storage
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    }
  },

  /**
   * Get the current user from localStorage
   * @returns User data or null if not logged in
   */
  getCurrentUser(): User | null {
    const userJson = localStorage.getItem('user');
    return userJson ? JSON.parse(userJson) : null;
  },

  /**
   * Check if the user is authenticated
   * @returns True if authenticated, false otherwise
   */
  isAuthenticated(): boolean {
    return !!localStorage.getItem('token');
  },

  /**
   * Get the authorization header for API requests
   * @returns Authorization header object or empty object
   */
  getAuthHeader(): { Authorization?: string } {
    const token = localStorage.getItem('token');
    return token ? { Authorization: `Bearer ${token}` } : {};
  }
};

export default authService;
