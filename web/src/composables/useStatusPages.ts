/**
 * Composable for Status Page Management
 * delegates to statusPages Pinia store.
 */
import { storeToRefs } from 'pinia';
import { useStatusPagesStore, type StatusPage, type Monitor } from '@/stores/statusPages';

export type { StatusPage, Monitor };

export function useStatusPages() {
  const store = useStatusPagesStore();
  const { statusPages, availableMonitors, loading, error } = storeToRefs(store);

  return {
    statusPages,
    availableMonitors,
    loading,
    error,
    fetchStatusPagesData: store.fetchStatusPagesData,
    createStatusPage: store.createStatusPage,
    updateStatusPage: store.updateStatusPage,
    deleteStatusPage: store.deleteStatusPage,
  };
}
