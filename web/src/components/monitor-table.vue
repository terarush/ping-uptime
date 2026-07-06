<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { Clock, Activity, Edit, Trash2, ExternalLink, Shield, ChevronLeft, ChevronRight } from '@lucide/vue';
import type { Monitor } from '@/stores/monitors';

const props = defineProps<{
  monitors: Monitor[];
  selectedId?: number | null;
}>();

defineEmits<{
  (e: 'select', id: number): void;
  (e: 'edit', monitor: Monitor): void;
  (e: 'delete', monitor: Monitor): void;
}>();

const currentPage = ref(1);
const pageSize = ref(10);

const paginatedMonitors = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return props.monitors.slice(start, end);
});

const totalPages = computed(() => {
  return Math.ceil(props.monitors.length / pageSize.value) || 1;
});

watch(() => totalPages.value, () => {
  if (currentPage.value > totalPages.value) {
    currentPage.value = totalPages.value;
  }
});

const sslStatusClass = (expiresAt: string) => {
  const days = daysUntilExpiryNum(expiresAt);
  if (days === null) return 'text-muted-foreground';
  if (days <= 7) return 'text-red-500 font-bold';
  if (days <= 30) return 'text-amber-500';
  return 'text-emerald-500';
};

const daysUntilExpiryNum = (expiresAt: string): number | null => {
  const expiry = new Date(expiresAt);
  if (isNaN(expiry.getTime())) return null;
  const diff = expiry.getTime() - Date.now();
  return Math.ceil(diff / (1000 * 60 * 60 * 24));
};

const daysUntilExpiry = (expiresAt: string): string => {
  const days = daysUntilExpiryNum(expiresAt);
  if (days === null) return '—';
  if (days < 0) return 'Expired';
  if (days === 0) return 'Today';
  if (days === 1) return '1 day';
  return `${days}d`;
};
</script>

<template>
  <div class="flex flex-col">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead class="w-[80px] text-center">Status</TableHead>
          <TableHead>Name / URL</TableHead>
          <TableHead>Type</TableHead>
          <TableHead>Tags</TableHead>
          <TableHead class="w-[110px]">SSL</TableHead>
          <TableHead>Interval / Timeout</TableHead>
          <TableHead class="w-[120px] text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow
          v-for="item in paginatedMonitors"
          :key="item.id"
          @click="$emit('select', item.id)"
          class="cursor-pointer"
          :data-state="item.id === selectedId ? 'selected' : undefined"
        >
          <!-- Status -->
          <TableCell>
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
          <TableCell class="font-medium">
            <div class="flex flex-col gap-0.5 max-w-sm sm:max-w-md md:max-w-lg lg:max-w-xl">
              <span class="text-xs font-bold truncate">{{ item.name }}</span>
              <a :href="item.url" target="_blank" @click.stop class="text-[10px] text-muted-foreground hover:underline inline-flex items-center gap-1">
                <span class="truncate">{{ item.url }}</span>
                <ExternalLink class="w-2.5 h-2.5 shrink-0" />
              </a>
            </div>
          </TableCell>

          <!-- Type -->
          <TableCell>
            <Badge variant="outline" class="uppercase text-[10px] font-bold py-0.5 px-2 bg-slate-50 dark:bg-slate-900 border-border/50">
              {{ item.type }}
            </Badge>
          </TableCell>

          <!-- Tags -->
          <TableCell>
            <div class="flex flex-wrap gap-1" v-if="item.tags && item.tags.length">
              <span
                v-for="t in item.tags"
                :key="t.id"
                class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded-full text-[9px] font-medium"
                :style="{ backgroundColor: t.color + '20', color: t.color, borderColor: t.color + '40' }"
                style="border-width: 1px"
              >
                <span class="w-1.5 h-1.5 rounded-full shrink-0" :style="{ backgroundColor: t.color }"></span>
                {{ t.name }}
              </span>
            </div>
            <span v-else class="text-[10px] text-muted-foreground">—</span>
          </TableCell>

          <!-- SSL -->
          <TableCell>
            <div v-if="item.check_ssl" class="flex items-center gap-1.5">
              <Shield class="w-3.5 h-3.5 text-emerald-500" />
              <span v-if="item.ssl_expires_at" class="text-[10px] font-semibold whitespace-nowrap"
                :class="sslStatusClass(item.ssl_expires_at)">
                {{ daysUntilExpiry(item.ssl_expires_at) }}
              </span>
              <span v-else class="text-[10px] text-muted-foreground">Pending</span>
            </div>
            <span v-else class="text-[10px] text-muted-foreground">—</span>
          </TableCell>

          <!-- Interval/Timeout -->
          <TableCell>
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
          <TableCell class="text-right">
            <div class="flex items-center justify-end gap-1">
              <Button variant="ghost" size="icon" class="h-8 w-8 text-foreground hover:bg-muted" @click.stop="$emit('edit', item)">
                <Edit class="w-4 h-4" />
              </Button>
              <Button variant="ghost" size="icon" class="h-8 w-8 text-destructive hover:bg-destructive/10 hover:text-destructive" @click.stop="$emit('delete', item)">
                <Trash2 class="w-4 h-4" />
              </Button>
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
        Showing {{ Math.min(props.monitors.length, (currentPage - 1) * pageSize + 1) }} to {{ Math.min(props.monitors.length, currentPage * pageSize) }} of {{ props.monitors.length }} entries
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
