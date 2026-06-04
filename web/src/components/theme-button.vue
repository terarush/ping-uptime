<script setup lang="ts">
import { ref, onMounted } from 'vue';
import Button from 'primevue/button';
import { Sun, Moon } from '@lucide/vue';

const isDark = ref(false);

const updateThemeClasses = (dark: boolean) => {
  const el = document.documentElement;
  if (dark) {
    el.classList.add('my-app-dark', 'p-dark');
  } else {
    el.classList.remove('my-app-dark', 'p-dark');
  }
};

const toggleDarkMode = () => {
  isDark.value = !isDark.value;
  updateThemeClasses(isDark.value);
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light');
};

onMounted(() => {
  const userTheme = localStorage.getItem('theme');
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

  isDark.value = userTheme === 'dark' || (!userTheme && prefersDark);
  updateThemeClasses(isDark.value);
});
</script>

<template>
  <Button
    @click="toggleDarkMode"
    severity="secondary"
    variant="text"
    rounded
    aria-label="Toggle Dark Mode"
  >
    <template #icon>
      <Sun v-if="isDark" class="w-5 h-5 text-yellow-500" />
      <Moon v-else class="w-5 h-5 text-slate-700" />
    </template>
  </Button>
</template>
