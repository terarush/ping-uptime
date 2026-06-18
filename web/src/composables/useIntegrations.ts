/**
 * Composable for Integration Management
 * delegates to integrations Pinia store.
 */
import { storeToRefs } from 'pinia';
import { useIntegrationsStore, type Integration } from '@/stores/integrations';

export type { Integration };

export function useIntegrations() {
  const store = useIntegrationsStore();
  const { integrations, loading, error } = storeToRefs(store);

  return {
    integrations,
    loading,
    error,
    fetchAll: store.fetchAll,
    create: store.create,
    update: store.update,
    deleteIntegration: store.deleteIntegration,
    test: store.test,
  };
}
