<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { AccordionHeader, AccordionTrigger as RekaTrigger } from 'reka-ui'
import { ChevronDown } from '@lucide/vue'
import { cn } from '@/lib/utils'

interface Props {
  class?: HTMLAttributes['class']
}

const props = defineProps<Props>()
</script>

<template>
  <AccordionHeader class="flex">
    <RekaTrigger
      :class="cn(
        'flex flex-1 items-center justify-between py-4 text-sm font-medium transition-all hover:underline [&[data-state=open]>svg]:rotate-180',
        props.class,
      )"
    >
      <slot />
      <ChevronDown class="h-4 w-4 shrink-0 text-muted-foreground transition-transform duration-200" />
    </RekaTrigger>
  </AccordionHeader>
</template>
