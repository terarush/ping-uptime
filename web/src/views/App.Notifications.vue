<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useNotificationChannels, type NotificationChannel } from '@/composables/useNotificationChannels';
import { useSettings } from '@/composables/useSettings';
import { notificationChannelSchema } from '@/validations/notification-channel';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import NotificationChannelTable from '@/components/notification-channel-table.vue';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { toast } from 'vue-sonner';
import {
  Bell,
  Plus,
  Trash2,
  Search,
  Loader2,
  RefreshCw,
  ExternalLink,
  MessageCircle
} from '@lucide/vue';
import gsap from 'gsap';

// Composable states
const {
  channels,
  loading,
  error,
  fetchChannelsData,
  createChannel,
  updateChannel,
  deleteChannel
} = useNotificationChannels();

const { settings: globalSettings, fetchSettingsData } = useSettings();

const searchQuery = ref('');

const hasDiscordBot = computed(() => {
  const s = globalSettings.value.find(x => x.key === 'discord_bot_token');
  return !!s?.value;
});

const hasTelegramBot = computed(() => {
  const s = globalSettings.value.find(x => x.key === 'telegram_bot_token');
  return !!s?.value;
});

const discordInviteLink = computed(() => {
  const clientId = globalSettings.value.find(x => x.key === 'discord_client_id');
  if (!clientId?.value) return '';
  return `https://discord.com/api/oauth2/authorize?client_id=${clientId.value}&permissions=0&scope=bot%20applications.commands`;
});

// Dialog states
const isFormDialogOpen = ref(false);
const isDeleteDialogOpen = ref(false);
const formLoading = ref(false);
const actionChannel = ref<NotificationChannel | null>(null);

// Form Fields
const formName = ref('');
const formType = ref('webhook');
const formEnabled = ref(true);

// Dynamic config inputs based on Type
const configWebhookUrl = ref('');
const configTelegramToken = ref('');
const configTelegramChatId = ref('');
const configEmailAddress = ref('');
const configDiscordChannelId = ref('');

const isEditMode = computed(() => !!actionChannel.value);

// Fetch notification channels wrapper
const fetchAll = async () => {
  try {
    await fetchChannelsData();
  } catch (err) {
    console.error('Failed to load alert channels list:', err);
  } finally {
    setTimeout(animateTableRows, 50);
  }
};

// Filtered channels
const filteredChannels = computed(() => {
  if (!channels.value) return [];
  const query = searchQuery.value.toLowerCase().trim();
  if (!query) return channels.value;
  return channels.value.filter(c =>
    c.name.toLowerCase().includes(query) ||
    c.type.toLowerCase().includes(query)
  );
});

// Watch type change to reset dynamic form fields
watch(formType, () => {
  configWebhookUrl.value = '';
  configTelegramToken.value = '';
  configTelegramChatId.value = '';
  configEmailAddress.value = '';
  configDiscordChannelId.value = '';
});

// Reset form
const resetForm = () => {
  formName.value = '';
  formType.value = 'webhook';
  formEnabled.value = true;
  configWebhookUrl.value = '';
  configTelegramToken.value = '';
  configTelegramChatId.value = '';
  configEmailAddress.value = '';
  configDiscordChannelId.value = '';
  actionChannel.value = null;
};

const openCreateDialog = () => {
  resetForm();
  isFormDialogOpen.value = true;
};

const openEditDialog = (channel: NotificationChannel) => {
  resetForm();
  actionChannel.value = channel;
  formName.value = channel.name;
  formType.value = channel.type;
  formEnabled.value = channel.enabled;

  // Parse config
  try {
    const parsed = JSON.parse(channel.config);
    if (channel.type === 'email') {
      configEmailAddress.value = parsed.email || '';
    } else if (channel.type === 'telegram') {
      configTelegramToken.value = parsed.bot_token || '';
      configTelegramChatId.value = parsed.chat_id || '';
    } else if (channel.type === 'discord_bot') {
      configDiscordChannelId.value = parsed.channel_id || '';
    } else {
      configWebhookUrl.value = parsed.webhook_url || parsed.url || '';
    }
  } catch (e) {
    console.error('Failed to parse channel config JSON:', e);
  }

  isFormDialogOpen.value = true;
};

const openDeleteDialog = (channel: NotificationChannel) => {
  actionChannel.value = channel;
  isDeleteDialogOpen.value = true;
};

