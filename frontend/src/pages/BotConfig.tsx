import React, { useState } from 'react';
import { Header } from '@/components/layout/Header';
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertTriangle, Plus, Save, Trash, Undo2, X } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog";
import { Form, FormField, FormItem, FormLabel, FormControl } from "@/components/ui/form";
import { Slider } from "@/components/ui/slider";

type Strategy = {
  id: string;
  name: string;
  description: string;
}

type TradingPair = {
  id: string;
  symbol: string;
  name: string;
  active: boolean;
}

type TakeProfitLevel = {
  id: string;
  percentage: number;
  positionSize: number;
}

const strategies: Strategy[] = [
  {
    id: 'trend',
    name: 'Trend Following',
    description: 'Follow market trends using technical indicators like moving averages'
  },
  {
    id: 'ml',
    name: 'Machine Learning',
    description: 'Use AI prediction models to determine optimal entry and exit points'
  },
  {
    id: 'arbitrage',
    name: 'Arbitrage',
    description: 'Exploit price differences of the same asset across different markets'
  },
  {
    id: 'sniper',
    name: 'Sniper',
    description: 'Target specific price points for quick entry and exit with precision'
  },
  {
    id: 'sandwich',
    name: 'Sandwich',
    description: 'Execute trades before and after large transactions to capitalize on price impacts'
  }
];

const tradingPairs: TradingPair[] = [
  { id: '1', symbol: 'BTC/USDT', name: 'Bitcoin/USDT', active: true },
  { id: '2', symbol: 'ETH/USDT', name: 'Ethereum/USDT', active: true },
  { id: '3', symbol: 'SOL/USDT', name: 'Solana/USDT', active: true },
  { id: '4', symbol: 'AVAX/USDT', name: 'Avalanche/USDT', active: false },
  { id: '5', symbol: 'BNB/USDT', name: 'Binance Coin/USDT', active: true },
  { id: '6', symbol: 'ADA/USDT', name: 'Cardano/USDT', active: false },
  { id: '7', symbol: 'DOT/USDT', name: 'Polkadot/USDT', active: false },
  { id: '8', symbol: 'LINK/USDT', name: 'Chainlink/USDT', active: true }
];

