<script setup lang="ts">
import { ref, onMounted } from 'vue';
import Button from 'primevue/button';
import { Sun, Moon } from '@lucide/vue';

const isDark = ref(false);

const toggleDarkMode = () => {
  const isDarkVal = document.documentElement.classList.toggle('my-app-dark');
  isDark.value = isDarkVal;
  localStorage.setItem('theme', isDarkVal ? 'dark' : 'light');
};

onMounted(() => {
  const userTheme = localStorage.getItem('theme');
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

  if (userTheme === 'dark' || (!userTheme && prefersDark)) {
    document.documentElement.classList.add('my-app-dark');
    isDark.value = true;
  } else {
    document.documentElement.classList.remove('my-app-dark');
    isDark.value = false;
  }
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
