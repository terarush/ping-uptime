import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface ApiToken {
  id: number;
  name: string;
  token_prefix: string;
  last_used_at: string | null;
  expires_at: string | null;
  is_revoked: boolean;
  created_at: string;
}

export interface CreateTokenResult extends ApiToken {
  raw_token: string;
}

export const useApiTokensStore = defineStore('apiTokens', () => {
  const tokens = ref<ApiToken[]>([]);
  const loading = ref(false);
  const error = ref('');

  const fetchAll = async () => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.get('/api-tokens');
      tokens.value = response.data?.data || response.data || [];
    } catch (err: any) {
      console.error('Failed to fetch API tokens:', err);
      error.value = err.response?.data?.error || 'Failed to load API tokens.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const create = async (name: string, expiresAt?: string) => {
    loading.value = true;
    error.value = '';
    try {
      const payload: any = { name };
      if (expiresAt) {
        payload.expires_at = expiresAt;
      }
      const response = await ExtendedFetch.post('/api-tokens', payload);
      return response.data?.data as CreateTokenResult;
    } catch (err: any) {
      console.error('Failed to create API token:', err);
      error.value = err.response?.data?.error || 'Failed to create API token.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const revoke = async (id: number) => {
    loading.value = true;
    error.value = '';
    try {
      await ExtendedFetch.delete(`/api-tokens/${id}`);
      tokens.value = tokens.value.filter(t => t.id !== id);
    } catch (err: any) {
      console.error('Failed to revoke API token:', err);
      error.value = err.response?.data?.error || 'Failed to revoke API token.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    tokens,
    loading,
    error,
    fetchAll,
    create,
    revoke,
  };
});
