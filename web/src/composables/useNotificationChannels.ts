/**
 * Composable for Notification Channel Management
 * delegates to notifications Pinia store.
 */
import { storeToRefs } from 'pinia';
import { useNotificationsStore, type NotificationChannel } from '@/stores/notifications';

export type { NotificationChannel };

export function useNotificationChannels() {
  const store = useNotificationsStore();
  const { channels, loading, error } = storeToRefs(store);

  return {
    channels,
    loading,
    error,
    fetchChannelsData: store.fetchChannelsData,
    createChannel: store.createChannel,
    updateChannel: store.updateChannel,
    deleteChannel: store.deleteChannel,
  };
}
