<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useBackup } from '@/composables/useBackup';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import {
  FolderArchive,
  Download,
  Upload,
  Loader2,
  RefreshCw,
  AlertTriangle,
  FileJson,
  Trash2,
  Clock,
} from '@lucide/vue';
import gsap from 'gsap';
import { toast } from 'vue-sonner';

const { records, loading, downloadBackup, importBackup, getHistory, deleteRecord } = useBackup();

const importing = ref(false);
const importFile = ref<File | null>(null);
const importConfirm = ref(false);
const dragOver = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);

const lastExportTime = computed(() => {
  if (records.value.length === 0) return null;
  const exportRecords = records.value.filter(r => r.file_name.startsWith('manual-export'));
  if (exportRecords.length === 0) return null;
  return exportRecords[0]?.created_at;
});

const handleDownload = async () => {
  try {
    await downloadBackup();
    toast.success('Backup downloaded');
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to download backup');
  }
};

const handleFileSelect = (e: Event) => {
  const target = e.target as HTMLInputElement;
  if (target.files && target.files.length > 0 && target.files[0]) {
    importFile.value = target.files[0];
  }
};

const handleDrop = (e: DragEvent) => {
  dragOver.value = false;
  if (e.dataTransfer?.files && e.dataTransfer.files.length > 0 && e.dataTransfer.files[0]) {
    importFile.value = e.dataTransfer.files[0];
  }
};

const handleImport = async () => {
  if (!importFile.value || !importConfirm.value) return;
  importing.value = true;
  try {
    await importBackup(importFile.value);
    toast.success('Configuration imported successfully');
    importFile.value = null;
    importConfirm.value = false;
    if (fileInput.value) fileInput.value.value = '';
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Import failed');
  } finally {
    importing.value = false;
  }
};

const handleDelete = async (id: number) => {
  try {
    await deleteRecord(id);
    toast.success('Backup record deleted');
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to delete record');
  }
};

const fmtSize = (bytes: number) => {
  if (bytes === 0) return '—';
  const units = ['B', 'KB', 'MB', 'GB'];
  let i = 0;
  let size = bytes;
  while (size >= 1024 && i < units.length - 1) {
    size /= 1024;
    i++;
  }
  return `${size.toFixed(1)} ${units[i]}`;
};

const fmtTime = (s: string) => new Date(s).toLocaleString();

const animateRows = () => {
  setTimeout(() => {
    gsap.fromTo('.backup-row',
      { opacity: 0, y: 12 },
      { opacity: 1, y: 0, duration: 0.4, stagger: 0.04, ease: 'power2.out' }
    );
  }, 50);
};

