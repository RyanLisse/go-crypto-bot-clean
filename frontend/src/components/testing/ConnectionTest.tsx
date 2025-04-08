import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useBackendStatus } from '@/hooks/useBackendStatus';
import { api } from '@/lib/api';
import { sendChatMessage, getAIMetrics } from '@/lib/aiClient';
import { toast } from 'sonner';
import { CheckCircle, XCircle, AlertTriangle, RefreshCw } from 'lucide-react';

/**
 * Component to test the connection to the backend and AI services
 */
export function ConnectionTest() {
  const [isTestingAPI, setIsTestingAPI] = useState(false);
  const [isTestingAI, setIsTestingAI] = useState(false);
  const [apiTestResult, setApiTestResult] = useState<'success' | 'error' | null>(null);
  const [aiTestResult, setAiTestResult] = useState<'success' | 'error' | 'fallback' | null>(null);
  const [apiResponse, setApiResponse] = useState<any>(null);
  const [aiResponse, setAiResponse] = useState<string | null>(null);
  const [testTime, setTestTime] = useState<Date | null>(null);
  
  const { isConnected, status, refetch } = useBackendStatus({
    refetchInterval: 0, // Disable auto-refetch for the test component
  });
  
  // Reset test results when connection status changes
  useEffect(() => {
    if (!isConnected) {
      setApiTestResult(null);
      setAiTestResult(null);
    }
  }, [isConnected]);
  
  // Test the API connection
  const testAPIConnection = async () => {
    setIsTestingAPI(true);
    setApiTestResult(null);
    setApiResponse(null);
    
    try {
      // Test the status endpoint
      const statusResponse = await api.getStatus();
      
      // Test the wallet endpoint
      const walletResponse = await api.getWallet();
      
      // Set the test results
      setApiTestResult('success');
      setApiResponse({
        status: statusResponse,
        wallet: walletResponse,
      });
      
      toast.success('API connection test successful');
    } catch (error) {
      console.error('API connection test failed:', error);
      setApiTestResult('error');
      setApiResponse(error);
      
      toast.error('API connection test failed');
    } finally {
      setIsTestingAPI(false);
      setTestTime(new Date());
    }
  };
  
  // Test the AI connection
  const testAIConnection = async () => {
    setIsTestingAI(true);
    setAiTestResult(null);
    setAiResponse(null);
    
    try {
      // Send a test message to the AI
      const response = await sendChatMessage('Test message: What is the current time?');
      
      // Check if we're using fallback mode
      const metrics = getAIMetrics();
      
      if (metrics.usingFallback) {
        setAiTestResult('fallback');
        toast.warning('AI connection test completed in fallback mode');
      } else {
        setAiTestResult('success');
        toast.success('AI connection test successful');
      }
      
      setAiResponse(response.message.content);
    } catch (error) {
      console.error('AI connection test failed:', error);
      setAiTestResult('error');
      setAiResponse(String(error));
      
      toast.error('AI connection test failed');
    } finally {
      setIsTestingAI(false);
      setTestTime(new Date());
    }
  };
  
  // Run all tests
  const runAllTests = async () => {
    await refetch();
    await testAPIConnection();
    await testAIConnection();
  };
  
  return (
    <Card className="border-2 border-black">
      <CardHeader className="border-b-2 border-black px-4 py-2">
        <CardTitle className="text-lg font-mono">Connection Tests</CardTitle>
      </CardHeader>
      <CardContent className="p-4">
        <div className="space-y-4">
          {/* Connection Status */}
          <div className="p-3 border-2 border-black rounded">
            <h3 className="text-sm font-bold mb-2">Backend Connection Status</h3>
            <div className="flex items-center">
              {isConnected ? (
                <CheckCircle className="h-5 w-5 text-brutal-success mr-2" />
              ) : (
                <XCircle className="h-5 w-5 text-brutal-error mr-2" />
              )}
              <span className="font-mono">
                {isConnected ? 'Connected' : 'Disconnected'}
              </span>
            </div>
            {status && (
              <div className="mt-2 text-xs font-mono">
                <div>Version: {status.version}</div>
                <div>Uptime: {status.uptime}</div>
                {status.processes && (
                  <div className="mt-1">
                    <div className="font-bold">Services:</div>
                    <div className="ml-2">
                      {Object.entries(status.processes).map(([name, proc]: [string, any]) => (
                        <div key={name} className="flex justify-between">
                          <span>{name}:</span>
                          <span className={proc.status === 'running' ? 'text-brutal-success' : 'text-brutal-error'}>
                            {proc.status}
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
          
          {/* Test Results */}
          <div className="p-3 border-2 border-black rounded">
            <h3 className="text-sm font-bold mb-2">Test Results</h3>
            <div className="space-y-2">
              <div className="flex justify-between items-center">
                <span className="font-mono">API Connection:</span>
                <div className="flex items-center">
                  {apiTestResult === 'success' && <CheckCircle className="h-4 w-4 text-brutal-success mr-1" />}
                  {apiTestResult === 'error' && <XCircle className="h-4 w-4 text-brutal-error mr-1" />}
                  {isTestingAPI && <RefreshCw className="h-4 w-4 animate-spin mr-1" />}
                  <span className={`font-mono ${
                    apiTestResult === 'success' ? 'text-brutal-success' : 
                    apiTestResult === 'error' ? 'text-brutal-error' : ''
                  }`}>
                    {apiTestResult === 'success' ? 'Success' : 
                     apiTestResult === 'error' ? 'Failed' : 
                     isTestingAPI ? 'Testing...' : 'Not Tested'}
                  </span>
                </div>
              </div>
              
              <div className="flex justify-between items-center">
                <span className="font-mono">AI Connection:</span>
                <div className="flex items-center">
                  {aiTestResult === 'success' && <CheckCircle className="h-4 w-4 text-brutal-success mr-1" />}
                  {aiTestResult === 'error' && <XCircle className="h-4 w-4 text-brutal-error mr-1" />}
                  {aiTestResult === 'fallback' && <AlertTriangle className="h-4 w-4 text-brutal-warning mr-1" />}
                  {isTestingAI && <RefreshCw className="h-4 w-4 animate-spin mr-1" />}
                  <span className={`font-mono ${
                    aiTestResult === 'success' ? 'text-brutal-success' : 
                    aiTestResult === 'error' ? 'text-brutal-error' :
                    aiTestResult === 'fallback' ? 'text-brutal-warning' : ''
                  }`}>
                    {aiTestResult === 'success' ? 'Success' : 
                     aiTestResult === 'error' ? 'Failed' : 
                     aiTestResult === 'fallback' ? 'Fallback Mode' :
                     isTestingAI ? 'Testing...' : 'Not Tested'}
                  </span>
                </div>
              </div>
              
              {testTime && (
                <div className="text-xs text-brutal-text/50 font-mono">
                  Last tested: {testTime.toLocaleTimeString()}
                </div>
              )}
            </div>
          </div>
          
          {/* Test Actions */}
          <div className="flex space-x-2">
            <Button 
              onClick={runAllTests} 
              disabled={isTestingAPI || isTestingAI}
              className="flex-1 bg-black text-white border-2 border-black hover:bg-gray-800"
            >
              {(isTestingAPI || isTestingAI) ? (
                <>
                  <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                  Testing...
                </>
              ) : (
                'Run All Tests'
              )}
            </Button>
            <Button 
              onClick={testAPIConnection} 
              disabled={isTestingAPI}
              className="flex-1 bg-black text-white border-2 border-black hover:bg-gray-800"
            >
              {isTestingAPI ? (
                <>
                  <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                  Testing API...
                </>
              ) : (
                'Test API'
              )}
            </Button>
            <Button 
              onClick={testAIConnection} 
              disabled={isTestingAI}
              className="flex-1 bg-black text-white border-2 border-black hover:bg-gray-800"
            >
              {isTestingAI ? (
                <>
                  <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                  Testing AI...
                </>
              ) : (
                'Test AI'
              )}
            </Button>
          </div>
          
          {/* Response Details */}
          {(apiResponse || aiResponse) && (
            <div className="p-3 border-2 border-black rounded">
              <h3 className="text-sm font-bold mb-2">Response Details</h3>
              <div className="space-y-2">
                {apiResponse && (
                  <div>
                    <h4 className="text-xs font-bold">API Response:</h4>
                    <pre className="text-xs font-mono bg-brutal-panel/30 p-2 rounded overflow-auto max-h-40">
                      {JSON.stringify(apiResponse, null, 2)}
                    </pre>
                  </div>
                )}
                
                {aiResponse && (
                  <div>
                    <h4 className="text-xs font-bold">AI Response:</h4>
                    <div className="text-xs font-mono bg-brutal-panel/30 p-2 rounded overflow-auto max-h-40">
                      {aiResponse}
                    </div>
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

export default ConnectionTest;
