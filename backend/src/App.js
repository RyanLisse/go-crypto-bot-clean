import React from 'react';
import './App.css';
import MexcApiTest from './components/MexcApiTest';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>Crypto Bot Frontend</h1>
      </header>
      <main>
        <MexcApiTest />
      </main>
      <footer className="App-footer">
        <p>&copy; 2023 Crypto Bot</p>
      </footer>
    </div>
  );
}

export default App;
