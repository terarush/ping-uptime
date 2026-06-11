<script setup lang="ts">
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { useScrollReveal } from '@/lib/useScrollReveal'
import { Monitor, Activity, Webhook, GlobeLock, BarChart, BellRing } from '@lucide/vue'

const { observe } = useScrollReveal()

const services = [
  { icon: Monitor, title: 'HTTP/HTTPS Monitoring', desc: 'Monitor any web endpoint. Follow redirects, check status codes, validate SSL certificates and expiry dates.', badge: 'Core' },
  { icon: Activity, title: 'Ping & TCP Checks', desc: 'ICMP ping and TCP port checks for services behind firewalls or non-HTTP infrastructure.', badge: 'Core' },
  { icon: Webhook, title: 'Webhook Integrations', desc: 'Send alerts to Slack, Discord, Telegram, PagerDuty, or any custom webhook endpoint.', badge: 'Integrations' },
  { icon: GlobeLock, title: 'SSL Certificate Monitoring', desc: 'Get notified before certificates expire. Avoid the "this site is not secure" nightmare.', badge: 'Security' },
  { icon: BarChart, title: 'Response Time Analytics', desc: 'Track latency trends over time. Spot degradation before it becomes downtime.', badge: 'Analytics' },
  { icon: BellRing, title: 'Status Pages', desc: 'Publish a public status page with live uptime data, incident history, and subscriber notifications.', badge: 'Pages' },
]
</script>

<template>
  <section id="services" class="section-grid py-20 md:py-28">
    <div class="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <!-- Section header -->
      <div class="mx-auto max-w-2xl text-center mb-16 reveal" :ref="(el: any) => observe(el)">
        <span class="inline-flex items-center gap-2 rounded-full border border-border bg-background/80 px-3.5 py-1 text-xs font-mono font-medium text-muted-foreground backdrop-blur-sm mb-5 tracking-wider">
          <span class="h-1.5 w-1.5 rounded-full bg-primary/60" />
          03 / Services
        </span>
        <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">What you can monitor</h2>
        <p class="mt-3 text-lg text-muted-foreground">Every protocol, every service, one dashboard</p>
        <p class="mt-2 text-sm text-muted-foreground/80 max-w-md mx-auto">
          From simple HTTP checks to complex multi-region monitoring — Ping-Uptime handles it all.
        </p>
      </div>

      <!-- Grid -->
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <Card
          v-for="(service, i) in services"
          :key="service.title"
          :ref="(el: any) => observe(el?.$el ?? el)"
          class="reveal group hover:-translate-y-0.5 transition-all duration-300"
          :class="`reveal-delay-${(i % 3) + 1}`"
        >
          <template #default>
            <CardHeader class="relative">
              <Badge variant="secondary" class="absolute top-4 right-4 text-[10px] px-2 py-0.5 font-mono">
                {{ service.badge }}
              </Badge>
              <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary mb-2 group-hover:bg-primary/20 transition-colors">
                <component :is="service.icon" class="h-5 w-5" />
              </div>
              <CardTitle class="text-base font-semibold">{{ service.title }}</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription class="text-sm leading-relaxed">{{ service.desc }}</CardDescription>
            </CardContent>
          </template>
        </Card>
      </div>
    </div>
  </section>
</template>
