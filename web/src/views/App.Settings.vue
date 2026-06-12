<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useAuth } from '@/composables/useAuth';
import { useAppTitle } from '@/composables/useAppTitle';
import { useSettings } from '@/composables/useSettings';
import { settingSchema } from '@/validations/setting';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { toast } from 'vue-sonner';
import {
  Settings,
  Loader2,
  RefreshCw,
  Save,
  ShieldAlert,
  Sliders,
  Mail,
  Lock,
  MessageSquare
} from '@lucide/vue';
import gsap from 'gsap';

// Auth checks
const { currentUser } = useAuth();
const isAdmin = computed(() => currentUser.value?.role === 'admin');

// Composable states
const {
  settings,
  loading,
  saveLoading,
  error,
  fetchSettingsData,
  saveSettingsData
} = useSettings();

const { fetchAppTitle } = useAppTitle();

// Standard Setting fields
const systemName = ref('ping-uptime');
const adminEmail = ref('');
const allowRegistration = ref('true');

// SMTP Setting fields
const smtpHost = ref('');
const smtpPort = ref('587');
const smtpUsername = ref('');
const smtpPassword = ref('');
const smtpSender = ref('');
const smtpEncryption = ref('TLS');

// Bot Setting fields
const discordBotToken = ref('');
const telegramBotToken = ref('');
const discordClientId = ref('');

// Fetch settings wrapper
const fetchAll = async () => {
  try {
    await fetchSettingsData();

    // Map fetched values
    const nameSetting = settings.value.find(s => s.key === 'system_name');
    if (nameSetting) systemName.value = nameSetting.value;

    const emailSetting = settings.value.find(s => s.key === 'admin_email');
    if (emailSetting) adminEmail.value = emailSetting.value;

    const regSetting = settings.value.find(s => s.key === 'allow_registration');
    if (regSetting) allowRegistration.value = regSetting.value;

    const hostSetting = settings.value.find(s => s.key === 'smtp_host');
    if (hostSetting) smtpHost.value = hostSetting.value;

    const portSetting = settings.value.find(s => s.key === 'smtp_port');
    if (portSetting) smtpPort.value = portSetting.value;

    const userSetting = settings.value.find(s => s.key === 'smtp_username');
    if (userSetting) smtpUsername.value = userSetting.value;

    const passSetting = settings.value.find(s => s.key === 'smtp_password');
    if (passSetting) smtpPassword.value = passSetting.value;

    const senderSetting = settings.value.find(s => s.key === 'smtp_sender');
    if (senderSetting) smtpSender.value = senderSetting.value;

    const encSetting = settings.value.find(s => s.key === 'smtp_encryption');
    if (encSetting) smtpEncryption.value = encSetting.value;

    const discordSetting = settings.value.find(s => s.key === 'discord_bot_token');
    if (discordSetting) discordBotToken.value = discordSetting.value;

    const telegramSetting = settings.value.find(s => s.key === 'telegram_bot_token');
    if (telegramSetting) telegramBotToken.value = telegramSetting.value;

    const clientIdSetting = settings.value.find(s => s.key === 'discord_client_id');
    if (clientIdSetting) discordClientId.value = clientIdSetting.value;

  } catch (err) {
    console.error('Failed to load settings:', err);
  }
};

