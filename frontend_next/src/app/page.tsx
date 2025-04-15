import { redirect } from 'next/navigation';

export default function Home() {
  // Just redirect to the sign-in page, Clerk middleware will handle the authentication check
  redirect('/sign-in');
}
