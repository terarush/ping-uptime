<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { ArrowRight } from '@lucide/vue'
import { useScrollReveal } from '@/lib/useScrollReveal'

const { observe } = useScrollReveal()
</script>

<template>
  <section class="relative overflow-hidden py-24 md:py-32">
    <!-- Ambient glow -->
    <div class="absolute left-[-8%] top-1/2 -translate-y-1/2 w-[30%] h-[50%] rounded-full bg-primary/5 blur-[120px] pointer-events-none" />
    <div class="absolute right-[-8%] bottom-0 w-[25%] h-[40%] rounded-full bg-primary/4 blur-[100px] pointer-events-none" />

    <div class="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-3xl text-center reveal" :ref="(el: any) => observe(el)">
        <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl text-foreground">
          Start monitoring today
        </h2>
        <p class="mt-4 text-base text-muted-foreground max-w-lg mx-auto leading-relaxed">
          Self-hosted uptime tracking with alerts, analytics, and a clean status dashboard.
        </p>
        <div class="mt-8 flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button
            size="lg"
            class="rounded-lg px-8 font-semibold text-base"
            as="a"
            href="https://demo-ping.terarush.dev"
            target="_blank"
          >
            Live Demo
            <ArrowRight class="h-4 w-4 ml-0.5" />
          </Button>
          <Button
            variant="outline"
            size="lg"
            class="rounded-lg px-8 font-semibold text-base"
            as="a"
            href="https://demo-ping.terarush.dev/api/docs/"
            target="_blank"
          >
            API Docs
          </Button>
        </div>
      </div>
    </div>
  </section>
</template>
