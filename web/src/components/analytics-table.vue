<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { ExternalLink, User, ChevronLeft, ChevronRight } from '@lucide/vue';
import type { MonitorStats } from '@/stores/analytics';

const props = defineProps<{
  stats: MonitorStats[];
  isAdmin: boolean;
  getMonitorOwner: (monitorId: number) => string | null;
}>();

const currentPage = ref(1);
const pageSize = ref(10);

const paginatedStats = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return props.stats.slice(start, end);
});

const totalPages = computed(() => {
  return Math.ceil(props.stats.length / pageSize.value) || 1;
});

watch(() => totalPages.value, () => {
  if (currentPage.value > totalPages.value) {
    currentPage.value = totalPages.value;
  }
});
</script>

<template>
  <div class="flex flex-col">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead class="w-[80px] text-center">Status</TableHead>
          <TableHead>Monitor</TableHead>
          <TableHead>Uptime %</TableHead>
          <TableHead>Checks (Failed/Total)</TableHead>
          <TableHead>Avg Latency</TableHead>
          <TableHead>Min/Max Latency</TableHead>
          <TableHead v-if="isAdmin">Owner</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="item in paginatedStats" :key="item.monitor_id">
          <!-- Status -->
          <TableCell>
            <div class="flex items-center justify-center">
              <span :class="[
                'h-2.5 w-2.5 rounded-full ring-4 shrink-0',
                item.status === 'paused'
                  ? 'bg-amber-500 ring-amber-500/20'
                  : item.status === 'operational'
                  ? 'bg-emerald-500 ring-emerald-500/20'
                  : item.status === 'degraded'
                  ? 'bg-amber-500 ring-amber-500/20'
                  : item.status === 'outage'
                  ? 'bg-red-500 ring-red-500/20 animate-pulse'
                  : 'bg-slate-400 ring-slate-400/20'
              ]"></span>
            </div>
          </TableCell>

          <!-- Monitor -->
          <TableCell class="font-medium">
            <div class="flex flex-col gap-0.5 max-w-sm sm:max-w-md md:max-w-lg lg:max-w-xl">
              <span class="text-xs font-bold truncate">{{ item.monitor_name }}</span>
              <a :href="item.monitor_url" target="_blank" @click.stop class="text-[10px] text-muted-foreground hover:underline inline-flex items-center gap-1">
                <span class="truncate">{{ item.monitor_url }}</span>
                <ExternalLink class="w-2.5 h-2.5 shrink-0" />
              </a>
            </div>
          </TableCell>

          <!-- Uptime % -->
          <TableCell>
            <Badge variant="outline" :class="[
              'text-[10px] font-bold py-0.5 px-2 transition-colors duration-200',
              item.uptime_pct >= 99
                ? 'border-emerald-500/30 text-emerald-500 bg-emerald-500/5'
                : item.uptime_pct >= 95
                ? 'border-amber-500/30 text-amber-500 bg-amber-500/5'
                : 'border-rose-500/30 text-rose-500 bg-rose-500/5'
            ]">
              {{ item.uptime_pct.toFixed(3) }}%
            </Badge>
          </TableCell>

          <!-- Checks (Failed/Total) -->
          <TableCell class="text-xs text-muted-foreground">
            <span class="text-rose-500 font-semibold">{{ item.failed_checks }}</span> / {{ item.total_checks }}
          </TableCell>

          <!-- Avg Latency -->
          <TableCell class="text-xs font-semibold text-blue-500">
            {{ Math.round(item.avg_latency) }} ms
          </TableCell>

          <!-- Min/Max Latency -->
          <TableCell class="text-xs text-muted-foreground">
            {{ Math.round(item.min_latency) }} ms / {{ Math.round(item.max_latency) }} ms
          </TableCell>

          <!-- Owner -->
          <TableCell v-if="isAdmin" class="text-xs text-muted-foreground">
            <div class="flex items-center gap-1 max-w-[180px]">
              <User class="w-3.5 h-3.5 text-muted-foreground/70 shrink-0" />
              <span class="truncate">{{ getMonitorOwner(item.monitor_id) || 'Loading...' }}</span>
            </div>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>

    <!-- Pagination Footer -->
    <div class="flex flex-col sm:flex-row items-center justify-between gap-4 p-4 border-t border-border/40 text-xs text-muted-foreground bg-muted/5">
      <!-- Page Size Selector -->
      <div class="flex items-center gap-2">
        <span>Show</span>
        <Select :model-value="String(pageSize)" @update:model-value="pageSize = Number($event); currentPage = 1">
          <SelectTrigger class="h-8 w-16 text-xs bg-background">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="10">10</SelectItem>
            <SelectItem value="20">20</SelectItem>
            <SelectItem value="50">50</SelectItem>
            <SelectItem value="100">100</SelectItem>
          </SelectContent>
        </Select>
        <span>entries</span>
      </div>

      <!-- Info Text -->
      <div class="font-medium">
        Showing {{ Math.min(props.stats.length, (currentPage - 1) * pageSize + 1) }} to {{ Math.min(props.stats.length, currentPage * pageSize) }} of {{ props.stats.length }} entries
      </div>

      <!-- Navigation Buttons -->
      <div class="flex items-center gap-1.5">
        <Button 
          variant="outline" 
          size="icon" 
          class="h-8 w-8" 
          :disabled="currentPage === 1" 
          @click="currentPage--"
        >
          <ChevronLeft class="w-4 h-4" />
        </Button>
        <span class="px-3 py-1 font-bold text-foreground">
          Page {{ currentPage }} of {{ totalPages }}
        </span>
        <Button 
          variant="outline" 
          size="icon" 
          class="h-8 w-8" 
          :disabled="currentPage === totalPages" 
          @click="currentPage++"
        >
          <ChevronRight class="w-4 h-4" />
        </Button>
      </div>
    </div>
  </div>
</template>
