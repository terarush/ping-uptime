<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useMaintenances, type Maintenance } from '@/composables/useMaintenances';
import { useMonitors } from '@/composables/useMonitors';
import MaintenanceTable from '@/components/maintenance-table.vue';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog';
import { toast } from 'vue-sonner';
import { Wrench, Plus, Loader2, RefreshCw } from '@lucide/vue';

const { items, loading, fetchAll, create, update, remove } = useMaintenances();
const { monitors, fetchMonitors } = useMonitors();

const isFormOpen = ref(false);
const isDeleteOpen = ref(false);
const formLoading = ref(false);
const actionItem = ref<Maintenance | null>(null);

const formName = ref('');
const formDesc = ref('');
const formStart = ref('');
const formEnd = ref('');
const formMonitorIDs = ref<number[]>([]);

const isEdit = computed(() => !!actionItem.value);

const resetForm = () => {
  formName.value = '';
  formDesc.value = '';
  formStart.value = '';
  formEnd.value = '';
  formMonitorIDs.value = [];
  actionItem.value = null;
};

const openCreate = () => {
  resetForm();
  isFormOpen.value = true;
};

const openEdit = (item: Maintenance) => {
  actionItem.value = item;
  formName.value = item.name;
  formDesc.value = item.description;
  formStart.value = item.start_at.slice(0, 16);
  formEnd.value = item.end_at.slice(0, 16);
  formMonitorIDs.value = item.monitor_ids || [];
  isFormOpen.value = true;
};

const openDelete = (item: Maintenance) => {
  actionItem.value = item;
  isDeleteOpen.value = true;
};

const handleSubmit = async () => {
  formLoading.value = true;
  try {
    const payload = {
      name: formName.value,
      description: formDesc.value,
      start_at: new Date(formStart.value).toISOString(),
      end_at: new Date(formEnd.value).toISOString(),
      monitor_ids: formMonitorIDs.value,
    };

    if (isEdit.value && actionItem.value) {
      await update(actionItem.value.id, payload);
      toast.success(`Maintenance "${formName.value}" updated`);
    } else {
      await create(payload);
      toast.success(`Maintenance "${formName.value}" created`);
    }

    isFormOpen.value = false;
    await fetchAll();
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to save maintenance');
  } finally {
    formLoading.value = false;
  }
};

const handleDelete = async () => {
  if (!actionItem.value) return;
  formLoading.value = true;
  try {
    await remove(actionItem.value.id);
    toast.success(`Maintenance "${actionItem.value.name}" deleted`);
    isDeleteOpen.value = false;
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to delete');
  } finally {
    formLoading.value = false;
  }
};

const toggleMonitor = (id: number) => {
  const idx = formMonitorIDs.value.indexOf(id);
  if (idx === -1) formMonitorIDs.value.push(id);
  else formMonitorIDs.value.splice(idx, 1);
};

onMounted(async () => {
  await Promise.all([fetchAll(), fetchMonitors()]);
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-amber-500/10 dark:bg-amber-500/5 blur-[100px] pointer-events-none"></div>

    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <Wrench class="w-6 h-6 text-primary" />
          <span>Maintenance Windows</span>
        </h2>
        <p class="text-xs text-muted-foreground">Schedule planned maintenance to suppress alerts.</p>
      </div>
      <div class="flex items-center gap-2">
        <Button variant="outline" size="sm" @click="fetchAll" class="h-9">
          <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
          Refresh
        </Button>
        <Button @click="openCreate" size="sm" class="h-9">
          <Plus class="w-4 h-4 mr-1.5" />
          Add Maintenance
        </Button>
      </div>
    </div>

    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <CardTitle class="text-sm font-bold">Scheduled Maintenances</CardTitle>
        <CardDescription class="text-xs">Active maintenance windows suppress alert notifications.</CardDescription>
      </CardHeader>
      <CardContent class="p-0">
        <div v-if="loading && items.length === 0" class="flex justify-center py-20">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
        </div>
        <div v-else-if="items.length === 0" class="flex flex-col items-center py-20 text-center">
          <Wrench class="w-12 h-12 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-bold text-foreground">No maintenance windows</p>
          <p class="text-xs text-muted-foreground mt-1">Schedule maintenance to avoid false alerts.</p>
        </div>
        <MaintenanceTable v-else :items="items" @edit="openEdit" @delete="openDelete" />
      </CardContent>
    </Card>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:open="isFormOpen">
      <DialogContent class="sm:max-w-106.25">
        <DialogHeader>
          <DialogTitle>{{ isEdit ? 'Edit Maintenance' : 'Schedule Maintenance' }}</DialogTitle>
          <DialogDescription>Set start and end times to suppress alerts for selected monitors.</DialogDescription>
        </DialogHeader>
        <form @submit.prevent="handleSubmit" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="mname">Name</Label>
            <Input id="mname" v-model="formName" placeholder="e.g. Database migration" required />
          </div>
          <div class="space-y-2">
            <Label for="mdesc">Description</Label>
            <Input id="mdesc" v-model="formDesc" placeholder="Brief description" />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="mstart">Start</Label>
              <Input id="mstart" v-model="formStart" type="datetime-local" required />
            </div>
            <div class="space-y-2">
              <Label for="mend">End</Label>
              <Input id="mend" v-model="formEnd" type="datetime-local" required />
            </div>
          </div>
          <div class="space-y-2">
            <Label>Affected Monitors</Label>
            <div class="flex flex-wrap gap-2 max-h-32 overflow-y-auto border border-border/40 rounded-lg p-2">
              <button v-for="m in monitors" :key="m.id" type="button"
                :class="['px-2.5 py-1 text-[11px] font-semibold rounded-full border transition-colors',
                  formMonitorIDs.includes(m.id)
                    ? 'bg-primary/10 text-primary border-primary/30'
                    : 'bg-muted text-muted-foreground border-border/40 hover:border-primary/30'
                ]"
                @click="toggleMonitor(m.id)">
                {{ m.name }}
              </button>
            </div>
          </div>
          <DialogFooter class="pt-4 border-t border-border/40">
            <Button type="button" variant="outline" @click="isFormOpen = false">Cancel</Button>
            <Button type="submit" :disabled="formLoading">
              <Loader2 v-if="formLoading" class="w-4 h-4 mr-1.5 animate-spin" />
              {{ isEdit ? 'Update' : 'Create' }}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>

    <!-- Delete Dialog -->
    <Dialog v-model:open="isDeleteOpen">
      <DialogContent class="sm:max-w-100">
        <DialogHeader>
          <DialogTitle class="text-destructive">Delete Maintenance</DialogTitle>
          <DialogDescription>Remove "{{ actionItem?.name }}"? This cannot be undone.</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="isDeleteOpen = false">Cancel</Button>
          <Button variant="destructive" :disabled="formLoading" @click="handleDelete">
            <Loader2 v-if="formLoading" class="w-4 h-4 mr-1.5 animate-spin" />
            Delete
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
