'use client';

import { useState, useEffect } from 'react';
import { useAddressValidator } from '@/hooks/useAddressValidator';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Loader2 } from 'lucide-react';

export default function AddressValidatorExample() {
  const { isValidating, error, validateAddress, getAddressInfo, getSupportedNetworks } = useAddressValidator();
  const [network, setNetwork] = useState<string>('');
  const [address, setAddress] = useState<string>('');
  const [networks, setNetworks] = useState<string[]>([]);
  const [isValid, setIsValid] = useState<boolean | null>(null);
  const [addressInfo, setAddressInfo] = useState<any>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    async function fetchNetworks() {
      try {
        const supportedNetworks = await getSupportedNetworks();
        setNetworks(supportedNetworks);
        if (supportedNetworks.length > 0) {
          setNetwork(supportedNetworks[0]);
        }
      } catch (error) {
        console.error('Failed to fetch networks:', error);
      } finally {
        setLoading(false);
      }
    }

    fetchNetworks();
  }, [getSupportedNetworks]);

  const handleValidate = async () => {
    if (!network || !address) return;
    
    const isAddressValid = await validateAddress(network, address);
    setIsValid(isAddressValid);
    
    if (isAddressValid) {
      const info = await getAddressInfo(network, address);
      setAddressInfo(info);
    } else {
      setAddressInfo(null);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[200px]">
        <Loader2 className="w-6 h-6 animate-spin" />
      </div>
    );
  }

  return (
    <Card className="w-full max-w-md mx-auto">
      <CardHeader>
        <CardTitle>Address Validator</CardTitle>
        <CardDescription>Validate cryptocurrency addresses across different networks</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="network">Network</Label>
          <Select 
            value={network} 
            onValueChange={(value) => {
              setNetwork(value);
              setIsValid(null);
              setAddressInfo(null);
            }}
          >
            <SelectTrigger id="network">
              <SelectValue placeholder="Select a network" />
            </SelectTrigger>
            <SelectContent>
              {networks.map((networkName) => (
                <SelectItem key={networkName} value={networkName}>
                  {networkName}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        
        <div className="space-y-2">
          <Label htmlFor="address">Wallet Address</Label>
          <Input 
            id="address" 
            value={address} 
            onChange={(e) => {
              setAddress(e.target.value);
              setIsValid(null);
              setAddressInfo(null);
            }} 
            placeholder="Enter a wallet address" 
          />
        </div>
        
        {error && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
        
        {isValid !== null && (
          <Alert variant={isValid ? "default" : "destructive"}>
            <AlertDescription>
              {isValid ? "Address is valid" : "Address is invalid"}
            </AlertDescription>
          </Alert>
        )}
        
        {addressInfo && (
          <div className="space-y-2 p-4 bg-muted rounded-md text-sm">
            <p><strong>Network:</strong> {addressInfo.network}</p>
            <p><strong>Type:</strong> {addressInfo.addressType || 'Unknown'}</p>
            {addressInfo.chainId && <p><strong>Chain ID:</strong> {addressInfo.chainId}</p>}
            {addressInfo.explorer && (
              <p>
                <strong>Explorer:</strong>{' '}
                <a 
                  href={`${addressInfo.explorer}${address}`} 
                  target="_blank" 
                  rel="noopener noreferrer"
                  className="text-primary underline"
                >
                  View on Explorer
                </a>
              </p>
            )}
          </div>
        )}
      </CardContent>
      <CardFooter>
        <Button 
          onClick={handleValidate} 
          disabled={!network || !address || isValidating}
          className="w-full"
        >
          {isValidating ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" /> 
              Validating...
            </>
          ) : (
            'Validate Address'
          )}
        </Button>
      </CardFooter>
    </Card>
  );
} 