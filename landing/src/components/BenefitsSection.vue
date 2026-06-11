<script setup lang="ts">
import { useScrollReveal } from '@/lib/useScrollReveal'
import { ArrowUp, Clock, Server, TrendingUp } from '@lucide/vue'

const { observe } = useScrollReveal()

const benefits = [
  {
    icon: ArrowUp,
    title: 'Reduce Downtime',
    desc: 'Get alerted within seconds of a failure. Mean time to detection drops from hours to instant.',
    stat: '99.9%',
    statLabel: 'uptime achieved',
  },
  {
    icon: Clock,
    title: 'Save Engineering Time',
    desc: 'Stop manual health checks. Automate monitoring across all your services from one dashboard.',
    stat: '10x',
    statLabel: 'faster incident response',
  },
  {
    icon: Server,
    title: 'Own Your Infrastructure',
    desc: 'Self-hosted means no vendor lock-in, no per-check pricing, no data leaving your network.',
    stat: '100%',
    statLabel: 'data sovereignty',
  },
  {
    icon: TrendingUp,
    title: 'Build User Trust',
    desc: 'Publish a live status page and show your users you take reliability seriously. Transparency wins.',
    stat: '24/7',
    statLabel: 'continuous monitoring',
  },
]
</script>

<template>
  <section class="py-20 md:py-28">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <!-- Section header -->
      <div class="mx-auto max-w-2xl text-center mb-16 reveal" ref="reveal">
        <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">Why self-host?</h2>
        <p class="mt-3 text-lg text-muted-foreground">The numbers speak for themselves</p>
      </div>

      <!-- Grid -->
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <div
          v-for="(benefit, i) in benefits"
          :key="benefit.title"
          :ref="(el: any) => observe(el)"
          class="reveal group relative"
          :class="`reveal-delay-${i + 1}`"
        >
          <div class="relative overflow-hidden rounded-xl border border-border/50 bg-card/60 dark:bg-card/40 p-6 hover:shadow-md transition-shadow">
            <!-- Stat -->
            <div class="mb-4">
              <span class="font-display text-3xl font-bold text-primary">{{ benefit.stat }}</span>
              <p class="text-xs text-muted-foreground mt-0.5">{{ benefit.statLabel }}</p>
            </div>

            <!-- Icon -->
            <div
              class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary mb-3 group-hover:bg-primary/20 transition-colors"
            >
              <component :is="benefit.icon" class="h-5 w-5" />
            </div>

            <h3 class="font-semibold text-sm mb-1">{{ benefit.title }}</h3>
            <p class="text-sm text-muted-foreground leading-relaxed">{{ benefit.desc }}</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
