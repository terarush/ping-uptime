<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { Edit, Trash2, ExternalLink, ChevronLeft, ChevronRight } from '@lucide/vue';
import type { StatusPage } from '@/stores/statusPages';

const props = defineProps<{
  pages: StatusPage[];
}>();

const emit = defineEmits<{
  (e: 'edit', page: StatusPage): void;
  (e: 'delete', page: StatusPage): void;
}>();

const currentPage = ref(1);
const pageSize = ref(10);

const paginatedPages = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return props.pages.slice(start, end);
});

const totalPages = computed(() => {
  return Math.ceil(props.pages.length / pageSize.value) || 1;
});

watch(() => totalPages.value, () => {
  if (currentPage.value > totalPages.value) {
    currentPage.value = totalPages.value;
  }
});

// Generate public link
const getPublicUrl = (slug: string) => {
  return `${window.location.origin}/status/${slug}`;
};
</script>

<template>
  <div class="flex flex-col">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Page Name / Slug</TableHead>
          <TableHead>Description</TableHead>
          <TableHead>Linked Monitors</TableHead>
          <TableHead class="w-[120px] text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="item in paginatedPages" :key="item.id">
          <!-- Name & Slug -->
          <TableCell class="font-medium">
            <div class="flex flex-col gap-0.5 max-w-sm">
              <span class="text-xs font-bold">{{ item.name }}</span>
              <a :href="getPublicUrl(item.slug)" target="_blank" class="text-[10px] text-primary hover:underline inline-flex items-center gap-1">
                <span>/status/{{ item.slug }}</span>
                <ExternalLink class="w-2.5 h-2.5" />
              </a>
            </div>
          </TableCell>

          <!-- Description -->
          <TableCell class="text-xs text-muted-foreground">
            <span class="line-clamp-2 max-w-md">{{ item.description || 'No description provided.' }}</span>
          </TableCell>

          <!-- Linked Monitors -->
          <TableCell class="text-xs">
            <Badge variant="secondary" class="font-bold py-0.5 px-2 bg-slate-100 dark:bg-slate-900 border-border/50 text-[10px]">
              {{ item.monitors?.length || 0 }} Monitors
            </Badge>
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
        Showing {{ Math.min(props.pages.length, (currentPage - 1) * pageSize + 1) }} to {{ Math.min(props.pages.length, currentPage * pageSize) }} of {{ props.pages.length }} entries
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
