import { useState } from 'react';

export interface AddressInfo {
  network: string;
  address: string;
  isValid: boolean;
  addressType?: string;
  chainId?: number;
  explorer?: string;
}

export interface UseAddressValidatorReturn {
  isValidating: boolean;
  error: string | null;
  validateAddress: (network: string, address: string) => Promise<boolean>;
  getAddressInfo: (network: string, address: string) => Promise<AddressInfo | null>;
  getSupportedNetworks: () => Promise<string[]>;
}

export function useAddressValidator(): UseAddressValidatorReturn {
  const [isValidating, setIsValidating] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const validateAddress = async (network: string, address: string): Promise<boolean> => {
    setIsValidating(true);
    setError(null);
    
    try {
      const response = await fetch('/api/address-validator/validate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ network, address }),
      });
      
      if (!response.ok) {
        throw new Error('Failed to validate address');
      }
      
      const data = await response.json();
      return data.valid;
    } catch (error) {
      console.error('Error validating address:', error);
      setError('Failed to validate address. Please try again.');
      return false;
    } finally {
      setIsValidating(false);
    }
  };

  const getAddressInfo = async (network: string, address: string): Promise<AddressInfo | null> => {
    setIsValidating(true);
    setError(null);
    
    try {
      const response = await fetch('/api/address-validator/info', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ network, address }),
      });
      
      if (!response.ok) {
        throw new Error('Failed to get address info');
      }
      
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error getting address info:', error);
      setError('Failed to get address info. Please try again.');
      return null;
    } finally {
      setIsValidating(false);
    }
  };

  const getSupportedNetworks = async (): Promise<string[]> => {
    setIsValidating(true);
    setError(null);
    
    try {
      const response = await fetch('/api/address-validator/networks');
      
      if (!response.ok) {
        throw new Error('Failed to get supported networks');
      }
      
      const data = await response.json();
      return data.networks || [];
    } catch (error) {
      console.error('Error getting supported networks:', error);
      setError('Failed to get supported networks. Please try again.');
      return [];
    } finally {
      setIsValidating(false);
    }
  };

  return {
    isValidating,
    error,
    validateAddress,
    getAddressInfo,
    getSupportedNetworks,
  };
}