// Save Settings
const handleSaveSettings = async () => {
  if (!isAdmin.value) return;
  error.value = '';

  const rawPayload = {
    system_name: systemName.value,
    admin_email: adminEmail.value,
    allow_registration: allowRegistration.value,
    smtp_host: smtpHost.value,
    smtp_port: smtpPort.value,
    smtp_username: smtpUsername.value,
    smtp_password: smtpPassword.value,
    smtp_sender: smtpSender.value,
    smtp_encryption: smtpEncryption.value,
    discord_bot_token: discordBotToken.value,
    telegram_bot_token: telegramBotToken.value,
  };

  // Perform client side Zod validation
  const validation = settingSchema.safeParse(rawPayload);
  if (!validation.success) {
    const firstError = validation.error.errors[0]?.message || 'Validation failed';
    toast.error(firstError);
    error.value = firstError;
    return;
  }

  try {
    const payloads = [
      { key: 'system_name', value: systemName.value, description: 'Friendly name of the status monitoring application.' },
      { key: 'admin_email', value: adminEmail.value, description: 'Global administrator email address to send alert backups.' },
      { key: 'allow_registration', value: allowRegistration.value, description: 'Allow new guest account registrations.' },
      { key: 'smtp_host', value: smtpHost.value, description: 'SMTP Host for sending email notifications.' },
      { key: 'smtp_port', value: smtpPort.value, description: 'SMTP Port (e.g. 587 or 465).' },
      { key: 'smtp_username', value: smtpUsername.value, description: 'SMTP Username/Authentication Email.' },
      { key: 'smtp_password', value: smtpPassword.value, description: 'SMTP Password/App Password.' },
      { key: 'smtp_sender', value: smtpSender.value, description: 'SMTP Sender Address.' },
      { key: 'smtp_encryption', value: smtpEncryption.value, description: 'SMTP Encryption: SSL, TLS, or None.' },
      { key: 'discord_bot_token', value: discordBotToken.value, description: 'Discord Bot Token for sending notifications.' },
      { key: 'telegram_bot_token', value: telegramBotToken.value, description: 'Telegram Bot Token for sending notifications.' },
      { key: 'discord_client_id', value: discordClientId.value, description: 'Discord Client ID for bot invite OAuth URL.' },
    ];

    await saveSettingsData(payloads);
    toast.success('System settings updated successfully!');
    fetchAppTitle();
    await fetchAll();
  } catch (err: any) {
    console.error('Failed to save settings:', err);
    const msg = err.response?.data?.error || 'Failed to save system settings.';
    toast.error(msg);
    error.value = msg;
  }
};

