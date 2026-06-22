<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useSubscribers } from '@/composables/useSubscribers';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import {
  Mail,
  Search,
  Loader2,
  RefreshCw,
  Users,
  Trash2,
  CheckCircle2,
  XCircle,
} from '@lucide/vue';
import gsap from 'gsap';
import { toast } from 'vue-sonner';

const {
  subscribers,
  statusPages,
  loading,
  error,
  fetchSubscribers,
  fetchStatusPages,
  unsubscribe,
} = useSubscribers();

const selectedPageID = ref<number | null>(null);
const searchQuery = ref('');
const loadingAction = ref(false);

const filteredSubscribers = computed(() => {
  if (!subscribers.value) return [];
  const query = searchQuery.value.toLowerCase().trim();
  if (!query) return subscribers.value;
  return subscribers.value.filter(s =>
    s.email.toLowerCase().includes(query)
  );
});

const animateRows = () => {
  setTimeout(() => {
    gsap.fromTo('.sub-row',
      { opacity: 0, y: 12 },
      { opacity: 1, y: 0, duration: 0.3, stagger: 0.03, ease: 'power2.out' }
    );
  }, 50);
};

const loadSubscribers = async () => {
  if (!selectedPageID.value) return;
  loadingAction.value = true;
  try {
    await fetchSubscribers(selectedPageID.value);
    animateRows();
  } catch {
    // handled in store
  } finally {
    loadingAction.value = false;
  }
};

const handleUnsubscribe = async (email: string) => {
  if (!selectedPageID.value) return;
  try {
    await unsubscribe(selectedPageID.value, email);
    toast.success('Unsubscribed successfully');
  } catch {
    toast.error('Failed to unsubscribe');
  }
};

onMounted(async () => {
  await fetchStatusPages();
  if (statusPages.value.length > 0) {
    selectedPageID.value = statusPages.value[0]!.id;
    await loadSubscribers();
  }
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-emerald-500/5 dark:bg-emerald-500/2 blur-[100px] pointer-events-none"></div>

    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <Mail class="w-6 h-6 text-primary" />
          <span>Subscribers</span>
        </h2>
        <p class="text-xs text-muted-foreground">Manage email subscribers for your status pages.</p>
      </div>

      <div class="flex items-center gap-3">
        <Select v-model="selectedPageID" @update:model-value="loadSubscribers">
          <SelectTrigger class="w-48 h-9">
            <SelectValue placeholder="Select status page" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="page in statusPages" :key="page.id" :value="page.id">
              {{ page.name }}
            </SelectItem>
          </SelectContent>
        </Select>

        <Button variant="outline" size="sm" @click="loadSubscribers" class="h-9">
          <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loadingAction }" />
          <span>Refresh</span>
        </Button>
      </div>
    </div>

    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <CardTitle class="text-sm font-bold text-foreground flex items-center gap-2">
              <Users class="w-4 h-4" />
              <span>{{ subscribers.length }} verified subscribers</span>
            </CardTitle>
            <CardDescription class="text-xs">Subscribers to the selected status page.</CardDescription>
          </div>

          <div class="relative w-full sm:w-72">
            <Search class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              v-model="searchQuery"
              placeholder="Search by email..."
              class="pl-9 h-9"
            />
          </div>
        </div>
      </CardHeader>

      <CardContent class="p-0">
        <div v-if="loading && subscribers.length === 0" class="flex flex-col items-center justify-center py-20 gap-3">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
          <p class="text-sm text-muted-foreground">Loading subscribers...</p>
        </div>

        <div v-else-if="!selectedPageID" class="flex flex-col items-center justify-center py-20 text-center">
          <Mail class="w-12 h-12 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-bold text-foreground">Select a status page</p>
          <p class="text-xs text-muted-foreground mt-1">Choose a status page above to view its subscribers.</p>
        </div>

        <div v-else-if="filteredSubscribers.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
          <Users class="w-12 h-12 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-bold text-foreground">No subscribers</p>
          <p class="text-xs text-muted-foreground mt-1">This status page has no verified subscribers yet.</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b border-border/40 text-xs text-muted-foreground">
                <th class="text-left font-medium px-4 py-3">Email</th>
                <th class="text-center font-medium px-4 py-3">Verified</th>
                <th class="text-right font-medium px-4 py-3">Subscribed</th>
                <th class="text-right font-medium px-4 py-3 w-20">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="sub in filteredSubscribers"
                :key="sub.id"
                class="sub-row border-b border-border/20 hover:bg-accent/30 transition-colors"
              >
                <td class="px-4 py-3 text-sm">{{ sub.email }}</td>
                <td class="px-4 py-3 text-center">
                  <CheckCircle2 v-if="sub.verified" class="w-4 h-4 text-emerald-500 inline" />
                  <XCircle v-else class="w-4 h-4 text-muted-foreground/50 inline" />
                </td>
                <td class="px-4 py-3 text-sm text-right text-muted-foreground">
                  {{ new Date(sub.created_at).toLocaleDateString() }}
                </td>
                <td class="px-4 py-3 text-right">
                  <Button variant="ghost" size="icon" class="h-8 w-8 text-red-500 hover:text-red-600 hover:bg-red-500/10" @click="handleUnsubscribe(sub.email)">
                    <Trash2 class="w-4 h-4" />
                  </Button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
