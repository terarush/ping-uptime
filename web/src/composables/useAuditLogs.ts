import { storeToRefs } from 'pinia';
import { useAuditLogsStore, type AuditLogFilter } from '@/stores/auditLogs';
export type { AuditLogFilter };
export function useAuditLogs() {
  const store = useAuditLogsStore();
  const { items, loading } = storeToRefs(store);
  return { items, loading, fetchAll: store.fetchAll };
}
