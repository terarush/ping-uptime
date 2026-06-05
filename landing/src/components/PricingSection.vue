<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Check } from '@lucide/vue'

const plans = [
  {
    name: 'Free',
    desc: 'Lorem ipsum dolor sit, amet ipsum consectetur adipisicing elit.',
    price: '$0',
    features: ['1 team member', '1 GB storage', 'Upto 2 pages', 'Community support', 'AI assistance'],
  },
  {
    name: 'Premium',
    desc: 'Lorem ipsum dolor sit, amet ipsum consectetur adipisicing elit.',
    price: '$45',
    popular: true,
    features: ['4 team member', '8 GB storage', 'Upto 6 pages', 'Priority support', 'AI assistance'],
  },
  {
    name: 'Enterprise',
    desc: 'Lorem ipsum dolor sit, amet ipsum consectetur adipisicing elit.',
    price: '$120',
    features: ['10 team member', '20 GB storage', 'Upto 10 pages', 'Phone & email support', 'AI assistance'],
  },
]
</script>

<template>
  <section class="py-20 md:py-28">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-2xl text-center mb-16">
        <h2 class="text-3xl font-bold tracking-tight sm:text-4xl">Pricing</h2>
        <p class="mt-3 text-lg text-muted-foreground">Get unlimitted access</p>
        <p class="mt-2 text-sm text-muted-foreground/80">Lorem ipsum dolor sit amet consectetur adipisicing reiciendis.</p>
      </div>

      <div class="grid gap-8 md:grid-cols-3 lg:max-w-5xl lg:mx-auto">
        <Card
          v-for="plan in plans"
          :key="plan.name"
          :class="[
            'relative border-border/50 bg-background/50 backdrop-blur-sm flex flex-col',
            plan.popular ? 'border-primary shadow-lg shadow-primary/10 scale-105' : '',
          ]"
        >
          <CardHeader v-if="plan.popular" class="pb-0">
            <Badge class="w-fit mx-auto text-xs">Most Popular</Badge>
          </CardHeader>
          <CardHeader>
            <CardTitle class="text-xl">{{ plan.name }}</CardTitle>
            <CardDescription class="text-sm">{{ plan.desc }}</CardDescription>
          </CardHeader>
          <CardContent class="flex-1">
            <div class="mb-6">
              <span class="text-4xl font-bold">{{ plan.price }}</span>
              <span class="text-sm text-muted-foreground ml-1">/month</span>
            </div>
            <ul class="space-y-3">
              <li
                v-for="feature in plan.features"
                :key="feature"
                class="flex items-center gap-2 text-sm text-muted-foreground"
              >
                <Check class="h-4 w-4 text-primary shrink-0" />
                {{ feature }}
              </li>
            </ul>
          </CardContent>
          <CardFooter>
            <Button
              :variant="plan.popular ? 'default' : 'outline'"
              class="w-full rounded-full"
            >
              {{ plan.name === 'Free' ? 'Start Free Trial' : plan.name === 'Premium' ? 'Get started' : 'Contact US' }}
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  </section>
</template>
