<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Button } from '@/components/ui/button'
import { ArrowRight, Activity } from '@lucide/vue'

const version = ref('...')
const isDark = ref(false)
let observer: MutationObserver | null = null

async function fetchVersion() {
  try {
    const res = await fetch('https://api.github.com/repos/terarush/ping-uptime/releases/latest')
    const data = await res.json()
    version.value = data.tag_name ?? 'latest'
  } catch {
    version.value = 'latest'
  }
}

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

function onMouseEnter() { stopInterval() }
function onMouseLeave() { startInterval() }

onMounted(() => {
  fetchVersion()
  isDark.value = document.documentElement.classList.contains('dark')
  observer = new MutationObserver(() => {
    isDark.value = document.documentElement.classList.contains('dark')
  })
  observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] })
  startInterval()
})

onUnmounted(() => {
  stopInterval()
  observer?.disconnect()
})
</script>

<template>
  <section id="hero" class="relative overflow-hidden pt-32 pb-20 md:pt-40 md:pb-28">
    <!-- Ambient background -->
    <div class="absolute inset-0 bg-linear-to-b from-primary/5 via-transparent to-transparent pointer-events-none"></div>
    <div class="absolute top-[-10%] left-[-10%] w-[50%] h-[50%] rounded-full bg-primary/10 blur-[120px] pointer-events-none animate-pulse"></div>
    <div class="absolute bottom-[-10%] right-[-10%] w-[50%] h-[50%] rounded-full bg-primary/10 blur-[120px] pointer-events-none"></div>
    <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[70%] h-[70%] rounded-full bg-primary/5 blur-[100px] pointer-events-none"></div>

    <div class="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto max-w-3xl text-center">
        <!-- Badge -->
        <div
          class="reveal visible mb-6 inline-flex items-center gap-2 rounded-full border border-border bg-background/50 px-4 py-1.5 text-xs font-medium text-muted-foreground backdrop-blur-sm"
        >
          <Activity class="h-3.5 w-3.5 text-primary" />
          <span>{{ version }} — Now open source</span>
        </div>

        <!-- Headline -->
        <h1 class="reveal visible font-display text-4xl font-bold tracking-tight sm:text-5xl md:text-6xl lg:text-7xl leading-[1.1]">
          <span class="text-foreground">Monitor uptime.</span>
          <br />
          <span class="bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">Ship with confidence.</span>
        </h1>

        <!-- Subtext -->
        <p class="reveal visible mt-6 text-base text-muted-foreground sm:text-lg max-w-2xl mx-auto leading-relaxed">
          Self-hosted, open-source uptime monitoring built with <span class="font-semibold text-foreground">Go &amp; Vue</span>.
          Track your websites, APIs, and services with real-time alerts, detailed analytics, and a clean status dashboard.
        </p>

        <!-- CTAs -->
        <div class="reveal visible mt-8 flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button size="lg" class="rounded-lg px-8 font-semibold text-base">
            Get Started
            <ArrowRight class="h-4 w-4 ml-0.5" />
          </Button>
          <Button
            variant="outline"
            size="lg"
            class="rounded-lg px-8 font-semibold text-base"
            as="a"
            href="https://github.com/terarush/ping-uptime"
            target="_blank"
          >
            GitHub
          </Button>
        </div>

        <!-- Live indicator -->
        <div class="reveal visible mt-6 flex items-center justify-center gap-4 text-xs text-muted-foreground">
          <span class="inline-flex items-center gap-1.5">
            <span class="relative flex h-2 w-2">
              <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary opacity-75"></span>
              <span class="relative inline-flex h-2 w-2 rounded-full bg-primary"></span>
            </span>
            Live monitoring active
          </span>
          <span class="inline-flex items-center gap-1.5">
            <span class="h-2 w-2 rounded-full bg-primary/40"></span>
            100% uptime track record
          </span>
        </div>
      </div>

      <!-- Carousel -->
      <div class="mt-16 mx-auto max-w-5xl reveal visible">
        <div
          class="group relative rounded-xl border border-border/50 bg-linear-to-b from-background to-muted/50 shadow-2xl overflow-hidden"
          @mouseenter="onMouseEnter"
          @mouseleave="onMouseLeave"
        >
          <div class="aspect-video relative overflow-hidden">
            <div
              class="flex h-full transition-transform duration-500 ease-in-out"
              :style="{ transform: `translateX(-${currentIndex * 100}%)` }"
            >
              <div v-for="(img, i) in images" :key="i" class="min-w-full h-full">
                <img :src="img" class="w-full h-full object-cover" draggable="false" />
              </div>
            </div>

            <!-- Invisible click zones — left/right edges -->
            <div
              @click="prev"
              class="absolute left-0 top-0 h-full w-1/3 cursor-pointer z-10"
            ></div>
            <div
              @click="next"
              class="absolute right-0 top-0 h-full w-1/3 cursor-pointer z-10"
            ></div>

            <div class="absolute bottom-3 left-1/2 -translate-x-1/2 flex items-center gap-2">
              <button
                v-for="(_, i) in images"
                :key="i"
                @click="goTo(i)"
                class="h-2 w-2 rounded-full transition-all duration-300"
                :class="i === currentIndex ? 'bg-primary w-4' : 'bg-muted-foreground/30 hover:bg-muted-foreground/50'"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
