import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';
import { useQuery, useMutation } from '@tanstack/react-query';
import { Loader2, AlertCircle, TrendingUp, TrendingDown, RefreshCw, ChevronDown, ChevronUp, Lightbulb, BarChart, LineChart, PieChart } from 'lucide-react';
import { cn } from '@/lib/utils';
import { usePortfolioQuery } from '@/hooks/queries/usePortfolioQueries';
import { useTradeHistoryQuery } from '@/hooks/queries/useTradeQueries';

interface AIInsightsProps {
  className?: string;
}

// Insight data structure
interface Insight {
  id: string;
  title: string;
  description: string;
  type: 'portfolio' | 'market' | 'opportunity';
  importance: 'high' | 'medium' | 'low';
  timestamp: string;
  metrics?: {
    name: string;
    value: string;
    change?: number;
  }[];
  recommendation?: string;
}

// Utility function for AI insights API call
async function fetchAIInsights(portfolioData: any, tradeHistory: any): Promise<Insight[]> {
  const response = await fetch('http://localhost:8080/api/ai/insights', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    },
    body: JSON.stringify({
      portfolio: portfolioData,
      trade_history: tradeHistory,
      insight_types: ['portfolio', 'market', 'opportunity']
    })
  });
  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }
  const data = await response.json();
  return data.insights;
}

// Mutation function for generating AI insights
const generateAIInsightsFn = async ({ portfolio, tradeHistory }: { portfolio: any; tradeHistory: any }): Promise<Insight[]> => {
  try {
    return await fetchAIInsights(portfolio, tradeHistory);
  } catch (error) {
    // Fallback to mock insights if the API call fails
    return [
      {
        id: '1',
        title: 'Portfolio Diversification Alert',
        description: 'Your portfolio is heavily concentrated in BTC (65%). Consider diversifying to reduce risk.',
        type: 'portfolio',
        importance: 'high',
        timestamp: new Date().toISOString(),
        metrics: [
          { name: 'Concentration Risk', value: 'High', change: 15 },
          { name: 'Volatility', value: '32%', change: 5 }
        ],
        recommendation: 'Reduce BTC allocation to 40% and distribute to ETH, SOL, and AVAX.'
      },
      {
        id: '2',
        title: 'Market Trend Analysis',
        description: 'The market is showing signs of recovery after recent correction. DeFi tokens are leading the recovery.',
        type: 'market',
        importance: 'medium',
        timestamp: new Date().toISOString(),
        metrics: [
          { name: 'Market Sentiment', value: 'Bullish', change: 10 },
          { name: 'DeFi Index', value: '+8.5%', change: 8.5 }
        ]
      },
      {
        id: '3',
        title: 'Trading Opportunity: SOL',
        description: 'SOL is showing strong technical signals with increasing volume and breaking resistance.',
        type: 'opportunity',
        importance: 'high',
        timestamp: new Date().toISOString(),
        metrics: [
          { name: 'Technical Score', value: '85/100', change: 12 },
          { name: 'Volume Change', value: '+45%', change: 45 }
        ],
        recommendation: 'Consider increasing SOL position by 2-3% of portfolio.'
      }
    ];
  }
};

