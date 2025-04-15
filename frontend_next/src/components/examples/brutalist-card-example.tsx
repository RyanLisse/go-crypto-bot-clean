'use client';

import React from 'react';
import { 
  BrutalistCard,
  BrutalistCardHeader,
  BrutalistCardTitle,
  BrutalistCardDescription,
  BrutalistCardContent,
  BrutalistCardFooter 
} from '../ui/brutalist-card';
import { Button } from '../ui/button';
import { ArrowUp, ArrowDown, TrendingUp } from 'lucide-react';

export function BrutalistCardExample() {
  const cryptoData = [
    { name: 'Bitcoin', symbol: 'BTC', price: 52341.23, change: 2.5, positive: true },
    { name: 'Ethereum', symbol: 'ETH', price: 2856.17, change: -1.2, positive: false },
    { name: 'Solana', symbol: 'SOL', price: 137.84, change: 5.7, positive: true },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 my-6">
      {cryptoData.map((crypto) => (
        <BrutalistCard key={crypto.symbol}>
          <BrutalistCardHeader>
            <BrutalistCardTitle>{crypto.name}</BrutalistCardTitle>
            <BrutalistCardDescription>{crypto.symbol}</BrutalistCardDescription>
          </BrutalistCardHeader>
          <BrutalistCardContent>
            <div className="flex flex-col space-y-2">
              <div className="text-3xl font-bold">${crypto.price.toLocaleString()}</div>
              <div className={`flex items-center ${crypto.positive ? 'text-green-600' : 'text-red-600'}`}>
                {crypto.positive ? <ArrowUp className="mr-1 h-4 w-4" /> : <ArrowDown className="mr-1 h-4 w-4" />}
                <span>{Math.abs(crypto.change)}%</span>
              </div>
            </div>
          </BrutalistCardContent>
          <BrutalistCardFooter>
            <Button className="w-full" variant={crypto.positive ? 'default' : 'outline'}>
              <TrendingUp className="mr-2 h-4 w-4" /> View Analytics
            </Button>
          </BrutalistCardFooter>
        </BrutalistCard>
      ))}
    </div>
  );
} 