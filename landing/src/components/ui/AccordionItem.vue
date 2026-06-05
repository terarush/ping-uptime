<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { AccordionItem, type AccordionItemProps } from 'reka-ui'
import { cn } from '@/lib/utils'

interface Props extends AccordionItemProps {
  class?: HTMLAttributes['class']
}

const props = defineProps<Props>()
</script>

<template>
  <AccordionItem v-bind="props" :class="cn('border-b', props.class)">
    <slot />
  </AccordionItem>
</template>
