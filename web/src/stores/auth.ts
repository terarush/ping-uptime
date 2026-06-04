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
  const currentUser = ref<User | null>(null);
  const isAuthenticated = ref(false);
  const isVerifying = ref(false);

  const verifyToken = async (): Promise<boolean> => {
    const token = Cookies.get('accessToken');
    if (!token) {
      currentUser.value = null;
      isAuthenticated.value = false;
      return false;
    }

    isVerifying.value = true;
    try {
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
      Cookies.remove('accessToken');
      Cookies.remove('refreshToken');
      return false;
    } finally {
      isVerifying.value = false;
    }
  };

  const logout = () => {
    currentUser.value = null;
    isAuthenticated.value = false;
    Cookies.remove('accessToken');
    Cookies.remove('refreshToken');
  };

  const setSession = (accessToken: string, refreshToken: string, user: User) => {
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
