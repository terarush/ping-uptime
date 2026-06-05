<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import {
  Mail,
  MessageSquare,
  MessageCircle,
  Send,
  Webhook,
  Edit,
  Trash2,
  ChevronLeft,
  ChevronRight
} from '@lucide/vue';
import type { NotificationChannel } from '@/stores/notifications';

const props = defineProps<{
  channels: NotificationChannel[];
}>();

const emit = defineEmits<{
  (e: 'edit', channel: NotificationChannel): void;
  (e: 'delete', channel: NotificationChannel): void;
  (e: 'toggle', channel: NotificationChannel): void;
}>();

const currentPage = ref(1);
const pageSize = ref(10);

const paginatedChannels = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return props.channels.slice(start, end);
});

const totalPages = computed(() => {
  return Math.ceil(props.channels.length / pageSize.value) || 1;
});

watch(() => totalPages.value, () => {
  if (currentPage.value > totalPages.value) {
    currentPage.value = totalPages.value;
  }
});

// Map channel type to icon
const getTypeIcon = (type: string) => {
  switch (type) {
    case 'email': return Mail;
    case 'slack': return MessageSquare;
    case 'discord': return MessageCircle;
    case 'telegram': return Send;
    default: return Webhook;
  }
};

// Map config info to readable text
const getConfigSummary = (item: NotificationChannel) => {
  try {
    const parsed = JSON.parse(item.config);
    if (item.type === 'email') return parsed.email || '—';
    if (item.type === 'telegram') return `Chat ID: ${parsed.chat_id || '—'}`;
    return parsed.webhook_url || parsed.url || '—';
  } catch (e) {
    return 'Invalid configuration';
  }
};
</script>

<template>
  <div class="flex flex-col">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead class="w-[80px]">Active</TableHead>
          <TableHead>Channel Name</TableHead>
          <TableHead>Type</TableHead>
          <TableHead>Destination/Target</TableHead>
          <TableHead class="w-[120px] text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="item in paginatedChannels" :key="item.id">
          <!-- Toggle switch / status -->
          <TableCell>
            <div class="flex items-center">
              <button
                @click="emit('toggle', item)"
                class="relative inline-flex h-5 w-9 shrink-0 cursor-pointer items-center rounded-full transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring bg-input"
                :class="{ 'bg-primary': item.enabled }"
              >
                <span
                  class="pointer-events-none block h-4 w-4 rounded-full bg-background shadow-lg ring-0 transition-transform"
                  :class="item.enabled ? 'translate-x-4.5' : 'translate-x-0.5'"
                />
              </button>
            </div>
          </TableCell>

          <!-- Channel Name -->
          <TableCell class="font-medium font-bold">
            {{ item.name }}
          </TableCell>

          <!-- Type -->
          <TableCell>
            <Badge variant="outline" class="uppercase text-[9px] font-bold py-0.5 px-2 bg-slate-50 dark:bg-slate-900 border-border/50 flex items-center gap-1.5 w-fit">
              <component :is="getTypeIcon(item.type)" class="w-3.5 h-3.5" />
              <span>{{ item.type }}</span>
            </Badge>
          </TableCell>

          <!-- Destination -->
          <TableCell class="text-xs text-muted-foreground font-mono">
            <span class="truncate block max-w-sm sm:max-w-md lg:max-w-lg">{{ getConfigSummary(item) }}</span>
          </TableCell>

          <!-- Actions -->
          <TableCell class="text-right">
            <div class="flex items-center justify-end gap-1">
              <Button variant="ghost" size="icon" class="h-8 w-8 text-foreground hover:bg-muted" @click="emit('edit', item)">
                <Edit class="w-4 h-4" />
              </Button>
              <Button variant="ghost" size="icon" class="h-8 w-8 text-destructive hover:bg-destructive/10 hover:text-destructive" @click="emit('delete', item)">
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
        Showing {{ Math.min(props.channels.length, (currentPage - 1) * pageSize + 1) }} to {{ Math.min(props.channels.length, currentPage * pageSize) }} of {{ props.channels.length }} entries
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
