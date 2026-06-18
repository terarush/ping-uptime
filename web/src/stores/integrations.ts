import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface Integration {
  id: number;
  name: string;
  type: string;
  config: string;
  enabled: boolean;
  user_id: number;
  created_at: string;
}

export const useIntegrationsStore = defineStore('integrations', () => {
  const integrations = ref<Integration[]>([]);
  const loading = ref(false);
  const error = ref('');

  const fetchAll = async () => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.get('/integrations');
      integrations.value = response.data?.data || response.data || [];
    } catch (err: any) {
      console.error('Failed to fetch integrations:', err);
      error.value = err.response?.data?.error || 'Failed to load integrations.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const create = async (payload: any) => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.post('/integrations', payload);
      return response.data;
    } catch (err: any) {
      console.error('Failed to create integration:', err);
      error.value = err.response?.data?.error || 'Failed to create integration.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const update = async (id: number, payload: any) => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.put(`/integrations/${id}`, payload);
      return response.data;
    } catch (err: any) {
      console.error('Failed to update integration:', err);
      error.value = err.response?.data?.error || 'Failed to update integration.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const deleteIntegration = async (id: number) => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.delete(`/integrations/${id}`);
      integrations.value = integrations.value.filter(i => i.id !== id);
      return response.data;
    } catch (err: any) {
      console.error('Failed to delete integration:', err);
      error.value = err.response?.data?.error || 'Failed to delete integration.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const test = async (id: number) => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.post(`/integrations/${id}/test`);
      return response.data;
    } catch (err: any) {
      console.error('Failed to test integration:', err);
      error.value = err.response?.data?.error || 'Failed to test integration.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    integrations,
    loading,
    error,
    fetchAll,
    create,
    update,
    deleteIntegration,
    test,
  };
});
