<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Smartphone, Users, Crosshair, Image, MousePointerClick, Type } from '@lucide/vue'

const features = [
  { icon: Smartphone, title: 'Mobile Friendly', desc: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. A odio velit cum aliquam, consectetur.' },
  { icon: Users, title: 'Social Proof', desc: 'Lorem ipsum dolor sit amet consectetur. Natus consectetur, odio ea accusamus aperiam.' },
  { icon: Crosshair, title: 'Targeted Content', desc: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. odio ea accusamus aperiam.' },
  { icon: Image, title: 'Strong Visuals', desc: 'Lorem elit. A odio velit cum aliquam. Natus consectetur dolores, odio ea accusamus aperiam.' },
  { icon: MousePointerClick, title: 'Clear CTA', desc: 'Lorem ipsum dolor sit amet consectetur adipisicing. odio ea accusamus consectetur.' },
  { icon: Type, title: 'Clear Headline', desc: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. A odio velit cum aliquam. Natus consectetur.' },
]
</script>

<template>
  <section class="py-20 md:py-28 bg-muted/30">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-2xl text-center mb-16">
        <h2 class="text-3xl font-bold tracking-tight sm:text-4xl">Features</h2>
        <p class="mt-3 text-lg text-muted-foreground">What Makes Us Different</p>
        <p class="mt-2 text-sm text-muted-foreground/80">
          Lorem ipsum dolor, sit amet consectetur adipisicing elit. Voluptatem fugiat,
          odit similique quasi sint reiciendis quidem iure veritatis optio facere tenetur.
        </p>
      </div>

      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <Card
          v-for="feature in features"
          :key="feature.title"
          class="border-border/50 bg-background hover:shadow-lg hover:border-primary/20 transition-all duration-300 group"
        >
          <CardHeader>
            <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary mb-2 group-hover:bg-primary/20 transition-colors">
              <component :is="feature.icon" class="h-5 w-5" />
            </div>
            <CardTitle class="text-base">{{ feature.title }}</CardTitle>
          </CardHeader>
          <CardContent>
            <CardDescription class="text-sm">{{ feature.desc }}</CardDescription>
          </CardContent>
        </Card>
      </div>
    </div>
  </section>
</template>
