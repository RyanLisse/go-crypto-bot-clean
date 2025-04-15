import React from 'react';
import { Metadata } from 'next';
import WalletConnectionFlow from '@/components/wallet/wallet-connection-flow';

export const metadata: Metadata = {
  title: 'Connect Wallet',
  description: 'Connect your wallet to the platform',
};

export default function WalletConnectionPage() {
  return (
    <div className="container mx-auto py-10">
      <div className="flex flex-col items-center space-y-6">
        <div className="text-center space-y-2">
          <h1 className="text-3xl font-bold tracking-tight">Connect Your Wallet</h1>
          <p className="text-muted-foreground">
            Connect your wallet to access all features of the platform
          </p>
        </div>
        <WalletConnectionFlow />
      </div>
    </div>
  );
}
