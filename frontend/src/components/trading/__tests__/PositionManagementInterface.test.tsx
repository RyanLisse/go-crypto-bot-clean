import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import PositionManagementInterface from '../PositionManagementInterface';

describe('PositionManagementInterface', () => {
  const mockPositions = [
    { id: '1', symbol: 'BTC/USD', size: 1.5, entryPrice: 50000, currentPrice: 51000, pnl: 1500, pnlPercentage: 3 },
    { id: '2', symbol: 'ETH/USD', size: 10, entryPrice: 3000, currentPrice: 2900, pnl: -1000, pnlPercentage: -3.33 }
  ];
  
  const mockOrders = [
    { id: '101', symbol: 'BTC/USD', side: 'buy', type: 'limit', price: 49500, size: 0.5, status: 'open' },
    { id: '102', symbol: 'ETH/USD', side: 'sell', type: 'limit', price: 3100, size: 5, status: 'open' }
  ];

  const mockClosePosition = vi.fn();
  const mockCancelOrder = vi.fn();
  const mockPlaceOrder = vi.fn();
  
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders position management interface with tabs', () => {
    render(
      <PositionManagementInterface 
        positions={mockPositions}
        orders={mockOrders}
        onClosePosition={mockClosePosition}
        onCancelOrder={mockCancelOrder}
        onPlaceOrder={mockPlaceOrder}
      />
    );
    
    // Check tabs are present
    expect(screen.getByText('Positions')).toBeInTheDocument();
    expect(screen.getByText('Orders')).toBeInTheDocument();
    expect(screen.getByText('Place Order')).toBeInTheDocument();
  });

  it('displays position information correctly', () => {
    render(
      <PositionManagementInterface 
        positions={mockPositions}
        orders={mockOrders}
        onClosePosition={mockClosePosition}
        onCancelOrder={mockCancelOrder}
        onPlaceOrder={mockPlaceOrder}
      />
    );
    
    // Check position data is displayed
    expect(screen.getByText('BTC/USD')).toBeInTheDocument();
    expect(screen.getByText('1.5')).toBeInTheDocument();
    expect(screen.getByText('$50,000.00')).toBeInTheDocument();
    expect(screen.getByText('$51,000.00')).toBeInTheDocument();
    expect(screen.getByText('$1,500.00')).toBeInTheDocument();
    expect(screen.getByText('+3.00%')).toBeInTheDocument();
    
    expect(screen.getByText('ETH/USD')).toBeInTheDocument();
    expect(screen.getByText('10')).toBeInTheDocument();
    expect(screen.getByText('$3,000.00')).toBeInTheDocument();
    expect(screen.getByText('$2,900.00')).toBeInTheDocument();
    expect(screen.getByText('-$1,000.00')).toBeInTheDocument();
    expect(screen.getByText('-3.33%')).toBeInTheDocument();
  });

  it('calls onClosePosition when close button is clicked', () => {
    render(
      <PositionManagementInterface 
        positions={mockPositions}
        orders={mockOrders}
        onClosePosition={mockClosePosition}
        onCancelOrder={mockCancelOrder}
        onPlaceOrder={mockPlaceOrder}
      />
    );
    
    // Click close button for first position
    const closeButtons = screen.getAllByText('Close');
    fireEvent.click(closeButtons[0]);
    
    expect(mockClosePosition).toHaveBeenCalledWith('1');
  });

  it('displays orders information correctly', () => {
    render(
      <PositionManagementInterface 
        positions={mockPositions}
        orders={mockOrders}
        onClosePosition={mockClosePosition}
        onCancelOrder={mockCancelOrder}
        onPlaceOrder={mockPlaceOrder}
      />
    );
    
    // Switch to Orders tab
    fireEvent.click(screen.getByText('Orders'));
    
    // Check order data is displayed
    expect(screen.getByText('BTC/USD')).toBeInTheDocument();
    expect(screen.getByText('BUY')).toBeInTheDocument();
    expect(screen.getByText('LIMIT')).toBeInTheDocument();
    expect(screen.getByText('$49,500.00')).toBeInTheDocument();
    expect(screen.getByText('0.5')).toBeInTheDocument();
    
    expect(screen.getByText('ETH/USD')).toBeInTheDocument();
    expect(screen.getByText('SELL')).toBeInTheDocument();
    expect(screen.getByText('LIMIT')).toBeInTheDocument();
    expect(screen.getByText('$3,100.00')).toBeInTheDocument();
    expect(screen.getByText('5')).toBeInTheDocument();
  });

  it('calls onCancelOrder when cancel button is clicked', () => {
    render(
      <PositionManagementInterface 
        positions={mockPositions}
        orders={mockOrders}
        onClosePosition={mockClosePosition}
        onCancelOrder={mockCancelOrder}
        onPlaceOrder={mockPlaceOrder}
      />
    );
    
    // Switch to Orders tab
    fireEvent.click(screen.getByText('Orders'));
    
    // Click cancel button for first order
    const cancelButtons = screen.getAllByText('Cancel');
    fireEvent.click(cancelButtons[0]);
    
    expect(mockCancelOrder).toHaveBeenCalledWith('101');
  });

  it('allows placing a new order', () => {
    render(
      <PositionManagementInterface 
        positions={mockPositions}
        orders={mockOrders}
        onClosePosition={mockClosePosition}
        onCancelOrder={mockCancelOrder}
        onPlaceOrder={mockPlaceOrder}
      />
    );
    
    // Switch to Place Order tab
    fireEvent.click(screen.getByText('Place Order'));
    
    // Fill order form
    fireEvent.change(screen.getByLabelText('Symbol'), { target: { value: 'BTC/USD' } });
    fireEvent.click(screen.getByLabelText('Buy'));
    fireEvent.change(screen.getByLabelText('Price'), { target: { value: '50000' } });
    fireEvent.change(screen.getByLabelText('Size'), { target: { value: '1' } });
    
    // Submit order
    fireEvent.click(screen.getByText('Place Order'));
    
    expect(mockPlaceOrder).toHaveBeenCalledWith({
      symbol: 'BTC/USD',
      side: 'buy',
      type: 'limit',
      price: 50000,
      size: 1
    });
  });

  it('displays empty state when no positions are available', () => {
    render(
      <PositionManagementInterface 
        positions={[]}
        orders={mockOrders}
        onClosePosition={mockClosePosition}
        onCancelOrder={mockCancelOrder}
        onPlaceOrder={mockPlaceOrder}
      />
    );
    
    expect(screen.getByText('No open positions')).toBeInTheDocument();
  });

  it('displays empty state when no orders are available', () => {
    render(
      <PositionManagementInterface 
        positions={mockPositions}
        orders={[]}
        onClosePosition={mockClosePosition}
        onCancelOrder={mockCancelOrder}
        onPlaceOrder={mockPlaceOrder}
      />
    );
    
    // Switch to Orders tab
    fireEvent.click(screen.getByText('Orders'));
    
    expect(screen.getByText('No open orders')).toBeInTheDocument();
  });
}); 