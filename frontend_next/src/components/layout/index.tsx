import React from 'react';
import { Sidebar } from '@/components/layout/Sidebar';
import { Header } from '@/components/layout/Header';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <>
      <Header />
      <Sidebar />
      <main style={{ flexGrow: 1, padding: 24, width: '100%' }}>
        {children}
      </main>
    </>
  );
};

export { Layout };