onMounted(async () => {
  await getHistory();
  animateRows();
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-amber-500/10 dark:bg-amber-500/5 blur-[100px] pointer-events-none"></div>

    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <FolderArchive class="w-6 h-6 text-primary" />
          <span>Backup & Export</span>
        </h2>
        <p class="text-xs text-muted-foreground">Manage full system configuration backups.</p>
      </div>
      <Button variant="outline" size="sm" @click="getHistory" class="h-9">
        <RefreshCw class="w-4 h-4 mr-1.5" :class="{ 'animate-spin': loading }" />
        Refresh
      </Button>
    </div>

    <!-- Warning Banner -->
    <Card class="border-amber-500/30 bg-amber-500/5 dark:bg-amber-500/10 z-10 relative">
      <CardContent class="p-4 flex items-start gap-3">
        <AlertTriangle class="w-5 h-5 text-amber-500 shrink-0 mt-0.5" />
        <div>
          <p class="text-sm font-bold text-amber-600 dark:text-amber-400">Warning</p>
          <p class="text-xs text-muted-foreground mt-1">This tool manages full system configuration backups. Importing a backup will overwrite existing configuration data including monitors, incidents, status pages, settings, notification channels, tags, integrations, SSL certificates, and API tokens.</p>
        </div>
      </CardContent>
    </Card>

    <!-- Export Section -->
    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <CardTitle class="text-sm font-bold flex items-center gap-2">
          <Download class="w-4 h-4 text-emerald-500" />
          <span>Export</span>
        </CardTitle>
        <CardDescription class="text-xs">Download a complete snapshot of your current configuration as a JSON file.</CardDescription>
      </CardHeader>
      <CardContent class="p-4">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <p class="text-sm text-muted-foreground">Export all system data including monitors, status pages, settings, and more.</p>
            <p v-if="lastExportTime" class="text-xs text-muted-foreground mt-1 flex items-center gap-1">
              <Clock class="w-3 h-3" />
              Last export: {{ fmtTime(lastExportTime) }}
            </p>
          </div>
          <Button @click="handleDownload" class="h-9 shrink-0">
            <Download class="w-4 h-4 mr-1.5" />
            Export Backup
          </Button>
        </div>
      </CardContent>
    </Card>

    <!-- Import Section -->
    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <CardTitle class="text-sm font-bold flex items-center gap-2">
          <Upload class="w-4 h-4 text-amber-500" />
          <span>Import</span>
        </CardTitle>
        <CardDescription class="text-xs">Restore configuration from a previously exported JSON file.</CardDescription>
      </CardHeader>
      <CardContent class="p-4 space-y-4">
        <!-- File Upload Area -->
        <div
          class="border-2 border-dashed rounded-lg p-8 text-center transition-colors cursor-pointer"
          :class="dragOver ? 'border-primary bg-primary/5' : 'border-border/50 hover:border-border'"
          @dragover.prevent="dragOver = true"
          @dragleave.prevent="dragOver = false"
          @drop.prevent="handleDrop"
          @click="fileInput?.click()"
        >
          <input
            ref="fileInput"
            type="file"
            accept=".json"
            class="hidden"
            @change="handleFileSelect"
          />
          <div v-if="!importFile" class="flex flex-col items-center gap-2">
            <FileJson class="w-10 h-10 text-muted-foreground/50" />
            <p class="text-sm font-medium text-foreground">Drop backup file here or click to browse</p>
            <p class="text-xs text-muted-foreground">JSON files only</p>
          </div>
          <div v-else class="flex flex-col items-center gap-2">
            <FileJson class="w-10 h-10 text-primary" />
            <p class="text-sm font-medium text-foreground">{{ importFile.name }}</p>
            <p class="text-xs text-muted-foreground">{{ fmtSize(importFile.size) }}</p>
            <Button variant="ghost" size="sm" class="h-7 text-xs" @click.stop="importFile = null; if (fileInput) fileInput.value = ''">
              Remove
            </Button>
          </div>
        </div>

        <!-- Confirmation Checkbox -->
        <div class="flex items-start gap-3 p-3 rounded-lg bg-destructive/5 border border-destructive/20">
          <input
            id="import-confirm"
            type="checkbox"
            v-model="importConfirm"
            class="mt-0.5 size-4 shrink-0 rounded-[4px] border border-destructive/40 accent-destructive"
          />
          <Label for="import-confirm" class="text-xs font-normal text-destructive cursor-pointer leading-relaxed">
            I understand this will overwrite existing configuration data. This action cannot be undone.
          </Label>
        </div>

        <Button
          @click="handleImport"
          :disabled="!importFile || !importConfirm || importing"
          class="h-9 w-full sm:w-auto"
          variant="destructive"
        >
          <Loader2 v-if="importing" class="w-4 h-4 mr-1.5 animate-spin" />
          <Upload v-else class="w-4 h-4 mr-1.5" />
          {{ importing ? 'Importing...' : 'Import Configuration' }}
        </Button>
      </CardContent>
    </Card>

    <!-- History Table -->
    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <CardTitle class="text-sm font-bold flex items-center gap-2">
          <Clock class="w-4 h-4 text-muted-foreground" />
          <span>Backup History</span>
        </CardTitle>
        <CardDescription class="text-xs">Record of all export and import operations.</CardDescription>
      </CardHeader>
      <CardContent class="p-0">
        <div v-if="loading && records.length === 0" class="flex flex-col items-center justify-center py-20 gap-3">
          <Loader2 class="w-8 h-8 text-primary animate-spin" />
          <p class="text-sm text-muted-foreground">Loading backup history...</p>
        </div>

        <div v-else-if="records.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
          <FolderArchive class="w-12 h-12 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-bold text-foreground">No backup history</p>
          <p class="text-xs text-muted-foreground mt-1">Export a backup to begin tracking history.</p>
        </div>

        <Table v-else>
          <TableHeader>
            <TableRow>
              <TableHead>File Name</TableHead>
              <TableHead>Size</TableHead>
              <TableHead>Date</TableHead>
              <TableHead class="text-right w-20">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="record in records" :key="record.id" class="backup-row text-xs hover:bg-muted/50 transition-colors">
              <TableCell class="font-medium">
                <span class="flex items-center gap-2">
                  <FileJson class="w-4 h-4 text-muted-foreground" />
                  {{ record.file_name }}
                </span>
              </TableCell>
              <TableCell class="text-muted-foreground">{{ fmtSize(record.file_size) }}</TableCell>
              <TableCell class="text-muted-foreground whitespace-nowrap">{{ fmtTime(record.created_at) }}</TableCell>
              <TableCell class="text-right">
                <Button variant="ghost" size="icon" class="h-7 w-7 text-destructive hover:text-destructive" @click="handleDelete(record.id)" title="Delete record">
                  <Trash2 class="h-3.5 w-3.5" />
                </Button>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </div>
</template>
