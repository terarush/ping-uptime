<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { useScrollReveal } from '@/lib/useScrollReveal'
import { Bell, BarChart3, Globe, Users, Shield, Zap } from '@lucide/vue'

const { observe } = useScrollReveal()

const features = [
  {
    icon: Bell,
    title: 'Real-Time Alerts',
    desc: 'Get notified via email, Slack, Discord, or Telegram the moment a service goes down. No delays, no excuses.',
  },
  {
    icon: BarChart3,
    title: 'Detailed Analytics',
    desc: 'Track response times, uptime percentages, and historical trends with interactive charts and exportable reports.',
  },
  {
    icon: Globe,
    title: 'Multi-Region Checks',
    desc: 'Monitor from multiple locations worldwide to ensure your services are accessible everywhere, not just locally.',
  },
  {
    icon: Users,
    title: 'Team Collaboration',
    desc: 'Invite team members, assign roles, and share status pages with stakeholders. Uptime is a team sport.',
  },
  {
    icon: Shield,
    title: 'Self-Hosted Privacy',
    desc: 'Your monitoring data stays on your infrastructure. No third-party cloud, no data leaks, full control.',
  },
  {
    icon: Zap,
    title: 'Lightning Fast Setup',
    desc: 'Deploy in minutes with a single binary. No complex dependencies — just Go, a database, and you\'re online.',
  },
]
</script>

<template>
  <section id="features" class="section-grid py-20 md:py-28 bg-muted/30">
    <div class="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <!-- Section header -->
      <div class="mx-auto max-w-2xl text-center mb-16 reveal" :ref="(el: any) => observe(el)">
        <span class="inline-flex items-center gap-2 rounded-full border border-border bg-background/80 px-3.5 py-1 text-xs font-mono font-medium text-muted-foreground backdrop-blur-sm mb-5 tracking-wider">
          <span class="h-1.5 w-1.5 rounded-full bg-primary/60" />
          02 / Features
        </span>
        <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">Everything you need</h2>
        <p class="mt-3 text-lg text-muted-foreground">Built for teams that take uptime seriously</p>
        <p class="mt-2 text-sm text-muted-foreground/80 max-w-md mx-auto">
          From solo devs to enterprise teams — every feature ships out of the box, zero configuration required.
        </p>
      </div>

      <!-- Grid -->
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <Card
          v-for="(feature, i) in features"
          :key="feature.title"
          :ref="(el: any) => observe(el?.$el ?? el)"
          class="reveal group hover:-translate-y-0.5 transition-all duration-300"
          :class="`reveal-delay-${(i % 3) + 1}`"
        >
          <template #default>
            <CardHeader>
              <div
                class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary mb-2 group-hover:bg-primary/20 transition-colors"
              >
                <component :is="feature.icon" class="h-5 w-5" />
              </div>
              <CardTitle class="text-base font-semibold">{{ feature.title }}</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription class="text-sm leading-relaxed">{{ feature.desc }}</CardDescription>
            </CardContent>
          </template>
        </Card>
      </div>
    </div>
  </section>
</template>
