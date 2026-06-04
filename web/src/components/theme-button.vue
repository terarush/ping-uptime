<script setup lang="ts">
/**
 * Theme Button Component
 * Supports toggling between Light and Dark mode.
 * Persists theme choice in localStorage and updates document classes for Tailwind CSS v4.
 */
import { ref, onMounted } from 'vue';
import { Sun, Moon } from '@lucide/vue';
import { Button } from '@/components/ui/button';
import { SidebarMenuButton } from '@/components/ui/sidebar';

interface ThemeToggleProps {
  variant?: 'default' | 'rounded';
  isSidebarItem?: boolean;
}

// Set default values for props
withDefaults(defineProps<ThemeToggleProps>(), {
  variant: 'default',
  isSidebarItem: false
});

// Reactive theme state
const theme = ref<'light' | 'dark'>('light');

/**
 * Updates document body classes for Tailind CSS dark mode
 */
const updateThemeClasses = (currentTheme: 'light' | 'dark') => {
  const el = document.documentElement;
  if (currentTheme === 'dark') {
    el.classList.add('dark');
  } else {
    el.classList.remove('dark');
  }
};

/**
 * Toggles theme state, updates HTML classes, and saves configuration to localStorage
 */
const toggleTheme = () => {
  theme.value = theme.value === 'dark' ? 'light' : 'dark';
  updateThemeClasses(theme.value);
  localStorage.setItem('theme', theme.value);
};

// Check client-side configuration settings on component mount
onMounted(() => {
  const savedTheme = localStorage.getItem('theme');
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  
  // Default to system preference if no localStorage value is set
  const initialTheme = savedTheme === 'dark' || (!savedTheme && prefersDark) ? 'dark' : 'light';
  theme.value = initialTheme;
  updateThemeClasses(initialTheme);
});
</script>

<template>
  <SidebarMenuButton
    v-if="isSidebarItem"
    @click="toggleTheme"
    tooltip="Toggle Theme"
    variant="default"
  >
    <Sun
      v-if="theme === 'dark'"
      class="w-4 h-4 text-yellow-500 shrink-0"
    />
    <Moon
      v-else
      class="w-4 h-4 text-slate-700 dark:text-slate-400 shrink-0"
    />
    <span class="group-data-[collapsible=icon]:hidden text-sm">
      {{ theme === 'dark' ? 'Light Mode' : 'Dark Mode' }}
    </span>
  </SidebarMenuButton>

  <Button
    v-else
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
