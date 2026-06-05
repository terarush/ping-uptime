import { storeToRefs } from 'pinia';
import { useAnalyticsStore, type ChartDataPoint, type MonitorStats } from '@/stores/analytics';

export type { ChartDataPoint, MonitorStats };

export function useAnalytics() {
  const analyticsStore = useAnalyticsStore();
  const { stats, chartPoints, chartMonitorID, chartWindow, loading, error } = storeToRefs(analyticsStore);

  return {
    stats,
    chartPoints,
    chartMonitorID,
    chartWindow,
    loading,
    error,
    fetchDashboardStats: analyticsStore.fetchDashboardStats,
    fetchChart: analyticsStore.fetchChart,
  };
}
