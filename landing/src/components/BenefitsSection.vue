<script setup lang="ts">
import { useScrollReveal } from '@/lib/useScrollReveal'
import { ArrowUp, Clock, Server, TrendingUp } from '@lucide/vue'

const { observe } = useScrollReveal()

const benefits = [
  {
    icon: ArrowUp,
    title: 'Reduce downtime',
    desc: 'Get alerted within seconds of a failure. Mean time to detection drops from hours to instant.',
    featured: true,
  },
  {
    icon: Clock,
    title: 'Save engineering time',
    desc: 'Stop manual health checks. Automate monitoring across all your services from one dashboard.',
    featured: false,
  },
  {
    icon: Server,
    title: 'Own your infrastructure',
    desc: 'Self-hosted means no vendor lock-in, no per-check pricing, no data leaving your network.',
    featured: false,
  },
  {
    icon: TrendingUp,
    title: 'Build user trust',
    desc: 'Publish a live status page and show your users you take reliability seriously.',
    featured: false,
  },
]
</script>

<template>
  <section id="benefits" class="section-grid py-20 md:py-28">
    <div class="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-2xl text-center mb-16 reveal" :ref="(el: any) => observe(el)">
        <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">Why self-host?</h2>
        <p class="mt-3 text-lg text-muted-foreground">Your monitoring, your infrastructure</p>
      </div>

      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="(benefit, i) in benefits"
          :key="benefit.title"
          :ref="(el: any) => observe(el)"
          class="reveal group"
          :class="[`reveal-delay-${Math.min(i + 1, 4)}`, benefit.featured ? 'md:col-span-2 lg:col-span-2' : '']"
        >
          <div :class="[
            'relative overflow-hidden rounded-xl border border-border/50 bg-card/60 dark:bg-card/40 hover:shadow-md hover:-translate-y-0.5 transition-all duration-300',
            benefit.featured ? 'p-8' : 'p-6'
          ]">
            <div class="flex items-start gap-4">
              <div
                class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:bg-primary/20 transition-colors"
              >
                <component :is="benefit.icon" class="h-5 w-5" />
              </div>
              <div>
                <h3 class="font-semibold text-sm mb-1">{{ benefit.title }}</h3>
                <p class="text-sm text-muted-foreground leading-relaxed">{{ benefit.desc }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
