import { storeToRefs } from 'pinia';
import { useMaintenancesStore, type Maintenance } from '@/stores/maintenances';

export type { Maintenance };

export function useMaintenances() {
  const store = useMaintenancesStore();
  const { items, loading, error } = storeToRefs(store);
  return {
    items, loading, error,
    fetchAll: store.fetchAll,
    create: store.create,
    update: store.update,
    remove: store.remove,
  };
}
