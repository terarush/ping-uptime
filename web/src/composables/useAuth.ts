import { useAuthStore, type User } from '@/stores/auth';
import { storeToRefs } from 'pinia';

export type { User };

export function useAuth() {
  const authStore = useAuthStore();
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
