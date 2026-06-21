<script setup lang="ts">
import { Card, CardContent } from '@/components/ui/card'
import { ref } from 'vue'
import { ChevronLeft, ChevronRight } from '@lucide/vue'

const testimonials = [
  { name: 'Maya Chen', role: 'Senior DevOps Engineer', initial: 'MC', quote: 'Cut our incident response time by more than half. The multi-region checks caught a regional outage before any of our users reported it.' },
  { name: 'James Okonkwo', role: 'CTO at Syncbase', initial: 'JO', quote: 'We evaluated five uptime monitors. Self-hosting this was the right call for our compliance requirements, and the setup took under an hour.' },
  { name: 'Priya Sharma', role: 'Infrastructure Lead', initial: 'PS', quote: 'The status page feature alone was worth the switch. Our users finally have transparent visibility into our system health.' },
  { name: 'Liam Gallagher', role: 'Platform Engineer', initial: 'LG', quote: 'Clean, fast, and no surprise billing. We monitor over 200 endpoints across three regions without breaking a sweat.' },
  { name: 'Sofia Torres', role: 'Principal Engineer', initial: 'ST', quote: 'The alerting pipeline is solid. Slack, Discord, email, and custom webhooks — it covers every channel our team uses.' },
  { name: 'Kenji Nakamura', role: 'VP of Engineering', initial: 'KN', quote: 'Having full control over our monitoring data was the deciding factor. Ping uptime delivered exactly what we needed.' },
]

const current = ref(0)
const total = Math.ceil(testimonials.length / 2)

function next() {
  current.value = (current.value + 1) % total
}

function prev() {
  current.value = (current.value - 1 + total) % total
}
</script>

<template>
  <section id="testimonials" class="py-20 md:py-28 bg-muted/30">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-2xl text-center mb-16">
        <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">Trusted by engineering teams</h2>
        <p class="mt-3 text-lg text-muted-foreground">From startups to enterprises</p>
      </div>

      <div class="relative">
        <div class="grid gap-6 md:grid-cols-2">
          <Card
            v-for="(item, index) in testimonials.slice(current * 2, current * 2 + 2)"
            :key="index"
            class="border-border/50 bg-card/60 dark:bg-card/40"
          >
            <CardContent class="p-6">
              <p class="text-sm text-muted-foreground leading-relaxed mb-6">{{ item.quote }}</p>
              <div class="flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10 text-primary text-sm font-bold">
                  {{ item.initial }}
                </div>
                <div>
                  <p class="text-sm font-semibold">{{ item.name }}</p>
                  <p class="text-xs text-muted-foreground">{{ item.role }}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <div class="flex items-center justify-center gap-4 mt-8">
          <button
            @click="prev"
            class="flex h-10 w-10 items-center justify-center rounded-full border border-border bg-background hover:bg-accent transition-colors"
          >
            <ChevronLeft class="h-4 w-4" />
          </button>
          <div class="flex gap-2">
            <span
              v-for="i in total"
              :key="i"
              :class="[
                'h-2 rounded-full transition-all',
                i - 1 === current ? 'w-6 bg-primary' : 'w-2 bg-muted-foreground/30',
              ]"
            />
          </div>
          <button
            @click="next"
            class="flex h-10 w-10 items-center justify-center rounded-full border border-border bg-background hover:bg-accent transition-colors"
          >
            <ChevronRight class="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>
  </section>
</template>
