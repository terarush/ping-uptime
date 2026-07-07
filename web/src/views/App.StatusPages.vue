<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useStatusPages, type StatusPage, type Monitor } from '@/composables/useStatusPages';
import { useAnalytics } from '@/composables/useAnalytics';
import { statusPageSchema } from '@/validations/status-page';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import StatusPageTable from '@/components/status-page-table.vue';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import {
  ShieldCheck,
  Plus,
  Trash2,
  Edit,
  ExternalLink,
  Search,
  Loader2,
  RefreshCw,
  Layout,
  Link2
} from '@lucide/vue';
import gsap from 'gsap';

// Use external status pages composable
const {
  statusPages,
  availableMonitors,
  loading,
  error,
  fetchStatusPagesData,
  createStatusPage,
  updateStatusPage,
  deleteStatusPage
} = useStatusPages();

// Analytics for real status data
const {
  stats: monitorStats,
  loading: analyticsLoading,
  fetchDashboardStats,
} = useAnalytics();

const searchQuery = ref('');
const success = ref('');

// Dialog states
const isFormDialogOpen = ref(false);
const isDeleteDialogOpen = ref(false);
const formLoading = ref(false);
const actionPage = ref<StatusPage | null>(null);

// Form fields
const formName = ref('');
const formSlug = ref('');
const formDescription = ref('');
const formSelectedMonitors = ref<number[]>([]);

const isEditMode = computed(() => !!actionPage.value);

// Merge analytics stats into status pages monitors
const pagesWithStats = computed(() => {
  return statusPages.value.map(page => {
    const monitors = page.monitors?.map(m => {
      const stats = monitorStats.value?.find(s => s.monitor_id === m.id);
      return stats
        ? { ...m, uptime_pct: stats.uptime_pct, status_label: stats.status, avg_latency: stats.avg_latency }
        : { ...m, uptime_pct: null, status_label: null, avg_latency: null };
    }) || [];
    return { ...page, monitors };
  });
});

// Fetch wrapper with animation callback and analytics
const fetchAll = async () => {
  try {
    await Promise.all([
      fetchStatusPagesData(),
      fetchDashboardStats('1m'),
    ]);
  } catch (err) {
    console.error('Failed to load status pages list:', err);
  } finally {
    setTimeout(animateTableRows, 50);
  }
};

// Filtered status pages list — use pagesWithStats so search + real data work together
const filteredPages = computed(() => {
  const pages = pagesWithStats.value;
  if (!pages) return [];
  const query = searchQuery.value.toLowerCase().trim();
  if (!query) return pages;
  return pages.filter(p =>
    p.name.toLowerCase().includes(query) ||
    p.slug.toLowerCase().includes(query) ||
    p.description.toLowerCase().includes(query)
  );
});

// Reset form fields
const resetForm = () => {
  formName.value = '';
  formSlug.value = '';
  formDescription.value = '';
  formSelectedMonitors.value = [];
  actionPage.value = null;
  error.value = '';
  success.value = '';
};

const openCreateDialog = async () => {
  resetForm();
  // Refetch monitors in case the initial load failed
  if (availableMonitors.value.length === 0) {
    try {
      const { data } = await (await import('@/lib/fetch')).default.get('/monitors');
      availableMonitors.value = data?.data || data || [];
    } catch { /* ignore */ }
  }
  isFormDialogOpen.value = true;
};

const openEditDialog = (page: StatusPage) => {
  resetForm();
  actionPage.value = page;
  formName.value = page.name;
  formSlug.value = page.slug;
  formDescription.value = page.description;
  formSelectedMonitors.value = page.monitors?.map(m => m.id) || [];
  isFormDialogOpen.value = true;
};

const openDeleteDialog = (page: StatusPage) => {
  actionPage.value = page;
  isDeleteDialogOpen.value = true;
};

