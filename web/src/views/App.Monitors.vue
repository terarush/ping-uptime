<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useAuth } from '@/composables/useAuth';
import { useMonitors, type Monitor } from '@/composables/useMonitors';
import { monitorSchema } from '@/validations/monitor';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import MonitorTable from '@/components/monitor-table.vue';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { toast } from 'vue-sonner';
import {
  Activity,
  Plus,
  Trash2,
  Search,
  Loader2,
  RefreshCw,
  TrendingUp,
  BarChart3
} from '@lucide/vue';
import gsap from 'gsap';
import { useAnalytics } from '@/composables/useAnalytics';
import Chart from '@/components/chart.vue';

// Auth checks
const { currentUser } = useAuth();
const isAdmin = computed(() => currentUser.value?.role === 'admin');

const {
  stats,
  chartPoints,
  chartMonitorID,
  chartWindow,
  loading: analyticsLoading,
  error: analyticsError,
  fetchDashboardStats,
  fetchChart,
} = useAnalytics();

// Composable monitor states
const {
  monitors,
  loading,
  error,
  fetchMonitors,
  createMonitor,
  updateMonitor,
  deleteMonitor
} = useMonitors();

const searchQuery = ref('');
const success = ref('');

const selectedMonitorId = ref<number | null>(null);
const selectedWindow = ref<'1h' | '1d' | '1w' | '1m' | '1y' | 'all'>('1m');

// Dialog states
const isFormDialogOpen = ref(false);
const isDeleteDialogOpen = ref(false);
const formLoading = ref(false);
const actionMonitor = ref<Monitor | null>(null);

// Form fields
const formName = ref('');
const formUrl = ref('');
const formType = ref('http');
const formInterval = ref(60);
const formTimeout = ref(10);
const formStatus = ref('active');

const isEditMode = computed(() => !!actionMonitor.value);

// Fetch wrapper with animation callback
const fetchAll = async () => {
  try {
    await fetchMonitors();
  } catch (err) {
    console.error('Failed to load monitors list:', err);
  } finally {
    setTimeout(animateTableRows, 50);
  }
};

// Filter monitors list
const filteredMonitors = computed(() => {
  if (!monitors.value) return [];
  const query = searchQuery.value.toLowerCase().trim();
  if (!query) return monitors.value;
  return monitors.value.filter(m =>
    m.name.toLowerCase().includes(query) ||
    m.url.toLowerCase().includes(query) ||
    m.type.toLowerCase().includes(query) ||
    m.status.toLowerCase().includes(query)
  );
});

// Reset form fields
const resetForm = () => {
  formName.value = '';
  formUrl.value = '';
  formType.value = 'http';
  formInterval.value = 60;
  formTimeout.value = 10;
  formStatus.value = 'active';
  actionMonitor.value = null;
  error.value = '';
  success.value = '';
};

const openCreateDialog = () => {
  resetForm();
  isFormDialogOpen.value = true;
};

const openEditDialog = (monitor: Monitor) => {
  resetForm();
  actionMonitor.value = monitor;
  formName.value = monitor.name;
  formUrl.value = monitor.url;
  formType.value = monitor.type;
  formInterval.value = monitor.interval;
  formTimeout.value = monitor.timeout;
  formStatus.value = monitor.status;
  isFormDialogOpen.value = true;
};

const openDeleteDialog = (monitor: Monitor) => {
  actionMonitor.value = monitor;
  isDeleteDialogOpen.value = true;
};

// Submit Form (Create / Update)
const handleFormSubmit = async () => {
  error.value = '';
  formLoading.value = true;

  const rawPayload = {
    name: formName.value,
    url: formUrl.value,
    type: formType.value,
    interval: isAdmin.value ? Number(formInterval.value) : 60,
    timeout: isAdmin.value ? Number(formTimeout.value) : 10,
    status: formStatus.value,
  };

  // Perform client-side Zod validation
  const validation = monitorSchema.safeParse(rawPayload);
  if (!validation.success) {
    const firstError = validation.error.errors[0]?.message || 'Validation failed';
    toast.error(firstError);
    error.value = firstError;
    formLoading.value = false;
    return;
  }

  try {
    if (isEditMode.value && actionMonitor.value) {
      await updateMonitor(actionMonitor.value.id, rawPayload);
      toast.success(`Monitor "${formName.value}" updated successfully!`);
    } else {
      await createMonitor(rawPayload);
      toast.success(`Monitor "${formName.value}" created successfully!`);
    }

    isFormDialogOpen.value = false;
    await fetchAll();
  } catch (err: any) {
    console.error('Failed to save monitor:', err);
    const msg = err.response?.data?.error || 'Failed to save monitor data.';
    toast.error(msg);
    error.value = msg;
  } finally {
    formLoading.value = false;
  }
};

