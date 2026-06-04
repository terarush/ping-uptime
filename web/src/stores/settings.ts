import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface Setting {
  key: string;
  value: string;
  description: string;
}

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<Setting[]>([]);
  const loading = ref(false);
  const saveLoading = ref(false);
  const error = ref('');

  const fetchSettingsData = async () => {
    loading.value = true;
    error.value = '';
    try {
      const response = await ExtendedFetch.get('/settings');
      settings.value = response.data?.data || response.data || [];
    } catch (err: any) {
      console.error('Failed to fetch settings:', err);
      error.value = err.response?.data?.error || 'Failed to load system settings.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const saveSettingsData = async (payloads: Array<{ key: string; value: string; description: string }>) => {
    saveLoading.value = true;
    error.value = '';
    try {
      // Gorm backend POST saves/updates settings
      for (const p of payloads) {
        await ExtendedFetch.post('/settings', p);
      }
    } catch (err: any) {
      console.error('Failed to save settings:', err);
      error.value = err.response?.data?.error || 'Failed to save system settings.';
      throw err;
    } finally {
      saveLoading.value = false;
    }
  };

  return {
    settings,
    loading,
    saveLoading,
    error,
    fetchSettingsData,
    saveSettingsData,
  };
});
