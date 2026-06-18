import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface NotificationLog {
  id: number;
  user_id: number;
  notification_id: number;
  channel_type: string;
  monitor_id: number | null;
  incident_id: number | null;
  event_type: string; // alert, resolution, test
  status: string; // sent, failed
  error_message: string;
  recipient: string;
  sent_at: string;
}

export const useNotificationLogsStore = defineStore('notificationLogs', () => {
  const logs = ref<NotificationLog[]>([]);
  const loading = ref(false);
  const error = ref('');

  const fetchLogs = async (filters?: Record<string, string>) => {
    loading.value = true;
    error.value = '';
    try {
      let url = '/notification-logs';
      if (filters) {
        const params = new URLSearchParams();
        Object.entries(filters).forEach(([k, v]) => {
          if (v) params.set(k, v);
        });
        const qs = params.toString();
        if (qs) url += '?' + qs;
      }
      const res = await ExtendedFetch.get(url);
      logs.value = res.data?.data || res.data || [];
    } catch (err: any) {
      console.error('Failed to fetch notification logs:', err);
      error.value = err.response?.data?.error || 'Failed to load notification logs.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    logs,
    loading,
    error,
    fetchLogs,
  };
});
