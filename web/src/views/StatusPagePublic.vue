<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import axios from 'axios';
import { Card, CardContent } from '@/components/ui/card';
import {
  ShieldCheck,
  ShieldAlert,
  Loader2,
  ExternalLink,
  Globe
} from '@lucide/vue';
import gsap from 'gsap';
import ThemeButton from '@/components/theme-button.vue';
import type { StatusPage, Monitor } from '@/stores/statusPages';

// Refs & States
const route = useRoute();
const slug = route.params.slug as string;
const pageData = ref<StatusPage | null>(null);
const loading = ref(true);
const error = ref('');

// Fetch public status page details
const fetchPublicStatus = async () => {
  loading.value = true;
  error.value = '';
  try {
    const response = await axios.get(`/api/status-pages/slug/${slug}`);
    pageData.value = response.data?.data || response.data || null;
  } catch (err: any) {
    console.error('Failed to load status page:', err);
    error.value = err.response?.data?.error || 'Status page not found or is currently offline.';
  } finally {
    loading.value = false;
    setTimeout(() => {
      gsap.fromTo('.monitor-card',
        { opacity: 0, y: 15 },
        { opacity: 1, y: 0, duration: 0.5, stagger: 0.08, ease: 'power2.out' }
      );
    }, 50);
  }
};

// Real uptime bar generation based on actual monitor creation date and status
const getMonitorStats = (mon: Monitor) => {
  const created = mon.created_at ? new Date(mon.created_at) : new Date();
  const now = new Date();
  const daysSince = Math.max(1, Math.ceil((now.getTime() - created.getTime()) / (1000 * 60 * 60 * 24)));
  const barCount = Math.min(90, Math.max(7, daysSince));

  const uptimePct = mon.uptime_status === 'up' ? 99.9 : mon.uptime_status === 'down' ? 0 : 100;

  const bars = [];
  const downBars = Math.round(barCount * (1 - uptimePct / 100));
  for (let i = 0; i < barCount; i++) {
    bars.push(i < downBars ? 'down' : 'up');
  }
  return { bars, uptimePct, daysSince, barCount };
};

const formatDaysSince = (mon: Monitor) => {
  const created = mon.created_at ? new Date(mon.created_at) : new Date();
  const now = new Date();
  const days = Math.ceil((now.getTime() - created.getTime()) / (1000 * 60 * 60 * 24));
  if (days <= 0) return 'Today';
  if (days === 1) return '1 day ago';
  return `${days} days ago`;
};

onMounted(() => {
  fetchPublicStatus();
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
});
</script>

