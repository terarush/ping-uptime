import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';
import { useMonitorsStore } from '@/stores/monitors';

export interface Incident {
  id: number;
  monitor_id: number;
  user_id: number;
  status: string; // active, resolved
  error_message: string;
  latency: number;
  created_at: string;
  resolved_at?: string;
}

export interface MonitorItem {
  id: number;
  name: string;
  url: string;
}

export const useIncidentsStore = defineStore('incidents', () => {
  const incidents = ref<Incident[]>([]);
  const monitorsMap = ref<Record<number, MonitorItem>>({});
  const loading = ref(false);
  const error = ref('');

  const fetchIncidentsData = async () => {
    loading.value = true;
    error.value = '';
    try {
      // 1. Fetch monitors to resolve names/URLs
      const monitorsResponse = await ExtendedFetch.get('/monitors');
      const monitorsList: MonitorItem[] = monitorsResponse.data?.data || monitorsResponse.data || [];
      const map: Record<number, MonitorItem> = {};
      for (const m of monitorsList) {
        map[m.id] = m;
      }
      monitorsMap.value = map;

      // 2. Fetch incidents log
      const incidentsResponse = await ExtendedFetch.get('/incidents');
      incidents.value = incidentsResponse.data?.data || incidentsResponse.data || [];

      // Make sure SSE is running
      const monitorsStore = useMonitorsStore();
      monitorsStore.initializeSSE();
    } catch (err: any) {
      console.error('Failed to fetch incident data:', err);
      error.value = err.response?.data?.error || 'Failed to load incident history.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    incidents,
    monitorsMap,
    loading,
    error,
    fetchIncidentsData,
  };
});
