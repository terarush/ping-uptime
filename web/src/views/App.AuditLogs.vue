<script setup lang="ts">
import { onMounted } from 'vue';
import { useAuditLogs } from '@/composables/useAuditLogs';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ClipboardList, RefreshCw, Loader2 } from '@lucide/vue';

const { items, loading, fetchAll } = useAuditLogs();
onMounted(fetchAll);

const fmtTime = (s: string) => new Date(s).toLocaleString();

const actionColor = (a: string) => {
  switch (a) {
    case 'created': return 'bg-emerald-500/10 text-emerald-500';
    case 'updated': return 'bg-blue-500/10 text-blue-500';
    case 'deleted': return 'bg-red-500/10 text-red-500';
    case 'resolved': return 'bg-emerald-500/10 text-emerald-500';
    default: return 'bg-slate-500/10 text-slate-500';
  }
};
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-indigo-500/10 dark:bg-indigo-500/5 blur-[100px] pointer-events-none"></div>

    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <ClipboardList class="w-6 h-6 text-primary" />
          <span>Audit Logs</span>
        </h2>
        <p class="text-xs text-muted-foreground">Track all entity changes in the system.</p>
      </div>
      <Button variant="outline" size="sm" @click="fetchAll" class="h-9">
        <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
        Refresh
      </Button>
    </div>

    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <CardTitle class="text-sm font-bold">Activity Log</CardTitle>
        <CardDescription class="text-xs">Last 200 events across all entities.</CardDescription>
      </CardHeader>
      <CardContent class="p-0">
        <div v-if="loading && items.length === 0" class="flex justify-center py-20">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
        </div>
        <div v-else-if="items.length === 0" class="flex flex-col items-center py-20 text-center">
          <ClipboardList class="w-12 h-12 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-bold text-foreground">No audit logs yet</p>
        </div>
        <Table v-else>
          <TableHeader>
            <TableRow>
              <TableHead>Time</TableHead>
              <TableHead>Action</TableHead>
              <TableHead>Entity</TableHead>
              <TableHead>Entity ID</TableHead>
              <TableHead>Details</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="item in items" :key="item.id" class="text-xs">
              <TableCell class="text-muted-foreground whitespace-nowrap">{{ fmtTime(item.created_at) }}</TableCell>
              <TableCell>
                <Badge variant="outline" :class="['uppercase text-[10px] font-bold', actionColor(item.action)]">
                  {{ item.action }}
                </Badge>
              </TableCell>
              <TableCell><code class="text-[10px] bg-muted px-1.5 py-0.5 rounded">{{ item.entity_type }}</code></TableCell>
              <TableCell>{{ item.entity_id }}</TableCell>
              <TableCell class="max-w-xs truncate text-muted-foreground">{{ item.details?.slice(0, 80) }}{{ item.details?.length > 80 ? '...' : '' }}</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </div>
</template>
