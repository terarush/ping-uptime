<script setup lang="ts">
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Clock, Activity, Edit, Trash2, ExternalLink } from '@lucide/vue';
import type { Monitor } from '@/stores/monitors';

defineProps<{
  monitors: Monitor[];
}>();

defineEmits<{
  (e: 'edit', monitor: Monitor): void;
  (e: 'delete', monitor: Monitor): void;
}>();
</script>

<template>
  <div class="overflow-x-auto">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead class="w-[80px] text-center text-xs font-bold uppercase text-muted-foreground">Status</TableHead>
          <TableHead class="text-xs font-bold uppercase text-muted-foreground">Name / URL</TableHead>
          <TableHead class="text-xs font-bold uppercase text-muted-foreground">Type</TableHead>
          <TableHead class="text-xs font-bold uppercase text-muted-foreground">Interval / Timeout</TableHead>
          <TableHead class="w-[120px] text-right text-xs font-bold uppercase text-muted-foreground">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="item in monitors" :key="item.id" class="monitor-row hover:bg-muted/40 transition-colors">
          <!-- Status -->
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

          <!-- Name & URL -->
          <TableCell class="font-medium text-foreground py-4">
            <div class="flex flex-col gap-0.5 max-w-sm sm:max-w-md md:max-w-lg lg:max-w-xl">
              <span class="text-xs font-bold truncate">{{ item.name }}</span>
              <a :href="item.url" target="_blank" class="text-[10px] text-muted-foreground hover:underline inline-flex items-center gap-1">
                <span class="truncate">{{ item.url }}</span>
                <ExternalLink class="w-2.5 h-2.5 shrink-0" />
              </a>
            </div>
          </TableCell>

          <!-- Type -->
          <TableCell class="py-4">
            <Badge variant="outline" class="uppercase text-[10px] font-bold py-0.5 px-2 bg-slate-50 dark:bg-slate-900 border-border/50">
              {{ item.type }}
            </Badge>
          </TableCell>

          <!-- Interval/Timeout -->
          <TableCell class="py-4 text-xs">
            <div class="flex items-center gap-3 text-muted-foreground">
              <span class="inline-flex items-center gap-1">
                <Clock class="w-3.5 h-3.5" />
                <span>{{ item.interval }}s</span>
              </span>
              <span class="inline-flex items-center gap-1">
                <Activity class="w-3.5 h-3.5" />
                <span>{{ item.timeout }}s</span>
              </span>
            </div>
          </TableCell>

          <!-- Actions -->
          <TableCell class="py-4 text-right">
            <div class="flex items-center justify-end gap-1">
              <Button variant="ghost" size="icon" class="h-8 w-8 text-foreground hover:bg-muted" @click="$emit('edit', item)">
                <Edit class="w-4 h-4" />
              </Button>
              <Button variant="ghost" size="icon" class="h-8 w-8 text-destructive hover:bg-destructive/10 hover:text-destructive" @click="$emit('delete', item)">
                <Trash2 class="w-4 h-4" />
              </Button>
            </div>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
  </div>
</template>
