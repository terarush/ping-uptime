import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface SSLCert {
  id: number;
  monitor_id: number;
  domain: string;
  issuer: string;
  valid_from: string;
  valid_until: string;
  days_remaining: number;
  status: string; // valid, expiring_soon, expired, error
  error_message: string;
  checked_at: string;
  created_at: string;
}

export const useSslMonitorsStore = defineStore('sslMonitors', () => {
  const certs = ref<SSLCert[]>([]);
  const loading = ref(false);
  const error = ref('');

  const fetchAll = async () => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.get('/ssl-monitors');
      certs.value = response.data?.data || response.data || [];
    } catch (err: any) {
      console.error('Failed to fetch SSL certs:', err);
      error.value = err.response?.data?.error || 'Failed to load SSL certificates.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const fetchExpiring = async (days: number = 30) => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.get(`/ssl-monitors/expiring?days=${days}`);
      certs.value = response.data?.data || response.data || [];
    } catch (err: any) {
      console.error('Failed to fetch expiring SSL certs:', err);
      error.value = err.response?.data?.error || 'Failed to load expiring SSL certificates.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const checkMonitor = async (monitorId: number) => {
    try {
      const response = await ExtendedFetch.post(`/ssl-monitors/check/${monitorId}`);
      const cert = response.data?.data || response.data;
      // Update in list or append
      const idx = certs.value.findIndex(c => c.monitor_id === monitorId);
      if (idx >= 0) {
        certs.value[idx] = cert;
      } else if (cert) {
        certs.value.push(cert);
      }
      return cert;
    } catch (err: any) {
      console.error('Failed to check SSL for monitor:', err);
      throw err;
    }
  };

  const checkAll = async () => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.post('/ssl-monitors/check-all');
      // Refetch the full list after check-all completes
      await fetchAll();
      return response.data?.data || response.data || [];
    } catch (err: any) {
      console.error('Failed to check all SSL certs:', err);
      error.value = err.response?.data?.error || 'Failed to check SSL certificates.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const deleteCert = async (id: number) => {
    try {
      await ExtendedFetch.delete(`/ssl-monitors/${id}`);
      certs.value = certs.value.filter(c => c.id !== id);
    } catch (err: any) {
      console.error('Failed to delete SSL cert:', err);
      throw err;
    }
  };

  return {
    certs,
    loading,
    error,
    fetchAll,
    fetchExpiring,
    checkMonitor,
    checkAll,
    deleteCert,
  };
});
