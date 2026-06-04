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
import { sidebarContent } from '@/content/sidebar';
import { siteConfig } from '@/content/config';

const route = useRoute();

const pageTitle = computed(() => {
  const path = route.path;

  // Find a match in the sidebarContent using raw path
  for (const group of sidebarContent) {
    const matchedItem = group.items.find(item => item.href === path);
    if (matchedItem) {
      return matchedItem.title;
    }
  }

  return 'Dashboard';
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
        <header class="flex h-16 shrink-0 items-center justify-between gap-2 border-b border-border/50 px-6 bg-card/50 backdrop-blur-md sticky top-0 z-50">
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
