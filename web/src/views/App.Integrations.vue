<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useIntegrations, type Integration } from '@/composables/useIntegrations';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { toast } from 'vue-sonner';
import {
  Puzzle,
  Plus,
  Trash2,
  Loader2,
  RefreshCw,
  Send,
  ToggleLeft,
  ToggleRight,
  Globe,
  Bot,
  Code,
  Bell,
  Phone,
  Webhook,
  Zap,
} from '@lucide/vue';
import gsap from 'gsap';

const {
  integrations,
  loading,
  error,
  fetchAll,
  create,
  update,
  deleteIntegration,
  test,
} = useIntegrations();

const searchQuery = ref('');

// Dialog states
const isFormDialogOpen = ref(false);
const isDeleteDialogOpen = ref(false);
const formLoading = ref(false);
const testLoading = ref<Record<number, boolean>>({});
const actionIntegration = ref<Integration | null>(null);

// Form fields
const formName = ref('');
const formType = ref('webhook');
const formConfigUrl = ref('');

const isEditMode = computed(() => !!actionIntegration.value);

const fetchAllIntegrations = async () => {
  try {
    await fetchAll();
  } catch (err) {
    console.error('Failed to load integrations:', err);
  } finally {
    setTimeout(animateCards, 50);
  }
};

const filteredIntegrations = computed(() => {
  if (!integrations.value) return [];
  const query = searchQuery.value.toLowerCase().trim();
  if (!query) return integrations.value;
  return integrations.value.filter(i =>
    i.name.toLowerCase().includes(query) ||
    i.type.toLowerCase().includes(query)
  );
});

const resetForm = () => {
  formName.value = '';
  formType.value = 'webhook';
  formConfigUrl.value = '';
  actionIntegration.value = null;
};

const openCreateDialog = () => {
  resetForm();
  isFormDialogOpen.value = true;
};

const openEditDialog = (integration: Integration) => {
  resetForm();
  actionIntegration.value = integration;
  formName.value = integration.name;
  formType.value = integration.type;
  try {
    const parsed = JSON.parse(integration.config);
    formConfigUrl.value = parsed.webhook_url || '';
  } catch {
    formConfigUrl.value = '';
  }
  isFormDialogOpen.value = true;
};

const openDeleteDialog = (integration: Integration) => {
  actionIntegration.value = integration;
  isDeleteDialogOpen.value = true;
};

const handleFormSubmit = async () => {
  formLoading.value = true;
  try {
    const configObj = { webhook_url: formConfigUrl.value };
    const rawPayload = {
      name: formName.value,
      type: formType.value,
      config: JSON.stringify(configObj),
    };

    if (isEditMode.value && actionIntegration.value) {
      await update(actionIntegration.value.id, {
        ...rawPayload,
        enabled: actionIntegration.value.enabled,
      });
      toast.success(`Integration "${formName.value}" updated`);
    } else {
      await create(rawPayload);
      toast.success(`Integration "${formName.value}" created`);
    }

    isFormDialogOpen.value = false;
    await fetchAllIntegrations();
  } catch (err: any) {
    const msg = err.response?.data?.error || 'Failed to save integration.';
    toast.error(msg);
  } finally {
    formLoading.value = false;
  }
};

const handleToggleEnabled = async (item: Integration) => {
  try {
    const payload = {
      name: item.name,
      type: item.type,
      config: item.config,
      enabled: !item.enabled,
    };
    await update(item.id, payload);
    item.enabled = !item.enabled;
    toast.success(`Integration "${item.name}" ${item.enabled ? 'enabled' : 'disabled'}`);
  } catch (err: any) {
    toast.error('Failed to toggle integration status.');
  }
};

const handleDeleteConfirm = async () => {
  if (!actionIntegration.value) return;
  formLoading.value = true;
  try {
    await deleteIntegration(actionIntegration.value.id);
    toast.success(`Integration "${actionIntegration.value.name}" deleted`);
    isDeleteDialogOpen.value = false;
    await fetchAllIntegrations();
  } catch (err: any) {
    const msg = err.response?.data?.error || 'Failed to delete integration.';
    toast.error(msg);
  } finally {
    formLoading.value = false;
  }
};

