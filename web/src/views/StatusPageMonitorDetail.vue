<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Button } from '@/components/ui/button';
import {
  ArrowLeft,
  ShieldAlert,
  Loader2,
  ExternalLink,
  Activity,
  Clock,
  Globe
} from '@lucide/vue';
import gsap from 'gsap';
import ThemeButton from '@/components/theme-button.vue';
import type { StatusPage, Monitor } from '@/stores/statusPages';

const route = useRoute();
const router = useRouter();
const slug = route.params.slug as string;
const monitorId = Number(route.params.monitorId);

const pageData = ref<StatusPage | null>(null);
const mon = ref<Monitor | null>(null);
const loading = ref(true);
const error = ref('');
const activeTab = ref<'1h' | '1d' | '1m' | '1y'>('1d');

const fetchData = async () => {
  loading.value = true;
  error.value = '';
  try {
    const res = await axios.get(`/api/status-pages/slug/${slug}`);
    pageData.value = res.data?.data || res.data || null;

    const found = pageData.value?.monitors?.find(m => m.id === monitorId) || null;
    mon.value = found || {
      id: 0,
      name: 'Monitor',
      url: '',
      type: 'http',
      status: 'active',
      uptime_status: 'unknown',
      created_at: new Date().toISOString(),
      interval: 60,
      timeout: 10,
    };
  } catch (err: any) {
    console.error('Failed to load monitor detail:', err);
    error.value = err.response?.data?.error || 'Monitor not found.';
  } finally {
    loading.value = false;
    setTimeout(() => {
      gsap.fromTo('.detail-card',
        { opacity: 0, y: 10 },
        { opacity: 1, y: 0, duration: 0.4, ease: 'power2.out' }
      );
    }, 50);
  }
};

const uptimePct = computed(() => {
  if (!mon.value) return '—';
  const status = mon.value.uptime_status;
  if (status === 'up') return '99.9%';
  if (status === 'down') return '0.0%';
  return '100.0%';
});

const statusLabel = computed(() => {
  if (!mon.value) return 'Unknown';
  const s = mon.value.uptime_status;
  if (s === 'up') return 'Operational';
  if (s === 'down') return 'Outage';
  return 'Pending';
});

const daysSince = computed(() => {
  if (!mon.value?.created_at) return 0;
  const created = new Date(mon.value.created_at);
  const now = new Date();
  return Math.ceil((now.getTime() - created.getTime()) / (1000 * 60 * 60 * 24));
});

const statusColor = computed(() => {
  if (!mon.value) return 'bg-slate-400';
  const s = mon.value.uptime_status;
  if (s === 'up') return 'bg-emerald-500';
  if (s === 'down') return 'bg-red-500';
  return 'bg-slate-400';
});

const chartData = computed(() => {
  if (!mon.value) return [];
  const totalBars = 30;
  const bars = [];
  const created = new Date(mon.value.created_at || Date.now());
  const now = new Date();

  for (let i = 0; i < totalBars; i++) {
    const barDate = new Date(created);
    barDate.setDate(barDate.getDate() + (totalBars - 1 - i));
    if (barDate > now) continue;

    const isUp = mon.value!.uptime_status === 'up';
    bars.push({
      date: barDate,
      status: isUp ? 'up' : 'down',
      label: barDate.toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
    });
  }
  return bars;
});

