import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';
import { useIncidentsStore } from '@/stores/incidents';

export interface Monitor {
  id: number;
  name: string;
  url: string;
  type: string;
  interval: number;
  timeout: number;
  status: string; // active, paused
  uptime_status: string; // up, down, unknown
  last_checked_at?: string;
  user_id: number;
  created_at: string;
}

export const useMonitorsStore = defineStore('monitors', () => {
  const monitors = ref<Monitor[]>([]);
  const loading = ref(false);
  const error = ref('');
  let sse: EventSource | null = null;

  const initializeSSE = () => {
    if (sse) return; // Already listening

    const url = `${window.location.origin}/api/monitors/events`;
    sse = new EventSource(url);

    sse.onmessage = (event) => {
      try {
        if (!event.data || event.data === '{}') return;
        const data = JSON.parse(event.data);
        if (!data || !data.type) return;

        const incidentsStore = useIncidentsStore();

        if (data.type === 'monitor.checked' || data.type === 'monitor.updated' || data.type === 'monitor.created') {
          const mon = data.payload;
          const idx = monitors.value.findIndex(m => m.id === mon.id);
          if (idx !== -1) {
            monitors.value[idx] = { ...monitors.value[idx], ...mon };
          } else {
            monitors.value.push(mon);
          }
        } else if (data.type === 'monitor.deleted') {
          const mon = data.payload;
          monitors.value = monitors.value.filter(m => m.id !== mon.id);
        } else if (data.type === 'incident.created') {
          const inc = data.payload;
          const idx = incidentsStore.incidents.findIndex(i => i.id === inc.id);
          if (idx === -1) {
            incidentsStore.incidents.unshift(inc);
          } else {
            incidentsStore.incidents[idx] = { ...incidentsStore.incidents[idx], ...inc };
          }
        } else if (data.type === 'incident.resolved') {
          const inc = data.payload;
          const idx = incidentsStore.incidents.findIndex(i => i.id === inc.id);
          if (idx !== -1) {
            incidentsStore.incidents[idx] = { ...incidentsStore.incidents[idx], ...inc };
          } else {
            incidentsStore.incidents.unshift(inc);
          }
        }
      } catch (err) {
        console.error('Failed to parse real-time event data:', err);
      }
    };

    sse.onerror = (err) => {
      console.warn('Real-time connection lost, auto-reconnecting...', err);
    };
  };

  const fetchMonitors = async () => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.get('/monitors');
      monitors.value = response.data?.data || response.data || [];
      // Initialize SSE listener right after fetching initial monitors
      initializeSSE();
    } catch (err: any) {
      console.error('Failed to fetch monitors:', err);
      error.value = err.response?.data?.error || 'Failed to load monitors.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const createMonitor = async (payload: Omit<Monitor, 'id' | 'user_id' | 'created_at' | 'uptime_status'>) => {
    error.value = '';
    try {
      const response = await ExtendedFetch.post('/monitors', payload);
      const newMon = response.data?.data || response.data;
      return newMon;
    } catch (err: any) {
      console.error('Failed to create monitor:', err);
      error.value = err.response?.data?.error || 'Failed to create monitor.';
      throw err;
    }
  };

  const updateMonitor = async (id: number, payload: Partial<Monitor>) => {
    error.value = '';
    try {
      const response = await ExtendedFetch.put(`/monitors/${id}`, payload);
      const updatedMon = response.data?.data || response.data;
      return updatedMon;
    } catch (err: any) {
      console.error('Failed to update monitor:', err);
      error.value = err.response?.data?.error || 'Failed to update monitor.';
      throw err;
    }
  };

  const deleteMonitor = async (id: number) => {
    error.value = '';
    try {
      await ExtendedFetch.delete(`/monitors/${id}`);
      monitors.value = monitors.value.filter(m => m.id !== id);
    } catch (err: any) {
      console.error('Failed to delete monitor:', err);
      error.value = err.response?.data?.error || 'Failed to delete monitor.';
      throw err;
    }
  };

  return {
    monitors,
    loading,
    error,
    fetchMonitors,
    createMonitor,
    updateMonitor,
    deleteMonitor,
    initializeSSE,
  };
});