// Submit form (Create / Update)
const handleFormSubmit = async () => {
  error.value = '';
  success.value = '';
  formLoading.value = true;

  const rawPayload = {
    name: formName.value,
    slug: formSlug.value,
    description: formDescription.value,
    monitor_ids: formSelectedMonitors.value,
  };

  // Perform client-side Zod validation
  const validation = statusPageSchema.safeParse(rawPayload);
  if (!validation.success) {
    const firstError = validation.error.errors[0]?.message || 'Validation failed';
    error.value = firstError;
    formLoading.value = false;
    return;
  }

  try {
    if (isEditMode.value && actionPage.value) {
      await updateStatusPage(actionPage.value.id, rawPayload);
      success.value = `Status page "${formName.value}" updated successfully!`;
    } else {
      await createStatusPage(rawPayload);
      success.value = `Status page "${formName.value}" created successfully!`;
    }

    isFormDialogOpen.value = false;
    await fetchAll();
  } catch (err: any) {
    console.error('Failed to save status page:', err);
    error.value = err.response?.data?.error || 'Failed to save status page data.';
  } finally {
    formLoading.value = false;
  }
};

// Delete status page
const handleDeleteConfirm = async () => {
  if (!actionPage.value) return;
  error.value = '';
  formLoading.value = true;

  try {
    await deleteStatusPage(actionPage.value.id);
    success.value = `Status page "${actionPage.value.name}" deleted successfully!`;
    isDeleteDialogOpen.value = false;
    await fetchAll();
  } catch (err: any) {
    console.error('Failed to delete status page:', err);
    error.value = err.response?.data?.error || 'Failed to delete status page.';
  } finally {
    formLoading.value = false;
  }
};

// Toggle monitor checkbox selection
const toggleMonitorSelection = (id: number) => {
  const index = formSelectedMonitors.value.indexOf(id);
  if (index > -1) {
    formSelectedMonitors.value.splice(index, 1);
  } else {
    formSelectedMonitors.value.push(id);
  }
};

// Generate public link
const getPublicUrl = (slug: string) => {
  return `${window.location.origin}/status/${slug}`;
};

// GSAP Animations
const animateTableRows = () => {
  gsap.fromTo('.page-row',
    { opacity: 0, y: 15 },
    { opacity: 1, y: 0, duration: 0.4, stagger: 0.05, ease: 'power2.out' }
  );
};

