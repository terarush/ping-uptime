import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  is_blocked: boolean;
  avatar?: string;
  created_at: string;
}

export const useUsersStore = defineStore('users', () => {
  const users = ref<User[]>([]);
  const loading = ref(false);
  const error = ref('');

  const fetchUsersData = async () => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.get('/users');
      users.value = response.data?.data || response.data || [];
    } catch (err: any) {
      console.error('Failed to fetch users:', err);
      error.value = err.response?.data?.error || 'Failed to load user accounts.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const createUser = async (payload: any) => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.post('/users', payload);
      return response.data;
    } catch (err: any) {
      console.error('Failed to create user:', err);
      error.value = err.response?.data?.error || 'Failed to register new user.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const updateUser = async (id: number, payload: any) => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.put(`/users/${id}`, payload);
      return response.data;
    } catch (err: any) {
      console.error('Failed to update user:', err);
      error.value = err.response?.data?.error || 'Failed to update user account.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const deleteUser = async (id: number) => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.delete(`/users/${id}`);
      users.value = users.value.filter(u => u.id !== id);
      return response.data;
    } catch (err: any) {
      console.error('Failed to delete user:', err);
      error.value = err.response?.data?.error || 'Failed to delete user account.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    users,
    loading,
    error,
    fetchUsersData,
    createUser,
    updateUser,
    deleteUser,
  };
});
