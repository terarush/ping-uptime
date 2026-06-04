/**
 * Authentication Composable
 * Wraps useAuthStore and storeToRefs to provide a clean reactive interface
 * for managing user session state throughout the application.
 */
import { useAuthStore, type User } from '@/stores/auth';
import { storeToRefs } from 'pinia';

export type { User };

export function useAuth() {
  // Initialize the Pinia Auth store instance
  const authStore = useAuthStore();
  
  // Extract reactive properties with storeToRefs to maintain reactivity
  const { currentUser, isAuthenticated, isVerifying } = storeToRefs(authStore);

  return {
    currentUser,
    isAuthenticated,
    isVerifying,
    verifyToken: authStore.verifyToken,
    logout: authStore.logout,
    setSession: authStore.setSession
  };
}
