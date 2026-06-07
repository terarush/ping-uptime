<script setup lang="ts">
import { useRouter } from 'vue-router';
import { Button } from '@/components/ui/button';
import { ArrowLeft, Home, FileQuestion } from '@lucide/vue';
import { siteConfig } from '@/content/config';
import { useAuth } from '@/composables/useAuth';
import ThemeButton from '@/components/theme-button.vue';
import { onMounted } from 'vue';
import gsap from 'gsap';

const router = useRouter();
const { isAuthenticated } = useAuth();

const goBack = () => {
  router.go(-1);
};

const goHome = () => {
  if (isAuthenticated.value) {
    router.push(siteConfig.appPath);
  } else {
    router.push('/');
  }
};

onMounted(() => {
  // Orb entry animations
  gsap.fromTo('.ambient-orb-1',
    { opacity: 0, scale: 0.6 },
    { opacity: 0.5, scale: 1, duration: 2.2, ease: 'power3.out' }
  );

  gsap.fromTo('.ambient-orb-2',
    { opacity: 0, scale: 0.6 },
    { opacity: 0.4, scale: 1, duration: 2.5, ease: 'power3.out', delay: 0.2 }
  );

  // Content fade-in
  gsap.fromTo('.notfound-content',
    { opacity: 0, y: 30, scale: 0.98 },
    { opacity: 1, y: 0, scale: 1, duration: 0.8, ease: 'power4.out', delay: 0.3 }
  );
});
</script>

<template>
  <div class="relative min-h-screen flex flex-col justify-center items-center px-4 overflow-hidden bg-background">
    <!-- Beautiful Ambient Background Orbs -->
    <div class="ambient-orb-1 absolute top-[-10%] left-[-10%] w-[60%] h-[60%] rounded-full bg-emerald-500/10 dark:bg-emerald-500/5 blur-[120px] pointer-events-none"></div>
    <div class="ambient-orb-2 absolute bottom-[-10%] right-[-10%] w-[60%] h-[60%] rounded-full bg-emerald-600/10 dark:bg-emerald-600/5 blur-[120px] pointer-events-none"></div>

    <!-- Theme Toggle -->
    <div class="absolute top-6 right-6">
      <ThemeButton variant="rounded" />
    </div>

    <!-- Main Content Container -->
    <div class="notfound-content w-full max-w-md flex flex-col items-center text-center space-y-8 z-10">
      <!-- 404 Icon & Numbers -->
      <div class="relative flex items-center justify-center">
        <div class="absolute w-24 h-24 rounded-full bg-emerald-500/10 dark:bg-emerald-500/5 blur-md animate-pulse"></div>
        <FileQuestion class="w-16 h-16 text-primary relative z-10" />
      </div>

      <div class="space-y-3">
        <h1 class="text-6xl font-black tracking-tighter text-primary">404</h1>
        <h2 class="text-xl font-bold text-foreground tracking-tight">Page Not Found</h2>
        <p class="text-xs text-muted-foreground max-w-sm leading-relaxed mx-auto">
          The page you are looking for might have been removed, had its name changed, or is temporarily unavailable.
        </p>
      </div>

      <!-- Action Buttons -->
      <div class="flex flex-col sm:flex-row gap-3 w-full justify-center">
        <!-- Go Back -->
        <Button @click="goBack" variant="outline" class="h-10 rounded-lg font-semibold flex items-center gap-2 border border-border/60 hover:bg-muted/40 cursor-pointer">
          <ArrowLeft class="w-4 h-4" />
          <span>Go Back</span>
        </Button>

        <!-- Go Home -->
        <Button @click="goHome" class="h-10 rounded-lg font-bold shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-[1.01] bg-primary hover:bg-primary/90 text-primary-foreground flex items-center gap-2 cursor-pointer">
          <Home class="w-4 h-4" />
          <span>{{ isAuthenticated ? 'Dashboard' : 'Sign In' }}</span>
        </Button>
      </div>
    </div>

    <!-- Brand Footer -->
    <div class="absolute bottom-6 text-center">
      <p class="text-[10px] text-muted-foreground select-none">
        Powered by {{ siteConfig.name }}. All rights reserved.
      </p>
    </div>
  </div>
</template>
