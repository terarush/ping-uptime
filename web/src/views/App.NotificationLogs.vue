<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useNotificationLogs } from '@/composables/useNotificationLogs';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import {
  ClipboardList,
  Search,
  Loader2,
  RefreshCw,
  CheckCircle2,
  X,
} from '@lucide/vue';
import gsap from 'gsap';

const { logs, loading, fetchLogs, formatDate, channelIcon } = useNotificationLogs();

const searchQuery = ref('');
const filterChannel = ref('');
const filterStatus = ref('');

const filteredLogs = computed(() => {
  if (!logs.value) return [];
  let result = logs.value;

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    result = result.filter(l =>
      l.recipient?.toLowerCase().includes(q) ||
      l.channel_type?.toLowerCase().includes(q) ||
      l.event_type?.toLowerCase().includes(q)
    );
  }
  if (filterChannel.value) {
    result = result.filter(l => l.channel_type === filterChannel.value);
  }
  if (filterStatus.value) {
    result = result.filter(l => l.status === filterStatus.value);
  }
  return result;
});

const fetchAll = async () => {
  try {
    await fetchLogs();
    animateRows();
  } catch {
    // handled in store
  }
};

const animateRows = () => {
  setTimeout(() => {
    gsap.fromTo('.log-row',
      { opacity: 0, y: 12 },
      { opacity: 1, y: 0, duration: 0.3, stagger: 0.03, ease: 'power2.out' }
    );
  }, 50);
};

onMounted(async () => {
  await fetchAll();
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-blue-500/5 dark:bg-blue-500/2 blur-[100px] pointer-events-none"></div>

    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <ClipboardList class="w-6 h-6 text-primary" />
          <span>Notification History</span>
        </h2>
        <p class="text-xs text-muted-foreground">Log of all sent notifications and delivery status.</p>
      </div>

      <Button variant="outline" size="sm" @click="fetchAll" class="h-9">
        <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
        <span>Refresh</span>
      </Button>
    </div>

    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <CardTitle class="text-sm font-bold text-foreground">Delivery Log</CardTitle>
            <CardDescription class="text-xs">Filter and review notification delivery status.</CardDescription>
          </div>

          <div class="flex items-center gap-2 flex-wrap">
            <Select v-model="filterChannel">
              <SelectTrigger class="w-36 h-9">
                <SelectValue placeholder="Channel" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="">All Channels</SelectItem>
                <SelectItem value="email">Email</SelectItem>
                <SelectItem value="discord">Discord</SelectItem>
                <SelectItem value="discord_bot">Discord Bot</SelectItem>
                <SelectItem value="telegram">Telegram</SelectItem>
                <SelectItem value="slack">Slack</SelectItem>
                <SelectItem value="webhook">Webhook</SelectItem>
              </SelectContent>
            </Select>

            <Select v-model="filterStatus">
              <SelectTrigger class="w-32 h-9">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="">All Status</SelectItem>
                <SelectItem value="sent">Sent</SelectItem>
                <SelectItem value="failed">Failed</SelectItem>
              </SelectContent>
            </Select>

            <div class="relative w-full sm:w-56">
              <Search class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input v-model="searchQuery" placeholder="Search..." class="pl-9 h-9" />
            </div>
          </div>
        </div>
      </CardHeader>

      <CardContent class="p-0">
        <div v-if="loading && logs.length === 0" class="flex flex-col items-center justify-center py-20 gap-3">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
          <p class="text-sm text-muted-foreground">Loading notification logs...</p>
        </div>

        <div v-else-if="filteredLogs.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
          <ClipboardList class="w-12 h-12 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-bold text-foreground">No notification logs</p>
          <p class="text-xs text-muted-foreground mt-1">Notifications will appear here when alerts are triggered.</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b border-border/40 text-xs text-muted-foreground">
                <th class="text-left font-medium px-4 py-3">Channel</th>
                <th class="text-left font-medium px-4 py-3">Event</th>
                <th class="text-left font-medium px-4 py-3">Recipient</th>
                <th class="text-center font-medium px-4 py-3">Status</th>
                <th class="text-right font-medium px-4 py-3">Sent At</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="log in filteredLogs"
                :key="log.id"
                class="log-row border-b border-border/20 hover:bg-accent/30 transition-colors"
              >
                <td class="px-4 py-3 text-sm">
                  <span class="inline-flex items-center gap-1.5">
                    <span class="text-base">{{ channelIcon(log.channel_type) }}</span>
                    <span class="capitalize">{{ log.channel_type.replace('_', ' ') }}</span>
                  </span>
                </td>
                <td class="px-4 py-3 text-sm capitalize">{{ log.event_type }}</td>
                <td class="px-4 py-3 text-sm text-muted-foreground">{{ log.recipient || '—' }}</td>
                <td class="px-4 py-3 text-center">
                  <span v-if="log.status === 'sent'" class="inline-flex items-center gap-1 text-emerald-500 text-xs font-medium">
                    <CheckCircle2 class="w-3.5 h-3.5" /> Sent
                  </span>
                  <span v-else class="inline-flex items-center gap-1 text-red-500 text-xs font-medium" :title="log.error_message">
                    <X class="w-3.5 h-3.5" /> Failed
                  </span>
                </td>
                <td class="px-4 py-3 text-sm text-right text-muted-foreground">
                  {{ formatDate(log.sent_at) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
