<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useIncidents } from '@/composables/useIncidents';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import IncidentTable from '@/components/incident-table.vue';
import {
  FileText,
  CheckCircle,
  Search,
  Loader2,
  RefreshCw
} from '@lucide/vue';
import gsap from 'gsap';

// Use external Incidents composable
const {
  incidents,
  monitorsMap,
  loading,
  error,
  fetchIncidentsData,
  formatDate,
  getDuration
} = useIncidents();

const searchQuery = ref('');

// Fetch both monitors and incidents log on mount
const fetchAll = async () => {
  try {
    await fetchIncidentsData();
  } catch (err) {
    console.error('Failed to load incidents logs:', err);
  } finally {
    setTimeout(animateTableRows, 50);
  }
};

// Filtered incidents
const filteredIncidents = computed(() => {
  if (!incidents.value) return [];
  const query = searchQuery.value.toLowerCase().trim();
  if (!query) return incidents.value;
  return incidents.value.filter(inc => {
    const monitor = monitorsMap.value[inc.monitor_id];
    const monitorName = monitor ? monitor.name.toLowerCase() : '';
    const monitorUrl = monitor ? monitor.url.toLowerCase() : '';
    const errMessage = inc.error_message ? inc.error_message.toLowerCase() : '';
    const statusStr = inc.status.toLowerCase();

    return monitorName.includes(query) ||
           monitorUrl.includes(query) ||
           errMessage.includes(query) ||
           statusStr.includes(query);
  });
});

// GSAP Animations
const animateTableRows = () => {
  gsap.fromTo('.incident-row',
    { opacity: 0, y: 12 },
    { opacity: 1, y: 0, duration: 0.4, stagger: 0.04, ease: 'power2.out' }
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
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-red-500/5 dark:bg-red-500/2 blur-[100px] pointer-events-none"></div>

    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <FileText class="w-6 h-6 text-primary" />
          <span>Incident Logs</span>
        </h2>
        <p class="text-xs text-muted-foreground">Historical records of checks that failed and recovery timeline logs.</p>
      </div>

      <!-- Actions -->
      <Button variant="outline" size="sm" @click="fetchAll" class="h-9">
        <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
        <span>Refresh</span>
      </Button>
    </div>

    <!-- Main Container Card -->
    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <CardTitle class="text-sm font-bold text-foreground">Outages & Alerts history</CardTitle>
            <CardDescription class="text-xs">Timeline of network failures and latency responses.</CardDescription>
          </div>

          <!-- Search Bar -->
          <div class="relative w-full sm:w-72">
            <Search class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              v-model="searchQuery"
              placeholder="Search incidents..."
              class="pl-9 h-9"
            />
          </div>
        </div>
      </CardHeader>

      <CardContent class="p-0">
        <div v-if="loading && incidents.length === 0" class="flex flex-col items-center justify-center py-20 gap-3">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
          <p class="text-sm text-muted-foreground">Loading incidents log...</p>
        </div>

        <div v-else-if="filteredIncidents.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
          <CheckCircle class="w-12 h-12 text-emerald-500/30 mb-3" />
          <p class="text-sm font-bold text-foreground">All systems operational</p>
          <p class="text-xs text-muted-foreground mt-1">No recorded incidents matching the filters.</p>
        </div>

        <IncidentTable
          v-else
          :incidents="filteredIncidents"
          :monitors-map="monitorsMap"
          :format-date="formatDate"
          :get-duration="getDuration"
        />
      </CardContent>
    </Card>
  </div>
</template>

