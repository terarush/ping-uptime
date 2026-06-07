<script setup lang="ts">
import { onMounted, computed, ref } from 'vue';
import { useAuth } from '@/composables/useAuth';
import { useMonitors, type Monitor } from '@/composables/useMonitors';
import { useIncidents } from '@/composables/useIncidents';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Activity,
  ShieldAlert,
  TrendingUp,
  Plus,
  RefreshCw,
  ExternalLink,
  CheckCircle2,
  AlertCircle,
  Loader2,
  Clock,
  Wifi,
  WifiOff,
  PauseCircle,
  ArrowRight,
  Zap,
  Globe,
  List,
  Bell,
  BarChart3,
} from '@lucide/vue';
import gsap from 'gsap';
import { useRouter } from 'vue-router';
import { siteConfig } from '@/content/config';

const { currentUser } = useAuth();
const router = useRouter();

const { monitors, fetchMonitors, loading: monitorsLoading } = useMonitors();
const { incidents, fetchIncidentsData, loading: incidentsLoading } = useIncidents();

const loading = computed(() => monitorsLoading.value || incidentsLoading.value);

// Mobile tab switcher
const activeTab = ref<'monitors' | 'events'>('monitors');

const fetchDashboardData = async () => {
  try {
    await Promise.all([fetchMonitors(), fetchIncidentsData()]);
  } catch (err: any) {
    console.error('Failed to load dashboard data:', err);
  } finally {
    setTimeout(animateDashboard, 50);
  }
};

// --- Computed Metrics ---
const totalMonitors = computed(() => monitors.value.length);
const upCount = computed(() => monitors.value.filter(m => m.status === 'active' && m.uptime_status === 'up').length);
const downCount = computed(() => monitors.value.filter(m => m.status === 'active' && m.uptime_status === 'down').length);
const pausedCount = computed(() => monitors.value.filter(m => m.status === 'paused').length);

const activeIncidents = computed(() => {
  return incidents.value.filter(i => {
    if (i.status !== 'active') return false;
    const m = monitors.value.find(x => x.id === i.monitor_id);
    if (m && (m.uptime_status === 'up' || m.status === 'paused')) {
      return false;
    }
    return true;
  });
});
const recentLogs = computed(() => incidents.value.slice(0, 5));
const recentMonitors = computed(() => [...monitors.value].slice(0, 10));

const avgLatency = computed(() => {
  const withLatency = monitors.value.filter(m => m.last_latency && m.last_latency > 0);
  if (withLatency.length === 0) return null;
  const sum = withLatency.reduce((acc, m) => acc + (m.last_latency ?? 0), 0);
  return Math.round(sum / withLatency.length);
});

const overallHealthLabel = computed(() => {
  if (monitors.value.length === 0) return { text: 'No monitors', color: 'text-muted-foreground', dot: 'bg-slate-400' };
  if (downCount.value > 0) return { text: 'Outage detected', color: 'text-rose-500', dot: 'bg-rose-500 animate-pulse' };
  const active = monitors.value.filter(m => m.status === 'active');
  const unknowns = active.filter(m => m.uptime_status === 'unknown');
  if (unknowns.length === active.length) return { text: 'Initializing...', color: 'text-amber-500', dot: 'bg-amber-500' };
  return { text: 'All systems operational', color: 'text-emerald-500', dot: 'bg-emerald-500' };
});

// --- Helpers ---
const getLastCheckedText = (item: Monitor) => {
  if (item.status === 'paused') return 'Paused';
  if (!item.last_checked_at) return 'Pending';
  const d = new Date(item.last_checked_at);
  const diffSec = Math.floor((Date.now() - d.getTime()) / 1000);
  if (diffSec < 5) return 'Just now';
  if (diffSec < 60) return `${diffSec}s ago`;
  const diffMin = Math.floor(diffSec / 60);
  return `${diffMin}m ago`;
};

