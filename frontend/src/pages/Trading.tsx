
import React, { useState } from 'react';
import { Header } from '@/components/layout/Header';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Search, TrendingUp, ArrowRightCircle } from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';

const Trading = () => {
  const [searchQuery, setSearchQuery] = useState('');
  
  // Mock data for available coins
  const availableCoins = [
    { 
      name: 'Bitcoin',
      symbol: 'BTC',
      price: 58432.21,
      volume24h: 24541530000,
      volumeChange: 23.5,
      marketCap: 1342000000000
    },
    { 
      name: 'Ethereum',
      symbol: 'ETH',
      price: 2843.67,
      volume24h: 15331230000,
      volumeChange: 14.2,
      marketCap: 342000000000
    },
    { 
      name: 'Solana',
      symbol: 'SOL',
      price: 142.86,
      volume24h: 3751330000,
      volumeChange: 32.8,
      marketCap: 62500000000
    },
    { 
      name: 'Binance Coin',
      symbol: 'BNB',
      price: 563.21,
      volume24h: 1895520000,
      volumeChange: 8.4,
      marketCap: 88300000000
    },
    { 
      name: 'Cardano',
      symbol: 'ADA',
      price: 0.89,
      volume24h: 1020330000,
      volumeChange: -5.2,
      marketCap: 31400000000
    },
  ];
  
  // Filter coins based on search query
  const filteredCoins = availableCoins.filter(coin => 
    coin.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
    coin.symbol.toLowerCase().includes(searchQuery.toLowerCase())
  );
  
  // Format large numbers with commas
  const formatNumber = (num: number) => {
    return num.toLocaleString();
  };
  
  return (
    <div className="flex-1 flex flex-col h-full overflow-auto">
      <Header />
      
      <div className="flex-1 p-4 md:p-6 space-y-4 md:space-y-6">
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-3">
          <h1 className="text-2xl font-bold text-brutal-text tracking-tight">TRADING</h1>
          
          <div className="w-full md:w-auto flex items-center gap-2">
            <div className="relative flex-1 md:w-64">
              <Search className="absolute left-2 top-2.5 h-4 w-4 text-brutal-text/50" />
              <Input
                placeholder="Search coins..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-8 w-full bg-brutal-panel border-brutal-border"
              />
            </div>
            
            <Button variant="default" className="bg-brutal-info text-white hover:bg-brutal-info/80">
              <TrendingUp className="mr-2 h-4 w-4" />
              Volume Filter
            </Button>
          </div>
        </div>
        
        <Card className="bg-brutal-panel border-brutal-border">
          <CardHeader className="pb-2">
            <CardTitle className="text-brutal-text flex items-center text-lg">
              <TrendingUp className="mr-2 h-5 w-5 text-brutal-info" />
              Available Coins
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow className="border-brutal-border">
                    <TableHead className="text-brutal-text/70">Coin</TableHead>
                    <TableHead className="text-brutal-text/70 text-right">Price</TableHead>
                    <TableHead className="text-brutal-text/70 text-right">24h Volume</TableHead>
                    <TableHead className="text-brutal-text/70 text-right">Volume Change</TableHead>
                    <TableHead className="text-brutal-text/70 text-right">Market Cap</TableHead>
                    <TableHead className="text-brutal-text/70 text-right">Action</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredCoins.map((coin) => (
                    <TableRow key={coin.symbol} className="border-brutal-border">
                      <TableCell className="font-medium text-brutal-text">
                        <div className="flex items-center">
                          <div className="w-6 h-6 rounded-full bg-brutal-info/20 mr-2 flex items-center justify-center text-xs">
                            {coin.symbol.substring(0, 1)}
                          </div>
                          <div>
                            <div>{coin.symbol}</div>
                            <div className="text-xs text-brutal-text/70">{coin.name}</div>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell className="text-right font-mono text-brutal-text">
                        ${coin.price.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                      </TableCell>
                      <TableCell className="text-right font-mono text-brutal-text">
                        ${formatNumber(coin.volume24h)}
                      </TableCell>
                      <TableCell className={`text-right font-mono ${coin.volumeChange >= 0 ? 'text-brutal-success' : 'text-brutal-error'}`}>
                        {coin.volumeChange >= 0 ? '+' : ''}{coin.volumeChange}%
                      </TableCell>
                      <TableCell className="text-right font-mono text-brutal-text">
                        ${formatNumber(coin.marketCap)}
                      </TableCell>
                      <TableCell className="text-right">
                        <Button 
                          size="sm"
                          variant="outline"
                          className="border-brutal-border hover:bg-brutal-info hover:text-white"
                        >
                          <ArrowRightCircle className="h-4 w-4 mr-1" />
                          Buy
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card className="bg-brutal-panel border-brutal-border">
            <CardHeader className="pb-2">
              <CardTitle className="text-brutal-text text-lg">
                Volume Alerts
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="p-3 border border-brutal-border rounded">
                  <div className="flex justify-between items-center">
                    <div className="flex items-center">
                      <div className="w-6 h-6 rounded-full bg-brutal-info/20 mr-2 flex items-center justify-center text-xs">
                        S
                      </div>
                      <div className="font-mono">
                        <div>SOL</div>
                        <div className="text-xs text-brutal-text/70">Volume +32.8%</div>
                      </div>
                    </div>
                    <Button 
                      size="sm"
                      className="bg-brutal-success text-white hover:bg-brutal-success/80"
                    >
                      Buy
                    </Button>
                  </div>
                </div>
                
                <div className="p-3 border border-brutal-border rounded">
                  <div className="flex justify-between items-center">
                    <div className="flex items-center">
                      <div className="w-6 h-6 rounded-full bg-brutal-info/20 mr-2 flex items-center justify-center text-xs">
                        B
                      </div>
                      <div className="font-mono">
                        <div>BTC</div>
                        <div className="text-xs text-brutal-text/70">Volume +23.5%</div>
                      </div>
                    </div>
                    <Button 
                      size="sm"
                      className="bg-brutal-success text-white hover:bg-brutal-success/80"
                    >
                      Buy
                    </Button>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
          
          <Card className="bg-brutal-panel border-brutal-border">
            <CardHeader className="pb-2">
              <CardTitle className="text-brutal-text text-lg">
                Manual Buy
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <label className="text-xs text-brutal-text/70 mb-1 block">
                    Coin Symbol
                  </label>
                  <Input 
                    placeholder="Enter coin symbol (e.g. BTC)" 
                    className="bg-brutal-background border-brutal-border"
                  />
                </div>
                
                <div>
                  <label className="text-xs text-brutal-text/70 mb-1 block">
                    Amount (USD)
                  </label>
                  <Input
                    type="number"
                    placeholder="Enter amount to buy"
                    className="bg-brutal-background border-brutal-border"
                  />
                </div>
                
                <div>
                  <label className="text-xs text-brutal-text/70 mb-1 block">
                    Stop Loss (%)
                  </label>
                  <Input
                    type="number"
                    placeholder="Enter stop loss percentage"
                    className="bg-brutal-background border-brutal-border"
                  />
                </div>
                
                <div>
                  <label className="text-xs text-brutal-text/70 mb-1 block">
                    Take Profit (%)
                  </label>
                  <Input
                    type="number"
                    placeholder="Enter take profit percentage"
                    className="bg-brutal-background border-brutal-border"
                  />
                </div>
                
                <Button className="w-full bg-brutal-info text-white hover:bg-brutal-info/80">
                  Execute Buy Order
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default Trading;
