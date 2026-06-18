<template>
  <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium"
    :style="{ backgroundColor: color + '20', color: color, borderColor: color + '40' }"
  >
    {{ name }}
  </span>
</template>

<script setup lang="ts">
defineProps<{ name: string; color: string }>()
</script>
