import React, { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Textarea } from '@/components/ui/textarea';
import { Input } from '@/components/ui/input';
import { WalletConnectionData } from './wallet-connection-flow';
import { Loader2, CheckCircle, Copy, AlertCircle } from 'lucide-react';

type WalletVerificationStepProps = {
  walletData: WalletConnectionData;
  onComplete: () => void;
  onCancel: () => void;
};

export function WalletVerificationStep({ walletData, onComplete, onCancel }: WalletVerificationStepProps) {
  const [challenge, setChallenge] = useState<string>('');
  const [signature, setSignature] = useState<string>('');
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [isVerified, setIsVerified] = useState<boolean>(false);

  useEffect(() => {
    if (walletData.type === 'web3' && walletData.walletId) {
      generateChallenge();
    } else if (walletData.type === 'exchange') {
      // Exchange wallets don't need verification
      onComplete();
    }
  }, [walletData]);

  const generateChallenge = async () => {
    if (!walletData.walletId) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/api/wallet-verification/challenge/${walletData.walletId}`, {
        method: 'POST',
      });
      
      if (!response.ok) {
        throw new Error('Failed to generate challenge');
      }
      
      const data = await response.json();
      setChallenge(data.challenge);
    } catch (error) {
      console.error('Error generating challenge:', error);
      setError('Failed to generate challenge. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const verifySignature = async () => {
    if (!walletData.walletId || !challenge || !signature) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/api/wallet-verification/verify/${walletData.walletId}`, {
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
        throw new Error('Failed to verify signature');
      }
      
      const data = await response.json();
      
      if (data.verified) {
        setIsVerified(true);
        setTimeout(() => {
          onComplete();
        }, 1500);
      } else {
        setError('Signature verification failed. Please try again.');
      }
    } catch (error) {
      console.error('Error verifying signature:', error);
      setError('Failed to verify signature. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  if (walletData.type === 'exchange') {
    return (
      <div className="flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <Alert>
        <AlertCircle className="h-4 w-4" />
        <AlertTitle>Verification Required</AlertTitle>
        <AlertDescription>
          Please sign the message below with your wallet to verify ownership.
        </AlertDescription>
      </Alert>

      {isVerified ? (
        <Alert>
          <CheckCircle className="h-4 w-4" />
          <AlertTitle>Verification Successful</AlertTitle>
          <AlertDescription>
            Your wallet has been verified successfully.
          </AlertDescription>
        </Alert>
      ) : (
        <>
          <Card>
            <CardContent className="pt-6">
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <p className="text-sm font-medium">Sign this message:</p>
                  <Button 
                    variant="ghost" 
                    size="sm" 
                    onClick={() => copyToClipboard(challenge)}
                    disabled={!challenge}
                  >
                    <Copy className="h-4 w-4 mr-2" />
                    Copy
                  </Button>
                </div>
                <Textarea 
                  value={challenge} 
                  readOnly 
                  className="h-32 font-mono text-xs"
                />
              </div>
            </CardContent>
          </Card>

          <div className="space-y-2">
            <label htmlFor="signature" className="text-sm font-medium">
              Paste your signature:
            </label>
            <Input
              id="signature"
              value={signature}
              onChange={(e) => setSignature(e.target.value)}
              placeholder="0x..."
              className="font-mono"
            />
          </div>

          {error && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <div className="flex justify-between">
            <Button variant="outline" onClick={onCancel} disabled={isLoading}>
              Cancel
            </Button>
            <Button 
              onClick={verifySignature} 
              disabled={!challenge || !signature || isLoading}
            >
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Verify Signature
            </Button>
          </div>
        </>
      )}
    </div>
  );
}
