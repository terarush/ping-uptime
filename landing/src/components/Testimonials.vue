<script setup lang="ts">
import { Card, CardContent } from '@/components/ui/card'
import { ref } from 'vue'
import { ChevronLeft, ChevronRight, Star } from '@lucide/vue'

const testimonials = [
  { name: 'John Doe', role: 'Product Manager', initial: 'JD', quote: 'Wow NextJs + Shadcn is awesome!. This template lets me change colors, fonts and images to match my brand identity.' },
  { name: 'Sophia Collins', role: 'Cybersecurity Analyst', initial: 'SC', quote: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna.' },
  { name: 'Adam Johnson', role: 'Chief Technology Officer', initial: 'AJ', quote: 'Lorem ipsum dolor sit amet,exercitation. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.' },
  { name: 'Ethan Parker', role: 'Data Scientist', initial: 'EP', quote: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod labore et dolore magna aliqua. Ut enim ad minim veniam.' },
  { name: 'Ava Mitchell', role: 'IT Project Manager', initial: 'AM', quote: 'Lorem ipsum dolor sit amet, tempor incididunt aliqua. Ut enim ad minim veniam, quis nostrud incididunt consectetur adipiscing elit.' },
  { name: 'Isabella Reed', role: 'DevOps Engineer', initial: 'IR', quote: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.' },
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
        <h2 class="text-3xl font-bold tracking-tight sm:text-4xl">Testimonials</h2>
        <p class="mt-3 text-lg text-muted-foreground">Hear What Our 1000+ Clients Say</p>
      </div>

      <div class="relative">
        <div class="grid gap-6 md:grid-cols-2">
          <Card
            v-for="(item, index) in testimonials.slice(current * 2, current * 2 + 2)"
            :key="index"
            class="border-border/50 bg-card/60 dark:bg-card/40"
          >
            <CardContent class="p-6">
              <div class="flex gap-1 mb-4">
                <Star v-for="i in 5" :key="i" class="h-4 w-4 fill-primary text-primary" />
              </div>
              <p class="text-sm text-muted-foreground leading-relaxed mb-6">&ldquo;{{ item.quote }}&rdquo;</p>
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
                'h-2 w-2 rounded-full transition-all',
                i - 1 === current ? 'bg-primary w-6' : 'bg-muted-foreground/30',
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
