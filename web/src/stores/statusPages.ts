import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface Monitor {
  id: number;
  name: string;
  url: string;
}

export interface StatusPage {
  id: number;
  name: string;
  slug: string;
  description: string;
  user_id: number;
  monitors: Monitor[];
  created_at: string;
}

export const useStatusPagesStore = defineStore('statusPages', () => {
  const statusPages = ref<StatusPage[]>([]);
  const availableMonitors = ref<Monitor[]>([]);
  const loading = ref(false);
  const error = ref('');

  const fetchStatusPagesData = async () => {
    loading.value = true;
    error.value = '';
    try {
      const pagesRes = await ExtendedFetch.get('/status-pages');
      statusPages.value = pagesRes.data?.data || pagesRes.data || [];

      const monitorsRes = await ExtendedFetch.get('/monitors');
      availableMonitors.value = monitorsRes.data?.data || monitorsRes.data || [];
    } catch (err: any) {
      console.error('Failed to fetch status pages data:', err);
      error.value = err.response?.data?.error || 'Failed to load status pages.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const createStatusPage = async (payload: { name: string; slug: string; description: string; monitor_ids: number[] }) => {
    error.value = '';
    try {
      const response = await ExtendedFetch.post('/status-pages', payload);
      return response.data?.data || response.data;
    } catch (err: any) {
      console.error('Failed to create status page:', err);
      error.value = err.response?.data?.error || 'Failed to create status page.';
      throw err;
    }
  };

  const updateStatusPage = async (id: number, payload: { name: string; slug: string; description: string; monitor_ids: number[] }) => {
    error.value = '';
    try {
      const response = await ExtendedFetch.put(`/status-pages/${id}`, payload);
      return response.data?.data || response.data;
    } catch (err: any) {
      console.error('Failed to update status page:', err);
      error.value = err.response?.data?.error || 'Failed to update status page.';
      throw err;
    }
  };

  const deleteStatusPage = async (id: number) => {
    error.value = '';
    try {
      await ExtendedFetch.delete(`/status-pages/${id}`);
      statusPages.value = statusPages.value.filter(p => p.id !== id);
    } catch (err: any) {
      console.error('Failed to delete status page:', err);
      error.value = err.response?.data?.error || 'Failed to delete status page.';
      throw err;
    }
  };

  return {
    statusPages,
    availableMonitors,
    loading,
    error,
    fetchStatusPagesData,
    createStatusPage,
    updateStatusPage,
    deleteStatusPage,
  };
});
