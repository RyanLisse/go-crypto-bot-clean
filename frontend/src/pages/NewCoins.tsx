
import React from 'react';
import { Header } from '@/components/layout/Header';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Star, Calendar, Coins, TrendingUp, AlertCircle } from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';

const NewCoins = () => {
  // Mock data for upcoming coin releases
  const upcomingCoins = [
    {
      name: 'DefiChain',
      symbol: 'DFI',
      releaseDate: '2025-04-15',
      description: 'Decentralized finance on Bitcoin',
      category: 'DeFi',
      hasMint: true,
      hasLiquidity: true,
      tags: ['defi', 'bitcoin'],
      riskLevel: 'Medium'
    },
    {
      name: 'MetaVerse Token',
      symbol: 'MVT',
      releaseDate: '2025-04-18',
      description: 'Virtual reality metaverse platform',
      category: 'Metaverse',
      hasMint: true,
      hasLiquidity: false,
      tags: ['gaming', 'metaverse'],
      riskLevel: 'High'
    },
    {
      name: 'GreenEnergy',
      symbol: 'GREN',
      releaseDate: '2025-04-22',
      description: 'Sustainable blockchain mining solution',
      category: 'Utility',
      hasMint: false,
      hasLiquidity: false,
      tags: ['eco', 'mining'],
      riskLevel: 'Medium'
    },
    {
      name: 'ArtifactNFT',
      symbol: 'ANFT',
      releaseDate: '2025-04-25',
      description: 'Museum artifact NFT marketplace',
      category: 'NFT',
      hasMint: true,
      hasLiquidity: true,
      tags: ['nft', 'art'],
      riskLevel: 'Low'
    },
    {
      name: 'DataChain',
      symbol: 'DATA',
      releaseDate: '2025-04-29',
      description: 'Decentralized data storage solution',
      category: 'Infrastructure',
      hasMint: true,
      hasLiquidity: true,
      tags: ['storage', 'web3'],
      riskLevel: 'Low'
    }
  ];

  // Mock data for recent launches
  const recentLaunches = [
    {
      name: 'SportBet',
      symbol: 'SBET',
      launchDate: '2025-04-07',
      initialPrice: 0.12,
      currentPrice: 0.28,
      change: 133.3,
      hasMint: true,
      hasLiquidity: true,
      volume24h: 4500000
    },
    {
      name: 'AudioStream',
      symbol: 'ASTR',
      launchDate: '2025-04-05',
      initialPrice: 2.50,
      currentPrice: 3.75,
      change: 50.0,
      hasMint: true,
      hasLiquidity: true,
      volume24h: 9800000
    },
    {
      name: 'QuantAI',
      symbol: 'QANT',
      launchDate: '2025-04-03',
      initialPrice: 5.20,
      currentPrice: 4.85,
      change: -6.7,
      hasMint: true,
      hasLiquidity: true,
      volume24h: 3200000
    }
  ];

  // Format date to be more readable
  const formatDate = (dateString: string) => {
    const options: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'short', day: 'numeric' };
    return new Date(dateString).toLocaleDateString(undefined, options);
  };

  // Calculate days until release
  const daysUntil = (dateString: string) => {
    const today = new Date();
    const releaseDate = new Date(dateString);
    const diffTime = Math.abs(releaseDate.getTime() - today.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays;
  };

  return (
    <div className="flex-1 flex flex-col h-full overflow-auto">
      <Header />
      
      <div className="flex-1 p-4 md:p-6 space-y-4 md:space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-brutal-text tracking-tight">NEW COINS</h1>
          <p className="text-brutal-text/70 text-sm">Track upcoming coin releases and recent launches</p>
        </div>
        
        <Card className="bg-brutal-panel border-brutal-border">
          <CardHeader className="pb-2">
            <CardTitle className="text-brutal-text flex items-center text-lg">
              <Calendar className="mr-2 h-5 w-5 text-brutal-info" />
              Upcoming Coin Releases
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow className="border-brutal-border">
                    <TableHead className="text-brutal-text/70">Coin</TableHead>
                    <TableHead className="text-brutal-text/70">Category</TableHead>
                    <TableHead className="text-brutal-text/70">Release Date</TableHead>
                    <TableHead className="text-brutal-text/70">Status</TableHead>
                    <TableHead className="text-brutal-text/70">Risk Level</TableHead>
                    <TableHead className="text-brutal-text/70">Description</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {upcomingCoins.map((coin) => (
                    <TableRow key={coin.symbol} className="border-brutal-border">
                      <TableCell className="font-medium text-brutal-text">
                        <div className="flex items-center">
                          <div className="w-8 h-8 rounded-full bg-brutal-info/20 mr-2 flex items-center justify-center text-xs">
                            {coin.symbol.substring(0, 2)}
                          </div>
                          <div>
                            <div className="font-bold">{coin.name}</div>
                            <div className="text-xs text-brutal-text/70">{coin.symbol}</div>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline" className="bg-brutal-background border-brutal-border">
                          {coin.category}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <div className="font-mono">
                          {formatDate(coin.releaseDate)}
                          <div className="text-xs text-brutal-info">
                            {daysUntil(coin.releaseDate)} days left
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex flex-col space-y-1">
                          <div className="flex items-center text-xs">
                            <div className={`w-2 h-2 rounded-full mr-1 ${coin.hasMint ? 'bg-brutal-success' : 'bg-brutal-error'}`}></div>
                            {coin.hasMint ? 'Minted' : 'Not Minted'}
                          </div>
                          <div className="flex items-center text-xs">
                            <div className={`w-2 h-2 rounded-full mr-1 ${coin.hasLiquidity ? 'bg-brutal-success' : 'bg-brutal-error'}`}></div>
                            {coin.hasLiquidity ? 'Has Liquidity' : 'No Liquidity'}
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className={`
                          px-2 py-1 rounded text-xs text-center
                          ${coin.riskLevel === 'Low' ? 'bg-brutal-success/20 text-brutal-success' : 
                            coin.riskLevel === 'Medium' ? 'bg-brutal-warning/20 text-brutal-warning' : 
                            'bg-brutal-error/20 text-brutal-error'}
                        `}>
                          {coin.riskLevel}
                        </div>
                      </TableCell>
                      <TableCell className="text-brutal-text/70 max-w-[200px] truncate">
                        {coin.description}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
        
        <Card className="bg-brutal-panel border-brutal-border">
          <CardHeader className="pb-2">
            <CardTitle className="text-brutal-text flex items-center text-lg">
              <Coins className="mr-2 h-5 w-5 text-brutal-success" />
              Recent Launches
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {recentLaunches.map((coin) => (
                <div 
                  key={coin.symbol}
                  className="border border-brutal-border rounded-md p-4 bg-brutal-background"
                >
                  <div className="flex justify-between items-start mb-3">
                    <div className="flex items-center">
                      <div className="w-10 h-10 rounded-full bg-brutal-info/20 mr-3 flex items-center justify-center text-sm">
                        {coin.symbol.substring(0, 2)}
                      </div>
                      <div>
                        <div className="font-bold text-brutal-text">{coin.name}</div>
                        <div className="text-xs text-brutal-text/70">{coin.symbol}</div>
                      </div>
                    </div>
                    <Badge 
                      className={`${coin.change >= 0 ? 'bg-brutal-success' : 'bg-brutal-error'} text-white`}
                    >
                      {coin.change >= 0 ? '+' : ''}{coin.change}%
                    </Badge>
                  </div>
                  
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-brutal-text/70">Launch Date:</span>
                      <span className="font-mono text-brutal-text">{formatDate(coin.launchDate)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-brutal-text/70">Initial Price:</span>
                      <span className="font-mono text-brutal-text">${coin.initialPrice.toFixed(4)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-brutal-text/70">Current Price:</span>
                      <span className="font-mono text-brutal-text">${coin.currentPrice.toFixed(4)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-brutal-text/70">24h Volume:</span>
                      <span className="font-mono text-brutal-text">${(coin.volume24h/1000000).toFixed(1)}M</span>
                    </div>
                    
                    <div className="flex items-center text-xs pt-2">
                      <div className="flex items-center mr-3">
                        <div className={`w-2 h-2 rounded-full mr-1 ${coin.hasMint ? 'bg-brutal-success' : 'bg-brutal-error'}`}></div>
                        {coin.hasMint ? 'Minted' : 'Not Minted'}
                      </div>
                      <div className="flex items-center">
                        <div className={`w-2 h-2 rounded-full mr-1 ${coin.hasLiquidity ? 'bg-brutal-success' : 'bg-brutal-error'}`}></div>
                        {coin.hasLiquidity ? 'Has Liquidity' : 'No Liquidity'}
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
        
        <Card className="bg-brutal-panel border-brutal-border">
          <CardHeader className="pb-2">
            <CardTitle className="text-brutal-text flex items-center text-lg">
              <AlertCircle className="mr-2 h-5 w-5 text-brutal-warning" />
              New Coin Risk Assessment
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              <div className="p-4 border border-brutal-border rounded-md">
                <div className="flex items-center mb-2">
                  <TrendingUp className="h-5 w-5 text-brutal-info mr-2" />
                  <div className="font-bold">Market Analysis</div>
                </div>
                <p className="text-sm text-brutal-text/70">
                  Always check market conditions before investing in new coins. 
                  High volatility periods increase risk.
                </p>
              </div>
              
              <div className="p-4 border border-brutal-border rounded-md">
                <div className="flex items-center mb-2">
                  <Coins className="h-5 w-5 text-brutal-success mr-2" />
                  <div className="font-bold">Liquidity Check</div>
                </div>
                <p className="text-sm text-brutal-text/70">
                  Low liquidity coins have higher slippage and are more susceptible 
                  to price manipulation.
                </p>
              </div>
              
              <div className="p-4 border border-brutal-border rounded-md">
                <div className="flex items-center mb-2">
                  <Star className="h-5 w-5 text-brutal-warning mr-2" />
                  <div className="font-bold">Team Verification</div>
                </div>
                <p className="text-sm text-brutal-text/70">
                  Research the team behind new coins. Anon teams carry 
                  higher risk of rugpulls.
                </p>
              </div>
              
              <div className="p-4 border border-brutal-border rounded-md">
                <div className="flex items-center mb-2">
                  <AlertCircle className="h-5 w-5 text-brutal-error mr-2" />
                  <div className="font-bold">Smart Contract</div>
                </div>
                <p className="text-sm text-brutal-text/70">
                  Review smart contract for security vulnerabilities or 
                  suspicious functions.
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default NewCoins;
