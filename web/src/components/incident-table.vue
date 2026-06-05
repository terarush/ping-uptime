<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { Clock, Zap, ChevronLeft, ChevronRight } from '@lucide/vue';
import type { Incident, MonitorItem } from '@/stores/incidents';

const props = defineProps<{
  incidents: Incident[];
  monitorsMap: Record<number, MonitorItem>;
  formatDate: (dStr?: string) => string;
  getDuration: (inc: Incident) => string;
}>();

const currentPage = ref(1);
const pageSize = ref(10);

const paginatedIncidents = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return props.incidents.slice(start, end);
});

const totalPages = computed(() => {
  return Math.ceil(props.incidents.length / pageSize.value) || 1;
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
          <TableHead class="w-[120px]">Status</TableHead>
          <TableHead>Monitor Name</TableHead>
          <TableHead>Error Details</TableHead>
          <TableHead>Triggered At</TableHead>
          <TableHead>Resolved At</TableHead>
          <TableHead class="w-[120px] text-right">Duration</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="item in paginatedIncidents" :key="item.id">
          <!-- Status Badge -->
          <TableCell>
            <Badge
              :variant="item.status === 'resolved' ? 'outline' : 'destructive'"
              :class="[
                'text-[10px] font-bold py-0.5 px-2 uppercase',
                item.status === 'resolved'
                  ? 'bg-emerald-500/5 text-emerald-600 dark:text-emerald-400 border-emerald-500/20'
                  : 'bg-red-500/10 text-red-500 border-red-500/20'
              ]"
            >
              <span class="flex items-center gap-1">
                <span :class="[
                  'w-1.5 h-1.5 rounded-full',
                  item.status === 'resolved' ? 'bg-emerald-500' : 'bg-red-500 animate-pulse'
                ]"></span>
                <span>{{ item.status === 'resolved' ? 'resolved' : 'down' }}</span>
              </span>
            </Badge>
          </TableCell>

          <!-- Monitor Details -->
          <TableCell class="font-bold">
            <div class="flex flex-col gap-0.5 max-w-[180px]">
              <span class="truncate">{{ monitorsMap[item.monitor_id]?.name || `Monitor #${item.monitor_id}` }}</span>
              <span class="text-[9px] text-muted-foreground font-normal truncate">
                {{ monitorsMap[item.monitor_id]?.url || 'Target Endpoint' }}
              </span>
            </div>
          </TableCell>

          <!-- Error Message Details -->
          <TableCell>
            <div class="flex flex-col gap-0.5 max-w-[250px]">
              <span class="font-semibold truncate">{{ item.error_message || 'Connection failure' }}</span>
              <span class="text-[9px] text-muted-foreground inline-flex items-center gap-0.5">
                <Zap class="w-3 h-3 text-amber-500" />
                <span>Ping response: {{ item.latency > 0 ? `${item.latency}ms` : 'Timeout' }}</span>
              </span>
            </div>
          </TableCell>

          <!-- Triggered At -->
          <TableCell class="text-muted-foreground">
            {{ formatDate(item.created_at) }}
          </TableCell>

          <!-- Resolved At -->
          <TableCell class="text-muted-foreground">
            <span :class="item.status === 'resolved' ? 'text-muted-foreground' : 'text-amber-500 font-medium'">
              {{ item.status === 'resolved' ? formatDate(item.resolved_at) : 'Active Outage' }}
            </span>
          </TableCell>

          <!-- Duration -->
          <TableCell class="text-right font-semibold">
            <div class="inline-flex items-center gap-1">
              <Clock class="w-3.5 h-3.5 text-muted-foreground/60" />
              <span>{{ getDuration(item) }}</span>
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
        Showing {{ Math.min(props.incidents.length, (currentPage - 1) * pageSize + 1) }} to {{ Math.min(props.incidents.length, currentPage * pageSize) }} of {{ props.incidents.length }} entries
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
