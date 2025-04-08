
import React from 'react';
import { Link } from 'react-router-dom';

export default function NotFound() {
  return (
    <div className="h-screen flex flex-col items-center justify-center bg-brutal-background">
      <div className="brutal-card p-10 max-w-md w-full">
        <h1 className="text-4xl font-bold mb-4 text-brutal-error">404</h1>
        <p className="text-xl text-brutal-text mb-8">PAGE NOT FOUND</p>
        <div className="border-t border-brutal-border pt-6">
          <Link to="/" className="text-brutal-info hover:underline">
            RETURN TO DASHBOARD
          </Link>
        </div>
      </div>
    </div>
  );
}
