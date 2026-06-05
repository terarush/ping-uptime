import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface ChartDataPoint {
  date: string;
  status: string;
  uptime_pct: number;
  failed: number;
  total: number;
  latency: number;
}

export interface MonitorStats {
  monitor_id: number;
  monitor_name: string;
  monitor_url: string;
  window: string;
  uptime_pct: number;
  total_checks: number;
  failed_checks: number;
  points: ChartDataPoint[];
  status: string;
  avg_latency: number;
  min_latency: number;
  max_latency: number;
}

export const useAnalyticsStore = defineStore('analytics', () => {
  const stats = ref<MonitorStats[]>([]);
  const chartPoints = ref<ChartDataPoint[]>([]);
  const chartMonitorID = ref<number | null>(null);
  const chartWindow = ref('1m');
  const loading = ref(false);
  const error = ref('');

  const fetchDashboardStats = async (window = '1m') => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.get(`/analytics/dashboard?window=${window}`);
      stats.value = response.data?.data?.data || [];
    } catch (err: any) {
      console.error('Failed to fetch dashboard stats:', err);
      error.value = err.response?.data?.error || 'Failed to load analytics.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const fetchChart = async (monitorID: number, window = '1m') => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.get(`/analytics/monitors/${monitorID}/chart?window=${window}`);
      chartPoints.value = response.data?.data?.data || response.data?.data?.points || [];
      chartMonitorID.value = monitorID;
      chartWindow.value = window;
    } catch (err: any) {
      console.error('Failed to fetch chart:', err);
      error.value = err.response?.data?.error || 'Failed to load chart.';
      chartPoints.value = [];
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    stats,
    chartPoints,
    chartMonitorID,
    chartWindow,
    loading,
    error,
    fetchDashboardStats,
    fetchChart,
  };
});
