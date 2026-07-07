<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useApiTokens, type ApiToken, type CreateTokenResult } from '@/composables/useApiTokens';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog';
import { Badge } from '@/components/ui/badge';
import { toast } from 'vue-sonner';
import { KeyRound, Plus, Loader2, RefreshCw, Trash2, Copy, CheckCircle, Search } from '@lucide/vue';
import gsap from 'gsap';

const {
  tokens,
  loading,
  error,
  fetchAll,
  create,
  revoke,
  formatDate,
} = useApiTokens();

const isCreateOpen = ref(false);
const isRevokeOpen = ref(false);
const isTokenDisplayOpen = ref(false);
const formLoading = ref(false);
const actionToken = ref<ApiToken | null>(null);
const lastCreatedToken = ref<CreateTokenResult | null>(null);
const copied = ref(false);

const formName = ref('');
const formExpiresAt = ref('');

const resetForm = () => {
  formName.value = '';
  formExpiresAt.value = '';
};

const openCreate = () => {
  resetForm();
  isCreateOpen.value = true;
};

const openRevoke = (token: ApiToken) => {
  actionToken.value = token;
  isRevokeOpen.value = true;
};

const handleCreate = async () => {
  if (!formName.value.trim()) {
    toast.error('Token name is required');
    return;
  }
  formLoading.value = true;
  try {
    const expiresAt = formExpiresAt.value
      ? new Date(formExpiresAt.value).toISOString()
      : undefined;
    const result = await create(formName.value.trim(), expiresAt);
    if (result) {
      lastCreatedToken.value = result;
      isCreateOpen.value = false;
      isTokenDisplayOpen.value = true;
    }
    await fetchAll();
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to create token');
  } finally {
    formLoading.value = false;
  }
};

const handleRevoke = async () => {
  if (!actionToken.value) return;
  formLoading.value = true;
  try {
    await revoke(actionToken.value.id);
    toast.success(`Token "${actionToken.value.name}" revoked`);
    actionToken.value.is_revoked = true;
    isRevokeOpen.value = false;
    actionToken.value = null;
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to revoke token');
  } finally {
    formLoading.value = false;
  }
};

const copyToken = async () => {
  if (!lastCreatedToken.value?.raw_token) return;
  try {
    await navigator.clipboard.writeText(lastCreatedToken.value.raw_token);
    copied.value = true;
    setTimeout(() => { copied.value = false; }, 3000);
  } catch {
    toast.error('Failed to copy to clipboard');
  }
};

const searchQuery = ref('');

const filteredTokens = computed(() => {
  if (!searchQuery.value.trim()) return tokens.value;
  const q = searchQuery.value.toLowerCase().trim();
  return tokens.value.filter(t => t.name.toLowerCase().includes(q));
});

const activeTokens = computed(() =>
  filteredTokens.value.filter(t => !t.is_revoked)
);

const revokedTokens = computed(() =>
  filteredTokens.value.filter(t => t.is_revoked)
);

