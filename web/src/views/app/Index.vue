<script setup lang="ts">
import { ref } from 'vue';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { 
  Activity, 
  ShieldCheck, 
  ShieldAlert, 
  Percent, 
  Plus, 
  RefreshCw, 
  ExternalLink,
  CheckCircle2,
  AlertCircle
} from '@lucide/vue';

// Mock monitors data
const monitors = ref([
  { id: 1, name: 'Main API Gateway', url: 'https://api.ping-uptime.com/health', status: 'up', latency: '24ms', uptime: '99.98%', lastCheck: 'Just now' },
  { id: 2, name: 'Web Dashboard App', url: 'https://ping-uptime.com', status: 'up', latency: '42ms', uptime: '100.00%', lastCheck: '1 min ago' },
  { id: 3, name: 'Auth Module Service', url: 'https://auth.ping-uptime.com/validate', status: 'up', latency: '18ms', uptime: '99.95%', lastCheck: '2 mins ago' },
  { id: 4, name: 'Database Primary Replica', url: 'https://db.ping-uptime.com/ping', status: 'down', latency: '0ms', uptime: '98.42%', lastCheck: '5 mins ago' },
  { id: 5, name: 'Discord Webhook Delivery', url: 'https://discord.com/api/webhooks', status: 'up', latency: '124ms', uptime: '100.00%', lastCheck: '3 mins ago' }
]);

// Mock logs data
const logs = ref([
  { id: 1, time: '13:42:10', type: 'up', message: 'Main API Gateway responded (200 OK)', latency: '24ms' },
  { id: 2, time: '13:40:05', type: 'down', message: 'Database Primary Replica failed (503 Service Unavailable)', latency: '0ms' },
  { id: 3, time: '13:38:55', type: 'up', message: 'Auth Module Service responded (200 OK)', latency: '18ms' },
  { id: 4, time: '13:35:12', type: 'up', message: 'Web Dashboard App responded (200 OK)', latency: '42ms' }
]);

const totalMonitors = ref(5);
const upCount = ref(4);
const downCount = ref(1);
const avgUptime = ref('99.67%');
</script>