const getMonitorName = (id: number) => {
  const m = monitors.value.find(x => x.id === id);
  return m ? m.name : `Monitor #${id}`;
};

const formatRelativeTime = (dateStr: string) => {
  const d = new Date(dateStr);
  const diffSec = Math.floor((Date.now() - d.getTime()) / 1000);
  if (diffSec < 60) return `${diffSec}s ago`;
  const diffMin = Math.floor(diffSec / 60);
  if (diffMin < 60) return `${diffMin}m ago`;
  const diffHr = Math.floor(diffMin / 60);
  if (diffHr < 24) return `${diffHr}h ago`;
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
};

const statusClass = (item: Monitor) => {
  if (item.status === 'paused') return 'bg-amber-500 ring-amber-500/25';
  if (item.uptime_status === 'up') return 'bg-emerald-500 ring-emerald-500/25';
  if (item.uptime_status === 'down') return 'bg-rose-500 ring-rose-500/25 animate-pulse';
  return 'bg-slate-400 ring-slate-400/15';
};

const goToAddMonitor = () => router.push(`${siteConfig.appPath}/monitors`);
const navTo = (path: string) => router.push(`${siteConfig.appPath}/${path}`);
const selectMonitorAndNavigate = (id: number) => router.push(`${siteConfig.appPath}/monitors?select=${id}`);

const animateDashboard = () => {
  gsap.fromTo('.stat-card', { opacity: 0, y: 14 }, { opacity: 1, y: 0, duration: 0.35, stagger: 0.06, ease: 'power2.out' });
  gsap.fromTo('.monitor-row', { opacity: 0, x: -10 }, { opacity: 1, x: 0, duration: 0.3, stagger: 0.035, ease: 'power2.out' });
  gsap.fromTo('.event-row', { opacity: 0, x: 10 }, { opacity: 1, x: 0, duration: 0.3, stagger: 0.035, ease: 'power2.out' });
};

onMounted(() => {
  fetchDashboardData();
});
</script>

