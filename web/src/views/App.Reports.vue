<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import ExtendedFetch from '@/lib/fetch';
import {
  BarChart3,
  Loader2,
  Download,
  TrendingUp,
  Activity,
  CheckCircle2,
  XCircle,
} from '@lucide/vue';
import gsap from 'gsap';

interface ReportData {
  from: string;
  to: string;
  stats?: {
    total_monitors: number;
    total_checks: number;
    uptime: number;
    avg_latency: number;
  };
  data?: any[];
  monitor_id?: number;
}

interface ReportStats {
  total_monitors: number;
  total_checks: number;
  uptime: number;
  avg_latency: number;
}

const dateFrom = ref('');
const dateTo = ref('');
const loading = ref(false);
const report = ref<ReportData | null>(null);
const error = ref('');
const monitors = ref<{ id: number; name: string }[]>([]);
const selectedMonitorID = ref<string>('');

const fetchMonitors = async () => {
  try {
    const res = await ExtendedFetch.get('/monitors');
    const data = res.data?.data || res.data || [];
    monitors.value = data.map((m: any) => ({ id: m.id, name: m.name }));
  } catch {
    // ignore
  }
};

const generateReport = async () => {
  if (!dateFrom.value || !dateTo.value) return;
  loading.value = true;
  error.value = '';
  try {
    const params = new URLSearchParams({ from: dateFrom.value, to: dateTo.value });
    if (selectedMonitorID.value) {
      params.set('monitor_id', selectedMonitorID.value);
    }
    const res = await ExtendedFetch.get(`/analytics/report?${params.toString()}`);
    report.value = res.data?.data || res.data;
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Failed to generate report.';
    report.value = null;
  } finally {
    loading.value = false;
  }
};

const exportCSV = () => {
  if (!report.value) return;
  const r = report.value;
  const s: ReportStats | undefined = r.stats;
  const rows = [
    ['Metric', 'Value'],
    ['From', r.from],
    ['To', r.to],
    ['Total Monitors', s?.total_monitors ?? ''],
    ['Total Checks', s?.total_checks ?? ''],
    ['Uptime %', s?.uptime ? (s.uptime * 100).toFixed(2) + '%' : ''],
    ['Avg Latency', s?.avg_latency ? s.avg_latency + 'ms' : ''],
  ];
  const csv = rows.map(r => r.join(',')).join('\n');
  const blob = new Blob([csv], { type: 'text/csv' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `uptime-report-${dateFrom.value}-${dateTo.value}.csv`;
  a.click();
  URL.revokeObjectURL(url);
};

const exportJSON = () => {
  if (!report.value) return;
  const blob = new Blob([JSON.stringify(report.value, null, 2)], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `uptime-report-${dateFrom.value}-${dateTo.value}.json`;
  a.click();
  URL.revokeObjectURL(url);
};

const uptimePercent = computed(() => {
  const s: ReportStats | undefined = report.value?.stats;
  if (!s?.uptime && s?.uptime !== 0) return null;
  return (s.uptime * 100).toFixed(2);
});

onMounted(async () => {
  await fetchMonitors();
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-indigo-500/5 dark:bg-indigo-500/2 blur-[100px] pointer-events-none"></div>

    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <BarChart3 class="w-6 h-6 text-primary" />
          <span>Reports</span>
        </h2>
        <p class="text-xs text-muted-foreground">Uptime reports and data export.</p>
      </div>
    </div>

    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <CardTitle class="text-sm font-bold text-foreground">Generate Report</CardTitle>
        <CardDescription class="text-xs">Select date range and optionally filter by monitor.</CardDescription>
      </CardHeader>
      <CardContent class="p-4 space-y-4">
        <div class="flex flex-col sm:flex-row items-end gap-4">
          <div class="space-y-1.5 flex-1">
            <Label class="text-xs">From</Label>
            <Input v-model="dateFrom" type="date" class="h-9" />
          </div>
          <div class="space-y-1.5 flex-1">
            <Label class="text-xs">To</Label>
            <Input v-model="dateTo" type="date" class="h-9" />
          </div>
          <div class="space-y-1.5 flex-1">
            <Label class="text-xs">Monitor (optional)</Label>
            <Select v-model="selectedMonitorID">
              <SelectTrigger class="h-9">
                <SelectValue placeholder="All monitors" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Monitors</SelectItem>
                <SelectItem v-for="m in monitors" :key="m.id" :value="String(m.id)">
                  {{ m.name }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <Button @click="generateReport" :disabled="loading || !dateFrom || !dateTo" class="h-9">
            <Loader2 v-if="loading" class="w-4 h-4 mr-1.5 animate-spin" />
            <TrendingUp v-else class="w-4 h-4 mr-1.5" />
            <span>Generate</span>
          </Button>
        </div>
      </CardContent>
    </Card>

    <div v-if="error" class="text-sm text-red-500 p-4 bg-red-500/10 rounded-lg border border-red-500/20">{{ error }}</div>

    <div v-if="report" class="z-10 relative">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <Card class="border-border/50 bg-card/60 backdrop-blur-md">
          <CardContent class="p-4 flex items-center gap-3">
            <div class="p-2 rounded-lg bg-primary/10"><Activity class="w-5 h-5 text-primary" /></div>
            <div>
              <p class="text-xs text-muted-foreground">Total Monitors</p>
              <p class="text-xl font-bold">{{ report.stats?.total_monitors || report.monitor_id ? 1 : 0 }}</p>
            </div>
          </CardContent>
        </Card>
        <Card class="border-border/50 bg-card/60 backdrop-blur-md">
          <CardContent class="p-4 flex items-center gap-3">
            <div class="p-2 rounded-lg bg-emerald-500/10"><CheckCircle2 class="w-5 h-5 text-emerald-500" /></div>
            <div>
              <p class="text-xs text-muted-foreground">Uptime</p>
              <p class="text-xl font-bold">{{ uptimePercent || '—' }}%</p>
            </div>
          </CardContent>
        </Card>
        <Card class="border-border/50 bg-card/60 backdrop-blur-md">
          <CardContent class="p-4 flex items-center gap-3">
            <div class="p-2 rounded-lg bg-blue-500/10"><TrendingUp class="w-5 h-5 text-blue-500" /></div>
            <div>
              <p class="text-xs text-muted-foreground">Avg Latency</p>
              <p class="text-xl font-bold">{{ report.stats?.avg_latency || 0 }}ms</p>
            </div>
          </CardContent>
        </Card>
        <Card class="border-border/50 bg-card/60 backdrop-blur-md">
          <CardContent class="p-4 flex items-center gap-3">
            <div class="p-2 rounded-lg bg-amber-500/10"><BarChart3 class="w-5 h-5 text-amber-500" /></div>
            <div>
              <p class="text-xs text-muted-foreground">Period</p>
              <p class="text-sm font-bold">{{ report.from }} → {{ report.to }}</p>
            </div>
          </CardContent>
        </Card>
      </div>

      <div class="flex gap-2 justify-end">
        <Button variant="outline" size="sm" class="h-9" @click="exportCSV">
          <Download class="w-4 h-4 mr-1.5" /> Export CSV
        </Button>
        <Button variant="outline" size="sm" class="h-9" @click="exportJSON">
          <Download class="w-4 h-4 mr-1.5" /> Export JSON
        </Button>
      </div>
    </div>
  </div>
</template>
