<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar';
import { Shield, Mail, Edit, Trash2, ChevronLeft, ChevronRight } from '@lucide/vue';

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  is_blocked: boolean;
  avatar?: string;
  created_at: string;
}

const props = defineProps<{
  users: User[];
  currentUser: { id: number } | null;
}>();

defineEmits<{
  (e: 'edit', user: User): void;
  (e: 'delete', user: User): void;
}>();

const currentPage = ref(1);
const pageSize = ref(10);

const paginatedUsers = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return props.users.slice(start, end);
});

const totalPages = computed(() => {
  return Math.ceil(props.users.length / pageSize.value) || 1;
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
          <TableHead class="w-14"></TableHead>
          <TableHead>User</TableHead>
          <TableHead>Role</TableHead>
          <TableHead>Status</TableHead>
          <TableHead class="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="user in paginatedUsers" :key="user.id">
          <!-- Avatar -->
          <TableCell>
            <Avatar class="h-9 w-9 ring-1 ring-primary/10">
              <AvatarImage :src="user.avatar || `https://api.dicebear.com/7.x/initials/svg?seed=${user.name}`" alt="Avatar" />
              <AvatarFallback class="font-bold text-xs uppercase">{{ user.name.slice(0, 2) }}</AvatarFallback>
            </Avatar>
          </TableCell>

          <!-- Profile Details -->
          <TableCell>
            <div class="flex flex-col gap-0.5">
              <span class="text-xs font-bold text-foreground flex items-center gap-1.5">
                {{ user.name }}
                <Badge v-if="currentUser && currentUser.id === user.id" variant="secondary" class="text-[9px] font-bold px-1 py-0 h-4 bg-muted text-muted-foreground">You</Badge>
              </span>
              <span class="text-[10px] text-muted-foreground flex items-center gap-1">
                <Mail class="w-2.5 h-2.5" />
                {{ user.email }}
              </span>
            </div>
          </TableCell>

          <!-- Role Badge -->
          <TableCell>
            <Badge variant="outline" :class="[
              'text-[10px] font-bold py-0.5 px-2 border/50',
              user.role === 'admin'
                ? 'bg-primary/5 text-primary border-primary/20'
                : 'bg-muted text-muted-foreground border-border'
            ]">
              <Shield v-if="user.role === 'admin'" class="w-2.5 h-2.5 inline mr-1" />
              {{ user.role === 'admin' ? 'Administrator' : 'Standard User' }}
            </Badge>
          </TableCell>

          <!-- Status Badge -->
          <TableCell>
            <Badge variant="outline" :class="[
              'text-[10px] font-bold py-0.5 px-2 border/50',
              user.is_blocked
                ? 'bg-destructive/5 text-destructive border-destructive/20'
                : 'bg-emerald-500/5 text-emerald-600 dark:text-emerald-400 border-emerald-500/20'
            ]">
              {{ user.is_blocked ? 'Blocked' : 'Active' }}
            </Badge>
          </TableCell>

          <!-- Action Menu -->
          <TableCell class="text-right">
            <div class="flex items-center justify-end gap-1.5">
              <!-- Edit Action -->
              <Button @click="$emit('edit', user)" size="icon" variant="ghost" class="h-8 w-8 text-muted-foreground hover:text-foreground">
                <Edit class="w-4 h-4" />
              </Button>

              <!-- Delete Action (Disabled for logged-in user itself) -->
              <Button
                @click="$emit('delete', user)"
                size="icon"
                variant="ghost"
                class="h-8 w-8 text-muted-foreground hover:text-destructive"
                :disabled="currentUser && currentUser.id === user.id"
              >
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
        Showing {{ Math.min(props.users.length, (currentPage - 1) * pageSize + 1) }} to {{ Math.min(props.users.length, currentPage * pageSize) }} of {{ props.users.length }} entries
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
