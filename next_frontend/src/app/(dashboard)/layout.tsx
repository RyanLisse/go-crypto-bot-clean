import React from 'react';

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <div style={{ padding: 32 }}>
      {/* Dashboard navigation can go here */}
      {children}
    </div>
  );
} 