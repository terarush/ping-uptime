<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { useAnalytics } from '@/composables/useAnalytics'
import { useMonitors } from '@/composables/useMonitors'
import { useSslMonitors } from '@/composables/useSslMonitors'
import { useIncidents } from '@/composables/useIncidents'
import { useNotificationChannels } from '@/composables/useNotificationChannels'
import { useIntegrations } from '@/composables/useIntegrations'
import { useStatusPages } from '@/composables/useStatusPages'
import { useMaintenances } from '@/composables/useMaintenances'
import { useTags } from '@/composables/useTags'
import { useApiTokens } from '@/composables/useApiTokens'
import { useTeams } from '@/composables/useTeams'
import { useUsers } from '@/composables/useUsers'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select'
import {
  Activity, Bell, ShieldCheck, Globe, FileText, Wrench, Tag, KeyRound,
  Users, UsersRound, Link, BarChart3, Loader2, RefreshCw, AlertCircle,
  CheckCircle2, User, TrendingUp, Clock, Search, X, Settings,
} from '@lucide/vue'
import { useRouter } from 'vue-router'
import { siteConfig } from '@/content/config'
import gsap from 'gsap'
import MultiChart, { type Series } from '@/components/multi-chart.vue'
import ChartVue from '@/components/chart.vue'
import type { DataPoint } from '@/components/chart.vue'
import AnalyticsTable from '@/components/analytics-table.vue'

const { currentUser } = useAuth()
const isAdmin = computed(() => currentUser.value?.role === 'admin')
const router = useRouter()

// -- data sources --
const { monitors, fetchMonitors, loading: monitorsLoading } = useMonitors()
const { certs, fetchAll: fetchSsl, loading: sslLoading } = useSslMonitors()
const { incidents, fetchIncidentsData, loading: incLoading } = useIncidents()
const { channels, fetchChannelsData, loading: notifChannelsLoading } = useNotificationChannels()
const { integrations, fetchAll: fetchIntegrations, loading: intLoading } = useIntegrations()
const { statusPages, fetchStatusPagesData, loading: spLoading } = useStatusPages()
const { items: maintenances, fetchAll: fetchMaintenances, loading: maintLoading } = useMaintenances()
const { tags, fetchAll: fetchTags, loading: tagsLoading } = useTags()
const { tokens, fetchAll: fetchTokens, loading: tokensLoading } = useApiTokens()
const { teams, fetchAll: fetchTeams, loading: teamsLoading } = useTeams()
const { users, fetchUsersData, loading: usersLoading } = useUsers()
const { stats, loading: analyticsLoading, fetchDashboardStats } = useAnalytics()

const loading = computed(() =>
  monitorsLoading.value || sslLoading.value || incLoading.value || analyticsLoading.value ||
  notifChannelsLoading.value || intLoading.value || spLoading.value || maintLoading.value ||
  tagsLoading.value || tokensLoading.value || teamsLoading.value || usersLoading.value
)

const getMonitorOwner = (monitorId: number) => {
  const mon = monitors.value.find(m => m.id === monitorId)
  if (!mon) return null
  const user = users.value.find(u => u.id === mon.user_id)
  return user ? `${user.name} (${user.email})` : `User #${mon.user_id}`
}

// -- time window --
const selectedWindow = ref<'1h' | '1d' | '1w' | '1m' | '1y' | 'all'>('1w')

// -- aggregate KPIs --
const totalEndpoints = computed(() => monitors.value.length)
const activeEndpoints = computed(() => monitors.value.filter(m => m.status === 'active').length)
const upCount = computed(() => monitors.value.filter(m => m.status === 'active' && m.uptime_status === 'up').length)
const downCount = computed(() => monitors.value.filter(m => m.status === 'active' && m.uptime_status === 'down').length)
const pausedCount = computed(() => monitors.value.filter(m => m.status === 'paused').length)

const avgUptime = computed(() => {
  if (!stats.value || stats.value.length === 0) return '—'
  const sum = stats.value.reduce((a, s) => a + s.uptime_pct, 0)
  return `${(sum / stats.value.length).toFixed(2)}%`
})