// Submit form (Create / Update)
const handleFormSubmit = async () => {
  formLoading.value = true;

  try {
    // Construct config JSON object
    let configObj: Record<string, string> = {};
    if (formType.value === 'email') {
      configObj = { email: configEmailAddress.value };
    } else if (formType.value === 'telegram') {
      configObj = { bot_token: configTelegramToken.value, chat_id: configTelegramChatId.value };
    } else if (formType.value === 'discord_bot') {
      configObj = { channel_id: configDiscordChannelId.value };
    } else {
      configObj = { webhook_url: configWebhookUrl.value };
    }

    const rawPayload = {
      name: formName.value,
      type: formType.value,
      config: JSON.stringify(configObj),
      enabled: formEnabled.value,
    };

    // Perform client side Zod validation
    const validation = notificationChannelSchema.safeParse(rawPayload);
    if (!validation.success) {
      const firstError = validation.error.errors[0]?.message || 'Validation failed';
      toast.error(firstError);
      formLoading.value = false;
      return;
    }

    if (isEditMode.value && actionChannel.value) {
      await updateChannel(actionChannel.value.id, rawPayload);
      toast.success(`Alert channel "${formName.value}" updated successfully!`);
    } else {
      await createChannel(rawPayload);
      toast.success(`Alert channel "${formName.value}" created successfully!`);
    }

    isFormDialogOpen.value = false;
    await fetchAll();
  } catch (err: any) {
    console.error('Failed to save alert channel:', err);
    const msg = err.response?.data?.error || 'Failed to save alert channel configuration.';
    toast.error(msg);
  } finally {
    formLoading.value = false;
  }
};

// Toggle channel state inline
const handleToggleEnabled = async (item: NotificationChannel) => {
  try {
    const payload = {
      name: item.name,
      type: item.type,
      config: item.config,
      enabled: !item.enabled,
    };
    await updateChannel(item.id, payload);
    item.enabled = !item.enabled;
    toast.success(`Alert channel "${item.name}" ${item.enabled ? 'enabled' : 'disabled'}!`);
  } catch (err: any) {
    console.error('Failed to toggle channel status:', err);
    toast.error('Failed to toggle channel status.');
  }
};

// Delete channel
const handleDeleteConfirm = async () => {
  if (!actionChannel.value) return;
  formLoading.value = true;

  try {
    await deleteChannel(actionChannel.value.id);
    toast.success(`Alert channel "${actionChannel.value.name}" deleted successfully!`);
    isDeleteDialogOpen.value = false;
    await fetchAll();
  } catch (err: any) {
    console.error('Failed to delete alert channel:', err);
    const msg = err.response?.data?.error || 'Failed to delete alert channel.';
    toast.error(msg);
  } finally {
    formLoading.value = false;
  }
};

// GSAP Animations
const animateTableRows = () => {
  gsap.fromTo('.channel-row',
    { opacity: 0, y: 15 },
    { opacity: 1, y: 0, duration: 0.4, stagger: 0.05, ease: 'power2.out' }
  );
};

