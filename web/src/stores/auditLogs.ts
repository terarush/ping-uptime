import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface AuditLog {
  id: number;
  user_id: number;
  action: string;
  entity_type: string;
  entity_id: number;
  details: string;
  created_at: string;
}

export const useAuditLogsStore = defineStore('auditLogs', () => {
  const items = ref<AuditLog[]>([]);
  const loading = ref(false);

  const fetchAll = async () => {
    loading.value = true;
    try {
      const resp = await ExtendedFetch.get('/audit-logs');
      items.value = resp.data?.data || resp.data || [];
    } finally {
      loading.value = false;
    }
  };

  return { items, loading, fetchAll };
});
