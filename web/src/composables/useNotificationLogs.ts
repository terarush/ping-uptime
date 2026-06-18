/**
 * Composable for Notification Logs
 * delegates to notificationLogs Pinia store.
 */
import { storeToRefs } from 'pinia';
import { useNotificationLogsStore, type NotificationLog } from '@/stores/notificationLogs';

export type { NotificationLog };

export function useNotificationLogs() {
  const store = useNotificationLogsStore();
  const { logs, loading, error } = storeToRefs(store);

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return '—';
    return new Date(dateStr).toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const channelIcon = (type: string) => {
    const icons: Record<string, string> = {
      email: '📧',
      discord: '💬',
      discord_bot: '🤖',
      telegram: '📱',
      slack: '💼',
      webhook: '🔗',
    };
    return icons[type] || '📨';
  };

  return {
    logs,
    loading,
    error,
    fetchLogs: store.fetchLogs,
    formatDate,
    channelIcon,
  };
}
