<script setup lang="ts">
import { Separator } from '@/components/ui/separator'

const footerLinks = {
  Platforms: ['iOS', 'Android', 'Web'],
  Help: ['Contact Us', 'FAQ', 'Feedback'],
  Resources: ['Github', 'Community', 'Company'],
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
            <li v-for="link in links" :key="link">
              <a
                href="#"
                class="text-xs text-muted-foreground hover:text-foreground transition-colors"
                >{{ link }}</a
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
            ><svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="lucide lucide-twitter h-4 w-4"
              aria-hidden="true"
            >
              <path
                d="M22 4s-.7 2.1-2 3.4c1.6 10-9.4 17.3-18 11.6 2.2.1 4.4-.6 6-2C3 15.5.5 9.6 3 5c2.2 2.6 5.6 4.1 9 4-.9-4.2 4-6.6 7-3.8 1.1 0 3-1.2 3-1.2z"
              ></path></svg
          ></a>
          <a
            href="https://linkedin.com/company/terarush"
            target="_blank"
            rel="noopener noreferrer"
            class="p-2 rounded-full bg-card border border-border shadow-sm hover:shadow-md hover:border-primary/50 transition-shadow text-muted-foreground hover:text-blue-500"
            aria-label="LinkedIn"
            ><svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="lucide lucide-linkedin h-4 w-4"
              aria-hidden="true"
            >
              <path
                d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"
              ></path>
              <rect width="4" height="12" x="2" y="9"></rect>
              <circle cx="4" cy="4" r="2"></circle></svg
          ></a>
          <a
            href="https://github.com/terarush"
            target="_blank"
            rel="noopener noreferrer"
            class="p-2 rounded-full bg-card border border-border shadow-sm hover:shadow-md hover:border-primary/50 transition-shadow text-muted-foreground hover:text-foreground"
            aria-label="GitHub"
            ><svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="lucide lucide-github h-4 w-4"
              aria-hidden="true"
            >
              <path
                d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"
              ></path>
              <path d="M9 18c-4.51 2-5-2-7-2"></path></svg
          ></a>
        </div>
      </div>
    </div>
  </footer>
</template>
