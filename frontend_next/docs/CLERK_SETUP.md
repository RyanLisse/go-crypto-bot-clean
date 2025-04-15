# Clerk Authentication Setup

This project uses Clerk for authentication. This document will guide you through the setup process.

## Setting Up Clerk Authentication

### 1. Register with Clerk

1. Go to [Clerk's website](https://clerk.dev/) and sign up for an account
2. Create a new application in your Clerk dashboard
3. Set up your authentication methods (Email, social logins, etc.)
4. Get your API keys from the Clerk dashboard

### 2. Environment Variables

Add the following environment variables to your `.env.local` file:

```
# Clerk Auth
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=your_clerk_publishable_key
CLERK_SECRET_KEY=your_clerk_secret_key

# Clerk auth URLs
NEXT_PUBLIC_CLERK_SIGN_IN_URL=/sign-in
NEXT_PUBLIC_CLERK_SIGN_UP_URL=/sign-up
NEXT_PUBLIC_CLERK_AFTER_SIGN_IN_URL=/dashboard
NEXT_PUBLIC_CLERK_AFTER_SIGN_UP_URL=/dashboard

# Fallback URLs
NEXT_PUBLIC_CLERK_SIGN_IN_FALLBACK_REDIRECT_URL=/dashboard
NEXT_PUBLIC_CLERK_SIGN_UP_FALLBACK_REDIRECT_URL=/dashboard
```

Replace the placeholder API keys with your actual Clerk API keys.

### 3. Project Structure

The authentication flow is set up with the following files:

- `src/app/sign-in/[[...sign-in]]/page.tsx`: The sign-in page
- `src/app/sign-up/[[...sign-up]]/page.tsx`: The sign-up page
- `src/middleware.ts`: Protects routes and handles authentication
- `src/app/layout.tsx`: Wraps the app with the ClerkProvider

### 4. Accessing the User

You can access the user in client components with the `useUser` hook:

```tsx
'use client';
import { useUser } from '@clerk/nextjs';

export default function MyComponent() {
  const { user } = useUser();
  
  return <div>Hello, {user?.fullName}</div>;
}
```

And in server components with the `currentUser` function:

```tsx
import { currentUser } from '@clerk/nextjs';

export default async function MyServerComponent() {
  const user = await currentUser();
  
  return <div>Hello, {user?.fullName}</div>;
}
```

### 5. Protecting Routes

Routes are protected using the middleware. By default, all routes require authentication except for:

- `/sign-in/*`
- `/sign-up/*`
- `/`
- `/api/webhooks/clerk`

To make a route public, add it to the `isPublicRoute` matcher in `src/middleware.ts`.

### 6. User Sign Out

The UserButton component from Clerk handles sign-out functionality automatically. If you need to programmatically sign out:

```tsx
'use client';
import { useClerk } from '@clerk/nextjs';

export default function SignOutButton() {
  const { signOut } = useClerk();
  
  return <button onClick={() => signOut()}>Sign out</button>;
}
```

## Further Reading

- [Clerk Documentation](https://clerk.com/docs)
- [Clerk + Next.js Integration](https://clerk.com/docs/quickstarts/nextjs)
- [Customizing Appearance](https://clerk.com/docs/customization/appearance) 