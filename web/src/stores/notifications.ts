import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface NotificationChannel {
  id: number;
  name: string;
  type: string; // email, webhook, telegram, discord, slack
  config: string; // JSON string
  enabled: boolean;
  user_id: number;
  created_at: string;
}

export const useNotificationsStore = defineStore('notifications', () => {
  const channels = ref<NotificationChannel[]>([]);
  const loading = ref(false);
  const error = ref('');

  const fetchChannelsData = async () => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.get('/notification-channels');
      channels.value = response.data?.data || response.data || [];
    } catch (err: any) {
      console.error('Failed to fetch notification channels:', err);
      error.value = err.response?.data?.error || 'Failed to load alert channels.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const createChannel = async (payload: any) => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.post('/notification-channels', payload);
      return response.data;
    } catch (err: any) {
      console.error('Failed to create notification channel:', err);
      error.value = err.response?.data?.error || 'Failed to create alert channel.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const updateChannel = async (id: number, payload: any) => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.put(`/notification-channels/${id}`, payload);
      return response.data;
    } catch (err: any) {
      console.error('Failed to update notification channel:', err);
      error.value = err.response?.data?.error || 'Failed to update alert channel.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const deleteChannel = async (id: number) => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.delete(`/notification-channels/${id}`);
      channels.value = channels.value.filter(c => c.id !== id);
      return response.data;
    } catch (err: any) {
      console.error('Failed to delete notification channel:', err);
      error.value = err.response?.data?.error || 'Failed to delete alert channel.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    channels,
    loading,
    error,
    fetchChannelsData,
    createChannel,
    updateChannel,
    deleteChannel,
  };
});
