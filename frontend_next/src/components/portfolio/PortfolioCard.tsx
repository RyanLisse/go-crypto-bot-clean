import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export interface PortfolioAsset {
  symbol: string;
  amount: number;
  value: number;
  change: number;
}

export interface PortfolioData {
  totalValue: number;
  change: number;
  assets: PortfolioAsset[];
}

interface PortfolioCardProps {
  data: PortfolioData;
}

export function PortfolioCard({ data }: PortfolioCardProps) {
  return (
    <Card className="portfolio-card border-2 border-black">
      <CardHeader className="border-b-2 border-black px-4 py-2">
        <CardTitle className="text-lg font-mono">Portfolio Summary</CardTitle>
      </CardHeader>
      <CardContent className="p-4">
        <div className="portfolio-value flex justify-between items-center mb-4 pb-2 border-b border-gray-200">
          <span className="label font-mono font-bold">Total Value:</span>
          <div className="flex items-center">
            <span className="value font-mono text-xl mr-2">${data.totalValue.toFixed(2)}</span>
            <span className={`change font-mono text-sm px-2 py-1 ${data.change >= 0 ? 'bg-green-100' : 'bg-red-100'}`}>
              {data.change >= 0 ? '+' : ''}{data.change.toFixed(2)}%
            </span>
          </div>
        </div>
        <div className="assets">
          <h4 className="text-md font-mono font-bold mb-2">Assets</h4>
          <ul className="space-y-2">
            {data.assets.map((asset) => (
              <li key={asset.symbol} className="flex justify-between items-center py-1 border-b border-gray-100">
                <span className="symbol font-mono font-bold">{asset.symbol}</span>
                <span className="amount font-mono text-sm">{asset.amount.toFixed(6)}</span>
                <span className="value font-mono">${asset.value.toFixed(2)}</span>
                <span className={`change font-mono text-sm px-2 py-1 ${asset.change >= 0 ? 'bg-green-100' : 'bg-red-100'}`}>
                  {asset.change >= 0 ? '+' : ''}{asset.change.toFixed(2)}%
                </span>
              </li>
            ))}
          </ul>
        </div>
      </CardContent>
    </Card>
  );
}

export { PortfolioCard };
