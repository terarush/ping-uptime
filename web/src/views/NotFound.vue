<script setup lang="ts">
import Button from 'primevue/button';
import { FileQuestion, Home, ArrowLeft } from '@lucide/vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const goHome = () => {
  router.push('/');
};

const goBack = () => {
  router.back();
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50 dark:bg-slate-950 p-6 transition-colors duration-300">
    <div class="max-w-md w-full text-center space-y-8 bg-white dark:bg-slate-900 p-8 rounded-2xl shadow-xl border border-slate-100 dark:border-slate-800 transition-all duration-300">
      <!-- Icon & Number -->
      <div class="relative flex justify-center">
        <div class="absolute inset-0 bg-primary/10 rounded-full blur-2xl w-24 h-24 mx-auto"></div>
        <div class="relative w-20 h-20 rounded-2xl bg-primary/10 dark:bg-primary/5 flex items-center justify-center border border-primary/20 dark:border-primary/10">
          <FileQuestion class="w-10 h-10 text-primary" />
        </div>
      </div>

      <div class="space-y-3">
        <h1 class="text-7xl font-black text-slate-800 dark:text-slate-100 tracking-tight">404</h1>
        <h2 class="text-xl font-bold text-slate-800 dark:text-slate-200">Page Not Found</h2>
        <p class="text-slate-500 dark:text-slate-400 text-sm max-w-xs mx-auto">
          The page you are looking for might have been removed, had its name changed, or is temporarily unavailable.
        </p>
      </div>

      <!-- Action Buttons -->
      <div class="flex flex-col sm:flex-row gap-3 justify-center pt-2">
        <Button
          @click="goBack"
          severity="secondary"
          outlined
          class="flex items-center justify-center gap-2 px-4 py-2.5 text-sm"
        >
          <ArrowLeft class="w-4 h-4" />
          <span>Go Back</span>
        </Button>

        <Button
          @click="goHome"
          class="flex items-center justify-center gap-2 px-4 py-2.5 text-sm"
        >
          <Home class="w-4 h-4" />
          <span>Go Home</span>
        </Button>
      </div>
    </div>
  </div>
</template>
