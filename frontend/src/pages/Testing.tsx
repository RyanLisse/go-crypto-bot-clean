import React from 'react';
import { PageLayout } from '@/components/layout/PageLayout';
import ConnectionTest from '@/components/testing/ConnectionTest';

export function Testing() {
  return (
    <PageLayout>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="md:col-span-2">
          <h1 className="text-2xl font-bold mb-4">System Tests</h1>
          <p className="mb-6">
            This page allows you to test the connection to the backend and AI services.
            Run the tests to verify that everything is working correctly.
          </p>
        </div>
        
        <div className="md:col-span-1">
          <ConnectionTest />
        </div>
        
        <div className="md:col-span-1">
          <div className="brutal-card">
            <div className="brutal-card-header mb-4">Test Instructions</div>
            <div className="p-4">
              <h3 className="text-lg font-bold mb-2">How to Test</h3>
              <ol className="list-decimal ml-5 space-y-2">
                <li>
                  <strong>Run All Tests</strong> - Tests both API and AI connections in sequence
                </li>
                <li>
                  <strong>Test API</strong> - Tests only the API connection to the backend
                </li>
                <li>
                  <strong>Test AI</strong> - Tests only the AI connection (will use fallback if backend is unavailable)
                </li>
              </ol>
              
              <h3 className="text-lg font-bold mt-6 mb-2">Troubleshooting</h3>
              <ul className="list-disc ml-5 space-y-2">
                <li>
                  <strong>Backend Disconnected</strong> - Make sure the backend is running using the <code>./run-dev.sh</code> script
                </li>
                <li>
                  <strong>API Test Failed</strong> - Check the backend logs for errors
                </li>
                <li>
                  <strong>AI in Fallback Mode</strong> - This is normal when the backend is unavailable, but the AI will have limited capabilities
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </PageLayout>
  );
}

export default Testing;
