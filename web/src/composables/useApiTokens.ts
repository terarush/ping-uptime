/**
 * Composable for API Token Management
 * delegates to apiTokens Pinia store.
 */
import { storeToRefs } from 'pinia';
import { useApiTokensStore, type ApiToken, type CreateTokenResult } from '@/stores/apiTokens';

export type { ApiToken, CreateTokenResult };

export function useApiTokens() {
  const store = useApiTokensStore();
  const { tokens, loading, error } = storeToRefs(store);

  const formatDate = (dateStr?: string | null) => {
    if (!dateStr) return '—';
    const d = new Date(dateStr);
    return d.toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return {
    tokens,
    loading,
    error,
    fetchAll: store.fetchAll,
    create: store.create,
    revoke: store.revoke,
    formatDate,
  };
}
