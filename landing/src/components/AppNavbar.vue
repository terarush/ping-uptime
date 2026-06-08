<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Menu, X, Moon, Sun } from '@lucide/vue'
import { Button } from '@/components/ui/button'

const isScrolled = ref(false)
const isMobileMenuOpen = ref(false)
const isDark = ref(false)

function handleScroll() {
  isScrolled.value = window.scrollY > 50
}

function toggleDark() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
  isDark.value = document.documentElement.classList.contains('dark')
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

const navLinks = [
  { label: 'Home', href: '#hero' },
  { label: 'Features', href: '#features' },
  { label: 'Team', href: '#team' },
  { label: 'Contact', href: '#contact' },
  // { label: 'FAQ', href: '#faq' },
]
</script>

<template>
  <nav
    :class="[
      'fixed top-0 z-50 w-full transition-all duration-300',
      isScrolled
        ? 'bg-background/80 backdrop-blur-xl border-b border-border/50 shadow-sm'
        : 'bg-transparent',
    ]"
  >
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="grid grid-cols-3 items-center h-16">
        <div class="flex justify-start items-center gap-2">
          <div class="relative flex h-3 w-3">
            <span
              class="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary opacity-75"
            ></span>
            <span class="relative inline-flex h-3 w-3 rounded-full bg-primary"></span>
          </div>
          <span class="text-lg font-bold tracking-tight">ping-uptime</span>
        </div>

        <div class="hidden md:flex justify-center items-center gap-6">
          <a
            v-for="link in navLinks"
            :key="link.label"
            :href="link.href"
            class="text-sm text-muted-foreground hover:text-foreground transition-colors font-medium"
          >
            {{ link.label }}
          </a>
        </div>

        <div class="flex justify-end items-center gap-2">
          <Button variant="ghost" size="icon" @click="toggleDark" class="rounded-full">
            <Sun v-if="!isDark" class="h-4 w-4" />
            <Moon v-else class="h-4 w-4" />
          </Button>

          <Button
            variant="ghost"
            size="icon"
            class="md:hidden rounded-full"
            @click="isMobileMenuOpen = !isMobileMenuOpen"
          >
            <Menu v-if="!isMobileMenuOpen" class="h-4 w-4" />
            <X v-else class="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>

    <div
      v-if="isMobileMenuOpen"
      class="md:hidden border-t border-border/50 bg-background/95 backdrop-blur-xl"
    >
      <div class="space-y-1 px-4 py-4">
        <a
          v-for="link in navLinks"
          :key="link.label"
          :href="link.href"
          class="block px-3 py-2 text-sm text-muted-foreground hover:text-foreground hover:bg-accent rounded-md transition-colors font-medium"
          @click="isMobileMenuOpen = false"
        >
          {{ link.label }}
        </a>
      </div>
    </div>
  </nav>
</template>
