<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useAuth } from '@/composables/useAuth';
import { useMonitors } from '@/composables/useMonitors';
import { useAnalytics } from '@/composables/useAnalytics';
import { useUsersStore } from '@/stores/users';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import {
  BarChart3,
  Search,
  Loader2,
  RefreshCw,
  TrendingUp,
  Activity,
  CheckCircle2,
  AlertCircle,
  Clock,
  Check,
  X
} from '@lucide/vue';
import gsap from 'gsap';
import MultiChart, { type Series } from '@/components/multi-chart.vue';
import AnalyticsTable from '@/components/analytics-table.vue';

const { currentUser } = useAuth();
const isAdmin = computed(() => currentUser.value?.role === 'admin');

const usersStore = useUsersStore();
const {
  stats,
  loading: analyticsLoading,
  fetchDashboardStats
} = useAnalytics();

const {
  monitors,
  loading: monitorsLoading,
  fetchMonitors
} = useMonitors();

const loading = computed(() => analyticsLoading.value || monitorsLoading.value || (isAdmin.value && usersStore.loading));

const searchQuery = ref('');
const selectedWindow = ref<'1h' | '1d' | '1w' | '1m' | '1y' | 'all'>('1d'); // Default to 1 day
const enabledMonitors = ref<Set<number>>(new Set());

// Colors palette for chart lines
const colorsPalette = [
  '#6366f1', // Indigo
  '#10b981', // Emerald
  '#3b82f6', // Blue
  '#f59e0b', // Amber
  '#ec4899', // Pink
  '#8b5cf6', // Violet
  '#06b6d4', // Cyan
  '#f97316', // Orange
  '#14b8a6', // Teal
  '#ef4444', // Red
];

const getMonitorColor = (monitorId: number, index: number): string => {
  return colorsPalette[index % colorsPalette.length] || '#6366f1';
};

// Fetch all necessary data
const fetchAllData = async () => {
  try {
    await Promise.all([
      fetchMonitors(),
      fetchDashboardStats(selectedWindow.value),
      isAdmin.value ? usersStore.fetchUsersData() : Promise.resolve()
    ]);

    // Default to enable all monitors on first load
    if (stats.value && enabledMonitors.value.size === 0) {
      stats.value.forEach(s => {
        enabledMonitors.value.add(s.monitor_id);
      });
    }
  } catch (err) {
    console.error('Failed to load analytics data:', err);
  } finally {
    setTimeout(animateItems, 50);
  }
};

// Toggle single monitor visibility
const toggleMonitor = (monitorId: number) => {
  if (enabledMonitors.value.has(monitorId)) {
    enabledMonitors.value.delete(monitorId);
  } else {
    enabledMonitors.value.add(monitorId);
  }
};

// Select All / Deselect All helpers
const selectAllMonitors = () => {
  stats.value.forEach(s => {
    enabledMonitors.value.add(s.monitor_id);
  });
};

const deselectAllMonitors = () => {
  enabledMonitors.value.clear();
};

// Filtered stats list based on search query
const filteredStats = computed(() => {
  if (!stats.value) return [];
  const query = searchQuery.value.toLowerCase().trim();
  if (!query) return stats.value;
  return stats.value.filter(s =>
    s.monitor_name.toLowerCase().includes(query) ||
    s.monitor_url.toLowerCase().includes(query)
  );
});

// Construct the series list for MultiChart
const chartSeries = computed<Series[]>(() => {
  if (!stats.value) return [];

  return stats.value
    .filter(s => enabledMonitors.value.has(s.monitor_id))
    .map((s, idx) => ({
      monitorId: s.monitor_id,
      name: s.monitor_name,
      color: getMonitorColor(s.monitor_id, idx),
      points: s.points || [],
    }));
});

// Aggregate calculations for visible (searched/filtered) monitors
const totalEndpoints = computed(() => stats.value.length);

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

