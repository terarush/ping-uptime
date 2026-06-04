<script setup lang="ts">
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar';
import { Shield, Mail, Edit, Trash2 } from '@lucide/vue';

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  is_blocked: boolean;
  avatar?: string;
  created_at: string;
}

defineProps<{
  users: User[];
  currentUser: { id: number } | null;
}>();

defineEmits<{
  (e: 'edit', user: User): void;
  (e: 'delete', user: User): void;
}>();
</script>

<template>
  <Table>
    <TableHeader>
      <TableRow>
        <TableHead class="w-14"></TableHead>
        <TableHead class="text-xs font-bold uppercase text-muted-foreground">User</TableHead>
        <TableHead class="text-xs font-bold uppercase text-muted-foreground">Role</TableHead>
        <TableHead class="text-xs font-bold uppercase text-muted-foreground">Status</TableHead>
        <TableHead class="text-xs font-bold uppercase text-muted-foreground text-right pr-6">Actions</TableHead>
      </TableRow>
    </TableHeader>
    <TableBody>
      <TableRow v-for="user in users" :key="user.id" class="user-row hover:bg-muted/40 transition-colors">
        <!-- Avatar -->
        <TableCell class="py-4 pl-6">
          <Avatar class="h-9 w-9 ring-1 ring-primary/10">
            <AvatarImage :src="user.avatar || `https://api.dicebear.com/7.x/initials/svg?seed=${user.name}`" alt="Avatar" />
            <AvatarFallback class="font-bold text-xs uppercase">{{ user.name.slice(0, 2) }}</AvatarFallback>
          </Avatar>
        </TableCell>

        <!-- Profile Details -->
        <TableCell class="py-4">
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
        <TableCell class="py-4">
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
        <TableCell class="py-4">
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
        <TableCell class="py-4 text-right pr-6">
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
</template>
