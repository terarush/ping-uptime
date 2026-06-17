<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Menu, X, Moon, Sun } from '@lucide/vue'
import { Button } from '@/components/ui/button'

const router = useRouter()

const isScrolled = ref(false)
const isMobileMenuOpen = ref(false)
const isDark = ref(false)
const activeSection = ref('')
const isInHero = ref(true)

function handleScroll() {
  isScrolled.value = window.scrollY > 50

  const sections = ['hero', 'benefits', 'features', 'services', 'contributors', 'contact']
  for (const id of sections) {
    const el = document.getElementById(id)
    if (el) {
      const rect = el.getBoundingClientRect()
      if (rect.top <= 150 && rect.bottom >= 150) {
        activeSection.value = id
        break
      }
    }
  }

  isInHero.value = window.scrollY < window.innerHeight * 0.6
}

function toggleDark() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
}

function navigateTo(href: string) {
  isMobileMenuOpen.value = false
  if (href.startsWith('#')) {
    const id = href.slice(1)
    const el = document.getElementById(id)
    if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  } else {
    router.push(href)
  }
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll, { passive: true })
  isDark.value = document.documentElement.classList.contains('dark')
  handleScroll()
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

const navLinks = [
  { label: 'Home', href: '#hero' },
  { label: 'Benefits', href: '#benefits' },
  { label: 'Features', href: '#features' },
  { label: 'Services', href: '#services' },
  { label: 'Contributors', href: '#contributors' },
  { label: 'Contact', href: '#contact' },
  { label: 'Contributing', href: '/contributing', external: true },
]
</script>

<template>
  <nav
    :class="[
      'fixed top-0 z-50 w-full transition-all duration-500',
      isScrolled
        ? 'bg-background/70 backdrop-blur-2xl border-b border-border/40 shadow-sm'
        : 'bg-transparent',
    ]"
  >
    <!-- Top gradient bar (subtle) -->
    <div
      :class="[
        'h-0.5 bg-linear-to-r from-transparent via-primary/30 to-transparent transition-opacity duration-700',
        isScrolled ? 'opacity-100' : 'opacity-0',
      ]"
    />

    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16 md:h-18">
        <!-- Brand -->
        <a
          href="#hero"
          @click.prevent="navigateTo('#hero')"
          class="flex items-center gap-2.5 shrink-0 group"
        >
          <div class="relative flex h-3 w-3">
            <span
              class="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary opacity-75"
            />
            <span class="relative inline-flex h-3 w-3 rounded-full bg-primary" />
          </div>
          <span class="font-display text-lg font-bold tracking-tight text-foreground"
            >ping-uptime</span
          >
        </a>

        <!-- Desktop nav -->
        <div class="hidden md:flex items-center gap-1">
          <div class="flex items-center gap-1 rounded-xl bg-muted/50 p-1">
            <a
              v-for="link in navLinks"
              :key="link.label"
              :href="link.href"
              @click.prevent="navigateTo(link.href)"
              class="relative px-4 py-2 text-sm font-medium rounded-lg transition-all duration-300"
              :class="[
                !link.external && activeSection === link.href.slice(1)
                  ? 'text-foreground bg-background shadow-xs'
                  : 'text-muted-foreground hover:text-foreground hover:bg-background/50',
              ]"
            >
              {{ link.label }}
            </a>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-0.5">
          <!-- Dark mode toggle -->
          <Button
            variant="ghost"
            size="icon"
            @click="toggleDark"
            class="rounded-full text-muted-foreground hover:text-foreground"
            aria-label="Toggle theme"
          >
            <Sun v-if="!isDark" class="h-4 w-4 transition-transform duration-300 rotate-0" />
            <Moon v-else class="h-4 w-4 transition-transform duration-300 rotate-0" />
          </Button>

          <!-- GitHub link -->
          <a
            href="https://github.com/terarush/ping-uptime"
            target="_blank"
            rel="noopener noreferrer"
            class="hidden sm:inline-flex items-center gap-2 px-3 py-2 text-xs font-medium text-muted-foreground hover:text-foreground transition-colors rounded-lg hover:bg-accent"
          >
            <img src="/github.svg" class="h-4 w-4 opacity-70 dark:invert" aria-hidden="true" />
            <span>Star on GitHub</span>
          </a>

          <!-- Mobile hamburger -->
          <Button
            variant="ghost"
            size="icon"
            class="md:hidden rounded-full text-muted-foreground hover:text-foreground"
            @click="isMobileMenuOpen = !isMobileMenuOpen"
            :aria-label="isMobileMenuOpen ? 'Close menu' : 'Open menu'"
          >
            <Menu v-if="!isMobileMenuOpen" class="h-4 w-4" />
            <X v-else class="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  </nav>

  <!-- Mobile drawer -->
  <!-- Backdrop -->
  <Transition
    enter-active-class="transition-opacity duration-300 ease-out"
    leave-active-class="transition-opacity duration-200 ease-in"
    enter-from-class="opacity-0"
    leave-to-class="opacity-0"
  >
    <div
      v-if="isMobileMenuOpen"
      class="fixed inset-0 bg-black/20 backdrop-blur-sm z-40 md:hidden"
      @click="isMobileMenuOpen = false"
    />
  </Transition>

  <!-- Panel -->
  <Transition
    enter-active-class="transition-all duration-300 ease-out"
    leave-active-class="transition-all duration-200 ease-in"
    enter-from-class="opacity-0 scale-95 -translate-y-4"
    leave-to-class="opacity-0 scale-95 -translate-y-4"
  >
    <div v-if="isMobileMenuOpen" class="fixed top-16 left-0 right-0 mx-4 mt-2 z-50 md:hidden">
      <div
        class="rounded-xl border border-border/50 bg-background/95 backdrop-blur-2xl shadow-xl overflow-hidden"
      >
        <div class="px-2 py-3 space-y-0.5">
          <a
            v-for="link in navLinks"
            :key="link.label"
            :href="link.href"
            @click.prevent="navigateTo(link.href)"
            class="flex items-center gap-3 px-4 py-3 text-sm font-medium rounded-lg transition-colors"
            :class="
              !link.external && activeSection === link.href.slice(1)
                ? 'text-foreground bg-accent'
                : 'text-muted-foreground hover:text-foreground hover:bg-accent/50'
            "
          >
            <span
              class="h-1.5 w-1.5 rounded-full"
              :class="
                activeSection && activeSection === link.href.slice(1) ? 'bg-primary' : 'bg-muted-foreground/30'
              "
            />
            {{ link.label }}
          </a>
        </div>

        <div class="border-t border-border/50 px-2 py-3">
          <a
            href="https://github.com/terarush/ping-uptime"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center gap-3 px-4 py-3 text-sm font-medium text-muted-foreground hover:text-foreground rounded-lg hover:bg-accent/50 transition-colors"
            @click="isMobileMenuOpen = false"
          >
            <img src="/github.svg" class="h-4 w-4 opacity-70" aria-hidden="true" />
            <span>Star on GitHub</span>
          </a>
        </div>
      </div>
    </div>
  </Transition>
</template>