// User mapping to resolve monitor owner
const getMonitorOwner = (monitorId: number) => {
  const monitor = monitors.value.find(m => m.id === monitorId);
  if (!monitor) return null;
  const user = usersStore.users.find(u => u.id === monitor.user_id);
  return user ? `${user.name} (${user.email})` : `User ID: ${monitor.user_id}`;
};

// Handle window change
watch(selectedWindow, async (newVal) => {
  try {
    await fetchDashboardStats(newVal);
  } catch (err) {
    console.error('Failed to update stats window:', err);
  }
});

// GSAP animations
const animateItems = () => {
  gsap.fromTo('.stat-card',
    { opacity: 0, y: 10 },
    { opacity: 1, y: 0, duration: 0.3, stagger: 0.05, ease: 'power2.out' }
  );
  gsap.fromTo('.monitor-checklist-item',
    { opacity: 0, x: -10 },
    { opacity: 1, x: 0, duration: 0.2, stagger: 0.02, ease: 'power2.out' }
  );
};

onMounted(async () => {
  await fetchAllData();

  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2, ease: 'power3.out' }
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
          <BarChart3 class="w-6 h-6 text-primary" />
          <span>System Analytics</span>
        </h2>
        <p class="text-xs text-muted-foreground">
          {{ isAdmin ? 'Comparative overview and performance timeline of all system monitors.' : 'Uptime history and response times for your configured endpoints.' }}
        </p>
      </div>

      <!-- Controls -->
      <div class="flex items-center gap-2">
        <Select v-model="selectedWindow">
          <SelectTrigger class="h-9 w-37.5">
            <SelectValue placeholder="Select Window" />
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

        <Button variant="outline" size="sm" @click="fetchAllData" class="h-9">
          <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
          <span>Refresh</span>
        </Button>
      </div>
    </div>

    <!-- Aggregate Overview Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 z-10 relative">
      <!-- Total Endpoints -->
      <Card class="stat-card border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
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
      <Card class="stat-card border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
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
      <Card class="stat-card border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
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
      <Card class="stat-card border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
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

    <!-- Main Chart Section: Multi-Monitor Interactive Canvas & Checklist -->
    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <CardTitle class="text-sm font-bold text-foreground flex items-center gap-2">
              <TrendingUp class="w-4 h-4 text-primary" />
              <span>Comparative Latency Timeline</span>
            </CardTitle>
            <CardDescription class="text-xs">Compare latency metrics across multiple selected monitors concurrently.</CardDescription>
          </div>
        </div>
      </CardHeader>

      <CardContent class="p-6">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
          <!-- Left: The Interactive Multi-Line Chart -->
          <div class="lg:col-span-8 space-y-4">
            <div class="relative bg-muted/20 border border-border/30 rounded-xl p-4 min-h-75 flex items-center justify-center">
              <div v-if="loading && stats.length === 0" class="flex flex-col items-center justify-center gap-2">
                <Loader2 class="w-6 h-6 text-primary animate-spin" />
                <span class="text-xs text-muted-foreground">Loading chart data...</span>
              </div>
              <div v-else-if="chartSeries.length === 0" class="flex flex-col items-center justify-center text-center text-muted-foreground p-8">
                <X class="w-8 h-8 mb-2 text-muted-foreground/40" />
                <p class="text-sm font-bold">No monitors selected</p>
                <p class="text-xs mt-1">Select one or more monitors from the checklist to draw comparative lines.</p>
              </div>
              <MultiChart v-else :seriesList="chartSeries" :height="320" />
            </div>

            <!-- Chart Legends -->
            <div v-if="chartSeries.length > 0" class="flex flex-wrap gap-x-4 gap-y-2 text-[11px] font-medium pt-2 border-t border-border/30">
              <div
                v-for="s in chartSeries"
                :key="s.monitorId"
                class="flex items-center gap-1.5"
              >
                <span class="w-2.5 h-2.5 rounded-full shrink-0" :style="{ backgroundColor: s.color }"></span>
                <span class="text-foreground truncate max-w-30">{{ s.name }}</span>
              </div>
            </div>
          </div>

          <!-- Right: Manage / Toggle Monitors list -->
          <div class="lg:col-span-4 border border-border/30 rounded-xl p-4 bg-muted/10 space-y-4">
            <div class="flex items-center justify-between border-b border-border/40 pb-2">
              <h4 class="text-xs font-black uppercase tracking-wider text-foreground">Manage Monitors</h4>
              <div class="flex items-center gap-2">
                <button
                  @click="selectAllMonitors"
                  class="text-[10px] text-primary hover:underline font-bold cursor-pointer"
                >
                  Select All
                </button>
                <span class="text-muted-foreground/30 text-[10px]">•</span>
                <button
                  @click="deselectAllMonitors"
                  class="text-[10px] text-rose-500 hover:underline font-bold cursor-pointer"
                >
                  Clear
                </button>
              </div>
            </div>

            <!-- Search input within checklists -->
            <div class="relative w-full">
              <Search class="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-muted-foreground" />
              <Input
                v-model="searchQuery"
                placeholder="Search checklist..."
                class="pl-8 h-8 text-xs bg-card"
              />
            </div>

            <!-- Scrollable checklist -->
            <div class="space-y-1 max-h-60.5 overflow-y-auto pr-1">
              <div v-if="loading && stats.length === 0" class="flex items-center justify-center py-10">
                <Loader2 class="w-5 h-5 text-primary animate-spin" />
              </div>
              <div v-else-if="filteredStats.length === 0" class="text-center text-xs text-muted-foreground py-10">
                No matching monitors
              </div>
              <label
                v-else
                v-for="(item, idx) in filteredStats"
                :key="item.monitor_id"
                class="monitor-checklist-item flex items-center justify-between p-2 rounded-lg hover:bg-muted/50 border border-transparent hover:border-border/30 cursor-pointer select-none transition-colors"
              >
                <div class="flex items-center gap-2 min-w-0">
                  <input
                    type="checkbox"
                    :checked="enabledMonitors.has(item.monitor_id)"
                    @change="toggleMonitor(item.monitor_id)"
                    class="rounded border-input text-primary focus:ring-primary h-3.5 w-3.5 shrink-0 cursor-pointer"
                  />
                  <!-- Color Indicator Dot -->
                  <span
                    class="w-2.5 h-2.5 rounded-full shrink-0"
                    :style="{ backgroundColor: getMonitorColor(item.monitor_id, idx) }"
                  ></span>
                  <span class="text-xs font-bold text-foreground truncate max-w-37.5">{{ item.monitor_name }}</span>
                </div>
                <span class="text-[10px] text-muted-foreground font-semibold shrink-0">{{ Math.round(item.avg_latency) }} ms</span>
              </label>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Bottom Section: Detailed Performance Table for all monitors -->
    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <div>
          <CardTitle class="text-sm font-bold text-foreground">Endpoint Analytics Breakdown</CardTitle>
          <CardDescription class="text-xs">Uptime percentage, failed checks, and latency ranges for active monitors.</CardDescription>
        </div>
      </CardHeader>

      <CardContent class="p-0">
        <div v-if="loading && stats.length === 0" class="flex items-center justify-center py-12">
          <Loader2 class="w-6 h-6 text-primary animate-spin" />
        </div>
        <div v-else-if="filteredStats.length === 0" class="text-center text-xs text-muted-foreground py-12">
          No monitors found.
        </div>
        <AnalyticsTable
          v-else
          :stats="filteredStats"
          :is-admin="isAdmin"
          :get-monitor-owner="getMonitorOwner"
        />
      </CardContent>
    </Card>
  </div>
</template>

<style scoped>
.monitor-checklist-item {
  transition: transform 0.15s cubic-bezier(0.16, 1, 0.3, 1), background-color 0.15s;
}
.monitor-checklist-item:hover {
  transform: translateX(1px);
}
</style>
