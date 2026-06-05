<script setup lang="ts">
const sponsors = [
  'Acmebrand', 'Acmelogo', 'Acmesponsor', 'Acmeipsum',
  'Acme', 'Accmee', 'Acmetech', 'Acmebrand',
]
</script>

<template>
  <section class="py-16 border-y border-border/50">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <p class="text-center text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-8">
        Our Platinum Sponsors
      </p>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-8 items-center justify-items-center">
        <div
          v-for="sponsor in sponsors"
          :key="sponsor"
          class="flex items-center justify-center h-12 px-6 rounded-lg bg-muted/50 border border-border/30 text-sm font-semibold text-muted-foreground/60"
        >
          {{ sponsor }}
        </div>
      </div>
    </div>
  </section>
</template>
