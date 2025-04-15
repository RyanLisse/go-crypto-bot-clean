import React from 'react';
import Link from 'next/link';

const navLinks = [
  { href: '/(dashboard)', label: 'Home' },
  { href: '/(dashboard)/portfolio', label: 'Portfolio' },
  { href: '/(dashboard)/trading', label: 'Trading' },
  { href: '/(dashboard)/new-coins', label: 'New Coins' },
  { href: '/(dashboard)/backtesting', label: 'Backtesting' },
  { href: '/(dashboard)/system', label: 'System' },
  { href: '/(dashboard)/config', label: 'Config' },
  { href: '/(dashboard)/settings', label: 'Settings' },
  { href: '/(dashboard)/testing', label: 'Testing' },
];

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <div style={{ padding: 32 }}>
      <nav style={{ marginBottom: 24 }}>
        {navLinks.map((link) => (
          <Link key={link.href} href={link.href} style={{ marginRight: 16 }}>
            {link.label}
          </Link>
        ))}
      </nav>
      {/* Dashboard navigation can go here */}
      {children}
    </div>
  );
} 