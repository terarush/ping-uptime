<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent } from '@/components/ui/card'
import { MapPin, Phone, Mail, Clock } from '@lucide/vue'

const contactInfo = [
  { icon: MapPin, label: 'Find us', value: '742 Evergreen Terrace, Springfield, IL 62704' },
  { icon: Phone, label: 'Call us', value: '+1 (619) 123-4567' },
  { icon: Mail, label: 'Mail US', value: 'leomirandadev@gmail.com' },
  { icon: Clock, label: 'Visit us', value: 'Monday - Friday\n8AM - 4PM' },
]
</script>

<template>
  <section id="contact" class="py-20 md:py-28 bg-muted/30">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-2xl text-center mb-16">
        <h2 class="text-3xl font-bold tracking-tight sm:text-4xl">Contact</h2>
        <p class="mt-3 text-lg text-muted-foreground">Connect With Us</p>
        <p class="mt-2 text-sm text-muted-foreground/80">
          Lorem ipsum dolor sit amet consectetur adipisicing elit. Voluptatum ipsam sint enim exercitationem ex autem corrupti quas tenetur
        </p>
      </div>

      <div class="grid gap-8 lg:grid-cols-3 max-w-5xl mx-auto">
        <div class="lg:col-span-2">
          <Card class="border-border/50 bg-background">
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
                <Textarea placeholder="Tell us more..." class="min-h-[120px] rounded-lg" />
              </div>
              <Button class="w-full rounded-full font-semibold">
                Send message
              </Button>
            </CardContent>
          </Card>
        </div>

        <div class="space-y-4">
          <div
            v-for="item in contactInfo"
            :key="item.label"
            class="flex gap-3 p-4 rounded-lg border border-border/50 bg-background"
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
