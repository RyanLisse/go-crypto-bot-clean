'use client';

import { ReactNode, useEffect } from 'react';
import { useRouter } from 'next/navigation';

interface ProtectedRouteProps {
  children: ReactNode;
}

export const ProtectedRoute = ({ children }: ProtectedRouteProps) => {
  const router = useRouter();
  
  useEffect(() => {
    // Check if user is authenticated
    const token = localStorage.getItem('token');
    
    if (!token) {
      // Redirect to login page if not authenticated
      router.push('/login');
    }
  }, [router]);
  
  // You can add additional logic here to display a loading state
  // while checking authentication
  
  return <>{children}</>;
};

export default ProtectedRoute; 