onMounted(() => {
  fetchAll();
  gsap.fromTo('.settings-card',
    { opacity: 0, y: 15 },
    { opacity: 1, y: 0, duration: 0.5, stagger: 0.1, ease: 'power2.out' }
  );
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-4xl mx-auto min-h-[calc(100vh-4rem)]">
    <!-- Ambient Background Orb -->
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-indigo-500/5 dark:bg-indigo-500/2 blur-[100px] pointer-events-none"></div>

    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <Settings class="w-6 h-6 text-primary" />
          <span>System Settings</span>
        </h2>
        <p class="text-xs text-muted-foreground">Configure application names, registrations, and SMTP configurations.</p>
      </div>

      <!-- Actions -->
      <Button variant="outline" size="sm" @click="fetchAll" class="h-9">
        <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
        <span>Refresh</span>
      </Button>
    </div>

    <!-- Role alert block if not admin -->
    <div v-if="!isAdmin" class="flex items-center gap-3 p-4 rounded-lg border border-amber-500/25 bg-amber-500/5 text-amber-600 dark:text-amber-400 z-10 relative">
      <ShieldAlert class="w-5 h-5 shrink-0" />
      <div class="text-xs leading-relaxed">
        <span class="font-bold">Read-Only Mode:</span> You are logged in as a standard user. You can view these configurations but modifying settings requires administrator privileges.
      </div>
    </div>

    <!-- Main Settings Form -->
    <form @submit.prevent="handleSaveSettings" class="space-y-6 z-10 relative">
      <!-- General Settings Card -->
      <Card class="settings-card border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
        <CardHeader class="pb-3 border-b border-border/40">
          <div class="flex items-center gap-2.5">
            <Sliders class="w-5 h-5 text-primary" />
            <div>
              <CardTitle class="text-sm font-bold text-foreground">Global Options</CardTitle>
              <CardDescription class="text-xs">General preferences for the service instance.</CardDescription>
            </div>
          </div>
        </CardHeader>

        <CardContent class="space-y-4 pt-6">
          <div v-if="loading" class="flex flex-col items-center justify-center py-10 gap-3">
            <Loader2 class="w-6 h-6 text-primary animate-spin" />
            <p class="text-xs text-muted-foreground">Loading configurations...</p>
          </div>

          <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- System Name -->
            <div class="space-y-2">
              <Label for="system_name">System Name</Label>
              <Input
                id="system_name"
                v-model="systemName"
                placeholder="ping-uptime"
                :disabled="!isAdmin || saveLoading"
                required
              />
              <span class="text-[10px] text-muted-foreground">The display name of the status dashboard header.</span>
            </div>

            <!-- Admin Alert Email -->
            <div class="space-y-2">
              <Label for="admin_email">Notification Backup Email</Label>
              <Input
                id="admin_email"
                v-model="adminEmail"
                type="email"
                placeholder="admin@company.com"
                :disabled="!isAdmin || saveLoading"
              />
              <span class="text-[10px] text-muted-foreground">Backups email address where pings failure logs will be CC'd.</span>
            </div>

            <!-- Allow Registration -->
            <div class="space-y-2">
              <Label for="allow_reg">Allow Public Registrations</Label>
              <Select v-model="allowRegistration" :disabled="!isAdmin || saveLoading">
                <SelectTrigger id="allow_reg" class="h-9">
                  <SelectValue placeholder="Select registration option" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="true">Yes, anyone can sign up</SelectItem>
                  <SelectItem value="false">No, administrators invite only</SelectItem>
                </SelectContent>
              </Select>
              <span class="text-[10px] text-muted-foreground">Toggles whether new users can register a dashboard account from the index.</span>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- SMTP Settings Card -->
      <Card class="settings-card border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
        <CardHeader class="pb-3 border-b border-border/40">
          <div class="flex items-center gap-2.5">
            <Mail class="w-5 h-5 text-primary" />
            <div>
              <CardTitle class="text-sm font-bold text-foreground">SMTP Outgoing Email Settings</CardTitle>
              <CardDescription class="text-xs">Configure mail server credentials to send status reports when monitors go down/up.</CardDescription>
            </div>
          </div>
        </CardHeader>

        <CardContent class="space-y-4 pt-6">
          <div v-if="loading" class="flex flex-col items-center justify-center py-10 gap-3">
            <Loader2 class="w-6 h-6 text-primary animate-spin" />
            <p class="text-xs text-muted-foreground">Loading configurations...</p>
          </div>

          <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- SMTP Host -->
            <div class="space-y-2">
              <Label for="smtp_host">SMTP Host</Label>
              <Input
                id="smtp_host"
                v-model="smtpHost"
                placeholder="e.g. smtp.gmail.com"
                :disabled="!isAdmin || saveLoading"
              />
              <span class="text-[10px] text-muted-foreground">Hostname of your outgoing mail service.</span>
            </div>

            <!-- SMTP Port -->
            <div class="space-y-2">
              <Label for="smtp_port">SMTP Port</Label>
              <Input
                id="smtp_port"
                v-model="smtpPort"
                placeholder="e.g. 587"
                :disabled="!isAdmin || saveLoading"
              />
              <span class="text-[10px] text-muted-foreground">Common ports are 587 (TLS), 465 (SSL), or 25 (None).</span>
            </div>

            <!-- SMTP Username -->
            <div class="space-y-2">
              <Label for="smtp_username">SMTP Username</Label>
              <Input
                id="smtp_username"
                v-model="smtpUsername"
                placeholder="e.g. alert@yourcompany.com"
                :disabled="!isAdmin || saveLoading"
              />
              <span class="text-[10px] text-muted-foreground">Authentication user identity for connecting to SMTP server.</span>
            </div>

            <!-- SMTP Password -->
            <div class="space-y-2">
              <Label for="smtp_password">SMTP Password</Label>
              <div class="relative">
                <Lock class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  id="smtp_password"
                  v-model="smtpPassword"
                  type="password"
                  placeholder="••••••••"
                  class="pl-9"
                  :disabled="!isAdmin || saveLoading"
                />
              </div>
              <span class="text-[10px] text-muted-foreground">Authentication password or secure App Password.</span>
            </div>

            <!-- SMTP Sender -->
            <div class="space-y-2">
              <Label for="smtp_sender">Sender Address ("From")</Label>
              <Input
                id="smtp_sender"
                v-model="smtpSender"
                placeholder="e.g. noreply@yourcompany.com"
                :disabled="!isAdmin || saveLoading"
              />
              <span class="text-[10px] text-muted-foreground">Email address that recipient will see on status emails.</span>
            </div>

            <!-- SMTP Encryption -->
            <div class="space-y-2">
              <Label for="smtp_enc">Security Encryption</Label>
              <Select v-model="smtpEncryption" :disabled="!isAdmin || saveLoading">
                <SelectTrigger id="smtp_enc" class="h-9">
                  <SelectValue placeholder="Select encryption" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="TLS">STARTTLS (Standard 587)</SelectItem>
                  <SelectItem value="SSL">SSL/TLS (Implicit 465)</SelectItem>
                  <SelectItem value="None">None (Unencrypted 25)</SelectItem>
                </SelectContent>
              </Select>
              <span class="text-[10px] text-muted-foreground">Matches host port security guidelines.</span>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Notification Bot Settings Card -->
      <Card class="settings-card border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
        <CardHeader class="pb-3 border-b border-border/40">
          <div class="flex items-center gap-2.5">
            <MessageSquare class="w-5 h-5 text-primary" />
            <div>
              <CardTitle class="text-sm font-bold text-foreground">Notification Bot Settings</CardTitle>
              <CardDescription class="text-xs">Configure Discord and Telegram bot tokens to send notifications when status changes.</CardDescription>
            </div>
          </div>
        </CardHeader>

        <CardContent class="space-y-4 pt-6">
          <div v-if="loading" class="flex flex-col items-center justify-center py-10 gap-3">
            <Loader2 class="w-6 h-6 text-primary animate-spin" />
            <p class="text-xs text-muted-foreground">Loading configurations...</p>
          </div>

          <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Discord Bot Token -->
            <div class="space-y-2">
              <Label for="discord_bot_token">Discord Bot Token</Label>
              <div class="relative">
                <Lock class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  id="discord_bot_token"
                  v-model="discordBotToken"
                  type="password"
                  placeholder="••••••••"
                  class="pl-9"
                  :disabled="!isAdmin || saveLoading"
                />
              </div>
              <span class="text-[10px] text-muted-foreground">The bot token used to authenticate with Discord API.</span>
            </div>

            <!-- Telegram Bot Token -->
            <div class="space-y-2">
              <Label for="telegram_bot_token">Telegram Bot Token</Label>
              <div class="relative">
                <Lock class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  id="telegram_bot_token"
                  v-model="telegramBotToken"
                  type="password"
                  placeholder="••••••••"
                  class="pl-9"
                  :disabled="!isAdmin || saveLoading"
                />
              </div>
              <span class="text-[10px] text-muted-foreground">The bot token from Telegram's @BotFather.</span>
            </div>

            <!-- Discord Client ID -->
            <div class="space-y-2">
              <Label for="discord_client_id">Discord Client ID</Label>
              <Input
                id="discord_client_id"
                v-model="discordClientId"
                placeholder="e.g. 123456789012345678"
                :disabled="!isAdmin || saveLoading"
              />
              <span class="text-[10px] text-muted-foreground">Used to generate OAuth invite URL for users to add the bot to their server.</span>
            </div>
          </div>
        </CardContent>

        <CardFooter v-if="isAdmin && !loading" class="flex justify-end pt-4 border-t border-border/40 bg-muted/10">
          <Button type="submit" :disabled="saveLoading" class="min-w-28 shadow-sm shadow-primary/10">
            <Loader2 v-if="saveLoading" class="w-4 h-4 mr-1.5 animate-spin" />
            <Save v-else class="w-4 h-4 mr-1.5" />
            <span>Save Options</span>
          </Button>
        </CardFooter>
      </Card>
    </form>
  </div>
</template>
