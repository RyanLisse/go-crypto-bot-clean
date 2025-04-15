import AddressValidatorExample from '@/components/AddressValidatorExample';
import { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Address Validator | Crypto Bot',
  description: 'Validate cryptocurrency wallet addresses across different networks',
};

export default function AddressValidatorPage() {
  return (
    <div className="container py-10">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Address Validator</h1>
        <p className="text-muted-foreground">
          Validate cryptocurrency wallet addresses and get additional information about them.
        </p>
      </div>
      
      <div className="grid gap-8">
        <AddressValidatorExample />
      </div>
    </div>
  );
} 