<script setup lang="ts">
import { useScrollReveal } from '@/lib/useScrollReveal'
import { Bell, BarChart3, Globe, Users, Shield, Zap } from '@lucide/vue'

const { observe } = useScrollReveal()

const features = [
  {
    icon: Bell,
    title: 'Real-time alerts',
    desc: 'Get notified via email, Slack, Discord, or Telegram the moment a service goes down.',
    tag: 'Notifications',
    featured: true,
  },
  {
    icon: Shield,
    title: 'Self-hosted privacy',
    desc: 'Your monitoring data stays on your infrastructure. No third-party cloud, no data leaks.',
    tag: 'Security',
    featured: true,
  },
  {
    icon: BarChart3,
    title: 'Detailed analytics',
    desc: 'Track response times, uptime percentages, and historical trends with interactive charts.',
    tag: 'Data',
    featured: false,
  },
  {
    icon: Globe,
    title: 'Multi-region checks',
    desc: 'Monitor from multiple locations worldwide to ensure global accessibility.',
    tag: 'Infrastructure',
    featured: false,
  },
  {
    icon: Users,
    title: 'Team collaboration',
    desc: 'Invite team members, assign roles, and share status pages with stakeholders.',
    tag: 'Team',
    featured: false,
  },
  {
    icon: Zap,
    title: 'Lightning fast setup',
    desc: 'Deploy in minutes with a single binary. Go plus a database, and you are online.',
    tag: 'Deploy',
    featured: false,
  },
]
</script>

<template>
  <section id="features" class="py-20 md:py-28 bg-muted/30">
    <div class="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-2xl text-center mb-16 reveal" :ref="(el: any) => observe(el)">
        <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">Everything you need</h2>
        <p class="mt-3 text-lg text-muted-foreground">Built for teams that take uptime seriously</p>
      </div>

      <!-- Bento grid: 2 featured (wider) + 4 compact -->
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <div
          v-for="(feature, i) in features"
          :key="feature.title"
          :ref="(el: any) => observe(el)"
          class="reveal group"
          :class="[
            `reveal-delay-${Math.min(i + 1, 5)}`,
            i < 2 ? 'lg:col-span-2' : '',
          ]"
        >
          <div :class="[
            'relative h-full rounded-xl border p-6 hover:shadow-md transition-all duration-300',
            feature.featured
              ? 'border-primary/20 bg-primary/[0.03] dark:bg-primary/[0.06]'
              : 'border-border/50 bg-card/60 dark:bg-card/40 hover:-translate-y-0.5'
          ]">
            <div
              class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary mb-4 group-hover:bg-primary/20 transition-colors"
            >
              <component :is="feature.icon" class="h-5 w-5" />
            </div>
            <div class="text-[10px] font-mono font-medium text-primary/60 uppercase tracking-wider mb-1">
              {{ feature.tag }}
            </div>
            <h3 class="font-semibold text-sm mb-1">{{ feature.title }}</h3>
            <p class="text-sm text-muted-foreground leading-relaxed">{{ feature.desc }}</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