export default function BotConfig() {
  const { toast } = useToast();
  const [selectedStrategy, setSelectedStrategy] = useState('ml');
  const [selectedPairs, setSelectedPairs] = useState<string[]>(
    tradingPairs.filter(pair => pair.active).map(pair => pair.id)
  );
  const [riskLevel, setRiskLevel] = useState(3);
  const [maxConcurrentTrades, setMaxConcurrentTrades] = useState(8);
  const [takeProfitLevels, setTakeProfitLevels] = useState<TakeProfitLevel[]>([
    { id: '1', percentage: 5, positionSize: 50 },
    { id: '2', percentage: 10, positionSize: 30 },
    { id: '3', percentage: 20, positionSize: 20 }
  ]);
  const [stopLossPercent, setStopLossPercent] = useState(3);
  const [dailyTradeLimit, setDailyTradeLimit] = useState(20);
  const [apiKey, setApiKey] = useState('••••••••••••••••••••••');
  const [apiSecret, setApiSecret] = useState('••••••••••••••••••••••');
  const [isNewLevelDialogOpen, setIsNewLevelDialogOpen] = useState(false);
  const [newLevelPercentage, setNewLevelPercentage] = useState(15);
  const [newLevelPositionSize, setNewLevelPositionSize] = useState(25);
  const [useNewCoins, setUseNewCoins] = useState(false);
  const [useMachineLearning, setUseMachineLearning] = useState(true);

  const handlePairToggle = (pairId: string) => {
    setSelectedPairs(prev => 
      prev.includes(pairId)
        ? prev.filter(id => id !== pairId)
        : [...prev, pairId]
    );
  };

  const handleSaveConfig = () => {
    // In a real app, this would send the config to the API
    console.log({
      strategy: selectedStrategy,
      tradingPairs: selectedPairs,
      riskLevel,
      maxConcurrentTrades,
      takeProfitLevels,
      stopLossPercent,
      dailyTradeLimit,
      useNewCoins,
      useMachineLearning
    });

    toast({
      title: "Configuration Saved",
      description: "Your bot configuration has been updated successfully.",
    });
  };

  const handleResetToDefault = () => {
    setSelectedStrategy('ml');
    setSelectedPairs(tradingPairs.filter(pair => pair.active).map(pair => pair.id));
    setRiskLevel(3);
    setMaxConcurrentTrades(8);
    setTakeProfitLevels([
      { id: '1', percentage: 5, positionSize: 50 },
      { id: '2', percentage: 10, positionSize: 30 },
      { id: '3', percentage: 20, positionSize: 20 }
    ]);
    setStopLossPercent(3);
    setDailyTradeLimit(20);
    setUseNewCoins(false);
    setUseMachineLearning(true);

    toast({
      title: "Reset to Default",
      description: "The configuration has been reset to system defaults.",
    });
  };

  const addTakeProfitLevel = () => {
    if (takeProfitLevels.length >= 5) {
      toast({
        title: "Limit Reached",
        description: "Maximum of 5 take-profit levels allowed",
        variant: "destructive"
      });
      return;
    }

    // Validate total position size doesn't exceed 100%
    const currentTotalSize = takeProfitLevels.reduce((sum, level) => sum + level.positionSize, 0);
    if (currentTotalSize + newLevelPositionSize > 100) {
      toast({
        title: "Invalid Position Size",
        description: "Total position size cannot exceed 100%",
        variant: "destructive"
      });
      return;
    }

    const newLevel = {
      id: Date.now().toString(),
      percentage: newLevelPercentage,
      positionSize: newLevelPositionSize
    };

    setTakeProfitLevels([...takeProfitLevels, newLevel]);
    setIsNewLevelDialogOpen(false);
    setNewLevelPercentage(15);
    setNewLevelPositionSize(25);
  };

  const removeTakeProfitLevel = (id: string) => {
    setTakeProfitLevels(takeProfitLevels.filter(level => level.id !== id));
  };

  const getTotalPositionSize = () => {
    return takeProfitLevels.reduce((sum, level) => sum + level.positionSize, 0);
  };

  return (
    <div className="flex-1 flex flex-col h-screen overflow-auto">
      <Header />
      
      <div className="flex-1 p-6 space-y-6">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h2 className="text-2xl font-bold tracking-tight text-brutal-text font-mono">BOT CONFIGURATION</h2>
            <p className="text-brutal-text/70">Configure trading strategies and risk parameters</p>
          </div>
          <div className="flex items-center gap-3">
            <Button 
              variant="outline"
              className="flex items-center gap-2 border-brutal-border text-brutal-text"
              onClick={handleResetToDefault}
            >
              <Undo2 className="h-4 w-4" />
              RESET
            </Button>
            <Button 
              className="flex items-center gap-2 bg-brutal-info text-black hover:bg-brutal-info/80"
              onClick={handleSaveConfig}
            >
              <Save className="h-4 w-4" />
              SAVE CONFIG
            </Button>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Strategy Selection */}
          <Card className="bg-brutal-panel border border-brutal-border">
            <CardHeader>
              <CardTitle className="text-brutal-text font-mono text-base">TRADING STRATEGY</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {strategies.map((strategy) => (
                  <div 
                    key={strategy.id}
                    className={`p-3 border cursor-pointer ${
                      selectedStrategy === strategy.id 
                        ? 'border-brutal-info bg-brutal-info/10' 
                        : 'border-brutal-border hover:bg-brutal-panel/80'
                    }`}
                    onClick={() => setSelectedStrategy(strategy.id)}
                  >
                    <div className="font-mono text-sm font-bold text-brutal-text">{strategy.name}</div>
                    <div className="text-xs text-brutal-text/70 mt-1">{strategy.description}</div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Risk Management */}
          <Card className="bg-brutal-panel border border-brutal-border">
            <CardHeader>
              <CardTitle className="text-brutal-text font-mono text-base">RISK MANAGEMENT</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="space-y-2">
                  <label className="text-xs text-brutal-text/70">Risk Level {riskLevel}/5</label>
                  <input 
                    type="range" 
                    min="1" 
                    max="5" 
                    value={riskLevel} 
                    onChange={(e) => setRiskLevel(parseInt(e.target.value))}
                    className="w-full h-1 bg-brutal-border rounded-lg appearance-none cursor-pointer"
                  />
                  <div className="flex justify-between text-[10px] text-brutal-text/50">
                    <span>Low Risk</span>
                    <span>High Risk</span>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-1">
                    <label className="text-xs text-brutal-text/70 block">Stop Loss %</label>
                    <input
                      type="number"
                      value={stopLossPercent}
                      onChange={(e) => setStopLossPercent(parseInt(e.target.value) || 0)}
                      className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                    />
                  </div>
                  <div className="space-y-1">
                    <label className="text-xs text-brutal-text/70 block">Max Concurrent Trades</label>
                    <input
                      type="number"
                      value={maxConcurrentTrades}
                      onChange={(e) => setMaxConcurrentTrades(parseInt(e.target.value) || 0)}
                      className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                    />
                  </div>
                  <div className="space-y-1">
                    <label className="text-xs text-brutal-text/70 block">Daily Trade Limit</label>
                    <input
                      type="number"
                      value={dailyTradeLimit}
                      onChange={(e) => setDailyTradeLimit(parseInt(e.target.value) || 0)}
                      className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                    />
                  </div>
                </div>

                <div className="mt-4 p-3 bg-brutal-warning/10 border border-brutal-warning/30 flex items-start gap-2">
                  <AlertTriangle className="h-4 w-4 text-brutal-warning mt-0.5" />
                  <div className="text-xs text-brutal-text/80">
                    Higher risk levels may result in increased trade frequency and larger position sizes.
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Multi-level Take Profit */}
          <Card className="bg-brutal-panel border border-brutal-border">
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle className="text-brutal-text font-mono text-base">TAKE PROFIT LEVELS</CardTitle>
              <Button 
                variant="outline" 
                size="sm" 
                onClick={() => setIsNewLevelDialogOpen(true)}
                className="h-8 border-brutal-border text-brutal-text"
              >
                <Plus className="h-3.5 w-3.5 mr-1" />
                ADD LEVEL
              </Button>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="space-y-2">
                  <div className="grid grid-cols-12 gap-2 text-xs text-brutal-text/70 font-mono mb-1">
                    <div className="col-span-1">#</div>
                    <div className="col-span-4">PROFIT %</div>
                    <div className="col-span-5">POSITION %</div>
                    <div className="col-span-2"></div>
                  </div>
                  
                  {takeProfitLevels.map((level, index) => (
                    <div key={level.id} className="grid grid-cols-12 gap-2 items-center">
                      <div className="col-span-1 text-xs font-mono text-brutal-text/70">{index + 1}</div>
                      <div className="col-span-4">
                        <input
                          type="number"
                          value={level.percentage}
                          onChange={(e) => {
                            const newLevels = [...takeProfitLevels];
                            newLevels[index].percentage = Number(e.target.value);
                            setTakeProfitLevels(newLevels);
                          }}
                          className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono text-xs"
                        />
                      </div>
                      <div className="col-span-5">
                        <input
                          type="number"
                          value={level.positionSize}
                          onChange={(e) => {
                            const newLevels = [...takeProfitLevels];
                            newLevels[index].positionSize = Number(e.target.value);
                            setTakeProfitLevels(newLevels);
                          }}
                          className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono text-xs"
                        />
                      </div>
                      <div className="col-span-2 text-right">
                        {takeProfitLevels.length > 1 && (
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => removeTakeProfitLevel(level.id)}
                            className="h-8 text-brutal-text/70 hover:text-brutal-warning"
                          >
                            <Trash className="h-3.5 w-3.5" />
                          </Button>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
                
                <div className={`flex items-center justify-between p-2 ${
                  getTotalPositionSize() === 100 
                    ? 'bg-brutal-info/10 border border-brutal-info/30' 
                    : 'bg-brutal-warning/10 border border-brutal-warning/30'
                }`}>
                  <span className="text-xs text-brutal-text/80">Total Position Size:</span>
                  <span className={`text-xs font-mono font-bold ${
                    getTotalPositionSize() === 100 
                      ? 'text-brutal-info' 
                      : 'text-brutal-warning'
                  }`}>
                    {getTotalPositionSize()}%
                  </span>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Advanced Strategy Options */}
          <Card className="bg-brutal-panel border border-brutal-border">
            <CardHeader>
              <CardTitle className="text-brutal-text font-mono text-base">ADVANCED OPTIONS</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-5">
                <div className="flex justify-between items-center">
                  <div>
                    <div className="text-sm font-mono text-brutal-text">Use New Coins</div>
                    <div className="text-xs text-brutal-text/70 mt-1">Automatically trade newly listed tokens</div>
                  </div>
                  <div 
                    onClick={() => setUseNewCoins(!useNewCoins)}
                    className={`w-12 h-6 flex items-center p-1 rounded-full cursor-pointer transition-colors ${
                      useNewCoins ? 'bg-brutal-info justify-end' : 'bg-brutal-border justify-start'
                    }`}
                  >
                    <div className="bg-white w-4 h-4 rounded-full"></div>
                  </div>
                </div>
                
                <div className="flex justify-between items-center">
                  <div>
                    <div className="text-sm font-mono text-brutal-text">Machine Learning</div>
                    <div className="text-xs text-brutal-text/70 mt-1">Use AI prediction models for trading decisions</div>
                  </div>
                  <div 
                    onClick={() => setUseMachineLearning(!useMachineLearning)}
                    className={`w-12 h-6 flex items-center p-1 rounded-full cursor-pointer transition-colors ${
                      useMachineLearning ? 'bg-brutal-info justify-end' : 'bg-brutal-border justify-start'
                    }`}
                  >
                    <div className="bg-white w-4 h-4 rounded-full"></div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Trading Pairs */}
          <Card className="bg-brutal-panel border border-brutal-border">
            <CardHeader>
              <CardTitle className="text-brutal-text font-mono text-base">TRADING PAIRS</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-2">
                {tradingPairs.map((pair) => (
                  <div 
                    key={pair.id} 
                    className={`p-2 border flex items-center gap-2 cursor-pointer ${
                      selectedPairs.includes(pair.id) 
                        ? 'border-brutal-info bg-brutal-info/10' 
                        : 'border-brutal-border hover:bg-brutal-panel/80'
                    }`}
                    onClick={() => handlePairToggle(pair.id)}
                  >
                    <div className={`w-4 h-4 ${
                      selectedPairs.includes(pair.id) 
                        ? 'bg-brutal-info' 
                        : 'border border-brutal-border'
                    }`}></div>
                    <span className="font-mono text-xs text-brutal-text">{pair.symbol}</span>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* API Configuration */}
          <Card className="bg-brutal-panel border border-brutal-border">
            <CardHeader>
              <CardTitle className="text-brutal-text font-mono text-base">API CONFIGURATION</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="space-y-1">
                  <label className="text-xs text-brutal-text/70 block">Exchange API Key</label>
                  <input
                    type="password"
                    value={apiKey}
                    onChange={(e) => setApiKey(e.target.value)}
                    className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs text-brutal-text/70 block">Exchange API Secret</label>
                  <input
                    type="password"
                    value={apiSecret}
                    onChange={(e) => setApiSecret(e.target.value)}
                    className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                  />
                </div>
                
                <div className="p-3 bg-brutal-info/10 border border-brutal-info/30 text-xs text-brutal-text/80">
                  API keys are securely stored and never exposed to frontend code.
                </div>
              </div>
            </CardContent>
          </Card>
          
          {/* Schedule Configuration */}
          <div className="col-span-1 lg:col-span-2">
            <Card className="bg-brutal-panel border border-brutal-border">
              <CardHeader>
                <CardTitle className="text-brutal-text font-mono text-base">TRADING SCHEDULE</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="grid grid-cols-7 gap-2">
                    {["MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN"].map((day, index) => (
                      <div key={day} className="flex flex-col items-center">
                        <div className="text-xs text-brutal-text/70 mb-2">{day}</div>
                        <div className={`w-full p-3 border text-center cursor-pointer ${
                          index < 5 ? 'border-brutal-info bg-brutal-info/10' : 'border-brutal-border'
                        }`}>
                          <span className="font-mono text-xs text-brutal-text">
                            {index < 5 ? 'ACTIVE' : 'OFF'}
                          </span>
                        </div>
                      </div>
                    ))}
                  </div>
                  
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-1">
                      <label className="text-xs text-brutal-text/70 block">Trading Hours (Start)</label>
                      <select 
                        className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                        defaultValue="00:00"
                      >
                        {[...Array(24)].map((_, i) => (
                          <option key={i} value={`${i.toString().padStart(2, '0')}:00`}>
                            {`${i.toString().padStart(2, '0')}:00`}
                          </option>
                        ))}
                      </select>
                    </div>
                    <div className="space-y-1">
                      <label className="text-xs text-brutal-text/70 block">Trading Hours (End)</label>
                      <select 
                        className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                        defaultValue="23:59"
                      >
                        {[...Array(24)].map((_, i) => (
                          <option key={i} value={`${i.toString().padStart(2, '0')}:00`}>
                            {`${i.toString().padStart(2, '0')}:00`}
                          </option>
                        ))}
                      </select>
                    </div>
                  </div>
                  
                  <div className="flex items-center gap-2">
                    <div className="w-4 h-4 bg-brutal-info"></div>
                    <span className="text-xs text-brutal-text/80">24/7 Trading Enabled</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>

      {/* Dialog for adding new take profit level */}
      <Dialog open={isNewLevelDialogOpen} onOpenChange={setIsNewLevelDialogOpen}>
        <DialogContent className="bg-brutal-panel border-brutal-border">
          <DialogHeader>
            <DialogTitle className="text-brutal-text font-mono">Add Take Profit Level</DialogTitle>
            <DialogDescription className="text-brutal-text/70">
              Configure the percentage and position size for this level.
            </DialogDescription>
          </DialogHeader>
          
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <label className="text-sm text-brutal-text font-mono">Take Profit Percentage</label>
              <Input
                type="number"
                value={newLevelPercentage}
                onChange={(e) => setNewLevelPercentage(Number(e.target.value))}
                className="bg-brutal-background border-brutal-border text-brutal-text"
              />
              <div className="text-xs text-brutal-text/70">Target percentage increase for profit taking</div>
            </div>
            
            <div className="space-y-2">
              <label className="text-sm text-brutal-text font-mono">Position Size Percentage</label>
              <Input
                type="number"
                value={newLevelPositionSize}
                onChange={(e) => setNewLevelPositionSize(Number(e.target.value))}
                className="bg-brutal-background border-brutal-border text-brutal-text"
              />
              <div className="text-xs text-brutal-text/70">
                Percentage of your position to sell at this level
              </div>
            </div>
            
            <div className={`p-3 text-xs ${
              getTotalPositionSize() + newLevelPositionSize > 100 
                ? 'bg-brutal-warning/10 border border-brutal-warning/30 text-brutal-warning' 
                : 'bg-brutal-info/10 border border-brutal-info/30 text-brutal-text/80'
            }`}>
              {getTotalPositionSize() + newLevelPositionSize > 100 
                ? `Total position size would exceed 100% (${getTotalPositionSize() + newLevelPositionSize}%)`
                : `New total position size: ${getTotalPositionSize() + newLevelPositionSize}%`
              }
            </div>
          </div>
          
          <DialogFooter>
            <Button 
              variant="outline" 
              onClick={() => setIsNewLevelDialogOpen(false)}
              className="border-brutal-border text-brutal-text"
            >
              Cancel
            </Button>
            <Button 
              onClick={addTakeProfitLevel}
              className="bg-brutal-info text-black hover:bg-brutal-info/80"
              disabled={getTotalPositionSize() + newLevelPositionSize > 100}
            >
              Add Level
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
