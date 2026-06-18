import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface Maintenance {
  id: number;
  name: string;
  description: string;
  start_at: string;
  end_at: string;
  status: string; // scheduled, ongoing, completed
  user_id: number;
  monitor_ids: number[];
  created_at: string;
  updated_at: string;
}

export const useMaintenancesStore = defineStore('maintenances', () => {
  const items = ref<Maintenance[]>([]);
  const loading = ref(false);
  const error = ref('');

  const fetchAll = async () => {
    loading.value = true;
    error.value = '';
    try {
      const resp = await ExtendedFetch.get('/maintenances');
      items.value = resp.data?.data || resp.data || [];
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to load maintenances.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const create = async (payload: any) => {
    const resp = await ExtendedFetch.post('/maintenances', payload);
    return resp.data?.data || resp.data;
  };

  const update = async (id: number, payload: any) => {
    const resp = await ExtendedFetch.put(`/maintenances/${id}`, payload);
    return resp.data?.data || resp.data;
  };

  const remove = async (id: number) => {
    await ExtendedFetch.delete(`/maintenances/${id}`);
    items.value = items.value.filter(i => i.id !== id);
  };

  return { items, loading, error, fetchAll, create, update, remove };
});
