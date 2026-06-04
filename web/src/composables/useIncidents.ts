/**
 * Composable for Incident & Outage History Management
 * delegates to incidents Pinia store.
 */
import { storeToRefs } from 'pinia';
import { useIncidentsStore, type Incident, type MonitorItem } from '@/stores/incidents';

export type { Incident, MonitorItem };

export function useIncidents() {
  const store = useIncidentsStore();
  const { incidents, monitorsMap, loading, error } = storeToRefs(store);

  // Format timestamp helper
  const formatDate = (dateStr?: string) => {
    if (!dateStr) return '—';
    const d = new Date(dateStr);
    return d.toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  };

  // Calculate outage duration helper
  const getDuration = (inc: Incident) => {
    const start = new Date(inc.created_at).getTime();
    const end = inc.resolved_at ? new Date(inc.resolved_at).getTime() : Date.now();
    const diffSec = Math.floor((end - start) / 1000);

    if (diffSec < 60) return `${diffSec}s`;
    const diffMin = Math.floor(diffSec / 60);
    if (diffMin < 60) return `${diffMin}m ${diffSec % 60}s`;
    const diffHour = Math.floor(diffMin / 60);
    return `${diffHour}h ${diffMin % 60}m`;
  };

  return {
    incidents,
    monitorsMap,
    loading,
    error,
    fetchIncidentsData: store.fetchIncidentsData,
    formatDate,
    getDuration,
  };
}
