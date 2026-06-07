<script setup lang="ts">
import { Card, CardContent } from '@/components/ui/card'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'

const team = [
  {
    name: "Ahmad Rafi'i",
    role: 'CEO & Founder',
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
    role: 'Junior',
    img: 'https://github.com/szmaou.png',
    linkedin: '#',
    github: '#',
  },
]
</script>

<template>
  <section id="team" class="py-20 md:py-28">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-2xl text-center mb-16">
        <h2 class="text-3xl font-bold tracking-tight sm:text-4xl">Team</h2>
        <p class="mt-3 text-lg text-muted-foreground">The Company Dream Team</p>
      </div>

      <div class="grid gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
        <Card
          v-for="member in team"
          :key="member.name"
          class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md text-center hover:shadow-md transition-shadow"
        >
          <CardContent class="p-6">
            <Avatar class="mx-auto h-20 w-20 mb-4">
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
                :href="member.linkedin"
                target="_blank"
                rel="noopener noreferrer"
                class="flex h-8 w-8 items-center justify-center rounded-full border border-border hover:bg-accent hover:text-primary transition-colors"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  class="h-3.5 w-3.5"
                >
                  <path
                    d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"
                  />
                  <rect x="2" y="9" width="4" height="12" />
                  <circle cx="4" cy="4" r="2" />
                </svg>
              </a>
              <a
                :href="member.github"
                target="_blank"
                rel="noopener noreferrer"
                class="flex h-8 w-8 items-center justify-center rounded-full border border-border hover:bg-accent hover:text-primary transition-colors"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  class="h-3.5 w-3.5"
                >
                  <path
                    d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"
                  />
                </svg>
              </a>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </section>
</template>