// Delete monitor
const handleDeleteConfirm = async () => {
  if (!actionMonitor.value) return;
  error.value = '';
  formLoading.value = true;

  try {
    await deleteMonitor(actionMonitor.value.id);
    toast.success(`Monitor "${actionMonitor.value.name}" deleted successfully!`);
    isDeleteDialogOpen.value = false;
    await fetchAll();
  } catch (err: any) {
    console.error('Failed to delete monitor:', err);
    const msg = err.response?.data?.error || 'Failed to delete monitor.';
    toast.error(msg);
    error.value = msg;
  } finally {
    formLoading.value = false;
  }
};

// GSAP Animations
const animateTableRows = () => {
  gsap.fromTo('.monitor-row',
    { opacity: 0, y: 15 },
    { opacity: 1, y: 0, duration: 0.4, stagger: 0.05, ease: 'power2.out' }
  );
};

onMounted(() => {
  fetchAll();
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <!-- Ambient Background Orb -->
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-emerald-500/10 dark:bg-emerald-500/5 blur-[100px] pointer-events-none"></div>

    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <Activity class="w-6 h-6 text-primary" />
          <span>Uptime Monitors</span>
        </h2>
        <p class="text-xs text-muted-foreground">Manage endpoints and ping check configurations.</p>
      </div>

      <!-- Actions -->
      <div class="flex items-center gap-2">
        <Button variant="outline" size="sm" @click="fetchAll" class="h-9">
          <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
          <span>Refresh</span>
        </Button>
        <Button @click="openCreateDialog" size="sm" class="h-9 shadow-sm shadow-primary/10">
          <Plus class="w-4 h-4 mr-1.5" />
          <span>Add Monitor</span>
        </Button>
      </div>
    </div>

    <!-- Main Card -->
    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <CardTitle class="text-sm font-bold text-foreground">Registered Monitors</CardTitle>
            <CardDescription class="text-xs">Uptime checks run based on configured intervals.</CardDescription>
          </div>

          <!-- Search Bar -->
          <div class="relative w-full sm:w-72">
            <Search class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              v-model="searchQuery"
              placeholder="Search monitors..."
              class="pl-9 h-9"
            />
          </div>
        </div>
      </CardHeader>

      <CardContent class="p-0">
        <div v-if="loading && monitors.length === 0" class="flex flex-col items-center justify-center py-20 gap-3">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
          <p class="text-sm text-muted-foreground">Loading monitors...</p>
        </div>

        <div v-else-if="filteredMonitors.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
          <Activity class="w-12 h-12 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-bold text-foreground">No monitors found</p>
          <p class="text-xs text-muted-foreground mt-1">Try refining your search or add a new monitor.</p>
        </div>

        <MonitorTable
          v-else
          :monitors="filteredMonitors"
          @edit="openEditDialog"
          @delete="openDeleteDialog"
        />
      </CardContent>
    </Card>

    <!-- Analytics Detail Section -->
    <Card v-if="selectedMonitorId" class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <CardTitle class="text-sm font-bold text-foreground flex items-center gap-2">
              <TrendingUp class="w-4 h-4 text-primary" />
              <span>Analytics</span>
            </CardTitle>
            <CardDescription class="text-xs">Uptime history and response timeline for the selected monitor.</CardDescription>
          </div>

          <div class="flex items-center gap-2">
            <Select v-model="selectedWindow">
              <SelectTrigger class="h-9 w-[140px]">
                <SelectValue placeholder="Window" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="1h">Last 1 hour</SelectItem>
                <SelectItem value="1d">Last 24 hours</SelectItem>
                <SelectItem value="1w">Last 7 days</SelectItem>
                <SelectItem value="1m">Last 30 days</SelectItem>
                <SelectItem value="1y">Last year</SelectItem>
                <SelectItem value="all">All time</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </CardHeader>

      <CardContent class="p-6 space-y-6">
        <div v-if="analyticsLoading && chartPoints.length === 0" class="flex flex-col items-center justify-center py-16 gap-3">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
          <p class="text-sm text-muted-foreground">Loading analytics...</p>
        </div>

        <template v-else>
          <Chart :points="chartPoints" :height="220" />

          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div class="rounded-lg border border-border/60 bg-background p-3">
              <p class="text-[10px] font-bold text-muted-foreground uppercase">Uptime</p>
              <p class="text-sm font-black text-foreground mt-1">{{ (stats[0]?.uptime_pct ?? 0).toFixed(2) }}%</p>
            </div>
            <div class="rounded-lg border border-border/60 bg-background p-3">
              <p class="text-[10px] font-bold text-muted-foreground uppercase">Total checks</p>
              <p class="text-sm font-black text-foreground mt-1">{{ stats[0]?.total_checks ?? 0 }}</p>
            </div>
            <div class="rounded-lg border border-border/60 bg-background p-3">
              <p class="text-[10px] font-bold text-muted-foreground uppercase">Failed</p>
              <p class="text-sm font-black text-foreground mt-1">{{ stats[0]?.failed_checks ?? 0 }}</p>
            </div>
            <div class="rounded-lg border border-border/60 bg-background p-3">
              <p class="text-[10px] font-bold text-muted-foreground uppercase">Status</p>
              <p class="text-sm font-black text-foreground mt-1">{{ stats[0]?.status ?? '—' }}</p>
            </div>
          </div>

          <div v-if="analyticsError" class="text-xs text-destructive">{{ analyticsError }}</div>
        </template>
      </CardContent>
    </Card>

    <!-- Create / Edit Dialog -->
    <Dialog v-model:open="isFormDialogOpen">
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{{ isEditMode ? 'Edit Uptime Monitor' : 'Create Uptime Monitor' }}</DialogTitle>
          <DialogDescription>
            Configure an endpoint to monitor its responsiveness and response latency.
          </DialogDescription>
        </DialogHeader>

        <form @submit.prevent="handleFormSubmit" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="name">Friendly Name</Label>
            <Input id="name" v-model="formName" placeholder="e.g. My Website Gateway" required />
          </div>

          <div class="space-y-2">
            <Label for="url">Target URL</Label>
            <Input id="url" v-model="formUrl" placeholder="e.g. https://example.com/health" type="url" required />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="type">Check Type</Label>
              <Select v-model="formType">
                <SelectTrigger id="type" class="h-9">
                  <SelectValue placeholder="Select type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="http">HTTP/HTTPS</SelectItem>
                  <SelectItem value="ping">ICMP Ping</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div class="space-y-2">
              <Label for="status">Status</Label>
              <Select v-model="formStatus">
                <SelectTrigger id="status" class="h-9">
                  <SelectValue placeholder="Select status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="paused">Paused</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4" v-if="isAdmin">
            <div class="space-y-2">
              <Label for="interval">Interval (seconds)</Label>
              <Input id="interval" v-model="formInterval" type="number" min="10" required />
            </div>

            <div class="space-y-2">
              <Label for="timeout">Timeout (seconds)</Label>
              <Input id="timeout" v-model="formTimeout" type="number" min="1" required />
            </div>
          </div>

          <p v-if="error" class="text-xs font-semibold text-destructive mt-2">{{ error }}</p>

          <DialogFooter class="pt-4 border-t border-border/40 mt-4">
            <Button type="button" variant="outline" @click="isFormDialogOpen = false">Cancel</Button>
            <Button type="submit" :disabled="formLoading" class="min-w-24">
              <Loader2 v-if="formLoading" class="w-4 h-4 mr-1.5 animate-spin" />
              <span>Save</span>
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>

    <!-- Delete Confirmation Dialog -->
    <Dialog v-model:open="isDeleteDialogOpen">
      <DialogContent class="sm:max-w-[400px]">
        <DialogHeader>
          <DialogTitle class="text-destructive flex items-center gap-2">
            <Trash2 class="w-5 h-5" />
            <span>Delete Monitor</span>
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete the monitor "{{ actionMonitor?.name }}"? This action will permanently remove its configurations and logged statistics.
          </DialogDescription>
        </DialogHeader>

        <p v-if="error" class="text-xs font-semibold text-destructive mt-2">{{ error }}</p>

        <DialogFooter class="pt-4 border-t border-border/40 mt-4">
          <Button type="button" variant="outline" @click="isDeleteDialogOpen = false">Cancel</Button>
          <Button type="button" variant="destructive" :disabled="formLoading" @click="handleDeleteConfirm" class="min-w-24">
            <Loader2 v-if="formLoading" class="w-4 h-4 mr-1.5 animate-spin" />
            <span>Delete</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