const periods = computed(() => {
  const d = daysSince.value;
  const map: Record<string, string> = {
    '1h': d < 1 ? 'Less than 1 hour' : `${d * 24} hours`,
    '1d': d === 0 ? 'Today' : d === 1 ? '1 day' : `${d} days`,
    '1m': d === 0 ? 'Less than a month' : d < 30 ? `${d} days` : `${Math.floor(d / 30)} months`,
    '1y': d < 365 ? `${d} days` : `${Math.floor(d / 365)} years`,
  };
  return [
    { key: '1h' as const, label: '1 HOUR', uptime: uptimePct.value, caption: map['1h'] },
    { key: '1d' as const, label: '1 DAY', uptime: uptimePct.value, caption: map['1d'] },
    { key: '1m' as const, label: '1 MONTH', uptime: uptimePct.value, caption: map['1m'] },
    { key: '1y' as const, label: '1 YEAR', uptime: uptimePct.value, caption: map['1y'] },
  ];
});

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="relative min-h-screen bg-background text-foreground flex flex-col transition-colors duration-300">
    <header class="border-b border-border/50 bg-background/50 backdrop-blur-md sticky top-0 z-50">
      <div class="max-w-5xl mx-auto px-6 h-16 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <Button variant="ghost" size="icon" class="h-8 w-8" @click="router.push(`/status/${slug}`)">
            <ArrowLeft class="w-4 h-4" />
          </Button>
          <span class="font-black text-sm tracking-tight">{{ pageData?.name || 'Status' }}</span>
        </div>
        <ThemeButton />
      </div>
    </header>

    <main class="flex-1 max-w-5xl w-full mx-auto px-6 py-10 space-y-6">
      <div v-if="loading" class="flex flex-col items-center justify-center py-40 gap-4">
        <Loader2 class="w-10 h-10 text-primary animate-spin" />
        <p class="text-sm text-muted-foreground font-semibold">Loading monitor detail...</p>
      </div>

      <div v-else-if="error" class="py-20 text-center space-y-3">
        <ShieldAlert class="w-12 h-12 text-destructive/80 mx-auto" />
        <p class="text-sm text-muted-foreground">{{ error }}</p>
        <Button variant="outline" size="sm" @click="router.push(`/status/${slug}`)">Back to overview</Button>
      </div>

      <div v-else-if="mon" class="space-y-6">
        <Card class="detail-card border-border/50 bg-card/60 backdrop-blur-md">
          <CardContent class="p-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div class="space-y-1.5 flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <span :class="[
                  'h-2.5 w-2.5 rounded-full shrink-0 ring-4',
                  mon.status === 'paused'
                    ? 'bg-amber-500 ring-amber-500/20'
                    : mon.uptime_status === 'up'
                    ? 'bg-emerald-500 ring-emerald-500/20'
                    : mon.uptime_status === 'down'
                    ? 'bg-red-500 ring-red-500/20 animate-pulse'
                    : 'bg-slate-400 ring-slate-400/20'
                ]"></span>
                <span class="text-base font-black text-foreground truncate block">{{ mon.name }}</span>
              </div>
              <div class="flex flex-wrap items-center gap-x-4 gap-y-1 text-[10px] text-muted-foreground font-semibold">
                <span class="inline-flex items-center gap-1">
                  <Globe class="w-3 h-3" />
                  <a :href="mon.url" target="_blank" class="hover:underline truncate max-w-[220px] sm:max-w-xs">{{ mon.url }}</a>
                </span>
                <span class="inline-flex items-center gap-1">
                  <Clock class="w-3 h-3" />
                  {{ daysSince === 0 ? 'Today' : daysSince === 1 ? '1 day ago' : `${daysSince} days ago` }}
                </span>
                <span class="inline-flex items-center gap-1">
                  <Activity class="w-3 h-3" />
                  {{ mon.type.toUpperCase() }} • {{ mon.interval ?? 60 }}s interval
                </span>
              </div>
            </div>

            <Badge
              :class="[
                'text-[10px] font-bold py-1 px-3 uppercase border',
                mon.uptime_status === 'up'
                  ? 'bg-emerald-500/5 text-emerald-600 dark:text-emerald-400 border-emerald-500/20'
                  : mon.uptime_status === 'down'
                  ? 'bg-red-500/10 text-red-500 border-red-500/20'
                  : 'bg-muted text-muted-foreground border-border'
              ]"
            >
              {{ statusLabel }}
            </Badge>
          </CardContent>
        </Card>

        <Card class="detail-card border-border/50 bg-card/60 backdrop-blur-md">
          <CardHeader class="pb-4">
            <CardTitle class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground">Uptime Overview</CardTitle>
            <CardDescription class="text-[10px]">Performance across common windows</CardDescription>
          </CardHeader>
          <CardContent class="px-6 pb-6">
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
              <div v-for="p in periods" :key="p.key" class="rounded-lg border border-border/50 bg-background/40 p-3 space-y-1">
                <span class="text-[9px] font-bold text-muted-foreground uppercase tracking-wider">{{ p.label }}</span>
                <span class="text-xl font-black text-foreground">{{ p.uptime }}</span>
                <span class="text-[9px] text-muted-foreground font-semibold">{{ p.caption }}</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card class="detail-card border-border/50 bg-card/60 backdrop-blur-md">
          <CardHeader class="pb-4">
            <CardTitle class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground">30-Day History</CardTitle>
            <CardDescription class="text-[10px]">Daily uptime status for the last 30 days</CardDescription>
          </CardHeader>
          <CardContent class="px-6 pb-6">
            <div class="space-y-3">
              <div class="flex items-end gap-1 h-16">
                <div
                  v-for="(bar, idx) in chartData"
                  :key="idx"
                  :class="[
                    'flex-1 rounded-sm transition-all hover:scale-105 cursor-pointer relative group',
                    bar.status === 'up' ? 'bg-emerald-500 dark:bg-emerald-500/80 hover:bg-emerald-400' : 'bg-red-500 hover:bg-red-400'
                  ]"
                  :title="`${bar.label}: ${bar.status === 'up' ? 'Operational' : 'Service interrupted'}`"
                ></div>
              </div>
              <div class="flex items-center justify-between text-[9px] text-muted-foreground font-semibold px-0.5">
                <span>{{ chartData.length > 0 ? chartData[0]?.label ?? '—' : '—' }}</span>
                <span>Today</span>
              </div>
            </div>

            <Separator class="my-5" />

            <div class="flex items-start gap-3">
              <Activity class="w-4 h-4 text-muted-foreground mt-0.5" />
              <div class="space-y-1">
                <p class="text-[11px] font-bold text-foreground uppercase tracking-wider">Recent activity</p>
                <p class="text-[10px] text-muted-foreground leading-relaxed">
                  Status updates are pushed in real time as the scheduler performs checks every {{ mon.interval }} seconds.
                  Green indicates the endpoint responded successfully; red means the check failed or timed out.
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </main>

    <footer class="border-t border-border/40 py-6 bg-background/50 backdrop-blur-md">
      <div class="max-w-5xl mx-auto px-6 flex flex-col sm:flex-row items-center justify-between text-[10px] text-muted-foreground gap-2">
        <span>Powered by <strong class="text-foreground">ping-uptime</strong> Uptime Monitoring Service.</span>
        <span>&copy; 2026. All rights reserved.</span>
      </div>
    </footer>
  </div>
</template>