onMounted(() => {
  fetchAll();
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <!-- Ambient Background Orb -->
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-emerald-500/10 dark:bg-emerald-500/5 blur-[100px] pointer-events-none"></div>

    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <ShieldCheck class="w-6 h-6 text-primary" />
          <span>Status Pages</span>
        </h2>
        <p class="text-xs text-muted-foreground">Manage external public status pages to showcase system metrics.</p>
      </div>

      <!-- Actions -->
      <div class="flex items-center gap-2">
        <Button variant="outline" size="sm" @click="fetchAll" class="h-9">
          <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
          <span>Refresh</span>
        </Button>
        <Button @click="openCreateDialog" size="sm" class="h-9 shadow-sm shadow-primary/10">
          <Plus class="w-4 h-4 mr-1.5" />
          <span>Create Page</span>
        </Button>
      </div>
    </div>

    <!-- Main Card -->
    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <CardTitle class="text-sm font-bold text-foreground">Status Pages Configurations</CardTitle>
            <CardDescription class="text-xs">Configure public links to display select services.</CardDescription>
          </div>

          <!-- Search Bar -->
          <div class="relative w-full sm:w-72">
            <Search class="absolute left-3 top-2 h-4 w-4 text-muted-foreground" />
            <Input
              v-model="searchQuery"
              placeholder="Search status pages..."
              class="pl-8 h-8"
            />
          </div>
        </div>
      </CardHeader>

      <CardContent class="p-0">
        <div v-if="loading && statusPages.length === 0" class="flex flex-col items-center justify-center py-20 gap-3">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
          <p class="text-sm text-muted-foreground">Loading status pages...</p>
        </div>

        <div v-else-if="filteredPages.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
          <Layout class="w-12 h-12 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-bold text-foreground">No status pages found</p>
          <p class="text-xs text-muted-foreground mt-1">Create one to display public service status details.</p>
        </div>

        <StatusPageTable
          v-else
          :pages="pagesWithStats"
          @edit="openEditDialog"
          @delete="openDeleteDialog"
        />
      </CardContent>
    </Card>

    <!-- Create / Edit Dialog -->
    <Dialog v-model:open="isFormDialogOpen">
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{{ isEditMode ? 'Edit Status Page' : 'Create Status Page' }}</DialogTitle>
          <DialogDescription>
            Group select monitors to compile a dashboard accessible by external visitors without login.
          </DialogDescription>
        </DialogHeader>

        <form @submit.prevent="handleFormSubmit" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="name">Page Name</Label>
            <Input id="name" v-model="formName" placeholder="e.g. Public Services Dashboard" required />
          </div>

          <div class="space-y-2">
            <Label for="slug">URL Slug</Label>
            <div class="flex items-center rounded-md border border-input bg-transparent shadow-sm">
              <span class="px-3 text-xs text-muted-foreground border-r border-input bg-muted/30 py-2 rounded-l-md font-mono">/status/</span>
              <input
                id="slug"
                v-model="formSlug"
                placeholder="e.g. operations"
                class="flex h-9 w-full bg-transparent px-3 py-1 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50"
                required
              />
            </div>
          </div>

          <div class="space-y-2">
            <Label for="description">Description</Label>
            <textarea
              id="description"
              v-model="formDescription"
              placeholder="e.g. Live tracking status of primary internal endpoints."
              rows="3"
              class="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 resize-none"
            ></textarea>
          </div>

          <!-- Monitor Selections -->
          <div class="space-y-3 pt-2">
            <Label class="text-xs font-bold uppercase tracking-wider text-muted-foreground">Select Monitors to Include</Label>

            <div v-if="availableMonitors.length === 0" class="text-xs text-muted-foreground py-4 border border-dashed border-border/80 rounded-lg text-center">
              No monitors available. Go to Monitors page to add some first.
            </div>

            <div v-else class="grid grid-cols-1 gap-2 max-h-40 overflow-y-auto border border-border/60 rounded-md p-3 bg-muted/10">
              <div
                v-for="mon in availableMonitors"
                :key="mon.id"
                @click="toggleMonitorSelection(mon.id)"
                class="flex items-center gap-3 p-2 rounded-md hover:bg-muted/40 cursor-pointer transition-colors border border-transparent"
                :class="{ 'bg-primary/5 border-primary/20': formSelectedMonitors.includes(mon.id) }"
              >
                <!-- Checked Circle -->
                <div class="h-4 w-4 rounded border flex items-center justify-center transition-all"
                  :class="formSelectedMonitors.includes(mon.id) ? 'bg-primary border-primary text-primary-foreground' : 'border-input'"
                >
                  <svg v-if="formSelectedMonitors.includes(mon.id)" class="w-2.5 h-2.5 fill-current" viewBox="0 0 20 20">
                    <path d="M0 11l2-2 5 5L18 3l2 2L7 18z"/>
                  </svg>
                </div>
                <div class="flex flex-col select-none">
                  <span class="text-xs font-bold text-foreground leading-none">{{ mon.name }}</span>
                  <span class="text-[9px] text-muted-foreground mt-0.5 truncate max-w-[280px]">{{ mon.url }}</span>
                </div>
              </div>
            </div>
          </div>

          <p v-if="error" class="text-xs font-semibold text-destructive mt-2">{{ error }}</p>

          <DialogFooter class="pt-4 border-t border-border/40 mt-4">
            <Button type="button" variant="outline" @click="isFormDialogOpen = false">Cancel</Button>
            <Button type="submit" :disabled="formLoading" class="min-w-24">
              <Loader2 v-if="formLoading" class="w-4 h-4 mr-1.5 animate-spin" />
              <span>Save</span>
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>

    <!-- Delete Confirmation Dialog -->
    <Dialog v-model:open="isDeleteDialogOpen">
      <DialogContent class="sm:max-w-[400px]">
        <DialogHeader>
          <DialogTitle class="text-destructive flex items-center gap-2">
            <Trash2 class="w-5 h-5" />
            <span>Delete Status Page</span>
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete the status page "{{ actionPage?.name }}"? This action will permanently remove its configurations.
          </DialogDescription>
        </DialogHeader>

        <p v-if="error" class="text-xs font-semibold text-destructive mt-2">{{ error }}</p>

        <DialogFooter class="pt-4 border-t border-border/40 mt-4">
          <Button type="button" variant="outline" @click="isDeleteDialogOpen = false">Cancel</Button>
          <Button type="button" variant="destructive" :disabled="formLoading" @click="handleDeleteConfirm" class="min-w-24">
            <Loader2 v-if="formLoading" class="w-4 h-4 mr-1.5 animate-spin" />
            <span>Delete</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