const avgLatency = computed(() => {
  if (!stats.value || stats.value.length === 0) return null
  const sum = stats.value.reduce((a, s) => a + s.avg_latency, 0)
  return Math.round(sum / stats.value.length)
})

const activeIncidents = computed(() => incidents.value.filter(i => i.status === 'active').length)
const certsExpiring = computed(() => certs.value.filter(c => c.status === 'expiring_soon' || c.status === 'expired').length)

const statusDistribution = computed(() => {
  const dist = { operational: 0, degraded: 0, outage: 0 }
  stats.value?.forEach(s => {
    if (s.status === 'operational') dist.operational++
    else if (s.status === 'degraded') dist.degraded++
    else if (s.status === 'outage') dist.outage++
  })
  return dist
})

// -- system health chart (aggregate uptime over time) --
const healthChartPoints = computed<DataPoint[]>(() => {
  if (!stats.value) return []
  const byDate = new Map<string, { uptime: number[], latency: number[], failed: number, total: number }>()
  stats.value.forEach(s => {
    (s.points || []).forEach(p => {
      if (!p.date) return
      if (!byDate.has(p.date)) byDate.set(p.date, { uptime: [], latency: [], failed: 0, total: 0 })
      const bucket = byDate.get(p.date)!
      if (p.status === 'up' && p.uptime_pct > 0) bucket.uptime.push(p.uptime_pct)
      if (p.latency > 0) bucket.latency.push(p.latency)
      bucket.failed += p.failed || 0
      bucket.total += p.total || 0
    })
  })
  return Array.from(byDate.entries()).sort(([a], [b]) => a.localeCompare(b)).map(([date, d]) => ({
    date,
    status: d.uptime.length > 0 ? 'up' : 'unknown',
    uptime_pct: d.uptime.length > 0 ? d.uptime.reduce((a, b) => a + b, 0) / d.uptime.length : 0,
    latency: d.latency.length > 0 ? d.latency.reduce((a, b) => a + b, 0) / d.latency.length : 0,
    failed: d.failed,
    total: d.total,
  }))
})

// -- multi-chart latency series --
const searchQuery = ref('')
const enabledMonitors = ref<Set<number>>(new Set())

const colorsPalette = [
  '#6366f1', '#10b981', '#3b82f6', '#f59e0b', '#ec4899',
  '#8b5cf6', '#06b6d4', '#f97316', '#14b8a6', '#ef4444',
]
const getMonitorColor = (id: number, idx: number): string => colorsPalette[idx % colorsPalette.length] || '#6366f1'

watch(stats, (s) => {
  if (s && s.length > 0 && enabledMonitors.value.size === 0) {
    s.forEach(st => enabledMonitors.value.add(st.monitor_id))
  }
}, { immediate: true })

const toggleMonitor = (id: number) => {
  if (enabledMonitors.value.has(id)) enabledMonitors.value.delete(id)
  else enabledMonitors.value.add(id)
}
const selectAllMonitors = () => stats.value?.forEach(s => enabledMonitors.value.add(s.monitor_id))
const deselectAllMonitors = () => enabledMonitors.value.clear()

const filteredStats = computed(() => {
  if (!stats.value) return []
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return stats.value
  return stats.value.filter(s =>
    s.monitor_name.toLowerCase().includes(q) || s.monitor_url.toLowerCase().includes(q)
  )
})

const chartSeries = computed<Series[]>(() => {
  if (!stats.value) return []
  return stats.value
    .filter(s => enabledMonitors.value.has(s.monitor_id))
    .map((s, idx) => ({
      monitorId: s.monitor_id,
      name: s.monitor_name,
      color: getMonitorColor(s.monitor_id, idx),
      points: s.points || [],
    }))
})

// -- recent incidents for timeline --
const recentIncidents = computed(() => [...incidents.value].slice(0, 8))

