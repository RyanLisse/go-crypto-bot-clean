
export interface PortfolioValue {
  date: string;
  value: number;
}

export interface Holding {
  coin: string;
  symbol: string;
  amount: string;
  price: number;
  value: number;
  allocation: number;
  change24h: number;
  change7d: number;
  cost: number;
  pnl: number;
}

export interface Transaction {
  id: string;
  type: 'BUY' | 'SELL';
  coin: string;
  amount: string;
  price: number;
  total: number;
  date: string;
  status: 'completed' | 'pending' | 'failed';
}

export interface PortfolioData {
  totalValue: number;
  totalPnL: number;
  totalPnLPercent: number;
  historicalData: PortfolioValue[];
  holdings: Holding[];
  recentTransactions: Transaction[];
  lastUpdated: string;
}
