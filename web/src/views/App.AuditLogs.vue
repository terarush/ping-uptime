<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useAuditLogs, type AuditLogFilter } from '@/composables/useAuditLogs';
import { useUsers } from '@/composables/useUsers';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { ClipboardList, RefreshCw, Loader2, Search } from '@lucide/vue';

const { items, loading, fetchAll } = useAuditLogs();
const { users, fetchUsersData } = useUsers();

// Filters
const filterAction = ref('');
const filterEntity = ref('');
const filterUserID = ref(0);
const filterFrom = ref('');
const filterTo = ref('');

const ENTITY_TYPES = ['monitor', 'incident', 'status_page', 'notification_channel'];
const ACTIONS = ['created', 'updated', 'deleted', 'resolved'];

// Build user map for name lookup
const userMap = computed(() => {
  const m: Record<number, string> = {};
  for (const u of users.value) {
    m[u.id] = u.name;
  }
  return m;
});

const applyFilters = async () => {
  const filter: AuditLogFilter = {};
  if (filterAction.value) filter.action = filterAction.value;
  if (filterEntity.value) filter.entity_type = filterEntity.value;
  if (filterUserID.value && filterUserID.value > 0) filter.user_id = filterUserID.value;
  if (filterFrom.value) filter.from = new Date(filterFrom.value).toISOString();
  if (filterTo.value) filter.to = new Date(filterTo.value).toISOString();
  await fetchAll(filter);
};

const clearFilters = async () => {
  filterAction.value = '';
  filterEntity.value = '';
  filterUserID.value = 0;
  filterFrom.value = '';
  filterTo.value = '';
  await fetchAll();
};

const fmtTime = (s: string) => new Date(s).toLocaleString();

const actionColor = (a: string) => {
  switch (a) {
    case 'created': return 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20';
    case 'updated': return 'bg-blue-500/10 text-blue-500 border-blue-500/20';
    case 'deleted': return 'bg-red-500/10 text-red-500 border-red-500/20';
    case 'resolved': return 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20';
    default: return 'bg-slate-500/10 text-slate-500 border-slate-500/20';
  }
};

onMounted(async () => {
  await Promise.all([fetchAll(), fetchUsersData()]);
});
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
      <Button variant="outline" size="sm" @click="applyFilters" class="h-9">
        <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
        Refresh
      </Button>
    </div>

    <!-- Filter Bar -->
    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardContent class="p-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-3 items-end">
          <div class="space-y-1.5">
            <label class="text-[10px] font-bold text-muted-foreground uppercase tracking-wider">Entity</label>
            <Select v-model="filterEntity">
              <SelectTrigger class="w-full h-9">
                <SelectValue placeholder="All entities" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="">All entities</SelectItem>
                <SelectItem v-for="et in ENTITY_TYPES" :key="et" :value="et">{{ et }}</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="space-y-1.5">
            <label class="text-[10px] font-bold text-muted-foreground uppercase tracking-wider">Action</label>
            <Select v-model="filterAction">
              <SelectTrigger class="w-full h-9">
                <SelectValue placeholder="All actions" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="">All actions</SelectItem>
                <SelectItem v-for="a in ACTIONS" :key="a" :value="a">{{ a }}</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="space-y-1.5">
            <label class="text-[10px] font-bold text-muted-foreground uppercase tracking-wider">User</label>
            <Select v-model="filterUserID">
              <SelectTrigger class="w-full h-9">
                <SelectValue placeholder="All users" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem :value="0">All users</SelectItem>
                <SelectItem v-for="u in users" :key="u.id" :value="u.id">{{ u.name }}</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="space-y-1.5">
            <label class="text-[10px] font-bold text-muted-foreground uppercase tracking-wider">From</label>
            <Input v-model="filterFrom" type="datetime-local" class="w-full h-9" />
          </div>
          <div class="space-y-1.5">
            <label class="text-[10px] font-bold text-muted-foreground uppercase tracking-wider">To</label>
            <Input v-model="filterTo" type="datetime-local" class="w-full h-9" />
          </div>
          <div class="flex items-end gap-2 sm:col-span-2 lg:col-span-3 xl:col-span-5">
            <Button size="sm" @click="applyFilters" class="h-9">Apply</Button>
            <Button variant="outline" size="sm" @click="clearFilters" class="h-9">Clear</Button>
          </div>
        </div>
      </CardContent>
    </Card>

    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <CardTitle class="text-sm font-bold">Activity Log</CardTitle>
        <CardDescription class="text-xs">{{ items.length }} event(s) recorded.</CardDescription>
      </CardHeader>
      <CardContent class="p-0">
        <div v-if="loading && items.length === 0" class="flex justify-center py-20">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
        </div>
        <div v-else-if="items.length === 0" class="flex flex-col items-center py-20 text-center">
          <ClipboardList class="w-12 h-12 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-bold text-foreground">No audit logs found</p>
          <p class="text-xs text-muted-foreground mt-1">Try adjusting filters or wait for events.</p>
        </div>
        <Table v-else>
          <TableHeader>
            <TableRow>
              <TableHead>Time</TableHead>
              <TableHead>User</TableHead>
              <TableHead>Action</TableHead>
              <TableHead>Entity</TableHead>
              <TableHead>ID</TableHead>
              <TableHead>Details</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="item in items" :key="item.id" class="text-xs">
              <TableCell class="text-muted-foreground whitespace-nowrap">{{ fmtTime(item.created_at) }}</TableCell>
              <TableCell class="text-muted-foreground whitespace-nowrap">
                {{ userMap[item.user_id] || `User #${item.user_id}` }}
              </TableCell>
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
