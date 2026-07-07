<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { useScrollReveal } from '@/lib/useScrollReveal'
import { MapPin, Phone, Mail } from '@lucide/vue'

const { observe } = useScrollReveal()

const contactInfo = [
  { icon: MapPin, label: 'Location', value: 'Malang, East Java\nIndonesia' },
  { icon: Phone, label: 'Phone', value: '+62 821 4333 8737' },
  { icon: Mail, label: 'Email', value: 'hello@terarush.dev' },
]
</script>

<template>
  <section id="contact" class="section-grid py-20 md:py-28 bg-muted/30">
    <div class="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-2xl text-center mb-16 reveal" :ref="(el: any) => observe(el)">
        <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">Get in touch</h2>
        <p class="mt-3 text-lg text-muted-foreground">Questions, ideas, or just want to say hi</p>
      </div>

      <div class="grid gap-8 lg:grid-cols-3 max-w-5xl mx-auto">
        <!-- Contact card -->
        <div class="lg:col-span-2 reveal" :ref="(el: any) => observe(el)">
          <Card class="border-border/50 bg-card/60 dark:bg-card/40">
            <CardContent class="p-8 text-center">
              <div class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-xl bg-primary/10 text-primary">
                <Mail class="h-7 w-7" />
              </div>
              <h3 class="font-display text-2xl font-bold mb-2">hello@terarush.dev</h3>
              <p class="text-sm text-muted-foreground mb-6 max-w-md mx-auto leading-relaxed">
                Send us an email. We typically reply within 24 hours.
              </p>
              <Button
                size="lg"
                class="rounded-lg px-8 font-semibold text-base"
                as="a"
                href="mailto:hello@terarush.dev"
              >
                Send Email
              </Button>
            </CardContent>
          </Card>
        </div>

        <!-- Contact info sidebar -->
        <div class="space-y-4 reveal reveal-delay-2" :ref="(el: any) => observe(el)">
          <div
            v-for="item in contactInfo"
            :key="item.label"
            class="flex gap-3 p-4 rounded-lg border border-border/50 bg-card/60 dark:bg-card/40"
          >
            <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <component :is="item.icon" class="h-5 w-5" />
            </div>
            <div>
              <p class="text-xs font-semibold text-muted-foreground">{{ item.label }}</p>
              <p class="text-sm text-foreground whitespace-pre-line">{{ item.value }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
