import { storeToRefs } from 'pinia';
import { useBackupStore } from '@/stores/backup';
export function useBackup() {
  const store = useBackupStore();
  const { records, loading } = storeToRefs(store);
  return { records, loading, exportBackup: store.exportBackup, downloadBackup: store.downloadBackup, importBackup: store.importBackup, getHistory: store.getHistory, deleteRecord: store.deleteRecord };
}
