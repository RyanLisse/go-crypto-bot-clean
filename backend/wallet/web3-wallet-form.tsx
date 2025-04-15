import React, { useState } from 'react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import * as z from 'zod';
import { Button } from '@/components/ui/button';
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { CheckCircle, AlertCircle } from 'lucide-react';

const formSchema = z.object({
  network: z.string().min(1, { message: 'Please select a network' }),
  address: z.string().min(1, { message: 'Wallet address is required' }),
});

type Web3WalletFormProps = {
  onSubmit: (data: z.infer<typeof formSchema>) => void;
  onCancel: () => void;
};

export function Web3WalletForm({ onSubmit, onCancel }: Web3WalletFormProps) {
  const [isValidating, setIsValidating] = useState(false);
  const [validationResult, setValidationResult] = useState<{ valid: boolean; message: string } | null>(null);
  
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      network: '',
      address: '',
    },
  });

  const validateAddress = async (network: string, address: string) => {
    if (!network || !address) return;
    
    setIsValidating(true);
    setValidationResult(null);
    
    try {
      const response = await fetch(`/api/address-validator/validate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ network, address }),
      });
      
      const data = await response.json();
      
      if (response.ok) {
        setValidationResult({
          valid: data.valid,
          message: data.valid 
            ? 'Address is valid' 
            : 'Address is not valid for the selected network',
        });
      } else {
        setValidationResult({
          valid: false,
          message: data.message || 'Failed to validate address',
        });
      }
    } catch (error) {
      setValidationResult({
        valid: false,
        message: 'Error validating address',
      });
    } finally {
      setIsValidating(false);
    }
  };

  const handleAddressBlur = () => {
    const network = form.getValues('network');
    const address = form.getValues('address');
    validateAddress(network, address);
  };

  const handleNetworkChange = (value: string) => {
    form.setValue('network', value);
    const address = form.getValues('address');
    if (address) {
      validateAddress(value, address);
    }
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="network"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Network</FormLabel>
              <Select onValueChange={handleNetworkChange} defaultValue={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Select a network" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="Ethereum">Ethereum</SelectItem>
                  <SelectItem value="Bitcoin">Bitcoin</SelectItem>
                  <SelectItem value="Polygon">Polygon</SelectItem>
                  <SelectItem value="Arbitrum">Arbitrum</SelectItem>
                  <SelectItem value="Optimism">Optimism</SelectItem>
                  <SelectItem value="Avalanche">Avalanche</SelectItem>
                </SelectContent>
              </Select>
              <FormDescription>
                Select the blockchain network of your wallet
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="address"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Wallet Address</FormLabel>
              <FormControl>
                <Input 
                  placeholder="Enter your wallet address" 
                  {...field} 
                  onBlur={handleAddressBlur}
                />
              </FormControl>
              <FormDescription>
                Your public wallet address on the selected network
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        
        {validationResult && (
          <Alert variant={validationResult.valid ? "default" : "destructive"}>
            {validationResult.valid ? (
              <CheckCircle className="h-4 w-4" />
            ) : (
              <AlertCircle className="h-4 w-4" />
            )}
            <AlertTitle>
              {validationResult.valid ? "Valid Address" : "Invalid Address"}
            </AlertTitle>
            <AlertDescription>
              {validationResult.message}
            </AlertDescription>
          </Alert>
        )}
        
        <div className="flex justify-between">
          <Button type="button" variant="outline" onClick={onCancel}>
            Back
          </Button>
          <Button 
            type="submit" 
            disabled={isValidating || (validationResult && !validationResult.valid)}
          >
            {isValidating ? 'Validating...' : 'Connect Wallet'}
          </Button>
        </div>
      </form>
    </Form>
  );
}
