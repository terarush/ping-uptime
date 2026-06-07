<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'

const benefits = [
  { num: '01', title: 'Build Brand Trust', desc: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. A odio velit cum aliquam. Natus consectetur dolores.' },
  { num: '02', title: 'More Leads', desc: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. A odio velit cum aliquam, natus consectetur.' },
  { num: '03', title: 'Higher Conversions', desc: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. Natus consectetur. A odio velit cum aliquam' },
  { num: '04', title: 'Test Marketing Ideas', desc: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. A odio velit cum aliquam. Natus consectetur dolores.' },
]
</script>

<template>
  <section id="features" class="py-20 md:py-28">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-2xl text-center mb-16">
        <h2 class="text-3xl font-bold tracking-tight sm:text-4xl">Benefits</h2>
        <p class="mt-3 text-lg text-muted-foreground">Your Shortcut to Success</p>
        <p class="mt-2 text-sm text-muted-foreground/80 max-w-lg mx-auto">
          Lorem ipsum dolor sit amet consectetur, adipisicing elit. Non ducimus reprehenderit
          architecto rerum similique facere odit deleniti necessitatibus quo quae.
        </p>
      </div>

      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <Card
          v-for="benefit in benefits"
          :key="benefit.num"
          class="relative overflow-hidden border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md hover:shadow-md transition-shadow group"
        >
          <CardHeader>
            <span class="text-4xl font-black text-primary/20 group-hover:text-primary/40 transition-colors">{{ benefit.num }}</span>
            <CardTitle class="text-lg mt-2">{{ benefit.title }}</CardTitle>
          </CardHeader>
          <CardContent>
            <CardDescription class="text-sm">{{ benefit.desc }}</CardDescription>
          </CardContent>
        </Card>
      </div>
    </div>
  </section>
</template>
