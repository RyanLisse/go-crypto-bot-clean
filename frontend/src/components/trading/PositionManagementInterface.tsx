import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { formatNumber } from '@/utils/formatters';
import { cn } from "@/lib/utils";

// Types
interface Position {
  id: string;
  symbol: string;
  size: number;
  entryPrice: number;
  currentPrice: number;
  pnl: number;
  pnlPercentage: number;
}

interface Order {
  id: string;
  symbol: string;
  side: 'buy' | 'sell';
  type: 'market' | 'limit';
  price: number;
  size: number;
  status: string;
}

interface NewOrder {
  symbol: string;
  side: 'buy' | 'sell';
  type: 'market' | 'limit';
  price: number;
  size: number;
}

interface PositionManagementInterfaceProps {
  positions: Position[];
  orders: Order[];
  onClosePosition: (positionId: string) => void;
  onCancelOrder: (orderId: string) => void;
  onPlaceOrder: (order: NewOrder) => void;
}

const PositionManagementInterface: React.FC<PositionManagementInterfaceProps> = ({
  positions,
  orders,
  onClosePosition,
  onCancelOrder,
  onPlaceOrder
}) => {
  const [selectedTab, setSelectedTab] = useState("positions");
  const [newOrder, setNewOrder] = useState<NewOrder>({
    symbol: '',
    side: 'buy',
    type: 'limit',
    price: 0,
    size: 0
  });

  const handleOrderSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onPlaceOrder({
      ...newOrder,
      price: Number(newOrder.price),
      size: Number(newOrder.size)
    });
    
    // Reset form
    setNewOrder({
      symbol: '',
      side: 'buy',
      type: 'limit',
      price: 0,
      size: 0
    });
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setNewOrder(prev => ({
      ...prev,
      [name]: value
    }));
  };

  return (
    <Card className="bg-brutal-card-bg border-brutal-border">
      <CardHeader className="pb-2">
        <CardTitle className="text-brutal-text text-lg">Position Management</CardTitle>
      </CardHeader>
      <CardContent>
        <Tabs defaultValue="positions" onValueChange={setSelectedTab}>
          <TabsList className="grid grid-cols-3 mb-4">
            <TabsTrigger value="positions">Positions</TabsTrigger>
            <TabsTrigger value="orders">Orders</TabsTrigger>
            <TabsTrigger value="place-order">Place Order</TabsTrigger>
          </TabsList>
          
          {/* Positions Tab */}
          <TabsContent value="positions" className="space-y-4">
            {positions.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="text-brutal-text-muted text-xs border-b border-brutal-border">
                      <th className="text-left py-2">Symbol</th>
                      <th className="text-right py-2">Size</th>
                      <th className="text-right py-2">Entry Price</th>
                      <th className="text-right py-2">Current Price</th>
                      <th className="text-right py-2">P&L</th>
                      <th className="text-right py-2">P&L %</th>
                      <th className="text-right py-2">Action</th>
                    </tr>
                  </thead>
                  <tbody>
                    {positions.map((position) => (
                      <tr 
                        key={position.id} 
                        className="border-b border-brutal-border text-brutal-text"
                        data-testid={`position-${position.id}`}
                      >
                        <td className="py-3">{position.symbol}</td>
                        <td className="text-right py-3">{position.size}</td>
                        <td className="text-right py-3">{formatNumber(position.entryPrice)}</td>
                        <td className="text-right py-3">{formatNumber(position.currentPrice)}</td>
                        <td 
                          className={cn(
                            "text-right py-3", 
                            position.pnl > 0 ? "text-green-500" : "text-red-500"
                          )}
                        >
                          {position.pnl > 0 ? '+' : ''}{formatNumber(position.pnl)}
                        </td>
                        <td 
                          className={cn(
                            "text-right py-3", 
                            position.pnlPercentage > 0 ? "text-green-500" : "text-red-500"
                          )}
                        >
                          {position.pnlPercentage > 0 ? '+' : ''}{formatNumber(position.pnlPercentage)}%
                        </td>
                        <td className="text-right py-3">
                          <Button 
                            variant="destructive" 
                            size="sm"
                            onClick={() => onClosePosition(position.id)}
                          >
                            Close
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="text-center py-8 text-brutal-text-muted">
                No open positions
              </div>
            )}
          </TabsContent>
          
          {/* Orders Tab */}
          <TabsContent value="orders" className="space-y-4">
            {orders.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="text-brutal-text-muted text-xs border-b border-brutal-border">
                      <th className="text-left py-2">Symbol</th>
                      <th className="text-center py-2">Side</th>
                      <th className="text-center py-2">Type</th>
                      <th className="text-right py-2">Price</th>
                      <th className="text-right py-2">Size</th>
                      <th className="text-center py-2">Status</th>
                      <th className="text-right py-2">Action</th>
                    </tr>
                  </thead>
                  <tbody>
                    {orders.map((order) => (
                      <tr 
                        key={order.id} 
                        className="border-b border-brutal-border text-brutal-text"
                        data-testid={`order-${order.id}`}
                      >
                        <td className="py-3">{order.symbol}</td>
                        <td 
                          className={cn(
                            "text-center py-3", 
                            order.side === 'buy' ? "text-green-500" : "text-red-500"
                          )}
                        >
                          {order.side.toUpperCase()}
                        </td>
                        <td className="text-center py-3">{order.type.toUpperCase()}</td>
                        <td className="text-right py-3">{formatNumber(order.price)}</td>
                        <td className="text-right py-3">{order.size}</td>
                        <td className="text-center py-3 capitalize">{order.status}</td>
                        <td className="text-right py-3">
                          <Button 
                            variant="destructive"
                            size="sm"
                            onClick={() => onCancelOrder(order.id)}
                          >
                            Cancel
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="text-center py-8 text-brutal-text-muted">
                No open orders
              </div>
            )}
          </TabsContent>
          
          {/* Place Order Tab */}
          <TabsContent value="place-order">
            <form onSubmit={handleOrderSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="symbol">Symbol</Label>
                <Input
                  id="symbol"
                  name="symbol"
                  placeholder="BTC/USD"
                  value={newOrder.symbol}
                  onChange={handleInputChange}
                  required
                  className="bg-brutal-input-bg border-brutal-border text-brutal-text"
                />
              </div>
              
              <div className="space-y-2">
                <Label>Side</Label>
                <RadioGroup 
                  defaultValue="buy" 
                  onValueChange={(value) => setNewOrder(prev => ({ ...prev, side: value as 'buy' | 'sell' }))}
                  className="flex gap-4"
                >
                  <div className="flex items-center space-x-2">
                    <RadioGroupItem value="buy" id="buy" />
                    <Label 
                      htmlFor="buy" 
                      className="text-green-500 cursor-pointer"
                    >
                      Buy
                    </Label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <RadioGroupItem value="sell" id="sell" />
                    <Label 
                      htmlFor="sell" 
                      className="text-red-500 cursor-pointer"
                    >
                      Sell
                    </Label>
                  </div>
                </RadioGroup>
              </div>
              
              <div className="space-y-2">
                <Label>Order Type</Label>
                <RadioGroup 
                  defaultValue="limit" 
                  onValueChange={(value) => setNewOrder(prev => ({ ...prev, type: value as 'market' | 'limit' }))}
                  className="flex gap-4"
                >
                  <div className="flex items-center space-x-2">
                    <RadioGroupItem value="limit" id="limit" />
                    <Label htmlFor="limit" className="cursor-pointer">Limit</Label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <RadioGroupItem value="market" id="market" />
                    <Label htmlFor="market" className="cursor-pointer">Market</Label>
                  </div>
                </RadioGroup>
              </div>
              
              {newOrder.type === 'limit' && (
                <div className="space-y-2">
                  <Label htmlFor="price">Price</Label>
                  <Input
                    id="price"
                    name="price"
                    type="number"
                    step="0.01"
                    min="0"
                    placeholder="0.00"
                    value={newOrder.price || ''}
                    onChange={handleInputChange}
                    required
                    className="bg-brutal-input-bg border-brutal-border text-brutal-text"
                  />
                </div>
              )}
              
              <div className="space-y-2">
                <Label htmlFor="size">Size</Label>
                <Input
                  id="size"
                  name="size"
                  type="number"
                  step="0.001"
                  min="0"
                  placeholder="0.00"
                  value={newOrder.size || ''}
                  onChange={handleInputChange}
                  required
                  className="bg-brutal-input-bg border-brutal-border text-brutal-text"
                />
              </div>
              
              <Button 
                type="submit" 
                className={cn(
                  "w-full", 
                  newOrder.side === 'buy' 
                    ? "bg-green-500 hover:bg-green-600" 
                    : "bg-red-500 hover:bg-red-600"
                )}
              >
                Place Order
              </Button>
            </form>
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  );
};

export default PositionManagementInterface; 