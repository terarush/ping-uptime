<script setup lang="ts">
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Edit, Trash2, Clock } from '@lucide/vue';
import type { Maintenance } from '@/stores/maintenances';

defineProps<{ items: Maintenance[] }>();
defineEmits<{
  (e: 'edit', item: Maintenance): void;
  (e: 'delete', item: Maintenance): void;
}>();

const statusClass = (status: string) => {
  switch (status) {
    case 'ongoing': return 'bg-amber-500/10 text-amber-500 border-amber-500/20';
    case 'completed': return 'bg-blue-500/10 text-blue-500 border-blue-500/20';
    default: return 'bg-slate-500/10 text-slate-500 border-slate-500/20';
  }
};

const fmtDate = (s: string) => {
  const d = new Date(s);
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' });
};
</script>

<template>
  <Table>
    <TableHeader>
      <TableRow>
        <TableHead>Name</TableHead>
        <TableHead>Start</TableHead>
        <TableHead>End</TableHead>
        <TableHead>Status</TableHead>
        <TableHead class="text-right">Actions</TableHead>
      </TableRow>
    </TableHeader>
    <TableBody>
      <TableRow v-for="item in items" :key="item.id">
        <TableCell class="font-medium">
          <div class="flex flex-col gap-0.5">
            <span class="text-xs font-bold">{{ item.name }}</span>
            <span class="text-[10px] text-muted-foreground truncate max-w-60">{{ item.description }}</span>
          </div>
        </TableCell>
        <TableCell>
          <span class="text-xs inline-flex items-center gap-1">
            <Clock class="w-3 h-3" /> {{ fmtDate(item.start_at) }}
          </span>
        </TableCell>
        <TableCell>
          <span class="text-xs inline-flex items-center gap-1">
            <Clock class="w-3 h-3" /> {{ fmtDate(item.end_at) }}
          </span>
        </TableCell>
        <TableCell>
          <Badge variant="outline" :class="['uppercase text-[10px] font-bold', statusClass(item.status)]">
            {{ item.status }}
          </Badge>
        </TableCell>
        <TableCell class="text-right">
          <div class="flex items-center justify-end gap-1">
            <Button variant="ghost" size="icon" class="h-8 w-8" @click="$emit('edit', item)">
              <Edit class="w-4 h-4" />
            </Button>
            <Button variant="ghost" size="icon" class="h-8 w-8 text-destructive" @click="$emit('delete', item)">
              <Trash2 class="w-4 h-4" />
            </Button>
          </div>
        </TableCell>
      </TableRow>
    </TableBody>
  </Table>
</template>
