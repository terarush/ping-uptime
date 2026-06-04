/**
 * Composable for Settings Management
 * delegates to settings Pinia store.
 */
import { storeToRefs } from 'pinia';
import { useSettingsStore, type Setting } from '@/stores/settings';

export type { Setting };

export function useSettings() {
  const store = useSettingsStore();
  const { settings, loading, saveLoading, error } = storeToRefs(store);

  return {
    settings,
    loading,
    saveLoading,
    error,
    fetchSettingsData: store.fetchSettingsData,
    saveSettingsData: store.saveSettingsData,
  };
}
