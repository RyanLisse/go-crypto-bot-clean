import React from 'react';

const TestPage = () => {
  return (
    <div className="p-6 bg-brutal-background text-brutal-text">
      <h1 className="text-2xl font-bold mb-4">Test Page</h1>
      <p className="mb-4">This is a test page to check if the styling is working correctly.</p>
      
      <div className="brutal-card mb-4">
        <div className="brutal-card-header">TEST CARD</div>
        <div className="text-2xl font-bold">$1,234.56</div>
        <div className="text-sm mt-2 text-brutal-success">+5.2%</div>
      </div>
      
      <div className="grid grid-cols-2 gap-4">
        <div className="bg-brutal-panel border border-brutal-border p-4">
          Panel 1
        </div>
        <div className="bg-brutal-panel border border-brutal-border p-4">
          Panel 2
        </div>
      </div>
    </div>
  );
};

export default TestPage;
