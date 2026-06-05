<script setup lang="ts">
import { Separator } from '@/components/ui/separator'
import { Globe, GitFork } from '@lucide/vue'

const footerLinks = {
  Platforms: ['iOS', 'Android', 'Web'],
  Help: ['Contact Us', 'FAQ', 'Feedback'],
  Socials: ['Twitch', 'Discord', 'Dribbble'],
}
</script>

<template>
  <footer class="border-t border-border/50 bg-background">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-12">
      <div class="grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
        <div>
          <div class="flex items-center gap-2 mb-4">
            <div class="relative flex h-3 w-3">
              <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary opacity-75"></span>
              <span class="relative inline-flex h-3 w-3 rounded-full bg-primary"></span>
            </div>
            <span class="text-lg font-bold tracking-tight">Shadcn</span>
          </div>
          <p class="text-xs text-muted-foreground leading-relaxed">
            Build beautiful landing pages with Shadcn + Vue.
          </p>
        </div>

        <div v-for="(links, title) in footerLinks" :key="title">
          <h4 class="text-xs font-semibold text-foreground mb-4 uppercase tracking-wider">{{ title }}</h4>
          <ul class="space-y-2">
            <li v-for="link in links" :key="link">
              <a href="#" class="text-xs text-muted-foreground hover:text-foreground transition-colors">{{ link }}</a>
            </li>
          </ul>
        </div>
      </div>

      <Separator class="my-8" />

      <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
        <p class="text-xs text-muted-foreground">
          &copy; 2024 Designed and developed by <a href="https://github.com/leoMirandaa" class="text-primary hover:underline">Leo Miranda</a>
        </p>
        <div class="flex items-center gap-4">
          <a href="#" class="text-muted-foreground hover:text-foreground transition-colors">
            <GitFork class="h-4 w-4" />
          </a>
          <a href="#" class="text-muted-foreground hover:text-foreground transition-colors">
            <Globe class="h-4 w-4" />
          </a>
          <a href="#" class="text-muted-foreground hover:text-foreground transition-colors">
            <Globe class="h-4 w-4" />
          </a>
        </div>
      </div>
    </div>
  </footer>
</template>
