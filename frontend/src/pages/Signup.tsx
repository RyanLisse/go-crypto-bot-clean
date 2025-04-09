import React from 'react';
import { SignUp } from '@clerk/clerk-react';
import { Box, Container, Typography } from '@mui/material';

const Signup: React.FC = () => {
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
          Create an Account
        </Typography>
        <SignUp path="/signup" routing="path" signInUrl="/login" />
      </Box>
    </Container>
  );
};

export default Signup;
