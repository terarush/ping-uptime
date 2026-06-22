/**
 * Composable for Subscriber Management
 * delegates to subscribers Pinia store.
 */
import { storeToRefs } from 'pinia';
import { useSubscribersStore, type Subscriber } from '@/stores/subscribers';

export type { Subscriber };

export function useSubscribers() {
  const store = useSubscribersStore();
  const { subscribers, statusPages, loading, error } = storeToRefs(store);

  return {
    subscribers,
    statusPages,
    loading,
    error,
    fetchSubscribers: store.fetchSubscribers,
    fetchStatusPages: store.fetchStatusPages,
    unsubscribe: store.unsubscribe,
  };
}
