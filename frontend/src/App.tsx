import React from 'react';
import { Routes, Route, BrowserRouter } from 'react-router-dom';
import { Box } from '@mui/material';
import { SignedIn, SignedOut, RedirectToSignIn, SignInButton, UserButton } from "@clerk/clerk-react";

import Layout from '@/components/Layout';
import Dashboard from '@/pages/Dashboard';
import Portfolio from '@/pages/Portfolio';
import Trading from '@/pages/Trading';
import Backtesting from '@/pages/Backtesting';
import Settings from '@/pages/Settings';
import NotFound from '@/pages/NotFound';
import Login from '@/pages/Login';
import Signup from '@/pages/Signup';

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Box sx={{ display: 'flex' }}>
        <Routes>
          <Route path="/" element={<Layout />}>
            {/* Public routes */}
            <Route index element={
              <>
                <SignedIn>
                  <Dashboard />
                </SignedIn>
                <SignedOut>
                  <RedirectToSignIn />
                </SignedOut>
              </>
            } />

            {/* Protected routes */}
            <Route path="portfolio" element={
              <>
                <SignedIn>
                  <Portfolio />
                </SignedIn>
                <SignedOut>
                  <RedirectToSignIn />
                </SignedOut>
              </>
            } />
            <Route path="trading" element={
              <>
                <SignedIn>
                  <Trading />
                </SignedIn>
                <SignedOut>
                  <RedirectToSignIn />
                </SignedOut>
              </>
            } />
            <Route path="backtesting" element={
              <>
                <SignedIn>
                  <Backtesting />
                </SignedIn>
                <SignedOut>
                  <RedirectToSignIn />
                </SignedOut>
              </>
            } />
            <Route path="settings" element={
              <>
                <SignedIn>
                  <Settings />
                </SignedIn>
                <SignedOut>
                  <RedirectToSignIn />
                </SignedOut>
              </>
            } />
            {/* Auth routes */}
          <Route path="login" element={<Login />} />
          <Route path="signup" element={<Signup />} />
          <Route path="*" element={<NotFound />} />
          </Route>
        </Routes>
      </Box>
    </BrowserRouter>
  );
};

export default App;