const handleTest = async (item: Integration) => {
  testLoading.value[item.id] = true;
  try {
    await test(item.id);
    toast.success(`Test notification sent to "${item.name}"`);
  } catch (err: any) {
    const msg = err.response?.data?.error || 'Test failed.';
    toast.error(msg);
  } finally {
    testLoading.value[item.id] = false;
  }
};

const typeConfig = (type: string) => {
  switch (type) {
    case 'slack': return { label: 'Slack', icon: Bell, color: 'text-purple-500', bg: 'bg-purple-500/10 border-purple-500/20' };
    case 'discord': return { label: 'Discord', icon: Bot, color: 'text-indigo-500', bg: 'bg-indigo-500/10 border-indigo-500/20' };
    case 'webhook': return { label: 'Webhook', icon: Webhook, color: 'text-cyan-500', bg: 'bg-cyan-500/10 border-cyan-500/20' };
    case 'github': return { label: 'GitHub', icon: Code, color: 'text-gray-500', bg: 'bg-gray-500/10 border-gray-500/20' };
    case 'pagerduty': return { label: 'PagerDuty', icon: Phone, color: 'text-orange-500', bg: 'bg-orange-500/10 border-orange-500/20' };
    default: return { label: type, icon: Globe, color: 'text-muted-foreground', bg: 'bg-muted/30 border-border/40' };
  }
};

const animateCards = () => {
  gsap.fromTo('.integration-card',
    { opacity: 0, y: 20 },
    { opacity: 1, y: 0, duration: 0.5, stagger: 0.06, ease: 'power2.out' }
  );
};

