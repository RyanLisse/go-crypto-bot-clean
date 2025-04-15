import { useState } from 'react';

export type WalletType = 'exchange' | 'web3';

export interface Wallet {
  id: string;
  type: WalletType;
  exchange?: string;
  network?: string;
  address?: string;
  status: 'pending' | 'active' | 'inactive' | 'verified' | 'failed';
  balances?: Record<string, {
    free: string;
    locked: string;
    total: string;
  }>;
  totalUSDValue?: string;
  lastUpdated?: string;
}

export interface UseWalletReturn {
  wallets: Wallet[];
  isLoading: boolean;
  error: string | null;
  fetchWallets: () => Promise<void>;
  connectExchangeWallet: (exchange: string, apiKey: string, apiSecret: string) => Promise<Wallet | null>;
  connectWeb3Wallet: (network: string, address: string) => Promise<Wallet | null>;
  verifyWallet: (walletId: string, challenge: string, signature: string) => Promise<boolean>;
  generateChallenge: (walletId: string) => Promise<string | null>;
  getWalletStatus: (walletId: string) => Promise<string | null>;
  disconnectWallet: (walletId: string) => Promise<boolean>;
  refreshWalletBalance: (walletId: string) => Promise<Wallet | null>;
}

export function useWallet(): UseWalletReturn {
  const [wallets, setWallets] = useState<Wallet[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const fetchWallets = async (): Promise<void> => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await fetch('/api/wallets');
      
      if (!response.ok) {
        throw new Error('Failed to fetch wallets');
      }
      
      const data = await response.json();
      setWallets(data.wallets || []);
    } catch (error) {
      console.error('Error fetching wallets:', error);
      setError('Failed to fetch wallets. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const connectExchangeWallet = async (
    exchange: string, 
    apiKey: string, 
    apiSecret: string
  ): Promise<Wallet | null> => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await fetch('/api/wallets/exchange', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          exchange,
          api_key: apiKey,
          api_secret: apiSecret,
        }),
      });
      
      if (!response.ok) {
        throw new Error('Failed to connect exchange wallet');
      }
      
      const wallet = await response.json();
      
      // Update wallets list
      setWallets(prev => [...prev, wallet]);
      
      return wallet;
    } catch (error) {
      console.error('Error connecting exchange wallet:', error);
      setError('Failed to connect exchange wallet. Please try again.');
      return null;
    } finally {
      setIsLoading(false);
    }
  };

  const connectWeb3Wallet = async (
    network: string, 
    address: string
  ): Promise<Wallet | null> => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await fetch('/api/wallets/web3', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          network,
          address,
        }),
      });
      
      if (!response.ok) {
        throw new Error('Failed to connect Web3 wallet');
      }
      
      const wallet = await response.json();
      
      // Update wallets list
      setWallets(prev => [...prev, wallet]);
      
      return wallet;
    } catch (error) {
      console.error('Error connecting Web3 wallet:', error);
      setError('Failed to connect Web3 wallet. Please try again.');
      return null;
    } finally {
      setIsLoading(false);
    }
  };

  const generateChallenge = async (walletId: string): Promise<string | null> => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/api/wallet-verification/challenge/${walletId}`, {
        method: 'POST',
      });
      
      if (!response.ok) {
        throw new Error('Failed to generate challenge');
      }
      
      const data = await response.json();
      return data.challenge;
    } catch (error) {
      console.error('Error generating challenge:', error);
      setError('Failed to generate challenge. Please try again.');
      return null;
    } finally {
      setIsLoading(false);
    }
  };

  const verifyWallet = async (
    walletId: string, 
    challenge: string, 
    signature: string
  ): Promise<boolean> => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/api/wallet-verification/verify/${walletId}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          challenge,
          signature,
        }),
      });
      
      if (!response.ok) {
        throw new Error('Failed to verify wallet');
      }
      
      const data = await response.json();
      
      if (data.verified) {
        // Update wallet status in the list
        setWallets(prev => 
          prev.map(wallet => 
            wallet.id === walletId 
              ? { ...wallet, status: 'verified' } 
              : wallet
          )
        );
      }
      
      return data.verified;
    } catch (error) {
      console.error('Error verifying wallet:', error);
      setError('Failed to verify wallet. Please try again.');
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  const getWalletStatus = async (walletId: string): Promise<string | null> => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/api/wallet-verification/status/${walletId}`);
      
      if (!response.ok) {
        throw new Error('Failed to get wallet status');
      }
      
      const data = await response.json();
      return data.status;
    } catch (error) {
      console.error('Error getting wallet status:', error);
      setError('Failed to get wallet status. Please try again.');
      return null;
    } finally {
      setIsLoading(false);
    }
  };

  const disconnectWallet = async (walletId: string): Promise<boolean> => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/api/wallets/${walletId}`, {
        method: 'DELETE',
      });
      
      if (!response.ok) {
        throw new Error('Failed to disconnect wallet');
      }
      
      // Remove wallet from the list
      setWallets(prev => prev.filter(wallet => wallet.id !== walletId));
      
      return true;
    } catch (error) {
      console.error('Error disconnecting wallet:', error);
      setError('Failed to disconnect wallet. Please try again.');
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  const refreshWalletBalance = async (walletId: string): Promise<Wallet | null> => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/api/wallets/${walletId}/balance`, {
        method: 'POST',
      });
      
      if (!response.ok) {
        throw new Error('Failed to refresh wallet balance');
      }
      
      const updatedWallet = await response.json();
      
      // Update wallet in the list
      setWallets(prev => 
        prev.map(wallet => 
          wallet.id === walletId ? updatedWallet : wallet
        )
      );
      
      return updatedWallet;
    } catch (error) {
      console.error('Error refreshing wallet balance:', error);
      setError('Failed to refresh wallet balance. Please try again.');
      return null;
    } finally {
      setIsLoading(false);
    }
  };

  return {
    wallets,
    isLoading,
    error,
    fetchWallets,
    connectExchangeWallet,
    connectWeb3Wallet,
    verifyWallet,
    generateChallenge,
    getWalletStatus,
    disconnectWallet,
    refreshWalletBalance,
  };
}
