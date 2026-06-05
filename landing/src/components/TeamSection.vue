<script setup lang="ts">
import { Card, CardContent } from '@/components/ui/card'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Globe, GitFork } from '@lucide/vue'

const team = [
  { name: 'LeoMiranda', role: 'Vue Frontend Developer, Creator Of This Website', img: 'https://i.pravatar.cc/250?img=58', social: ['linkedin', 'github', 'twitter'] },
  { name: 'Elizabeth Moore', role: 'UI/UX Designer', img: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=1528', social: ['linkedin', 'twitter'] },
  { name: 'David Diaz', role: 'Machine Learning Engineer, TensorFlow Tinkerer', img: 'https://images.unsplash.com/photo-1527980965255-d3b416303d12?w=1760', social: ['linkedin', 'github'] },
  { name: 'Sarah Robinson', role: 'Cloud Native Developer, Kubernetes Orchestrator', img: 'https://images.unsplash.com/photo-1573497161161-c3e73707e25c?w=1587', social: ['linkedin', 'github', 'twitter'] },
  { name: 'Michael Holland', role: 'DevOps Engineer, CI/CD Pipeline Mastermind', img: 'https://images.unsplash.com/photo-1616805765352-beedbad46b2a?w=1887', social: ['linkedin'] },
  { name: 'Zoe Garcia', role: 'JavaScript Evangelist, Deno Champion', img: 'https://images.unsplash.com/photo-1573497019940-1c28c88b4f3e?w=1587', social: ['linkedin', 'github'] },
  { name: 'Evan James', role: 'Backend Developer', img: 'https://images.unsplash.com/photo-1633332755192-727a05c4013d?w=1480', social: ['linkedin', 'github', 'twitter'] },
  { name: 'Pam Taylor', role: 'Fullstack Developer, UX Researcher', img: 'https://images.unsplash.com/photo-1573497019236-17f8177b81e8?w=1740', social: ['twitter'] },
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
          class="border-border/50 bg-background/50 backdrop-blur-sm text-center hover:shadow-lg transition-all duration-300"
        >
          <CardContent class="p-6">
            <Avatar class="mx-auto h-20 w-20 mb-4">
              <AvatarImage :src="member.img" :alt="member.name" />
              <AvatarFallback>{{ member.name.split(' ').map(n => n[0]).join('') }}</AvatarFallback>
            </Avatar>
            <h3 class="font-semibold text-sm">{{ member.name }}</h3>
            <p class="text-xs text-muted-foreground mt-1 leading-relaxed">{{ member.role }}</p>
            <div class="flex items-center justify-center gap-2 mt-4">
              <a
                v-for="social in member.social"
                :key="social"
                href="#"
                class="flex h-8 w-8 items-center justify-center rounded-full border border-border hover:bg-accent hover:text-primary transition-colors"
              >
                <GitFork v-if="social === 'github'" class="h-3.5 w-3.5" />
                <Globe v-else class="h-3.5 w-3.5" />
              </a>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </section>
</template>