onMounted(() => {
  fetchAllIntegrations();
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <!-- Ambient Background Orb -->
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-violet-500/10 dark:bg-violet-500/5 blur-[100px] pointer-events-none"></div>

    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <Puzzle class="w-6 h-6 text-primary" />
          <span>Integrations</span>
        </h2>
        <p class="text-xs text-muted-foreground">Connect external services to receive alerts and events.</p>
      </div>

      <div class="flex items-center gap-2">
        <Button variant="outline" size="sm" @click="fetchAllIntegrations" class="h-9">
          <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
          <span>Refresh</span>
        </Button>
        <Button @click="openCreateDialog" size="sm" class="h-9 shadow-sm shadow-primary/10">
          <Plus class="w-4 h-4 mr-1.5" />
          <span>Connect</span>
        </Button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading && integrations.length === 0" class="flex flex-col items-center justify-center py-20 gap-3 z-10 relative">
      <Loader2 class="w-8 h-8 text-primary animate-spin" />
      <p class="text-sm text-muted-foreground">Loading integrations...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredIntegrations.length === 0" class="flex flex-col items-center justify-center py-20 text-center z-10 relative">
      <Puzzle class="w-12 h-12 text-muted-foreground/30 mb-3" />
      <p class="text-sm font-bold text-foreground">No integrations configured</p>
      <p class="text-xs text-muted-foreground mt-1">Connect a service to receive alerts and notifications.</p>
      <Button @click="openCreateDialog" variant="outline" size="sm" class="mt-4 h-9">
        <Plus class="w-4 h-4 mr-1.5" />
        Connect Integration
      </Button>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="flex flex-col items-center justify-center py-20 text-center z-10 relative">
      <p class="text-sm text-destructive font-bold">Failed to load integrations</p>
      <p class="text-xs text-muted-foreground mt-1">{{ error }}</p>
      <Button @click="fetchAllIntegrations" variant="outline" size="sm" class="mt-4 h-9">
        <RefreshCw class="w-4 h-4 mr-1.5" />
        Retry
      </Button>
    </div>

    <!-- Integration Cards Grid -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 z-10 relative">
      <Card
        v-for="integration in filteredIntegrations"
        :key="integration.id"
        class="integration-card border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md hover:border-primary/30 transition-all duration-300"
      >
        <CardContent class="p-5">
          <div class="flex items-start justify-between mb-4">
            <div class="flex items-center gap-3">
              <div
                :class="[
                  'w-10 h-10 rounded-lg flex items-center justify-center border',
                  typeConfig(integration.type).bg,
                ]"
              >
                <component :is="typeConfig(integration.type).icon" :class="['w-5 h-5', typeConfig(integration.type).color]" />
              </div>
              <div>
                <CardTitle class="text-sm font-bold text-foreground">{{ integration.name }}</CardTitle>
                <span :class="['text-[10px] font-medium', typeConfig(integration.type).color]">
                  {{ typeConfig(integration.type).label }}
                </span>
              </div>
            </div>
            <Badge
              :variant="integration.enabled ? 'secondary' : 'outline'"
              :class="[
                'text-[10px] capitalize',
                integration.enabled
                  ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20'
                  : 'text-muted-foreground/50',
              ]"
            >
              {{ integration.enabled ? 'Active' : 'Disabled' }}
            </Badge>
          </div>

          <div class="flex items-center gap-2 mt-4 pt-3 border-t border-border/20">
            <!-- Test Button -->
            <Button
              variant="outline"
              size="sm"
              class="h-8 text-xs gap-1"
              :disabled="testLoading[integration.id]"
              @click="handleTest(integration)"
            >
              <Send v-if="!testLoading[integration.id]" class="w-3.5 h-3.5" />
              <Loader2 v-else class="w-3.5 h-3.5 animate-spin" />
              Test
            </Button>

            <!-- Enabled Toggle -->
            <Button
              variant="ghost"
              size="sm"
              class="h-8 w-8 p-0"
              :class="integration.enabled ? 'text-emerald-500' : 'text-muted-foreground'"
              @click="handleToggleEnabled(integration)"
            >
              <ToggleRight v-if="integration.enabled" class="w-5 h-5" />
              <ToggleLeft v-else class="w-5 h-5" />
            </Button>

            <!-- Edit Button -->
            <Button
              variant="ghost"
              size="sm"
              class="h-8 w-8 p-0"
              @click="openEditDialog(integration)"
            >
              <Zap class="w-4 h-4" />
            </Button>

            <div class="flex-1"></div>

            <!-- Delete Button -->
            <Button
              variant="ghost"
              size="sm"
              class="h-8 w-8 p-0 text-destructive hover:bg-destructive/10"
              @click="openDeleteDialog(integration)"
            >
              <Trash2 class="w-4 h-4" />
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Create / Edit Dialog -->
    <Dialog v-model:open="isFormDialogOpen">
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{{ isEditMode ? 'Edit Integration' : 'Connect Integration' }}</DialogTitle>
          <DialogDescription>
            {{ isEditMode ? 'Update integration configuration.' : 'Configure a new external service integration.' }}
          </DialogDescription>
        </DialogHeader>

        <form @submit.prevent="handleFormSubmit" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="name">Integration Name</Label>
            <Input id="name" v-model="formName" placeholder="e.g. Team Slack Alerts" required maxlength="100" />
          </div>

          <div class="space-y-2">
            <Label for="type">Type</Label>
            <Select v-model="formType">
              <SelectTrigger id="type" class="h-9">
                <SelectValue placeholder="Select type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="webhook">Custom Webhook</SelectItem>
                <SelectItem value="slack">Slack</SelectItem>
                <SelectItem value="discord">Discord</SelectItem>
                <SelectItem value="github">GitHub</SelectItem>
                <SelectItem value="pagerduty">PagerDuty</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="space-y-2">
            <Label for="webhook_url">Webhook URL</Label>
            <Input
              id="webhook_url"
              v-model="formConfigUrl"
              type="url"
              placeholder="https://hooks.example.com/..."
              required
            />
          </div>

          <DialogFooter class="pt-4 border-t border-border/40 mt-4">
            <Button type="button" variant="outline" @click="isFormDialogOpen = false">Cancel</Button>
            <Button type="submit" :disabled="formLoading" class="min-w-24">
              <Loader2 v-if="formLoading" class="w-4 h-4 mr-1.5 animate-spin" />
              <span>{{ isEditMode ? 'Update' : 'Connect' }}</span>
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
            <span>Delete Integration</span>
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete "{{ actionIntegration?.name }}"? This action is permanent.
          </DialogDescription>
        </DialogHeader>

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
