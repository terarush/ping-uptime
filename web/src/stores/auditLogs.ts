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

export interface AuditLogFilter {
  user_id?: number;
  entity_type?: string;
  action?: string;
  from?: string;
  to?: string;
  limit?: number;
}

export const useAuditLogsStore = defineStore('auditLogs', () => {
  const items = ref<AuditLog[]>([]);
  const loading = ref(false);
  const currentFilter = ref<AuditLogFilter>({});

  const fetchAll = async (filter?: AuditLogFilter) => {
    loading.value = true;
    currentFilter.value = filter || {};
    try {
      const params = new URLSearchParams();
      if (filter?.user_id) params.set('user_id', String(filter.user_id));
      if (filter?.entity_type) params.set('entity_type', filter.entity_type);
      if (filter?.action) params.set('action', filter.action);
      if (filter?.from) params.set('from', filter.from);
      if (filter?.to) params.set('to', filter.to);
      if (filter?.limit) params.set('limit', String(filter.limit));
      const qs = params.toString();
      const resp = await ExtendedFetch.get(`/audit-logs${qs ? '?' + qs : ''}`);
      items.value = resp.data?.data || resp.data || [];
    } finally {
      loading.value = false;
    }
  };

  return { items, loading, currentFilter, fetchAll };
});
