/**
 * Composable for SSL Certificate Monitor Management
 * delegates to sslMonitors Pinia store.
 */
import { storeToRefs } from 'pinia';
import { useSslMonitorsStore, type SSLCert } from '@/stores/sslMonitors';

export type { SSLCert };

export function useSslMonitors() {
  const store = useSslMonitorsStore();
  const { certs, loading, error } = storeToRefs(store);

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

  // Days remaining color class
  const daysClass = (days: number) => {
    if (days <= 0) return 'text-red-500 font-bold';
    if (days < 14) return 'text-red-400 font-semibold';
    if (days <= 30) return 'text-amber-400 font-semibold';
    return 'text-emerald-500';
  };

  // Status badge variant
  const statusVariant = (status: string) => {
    switch (status) {
      case 'valid': return 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20';
      case 'expiring_soon': return 'bg-amber-500/10 text-amber-500 border-amber-500/20';
      case 'expired': return 'bg-red-500/10 text-red-500 border-red-500/20';
      case 'error': return 'bg-destructive/10 text-destructive border-destructive/20';
      default: return 'bg-muted text-muted-foreground border-border';
    }
  };

  return {
    certs,
    loading,
    error,
    fetchAll: store.fetchAll,
    fetchExpiring: store.fetchExpiring,
    checkMonitor: store.checkMonitor,
    checkAll: store.checkAll,
    deleteCert: store.deleteCert,
    formatDate,
    daysClass,
    statusVariant,
  };
}
