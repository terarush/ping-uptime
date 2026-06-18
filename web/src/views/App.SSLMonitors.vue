<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useSslMonitors } from '@/composables/useSslMonitors';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ShieldCheck, RefreshCw, Loader2, Play, Trash2 } from '@lucide/vue';
import gsap from 'gsap';

const {
  certs,
  loading,
  error,
  fetchAll,
  checkAll,
  deleteCert,
  checkMonitor,
  formatDate,
  daysClass,
  statusVariant,
} = useSslMonitors();

// Fetch monitors too for name resolution
const monitorsMap = ref<Record<number, { id: number; name: string; url: string }>>({});

const loadData = async () => {
  try {
    // Load monitors for name resolution
    const { default: ExtendedFetch } = await import('@/lib/fetch');
    const monResp = await ExtendedFetch.get('/monitors');
    const monList = monResp.data?.data || monResp.data || [];
    const map: Record<number, any> = {};
    for (const m of monList) {
      map[m.id] = m;
    }
    monitorsMap.value = map;

    await fetchAll();
  } catch (err) {
    console.error('Failed to load SSL data:', err);
  } finally {
    setTimeout(animateRows, 50);
  }
};

const handleCheckAll = async () => {
  await checkAll();
  setTimeout(animateRows, 50);
};

const handleCheckMonitor = async (monitorId: number) => {
  await checkMonitor(monitorId);
};

const handleDelete = async (id: number) => {
  await deleteCert(id);
};

// GSAP animation
const animateRows = () => {
  gsap.fromTo('.ssl-row',
    { opacity: 0, y: 12 },
    { opacity: 1, y: 0, duration: 0.4, stagger: 0.04, ease: 'power2.out' }
  );
};

onMounted(() => {
  loadData();
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <!-- Ambient Background Orb (green-tinted) -->
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-emerald-500/10 dark:bg-emerald-500/5 blur-[100px] pointer-events-none"></div>

    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <ShieldCheck class="w-6 h-6 text-primary" />
          <span>SSL Certificates</span>
        </h2>
        <p class="text-xs text-muted-foreground">Monitor TLS certificate expiry for all monitored endpoints.</p>
      </div>

      <!-- Actions -->
      <div class="flex gap-2">
        <Button variant="outline" size="sm" @click="handleCheckAll" :disabled="loading" class="h-9">
          <Play class="w-4 h-4 mr-1.5" />
          <span>Check All</span>
        </Button>
        <Button variant="outline" size="sm" @click="loadData" class="h-9">
          <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
          <span>Refresh</span>
        </Button>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="error && certs.length === 0" class="z-10 relative">
      <Card class="border-destructive/50 bg-destructive/5">
        <CardContent class="py-6 text-center text-sm text-destructive">
          {{ error }}
        </CardContent>
      </Card>
    </div>

    <!-- Main Container Card -->
    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <CardTitle class="text-sm font-bold text-foreground">Certificate Status</CardTitle>
        <CardDescription class="text-xs">TLS certificate checks and expiry tracking.</CardDescription>
      </CardHeader>
      <CardContent class="p-0">
        <!-- Loading State -->
        <div v-if="loading && certs.length === 0" class="flex flex-col items-center justify-center py-20 gap-3">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
          <p class="text-sm text-muted-foreground">Checking SSL certificates...</p>
        </div>

        <!-- Empty State -->
        <div v-else-if="certs.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
          <ShieldCheck class="w-12 h-12 text-emerald-500/30 mb-3" />
          <p class="text-sm font-bold text-foreground">No SSL certificates</p>
          <p class="text-xs text-muted-foreground mt-1">Check a monitor to begin tracking its TLS certificate.</p>
        </div>

        <!-- Table -->
        <Table v-else>
          <TableHeader>
            <TableRow>
              <TableHead>Domain</TableHead>
              <TableHead>Monitor</TableHead>
              <TableHead>Issuer</TableHead>
              <TableHead>Valid Until</TableHead>
              <TableHead>Days Remaining</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Last Checked</TableHead>
              <TableHead class="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="cert in certs" :key="cert.id" class="ssl-row text-xs hover:bg-muted/50 transition-colors">
              <TableCell class="font-medium">{{ cert.domain || '—' }}</TableCell>
              <TableCell class="text-muted-foreground">
                {{ monitorsMap[cert.monitor_id]?.name || `Monitor #${cert.monitor_id}` }}
              </TableCell>
              <TableCell class="text-muted-foreground max-w-[200px] truncate">{{ cert.issuer || '—' }}</TableCell>
              <TableCell class="text-muted-foreground whitespace-nowrap">
                {{ cert.valid_until ? formatDate(cert.valid_until) : '—' }}
              </TableCell>
              <TableCell>
                <span :class="['text-sm', daysClass(cert.days_remaining)]">
                  {{ cert.status === 'error' ? '—' : cert.days_remaining + 'd' }}
                </span>
              </TableCell>
              <TableCell>
                <Badge variant="outline" :class="['text-[10px] font-bold uppercase tracking-wide', statusVariant(cert.status)]">
                  {{ cert.status === 'expiring_soon' ? 'Expiring' : cert.status }}
                </Badge>
              </TableCell>
              <TableCell class="text-muted-foreground whitespace-nowrap">{{ formatDate(cert.checked_at) }}</TableCell>
              <TableCell class="text-right">
                <div class="flex justify-end gap-1">
                  <Button variant="ghost" size="icon" class="h-7 w-7" @click="handleCheckMonitor(cert.monitor_id)" title="Check Now">
                    <Play class="h-3.5 w-3.5" />
                  </Button>
                  <Button variant="ghost" size="icon" class="h-7 w-7 text-destructive hover:text-destructive" @click="handleDelete(cert.id)" title="Delete">
                    <Trash2 class="h-3.5 w-3.5" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </div>
</template>
