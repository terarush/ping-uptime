<script setup lang="ts">
import SidebarProvider from '@/components/ui/sidebar/SidebarProvider.vue';
import SidebarInset from '@/components/ui/sidebar/SidebarInset.vue';
import SidebarTrigger from '@/components/ui/sidebar/SidebarTrigger.vue';
import SidebarComponent from '@/components/sidebar.vue';
import Avatar from '@/components/ui/avatar/Avatar.vue';
import AvatarImage from '@/components/ui/avatar/AvatarImage.vue';
import AvatarFallback from '@/components/ui/avatar/AvatarFallback.vue';
import { useRoute } from 'vue-router';
import { computed } from 'vue';

const route = useRoute();

const pageTitle = computed(() => {
  const path = route.path;
  if (!path || path === '/' || path === '/dashboard') {
    return 'Dashboard';
  }

  const segments = path.split('/').filter(s => s && s !== 'dashboard');
  if (segments.length === 0) return 'Dashboard';

  const formattedSegments = segments.map(segment => {
    return segment
      .split('-')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ');
  });

  return formattedSegments.join(' / ');
});
</script>

<template>
  <SidebarProvider>
    <div class="flex min-h-screen w-full bg-background">
      <!-- App Sidebar -->
      <SidebarComponent />

      <!-- Content Area -->
      <SidebarInset class="flex flex-col flex-1">
        <!-- Top Navbar -->
        <header class="flex h-16 shrink-0 items-center justify-between gap-2 border-b border-border/50 px-6 bg-card/50 backdrop-blur-md sticky top-0 z-10">
          <div class="flex items-center gap-4">
            <SidebarTrigger />
            <div class="h-4 w-px bg-border/60"></div>
            <h1 class="text-sm font-bold text-foreground tracking-tight">{{ pageTitle }}</h1>
          </div>

          <!-- User Profile & Quick Actions -->
          <div class="flex items-center gap-4">
            <div class="flex flex-col text-right sm:flex">
              <span class="text-xs font-bold text-foreground">Administrator</span>
              <span class="text-[10px] text-muted-foreground">admin@ping-uptime.com</span>
            </div>
            <Avatar class="h-8 w-8 ring-2 ring-primary/10">
              <AvatarImage src="https://images.unsplash.com/photo-1534528741775-53994a69daeb?q=80&w=100&auto=format&fit=crop" alt="Admin" />
              <AvatarFallback>AD</AvatarFallback>
            </Avatar>
          </div>
        </header>

        <!-- Page Main Content -->
        <main class="flex-1 overflow-y-auto bg-slate-50/50 dark:bg-slate-950/30">
          <RouterView />
        </main>
      </SidebarInset>
    </div>
  </SidebarProvider>
</template>
