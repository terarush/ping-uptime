<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Card, CardContent } from '@/components/ui/card'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { useScrollReveal } from '@/lib/useScrollReveal'

const { observe } = useScrollReveal()

interface Contributor {
  login: string
  avatar_url: string
  html_url: string
  contributions: number
  type: string
}

const apiContributors = ref<Contributor[]>([])
const error = ref('')

const contributors = computed(() => [...apiContributors.value])

onMounted(async () => {
  try {
    const res = await fetch('https://api.github.com/repos/terarush/ping-uptime/contributors?per_page=30')
    if (!res.ok) throw new Error(`GitHub API: ${res.status}`)
    apiContributors.value = await res.json()
  } catch (e) {
    error.value = 'Failed to load contributors'
  }
})
</script>

<template>
  <section id="contributors" class="section-grid py-20 md:py-28 bg-muted/30">
    <div class="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <!-- Section header -->
      <div class="mx-auto max-w-2xl text-center mb-16 reveal" :ref="(el: any) => observe(el)">
        <span class="inline-flex items-center gap-2 rounded-full border border-border bg-background/80 px-3.5 py-1 text-xs font-mono font-medium text-muted-foreground backdrop-blur-sm mb-5 tracking-wider">
          <span class="h-1.5 w-1.5 rounded-full bg-primary/60" />
          05 / Contributors
        </span>
        <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">Contributors</h2>
        <p class="mt-3 text-lg text-muted-foreground">People who helped build ping-uptime</p>
      </div>

      <!-- Error state -->
      <p v-if="error" class="text-center text-sm text-muted-foreground">{{ error }}</p>

      <!-- Empty / loading state -->
      <p v-else-if="!contributors.length" class="text-center text-sm text-muted-foreground">Loading contributors...</p>

      <!-- Grid -->
      <div v-else class="grid gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
        <Card
          v-for="(c, i) in contributors"
          :key="c.login"
          :ref="(el: any) => observe(el?.$el ?? el)"
          class="reveal group border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md text-center hover:shadow-md transition-shadow"
          :class="`reveal-delay-${i + 1}`"
        >
          <template #default>
            <a :href="c.html_url" target="_blank" rel="noopener noreferrer" class="block p-6">
              <Avatar class="mx-auto h-20 w-20 mb-4 ring-2 ring-border/50 group-hover:ring-primary/30 transition-all">
                <AvatarImage :src="c.avatar_url" :alt="c.login" />
                <AvatarFallback>{{ c.login.slice(0, 2).toUpperCase() }}</AvatarFallback>
              </Avatar>
              <h3 class="font-semibold text-sm">{{ c.login }}</h3>
              <p class="text-xs text-muted-foreground mt-1 leading-relaxed">
                {{ c.contributions }} commit{{ c.contributions !== 1 ? 's' : '' }}
              </p>
              <span
                class="mt-3 inline-flex items-center gap-1 text-xs font-medium text-primary/80 hover:text-primary transition-colors"
              >
                View profile &rarr;
              </span>
            </a>
          </template>
        </Card>
      </div>
    </div>
  </section>
</template>