onMounted(async () => {
  await fetchAll();
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
  gsap.fromTo('.token-row',
    { opacity: 0, y: 12 },
    { opacity: 1, y: 0, duration: 0.4, stagger: 0.04, ease: 'power2.out' }
  );
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <!-- Ambient Background Orb -->
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-indigo-500/10 dark:bg-indigo-500/5 blur-[100px] pointer-events-none"></div>

    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <KeyRound class="w-6 h-6 text-primary" />
          <span>API Tokens</span>
        </h2>
        <p class="text-xs text-muted-foreground">Manage personal access tokens for API authentication.</p>
      </div>

      <div class="flex items-center gap-2">
        <Button variant="outline" size="sm" @click="fetchAll" class="h-9">
          <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
          Refresh
        </Button>
        <Button @click="openCreate" size="sm" class="h-9">
          <Plus class="w-4 h-4 mr-1.5" />
          Generate Token
        </Button>
      </div>
    </div>

    <!-- Active Tokens Card -->
    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <CardTitle class="text-sm font-bold">Active Tokens</CardTitle>
            <CardDescription class="text-xs">{{ activeTokens.length }} active token{{ activeTokens.length !== 1 ? 's' : '' }}</CardDescription>
          </div>
          <div class="relative w-full sm:w-64">
            <Search class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input v-model="searchQuery" placeholder="Search tokens..." class="pl-9 h-9" />
          </div>
        </div>
      </CardHeader>
      <CardContent class="p-0">
        <div v-if="loading && tokens.length === 0" class="flex justify-center py-20">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
        </div>

        <div v-else-if="tokens.length === 0" class="flex flex-col items-center py-20 text-center">
          <KeyRound class="w-12 h-12 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-bold text-foreground">No API tokens</p>
          <p class="text-xs text-muted-foreground mt-1">Generate a token to authenticate API requests.</p>
        </div>

        <!-- Token Table -->
        <table v-else class="w-full text-sm">
          <thead>
            <tr class="border-b border-border/40 text-xs text-muted-foreground">
              <th class="text-left font-medium px-4 py-3">Name</th>
              <th class="text-left font-medium px-4 py-3 hidden sm:table-cell">Token</th>
              <th class="text-left font-medium px-4 py-3 hidden md:table-cell">Last Used</th>
              <th class="text-left font-medium px-4 py-3 hidden md:table-cell">Expires</th>
              <th class="text-left font-medium px-4 py-3">Status</th>
              <th class="text-right font-medium px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="token in tokens" :key="token.id" class="token-row border-b border-border/20 hover:bg-muted/30 transition-colors">
              <td class="px-4 py-3">
                <span class="font-medium text-foreground">{{ token.name }}</span>
              </td>
              <td class="px-4 py-3 hidden sm:table-cell">
                <code class="text-xs bg-muted px-2 py-1 rounded font-mono">{{ token.token_prefix }}...</code>
              </td>
              <td class="px-4 py-3 hidden md:table-cell text-muted-foreground text-xs">
                {{ formatDate(token.last_used_at) }}
              </td>
              <td class="px-4 py-3 hidden md:table-cell text-xs">
                <span v-if="token.expires_at" class="text-muted-foreground">
                  {{ formatDate(token.expires_at) }}
                </span>
                <span v-else class="text-muted-foreground/50">Never</span>
              </td>
              <td class="px-4 py-3">
                <Badge v-if="token.is_revoked" variant="destructive" class="text-[10px]">Revoked</Badge>
                <Badge v-else variant="secondary" class="text-[10px] bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20">Active</Badge>
              </td>
              <td class="px-4 py-3 text-right">
                <Button v-if="!token.is_revoked" variant="ghost" size="sm" class="h-8 w-8 p-0 text-destructive hover:bg-destructive/10" @click="openRevoke(token)">
                  <Trash2 class="w-4 h-4" />
                </Button>
              </td>
            </tr>
          </tbody>
        </table>
      </CardContent>
    </Card>

    <!-- Create Token Dialog -->
    <Dialog v-model:open="isCreateOpen">
      <DialogContent class="sm:max-w-106.25">
        <DialogHeader>
          <DialogTitle>Generate API Token</DialogTitle>
          <DialogDescription>Create a personal access token for API authentication.</DialogDescription>
        </DialogHeader>
        <form @submit.prevent="handleCreate" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="tname">Token Name</Label>
            <Input id="tname" v-model="formName" placeholder="e.g. CI/CD Pipeline" required maxlength="100" />
          </div>
          <div class="space-y-2">
            <Label for="texpires">Expiry Date (optional)</Label>
            <Input id="texpires" v-model="formExpiresAt" type="date" />
          </div>
          <DialogFooter class="pt-4 border-t border-border/40">
            <Button type="button" variant="outline" @click="isCreateOpen = false">Cancel</Button>
            <Button type="submit" :disabled="formLoading">
              <Loader2 v-if="formLoading" class="w-4 h-4 mr-1.5 animate-spin" />
              Generate
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>

    <!-- Token Display Dialog -->
    <Dialog v-model:open="isTokenDisplayOpen">
      <DialogContent class="sm:max-w-106.25">
        <DialogHeader>
          <DialogTitle class="flex items-center gap-2">
            <CheckCircle class="w-5 h-5 text-emerald-500" />
            Token Created
          </DialogTitle>
          <DialogDescription>
            Copy this token now. <strong class="text-destructive">You will not be able to see it again.</strong>
          </DialogDescription>
        </DialogHeader>
        <div v-if="lastCreatedToken" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label>Token</Label>
            <div class="relative">
              <code class="block w-full text-xs font-mono bg-muted p-4 pr-12 rounded-lg border border-border/40 break-all select-all">
                {{ lastCreatedToken.raw_token }}
              </code>
              <Button
                variant="ghost"
                size="sm"
                class="absolute top-2 right-2 h-8"
                @click="copyToken"
              >
                <Copy v-if="!copied" class="w-4 h-4" />
                <CheckCircle v-else class="w-4 h-4 text-emerald-500" />
              </Button>
            </div>
          </div>
          <div class="bg-destructive/5 border border-destructive/20 rounded-lg p-3">
            <p class="text-xs text-destructive font-medium">
              This token will not be shown again. Store it securely.
            </p>
          </div>
        </div>
        <DialogFooter>
          <Button @click="isTokenDisplayOpen = false">Done</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Revoke Confirmation Dialog -->
    <Dialog v-model:open="isRevokeOpen">
      <DialogContent class="sm:max-w-100">
        <DialogHeader>
          <DialogTitle class="text-destructive">Revoke Token</DialogTitle>
          <DialogDescription>
            Revoke "{{ actionToken?.name }}"? Any services using this token will immediately lose access.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="isRevokeOpen = false">Cancel</Button>
          <Button variant="destructive" :disabled="formLoading" @click="handleRevoke">
            <Loader2 v-if="formLoading" class="w-4 h-4 mr-1.5 animate-spin" />
            Revoke
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
