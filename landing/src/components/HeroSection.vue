<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Button } from '@/components/ui/button'
import { ArrowRight, Sparkles, ChevronLeft, ChevronRight } from '@lucide/vue'

const isDark = ref(false)
let observer: MutationObserver | null = null

const lightImages = ['/dashboard_preview_light.png', '/analycis_preview_light.png']
const darkImages = ['/dashboard_preview_dark.png', '/analycis_preview_dark.png']

const images = computed(() => (isDark.value ? darkImages : lightImages))
const currentIndex = ref(0)
let intervalId: ReturnType<typeof setInterval> | null = null

function startInterval() {
  stopInterval()
  intervalId = setInterval(() => {
    currentIndex.value = (currentIndex.value + 1) % images.value.length
  }, 4000)
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

function onMouseEnter() {
  stopInterval()
}

function onMouseLeave() {
  startInterval()
}

onMounted(() => {
  isDark.value = document.documentElement.classList.contains('dark')
  observer = new MutationObserver(() => {
    isDark.value = document.documentElement.classList.contains('dark')
  })
  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['class'],
  })
  startInterval()
})

onUnmounted(() => {
  stopInterval()
  observer?.disconnect()
})
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
                <img
                  :src="img"
                  class="w-full h-full object-cover"
                />
              </div>
            </div>

            <button
              @click="prev()"
              class="absolute left-2 top-1/2 -translate-y-1/2 flex items-center justify-center h-9 w-9 rounded-full bg-background/80 backdrop-blur-sm text-foreground shadow-sm opacity-0 group-hover:opacity-100 transition-opacity hover:bg-background"
            >
              <ChevronLeft class="h-5 w-5" />
            </button>
            <button
              @click="next()"
              class="absolute right-2 top-1/2 -translate-y-1/2 flex items-center justify-center h-9 w-9 rounded-full bg-background/80 backdrop-blur-sm text-foreground shadow-sm opacity-0 group-hover:opacity-100 transition-opacity hover:bg-background"
            >
              <ChevronRight class="h-5 w-5" />
            </button>

            <div class="absolute bottom-3 left-1/2 -translate-x-1/2 flex items-center gap-2">
              <button
                v-for="(_, i) in images"
                :key="i"
                @click="goTo(i)"
                class="h-2 w-2 rounded-full transition-all"
                :class="i === currentIndex
                  ? 'bg-primary w-4'
                  : 'bg-muted-foreground/30 hover:bg-muted-foreground/50'"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
