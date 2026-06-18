import { storeToRefs } from 'pinia';
import { useAuditLogsStore } from '@/stores/auditLogs';
export function useAuditLogs() {
  const store = useAuditLogsStore();
  const { items, loading } = storeToRefs(store);
  return { items, loading, fetchAll: store.fetchAll };
}
