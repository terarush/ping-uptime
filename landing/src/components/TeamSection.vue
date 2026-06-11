<script setup lang="ts">
import { Card, CardContent } from '@/components/ui/card'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { useScrollReveal } from '@/lib/useScrollReveal'

const { observe } = useScrollReveal()

const team = [
  {
    name: "Ahmad Rafi'i",
    role: 'Founder & CEO',
    img: 'https://github.com/rafia9005.png',
    linkedin: 'https://linkedin.com/in/rafia9005',
    github: 'https://github.com/rafia9005',
  },
  {
    name: 'Rakha Adi',
    role: 'Co-Founder',
    img: 'https://github.com/Rakhasptro.png',
    linkedin: 'https://linkedin.com/in/rakhaadi',
    github: 'https://github.com/Rakhasptro',
  },
  {
    name: 'Bintang Adi Alvaro',
    role: 'Graphics Designer',
    img: 'https://github.com/kentangaja.png',
    linkedin: '#',
    github: 'https://github.com/kentangaja',
  },
  {
    name: 'Hasan',
    role: 'Junior Backend',
    img: 'https://github.com/szmaou.png',
    linkedin: '#',
    github: 'https://github.com/szmaou',
  },
]
</script>

<template>
  <section id="team" class="py-20 md:py-28 bg-muted/30">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <!-- Section header -->
      <div class="mx-auto max-w-2xl text-center mb-16 reveal" ref="reveal">
        <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">Team</h2>
        <p class="mt-3 text-lg text-muted-foreground">The people behind ping-uptime</p>
      </div>

      <!-- Grid -->
      <div class="grid gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
        <Card
          v-for="(member, i) in team"
          :key="member.name"
          :ref="(el: any) => observe(el?.$el ?? el)"
          class="reveal group border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md text-center hover:shadow-md transition-shadow"
          :class="`reveal-delay-${i + 1}`"
        >
          <template #default>
            <CardContent class="p-6">
              <Avatar class="mx-auto h-20 w-20 mb-4 ring-2 ring-border/50 group-hover:ring-primary/30 transition-all">
                <AvatarImage :src="member.img" :alt="member.name" />
                <AvatarFallback>{{
                  member.name
                    .split(' ')
                    .map((n) => n[0])
                    .join('')
                }}</AvatarFallback>
              </Avatar>
              <h3 class="font-semibold text-sm">{{ member.name }}</h3>
              <p class="text-xs text-muted-foreground mt-1 leading-relaxed">{{ member.role }}</p>
              <div class="flex items-center justify-center gap-2 mt-4">
                <a
                  v-if="member.linkedin !== '#'"
                  :href="member.linkedin"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="flex h-8 w-8 items-center justify-center rounded-full border border-border hover:bg-accent hover:text-primary transition-colors"
                >
                  <img src="/linkedin.svg" class="h-3.5 w-3.5 opacity-70" aria-hidden="true" />
                </a>
                <a
                  :href="member.github"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="flex h-8 w-8 items-center justify-center rounded-full border border-border hover:bg-accent hover:text-primary transition-colors"
                >
                  <img src="/github.svg" class="h-3.5 w-3.5 opacity-70" aria-hidden="true" />
                </a>
              </div>
            </CardContent>
          </template>
        </Card>
      </div>
    </div>
  </section>
</template>
