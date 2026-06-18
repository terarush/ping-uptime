/**
 * Composable for Tag Management
 * delegates to tags Pinia store.
 */
import { storeToRefs } from 'pinia';
import { useTagsStore, type Tag } from '@/stores/tags';

export type { Tag };

export function useTags() {
  const store = useTagsStore();
  const { tags, loading, error } = storeToRefs(store);

  return {
    tags,
    loading,
    error,
    fetchAll: store.fetchAll,
    create: store.create,
    update: store.update,
    deleteTag: store.deleteTag,
    attachToMonitor: store.attachToMonitor,
    getMonitorTags: store.getMonitorTags,
    detachFromMonitor: store.detachFromMonitor,
  };
}
