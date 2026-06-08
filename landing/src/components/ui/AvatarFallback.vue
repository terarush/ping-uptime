<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { AvatarFallback as RekaFallback } from 'reka-ui'
import { cn } from '@/lib/utils'

interface Props {
  class?: HTMLAttributes['class']
}

const props = defineProps<Props>()
</script>

<template>
  <RekaFallback
    :class="cn(
      'flex h-full w-full items-center justify-center rounded-full bg-muted',
      props.class,
    )"
  >
    <slot />
  </RekaFallback>
</template>
