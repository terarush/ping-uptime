import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface Tag {
  id: number;
  name: string;
  color: string;
  created_at: string;
}

export const useTagsStore = defineStore('tags', () => {
  const tags = ref<Tag[]>([]);
  const loading = ref(false);
  const error = ref('');

  const fetchAll = async () => {
    loading.value = true;
    error.value = '';
    try {
      const res = await ExtendedFetch.get('/tags');
      tags.value = res.data?.data || res.data || [];
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to load tags.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const create = async (payload: { name: string; color?: string }) => {
    error.value = '';
    try {
      const res = await ExtendedFetch.post('/tags', payload);
      const tag = res.data?.data || res.data;
      if (tag) tags.value.push(tag);
      return tag;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to create tag.';
      throw err;
    }
  };

  const update = async (id: number, payload: { name?: string; color?: string }) => {
    error.value = '';
    try {
      const res = await ExtendedFetch.put(`/tags/${id}`, payload);
      const updated = res.data?.data || res.data;
      const idx = tags.value.findIndex(t => t.id === id);
      if (idx !== -1) tags.value[idx] = updated;
      return updated;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to update tag.';
      throw err;
    }
  };

  const deleteTag = async (id: number) => {
    error.value = '';
    try {
      await ExtendedFetch.delete(`/tags/${id}`);
      tags.value = tags.value.filter(t => t.id !== id);
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to delete tag.';
      throw err;
    }
  };

  const attachToMonitor = async (monitorID: number, tagIDs: number[]) => {
    error.value = '';
    try {
      await ExtendedFetch.post(`/monitors/${monitorID}/tags`, { tag_ids: tagIDs });
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to attach tags.';
      throw err;
    }
  };

  const getMonitorTags = async (monitorID: number) => {
    error.value = '';
    try {
      const res = await ExtendedFetch.get(`/monitors/${monitorID}/tags`);
      return res.data?.data || res.data || [];
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to get monitor tags.';
      throw err;
    }
  };

  const detachFromMonitor = async (monitorID: number, tagID: number) => {
    error.value = '';
    try {
      await ExtendedFetch.delete(`/monitors/${monitorID}/tags/${tagID}`);
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to detach tag.';
      throw err;
    }
  };

  return {
    tags,
    loading,
    error,
    fetchAll,
    create,
    update,
    deleteTag,
    attachToMonitor,
    getMonitorTags,
    detachFromMonitor,
  };
});
