<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { AccordionRoot, type AccordionRootEmits, type AccordionRootProps } from 'reka-ui'
import { cn } from '@/lib/utils'

interface Props extends AccordionRootProps {
  class?: HTMLAttributes['class']
}

const props = defineProps<Props>()
const emits = defineEmits<AccordionRootEmits>()
</script>

<template>
  <AccordionRoot
    v-bind="props"
    :class="cn('w-full', props.class)"
    @update:model-value="(val: string | string[] | undefined) => emits('update:modelValue', val)"
  >
    <slot />
  </AccordionRoot>
</template>
