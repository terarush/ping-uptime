/**
 * Auth Store (Pinia)
 * Manages JWT tokens, validates active session status,
 * and maintains authenticated user profile records.
 */
import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';
import Cookies from 'js-cookie';

export interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export const useAuthStore = defineStore('auth', () => {
  // --- Reactive State Variables ---
  const currentUser = ref<User | null>(null);
  const isAuthenticated = ref(false);
  const isVerifying = ref(false);

  /**
   * Checks for a valid cookie token and calls the server to verify it
   * @returns {Promise<boolean>} True if the token is valid, false otherwise
   */
  const verifyToken = async (): Promise<boolean> => {
    const token = Cookies.get('accessToken');
    if (!token) {
      currentUser.value = null;
      isAuthenticated.value = false;
      return false;
    }

    isVerifying.value = true;
    try {
      // Validate token with verify backend endpoint
      const response = await ExtendedFetch.get('/users/verify');
      if (response.data?.data) {
        currentUser.value = response.data.data;
        isAuthenticated.value = true;
        return true;
      }
      throw new Error('No user data returned');
    } catch (err) {
      console.warn('Token verification failed on server, clearing session:', err);
      currentUser.value = null;
      isAuthenticated.value = false;
      
      // Wipe invalid cookies
      Cookies.remove('accessToken');
      Cookies.remove('refreshToken');
      return false;
    } finally {
      isVerifying.value = false;
    }
  };

  /**
   * Logs out the user by clearing local reactive states and wiping cookies
   */
  const logout = () => {
    currentUser.value = null;
    isAuthenticated.value = false;
    Cookies.remove('accessToken');
    Cookies.remove('refreshToken');
  };

  /**
   * Saves access and refresh tokens as cookies and initializes the user session
   */
  const setSession = (accessToken: string, refreshToken: string, user: User) => {
    // Save tokens as secure cookies with a 7-day expiration
    Cookies.set('accessToken', accessToken, { expires: 7 });
    Cookies.set('refreshToken', refreshToken, { expires: 7 });
    currentUser.value = user;
    isAuthenticated.value = true;
  };

  return {
    currentUser,
    isAuthenticated,
    isVerifying,
    verifyToken,
    logout,
    setSession
  };
});
