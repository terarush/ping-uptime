<script setup lang="ts">
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Clock, Zap } from '@lucide/vue';
import type { Incident, MonitorItem } from '@/stores/incidents';

defineProps<{
  incidents: Incident[];
  monitorsMap: Record<number, MonitorItem>;
  formatDate: (dStr?: string) => string;
  getDuration: (inc: Incident) => string;
}>();
</script>

<template>
  <div class="overflow-x-auto">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead class="w-[120px] text-xs font-bold uppercase text-muted-foreground">Status</TableHead>
          <TableHead class="text-xs font-bold uppercase text-muted-foreground">Monitor Name</TableHead>
          <TableHead class="text-xs font-bold uppercase text-muted-foreground">Error Details</TableHead>
          <TableHead class="text-xs font-bold uppercase text-muted-foreground">Triggered At</TableHead>
          <TableHead class="text-xs font-bold uppercase text-muted-foreground">Resolved At</TableHead>
          <TableHead class="w-[120px] text-right text-xs font-bold uppercase text-muted-foreground">Duration</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="item in incidents" :key="item.id" class="incident-row hover:bg-muted/40 transition-colors">
          <!-- Status Badge -->
          <TableCell class="py-4">
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
          <TableCell class="font-bold text-foreground py-4 text-xs">
            <div class="flex flex-col gap-0.5 max-w-[180px]">
              <span class="truncate">{{ monitorsMap[item.monitor_id]?.name || `Monitor #${item.monitor_id}` }}</span>
              <span class="text-[9px] text-muted-foreground font-normal truncate">
                {{ monitorsMap[item.monitor_id]?.url || 'Target Endpoint' }}
              </span>
            </div>
          </TableCell>

          <!-- Error Message Details -->
          <TableCell class="py-4 text-xs">
            <div class="flex flex-col gap-0.5 max-w-[250px]">
              <span class="font-semibold text-foreground truncate">{{ item.error_message || 'Connection failure' }}</span>
              <span class="text-[9px] text-muted-foreground inline-flex items-center gap-0.5">
                <Zap class="w-3 h-3 text-amber-500" />
                <span>Ping response: {{ item.latency > 0 ? `${item.latency}ms` : 'Timeout' }}</span>
              </span>
            </div>
          </TableCell>

          <!-- Triggered At -->
          <TableCell class="py-4 text-xs text-muted-foreground">
            {{ formatDate(item.created_at) }}
          </TableCell>

          <!-- Resolved At -->
          <TableCell class="py-4 text-xs text-muted-foreground">
            <span :class="item.status === 'resolved' ? 'text-muted-foreground' : 'text-amber-500 font-medium'">
              {{ item.status === 'resolved' ? formatDate(item.resolved_at) : 'Active Outage' }}
            </span>
          </TableCell>

          <!-- Duration -->
          <TableCell class="py-4 text-right font-semibold text-xs text-foreground">
            <div class="inline-flex items-center gap-1">
              <Clock class="w-3.5 h-3.5 text-muted-foreground/60" />
              <span>{{ getDuration(item) }}</span>
            </div>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
  </div>
</template>
