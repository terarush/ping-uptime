<script setup lang="ts">
import { Separator } from '@/components/ui/separator'
import { useRouter } from 'vue-router'

const router = useRouter()

function navigateTo(url: string) {
  if (url.startsWith('/')) {
    router.push(url)
  }
}

const footerLinks = {
  Navigate: [
    { name: 'Benefits', url: '#benefits' },
    { name: 'Features', url: '#features' },
    { name: 'Services', url: '#services' },
  ],
  Help: [
    { name: 'Contact Us', url: '#contact' },
    { name: 'FAQ', url: '#' },
    { name: 'Feedback', url: '#' },
  ],
  Resources: [
    { name: 'How to Contribute', url: '/contributing' },
    { name: 'Repository', url: 'https://github.com/terarush/ping-uptime' },
    { name: 'Community', url: 'https://discord.gg/G32T9q55xh' },
    { name: 'Company', url: 'https://terarush.dev/' },
  ],
}
</script>

<template>
  <footer class="border-t border-border/50 bg-background">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-12">
      <div class="grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
        <div>
          <div class="flex items-center gap-2 mb-4">
            <div class="relative flex h-3 w-3">
              <span
                class="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary opacity-75"
              ></span>
              <span class="relative inline-flex h-3 w-3 rounded-full bg-primary"></span>
            </div>
            <span class="text-lg font-bold tracking-tight">ping-uptime</span>
          </div>
          <p class="text-xs text-muted-foreground leading-relaxed">
            Self-hosted uptime monitoring built with Go and Vue. Monitor websites, APIs, and
            services with real-time alerts and status pages.
          </p>
        </div>

        <div v-for="(links, title) in footerLinks" :key="title">
          <h4 class="text-xs font-semibold text-foreground mb-4 uppercase tracking-wider">
            {{ title }}
          </h4>
          <ul class="space-y-2">
            <li v-for="link in links" :key="link.name">
              <a
                :href="link.url"
                :target="link.url.startsWith('#') || link.url.startsWith('/') ? undefined : '_blank'"
                :rel="link.url.startsWith('#') || link.url.startsWith('/') ? undefined : 'noopener noreferrer'"
                @click.prevent="link.url.startsWith('/') ? navigateTo(link.url) : undefined"
                class="text-xs text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                >{{ link.name }}</a
              >
            </li>
          </ul>
        </div>
      </div>

      <Separator class="my-8" />

      <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
        <span class="text-xs text-muted-foreground mr-2"
          >&copy; 2026 Terarush. All rights reserved.</span
        >
        <div class="flex items-center gap-4">
          <span class="text-xs text-muted-foreground mr-2">Follow us:</span>
          <a
            href="https://twitter.com/terarush"
            target="_blank"
            rel="noopener noreferrer"
            class="p-2 rounded-full bg-card border border-border shadow-sm hover:shadow-md hover:border-primary/50 transition-shadow text-muted-foreground hover:text-foreground"
            aria-label="Twitter"
          >
            <img src="/twitter.svg" class="h-4 w-4 opacity-70 dark:invert" aria-hidden="true" />
          </a>
          <a
            href="https://linkedin.com/company/terarush"
            target="_blank"
            rel="noopener noreferrer"
            class="p-2 rounded-full bg-card border border-border shadow-sm hover:shadow-md hover:border-primary/50 transition-shadow text-muted-foreground hover:text-blue-500"
            aria-label="LinkedIn"
          >
            <img src="/linkedin.svg" class="h-4 w-4 opacity-70 dark:invert" aria-hidden="true" />
          </a>
          <a
            href="https://github.com/terarush"
            target="_blank"
            rel="noopener noreferrer"
            class="p-2 rounded-full bg-card border border-border shadow-sm hover:shadow-md hover:border-primary/50 transition-shadow text-muted-foreground hover:text-foreground"
            aria-label="GitHub"
          >
            <img src="/github.svg" class="h-4 w-4 opacity-70 dark:invert" aria-hidden="true" />
          </a>
        </div>
      </div>
    </div>
  </footer>
</template>
