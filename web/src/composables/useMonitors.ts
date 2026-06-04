/**
 * Composable for Monitor Management
 * delegates to monitors Pinia store.
 */
import { storeToRefs } from 'pinia';
import { useMonitorsStore, type Monitor } from '@/stores/monitors';

export type { Monitor };

export function useMonitors() {
  const store = useMonitorsStore();
  const { monitors, loading, error } = storeToRefs(store);

  return {
    monitors,
    loading,
    error,
    fetchMonitors: store.fetchMonitors,
    createMonitor: store.createMonitor,
    updateMonitor: store.updateMonitor,
    deleteMonitor: store.deleteMonitor,
  };
}