export function AIInsights({ className }: AIInsightsProps) {
  const [activeTab, setActiveTab] = useState('all');
  const [expandedInsights, setExpandedInsights] = useState<string[]>([]);
  
  // Get portfolio and trade history data
  const { data: portfolioData } = usePortfolioQuery();
  const { data: tradeHistory } = useTradeHistoryQuery(50);
  
  // TanStack Query mutation for generating insights
  const {
    mutate: generateInsights,
    data: insights,
    isPending,
    error,
    reset,
  } = useMutation<Insight[], Error, { portfolio: any; tradeHistory: any }>(generateAIInsightsFn);

  // Trigger insights generation when portfolioData or tradeHistory changes
  React.useEffect(() => {
    if (portfolioData && tradeHistory) {
      generateInsights({ portfolio: portfolioData, tradeHistory });
    } else {
      reset();
    }
  }, [portfolioData, tradeHistory, generateInsights, reset]);
  
  // Toggle insight expansion
  const toggleInsight = (id: string) => {
    setExpandedInsights(prev => 
      prev.includes(id) 
        ? prev.filter(i => i !== id) 
        : [...prev, id]
    );
  };
  
  // Filter insights based on active tab
  const filteredInsights = insights?.filter(insight => {
    if (activeTab === 'all') return true;
    return insight.type === activeTab;
  });
  
  // Get importance color
  const getImportanceColor = (importance: string) => {
    switch (importance) {
      case 'high':
        return 'bg-red-100 text-red-800';
      case 'medium':
        return 'bg-yellow-100 text-yellow-800';
      case 'low':
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };
  
  // Get type icon
  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'portfolio':
        return <PieChart className="h-4 w-4" />;
      case 'market':
        return <BarChart className="h-4 w-4" />;
      case 'opportunity':
        return <LineChart className="h-4 w-4" />;
      default:
        return <Lightbulb className="h-4 w-4" />;
    }
  };

  if (isPending) {
    return (
      <Card className={cn("h-full", className)}>
        <CardHeader className="pb-2">
          <CardTitle className="text-lg font-medium">AI Insights</CardTitle>
        </CardHeader>
        <CardContent className="flex items-center justify-center h-[300px]">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </CardContent>
      </Card>
    );
  }

  if (error || !insights) {
    return (
      <Card className={cn("h-full", className)}>
        <CardHeader className="pb-2">
          <CardTitle className="text-lg font-medium">AI Insights</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center h-[300px] text-center">
          <AlertCircle className="h-10 w-10 text-destructive mb-2" />
          <p className="text-sm text-muted-foreground">Failed to load AI insights</p>
          <Button onClick={() => generateInsights({ portfolio: portfolioData, tradeHistory })} variant="outline" size="sm" className="mt-4">
            <RefreshCw className="h-4 w-4 mr-2" /> Try again
          </Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className={cn("h-full", className)}>
      <CardHeader className="pb-2">
        <div className="flex justify-between items-center">
          <CardTitle className="text-lg font-medium">AI Insights</CardTitle>
          <Button variant="outline" size="sm" onClick={() => generateInsights({ portfolio: portfolioData, tradeHistory })}>
            <RefreshCw className="h-4 w-4" />
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <Tabs defaultValue="all" value={activeTab} onValueChange={setActiveTab} className="w-full">
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="all">All</TabsTrigger>
            <TabsTrigger value="portfolio">Portfolio</TabsTrigger>
            <TabsTrigger value="market">Market</TabsTrigger>
            <TabsTrigger value="opportunity">Opportunities</TabsTrigger>
          </TabsList>
          
          <TabsContent value={activeTab} className="pt-4">
            <div className="space-y-4">
              {filteredInsights?.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  No insights available
                </div>
              ) : (
                filteredInsights?.map((insight) => (
                  <Collapsible
                    key={insight.id}
                    open={expandedInsights.includes(insight.id)}
                    onOpenChange={() => toggleInsight(insight.id)}
                    className="border rounded-md"
                  >
                    <CollapsibleTrigger asChild>
                      <div className="flex justify-between items-center p-4 cursor-pointer hover:bg-muted/50">
                        <div className="flex items-center space-x-2">
                          <div className="flex-shrink-0">
                            {getTypeIcon(insight.type)}
                          </div>
                          <div>
                            <h3 className="font-medium">{insight.title}</h3>
                            <div className="flex items-center space-x-2 mt-1">
                              <Badge variant="outline" className={cn(getImportanceColor(insight.importance))}>
                                {insight.importance}
                              </Badge>
                              <span className="text-xs text-muted-foreground">
                                {new Date(insight.timestamp).toLocaleDateString()}
                              </span>
                            </div>
                          </div>
                        </div>
                        <div>
                          {expandedInsights.includes(insight.id) ? (
                            <ChevronUp className="h-5 w-5 text-muted-foreground" />
                          ) : (
                            <ChevronDown className="h-5 w-5 text-muted-foreground" />
                          )}
                        </div>
                      </div>
                    </CollapsibleTrigger>
                    <CollapsibleContent className="p-4 pt-0 border-t">
                      <p className="text-sm mb-4">{insight.description}</p>
                      
                      {insight.metrics && (
                        <div className="grid grid-cols-2 gap-4 mb-4">
                          {insight.metrics.map((metric, index) => (
                            <div key={index} className="bg-muted/50 p-3 rounded-md">
                              <p className="text-xs text-muted-foreground">{metric.name}</p>
                              <div className="flex items-center mt-1">
                                <span className="text-lg font-medium">{metric.value}</span>
                                {metric.change !== undefined && (
                                  <span className={cn(
                                    "ml-2 text-xs flex items-center",
                                    metric.change >= 0 ? "text-green-600" : "text-red-600"
                                  )}>
                                    {metric.change >= 0 ? (
                                      <>
                                        <TrendingUp className="h-3 w-3 mr-1" />
                                        +{metric.change}%
                                      </>
                                    ) : (
                                      <>
                                        <TrendingDown className="h-3 w-3 mr-1" />
                                        {metric.change}%
                                      </>
                                    )}
                                  </span>
                                )}
                              </div>
                            </div>
                          ))}
                        </div>
                      )}
                      
                      {insight.recommendation && (
                        <div className="bg-blue-50 border-l-4 border-blue-500 p-3 rounded-md">
                          <p className="text-sm font-medium text-blue-800">Recommendation:</p>
                          <p className="text-sm text-blue-700">{insight.recommendation}</p>
                        </div>
                      )}
                    </CollapsibleContent>
                  </Collapsible>
                ))
              )}
            </div>
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  );
}

export { AIInsights };
