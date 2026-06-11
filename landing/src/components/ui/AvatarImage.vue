<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { AvatarImage as RekaImage } from 'reka-ui'
import { cn } from '@/lib/utils'

interface Props {
  class?: HTMLAttributes['class']
  src?: string
  alt?: string
}

const props = defineProps<Props>()
</script>

<template>
  <RekaImage
    :src="props.src ?? ''"
    :alt="props.alt ?? ''"
    :class="cn('aspect-square h-full w-full', props.class)"
  />
</template>