<template>
  <div class="min-h-[calc(100vh-4rem)] p-3 sm:p-5 lg:p-8 space-y-4 sm:space-y-6 max-w-7xl mx-auto">

    <!-- ── Header ──────────────────────────────────────────── -->
    <div class="flex items-start justify-between gap-3">
      <div class="min-w-0">
        <h2 class="text-lg sm:text-2xl font-black tracking-tight text-foreground truncate">
          Welcome, <span class="text-primary">{{ currentUser?.name || 'User' }}</span>
        </h2>
        <!-- Overall health inline pill -->
        <div class="flex items-center gap-1.5 mt-1">
          <span class="w-2 h-2 rounded-full shrink-0" :class="overallHealthLabel.dot" />
          <span class="text-xs font-semibold" :class="overallHealthLabel.color">{{ overallHealthLabel.text }}</span>
        </div>
      </div>

      <div class="flex items-center gap-1.5 shrink-0">
        <Button variant="outline" size="sm" @click="fetchDashboardData" class="h-8 px-2.5 text-xs gap-1.5">
          <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" />
          <span class="hidden sm:inline">Refresh</span>
        </Button>
        <Button @click="goToAddMonitor" size="sm" class="h-8 px-2.5 text-xs gap-1.5">
          <Plus class="w-3.5 h-3.5" />
          <span class="hidden sm:inline">Add Monitor</span>
        </Button>
      </div>
    </div>

    <!-- ── Loading ─────────────────────────────────────────── -->
    <div v-if="loading && monitors.length === 0" class="flex flex-col items-center justify-center py-24 gap-3">
      <Loader2 class="w-8 h-8 text-primary animate-spin" />
      <p class="text-sm text-muted-foreground">Loading...</p>
    </div>

    <template v-else>
      <!-- ── Stat Cards ──────────────────────────────────────── -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 sm:gap-3 lg:gap-4">
        <!-- Monitors -->
        <Card class="stat-card border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md hover:shadow-md transition-shadow">
          <CardContent class="p-3 sm:p-4 lg:p-5">
            <div class="flex items-start justify-between gap-2">
              <div class="min-w-0">
                <p class="text-[9px] sm:text-[10px] font-bold text-muted-foreground uppercase tracking-wider">Monitors</p>
                <p class="text-xl sm:text-2xl lg:text-3xl font-black text-foreground mt-0.5">{{ totalMonitors }}</p>
                <p class="text-[9px] sm:text-[10px] text-muted-foreground mt-0.5">{{ pausedCount }} paused</p>
              </div>
              <div class="p-1.5 sm:p-2 rounded-lg bg-slate-100 dark:bg-slate-800/60 shrink-0">
                <Globe class="w-3.5 h-3.5 sm:w-4 sm:h-4 text-slate-500" />
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Online -->
        <Card class="stat-card border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md hover:shadow-md transition-shadow">
          <CardContent class="p-3 sm:p-4 lg:p-5">
            <div class="flex items-start justify-between gap-2">
              <div class="min-w-0">
                <p class="text-[9px] sm:text-[10px] font-bold text-muted-foreground uppercase tracking-wider">Online</p>
                <p class="text-xl sm:text-2xl lg:text-3xl font-black text-emerald-500 mt-0.5">{{ upCount }}</p>
                <p class="text-[9px] sm:text-[10px] text-emerald-600/70 dark:text-emerald-400/60 mt-0.5">Operational</p>
              </div>
              <div class="p-1.5 sm:p-2 rounded-lg bg-emerald-500/10 shrink-0">
                <Wifi class="w-3.5 h-3.5 sm:w-4 sm:h-4 text-emerald-500" />
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Offline -->
        <Card class="stat-card border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md hover:shadow-md transition-shadow">
          <CardContent class="p-3 sm:p-4 lg:p-5">
            <div class="flex items-start justify-between gap-2">
              <div class="min-w-0">
                <p class="text-[9px] sm:text-[10px] font-bold text-muted-foreground uppercase tracking-wider">Offline</p>
                <p
                  class="text-xl sm:text-2xl lg:text-3xl font-black mt-0.5"
                  :class="downCount > 0 ? 'text-rose-500' : 'text-foreground'"
                >{{ downCount }}</p>
                <p class="text-[9px] sm:text-[10px] mt-0.5" :class="downCount > 0 ? 'text-rose-500/80' : 'text-muted-foreground'">
                  {{ downCount > 0 ? 'Action needed' : 'All clear' }}
                </p>
              </div>
              <div class="p-1.5 sm:p-2 rounded-lg shrink-0" :class="downCount > 0 ? 'bg-rose-500/10' : 'bg-slate-100 dark:bg-slate-800/60'">
                <WifiOff class="w-3.5 h-3.5 sm:w-4 sm:h-4" :class="downCount > 0 ? 'text-rose-500' : 'text-slate-400'" />
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Avg Latency -->
        <Card class="stat-card border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md hover:shadow-md transition-shadow">
          <CardContent class="p-3 sm:p-4 lg:p-5">
            <div class="flex items-start justify-between gap-2">
              <div class="min-w-0">
                <p class="text-[9px] sm:text-[10px] font-bold text-muted-foreground uppercase tracking-wider">Avg Latency</p>
                <p class="text-xl sm:text-2xl lg:text-3xl font-black text-indigo-500 mt-0.5">
                  {{ avgLatency !== null ? avgLatency : '—' }}<span v-if="avgLatency !== null" class="text-sm sm:text-base font-bold text-indigo-400">ms</span>
                </p>
                <p class="text-[9px] sm:text-[10px] text-muted-foreground mt-0.5">Latest checks</p>
              </div>
              <div class="p-1.5 sm:p-2 rounded-lg bg-indigo-500/10 shrink-0">
                <Zap class="w-3.5 h-3.5 sm:w-4 sm:h-4 text-indigo-500" />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- ── Mobile Tab switcher ────────────────────────────── -->
      <div class="flex rounded-lg border border-border/50 bg-muted/30 p-0.5 gap-0.5 lg:hidden">
        <button
          @click="activeTab = 'monitors'"
          class="flex-1 flex items-center justify-center gap-1.5 h-8 rounded-md text-xs font-bold transition-all"
          :class="activeTab === 'monitors' ? 'bg-background text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'"
        >
          <List class="w-3.5 h-3.5" />
          Monitors
          <span v-if="downCount > 0" class="w-1.5 h-1.5 rounded-full bg-rose-500 animate-pulse shrink-0" />
        </button>
        <button
          @click="activeTab = 'events'"
          class="flex-1 flex items-center justify-center gap-1.5 h-8 rounded-md text-xs font-bold transition-all"
          :class="activeTab === 'events' ? 'bg-background text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'"
        >
          <Bell class="w-3.5 h-3.5" />
          Events
          <span v-if="activeIncidents.length > 0" class="bg-rose-500 text-white text-[9px] font-bold rounded-full px-1 leading-tight shrink-0">
            {{ activeIncidents.length }}
          </span>
        </button>
      </div>

      <!-- ── Main 2-col grid ────────────────────────────────── -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 lg:gap-6 items-start">

        <!-- Monitor list (hidden on mobile when events tab is active) -->
        <div class="lg:col-span-2" :class="activeTab === 'events' ? 'hidden lg:block' : 'block'">
          <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
            <CardHeader class="p-4 sm:p-5 border-b border-border/40">
              <div class="flex items-center justify-between gap-2">
                <div>
                  <CardTitle class="text-sm font-bold">Endpoint Status</CardTitle>
                  <CardDescription class="text-[10px] sm:text-xs mt-0.5">Live uptime status and latency per endpoint.</CardDescription>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  class="h-7 sm:h-8 text-[10px] sm:text-xs gap-1 text-muted-foreground shrink-0"
                  @click="navTo('monitors')"
                >
                  All <ArrowRight class="w-3 h-3" />
                </Button>
              </div>
            </CardHeader>

            <CardContent class="p-0">
              <!-- Empty -->
              <div v-if="monitors.length === 0" class="flex flex-col items-center justify-center py-14 text-center px-6">
                <div class="w-12 h-12 rounded-2xl bg-muted flex items-center justify-center mb-3">
                  <Activity class="w-6 h-6 text-muted-foreground/40" />
                </div>
                <p class="text-sm font-bold">No monitors yet</p>
                <p class="text-xs text-muted-foreground mt-1">Add your first endpoint to start monitoring.</p>
                <Button @click="goToAddMonitor" size="sm" class="mt-4 h-8 text-xs gap-1.5">
                  <Plus class="w-3.5 h-3.5" /> Add Monitor
                </Button>
              </div>

              <!-- List -->
              <div v-else class="divide-y divide-border/30">
                <div
                  v-for="item in recentMonitors"
                  :key="item.id"
                  @click="selectMonitorAndNavigate(item.id)"
                  class="monitor-row flex items-center gap-3 px-4 py-3 sm:py-3.5 hover:bg-muted/40 active:bg-muted/60 transition-colors cursor-pointer group"
                >
                  <!-- Status dot -->
                  <span class="h-2.5 w-2.5 rounded-full ring-4 shrink-0" :class="statusClass(item)" />

                  <!-- Name + URL -->
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-1.5 flex-wrap">
                      <span class="text-xs sm:text-sm font-bold text-foreground truncate">{{ item.name }}</span>
                      <Badge v-if="item.status === 'paused'" class="text-[8px] font-bold h-3.5 px-1 bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20 shrink-0">
                        Paused
                      </Badge>
                    </div>
                    <a
                      :href="item.url"
                      target="_blank"
                      @click.stop
                      class="text-[10px] text-muted-foreground hover:text-primary inline-flex items-center gap-0.5 group/link"
                    >
                      <span class="truncate max-w-45 sm:max-w-xs">{{ item.url }}</span>
                      <ExternalLink class="w-2.5 h-2.5 shrink-0 opacity-0 group-hover/link:opacity-100 transition-opacity" />
                    </a>
                  </div>

                  <!-- Right side: latency + type badge + last checked -->
                  <div class="shrink-0 flex flex-col items-end gap-0.5 text-right">
                    <span
                      class="text-xs font-bold tabular-nums"
                      :class="item.last_latency && item.last_latency > 0 ? 'text-foreground' : 'text-muted-foreground/40'"
                    >
                      {{ item.last_latency && item.last_latency > 0 ? `${item.last_latency}ms` : '—' }}
                    </span>
                    <div class="flex items-center gap-1">
                      <Badge variant="outline" class="text-[8px] font-bold h-3.5 px-1 hidden sm:flex border-border/40 text-muted-foreground/60">
                        {{ item.interval }}s
                      </Badge>
                      <span class="text-[9px] sm:text-[10px] text-muted-foreground">{{ getLastCheckedText(item) }}</span>
                    </div>
                  </div>
                </div>

                <!-- Footer: see all -->
                <div v-if="monitors.length > 10" class="flex items-center justify-center py-2.5">
                  <Button variant="ghost" size="sm" class="h-7 text-xs text-muted-foreground gap-1" @click="navTo('monitors')">
                    +{{ monitors.length - 10 }} more monitors <ArrowRight class="w-3 h-3" />
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Right column: events + quick nav (hidden on mobile when monitors tab is active) -->
        <div class="space-y-3 sm:space-y-4" :class="activeTab === 'monitors' ? 'hidden lg:flex lg:flex-col' : 'flex flex-col'">

          <!-- Active Incident Alert -->
          <div
            v-if="activeIncidents.length > 0"
            class="rounded-xl border border-rose-500/30 bg-rose-500/5 p-3 sm:p-4 space-y-2.5"
          >
            <div class="flex items-center gap-2">
              <AlertCircle class="w-4 h-4 text-rose-500 animate-pulse shrink-0" />
              <p class="text-xs font-bold text-rose-600 dark:text-rose-400">
                {{ activeIncidents.length }} Active Incident{{ activeIncidents.length > 1 ? 's' : '' }}
              </p>
            </div>
            <div class="space-y-1.5">
              <div
                v-for="inc in activeIncidents.slice(0, 4)"
                :key="inc.id"
                class="flex items-center justify-between gap-2 text-[10px] sm:text-xs"
              >
                <div class="flex items-center gap-1.5 min-w-0">
                  <span class="w-1.5 h-1.5 rounded-full bg-rose-500 shrink-0" />
                  <span class="font-semibold text-rose-700 dark:text-rose-300 truncate">{{ getMonitorName(inc.monitor_id) }}</span>
                </div>
                <span class="text-rose-500/70 shrink-0">{{ inc.latency > 0 ? `${inc.latency}ms` : '' }}</span>
              </div>
            </div>
            <Button size="sm" variant="outline" class="w-full h-7 text-[10px] border-rose-500/30 text-rose-600 dark:text-rose-400 hover:bg-rose-500/10" @click="navTo('incidents')">
              View all incidents
            </Button>
          </div>

          <!-- Events Feed -->
          <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
            <CardHeader class="p-4 sm:p-5 border-b border-border/40">
              <div class="flex items-center justify-between gap-2">
                <div>
                  <CardTitle class="text-sm font-bold">Recent Events</CardTitle>
                  <CardDescription class="text-[10px] sm:text-xs mt-0.5">Incident & recovery log.</CardDescription>
                </div>
                <Button variant="ghost" size="sm" class="h-7 text-[10px] gap-1 text-muted-foreground shrink-0" @click="navTo('incidents')">
                  All <ArrowRight class="w-3 h-3" />
                </Button>
              </div>
            </CardHeader>
            <CardContent class="p-0 max-h-[260px] overflow-y-auto">
              <!-- Empty -->
              <div v-if="incidents.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
                <CheckCircle2 class="w-8 h-8 text-emerald-500/25 mb-2" />
                <p class="text-xs font-bold">All systems nominal</p>
                <p class="text-[10px] text-muted-foreground mt-0.5">No incidents recorded.</p>
              </div>

              <div v-else class="divide-y divide-border/30">
                <div
                  v-for="log in recentLogs"
                  :key="log.id"
                  class="event-row flex items-center justify-between gap-3 px-4 py-2.5 hover:bg-muted/30 transition-colors"
                >
                  <div class="flex items-center gap-2.5 min-w-0">
                    <!-- Icon -->
                    <div class="shrink-0">
                      <div v-if="log.status === 'resolved'" class="w-4.5 h-4.5 rounded-full bg-emerald-500/15 flex items-center justify-center">
                        <CheckCircle2 class="w-2.5 h-2.5 text-emerald-500" />
                      </div>
                      <div v-else class="w-4.5 h-4.5 rounded-full bg-rose-500/15 flex items-center justify-center">
                        <AlertCircle class="w-2.5 h-2.5 text-rose-500" />
                      </div>
                    </div>

                    <!-- Name + Status -->
                    <div class="min-w-0">
                      <p class="text-xs font-bold text-foreground truncate max-w-[130px] sm:max-w-[160px]">
                        {{ getMonitorName(log.monitor_id) }}
                      </p>
                      <p class="text-[9px] font-semibold text-muted-foreground">
                        <span :class="log.status === 'resolved' ? 'text-emerald-600 dark:text-emerald-400' : 'text-rose-500'">
                          {{ log.status === 'resolved' ? 'Recovered' : 'Down' }}
                        </span>
                        <span class="mx-1">·</span>
                        <span>{{ formatRelativeTime(log.created_at) }}</span>
                      </p>
                    </div>
                  </div>

                  <!-- Latency -->
                  <div class="shrink-0 text-right flex flex-col items-end gap-0.5">
                    <span v-if="log.latency > 0" class="text-[10px] font-mono text-muted-foreground font-bold">{{ log.latency }}ms</span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- Quick navigation -->
          <div class="grid grid-cols-3 gap-2">
            <button
              @click="navTo('monitors')"
              class="flex flex-col items-center gap-1.5 py-3 px-2 rounded-xl border border-border/50 bg-card/60 dark:bg-card/40 hover:bg-muted/60 hover:border-border transition-all text-muted-foreground hover:text-foreground"
            >
              <Activity class="w-4 h-4" />
              <span class="text-[9px] sm:text-[10px] font-bold">Monitors</span>
            </button>
            <button
              @click="navTo('analytics')"
              class="flex flex-col items-center gap-1.5 py-3 px-2 rounded-xl border border-border/50 bg-card/60 dark:bg-card/40 hover:bg-muted/60 hover:border-border transition-all text-muted-foreground hover:text-foreground"
            >
              <BarChart3 class="w-4 h-4" />
              <span class="text-[9px] sm:text-[10px] font-bold">Analytics</span>
            </button>
            <button
              @click="navTo('incidents')"
              class="flex flex-col items-center gap-1.5 py-3 px-2 rounded-xl border border-border/50 bg-card/60 dark:bg-card/40 hover:bg-muted/60 hover:border-border transition-all text-muted-foreground hover:text-foreground"
            >
              <AlertCircle class="w-4 h-4" />
              <span class="text-[9px] sm:text-[10px] font-bold">Incidents</span>
            </button>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
