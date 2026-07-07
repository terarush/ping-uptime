<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useTags } from '@/composables/useTags';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import {
  Tag,
  Search,
  Loader2,
  RefreshCw,
  Plus,
  Trash2,
  Pencil,
  X,
  Check,
} from '@lucide/vue';
import gsap from 'gsap';
import { toast } from 'vue-sonner';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';

const { tags, loading, fetchAll, create, update, deleteTag } = useTags();

const searchQuery = ref('');
const showDialog = ref(false);
const editingTag = ref<{ id: number; name: string; color: string } | null>(null);
const formName = ref('');
const formColor = ref('#6366f1');
const saving = ref(false);
const showDeleteConfirm = ref(false);
const deletingTagID = ref<number | null>(null);

const COLORS = [
  '#6366f1', // indigo
  '#ef4444', // red
  '#22c55e', // green
  '#eab308', // yellow
  '#f97316', // orange
  '#06b6d4', // cyan
  '#a855f7', // purple
  '#ec4899', // pink
];

const filteredTags = computed(() => {
  if (!tags.value) return [];
  const query = searchQuery.value.toLowerCase().trim();
  if (!query) return tags.value;
  return tags.value.filter(t => t.name.toLowerCase().includes(query));
});

const animateRows = () => {
  setTimeout(() => {
    gsap.fromTo('.tag-row',
      { opacity: 0, y: 12 },
      { opacity: 1, y: 0, duration: 0.3, stagger: 0.03, ease: 'power2.out' }
    );
  }, 50);
};

const loadTags = async () => {
  try {
    await fetchAll();
    animateRows();
  } catch {
    // handled in store
  }
};

const openCreate = () => {
  editingTag.value = null;
  formName.value = '';
  formColor.value = '#6366f1';
  showDialog.value = true;
};

const openEdit = (tag: { id: number; name: string; color: string }) => {
  editingTag.value = tag;
  formName.value = tag.name;
  formColor.value = tag.color;
  showDialog.value = true;
};

const closeDialog = () => {
  showDialog.value = false;
  editingTag.value = null;
};

const saveTag = async () => {
  if (!formName.value.trim()) {
    toast.error('Tag name cannot be empty');
    return;
  }
  saving.value = true;
  try {
    if (editingTag.value) {
      await update(editingTag.value.id, { name: formName.value.trim(), color: formColor.value });
      toast.success('Tag updated');
    } else {
      await create({ name: formName.value.trim(), color: formColor.value });
      toast.success('Tag created');
    }
    closeDialog();
    animateRows();
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to save tag');
  } finally {
    saving.value = false;
  }
};

const confirmDelete = (id: number) => {
  deletingTagID.value = id;
  showDeleteConfirm.value = true;
};

const handleDelete = async () => {
  if (deletingTagID.value === null) return;
  try {
    await deleteTag(deletingTagID.value);
    toast.success('Tag deleted');
    showDeleteConfirm.value = false;
    deletingTagID.value = null;
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to delete tag');
  }
};

