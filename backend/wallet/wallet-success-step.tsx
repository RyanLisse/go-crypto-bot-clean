import React from 'react';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { WalletConnectionData } from './wallet-connection-flow';
import { CheckCircle } from 'lucide-react';
import Link from 'next/link';

type WalletSuccessStepProps = {
  walletData: WalletConnectionData;
  onDone: () => void;
};

export function WalletSuccessStep({ walletData, onDone }: WalletSuccessStepProps) {
  return (
    <div className="space-y-6">
      <Alert className="border-green-500 bg-green-50 dark:bg-green-950">
        <CheckCircle className="h-5 w-5 text-green-500" />
        <AlertTitle className="text-green-700 dark:text-green-300">Success!</AlertTitle>
        <AlertDescription className="text-green-700 dark:text-green-300">
          Your wallet has been successfully connected.
        </AlertDescription>
      </Alert>

      <div className="space-y-4">
        <div>
          <h3 className="text-lg font-medium">Wallet Details</h3>
          <div className="mt-2 space-y-2">
            <div className="flex justify-between">
              <span className="text-sm text-muted-foreground">Type:</span>
              <span className="text-sm font-medium">
                {walletData.type === 'exchange' ? 'Exchange' : 'Web3'}
              </span>
            </div>
            
            {walletData.type === 'exchange' && walletData.exchange && (
              <div className="flex justify-between">
                <span className="text-sm text-muted-foreground">Exchange:</span>
                <span className="text-sm font-medium">{walletData.exchange}</span>
              </div>
            )}
            
            {walletData.type === 'web3' && (
              <>
                {walletData.network && (
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">Network:</span>
                    <span className="text-sm font-medium">{walletData.network}</span>
                  </div>
                )}
                {walletData.address && (
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">Address:</span>
                    <span className="text-sm font-medium truncate max-w-[200px]">
                      {walletData.address}
                    </span>
                  </div>
                )}
              </>
            )}
            
            <div className="flex justify-between">
              <span className="text-sm text-muted-foreground">Status:</span>
              <span className="text-sm font-medium text-green-600">Connected</span>
            </div>
          </div>
        </div>

        <div className="space-y-2">
          <p className="text-sm text-muted-foreground">
            You can now use this wallet to interact with the platform.
          </p>
          <p className="text-sm text-muted-foreground">
            View your wallet details and balances in the portfolio section.
          </p>
        </div>
      </div>

      <div className="flex justify-between">
        <Button variant="outline" onClick={onDone}>
          Connect Another Wallet
        </Button>
        <Button asChild>
          <Link href="/portfolio">View Portfolio</Link>
        </Button>
      </div>
    </div>
  );
}
