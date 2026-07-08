<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Button } from '@/components/ui/button'
import { ArrowRight } from '@lucide/vue'

const isDark = ref(false)
let observer: MutationObserver | null = null
let reducedMotion = false

const lightImages = ['/dashboard_preview_light.png', '/analycis_preview_light.png']
const darkImages = ['/dashboard_preview_dark.png', '/analycis_preview_dark.png']

const images = computed(() => (isDark.value ? darkImages : lightImages))
const currentIndex = ref(0)
let intervalId: ReturnType<typeof setInterval> | null = null

function startInterval() {
  stopInterval()
  intervalId = setInterval(() => {
    currentIndex.value = (currentIndex.value + 1) % images.value.length
  }, 5000)
}

function stopInterval() {
  if (intervalId !== null) {
    clearInterval(intervalId)
    intervalId = null
  }
}

function prev() {
  currentIndex.value = (currentIndex.value - 1 + images.value.length) % images.value.length
  startInterval()
}

function next() {
  currentIndex.value = (currentIndex.value + 1) % images.value.length
  startInterval()
}

function goTo(i: number) {
  currentIndex.value = i
  startInterval()
}

onMounted(() => {
  reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  isDark.value = document.documentElement.classList.contains('dark')
  observer = new MutationObserver(() => {
    isDark.value = document.documentElement.classList.contains('dark')
  })
  observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] })
  if (!reducedMotion) startInterval()
})

onUnmounted(() => {
  stopInterval()
  observer?.disconnect()
})
</script>

<template>
  <section id="hero" class="relative overflow-hidden pt-24 pb-16 md:pt-24 md:pb-24">
    <!-- Ambient glow -->
    <div class="absolute inset-0 bg-linear-to-b from-primary/[0.04] via-transparent to-transparent pointer-events-none" />
    <div class="absolute top-[-12%] left-[15%] w-[40%] h-[40%] rounded-full bg-primary/8 blur-[140px] pointer-events-none" />

    <div class="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">

      <div class="mx-auto max-w-3xl text-center">
        <!-- Headline -->
        <h1 class="font-display text-4xl font-bold tracking-tight sm:text-5xl md:text-6xl lg:text-7xl leading-[1.08] animate-fade-in-up">
          <span class="text-foreground">Monitor uptime.</span>
          <br />
          <span class="bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">Ship with confidence.</span>
        </h1>

        <!-- Subtext -->
        <p class="mt-6 text-base text-muted-foreground sm:text-lg max-w-2xl mx-auto leading-relaxed motion-safe:animate-fade-in-up motion-safe:[animation-delay:150ms] motion-safe:opacity-0 motion-safe:[animation-fill-mode:forwards]">
          Self-hosted, open-source uptime monitoring built with Go and Vue.
          Track websites, APIs, and services with real-time alerts, detailed analytics, and a clean status dashboard.
        </p>

        <!-- CTAs -->
        <div class="mt-8 flex flex-col sm:flex-row items-center justify-center gap-4 motion-safe:animate-fade-in-up motion-safe:[animation-delay:250ms] motion-safe:opacity-0 motion-safe:[animation-fill-mode:forwards]">
          <Button
            size="lg"
            class="rounded-lg px-8 font-semibold text-base active:scale-[0.97] transition-transform motion-reduce:active:scale-100"
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
            class="rounded-lg px-8 font-semibold text-base active:scale-[0.97] transition-transform motion-reduce:active:scale-100"
            as="a"
            href="https://demo-ping.terarush.dev/api/docs/"
            target="_blank"
          >
            API Docs
          </Button>
          <Button
            variant="outline"
            size="lg"
            class="rounded-lg px-8 font-semibold text-base active:scale-[0.97] transition-transform motion-reduce:active:scale-100"
            as="a"
            href="https://github.com/terarush/ping-uptime"
            target="_blank"
          >
            View on GitHub
          </Button>
        </div>
      </div>

      <!-- Screenshot carousel -->
      <div class="mt-16 mx-auto max-w-5xl motion-safe:animate-fade-in motion-safe:[animation-delay:400ms] motion-safe:opacity-0 motion-safe:[animation-fill-mode:forwards]">
        <div
          class="group relative rounded-xl border border-border/50 bg-card shadow-2xl overflow-hidden"
        >
          <div class="aspect-video relative overflow-hidden">
            <div
              class="flex h-full transition-transform duration-500 ease-out"
              :style="{ transform: `translateX(-${currentIndex * 100}%)` }"
            >
              <div v-for="(img, i) in images" :key="i" class="min-w-full h-full">
                <img :src="img" class="w-full h-full object-cover" draggable="false" />
              </div>
            </div>

            <!-- Click zones -->
            <div
              @click="prev"
              class="absolute left-0 top-0 h-full w-1/3 cursor-pointer z-10"
            />
            <div
              @click="next"
              class="absolute right-0 top-0 h-full w-1/3 cursor-pointer z-10"
            />

            <div class="absolute bottom-3 left-1/2 -translate-x-1/2 flex items-center gap-2">
              <button
                v-for="(_, i) in images"
                :key="i"
                @click="goTo(i)"
                class="h-2 rounded-full transition-all duration-300"
                :class="i === currentIndex ? 'w-5 bg-primary' : 'w-2 bg-muted-foreground/30 hover:bg-muted-foreground/50'"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
