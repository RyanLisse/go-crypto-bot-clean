import React, { useState } from 'react';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { ExchangeWalletForm } from './exchange-wallet-form';
import { Web3WalletForm } from './web3-wallet-form';
import { WalletVerificationStep } from './wallet-verification-step';
import { WalletSuccessStep } from './wallet-success-step';

export type WalletType = 'exchange' | 'web3';

export type WalletConnectionStep = 'select' | 'connect' | 'verify' | 'success';

export interface WalletConnectionData {
  type: WalletType;
  exchange?: string;
  apiKey?: string;
  apiSecret?: string;
  network?: string;
  address?: string;
  walletId?: string;
}

export default function WalletConnectionFlow() {
  const [step, setStep] = useState<WalletConnectionStep>('select');
  const [walletData, setWalletData] = useState<WalletConnectionData>({
    type: 'exchange',
  });
  const [activeTab, setActiveTab] = useState<WalletType>('exchange');

  const handleTabChange = (value: string) => {
    setActiveTab(value as WalletType);
    setWalletData({ ...walletData, type: value as WalletType });
  };

  const handleExchangeSubmit = async (data: { exchange: string; apiKey: string; apiSecret: string }) => {
    setWalletData({
      ...walletData,
      type: 'exchange',
      exchange: data.exchange,
      apiKey: data.apiKey,
      apiSecret: data.apiSecret,
    });
    
    try {
      // Call API to connect exchange wallet
      const response = await fetch('/api/wallets/exchange', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          exchange: data.exchange,
          api_key: data.apiKey,
          api_secret: data.apiSecret,
        }),
      });
      
      if (!response.ok) {
        throw new Error('Failed to connect exchange wallet');
      }
      
      const result = await response.json();
      setWalletData({
        ...walletData,
        walletId: result.id,
      });
      
      // Move to verification step
      setStep('verify');
    } catch (error) {
      console.error('Error connecting exchange wallet:', error);
      // Handle error (show error message)
    }
  };

  const handleWeb3Submit = async (data: { network: string; address: string }) => {
    setWalletData({
      ...walletData,
      type: 'web3',
      network: data.network,
      address: data.address,
    });
    
    try {
      // Call API to connect Web3 wallet
      const response = await fetch('/api/wallets/web3', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          network: data.network,
          address: data.address,
        }),
      });
      
      if (!response.ok) {
        throw new Error('Failed to connect Web3 wallet');
      }
      
      const result = await response.json();
      setWalletData({
        ...walletData,
        walletId: result.id,
      });
      
      // Move to verification step
      setStep('verify');
    } catch (error) {
      console.error('Error connecting Web3 wallet:', error);
      // Handle error (show error message)
    }
  };

  const handleVerificationComplete = () => {
    setStep('success');
  };

  const handleReset = () => {
    setStep('select');
    setWalletData({ type: activeTab });
  };

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>Connect Wallet</CardTitle>
        <CardDescription>
          {step === 'select' && 'Choose a wallet type to connect'}
          {step === 'connect' && 'Enter your wallet details'}
          {step === 'verify' && 'Verify your wallet ownership'}
          {step === 'success' && 'Wallet connected successfully'}
        </CardDescription>
      </CardHeader>
      <CardContent>
        {step === 'select' && (
          <Tabs defaultValue="exchange" value={activeTab} onValueChange={handleTabChange}>
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="exchange">Exchange</TabsTrigger>
              <TabsTrigger value="web3">Web3</TabsTrigger>
            </TabsList>
            <TabsContent value="exchange" className="mt-4">
              <p className="text-sm text-muted-foreground mb-4">
                Connect your exchange wallet by providing your API key and secret.
              </p>
              <Button 
                className="w-full" 
                onClick={() => setStep('connect')}
              >
                Connect Exchange Wallet
              </Button>
            </TabsContent>
            <TabsContent value="web3" className="mt-4">
              <p className="text-sm text-muted-foreground mb-4">
                Connect your Web3 wallet by providing your wallet address.
              </p>
              <Button 
                className="w-full" 
                onClick={() => setStep('connect')}
              >
                Connect Web3 Wallet
              </Button>
            </TabsContent>
          </Tabs>
        )}

        {step === 'connect' && activeTab === 'exchange' && (
          <ExchangeWalletForm onSubmit={handleExchangeSubmit} onCancel={handleReset} />
        )}

        {step === 'connect' && activeTab === 'web3' && (
          <Web3WalletForm onSubmit={handleWeb3Submit} onCancel={handleReset} />
        )}

        {step === 'verify' && (
          <WalletVerificationStep 
            walletData={walletData} 
            onComplete={handleVerificationComplete} 
            onCancel={handleReset} 
          />
        )}

        {step === 'success' && (
          <WalletSuccessStep walletData={walletData} onDone={handleReset} />
        )}
      </CardContent>
      {step !== 'connect' && step !== 'verify' && step !== 'success' && (
        <CardFooter className="flex justify-between">
          <p className="text-xs text-muted-foreground">
            Your wallet data is securely stored and encrypted.
          </p>
        </CardFooter>
      )}
    </Card>
  );
}
