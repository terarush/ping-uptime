<script setup lang="ts">
import { onMounted, computed } from 'vue';
import { useAuth } from '@/composables/useAuth';
import { useMonitors, type Monitor } from '@/composables/useMonitors';
import { useIncidents, type Incident } from '@/composables/useIncidents';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import {
  Activity,
  ShieldCheck,
  ShieldAlert,
  Percent,
  Plus,
  RefreshCw,
  ExternalLink,
  CheckCircle2,
  AlertCircle,
  Loader2,
  Clock
} from '@lucide/vue';
import gsap from 'gsap';
import { useRouter } from 'vue-router';
import { siteConfig } from '@/content/config';

const { currentUser } = useAuth();
const router = useRouter();

// Use external global stores via composables
const { monitors, fetchMonitors, loading: monitorsLoading } = useMonitors();
const { incidents, fetchIncidentsData, loading: incidentsLoading } = useIncidents();

const loading = computed(() => monitorsLoading.value || incidentsLoading.value);

// Fetch dashboard data
const fetchDashboardData = async () => {
  try {
    await Promise.all([
      fetchMonitors(),
      fetchIncidentsData()
    ]);
  } catch (err: any) {
    console.error('Failed to load dashboard data:', err);
  } finally {
    setTimeout(animateDashboard, 50);
  }
};

// Calculations
const activeIncidentMonitorIDs = computed(() => {
  return incidents.value
    .filter(inc => inc.status === 'active')
    .map(inc => inc.monitor_id);
});

const totalMonitors = computed(() => monitors.value.length);

const upCount = computed(() => {
  return monitors.value.filter(m => m.status === 'active' && m.uptime_status === 'up').length;
});

const downCount = computed(() => {
  return monitors.value.filter(m => m.status === 'active' && m.uptime_status === 'down').length;
});

const pausedCount = computed(() => {
  return monitors.value.filter(m => m.status === 'paused').length;
});

const avgUptime = computed(() => {
  if (monitors.value.length === 0) return '100.00%';
  // Simulated overall calculation based on failures count
  const totalIncidents = incidents.value.length;
  if (totalIncidents === 0) return '100.00%';
  
  const uptimeFactor = Math.max(90, 100 - (totalIncidents * 0.15));
  return `${uptimeFactor.toFixed(2)}%`;
});

// Resolve latency response for a monitor
const getMonitorLatency = (monitorID: number) => {
  const monitorIncidents = incidents.value.filter(inc => inc.monitor_id === monitorID);
  if (monitorIncidents.length === 0) return '—';
  // Get latest incident latency
  const latest = monitorIncidents[0];
  if (!latest) return '—';
  return latest.latency > 0 ? `${latest.latency}ms` : '—';
};

// Check if monitor is currently operational
const isMonitorUp = (item: Monitor) => {
  return item.status === 'active' && item.uptime_status === 'up';
};

// Get last checked relative text
const getLastCheckedText = (item: Monitor) => {
  if (item.status === 'paused') return 'Paused';
  if (!item.last_checked_at) return 'Pending check';
  const d = new Date(item.last_checked_at);
  const diffSec = Math.floor((Date.now() - d.getTime()) / 1000);
  if (diffSec < 5) return 'Just now';
  if (diffSec < 60) return `${diffSec}s ago`;
  const diffMin = Math.floor(diffSec / 60);
  return `${diffMin}m ago`;
};

// Format incident trigger dates
const formatTime = (dateStr: string) => {
  const d = new Date(dateStr);
  return d.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false });
};

// Get recent 4 logs
const recentLogs = computed(() => {
  return incidents.value.slice(0, 4);
});

// Resolve monitor name for log
const getMonitorName = (id: number) => {
  const m = monitors.value.find(x => x.id === id);
  return m ? m.name : `Monitor #${id}`;
};

// Navigate to Add Monitor page
const goToAddMonitor = () => {
  router.push(`${siteConfig.appPath}/monitors`);
};

// GSAP Animations
const animateDashboard = () => {
  gsap.fromTo('.stat-card',
    { opacity: 0, y: 15 },
    { opacity: 1, y: 0, duration: 0.4, stagger: 0.05, ease: 'power2.out' }
  );
  gsap.fromTo('.dashboard-row',
    { opacity: 0, y: 15 },
    { opacity: 1, y: 0, duration: 0.5, stagger: 0.05, ease: 'power2.out' }
  );
};

