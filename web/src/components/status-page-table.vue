<script setup lang="ts">
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Edit, Trash2, ExternalLink } from '@lucide/vue';
import type { StatusPage } from '@/stores/statusPages';

defineProps<{
  pages: StatusPage[];
}>();

const emit = defineEmits<{
  (e: 'edit', page: StatusPage): void;
  (e: 'delete', page: StatusPage): void;
}>();

// Generate public link
const getPublicUrl = (slug: string) => {
  return `${window.location.origin}/status/${slug}`;
};
</script>

<template>
  <div class="overflow-x-auto">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead class="text-xs font-bold uppercase text-muted-foreground">Page Name / Slug</TableHead>
          <TableHead class="text-xs font-bold uppercase text-muted-foreground">Description</TableHead>
          <TableHead class="text-xs font-bold uppercase text-muted-foreground">Linked Monitors</TableHead>
          <TableHead class="w-[120px] text-right text-xs font-bold uppercase text-muted-foreground">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="item in pages" :key="item.id" class="page-row hover:bg-muted/40 transition-colors">
          <!-- Name & Slug -->
          <TableCell class="font-medium text-foreground py-4">
            <div class="flex flex-col gap-0.5 max-w-sm">
              <span class="text-xs font-bold">{{ item.name }}</span>
              <a :href="getPublicUrl(item.slug)" target="_blank" class="text-[10px] text-primary hover:underline inline-flex items-center gap-1">
                <span>/status/{{ item.slug }}</span>
                <ExternalLink class="w-2.5 h-2.5" />
              </a>
            </div>
          </TableCell>

          <!-- Description -->
          <TableCell class="py-4 text-xs text-muted-foreground">
            <span class="line-clamp-2 max-w-md">{{ item.description || 'No description provided.' }}</span>
          </TableCell>

          <!-- Linked Monitors -->
          <TableCell class="py-4 text-xs">
            <Badge variant="secondary" class="font-bold py-0.5 px-2 bg-slate-100 dark:bg-slate-900 border-border/50 text-[10px]">
              {{ item.monitors?.length || 0 }} Monitors
            </Badge>
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
