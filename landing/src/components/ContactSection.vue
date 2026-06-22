<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent } from '@/components/ui/card'
import { useScrollReveal } from '@/lib/useScrollReveal'
import { MapPin, Phone, Mail, Clock } from '@lucide/vue'

const { observe } = useScrollReveal()

const contactInfo = [
  { icon: MapPin, label: 'Location', value: 'Malang, East Java\nIndonesia' },
  { icon: Phone, label: 'Phone', value: '+62 821 4333 8737\nMon-Fri, 9AM-10PM' },
  { icon: Mail, label: 'Email', value: 'hello@terarush.dev\nWe reply within 24h' },
  { icon: Clock, label: 'Hours', value: 'Monday - Friday\n8AM - 4PM (WIB)' },
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
        <!-- Form -->
        <div class="lg:col-span-2 reveal" :ref="(el: any) => observe(el)">
          <Card class="border-border/50 bg-card/60 dark:bg-card/40">
            <CardContent class="p-6 space-y-4">
              <div class="grid gap-4 sm:grid-cols-2">
                <div>
                  <label class="text-xs font-medium text-foreground mb-1.5 block" for="first-name">First name</label>
                  <Input id="first-name" placeholder="Alex" class="rounded-lg" />
                </div>
                <div>
                  <label class="text-xs font-medium text-foreground mb-1.5 block" for="last-name">Last name</label>
                  <Input id="last-name" placeholder="Rivera" class="rounded-lg" />
                </div>
              </div>
              <div>
                <label class="text-xs font-medium text-foreground mb-1.5 block" for="email">Email</label>
                <Input id="email" type="email" placeholder="alex@example.com" class="rounded-lg" />
              </div>
              <div>
                <label class="text-xs font-medium text-foreground mb-1.5 block" for="subject">Subject</label>
                <Input id="subject" placeholder="How can we help?" class="rounded-lg" />
              </div>
              <div>
                <label class="text-xs font-medium text-foreground mb-1.5 block" for="message">Message</label>
                <Textarea id="message" placeholder="Tell us more..." class="min-h-30 rounded-lg" />
              </div>
              <Button class="w-full rounded-lg font-semibold">Send message</Button>
            </CardContent>
          </Card>
        </div>

        <!-- Contact info -->
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
