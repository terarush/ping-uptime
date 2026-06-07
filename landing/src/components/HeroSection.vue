<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { ArrowRight, Sparkles } from '@lucide/vue'
</script>

<template>
  <section id="hero" class="relative overflow-hidden pt-32 pb-20 md:pt-40 md:pb-28">
    <div
      class="absolute inset-0 bg-linear-to-b from-primary/5 via-transparent to-transparent pointer-events-none"
    ></div>
    <div
      class="absolute top-[-10%] left-[-10%] w-[50%] h-[50%] rounded-full bg-primary/10 blur-[120px] pointer-events-none"
    ></div>
    <div
      class="absolute bottom-[-10%] right-[-10%] w-[50%] h-[50%] rounded-full bg-primary/10 blur-[120px] pointer-events-none"
    ></div>

    <div class="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-3xl text-center">
        <div
          class="mb-6 inline-flex items-center gap-2 rounded-full border border-border bg-background/50 px-4 py-1.5 text-xs font-medium text-muted-foreground backdrop-blur-sm"
        >
          <Sparkles class="h-3.5 w-3.5 text-primary" />
          <span>New</span>
          <span class="text-primary">Design is out now!</span>
        </div>

        <h1 class="text-4xl font-bold tracking-tight sm:text-5xl md:text-6xl lg:text-6xl">
          <span class="text-foreground">Uptime Monitoring Service</span>
          <br />
          <!-- <span class="text-primary">ping-uptime</span> -->
        </h1>

        <p
          class="mt-6 text-base text-muted-foreground sm:text-lg max-w-2xl mx-auto leading-relaxed"
        >
          built with Go (Echo v4) and Vue 3. Monitor your websites, APIs, and services
        </p>

        <div class="mt-8 flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button size="lg" class="rounded-lg px-8 font-semibold">
            Get Started
            <ArrowRight class="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            size="lg"
            class="rounded-lg px-8 font-semibold"
            as="a"
            href="https://github.com/terarush/ping-uptime"
            target="_blank"
          >
            Github repository
          </Button>
        </div>
      </div>

      <div class="mt-16 mx-auto max-w-5xl">
        <div
          class="relative rounded-xl border border-border/50 bg-linear-to-b from-background to-muted/50 shadow-2xl overflow-hidden"
        >
          <div
            class="aspect-video bg-linear-to-br from-primary/5 via-background to-primary/5 flex items-center justify-center"
          >
            <div class="text-center p-8">
              <div
                class="inline-flex items-center gap-2 rounded-lg bg-primary/10 px-4 py-2 text-sm text-primary font-medium"
              >
                <div class="relative flex h-2 w-2">
                  <span
                    class="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary"
                  ></span>
                  <span class="relative inline-flex h-2 w-2 rounded-full bg-primary"></span>
                </div>
                Dashboard Preview
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
