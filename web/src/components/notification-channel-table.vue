<script setup lang="ts">
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Mail,
  MessageSquare,
  MessageCircle,
  Send,
  Webhook,
  Edit,
  Trash2
} from '@lucide/vue';
import type { NotificationChannel } from '@/stores/notifications';

defineProps<{
  channels: NotificationChannel[];
}>();

const emit = defineEmits<{
  (e: 'edit', channel: NotificationChannel): void;
  (e: 'delete', channel: NotificationChannel): void;
  (e: 'toggle', channel: NotificationChannel): void;
}>();

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
  <div class="overflow-x-auto">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead class="w-[80px] text-xs font-bold uppercase text-muted-foreground">Active</TableHead>
          <TableHead class="text-xs font-bold uppercase text-muted-foreground">Channel Name</TableHead>
          <TableHead class="text-xs font-bold uppercase text-muted-foreground">Type</TableHead>
          <TableHead class="text-xs font-bold uppercase text-muted-foreground">Destination/Target</TableHead>
          <TableHead class="w-[120px] text-right text-xs font-bold uppercase text-muted-foreground">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="item in channels" :key="item.id" class="channel-row hover:bg-muted/40 transition-colors">
          <!-- Toggle switch / status -->
          <TableCell class="py-4">
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
          <TableCell class="font-medium text-foreground py-4 text-xs font-bold">
            {{ item.name }}
          </TableCell>

          <!-- Type -->
          <TableCell class="py-4">
            <Badge variant="outline" class="uppercase text-[9px] font-bold py-0.5 px-2 bg-slate-50 dark:bg-slate-900 border-border/50 flex items-center gap-1.5 w-fit">
              <component :is="getTypeIcon(item.type)" class="w-3.5 h-3.5" />
              <span>{{ item.type }}</span>
            </Badge>
          </TableCell>

          <!-- Destination -->
          <TableCell class="py-4 text-xs text-muted-foreground font-mono">
            <span class="truncate block max-w-sm sm:max-w-md lg:max-w-lg">{{ getConfigSummary(item) }}</span>
          </TableCell>

          <!-- Actions -->
          <TableCell class="py-4 text-right">
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
  </div>
</template>
