<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useRoute } from 'vue-router';
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
  BarChart3,
  CheckCircle2,
  Clock,
  Shield,
  X
} from '@lucide/vue';
import gsap from 'gsap';
import { useAnalytics } from '@/composables/useAnalytics';
import MultiChart from '@/components/multi-chart.vue';

// Auth checks
const { currentUser } = useAuth();
const isAdmin = computed(() => currentUser.value?.role === 'admin');
const route = useRoute();

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
const selectedWindow = ref<'1h' | '1d' | '1w' | '1m' | '1y' | 'all'>('1d');

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
const formCheckSsl = ref(false);
const formStatus = ref('active');

const isEditMode = computed(() => !!actionMonitor.value);

// Fetch wrapper with animation callback
const fetchAll = async () => {
  try {
    await Promise.all([
      fetchMonitors(),
      fetchDashboardStats(selectedWindow.value)
    ]);

    // Auto-select first monitor if none is selected and monitors exist
    const firstMonitor = monitors.value?.[0];
    if (!selectedMonitorId.value && firstMonitor) {
      selectedMonitorId.value = firstMonitor.id;
    }
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

const selectedMonitorStats = computed(() => {
  if (!selectedMonitorId.value) return null;
  return stats.value.find(s => s.monitor_id === selectedMonitorId.value) || null;
});

// Aggregate overview stats
const totalEndpoints = computed(() => stats.value.length || monitors.value.length);

const avgUptime = computed(() => {
  if (stats.value.length === 0) return '0.00%';
  const sum = stats.value.reduce((acc, s) => acc + s.uptime_pct, 0);
  return `${(sum / stats.value.length).toFixed(2)}%`;
});

const avgLatency = computed(() => {
  if (stats.value.length === 0) return '0 ms';
  const sum = stats.value.reduce((acc, s) => acc + s.avg_latency, 0);
  return `${Math.round(sum / stats.value.length)} ms`;
});

const statusDistribution = computed(() => {
  const dist = { operational: 0, degraded: 0, outage: 0 };
  stats.value.forEach(s => {
    if (s.status === 'operational') dist.operational++;
    else if (s.status === 'degraded') dist.degraded++;
    else if (s.status === 'outage') dist.outage++;
  });
  return dist;
});

const chartSeries = computed(() => {
  if (!selectedMonitorId.value || !chartPoints.value.length) return [];
  const monitorName = monitors.value.find(m => m.id === selectedMonitorId.value)?.name || 'Monitor';
  return [
    {
      monitorId: selectedMonitorId.value,
      name: monitorName,
      color: '#6366f1',
      points: chartPoints.value,
    }
  ];
});

// Reset form fields
const resetForm = () => {
  formName.value = '';
  formUrl.value = '';
  formType.value = 'http';
  formInterval.value = 60;
  formTimeout.value = 10;
  formCheckSsl.value = false;
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
  formCheckSsl.value = monitor.check_ssl ?? false;
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
    check_ssl: formCheckSsl.value,
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
    if (selectedMonitorId.value === actionMonitor.value.id) {
      selectedMonitorId.value = null;
    }
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

// Watcher for window and monitor selection changes
watch([selectedMonitorId, selectedWindow], async ([newId, newWindow]) => {
  try {
    const tasks: Promise<any>[] = [fetchDashboardStats(newWindow)];
    if (newId) {
      tasks.push(fetchChart(newId, newWindow));
    }
    await Promise.all(tasks);
  } catch (err) {
    console.error('Failed to load analytics data:', err);
  }
});

onMounted(async () => {
  // Check if a monitor ID is passed in the query params to select it automatically
  if (route.query.select) {
    const id = Number(route.query.select);
    if (!isNaN(id)) {
      selectedMonitorId.value = id;
    }
  }

  await fetchAll();

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

    <!-- Aggregate Overview Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 z-10 relative">
      <!-- Total Endpoints -->
      <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
        <CardContent class="p-5 flex items-center justify-between">
          <div class="space-y-1">
            <p class="text-xs font-medium text-muted-foreground">Total Endpoints</p>
            <p class="text-2xl font-bold tracking-tight">{{ totalEndpoints }}</p>
          </div>
          <div class="p-3 bg-muted rounded-xl">
            <Activity class="w-5 h-5 text-muted-foreground" />
          </div>
        </CardContent>
      </Card>

      <!-- Avg Uptime -->
      <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
        <CardContent class="p-5 flex items-center justify-between">
          <div class="space-y-1">
            <p class="text-xs font-medium text-muted-foreground">Average Uptime</p>
            <p class="text-2xl font-bold tracking-tight text-emerald-500">{{ avgUptime }}</p>
          </div>
          <div class="p-3 bg-emerald-500/10 rounded-xl">
            <CheckCircle2 class="w-5 h-5 text-emerald-500" />
          </div>
        </CardContent>
      </Card>

      <!-- Avg Latency -->
      <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
        <CardContent class="p-5 flex items-center justify-between">
          <div class="space-y-1">
            <p class="text-xs font-medium text-muted-foreground">Average Latency</p>
            <p class="text-2xl font-bold tracking-tight text-blue-500">{{ avgLatency }}</p>
          </div>
          <div class="p-3 bg-blue-500/10 rounded-xl">
            <Clock class="w-5 h-5 text-blue-500" />
          </div>
        </CardContent>
      </Card>

      <!-- Status Breakdown -->
      <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
        <CardContent class="p-5 flex items-center justify-between">
          <div class="space-y-1.5">
            <p class="text-xs font-medium text-muted-foreground">Status Distribution</p>
            <div class="flex items-center gap-1.5 flex-wrap">
              <span class="inline-flex items-center gap-1 text-[11px] font-bold px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-500 border border-emerald-500/20">
                {{ statusDistribution.operational }}
                <span class="text-[9px] font-medium text-emerald-500/70 uppercase tracking-wider">Up</span>
              </span>
              <span class="inline-flex items-center gap-1 text-[11px] font-bold px-2 py-0.5 rounded-full bg-amber-500/10 text-amber-500 border border-amber-500/20">
                {{ statusDistribution.degraded }}
                <span class="text-[9px] font-medium text-amber-500/70 uppercase tracking-wider">Degraded</span>
              </span>
              <span class="inline-flex items-center gap-1 text-[11px] font-bold px-2 py-0.5 rounded-full bg-rose-500/10 text-rose-500 border border-rose-500/20" :class="{ 'animate-pulse': statusDistribution.outage > 0 }">
                {{ statusDistribution.outage }}
                <span class="text-[9px] font-medium text-rose-500/70 uppercase tracking-wider">Down</span>
              </span>
            </div>
          </div>
          <div class="p-3 bg-muted rounded-xl shrink-0">
            <BarChart3 class="w-5 h-5 text-muted-foreground" />
          </div>
        </CardContent>
      </Card>
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
          :selected-id="selectedMonitorId"
          @select="selectedMonitorId = $event"
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
              <SelectTrigger class="h-9 w-35">
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
          <div class="relative bg-muted/20 border border-border/30 rounded-xl p-4 min-h-55 flex items-center justify-center">
            <div v-if="analyticsLoading && chartPoints.length === 0" class="flex flex-col items-center justify-center gap-2">
              <Loader2 class="w-6 h-6 text-primary animate-spin" />
              <span class="text-xs text-muted-foreground">Loading chart data...</span>
            </div>
            <div v-else-if="chartSeries.length === 0" class="flex flex-col items-center justify-center text-center text-muted-foreground p-8">
              <X class="w-8 h-8 mb-2 text-muted-foreground/40" />
              <p class="text-sm font-bold">No data available</p>
            </div>
            <MultiChart v-else :seriesList="chartSeries" :height="220" />
          </div>

          <div class="grid grid-cols-2 md:grid-cols-5 gap-4">
            <!-- Uptime Card -->
            <div class="rounded-lg border border-border/60 bg-background p-3 hover:border-primary/40 transition-colors">
              <p class="text-[10px] font-bold text-muted-foreground uppercase tracking-wider">Uptime</p>
              <p class="text-lg font-black text-foreground mt-1">{{ (selectedMonitorStats?.uptime_pct ?? 0).toFixed(2) }}%</p>
              <p class="text-[9px] text-muted-foreground mt-0.5">{{ selectedMonitorStats?.failed_checks ?? 0 }} down events</p>
            </div>

            <!-- Avg Latency Card -->
            <div class="rounded-lg border border-border/60 bg-background p-3 hover:border-primary/40 transition-colors">
              <p class="text-[10px] font-bold text-muted-foreground uppercase tracking-wider">Avg Latency</p>
              <p class="text-lg font-black text-foreground mt-1">{{ (selectedMonitorStats?.avg_latency ?? 0).toFixed(1) }}ms</p>
              <p class="text-[9px] text-muted-foreground mt-0.5">Average response speed</p>
            </div>

            <!-- Min / Max Latency Card -->
            <div class="rounded-lg border border-border/60 bg-background p-3 hover:border-primary/40 transition-colors">
              <p class="text-[10px] font-bold text-muted-foreground uppercase tracking-wider">Min / Max Speed</p>
              <p class="text-lg font-black text-foreground mt-1">
                {{ Math.round(selectedMonitorStats?.min_latency ?? 0) }} / {{ Math.round(selectedMonitorStats?.max_latency ?? 0) }}<span class="text-xs font-normal text-muted-foreground ml-0.5">ms</span>
              </p>
              <p class="text-[9px] text-muted-foreground mt-0.5">Fastest & slowest pings</p>
            </div>

            <!-- Total Checks Card -->
            <div class="rounded-lg border border-border/60 bg-background p-3 hover:border-primary/40 transition-colors">
              <p class="text-[10px] font-bold text-muted-foreground uppercase tracking-wider">Total Checks</p>
              <p class="text-lg font-black text-foreground mt-1">{{ selectedMonitorStats?.total_checks ?? 0 }}</p>
              <p class="text-[9px] text-muted-foreground mt-0.5">Pings executed in window</p>
            </div>

            <!-- Status Card -->
            <div class="rounded-lg border border-border/60 bg-background p-3 hover:border-primary/40 transition-colors">
              <p class="text-[10px] font-bold text-muted-foreground uppercase tracking-wider">Status</p>
              <div class="flex items-center gap-2 mt-1.5">
                <span :class="[
                  'h-2.5 w-2.5 rounded-full ring-4 shrink-0',
                  selectedMonitorStats?.status === 'operational'
                    ? 'bg-emerald-500 ring-emerald-500/20'
                    : selectedMonitorStats?.status === 'degraded'
                    ? 'bg-amber-500 ring-amber-500/20'
                    : selectedMonitorStats?.status === 'outage'
                    ? 'bg-red-500 ring-red-500/20 animate-pulse'
                    : 'bg-slate-400 ring-slate-400/20'
                ]"></span>
                <span class="text-xs font-black capitalize text-foreground">{{ selectedMonitorStats?.status ?? '—' }}</span>
              </div>
              <p class="text-[9px] text-muted-foreground mt-1">Current system state</p>
            </div>
          </div>

          <div v-if="analyticsError" class="text-xs text-destructive">{{ analyticsError }}</div>
        </template>
      </CardContent>
    </Card>

    <!-- Create / Edit Dialog -->
    <Dialog v-model:open="isFormDialogOpen">
      <DialogContent class="sm:max-w-106.25">
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
                  <SelectItem value="heartbeat">Heartbeat</SelectItem>
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

          <div class="flex items-center gap-2 py-1.5">
            <input id="check_ssl" type="checkbox" v-model="formCheckSsl"
              class="h-4 w-4 rounded border-border text-primary focus-ring cursor-pointer" />
            <Label for="check_ssl" class="text-xs font-medium flex items-center gap-1.5 cursor-pointer">
              <Shield class="w-3.5 h-3.5 text-muted-foreground" />
              Check SSL certificate expiry
            </Label>
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
      <DialogContent class="sm:max-w-100">
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
