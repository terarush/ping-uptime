<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { MessageCircle, ArrowRight } from '@lucide/vue'
</script>

<template>
  <section class="py-20 md:py-28 bg-muted/30">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-primary/20 via-primary/10 to-background border border-primary/20 p-8 md:p-16 text-center">
        <div class="absolute top-[-20%] right-[-10%] w-64 h-64 rounded-full bg-primary/20 blur-[100px] pointer-events-none"></div>
        <div class="absolute bottom-[-20%] left-[-10%] w-64 h-64 rounded-full bg-primary/20 blur-[100px] pointer-events-none"></div>

        <div class="relative">
          <div class="inline-flex items-center justify-center h-16 w-16 rounded-full bg-primary/10 mb-6">
            <MessageCircle class="h-8 w-8 text-primary" />
          </div>
          <h2 class="text-3xl font-bold tracking-tight sm:text-4xl mb-4">
            Ready to join this <span class="text-primary">Community</span>?
          </h2>
          <p class="text-muted-foreground max-w-lg mx-auto mb-8">
            Join our vibrant Discord community! Connect, share, and grow with like-minded enthusiasts. Click to dive in!
          </p>
          <Button size="lg" class="rounded-full px-8 font-semibold">
            Join Discord
            <ArrowRight class="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  </section>
</template>
