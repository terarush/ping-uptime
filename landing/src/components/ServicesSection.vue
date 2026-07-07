<script setup lang="ts">
import { Badge } from '@/components/ui/badge'
import { useScrollReveal } from '@/lib/useScrollReveal'
import { Monitor, Activity, Webhook, GlobeLock, BarChart, BellRing } from '@lucide/vue'

const { observe } = useScrollReveal()

interface Service {
  icon: any
  title: string
  desc: string
  badge: string
  group: string
}

const services: Service[] = [
  { icon: Monitor, title: 'HTTP/HTTPS Monitoring', desc: 'Check any web endpoint. Follow redirects, validate status codes and SSL certificates.', badge: 'Core', group: 'Core' },
  { icon: Activity, title: 'Ping and TCP Checks', desc: 'ICMP ping and TCP port checks for services behind firewalls or non-HTTP infrastructure.', badge: 'Core', group: 'Core' },
  { icon: Webhook, title: 'Webhook Integrations', desc: 'Send alerts to Slack, Discord, Telegram, PagerDuty, or any custom webhook endpoint.', badge: 'Integrations', group: 'Platform' },
  { icon: GlobeLock, title: 'SSL Certificate Monitoring', desc: 'Get notified before certificates expire. Avoid the "this site is not secure" nightmare.', badge: 'Security', group: 'Platform' },
  { icon: BarChart, title: 'Response Time Analytics', desc: 'Track latency trends over time. Spot degradation before it becomes downtime.', badge: 'Analytics', group: 'Platform' },
  { icon: BellRing, title: 'Status Pages', desc: 'Publish a public status page with live uptime data, incident history, and subscriber notifications.', badge: 'Pages', group: 'Platform' },
]
</script>

<template>
  <section id="services" class="section-grid py-20 md:py-28">
    <div class="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-2xl text-center mb-16 reveal" :ref="(el: any) => observe(el)">
        <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">What you can monitor</h2>
        <p class="mt-3 text-lg text-muted-foreground">Every protocol, every service, one dashboard</p>
      </div>

      <div class="mx-auto max-w-5xl">
        <!-- Core (compact 2-col) -->
        <div class="grid gap-4 md:grid-cols-2 mb-4">
          <div
            v-for="(service, i) in services.filter(s => s.group === 'Core')"
            :key="service.title"
            :ref="(el: any) => observe(el)"
            class="reveal group"
            :class="`reveal-delay-${i + 1}`"
          >
            <div class="flex items-start gap-4 rounded-xl border border-border/50 bg-card/60 dark:bg-card/40 p-5 hover:shadow-md hover:border-primary/20 transition-all duration-300">
              <div
                class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary group-hover:bg-primary/[0.15] transition-colors"
              >
                <component :is="service.icon" class="h-6 w-6" />
              </div>
              <div class="min-w-0 flex-1">
                <div class="flex items-center gap-2 mb-1">
                  <h3 class="font-semibold text-sm">{{ service.title }}</h3>
                  <Badge variant="secondary" class="text-[10px] px-2 py-0.5 shrink-0">
                    {{ service.badge }}
                  </Badge>
                </div>
                <p class="text-sm text-muted-foreground leading-relaxed">{{ service.desc }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Platform (compact 3-col) -->
        <div class="grid gap-4 md:grid-cols-3">
          <div
            v-for="(service, i) in services.filter(s => s.group === 'Platform')"
            :key="service.title"
            :ref="(el: any) => observe(el)"
            class="reveal group"
            :class="`reveal-delay-${i + 2}`"
          >
            <div class="flex flex-col gap-3 rounded-xl border border-border/50 bg-card/60 dark:bg-card/40 p-5 hover:shadow-md hover:border-primary/20 transition-all duration-300 h-full">
              <div
                class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:bg-primary/[0.15] transition-colors shrink-0"
              >
                <component :is="service.icon" class="h-5 w-5" />
              </div>
              <div>
                <div class="flex items-center gap-2 mb-1">
                  <h3 class="font-semibold text-sm">{{ service.title }}</h3>
                  <Badge variant="secondary" class="text-[10px] px-2 py-0.5 shrink-0">
                    {{ service.badge }}
                  </Badge>
                </div>
                <p class="text-sm text-muted-foreground leading-relaxed">{{ service.desc }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