onMounted(() => {
  fetchDashboardData();
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

    <!-- Top Greeting and Actions -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground">Welcome back, {{ currentUser?.name || 'User' }}!</h2>
        <p class="text-xs text-muted-foreground">Here is the real-time status of your self-hosted services.</p>
      </div>
      <div class="flex items-center gap-2">
        <Button variant="outline" size="sm" @click="fetchDashboardData" class="h-9">
          <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
          <span>Refresh</span>
        </Button>
        <Button @click="goToAddMonitor" size="sm" class="h-9 shadow-sm shadow-primary/10">
          <Plus class="w-4 h-4 mr-1.5" />
          <span>Add Monitor</span>
        </Button>
      </div>
    </div>

    <!-- Loading Placeholder -->
    <div v-if="loading && monitors.length === 0" class="flex flex-col items-center justify-center py-40 gap-4 z-10 relative">
      <Loader2 class="w-10 h-10 text-primary animate-spin" />
      <p class="text-sm text-muted-foreground font-semibold">Loading system statistics...</p>
    </div>

    <div v-else class="space-y-8 z-10 relative">
      <!-- Quick Stats Grid -->
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <!-- Total Monitors Card -->
        <Card class="stat-card hover:shadow-md transition-shadow border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
          <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Total Monitors</CardTitle>
            <div class="p-2 rounded-lg bg-slate-100 dark:bg-slate-900">
              <Activity class="w-4 h-4 text-slate-500" />
            </div>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-black text-foreground">{{ totalMonitors }}</div>
            <p class="text-[10px] text-muted-foreground mt-1">{{ pausedCount }} paused checks</p>
          </CardContent>
        </Card>

        <!-- Up Monitors Card -->
        <Card class="stat-card hover:shadow-md transition-shadow border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
          <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Active Services</CardTitle>
            <div class="p-2 rounded-lg bg-emerald-500/10 dark:bg-emerald-500/5">
              <ShieldCheck class="w-4 h-4 text-emerald-500" />
            </div>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-black text-emerald-500">{{ upCount }}</div>
            <p class="text-[10px] text-muted-foreground mt-1">Status: Operational</p>
          </CardContent>
        </Card>

        <!-- Down Monitors Card -->
        <Card class="stat-card hover:shadow-md transition-shadow border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
          <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Down Services</CardTitle>
            <div class="p-2 rounded-lg bg-red-500/10 dark:bg-red-500/5">
              <ShieldAlert class="w-4 h-4 text-red-500" />
            </div>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-black" :class="downCount > 0 ? 'text-red-500 animate-pulse' : 'text-foreground'">
              {{ downCount }}
            </div>
            <p class="text-[10px] text-muted-foreground mt-1">
              {{ downCount > 0 ? 'Immediate action required' : 'No outages reported' }}
            </p>
          </CardContent>
        </Card>

        <!-- Avg Uptime Card -->
        <Card class="stat-card hover:shadow-md transition-shadow border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
          <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Avg Uptime</CardTitle>
            <div class="p-2 rounded-lg bg-indigo-500/10 dark:bg-indigo-500/5">
              <Percent class="w-4 h-4 text-indigo-500" />
            </div>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-black text-foreground">{{ avgUptime }}</div>
            <p class="text-[10px] text-muted-foreground mt-1">Overall calculated metrics</p>
          </CardContent>
        </Card>
      </div>

      <!-- Main Grid Content: Monitors list & Incidents log -->
      <div class="grid gap-6 lg:grid-cols-3">
        <!-- Monitors List Card -->
        <Card class="lg:col-span-2 border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
          <CardHeader class="pb-3 border-b border-border/40">
            <CardTitle class="text-sm font-bold text-foreground">Operational Status</CardTitle>
            <CardDescription class="text-xs">Uptime statistics and latency responses for registered endpoints.</CardDescription>
          </CardHeader>
          <CardContent class="p-0">
            <div v-if="monitors.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
              <Activity class="w-10 h-10 text-muted-foreground/30 mb-3" />
              <p class="text-xs font-bold text-foreground">No monitors registered</p>
              <p class="text-[10px] text-muted-foreground mt-0.5">Click "Add Monitor" to configure your first URL checker.</p>
            </div>

            <Table v-else>
              <TableHeader>
                <TableRow>
                  <TableHead class="w-[80px] text-center text-xs font-bold uppercase text-muted-foreground">Status</TableHead>
                  <TableHead class="text-xs font-bold uppercase text-muted-foreground">Monitor</TableHead>
                  <TableHead class="text-xs font-bold uppercase text-muted-foreground">Latency</TableHead>
                  <TableHead class="text-xs font-bold uppercase text-muted-foreground">Interval</TableHead>
                  <TableHead class="text-xs font-bold uppercase text-muted-foreground">Last Checked</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="item in monitors" :key="item.id" class="dashboard-row hover:bg-muted/40 transition-colors">
                  <!-- Status indicator -->
                  <TableCell class="py-4">
                    <div class="flex items-center justify-center">
                      <span :class="[
                        'h-2.5 w-2.5 rounded-full ring-4 shrink-0',
                        item.status === 'paused'
                          ? 'bg-amber-500 ring-amber-500/20'
                          : item.uptime_status === 'up'
                          ? 'bg-emerald-500 ring-emerald-500/20'
                          : item.uptime_status === 'down'
                          ? 'bg-red-500 ring-red-500/20 animate-pulse'
                          : 'bg-slate-400 ring-slate-400/20'
                      ]"></span>
                    </div>
                  </TableCell>

                  <!-- Monitor detail -->
                  <TableCell class="font-medium text-foreground py-4">
                    <div class="flex flex-col gap-0.5 max-w-[200px] sm:max-w-xs md:max-w-sm lg:max-w-md">
                      <span class="text-xs font-bold truncate">{{ item.name }}</span>
                      <a :href="item.url" target="_blank" class="text-[10px] text-muted-foreground hover:underline inline-flex items-center gap-1">
                        <span class="truncate">{{ item.url }}</span>
                        <ExternalLink class="w-2.5 h-2.5 shrink-0" />
                      </a>
                    </div>
                  </TableCell>

                  <!-- Latency response -->
                  <TableCell class="py-4">
                    <span class="text-xs font-semibold text-foreground">
                      {{ getMonitorLatency(item.id) }}
                    </span>
                  </TableCell>

                  <!-- Interval -->
                  <TableCell class="py-4 text-xs text-muted-foreground">
                    <Badge variant="outline" class="font-semibold text-[9px] border-border/50 bg-slate-50 dark:bg-slate-900">
                      {{ item.interval }}s
                    </Badge>
                  </TableCell>

                  <!-- Last Checked -->
                  <TableCell class="text-muted-foreground text-xs py-4">
                    {{ getLastCheckedText(item) }}
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <!-- Incident Log Feed -->
        <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
          <CardHeader class="pb-3 border-b border-border/40">
            <CardTitle class="text-sm font-bold text-foreground">Recent Events</CardTitle>
            <CardDescription class="text-xs">Latest pings and recovery transitions logs.</CardDescription>
          </CardHeader>
          <CardContent class="pt-4 space-y-4">
            <div v-if="incidents.length === 0" class="flex flex-col items-center justify-center py-16 text-center">
              <CheckCircle2 class="w-10 h-10 text-emerald-500/20 mb-2" />
              <p class="text-[10px] font-bold text-foreground">No events recorded</p>
              <p class="text-[9px] text-muted-foreground mt-0.5">All monitored systems are functioning smoothly.</p>
            </div>

            <div v-else v-for="log in recentLogs" :key="log.id" class="dashboard-row flex items-start gap-3 text-xs leading-relaxed">
              <!-- Event Icon -->
              <div class="mt-0.5">
                <CheckCircle2 v-if="log.status === 'resolved'" class="w-4 h-4 text-emerald-500" />
                <AlertCircle v-else class="w-4 h-4 text-red-500 animate-pulse" />
              </div>

              <!-- Log Message detail -->
              <div class="flex-1 flex flex-col gap-0.5">
                <span class="font-bold text-foreground leading-tight">
                  {{ getMonitorName(log.monitor_id) }} is {{ log.status === 'resolved' ? 'resolved' : 'down' }}
                </span>
                <span class="text-[10px] text-muted-foreground line-clamp-2">
                  {{ log.error_message || 'Successful ping check' }}
                </span>
                <div class="flex items-center gap-2 text-[9px] text-muted-foreground/80 mt-0.5">
                  <span class="inline-flex items-center gap-0.5">
                    <Clock class="w-3 h-3" />
                    <span>{{ formatTime(log.created_at) }}</span>
                  </span>
                  <span>•</span>
                  <span>Latency: {{ log.latency > 0 ? `${log.latency}ms` : '—' }}</span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>
