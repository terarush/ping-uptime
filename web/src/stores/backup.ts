import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface BackupRecord {
  id: number;
  file_name: string;
  file_size: number;
  user_id: number;
  created_at: string;
}

export const useBackupStore = defineStore('backup', () => {
  const records = ref<BackupRecord[]>([]);
  const loading = ref(false);

  const exportBackup = async () => {
    const resp = await ExtendedFetch.get('/backup/export');
    return resp.data?.data || resp.data;
  };

  const downloadBackup = async () => {
    const resp = await ExtendedFetch.get('/backup/download', { responseType: 'blob' });
    const url = URL.createObjectURL(resp.data);
    const a = document.createElement('a');
    a.href = url;
    const disposition = resp.headers?.['content-disposition'] || '';
    const match = disposition.match(/filename=([^;]+)/);
    a.download = match?.[1] || 'ping-uptime-backup.json';
    a.click();
    URL.revokeObjectURL(url);
  };

  const importBackup = async (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    const resp = await ExtendedFetch.uploadForm('/backup/import', formData);
    await getHistory();
    return resp;
  };

  const getHistory = async () => {
    loading.value = true;
    try {
      const resp = await ExtendedFetch.get('/backup/history');
      records.value = resp.data?.data || resp.data || [];
    } finally {
      loading.value = false;
    }
  };

  const deleteRecord = async (id: number) => {
    await ExtendedFetch.delete(`/backup/history/${id}`);
    records.value = records.value.filter(r => r.id !== id);
  };

  return { records, loading, exportBackup, downloadBackup, importBackup, getHistory, deleteRecord };
});
