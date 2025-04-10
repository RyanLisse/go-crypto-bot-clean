import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";

interface TradingPairSelectorProps {
  availablePairs: string[];
  selectedPair: string;
  onChange: (pair: string) => void;
  isLoading?: boolean;
}

const TradingPairSelector: React.FC<TradingPairSelectorProps> = ({
  availablePairs,
  selectedPair,
  onChange,
  isLoading = false,
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  
  const filteredPairs = availablePairs.filter(pair => 
    pair.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handlePairClick = (pair: string) => {
    if (pair !== selectedPair) {
      onChange(pair);
    }
  };

  return (
    <Card className="bg-brutal-card-bg border-brutal-border">
      <CardHeader className="pb-2">
        <CardTitle className="text-brutal-text text-lg">Trading Pairs</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          <Input
            className="bg-brutal-input-bg border-brutal-border text-brutal-text"
            placeholder="Search pairs..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          
          <div className="max-h-[300px] overflow-y-auto pr-1">
            {isLoading ? (
              <div className="flex items-center justify-center py-4">
                <Loader2 className="h-6 w-6 animate-spin text-brutal-primary" />
                <span className="ml-2 text-brutal-text">Loading trading pairs...</span>
              </div>
            ) : filteredPairs.length > 0 ? (
              <div className="space-y-1">
                {filteredPairs.map((pair) => (
                  <button
                    key={pair}
                    className={cn(
                      "w-full text-left px-3 py-2 rounded-md transition-colors",
                      pair === selectedPair
                        ? "bg-brutal-active-bg text-brutal-active-text"
                        : "bg-brutal-card-bg hover:bg-brutal-hover-bg text-brutal-text"
                    )}
                    onClick={() => handlePairClick(pair)}
                    data-testid={pair === selectedPair ? "selected-pair" : "pair-option"}
                  >
                    {pair}
                  </button>
                ))}
              </div>
            ) : (
              <div className="text-center py-4 text-brutal-text-muted">
                No trading pairs available
              </div>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default TradingPairSelector; 