onMounted(() => {
  fetchAll();
  fetchSettingsData();
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
          <Bell class="w-6 h-6 text-primary" />
          <span>Notification Settings</span>
        </h2>
        <p class="text-xs text-muted-foreground">Setup alerts and automated notification integrations.</p>
      </div>

      <!-- Actions -->
      <div class="flex items-center gap-2">
        <Button variant="outline" size="sm" @click="fetchAll" class="h-9">
          <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
          <span>Refresh</span>
        </Button>
        <Button @click="openCreateDialog" size="sm" class="h-9 shadow-sm shadow-primary/10">
          <Plus class="w-4 h-4 mr-1.5" />
          <span>Add Channel</span>
        </Button>
      </div>
    </div>

    <!-- Discord Bot Invite Banner -->
    <div v-if="discordInviteLink" class="flex items-center gap-3 p-4 rounded-lg border border-indigo-500/25 bg-indigo-500/5 z-10 relative">
      <MessageCircle class="w-5 h-5 shrink-0 text-indigo-500" />
      <div class="flex-1 min-w-0">
        <p class="text-xs font-bold text-foreground">Discord Bot</p>
        <p class="text-[11px] text-muted-foreground">Invite the bot to your server to receive alerts directly in a Discord channel.</p>
      </div>
      <a :href="discordInviteLink" target="_blank" rel="noopener noreferrer">
        <Button size="sm" class="h-8 gap-1.5">
          <ExternalLink class="w-3.5 h-3.5" />
          <span>Invite Bot</span>
        </Button>
      </a>
    </div>

    <!-- Main Card -->
    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <CardTitle class="text-sm font-bold text-foreground">Alert Integrations</CardTitle>
            <CardDescription class="text-xs">Channels to notify when monitors transition to DOWN or UP.</CardDescription>
          </div>

          <!-- Search Bar -->
          <div class="relative w-full sm:w-72">
            <Search class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              v-model="searchQuery"
              placeholder="Search channels..."
              class="pl-9 h-9"
            />
          </div>
        </div>
      </CardHeader>

      <CardContent class="p-0">
        <div v-if="loading && channels.length === 0" class="flex flex-col items-center justify-center py-20 gap-3">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
          <p class="text-sm text-muted-foreground">Loading configurations...</p>
        </div>

        <div v-else-if="filteredChannels.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
          <Bell class="w-12 h-12 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-bold text-foreground">No channels configured</p>
          <p class="text-xs text-muted-foreground mt-1">Setup alert destinations to receive outage reports.</p>
        </div>

        <!-- Notification Table Component -->
        <NotificationChannelTable
          v-else
          :channels="filteredChannels"
          @edit="openEditDialog"
          @delete="openDeleteDialog"
          @toggle="handleToggleEnabled"
        />
      </CardContent>
    </Card>

    <!-- Create / Edit Dialog -->
    <Dialog v-model:open="isFormDialogOpen">
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{{ isEditMode ? 'Edit Alert Channel' : 'Add Alert Channel' }}</DialogTitle>
          <DialogDescription>
            Configure integration details to receive instant outage updates.
          </DialogDescription>
        </DialogHeader>

        <form @submit.prevent="handleFormSubmit" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="name">Channel Name</Label>
            <Input id="name" v-model="formName" placeholder="e.g. Server Admin Discord Webhook" required />
          </div>

          <div class="space-y-2">
            <Label for="type">Integration Type</Label>
            <Select v-model="formType">
              <SelectTrigger id="type" class="w-full h-9">
                <SelectValue placeholder="Select type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="webhook">Custom Webhook</SelectItem>
                <SelectItem value="slack">Slack Webhook</SelectItem>
                <SelectItem v-if="hasDiscordBot" value="discord_bot">Discord Bot</SelectItem>
                <SelectItem v-if="hasTelegramBot" value="telegram">Telegram Bot</SelectItem>
                <SelectItem value="email">Email</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="flex items-center gap-2 select-none">
            <Switch id="enabled" v-model:checked="formEnabled" />
            <Label for="enabled" class="cursor-pointer">Enabled</Label>
          </div>

          <!-- Dynamic input configurations -->
          <div v-if="formType === 'email'" class="space-y-2">
            <p class="text-xs text-muted-foreground bg-muted/30 p-3 rounded-lg border border-border/40">
              Email alerts will be sent directly to your registered account email address. No additional configuration is required.
            </p>
          </div>

          <div v-else-if="formType === 'telegram'" class="space-y-2">
            <Label for="tg_chat">Telegram Chat/User ID</Label>
            <Input id="tg_chat" v-model="configTelegramChatId" placeholder="e.g. -100123456789 or 987654321" required />
            <span class="text-[10px] text-muted-foreground">The chat or user ID to receive telegram alert messages (uses global bot token).</span>
          </div>

          <div v-else-if="formType === 'discord_bot'" class="space-y-2">
            <Label for="discord_channel">Discord Channel ID</Label>
            <Input id="discord_channel" v-model="configDiscordChannelId" placeholder="e.g. 112233445566778899" required />
            <span class="text-[10px] text-muted-foreground">The Discord channel ID where the bot will post alerts (uses global bot token).</span>
          </div>

          <div v-else class="space-y-2">
            <Label for="webhook">Webhook URL</Label>
            <Input id="webhook" v-model="configWebhookUrl" type="url" placeholder="https://hooks.slack.com/services/..." required />
          </div>

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
            <span>Delete Alert Channel</span>
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete the alert channel "{{ actionChannel?.name }}"? This action will permanently remove its integration.
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
