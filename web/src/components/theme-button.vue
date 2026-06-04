<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Sun, Moon } from '@lucide/vue';
import { Button } from '@/components/ui/button';

interface ThemeToggleProps {
  variant?: 'default' | 'rounded';
}

withDefaults(defineProps<ThemeToggleProps>(), {
  variant: 'default'
});

const theme = ref<'light' | 'dark'>('light');

const updateThemeClasses = (currentTheme: 'light' | 'dark') => {
  const el = document.documentElement;
  if (currentTheme === 'dark') {
    el.classList.add('dark');
  } else {
    el.classList.remove('dark');
  }
};

const toggleTheme = () => {
  theme.value = theme.value === 'dark' ? 'light' : 'dark';
  updateThemeClasses(theme.value);
  localStorage.setItem('theme', theme.value);
};

onMounted(() => {
  const savedTheme = localStorage.getItem('theme');
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  
  const initialTheme = savedTheme === 'dark' || (!savedTheme && prefersDark) ? 'dark' : 'light';
  theme.value = initialTheme;
  updateThemeClasses(initialTheme);
});
</script>

<template>
  <Button
    :variant="variant === 'rounded' ? 'ghost' : 'outline'"
    :size="variant === 'rounded' ? 'icon' : 'default'"
    aria-label="Toggle theme"
    @click="toggleTheme"
    :class="variant === 'rounded' ? 'rounded-full border p-2' : ''"
  >
    <Sun
      :class="[
        'h-[1.2rem] w-[1.2rem] transition-all',
        theme === 'dark' ? 'scale-0 rotate-90 absolute' : 'scale-100 rotate-0'
      ]"
    />
    <Moon
      :class="[
        'h-[1.2rem] w-[1.2rem] transition-all',
        theme === 'dark' ? 'scale-100 rotate-0' : 'scale-0 -rotate-90 absolute'
      ]"
    />
  </Button>
</template>
