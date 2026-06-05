<script setup lang="ts">
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Globe, Share2, Mail, Search } from '@lucide/vue'

const services = [
  { icon: Globe, title: 'Custom Domain Integration', desc: 'Lorem ipsum dolor sit, amet consectetur adipisicing elit adipisicing.', badge: 'PRO' },
  { icon: Share2, title: 'Social Media Integrations', desc: 'Lorem ipsum dolor sit, amet consectetur adipisicing elit. Molestiae, dicta.', badge: 'PRO' },
  { icon: Mail, title: 'Email Marketing Integrations', desc: 'Lorem dolor sit amet adipisicing.', badge: 'PRO' },
  { icon: Search, title: 'SEO Optimization', desc: 'Lorem ipsum dolor sit amet consectetur.', badge: 'PRO' },
]
</script>

<template>
  <section class="py-20 md:py-28">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-2xl text-center mb-16">
        <h2 class="text-3xl font-bold tracking-tight sm:text-4xl">Services</h2>
        <p class="mt-3 text-lg text-muted-foreground">Grow Your Business</p>
        <p class="mt-2 text-sm text-muted-foreground/80">
          From marketing and sales to operations and strategy, we have the expertise to help you achieve your goals.
        </p>
      </div>

      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <Card
          v-for="service in services"
          :key="service.title"
          class="relative border-border/50 bg-background/50 backdrop-blur-sm hover:shadow-lg transition-all duration-300 group"
        >
          <CardHeader>
            <Badge variant="secondary" class="absolute top-4 right-4 text-[10px] px-2 py-0.5">
              {{ service.badge }}
            </Badge>
            <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary mb-2">
              <component :is="service.icon" class="h-5 w-5" />
            </div>
            <CardTitle class="text-base">{{ service.title }}</CardTitle>
          </CardHeader>
          <CardContent>
            <CardDescription class="text-sm">{{ service.desc }}</CardDescription>
          </CardContent>
        </Card>
      </div>
    </div>
  </section>
</template>
