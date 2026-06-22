import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface Subscriber {
  id: number;
  email: string;
  status_page_id: number;
  verified: boolean;
  created_at: string;
}

export const useSubscribersStore = defineStore('subscribers', () => {
  const subscribers = ref<Subscriber[]>([]);
  const statusPages = ref<{ id: number; name: string }[]>([]);
  const loading = ref(false);
  const error = ref('');

  const fetchSubscribers = async (pageID: number) => {
    loading.value = true;
    error.value = '';
    try {
      const res = await ExtendedFetch.get(`/status-pages/${pageID}/subscribers`);
      subscribers.value = res.data?.data || res.data || [];
    } catch (err: any) {
      console.error('Failed to fetch subscribers:', err);
      error.value = err.response?.data?.error || 'Failed to load subscribers.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const fetchStatusPages = async () => {
    try {
      const res = await ExtendedFetch.get('/status-pages');
      const pages = res.data?.data || res.data || [];
      statusPages.value = pages.map((p: any) => ({ id: p.id, name: p.name }));
    } catch (err: any) {
      console.error('Failed to fetch status pages:', err);
    }
  };

  const unsubscribe = async (pageID: number, email: string) => {
    error.value = '';
    try {
      await ExtendedFetch.post(`/status-pages/${pageID}/subscribers/unsubscribe`, { email });
      subscribers.value = subscribers.value.filter(s => s.email !== email);
    } catch (err: any) {
      console.error('Failed to unsubscribe:', err);
      error.value = err.response?.data?.error || 'Failed to unsubscribe.';
      throw err;
    }
  };

  return {
    subscribers,
    statusPages,
    loading,
    error,
    fetchSubscribers,
    fetchStatusPages,
    unsubscribe,
  };
});
