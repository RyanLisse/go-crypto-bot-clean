import React from 'react';
import { SignIn } from '@clerk/clerk-react';
import { Box, Container, Typography } from '@mui/material';

const Login: React.FC = () => {
  return (
    <Container maxWidth="sm">
      <Box
        sx={{
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Typography component="h1" variant="h5" sx={{ mb: 4 }}>
          Sign in to Crypto Trading Bot
        </Typography>
        <SignIn path="/login" routing="path" signUpUrl="/signup" />
      </Box>
    </Container>
  );
};

export default Login;