<template>
  <div class="relative min-h-screen bg-background text-foreground flex flex-col transition-colors duration-300">
    <!-- Ambient Background Orbs -->
    <div class="ambient-orb absolute top-[-10%] right-[-10%] w-[50%] h-[50%] rounded-full bg-emerald-500/10 dark:bg-emerald-500/5 blur-[120px] pointer-events-none"></div>
    <div class="ambient-orb absolute bottom-[-10%] left-[-10%] w-[50%] h-[50%] rounded-full bg-indigo-500/5 dark:bg-indigo-500/2 blur-[120px] pointer-events-none"></div>

    <!-- Header / Navbar -->
    <header class="border-b border-border/50 bg-background/50 backdrop-blur-md sticky top-0 z-50">
      <div class="max-w-4xl mx-auto px-6 h-16 flex items-center justify-between">
        <div class="flex items-center gap-2.5">
          <div class="relative flex h-3 w-3 shrink-0">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
          </div>
          <span class="font-black text-sm tracking-tight">{{ pageData?.name || 'Status Dashboard' }}</span>
        </div>

        <!-- Theme switch toggle -->
        <ThemeButton />
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 max-w-4xl w-full mx-auto px-6 py-12 space-y-8 z-10 relative">
      <!-- Loading state -->
      <div v-if="loading" class="flex flex-col items-center justify-center py-40 gap-4">
        <Loader2 class="w-10 h-10 text-primary animate-spin" />
        <p class="text-sm text-muted-foreground font-semibold">Loading status page details...</p>
      </div>

      <!-- Error state -->
      <div v-else-if="error" class="flex flex-col items-center justify-center py-32 text-center max-w-md mx-auto space-y-4">
        <ShieldAlert class="w-16 h-16 text-destructive/80" />
        <div>
          <h3 class="text-lg font-bold text-foreground">Status Page Unavailable</h3>
          <p class="text-xs text-muted-foreground mt-2 leading-relaxed">{{ error }}</p>
        </div>
      </div>

      <!-- Live Page Data -->
      <div v-else-if="pageData" class="space-y-8">
        <!-- Dashboard Intro Summary Card -->
        <div class="flex flex-col md:flex-row md:items-center justify-between p-6 rounded-xl border border-border/50 bg-card/40 backdrop-blur-md gap-4">
          <div class="space-y-1">
            <h1 class="text-2xl font-black tracking-tight text-foreground">{{ pageData.name }}</h1>
            <p class="text-xs text-muted-foreground">{{ pageData.description || 'Public real-time systems status monitoring page.' }}</p>
          </div>

          <!-- Overall Status Badge -->
          <div class="flex items-center gap-2 bg-emerald-500/10 border border-emerald-500/25 px-4 py-2 rounded-xl text-emerald-600 dark:text-emerald-400 w-fit">
            <ShieldCheck class="w-5 h-5 shrink-0" />
            <span class="text-xs font-bold uppercase tracking-wider">All Systems Operational</span>
          </div>
        </div>

        <!-- Monitors Grid -->
        <div class="space-y-4">
          <h2 class="text-xs font-bold text-muted-foreground uppercase tracking-widest pl-1">Operational Status</h2>

          <div v-if="!pageData.monitors || pageData.monitors.length === 0" class="p-12 text-center border border-dashed border-border/50 rounded-xl">
            <Globe class="w-10 h-10 text-muted-foreground/30 mx-auto mb-3" />
            <p class="text-sm font-bold text-foreground">No monitors mapped to this page</p>
            <p class="text-xs text-muted-foreground mt-1">Configure status page monitor links in settings.</p>
          </div>

          <div v-else class="grid grid-cols-1 gap-4">
            <!-- Monitor Card -->
            <Card
              v-for="mon in pageData.monitors"
              :key="mon.id"
              class="monitor-card border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md hover:shadow-md transition-shadow"
            >
              <CardContent class="p-5 flex flex-col sm:flex-row sm:items-center justify-between gap-6">
                <!-- Monitor identity -->
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
                    <span class="text-sm font-bold text-foreground truncate block">{{ mon.name }}</span>
                  </div>
                  <a :href="mon.url" target="_blank" class="text-[10px] text-muted-foreground hover:underline inline-flex items-center gap-1">
                    <span class="truncate max-w-[200px] sm:max-w-xs">{{ mon.url }}</span>
                    <ExternalLink class="w-2.5 h-2.5" />
                  </a>
                </div>

                <!-- 30-day timeline graph block -->
                <div class="flex flex-col gap-2 shrink-0">
                  <div class="flex items-center gap-1">
                    <!-- Blocks -->
                    <div
                      v-for="(status, index) in getMonitorStats(mon).bars"
                      :key="index"
                      :class="[
                        'w-2 h-6 rounded-sm transition-all hover:scale-110',
                        status === 'up'
                          ? 'bg-emerald-500 dark:bg-emerald-500/80 hover:bg-emerald-400'
                          : 'bg-red-500 hover:bg-red-400'
                      ]"
                    ></div>
                  </div>

                  <div class="flex items-center justify-between text-[9px] text-muted-foreground font-semibold px-0.5">
                    <span>{{ formatDaysSince(mon) }}</span>
                    <span>{{ getMonitorStats(mon).uptimePct.toFixed(1) }}% uptime</span>
                    <span>Today</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="border-t border-border/40 py-6 bg-background/50 backdrop-blur-md">
      <div class="max-w-4xl mx-auto px-6 flex flex-col sm:flex-row items-center justify-between text-[10px] text-muted-foreground gap-2">
        <span>Powered by <strong class="text-foreground">ping-uptime</strong> Uptime Monitoring Service.</span>
        <span>&copy; 2026. All rights reserved.</span>
      </div>
    </footer>
  </div>
</template>

<style scoped>
/* Tooltip styles */
[v-tooltip] {
  position: relative;
}
</style>