<template>
  <div class="p-8 space-y-8 max-w-[1600px] mx-auto">
    <!-- Top Greeting and Actions -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground">Welcome back, Admin!</h2>
        <p class="text-xs text-muted-foreground">Here is the real-time status of your self-hosted services.</p>
      </div>
      <div class="flex items-center gap-2">
        <button class="inline-flex items-center justify-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-lg border border-border bg-card text-foreground hover:bg-muted transition-colors">
          <RefreshCw class="w-3.5 h-3.5" />
          <span>Refresh</span>
        </button>
        <button class="inline-flex items-center justify-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-lg bg-primary hover:bg-primary/90 text-primary-foreground transition-colors shadow-sm shadow-primary/10">
          <Plus class="w-3.5 h-3.5" />
          <span>Add Monitor</span>
        </button>
      </div>
    </div>

    <!-- Quick Stats Grid -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <!-- Total Monitors Card -->
      <Card class="hover:shadow-md transition-shadow border-border/50">
        <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Total Monitors</CardTitle>
          <div class="p-2 rounded-lg bg-slate-100 dark:bg-slate-900">
            <Activity class="w-4 h-4 text-slate-500" />
          </div>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-black text-foreground">{{ totalMonitors }}</div>
          <p class="text-[10px] text-muted-foreground mt-1">Configured endpoints</p>
        </CardContent>
      </Card>

      <!-- Up Monitors Card -->
      <Card class="hover:shadow-md transition-shadow border-border/50">
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
      <Card class="hover:shadow-md transition-shadow border-border/50">
        <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Down Services</CardTitle>
          <div class="p-2 rounded-lg bg-red-500/10 dark:bg-red-500/5">
            <ShieldAlert class="w-4 h-4 text-red-500" />
          </div>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-black text-red-500">{{ downCount }}</div>
          <p class="text-[10px] text-muted-foreground mt-1">Action required immediately</p>
        </CardContent>
      </Card>

      <!-- Avg Uptime Card -->
      <Card class="hover:shadow-md transition-shadow border-border/50">
        <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Avg Uptime</CardTitle>
          <div class="p-2 rounded-lg bg-indigo-500/10 dark:bg-indigo-500/5">
            <Percent class="w-4 h-4 text-indigo-500" />
          </div>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-black text-foreground">{{ avgUptime }}</div>
          <p class="text-[10px] text-muted-foreground mt-1">Last 30 days average</p>
        </CardContent>
      </Card>
    </div>

    <!-- Main Grid Content: Monitors list & Incidents log -->
    <div class="grid gap-6 lg:grid-cols-3">
      <!-- Monitors List Card -->
      <Card class="lg:col-span-2 border-border/50">
        <CardHeader class="pb-3 border-b border-border/40">
          <CardTitle class="text-sm font-bold text-foreground">Operational Status</CardTitle>
          <CardDescription class="text-xs">Uptime statistics and latency responses for registered endpoints.</CardDescription>
        </CardHeader>
        <CardContent class="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead class="w-[30px]"></TableHead>
                <TableHead class="text-xs font-bold uppercase text-muted-foreground">Monitor</TableHead>
                <TableHead class="text-xs font-bold uppercase text-muted-foreground">Latency</TableHead>
                <TableHead class="text-xs font-bold uppercase text-muted-foreground">Uptime</TableHead>
                <TableHead class="text-xs font-bold uppercase text-muted-foreground">Last Checked</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="item in monitors" :key="item.id" class="hover:bg-muted/40 transition-colors">
                <!-- Status indicator -->
                <TableCell class="py-4">
                  <div class="flex items-center justify-center">
                    <span :class="[
                      'h-2.5 w-2.5 rounded-full ring-4 shrink-0',
                      item.status === 'up' 
                        ? 'bg-emerald-500 ring-emerald-500/20' 
                        : 'bg-red-500 ring-red-500/20 animate-pulse'
                    ]"></span>
                  </div>
                </TableCell>
                
                <!-- Monitor detail -->
                <TableCell class="font-medium text-foreground py-4">
                  <div class="flex flex-col gap-0.5 max-w-[280px]">
                    <span class="text-xs font-bold truncate">{{ item.name }}</span>
                    <a :href="item.url" target="_blank" class="text-[10px] text-muted-foreground hover:underline inline-flex items-center gap-1">
                      <span class="truncate">{{ item.url }}</span>
                      <ExternalLink class="w-2.5 h-2.5 shrink-0" />
                    </a>
                  </div>
                </TableCell>

                <!-- Latency response -->
                <TableCell class="py-4">
                  <span class="text-xs font-semibold" :class="item.status === 'up' ? 'text-foreground' : 'text-muted-foreground'">
                    {{ item.status === 'up' ? item.latency : '—' }}
                  </span>
                </TableCell>

                <!-- Uptime -->
                <TableCell class="py-4">
                  <Badge variant="outline" :class="[
                    'text-[10px] font-bold py-0.5 px-1.5 border/50',
                    item.status === 'up' 
                      ? 'bg-emerald-500/5 text-emerald-600 dark:text-emerald-400 border-emerald-500/20' 
                      : 'bg-red-500/5 text-red-600 dark:text-red-400 border-red-500/20'
                  ]">
                    {{ item.uptime }}
                  </Badge>
                </TableCell>

                <!-- Last Checked -->
                <TableCell class="text-muted-foreground text-xs py-4">
                  {{ item.lastCheck }}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <!-- Incident Log Feed -->
      <Card class="border-border/50">
        <CardHeader class="pb-3 border-b border-border/40">
          <CardTitle class="text-sm font-bold text-foreground">Recent Events</CardTitle>
          <CardDescription class="text-xs">Latest pings and recovery transitions logs.</CardDescription>
        </CardHeader>
        <CardContent class="pt-4 space-y-4">
          <div v-for="log in logs" :key="log.id" class="flex items-start gap-3 text-xs leading-relaxed">
            <!-- Event Icon -->
            <div class="mt-0.5">
              <CheckCircle2 v-if="log.type === 'up'" class="w-4 h-4 text-emerald-500" />
              <AlertCircle v-else class="w-4 h-4 text-red-500" />
            </div>
            
            <!-- Log Message detail -->
            <div class="flex-1 flex flex-col gap-0.5">
              <span class="font-semibold text-foreground">{{ log.message }}</span>
              <div class="flex items-center gap-2 text-[10px] text-muted-foreground">
                <span>{{ log.time }}</span>
                <span>•</span>
                <span>Latency: {{ log.type === 'up' ? log.latency : '—' }}</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
