import { useState } from 'react';

// Define types for toast messages
export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface ToastMessage {
  id: string;
  type: ToastType;
  title?: string;
  message: string;
  duration?: number;
}

// Default options for toast messages
const DEFAULT_TOAST_DURATION = 5000; // 5 seconds

// Hook for managing toast messages
export function useToast() {
  const [toasts, setToasts] = useState<ToastMessage[]>([]);

  // Function to add a toast
  const toast = (options: Omit<ToastMessage, 'id'>) => {
    const id = Math.random().toString(36).substring(2, 9);
    const newToast: ToastMessage = {
      id,
      duration: DEFAULT_TOAST_DURATION,
      ...options,
    };

    setToasts((prevToasts) => [...prevToasts, newToast]);

    // Auto-remove toast after duration
    if (newToast.duration !== Infinity) {
      setTimeout(() => {
        dismiss(id);
      }, newToast.duration);
    }

    return id;
  };

  // Function to dismiss a toast
  const dismiss = (id: string) => {
    setToasts((prevToasts) => prevToasts.filter((toast) => toast.id !== id));
  };

  // Function to dismiss all toasts
  const dismissAll = () => {
    setToasts([]);
  };

  // Convenience methods for different toast types
  const success = (message: string, options: Partial<Omit<ToastMessage, 'id' | 'type' | 'message'>> = {}) => 
    toast({ type: 'success', message, ...options });

  const error = (message: string, options: Partial<Omit<ToastMessage, 'id' | 'type' | 'message'>> = {}) => 
    toast({ type: 'error', message, ...options });

  const info = (message: string, options: Partial<Omit<ToastMessage, 'id' | 'type' | 'message'>> = {}) => 
    toast({ type: 'info', message, ...options });

  const warning = (message: string, options: Partial<Omit<ToastMessage, 'id' | 'type' | 'message'>> = {}) => 
    toast({ type: 'warning', message, ...options });

  return {
    toasts,
    toast,
    dismiss,
    dismissAll,
    success,
    error,
    info,
    warning,
  };
}

// Singleton for use outside of React components
const toastMethods = {
  success: (message: string, options: Partial<Omit<ToastMessage, 'id' | 'type' | 'message'>> = {}) => {},
  error: (message: string, options: Partial<Omit<ToastMessage, 'id' | 'type' | 'message'>> = {}) => {},
  info: (message: string, options: Partial<Omit<ToastMessage, 'id' | 'type' | 'message'>> = {}) => {},
  warning: (message: string, options: Partial<Omit<ToastMessage, 'id' | 'type' | 'message'>> = {}) => {},
  dismiss: (id: string) => {},
  dismissAll: () => {},
};

// Export for use outside React components
export const toast = toastMethods; 