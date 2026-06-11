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
  <section id="contact" class="py-20 md:py-28 bg-muted/30">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <!-- Section header -->
      <div class="mx-auto max-w-2xl text-center mb-16 reveal" ref="reveal">
        <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">Get in touch</h2>
        <p class="mt-3 text-lg text-muted-foreground">Questions? Ideas? We're listening</p>
        <p class="mt-2 text-sm text-muted-foreground/80 max-w-md mx-auto">
          Whether you need help with setup, want to contribute, or just want to say hi — drop us a message.
        </p>
      </div>

      <div class="grid gap-8 lg:grid-cols-3 max-w-5xl mx-auto">
        <!-- Form -->
        <div class="lg:col-span-2 reveal" ref="reveal">
          <Card class="border-border/50 bg-card/60 dark:bg-card/40">
            <CardContent class="p-6 space-y-4">
              <div class="grid gap-4 sm:grid-cols-2">
                <div>
                  <label class="text-xs font-medium text-foreground mb-1.5 block">First Name</label>
                  <Input placeholder="John" class="rounded-lg" />
                </div>
                <div>
                  <label class="text-xs font-medium text-foreground mb-1.5 block">Last Name</label>
                  <Input placeholder="Doe" class="rounded-lg" />
                </div>
              </div>
              <div>
                <label class="text-xs font-medium text-foreground mb-1.5 block">Email</label>
                <Input type="email" placeholder="john@example.com" class="rounded-lg" />
              </div>
              <div>
                <label class="text-xs font-medium text-foreground mb-1.5 block">Subject</label>
                <Input placeholder="How can we help?" class="rounded-lg" />
              </div>
              <div>
                <label class="text-xs font-medium text-foreground mb-1.5 block">Message</label>
                <Textarea placeholder="Tell us more..." class="min-h-30 rounded-lg" />
              </div>
              <Button class="w-full rounded-lg font-semibold">Send message</Button>
            </CardContent>
          </Card>
        </div>

        <!-- Contact info -->
        <div class="space-y-4 reveal reveal-delay-2" ref="reveal">
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