const formatRelative = (dateStr: string) => {
  const d = new Date(dateStr)
  const s = Math.floor((Date.now() - d.getTime()) / 1000)
  if (s < 60) return `${s}s ago`
  const m = Math.floor(s / 60)
  if (m < 60) return `${m}m ago`
  const h = Math.floor(m / 60)
  if (h < 24) return `${h}h ago`
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

const navigateTo = (path: string) => router.push(`${siteConfig.appPath}/${path}`)

const fetchAllData = async () => {
  await Promise.all([
    fetchMonitors(),
    fetchSsl(),
    fetchIncidentsData(),
    fetchChannelsData(),
    fetchIntegrations(),
    fetchStatusPagesData(),
    fetchMaintenances(),
    fetchTags(),
    fetchTokens(),
    fetchTeams(),
    isAdmin.value ? fetchUsersData() : Promise.resolve(),
    fetchDashboardStats(selectedWindow.value),
  ])
  setTimeout(animateCards, 50)
}

watch(selectedWindow, async (w) => {
  await fetchDashboardStats(w)
})

const animateCards = () => {
  gsap.fromTo('.ov-card', { opacity: 0, y: 12 }, { opacity: 1, y: 0, duration: 0.3, stagger: 0.03, ease: 'power2.out' })
}

onMounted(fetchAllData)
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-emerald-500/10 dark:bg-emerald-500/5 blur-[100px] pointer-events-none"></div>

    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <BarChart3 class="w-6 h-6 text-primary" />
          <span>System Overview</span>
        </h2>
        <p class="text-xs text-muted-foreground">
          Aggregate health, trends, and performance analysis across all endpoints.
        </p>
      </div>
      <div class="flex items-center gap-2">
        <Select v-model="selectedWindow">
          <SelectTrigger class="h-9 w-37.5">
            <SelectValue placeholder="Window" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="1h">Last hour</SelectItem>
            <SelectItem value="1d">Last 24 hours</SelectItem>
            <SelectItem value="1w">Last 7 days</SelectItem>
            <SelectItem value="1m">Last 30 days</SelectItem>
            <SelectItem value="1y">Last year</SelectItem>
            <SelectItem value="all">All time</SelectItem>
          </SelectContent>
        </Select>
        <Button variant="outline" size="sm" @click="fetchAllData" class="h-9">
          <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
          <span>Refresh</span>
        </Button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading && monitors.length === 0" class="flex items-center justify-center py-24">
      <Loader2 class="w-8 h-8 text-primary animate-spin" />
    </div>

    <template v-else>
      <!-- ── KPI row ── -->
      <div class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-6 gap-3">
        <Card class="ov-card border-border/50 bg-card/60 backdrop-blur-md">
          <CardContent class="p-4 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold text-muted-foreground uppercase">Total Endpoints</p>
              <p class="text-xl font-black text-foreground mt-0.5">{{ totalEndpoints }}</p>
              <p class="text-[9px] text-muted-foreground">{{ activeEndpoints }} active</p>
            </div>
            <div class="p-2 bg-muted rounded-lg"><Activity class="w-4 h-4 text-muted-foreground" /></div>
          </CardContent>
        </Card>
        <Card class="ov-card border-border/50 bg-card/60 backdrop-blur-md">
          <CardContent class="p-4 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold text-muted-foreground uppercase">Avg Uptime</p>
              <p class="text-xl font-black text-emerald-500 mt-0.5">{{ avgUptime }}</p>
              <p class="text-[9px] text-muted-foreground">{{ upCount }}/{{ activeEndpoints }} up</p>
            </div>
            <div class="p-2 bg-emerald-500/10 rounded-lg"><CheckCircle2 class="w-4 h-4 text-emerald-500" /></div>
          </CardContent>
        </Card>
        <Card class="ov-card border-border/50 bg-card/60 backdrop-blur-md">
          <CardContent class="p-4 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold text-muted-foreground uppercase">Avg Latency</p>
              <p class="text-xl font-black text-blue-500 mt-0.5">{{ avgLatency ?? '—' }}<span v-if="avgLatency" class="text-sm">ms</span></p>
              <p class="text-[9px] text-muted-foreground">Across {{ selectedWindow }}</p>
            </div>
            <div class="p-2 bg-blue-500/10 rounded-lg"><Clock class="w-4 h-4 text-blue-500" /></div>
          </CardContent>
        </Card>
        <Card class="ov-card border-border/50 bg-card/60 backdrop-blur-md">
          <CardContent class="p-4 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold text-muted-foreground uppercase">Down</p>
              <p class="text-xl font-black mt-0.5" :class="downCount > 0 ? 'text-rose-500' : 'text-foreground'">{{ downCount }}</p>
              <p class="text-[9px]" :class="downCount > 0 ? 'text-rose-500/80' : 'text-muted-foreground'">{{ downCount > 0 ? 'Outage' : 'All clear' }}</p>
            </div>
            <div class="p-2 rounded-lg" :class="downCount > 0 ? 'bg-rose-500/10' : 'bg-muted'">
              <AlertCircle class="w-4 h-4" :class="downCount > 0 ? 'text-rose-500' : 'text-muted-foreground'" />
            </div>
          </CardContent>
        </Card>
        <Card class="ov-card border-border/50 bg-card/60 backdrop-blur-md">
          <CardContent class="p-4 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold text-muted-foreground uppercase">Active Incidents</p>
              <p class="text-xl font-black mt-0.5" :class="activeIncidents > 0 ? 'text-rose-500' : 'text-foreground'">{{ activeIncidents }}</p>
              <p class="text-[9px] text-muted-foreground">{{ incidents.length }} total</p>
            </div>
            <div class="p-2 rounded-lg" :class="activeIncidents > 0 ? 'bg-rose-500/10' : 'bg-muted'">
              <AlertCircle class="w-4 h-4" :class="activeIncidents > 0 ? 'text-rose-500' : 'text-muted-foreground'" />
            </div>
          </CardContent>
        </Card>
        <Card class="ov-card border-border/50 bg-card/60 backdrop-blur-md">
          <CardContent class="p-4 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold text-muted-foreground uppercase">SSL Expiring</p>
              <p class="text-xl font-black mt-0.5" :class="certsExpiring > 0 ? 'text-amber-500' : 'text-foreground'">{{ certsExpiring }}</p>
              <p class="text-[9px] text-muted-foreground">{{ certs.length }} total certs</p>
            </div>
            <div class="p-2 rounded-lg" :class="certsExpiring > 0 ? 'bg-amber-500/10' : 'bg-muted'">
              <ShieldCheck class="w-4 h-4" :class="certsExpiring > 0 ? 'text-amber-500' : 'text-muted-foreground'" />
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- ── Status Distribution + Incident Timeline row ── -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card class="ov-card border-border/50 bg-card/60 backdrop-blur-md lg:col-span-1">
          <CardHeader class="pb-3 border-b border-border/40">
            <CardTitle class="text-sm font-bold">Status Distribution</CardTitle>
            <CardDescription class="text-xs">Operational status across all endpoints.</CardDescription>
          </CardHeader>
          <CardContent class="p-5 space-y-4">
            <div class="flex items-center justify-between p-3 rounded-lg bg-emerald-500/5 border border-emerald-500/10">
              <span class="text-xs font-bold text-emerald-500">Operational</span>
              <span class="text-lg font-black text-emerald-500">{{ statusDistribution.operational }}</span>
            </div>
            <div class="flex items-center justify-between p-3 rounded-lg bg-amber-500/5 border border-amber-500/10">
              <span class="text-xs font-bold text-amber-500">Degraded</span>
              <span class="text-lg font-black text-amber-500">{{ statusDistribution.degraded }}</span>
            </div>
            <div class="flex items-center justify-between p-3 rounded-lg bg-rose-500/5 border border-rose-500/10">
              <span class="text-xs font-bold text-rose-500">Outage</span>
              <span class="text-lg font-black text-rose-500">{{ statusDistribution.outage }}</span>
            </div>
          </CardContent>
        </Card>

        <!-- System Health Chart -->
        <Card class="ov-card border-border/50 bg-card/60 backdrop-blur-md lg:col-span-2">
          <CardHeader class="pb-3 border-b border-border/40">
            <div class="flex items-center justify-between">
              <div>
                <CardTitle class="text-sm font-bold flex items-center gap-2">
                  <TrendingUp class="w-4 h-4 text-primary" />
                  <span>System Health Trend</span>
                </CardTitle>
                <CardDescription class="text-xs">Aggregate uptime percentage over time.</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent class="p-4">
            <div v-if="healthChartPoints.length === 0" class="flex items-center justify-center py-12 text-xs text-muted-foreground">
              No health data for this window.
            </div>
            <ChartVue v-else :points="healthChartPoints" :height="200" />
          </CardContent>
        </Card>
      </div>

      <!-- ── Latency Multi-Chart ── -->
      <Card class="ov-card border-border/50 bg-card/60 backdrop-blur-md">
        <CardHeader class="pb-3 border-b border-border/40">
          <div>
            <CardTitle class="text-sm font-bold flex items-center gap-2">
              <Activity class="w-4 h-4 text-primary" />
              <span>Latency Comparison</span>
            </CardTitle>
            <CardDescription class="text-xs">Per-monitor latency over time. Toggle monitors below.</CardDescription>
          </div>
        </CardHeader>
        <CardContent class="p-6">
          <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
            <div class="lg:col-span-8">
              <div class="relative bg-muted/20 border border-border/30 rounded-xl p-4 min-h-[300px] flex items-center justify-center">
                <div v-if="analyticsLoading && stats?.length === 0" class="flex flex-col items-center gap-2">
                  <Loader2 class="w-6 h-6 text-primary animate-spin" />
                  <span class="text-xs text-muted-foreground">Loading chart...</span>
                </div>
                <div v-else-if="chartSeries.length === 0" class="flex flex-col items-center text-center text-muted-foreground p-8">
                  <X class="w-8 h-8 mb-2 text-muted-foreground/40" />
                  <p class="text-sm font-bold">No monitors selected</p>
                  <p class="text-xs mt-1">Select monitors from the checklist.</p>
                </div>
                <MultiChart v-else :seriesList="chartSeries" :height="320" />
              </div>
              <div v-if="chartSeries.length > 0" class="flex flex-wrap gap-x-4 gap-y-2 text-[11px] font-medium pt-2 border-t border-border/30 mt-2">
                <div v-for="s in chartSeries" :key="s.monitorId" class="flex items-center gap-1.5">
                  <span class="w-2.5 h-2.5 rounded-full shrink-0" :style="{ backgroundColor: s.color }"></span>
                  <span class="text-foreground truncate max-w-30">{{ s.name }}</span>
                </div>
              </div>
            </div>
            <div class="lg:col-span-4 border border-border/30 rounded-xl p-4 bg-muted/10 space-y-4">
              <div class="flex items-center justify-between border-b border-border/40 pb-2">
                <h4 class="text-xs font-black uppercase tracking-wider text-foreground">Monitors</h4>
                <div class="flex items-center gap-2">
                  <button @click="selectAllMonitors" class="text-[10px] text-primary hover:underline font-bold cursor-pointer">All</button>
                  <span class="text-muted-foreground/30 text-[10px]">·</span>
                  <button @click="deselectAllMonitors" class="text-[10px] text-rose-500 hover:underline font-bold cursor-pointer">Clear</button>
                </div>
              </div>
              <div class="relative w-full">
                <Search class="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-muted-foreground" />
                <Input v-model="searchQuery" placeholder="Search..." class="pl-8 h-8 text-xs bg-card" />
              </div>
              <div class="space-y-1 max-h-[220px] overflow-y-auto pr-1">
                <div v-if="analyticsLoading && stats?.length === 0" class="flex items-center justify-center py-10">
                  <Loader2 class="w-5 h-5 text-primary animate-spin" />
                </div>
                <div v-else-if="filteredStats.length === 0" class="text-center text-xs text-muted-foreground py-10">No matches</div>
                <label v-else v-for="(item, idx) in filteredStats" :key="item.monitor_id"
                  class="flex items-center justify-between p-2 rounded-lg hover:bg-muted/50 border border-transparent hover:border-border/30 cursor-pointer select-none transition-colors">
                  <div class="flex items-center gap-2 min-w-0">
                    <input type="checkbox" :checked="enabledMonitors.has(item.monitor_id)"
                      @change="toggleMonitor(item.monitor_id)"
                      class="rounded border-input text-primary focus:ring-primary h-3.5 w-3.5 shrink-0 cursor-pointer" />
                    <span class="w-2.5 h-2.5 rounded-full shrink-0" :style="{ backgroundColor: getMonitorColor(item.monitor_id, idx) }"></span>
                    <span class="text-xs font-bold text-foreground truncate max-w-[120px]">{{ item.monitor_name }}</span>
                  </div>
                  <span class="text-[10px] text-muted-foreground font-semibold shrink-0">{{ Math.round(item.avg_latency) }} ms</span>
                </label>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- ── Resource Navigation ── -->
      <Card class="ov-card border-border/50 bg-card/60 backdrop-blur-md">
        <CardHeader class="pb-3 border-b border-border/40">
          <CardTitle class="text-sm font-bold">Quick Navigation</CardTitle>
          <CardDescription class="text-xs">Jump to any management section.</CardDescription>
        </CardHeader>
        <CardContent class="p-4">
          <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-2">
            <button type="button" @click="navigateTo('monitors')" class="flex items-center gap-2 p-2.5 rounded-lg border border-border/30 bg-muted/10 hover:bg-muted/20 transition-colors text-left">
              <div class="p-1.5 rounded-md shrink-0 bg-emerald-500/10"><Activity class="w-3.5 h-3.5 text-emerald-500" /></div>
              <div class="min-w-0">
                <p class="text-[11px] font-bold truncate">Monitors</p>
                <p class="text-[9px] text-muted-foreground">{{ monitors.length }}</p>
              </div>
            </button>
            <button type="button" @click="navigateTo('ssl-monitors')" class="flex items-center gap-2 p-2.5 rounded-lg border border-border/30 bg-muted/10 hover:bg-muted/20 transition-colors text-left">
              <div class="p-1.5 rounded-md shrink-0 bg-blue-500/10"><ShieldCheck class="w-3.5 h-3.5 text-blue-500" /></div>
              <div class="min-w-0">
                <p class="text-[11px] font-bold truncate">SSL</p>
                <p class="text-[9px] text-muted-foreground">{{ certs.length }}</p>
              </div>
            </button>
            <button type="button" @click="navigateTo('status-pages')" class="flex items-center gap-2 p-2.5 rounded-lg border border-border/30 bg-muted/10 hover:bg-muted/20 transition-colors text-left">
              <div class="p-1.5 rounded-md shrink-0 bg-indigo-500/10"><Globe class="w-3.5 h-3.5 text-indigo-500" /></div>
              <div class="min-w-0">
                <p class="text-[11px] font-bold truncate">Status Pages</p>
                <p class="text-[9px] text-muted-foreground">{{ statusPages.length }}</p>
              </div>
            </button>
            <button type="button" @click="navigateTo('incidents')" class="flex items-center gap-2 p-2.5 rounded-lg border border-border/30 bg-muted/10 hover:bg-muted/20 transition-colors text-left">
              <div class="p-1.5 rounded-md shrink-0 bg-rose-500/10"><AlertCircle class="w-3.5 h-3.5 text-rose-500" /></div>
              <div class="min-w-0">
                <p class="text-[11px] font-bold truncate">Incidents</p>
                <p class="text-[9px] text-muted-foreground">{{ incidents.length }}</p>
              </div>
            </button>
            <button type="button" @click="navigateTo('maintenances')" class="flex items-center gap-2 p-2.5 rounded-lg border border-border/30 bg-muted/10 hover:bg-muted/20 transition-colors text-left">
              <div class="p-1.5 rounded-md shrink-0 bg-amber-500/10"><Wrench class="w-3.5 h-3.5 text-amber-500" /></div>
              <div class="min-w-0">
                <p class="text-[11px] font-bold truncate">Maintenances</p>
                <p class="text-[9px] text-muted-foreground">{{ maintenances.length }}</p>
              </div>
            </button>
            <button type="button" @click="navigateTo('settings/notifications')" class="flex items-center gap-2 p-2.5 rounded-lg border border-border/30 bg-muted/10 hover:bg-muted/20 transition-colors text-left">
              <div class="p-1.5 rounded-md shrink-0 bg-purple-500/10"><Bell class="w-3.5 h-3.5 text-purple-500" /></div>
              <div class="min-w-0">
                <p class="text-[11px] font-bold truncate">Channels</p>
                <p class="text-[9px] text-muted-foreground">{{ channels.length }}</p>
              </div>
            </button>
            <button type="button" @click="navigateTo('integrations')" class="flex items-center gap-2 p-2.5 rounded-lg border border-border/30 bg-muted/10 hover:bg-muted/20 transition-colors text-left">
              <div class="p-1.5 rounded-md shrink-0 bg-cyan-500/10"><Link class="w-3.5 h-3.5 text-cyan-500" /></div>
              <div class="min-w-0">
                <p class="text-[11px] font-bold truncate">Integrations</p>
                <p class="text-[9px] text-muted-foreground">{{ integrations.length }}</p>
              </div>
            </button>
            <button type="button" @click="navigateTo('tags')" class="flex items-center gap-2 p-2.5 rounded-lg border border-border/30 bg-muted/10 hover:bg-muted/20 transition-colors text-left">
              <div class="p-1.5 rounded-md shrink-0 bg-pink-500/10"><Tag class="w-3.5 h-3.5 text-pink-500" /></div>
              <div class="min-w-0">
                <p class="text-[11px] font-bold truncate">Tags</p>
                <p class="text-[9px] text-muted-foreground">{{ tags.length }}</p>
              </div>
            </button>
            <button type="button" @click="navigateTo('api-tokens')" class="flex items-center gap-2 p-2.5 rounded-lg border border-border/30 bg-muted/10 hover:bg-muted/20 transition-colors text-left">
              <div class="p-1.5 rounded-md shrink-0 bg-orange-500/10"><KeyRound class="w-3.5 h-3.5 text-orange-500" /></div>
              <div class="min-w-0">
                <p class="text-[11px] font-bold truncate">API Tokens</p>
                <p class="text-[9px] text-muted-foreground">{{ tokens.length }}</p>
              </div>
            </button>
            <button v-if="isAdmin" type="button" @click="navigateTo('users')" class="flex items-center gap-2 p-2.5 rounded-lg border border-border/30 bg-muted/10 hover:bg-muted/20 transition-colors text-left">
              <div class="p-1.5 rounded-md shrink-0 bg-sky-500/10"><Users class="w-3.5 h-3.5 text-sky-500" /></div>
              <div class="min-w-0">
                <p class="text-[11px] font-bold truncate">Users</p>
                <p class="text-[9px] text-muted-foreground">{{ users.length }}</p>
              </div>
            </button>
            <button v-if="isAdmin" type="button" @click="navigateTo('teams')" class="flex items-center gap-2 p-2.5 rounded-lg border border-border/30 bg-muted/10 hover:bg-muted/20 transition-colors text-left">
              <div class="p-1.5 rounded-md shrink-0 bg-violet-500/10"><UsersRound class="w-3.5 h-3.5 text-violet-500" /></div>
              <div class="min-w-0">
                <p class="text-[11px] font-bold truncate">Teams</p>
                <p class="text-[9px] text-muted-foreground">{{ teams.length }} ({{ teams.reduce((s, t) => s + (t.member_count || 0), 0) }} members)</p>
              </div>
            </button>
          </div>
        </CardContent>
      </Card>

      <!-- ── Incident Timeline + Endpoint Analytics Table ── -->
      <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">
        <!-- Incident timeline -->
        <Card class="ov-card border-border/50 bg-card/60 backdrop-blur-md lg:col-span-2">
          <CardHeader class="pb-3 border-b border-border/40">
            <div class="flex items-center justify-between">
              <div>
                <CardTitle class="text-sm font-bold">Incident Timeline</CardTitle>
                <CardDescription class="text-xs">Recent 8 events across monitors.</CardDescription>
              </div>
              <Button variant="ghost" size="sm" class="h-7 text-xs text-muted-foreground" @click="navigateTo('incidents')">
                View all
              </Button>
            </div>
          </CardHeader>
          <CardContent class="p-0 max-h-[280px] overflow-y-auto">
            <div v-if="incidents.length === 0" class="flex flex-col items-center justify-center py-14 text-center">
              <CheckCircle2 class="w-8 h-8 text-emerald-500/25 mb-2" />
              <p class="text-xs font-bold">No incidents</p>
              <p class="text-[10px] text-muted-foreground">All systems operational.</p>
            </div>
            <div v-else class="divide-y divide-border/30">
              <div v-for="inc in recentIncidents" :key="inc.id"
                class="flex items-center justify-between gap-3 px-4 py-3 hover:bg-muted/30 transition-colors">
                <div class="flex items-center gap-2.5 min-w-0">
                  <div class="shrink-0">
                    <div v-if="inc.status === 'resolved'" class="w-4.5 h-4.5 rounded-full bg-emerald-500/15 flex items-center justify-center">
                      <CheckCircle2 class="w-2.5 h-2.5 text-emerald-500" />
                    </div>
                    <div v-else class="w-4.5 h-4.5 rounded-full bg-rose-500/15 flex items-center justify-center">
                      <AlertCircle class="w-2.5 h-2.5 text-rose-500" />
                    </div>
                  </div>
                  <div class="min-w-0">
                    <p class="text-xs font-bold text-foreground truncate max-w-[180px]">Monitor #{{ inc.monitor_id }}</p>
                    <p class="text-[9px] font-semibold text-muted-foreground">
                      <span :class="inc.status === 'resolved' ? 'text-emerald-500' : 'text-rose-500'">{{ inc.status === 'resolved' ? 'Recovered' : 'Down' }}</span>
                      <span class="mx-1">·</span>
                      {{ formatRelative(inc.created_at) }}
                    </p>
                  </div>
                </div>
                <span v-if="inc.latency > 0" class="text-[10px] font-mono text-muted-foreground font-bold shrink-0">{{ inc.latency }}ms</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Endpoint analytics table -->
        <Card class="ov-card border-border/50 bg-card/60 backdrop-blur-md lg:col-span-3">
          <CardHeader class="pb-3 border-b border-border/40">
            <div class="flex items-center justify-between">
              <div>
                <CardTitle class="text-sm font-bold">Endpoint Performance</CardTitle>
                <CardDescription class="text-xs">Uptime, latency, and failure breakdown per monitor.</CardDescription>
              </div>
              <Button variant="ghost" size="sm" class="h-7 text-xs text-muted-foreground" @click="navigateTo('analytics')">
                Full analytics
              </Button>
            </div>
          </CardHeader>
          <CardContent class="p-0">
            <div v-if="analyticsLoading && stats?.length === 0" class="flex items-center justify-center py-12">
              <Loader2 class="w-6 h-6 text-primary animate-spin" />
            </div>
            <div v-else-if="!stats?.length" class="text-center text-xs text-muted-foreground py-12">
              No endpoint data for this window.
            </div>
            <AnalyticsTable v-else :stats="stats" :is-admin="isAdmin" :get-monitor-owner="getMonitorOwner" />
          </CardContent>
        </Card>
      </div>
    </template>
  </div>
</template>

<style scoped>
.ov-card {
  transition: transform 0.15s cubic-bezier(0.16, 1, 0.3, 1);
}
.ov-card:hover {
  transform: translateY(-1px);
}
</style>