onMounted(async () => {
  await loadTags();
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-indigo-500/5 dark:bg-indigo-500/2 blur-[100px] pointer-events-none"></div>

    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <Tag class="w-6 h-6 text-primary" />
          <span>Tags &amp; Labels</span>
        </h2>
        <p class="text-xs text-muted-foreground">Organise monitors with tags and coloured labels.</p>
      </div>

      <div class="flex items-center gap-3">
        <Button variant="outline" size="sm" @click="loadTags" class="h-9">
          <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
          <span>Refresh</span>
        </Button>
        <Button size="sm" @click="openCreate" class="h-9">
          <Plus class="w-4 h-4 mr-1.5" />
          <span>New Tag</span>
        </Button>
      </div>
    </div>

    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <CardTitle class="text-sm font-bold text-foreground flex items-center gap-2">
              <Tag class="w-4 h-4" />
              <span>{{ tags.length }} tags</span>
            </CardTitle>
            <CardDescription class="text-xs">All tags in your account.</CardDescription>
          </div>

          <div class="relative w-full sm:w-72">
            <Search class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              v-model="searchQuery"
              placeholder="Search tags..."
              class="pl-9 h-9"
            />
          </div>
        </div>
      </CardHeader>

      <CardContent class="p-0">
        <div v-if="loading && tags.length === 0" class="flex flex-col items-center justify-center py-20 gap-3">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
          <p class="text-sm text-muted-foreground">Loading tags...</p>
        </div>

        <div v-else-if="filteredTags.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
          <Tag class="w-12 h-12 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-bold text-foreground">No tags found</p>
          <p class="text-xs text-muted-foreground mt-1">{{ searchQuery ? 'No tags match your search.' : 'Create your first tag to organise monitors.' }}</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b border-border/40 text-xs text-muted-foreground">
                <th class="text-left font-medium px-4 py-3">Name</th>
                <th class="text-left font-medium px-4 py-3">Colour</th>
                <th class="text-right font-medium px-4 py-3">Created</th>
                <th class="text-right font-medium px-4 py-3 w-24">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="t in filteredTags"
                :key="t.id"
                class="tag-row border-b border-border/20 hover:bg-accent/30 transition-colors"
              >
                <td class="px-4 py-3 text-sm font-medium">{{ t.name }}</td>
                <td class="px-4 py-3">
                  <div class="flex items-center gap-2">
                    <div
                      class="w-5 h-5 rounded-full border border-border/40 shrink-0"
                      :style="{ backgroundColor: t.color }"
                    ></div>
                    <span class="text-xs text-muted-foreground font-mono">{{ t.color }}</span>
                  </div>
                </td>
                <td class="px-4 py-3 text-sm text-right text-muted-foreground">
                  {{ new Date(t.created_at).toLocaleDateString() }}
                </td>
                <td class="px-4 py-3 text-right">
                  <div class="flex items-center justify-end gap-1">
                    <Button variant="ghost" size="icon" class="h-8 w-8" @click="openEdit(t)">
                      <Pencil class="w-4 h-4" />
                    </Button>
                    <Button variant="ghost" size="icon" class="h-8 w-8 text-red-500 hover:text-red-600 hover:bg-red-500/10" @click="confirmDelete(t.id)">
                      <Trash2 class="w-4 h-4" />
                    </Button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>

    <!-- Create / Edit Dialog -->
    <Dialog :open="showDialog" @update:open="showDialog = $event">
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{{ editingTag ? 'Edit Tag' : 'Create Tag' }}</DialogTitle>
          <DialogDescription>
            {{ editingTag ? 'Update the tag name or colour.' : 'Add a new tag to organise monitors.' }}
          </DialogDescription>
        </DialogHeader>

        <form @submit.prevent="saveTag" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="tagname">Name</Label>
            <Input id="tagname" v-model="formName" placeholder="e.g. production, staging, critical" class="h-9" />
          </div>

          <div class="space-y-2">
            <Label for="tagcolor">Colour</Label>
            <div class="flex items-center gap-3">
              <div
                class="w-8 h-8 rounded-full border-2 shrink-0"
                :style="{ backgroundColor: formColor, borderColor: formColor }"
              ></div>
              <Input id="tagcolor" v-model="formColor" placeholder="#6366f1" class="h-9 font-mono flex-1" />
            </div>
            <div class="flex flex-wrap gap-2 mt-3">
              <button
                v-for="c in COLORS"
                :key="c"
                type="button"
                class="w-8 h-8 rounded-full border-2 transition-all"
                :class="formColor === c ? 'border-foreground scale-110' : 'border-border/40 hover:scale-110'"
                :style="{ backgroundColor: c }"
                @click="formColor = c"
              ></button>
            </div>
          </div>

          <DialogFooter class="pt-4 border-t border-border/40 mt-4">
            <Button type="button" variant="outline" @click="closeDialog" class="h-9">Cancel</Button>
            <Button type="submit" :disabled="saving" class="h-9">{{ saving ? 'Saving...' : 'Save' }}</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>

    <!-- Delete Confirmation Dialog -->
    <Dialog :open="showDeleteConfirm" @update:open="showDeleteConfirm = $event">
      <DialogContent class="sm:max-w-[400px]">
        <DialogHeader>
          <DialogTitle>Delete Tag</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this tag? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" size="sm" @click="showDeleteConfirm = false" class="h-9">
            <X class="w-4 h-4 mr-1.5" />
            <span>Cancel</span>
          </Button>
          <Button variant="destructive" size="sm" @click="handleDelete" class="h-9">
            <Trash2 class="w-4 h-4 mr-1.5" />
            <span>Delete</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
